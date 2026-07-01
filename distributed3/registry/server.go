package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const ServicePort = ":3000"
const ServicesURL = "http://localhost" + ServicePort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
	// mutex *sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()

	return nil
}

func (r *registry) remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceURL == url {
			fmt.Println("deubg!!!")
			r.mutex.Lock()
			reg.registrations = append(reg.registrations[:i], r.registrations[i+1:]...) // 第二个参数是元素，使用...将slice展开
			r.mutex.Unlock()
			return nil
		}
	}
	// return fmt.Errorf("Service at URL %s not found.", url) go中有一个编码约定是错误信息字符串应该以小写字母开头且不以标点符号结尾
	return fmt.Errorf("service at URL %s not found", url)
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

type RegistryService struct{}

// 接收registryservice/main.go的请求到达"/services"的请求，根据不同r.Method进行处理
func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received.")
	switch r.Method {
	case http.MethodPost: // POST请求->注册新服务
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL: %s\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r) // 将注册的服务记录添加到reg中
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		url := string(payload)
		fmt.Println("debug???")
		log.Printf("Removing service at URL: %s", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
