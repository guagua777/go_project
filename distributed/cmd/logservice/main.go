package logservice

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/guagua777/distributed/log"
	"github.com/guagua777/distributed/service"
)

func main() {
	log.Run("./distributed.log") // log写入的地址
	host, port := "localhost", "4000"

	fmt.Sprintf("http://%s:%s", host, port) // local:4000

	ctx, err := service.Start(
		context.Background(),
		"Log Service",
		host,
		port,
		log.RegisterHandler,
	)

	if err != nil {
		stlog.Fatalln(err)
	}

	<-ctx.Done()

	fmt.Println("Shutting down log service.")
}
