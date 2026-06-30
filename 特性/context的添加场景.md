go中，什么时候，方法里面需要添加Context参数？什么场景下需要这个
一、核心结论
只要方法满足任意一条，就必须把 context.Context 作为第一个入参（约定命名为 ctx）：
会发起网络请求（HTTP、gRPC、RPC、数据库查询、Redis、MQ 调用）
需要控制超时、截止时间
需要传递请求元数据（TraceID、用户身份、权限令牌、请求 ID）
支持主动取消调用（客户端关闭请求、手动终止任务）
函数调用链条跨多层业务函数（HTTP Handler → Service → DAO）
简单一句话：凡是 IO 操作、可被取消、带请求链路信息的函数，都要加 ctx；纯 CPU 计算、无外部调用的普通工具函数不需要。
二、必须加 Context 的典型场景
1. HTTP / API 接口（强制）
Go HTTP 框架（标准库、Gin、Echo），request 自带 r.Context()。
客户端断开连接（浏览器关闭、超时断开）
设置请求超时
向下游 Service 传递 trace 信息
go
运行
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    user, err := queryUser(ctx, 123) // 向下传递ctx
}
2. 数据库操作（MySQL/Postgres/SQL）
database/sql、GORM、Ent 都支持带 ctx 的查询：
防止慢查询长时间占用连接
设置 SQL 执行超时
请求取消时立刻终止 SQL 执行，释放连接
go
运行
db.WithContext(ctx).Find(&user)
3. 缓存、中间件调用（Redis、Mongo、Elasticsearch）
所有客户端 SDK 都支持 WithContext：
避免网络阻塞无限等待
链路追踪传递 traceId
4. gRPC / RPC 服务（强制规范）
gRPC 默认把上下文在客户端与服务端之间传递，用来：
超时控制
元数据（metadata）传递
调用取消
5. 协程异步任务
如果 goroutine 是跟随当前请求生命周期：
客户端取消请求 → goroutine 跟着退出，防止协程泄漏。
go
运行
go func(ctx context.Context) {
    select {
    case <-time.After(5 * time.Second):
    case <-ctx.Done(): // 请求被取消，直接退出
        return
    }
}(ctx)
6. 多层函数调用链路（Handler → Service → Repository）
整条调用链必须透传同一个 ctx，保证整条链路统一超时、统一取消、统一携带链路信息。
三、不需要加 Context 的场景
满足下面所有条件，可以完全不用 ctx：
纯内存计算，无任何 IO（没有网络、DB、缓存、文件读写）
没有超时控制需求
不需要传递 trace、用户信息、请求元数据
不会被外部请求主动取消
示例：
go
运行
// 纯字符串处理、数值计算，不需要ctx
func MD5Encrypt(s string) string { ... }
func CalcTotal(a, b int) int { ... }
文件读写小补充：
普通本地小文件同步读写一般不加；大文件异步读取、带超时控制的文件操作，依然建议带上 ctx。
四、Go 编码规范（非常重要）
ctx 统一作为第一个参数，变量名固定为 ctx
go
运行
// 规范写法
func GetUser(ctx context.Context, id int64) (*User, error)
不要在函数内部新建空上下文 context.Background() 替代上层 ctx，会丢失取消与超时控制
不要把 ctx 塞到结构体里传递（反模式），优先作为函数入参透传
根上下文起点：
HTTP 入口：r.Context()
定时任务 / 独立后台任务：context.Background() / context.TODO()
五、两种根上下文区别
context.Background()：主函数、初始化、测试用例、常驻后台任务（根上下文首选）
context.TODO()：暂时不确定用哪个上下文，临时占位（重构过渡期使用）
六、精简总结
✅ 要 ctx：DB/Redis/HTTP/RPC/IO、超时控制、可取消、链路追踪、请求级协程、多层业务调用链。
❌ 不要 ctx：纯内存计算、无 IO、不需要取消、不需要传递请求信息。
如果你需要，我可以给你一套 Gin + GORM 整条链路透传 ctx 的标准代码模板。