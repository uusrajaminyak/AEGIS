module github.com/uusrajaminyak/aegis-agent

go 1.25.0

require (
	github.com/uusrajaminyak/aegis-backend v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.79.1
)

require (
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260209200024-4cfbd4190f57 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/uusrajaminyak/aegis-backend => ../AEGIS_backend
