package main

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	for i := 0; i < 3; i++ {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// 业务逻辑
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("任务失败:", err)
	}
}
