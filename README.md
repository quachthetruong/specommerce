### Quick Start

#### 1. Start Infrastructure
```bash
docker-compose up -d
```

#### 2. Run Backend Services
```bash
make run-order    # Terminal 1
make run-payment  # Terminal 2  
make run-campaign # Terminal 3
```

#### 3. Run Admin UI
```bash
make run-ui       # Terminal 4
```

#### 4. Open Admin Portal
```bash
open http://localhost:3000/orders
```

#### 5. Mock Orders
```bash
go run simple_mock_orders.go
```

**Access Points:**
- Admin Portal: http://localhost:3000
- APIs: :8080 (orders), :8081 (payments), :8082 (campaigns)
- Kafka UI: http://localhost:3030


### Other Commands

####  Simulate Orders
```bash
go run simple_mock_orders.go
```

Note: Read the [docs](docs.md) for more details on services, APIs, and database schemas.