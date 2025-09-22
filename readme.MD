/your-project
│── api-gateway/                 # API Gateway (REST/GraphQL/gRPC-Gateway)
│   ├── main.go
│   ├── proto/                   # Proto untuk expose ke frontend
│   ├── handler/                 # Handler gateway
│   └── config/
│
│── services/                    # Semua microservices
│   ├── user-service/
│   │   ├── cmd/
│   │   │   └── main.go          # entrypoint gRPC server
│   │   ├── proto/               # .proto untuk user
│   │   ├── internal/
│   │   │   ├── handler/         # implementasi gRPC handler
│   │   │   ├── repository/      # database layer
│   │   │   ├── service/         # business logic
│   │   │   └── utils/           # khusus service ini
│   │   ├── config/
│   │   └── pb/                  # hasil generate .proto
│   │
│   ├── product-service/
│   │   ├── cmd/main.go
│   │   ├── proto/
│   │   ├── internal/{handler,repository,service,utils}
│   │   └── pb/
│   │
│   └── order-service/
│       ├── cmd/main.go
│       ├── proto/
│       ├── internal/{handler,repository,service,utils}
│       └── pb/
│
│── shared/                      # common/shared libs
│   ├── protos/ 
        buf.yml                  # proto yg dipakai lintas service (ex: common types)
│   ├── pb/                      # hasil generate proto shared
│   ├── utils/                   # util global (logger, config loader, middleware, tracing, etc)
│   ├── database/                # db connection builder (Postgres, Redis)
│   ├── events/                  # pubsub/kafka producer-consumer
│   └── auth/                    # jwt, session, dll
│
│── deployments/                 # docker-compose, k8s manifests
│── Makefile
│── go.mod
│── README.md
