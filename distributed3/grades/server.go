/*
与log里的server一样，都会调用service/service.go中的Start()启动及注册web服务
*/
package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// 由于Start()要求传入一个registerHandlerFunc注册服务函数，因此这里需要实现它，之后传给Start
func RegisterHandlers() {
	handler := new(studentsHandler)
	http.Handle("/students", handler)  // 对应学生这个集合资源
	http.Handle("/students/", handler) // 对应单个资源
}

type studentsHandler struct{}

/*
需要处理集合资源：/students
单个资源（某一个学生）：/students/{id}
或：/students/{id}/grades (获取某个具体学生的所有分数)
*/
func (sh studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/") // 分割的返回值是分出多少块，一个/分出两个块
	switch len(pathSegments) {
	case 2:
		sh.getAll(w, r)
	case 3:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.getOne(w, r, id)
	case 4:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.addGrade(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// 给类型添加方法
func (sh studentsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	// 同一个包下的变量students，可以直接使用
	data, err := sh.toJSON(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data)
}

func (sh studentsHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	data, err := sh.toJSON(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to serialize student: %q", err)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data)
}

func (sh studentsHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	// 定义变量
	var g Grade
	dec := json.NewDecoder(r.Body) // 解码body中的内容
	err = dec.Decode(&g)           // 放到g中
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	student.Grades = append(student.Grades, g)
	fmt.Print(student.Grades)
	w.WriteHeader(http.StatusCreated) // 201状态
	data, err := sh.toJSON(g)
	fmt.Println(data)
	fmt.Println(g)
	if err != nil {
		log.Println(err)
		return // 如果append成功，toJSON失败，后面都就不该执行了
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data) // 数据写回
}

// 是不是可以使用泛型？
// 可以，因为toJSON()函数的参数是一个interface{}，所以可以接受任意类型的变量
// 但是，如果要使用泛型，需要在toJSON()函数中添加一个type参数，例如：
//
//	func (sh studentsHandler) toJSON[T any](obj T) ([]byte, error) {
//		var b bytes.Buffer
//		enc := json.NewEncoder(&b) // 建立编码器
//		err := enc.Encode(obj)     // 对传进来的变量进行编码
//		if err != nil {
//			return nil, fmt.Errorf("failed to serialize %T: %q", obj, err)
//		}
//		return b.Bytes(), nil
//	}
func (sh studentsHandler) toJSON(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b) // 建立编码器
	err := enc.Encode(obj)     // 对传进来的变量进行编码
	if err != nil {
		return nil, fmt.Errorf("failed to serialize students %q", err)
	}
	return b.Bytes(), nil
}
