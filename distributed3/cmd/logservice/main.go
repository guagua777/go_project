package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/guagua777/distributed/log"
	"github.com/guagua777/distributed/registry"
	"github.com/guagua777/distributed/service"
)

func main() {
	log.Run("./distributed.log") // log写入的地址
	host, port := "localhost", "4000"

	ServiceURL := fmt.Sprintf("http://%s:%s", host, port) // local:4000

	reg := registry.Registration{
		ServiceName: registry.LogService,
		ServiceURL:  ServiceURL,
	}

	ctx, err := service.Start(
		context.Background(),
		reg,
		host,
		port,
		log.RegisterHandler,
	)

	if err != nil {
		stlog.Fatalln(err)
	}

	// 等待上下文取消
	<-ctx.Done()

	fmt.Println("Shutting down log service.")
}
