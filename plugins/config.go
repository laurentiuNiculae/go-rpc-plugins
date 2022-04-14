package plugins

type Config struct {
	Name              string
	IntegrationPoints []IntegrationPoint
}

type IntegrationPoint struct {
	Interface      string
	GrpcConnection GrpcConnection
}

type GrpcConnection struct {
	Addr string
	Port string
}
