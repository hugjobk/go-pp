package pp_test

import (
	"pp"
	"testing"
)

type Class struct {
	Id       int
	Name     string
	Students []Student
}

type Student struct {
	Id     int
	Name   string
	friend *Student
}

var (
	s1 = Student{
		Id:     1,
		Name:   "Student 1",
		friend: nil,
	}
	s2 = Student{
		Id:     2,
		Name:   "Student 2",
		friend: &s1,
	}
	s3 = Student{
		Id:     3,
		Name:   "Student 3",
		friend: &s2,
	}
	c = Class{
		Id:       1,
		Name:     "12A1",
		Students: []Student{s1,  s2,  s3},
	}
)

func TestPrintln(t *testing.T) {
	pp.Println(c)
}

func TestPrintIndentln(t *testing.T) {
	pp.PrintIndentln(c, "    ")
}
