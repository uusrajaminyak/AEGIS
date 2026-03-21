module github.com/uusrajaminyak/aegis-agent

go 1.25.0

require (
	github.com/bi-zone/etw v0.0.0-20210519083747-fe9042eb0ea8
	github.com/shirou/gopsutil/v3 v3.24.5
	github.com/uusrajaminyak/aegis-backend v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.41.0
	google.golang.org/grpc v1.79.1
)

require (
	github.com/0xrawsec/golang-etw v1.6.2 // indirect
	github.com/0xrawsec/golang-utils v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260209200024-4cfbd4190f57 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/uusrajaminyak/aegis-backend => ../AEGIS_backend
