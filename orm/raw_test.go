package orm

import (
	//"fmt"
	//"reflect"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type Student struct {
	Id   int    `orm:"primary key"`
	Name string `orm:"not null"`
	Age  int
}

var (
	student1 = &Student{Id: 1, Name: "of", Age: 26}
	student2 = &Student{Id: 2, Name: "dqj", Age: 27}
	student3 = &Student{Id: 3, Name: "mt", Age: 23}
)

func testInitSession(t *testing.T) *Session {
	driver := "mysql"
	source := "root:Root;123456@tcp(127.0.0.1:3306)/mall"
	o, err := NewOrm(driver, source)
	if err != nil {
		t.Fatal("test orm start fail")
		return nil
	}
	s := o.NewSession()
	s.Model(&Student{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(student1, student2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("test TestInit fail")
		//return
	}
	return s
}

func testSession_Insert(t *testing.T) {
	s := testInitSession(t)
	affected, err := s.Insert(student3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSession_Find(t *testing.T) {
	s := testInitSession(t)
	var students []Student
	if err := s.Find(&students); err != nil || len(students) != 2 {
		t.Fatal("failed to query all")
	}
	fmt.Println("sucess")
}
