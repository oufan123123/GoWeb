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
	student4 = &Student{Id: 4, Name: "lxy", Age: 26}
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
	//err1 := s.DropTable()
	//err2 := s.CreateTable()
	//_, err3 := s.Insert(student1, student2, student3, student4)
	//if err1 != nil || err2 != nil || err3 != nil {
	//t.Fatal("test TestInit fail")
	//return
	//}
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

func TestSession_Limit(t *testing.T) {
	s := testInitSession(t)
	var stus []Student
	if err := s.Limit(1).Find(&stus); err != nil {
		t.Fatal("failed to find students")
	}
	if len(stus) != 1 {
		t.Fatal("students' len is not 1")
	}
}

func TestSession_Update(t *testing.T) {
	s := testInitSession(t)
	m := map[string]interface{}{
		"Name": "qtt",
		"Age":  16,
	}
	var sl []interface{}

	affect, _ := s.Where("id = 3").Update(append(sl, m)...)
	stu := &Student{}
	s.OrderBy("Age asc").First(stu)

	if affect != 1 || stu.Age != 16 {
		t.Fatal("update test fail")
	}
}

func TestSession_DeleteAndCount(t *testing.T) {
	s := testInitSession(t)
	affect, _ := s.Where("Id = ?", 1).Delete()
	count, _ := s.Count()
	if affect != 1 || count != 3 {
		t.Fatal("test delete and count fail")
	}

}
