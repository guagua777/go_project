package main

import "context"

// Stream generates values with DoSomething and sends them to out
// until DoSomething returns an error or ctx.Done is closed.
// Stream 从 DoSomething 生成值并将其发送到 out
// 直到 DoSomething 返回错误或 ctx.Done 关闭。
func Stream(ctx context.Context, out chan<- int) error {
	for {
		v, err := DoSomething(ctx)
		if err != nil {
			return err
		}
		select {
		// 接受取消信号
		case <-ctx.Done():
			return ctx.Err()
		// 将 v 发送到 out 通道
		case out <- v:
		}
	}
}

func DoSomething(ctx context.Context) (int, error) {
	return 0, nil
}
