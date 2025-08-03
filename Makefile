.PHONY: build run stop clean help

# Build all services
build:
	cd orderservice && go build -o order-server ./cmd/server
	cd paymentservice && go build -o payment-server ./cmd/server
	cd campaignservice && go build -o campaign-server ./cmd/server

# Build UI
build-ui:
	cd admin-ui && npm install && npm run build

# Build everything
all: build build-ui

# Run services in foreground (individual commands)
run-order:
	./orderservice/order-server

run-payment:
	./paymentservice/payment-server

run-campaign:
	./campaignservice/campaign-server

# Run all services in background (original)
run-bg:
	./orderservice/order-server &
	./paymentservice/payment-server &
	./campaignservice/campaign-server &
	@echo "Services started on ports 8080, 8081, 8082"

# Default run - shows usage
run:
	@echo "Run services in separate terminals:"
	@echo "  Terminal 1: make run-order"
	@echo "  Terminal 2: make run-payment" 
	@echo "  Terminal 3: make run-campaign"
	@echo "  Terminal 4: make run-ui"
	@echo ""
	@echo "Or run all in background: make run-bg"

# Run UI
run-ui:
	cd admin-ui && npm run serve

# Stop services
stop:
	pkill order-server || true
	pkill payment-server || true
	pkill campaign-server || true

# Clean
clean:
	rm -f orderservice/order-server
	rm -f paymentservice/payment-server
	rm -f campaignservice/campaign-server
	rm -rf admin-ui/build

# Help
help:
	@echo "Commands:"
	@echo "  make build     - Build Go services"
	@echo "  make build-ui  - Build React UI"
	@echo "  make all       - Build everything"
	@echo "  make run       - Start services"
	@echo "  make run-ui    - Start UI"
	@echo "  make stop      - Stop services"
	@echo "  make clean     - Clean builds"