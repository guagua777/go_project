package registry

import (
	"encoding/json"
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
