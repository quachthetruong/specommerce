# Order Service

Microservice quản lý đơn hàng được viết bằng Go.

## Yêu cầu hệ thống

- Go 1.23+
- PostgreSQL

## Cài đặt

1. Clone repository:

```bash
git clone <repository-url>
cd orderservice
```

2. Cài đặt dependencies:

```bash
make deps
```

3. Cài đặt Swagger CLI (tùy chọn):

```bash
make install-swagger
```

## Cấu hình

Chỉnh sửa file `assets/config.yml` để cấu hình database và các thông số khác:

```yaml
env: local
db:
  user: postgres
  password: postgres
  host: localhost
  port: 5432
  dbName: order-db
  enableSsl: false
  autoMigrate: true

server:
  name: "order-service"
  port: 8080
```

## Sử dụng Makefile

### Build ứng dụng

```bash
make build
```

### Chạy ứng dụng

```bash
make run
```

### Chạy trong chế độ development

```bash
make dev
```

### Chạy tests

```bash
make test
```

### Chạy tests với coverage

```bash
make coverage
```

### Generate Swagger documentation

```bash
make generate-swagger
```

### Format Swagger documentation

```bash
make format-swagger
```

### Lint code

```bash
make lint
```

### Clean build artifacts

```bash
make clean
```

### Xem tất cả commands

```bash
make help
```

## API Documentation

Sau khi chạy ứng dụng, truy cập Swagger UI tại:

```
http://localhost:8080/docs/
```

## Cấu trúc dự án

```
orderservice/
├── cmd/                    # Entry points
│   └── server/            # Main server
├── config/                # Configuration structures
├── di/                    # Dependency injection
├── docs/                  # Generated documentation
├── internal/              # Internal packages
├── pkg/                   # Public packages
├── server/                # HTTP server
├── assets/                # Static assets
│   ├── config.yml         # Configuration file
│   └── migrations/        # Database migrations
├── go.mod                 # Go modules
├── go.sum                 # Go modules checksum
├── Makefile               # Build automation
└── README.md              # This file
```

## Development

### Thêm API endpoints

1. Thêm route trong `server/routes.go`
2. Thêm handler trong `server/`
3. Thêm Swagger annotations
4. Generate documentation: `make generate-swagger`

### Database migrations

1. Tạo migration files trong `assets/migrations/`
2. Chạy migrations thủ công hoặc để `autoMigrate: true` trong config

## Troubleshooting

### Lỗi "cannot read config from file"

- Kiểm tra file `assets/config.yml` có tồn tại không
- Kiểm tra cú pháp YAML

### Lỗi database connection

- Kiểm tra PostgreSQL đang chạy
- Kiểm tra thông tin kết nối trong `config.yml`

### Lỗi Swagger

- Cài đặt swag CLI: `make install-swagger`
- Generate lại docs: `make generate-swagger`
