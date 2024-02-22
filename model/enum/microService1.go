package enum

type Gateway string

const (
	ServiceName              Gateway = "microService1"
	ServiceNameLowCase       Gateway = "microservice1"
	MicroService1Json        Gateway = "config/microService1.json"
	MicroService1JsonDaoPath Gateway = "../../../../microService1/config/microService1.json"
)

func (g Gateway) String() string {
	return string(g)
}
