module github.com/ojihalaw/sample-grpc/api-gateway

go 1.23.3

require google.golang.org/grpc v1.75.1

require google.golang.org/genproto/googleapis/api v0.0.0-20250818200422-3122310a409c // indirect

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250818200422-3122310a409c // indirect
	google.golang.org/protobuf v1.36.9 // indirect
)

require github.com/ojihalaw/sample-grpc/product-service v0.0.0

replace github.com/ojihalaw/sample-grpc/product-service => ../product-service
