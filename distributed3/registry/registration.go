package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceURL       string
	RequiredServices []ServiceName // 当前服务所依赖的其他服务
	ServiceUpdateURL string        // 服务注册中心与当前服务沟通的URL，服务注册中心调用该url
}

type ServiceName string

const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
	PortalService  = ServiceName("PortalService")
)

// 这两个是干什么用的？
// 每一条更新(内部使用)
type patchEntry struct {
	Name ServiceName
	URL  string // 这个url是什么url？
}

// (内部使用)
type patch struct {
	Added   []patchEntry // 增加的条目
	Removed []patchEntry // 减少的条目
}
