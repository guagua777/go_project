package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// 注册服务，目的是给cmd下的registryservice(web service)发送一个post请求
func RegisterService(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
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
