package grades

import (
	"fmt"
	"sync"
)

/**
语义：
1. 定义学生结构体
2. 定义学生集合
**/

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

func (s Student) Average() float32 {
	var result float32
	for _, grade := range s.Grades {
		result += grade.Score
	}
	return result / float32(len(s.Grades))
}

// 使用切片定义类型
// 为什么要定义这个类型？
// 为什么使用切片？
// 因为如果使用数组的话，不知道里面的n该填几，所以要使用切片
type Students []Student

// 定义变量
var (
	students      Students
	studentsMutex sync.Mutex // webservice可能是并发访问的，因此需要加一个锁保证并发访问的安全
)

func (ss Students) GetByID(id int) (*Student, error) {
	for i := range ss {
		if ss[i].ID == id {
			return &ss[i], nil
		}
	}
	return nil, fmt.Errorf("Student with ID %d not found", id)
}

type GradeType string

const (
	GradeQuiz = GradeType("Quiz")
	GradeTest = GradeType("Test")
	GradeExam = GradeType("Exam")
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}
