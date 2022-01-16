package main

import (
	"database/sql"
	"fmt"
	"orm"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func OrmTest(t *testing.T) *orm.Session {
	driver := "mysql"
	source := "root:Root;123456@tcp(127.0.0.1:3306)/mall"
	o, err := orm.NewOrm(driver, source)
	if err != nil {
		t.Fatal("test orm start fail")
		return nil
	}
	s := o.NewSession()
	return s
}

func SessionTest(t *testing.T) {
	s := OrmTest(t)
	//sqlCreate := "create table user(user_id int not null)"
	//sqlInsert := "insert into user(user_id, user_name, mobile, password, is_locked, address, is_delete) values(?,?,?,?,?,?,?)"
	sqlSelect := "select user_name from user where user_id=?"

	//sqlUpdate := "update user set password=? where user_id=?"

	sqlSelectAll := "select * from user"

	rows, _ := s.Raw(sqlSelectAll).QueryRows()

	if rows != nil {
		for rows.Next() {
			var user_id, user_name, mobile, password, address string
			var is_locked, is_delete bool
			var createAt sql.NullTime
			var updateAt sql.NullTime // use pointer to receive may nil
			err := rows.Scan(&user_id, &user_name, &mobile, &password, &is_locked, &address, &is_delete, &createAt, &updateAt)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			ct, _ := createAt.Value()
			ut, _ := updateAt.Value()
			fmt.Printf("select user:%s, %s, %s,%s, %s, %s,%s, %v, %v", user_id, user_name, mobile, password, is_locked, address, is_delete, ct, ut)
		}
	}
	s.Clear()

	row := s.Raw(sqlSelect, 2).QueryRow()
	var user_name string
	err := row.Scan(&user_name)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("user_name:%s", user_name)
	}
	s.Clear()

}

type User struct {
	Name     string `orm:"primary key"`
	Age      int    `orm:"not null"`
	Birthday time.Time
}

var mysqlDialect, _ = orm.GetDialect("mysql")

func ParseTest(t *testing.T) {
	schema := orm.Parse(&User{}, mysqlDialect)

	if schema.Name != "User" || len(schema.FieldNames) != 3 {
		t.Fatal("func Parse test fail")
		return
	}
	if schema.GetField("Age").Tag != "not null" {
		t.Fatal("func Parse test fail for tag parse")
		return
	}
	if schema.GetField("Birthday").Type != "datetime" {
		t.Fatal("func Parse test fail for time parse")
		return
	}
	fmt.Println("test success")
}

func CreateTableTest(t *testing.T) {
	s := OrmTest(t)
	s.Model(&User{})
	err := s.DropTable()
	if err != nil {
		orm.Info("drop table error")
	}
	err = s.CreateTable()
	if err != nil {
		orm.Info("create table error")
	}
	err = s.HasTable()

}
