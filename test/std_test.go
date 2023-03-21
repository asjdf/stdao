package test

import (
	"github.com/asjdf/stdao"
	"testing"
)

func TestStdDAO(t *testing.T) {
	UserDAO := stdao.Create(&user{})
	err := UserDAO.Init(db)
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
}

func TestStdFind(t *testing.T) {
	UserDAO := stdao.Create(&user{})
	err := UserDAO.Init(db)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		user1 := &user{
			Name: "Atom",
			Age:  uint(i),
		}
		result := UserDAO.Save(user1)
		if result.Error != nil || result.RowsAffected != 1 {
			t.Error(result.Error, result.RowsAffected)
		}
	}

	findUser1 := &user{}
	findUser1.ID = 5
	result := UserDAO.Find(findUser1)
	if result.Error != nil {
		t.Error(result.Error)
	}
	if 4 != findUser1.Age {
		t.Error("findUser1 is not equal with user1")
	}

	findUser1 = &user{}
	findUser1.Age = 5
	result = UserDAO.Find(findUser1)
	if result.Error != nil {
		t.Error(result.Error)
	}
	if 5 != findUser1.Age {
		t.Error("findUser1 is not equal with user1")
	}
}
