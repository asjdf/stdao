package test

import (
	"github.com/asjdf/gorm-cache/config"
	"github.com/asjdf/stdao"
	"testing"
)

func TestStdWithCacheDAO(t *testing.T) {
	UserDAO := stdao.CreateWithCache(&user{})
	err := UserDAO.Init(db, &config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         config.CacheStorageMemory,
		InvalidateWhenUpdate: true,
		CacheTTL:             10000,
		CacheMaxItemCnt:      5000,
		CacheSize:            1000,
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

	findUser1 := &user{}
	findUser1.ID = user1.ID
	result = UserDAO.Find(findUser1)
	if result.Error != nil {
		t.Error(result.Error)
	}
	if user1.Name != findUser1.Name || user1.Age != findUser1.Age {
		t.Error("findUser1 is not equal with user1")
	}

	findUser2 := &user{}
	findUser2.ID = user1.ID
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
		CacheStorage:         config.CacheStorageMemory,
		InvalidateWhenUpdate: true,
		CacheTTL:             10000,
		CacheMaxItemCnt:      5000,
		CacheSize:            1000,
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
