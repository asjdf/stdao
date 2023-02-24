package test

import (
	"github.com/asjdf/gorm-cache/config"
	"github.com/asjdf/gorm-cache/storage"
	"github.com/asjdf/stdao"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"sync"
	"testing"
	"time"
)

func TestStdWithCacheDAO(t *testing.T) {
	UserDAO := stdao.CreateWithCache(&user{})
	err := UserDAO.Init(db.Debug(), &config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         storage.NewMem(storage.DefaultMemStoreConfig),
		InvalidateWhenUpdate: true,
		CacheTTL:             10000,
		CacheMaxItemCnt:      5000,
		DebugMode:            true,
	})
	if err != nil {
		t.Error(err)
	}

	user1 := &user{
		Name: "Atom",
		Age:  114514,
	}
	result := UserDAO.Save(user1)
	if result.Error != nil || result.RowsAffected != 1 {
		t.Error(result.Error, result.RowsAffected)
	}
	user2 := &user{
		Name: "Atom",
		Age:  114514,
	}
	result = UserDAO.Save(user2)
	if result.Error != nil || result.RowsAffected != 1 {
		t.Error(result.Error, result.RowsAffected)
	}

	findUser1 := &user{}
	findUser1.Name = user1.Name
	result = UserDAO.Find(findUser1)
	if result.Error != nil {
		t.Error(result.Error)
	}
	if user1.Name != findUser1.Name || user1.Age != findUser1.Age {
		t.Error("findUser1 is not equal with user1")
	}

	findUser2 := &user{}
	findUser2.Name = user1.Name
	result = UserDAO.Find(findUser2)
	if result.Error != nil {
		t.Error(result.Error)
	}
	if user1.Name != findUser2.Name || user1.Age != findUser2.Age {
		t.Error("findUser2 is not equal with user1")
	}
}

func BenchmarkStdWithCacheDAO(b *testing.B) {
	UserDAO := stdao.CreateWithCache(&user{})
	err := UserDAO.Init(db, &config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         storage.NewMem(storage.DefaultMemStoreConfig),
		InvalidateWhenUpdate: true,
		CacheTTL:             10000,
		CacheMaxItemCnt:      5000,
	})
	if err != nil {
		b.Error(err)
	}

	user1 := &user{
		Name: "Atom",
		Age:  114514,
	}
	result := UserDAO.Save(user1)
	if result.Error != nil || result.RowsAffected != 1 {
		b.Error(result.Error, result.RowsAffected)
	}
	for i := 0; i < b.N; i++ {
		findUser1 := &user{}
		findUser1.ID = user1.ID
		result = UserDAO.Find(findUser1)
		if result.Error != nil {
			b.Error(result.Error)
		}
		if user1.Name != findUser1.Name || user1.Age != findUser1.Age {
			b.Error("findUser1 is not equal with user1")
		}
	}
}

// 测试检验是否存在缓存穿透
func TestCachePenetration(t *testing.T) {
	UserDAO := stdao.CreateWithCache(&user{})
	err := UserDAO.Init(db.Debug(), &config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         storage.NewMem(storage.DefaultMemStoreConfig),
		InvalidateWhenUpdate: true,
		CacheTTL:             10000,
		CacheMaxItemCnt:      5000,
		DebugMode:            true,
	})
	if err != nil {
		t.Error(err)
	}

	// first query is not from cache
	u := &user{}
	u.ID = 2 // this id is not exist in database
	result := UserDAO.First(u)
	if result.Error != gorm.ErrRecordNotFound {
		t.Error("should return gorm.ErrRecordNotFound")
	}
	if UserDAO.HitCount() != 0 {
		t.Error("should not hit cache")
	}

	// second query is from cache
	result = UserDAO.First(u)
	if result.Error != gorm.ErrRecordNotFound {
		t.Error("should return gorm.ErrRecordNotFound")
	}
	if UserDAO.HitCount() != 1 {
		t.Error("should hit cache")
	}
}

// 测试检验是否存在缓存击穿
func TestCacheHotspotInvalid(t *testing.T) {
	UserDAO := stdao.CreateWithCache(&user{})
	err := UserDAO.Init(db, &config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         storage.NewMem(storage.DefaultMemStoreConfig),
		InvalidateWhenUpdate: true,
		CacheTTL:             10000,
		CacheMaxItemCnt:      5000,
	})
	if err != nil {
		t.Error(err)
	}

	var trigger sync.WaitGroup
	trigger.Add(1)
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			trigger.Wait()
			defer wg.Done()
			u := &user{}
			u.ID = 2 // this id is not exist in database
			result := UserDAO.First(u)
			if result.Error != gorm.ErrRecordNotFound {
				t.Error("should return gorm.ErrRecordNotFound, but got", result.Error)
			}
		}()
	}
	start := time.Now()
	trigger.Done()
	wg.Wait()
	t.Logf("cost %v", time.Since(start))

	assert.Equal(t, uint64(1), UserDAO.MissCount())
	assert.Equal(t, uint64(99), UserDAO.HitCount())
}
