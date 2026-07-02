package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"sync"
)

// 注册服务，目的是给cmd下的registryservice(web service)发送一个post请求
func RegisterService(r Registration) error {
	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL) // 服务注册中心向这个URL来更新一些信息
	if err != nil {
		return err
	}

	// 向自身服务注册handler
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})
	setProvidersHandle()

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}

	// 发送请求
	// const ServicesURL untyped string = "http://localhost:3000/services"
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service. Registry service "+
			"responded with code %v", res.StatusCode)
	}

	return nil
}

func setProvidersHandle() {
	http.HandleFunc("/providers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// 获取所有的providers
			w.WriteHeader(http.StatusOK)
			data, err := json.Marshal(prov.services)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(data)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

	})
}

type serviceUpdateHandler struct{}

// 下发通知的逻辑
func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // 要求必须是POST请求
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body) // 解码Body
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	fmt.Printf("Update received %v\n", p)
	prov.Update(p) // 更新
}

func ShutdownService(url string) error {
	// req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte(url))) // 第二个参数错了
	req, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deregister service. Registry "+
			"service responded with code %v", res.StatusCode)
	}
	return nil
}

// 这个是干什么用的？
// 比如gradingservice依赖于logservice，所以logservice就给gradingservice提供了服务，logservice就是gradingservice的provider
type providers struct {
	services map[ServiceName][]string // 服务所对应的URL，有可能多个URL都能提供同一个服务
	mutex    *sync.RWMutex
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}

func (p *providers) Update(pat patch) { // 当prov收到patch的时候进行更新
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 新增的更新
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}

	// 删除的更新
	for _, patchEntry := range pat.Removed {
		if providerURLs, ok := p.services[patchEntry.Name]; ok {
			for i := range providerURLs {
				p.services[patchEntry.Name] = append(providerURLs[:i], providerURLs[i+1:]...)
			}
		}
	}
}

// 使用服务的名称来找到它所依赖的URL
/*
这里返回的URL只有一个string，因为我们项目简易，一个service只有一个URL，实际上如果是多个URL，这里应该是一个[]string
*/
func (p providers) get(name ServiceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("no providers available for service %v", name)
	}
	// 这里为什么要这样做？因为我们要随机选择一个provider
	// 如果有多个provider，我们希望每个provider都有机会被选择
	// 所以，我们使用随机数来选择一个provider
	idx := int(rand.Float32() * float32(len(providers))) // 使用随机数rand.Float32()生成一个介于 0.0 和 1.0 之间的随机浮点数
	// float32(len(providers)) 将 providers 列表的长度转换为浮点数
	// int(...) 将这个浮点数转换为整数，从而得到一个介于 0 和 len(providers) - 1 之间的整数索引
	return providers[idx], nil
}

// 套一个函数来调用get
func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}
