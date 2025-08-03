.PHONY: build run stop clean help

build:
	cd orderservice && go build -o order-server ./cmd/server
	cd paymentservice && go build -o payment-server ./cmd/server
	cd campaignservice && go build -o campaign-server ./cmd/server

build-ui:
	cd adminportal && npm install && npm run build

all: build build-ui

run-order:
	./orderservice/order-server

run-payment:
	./paymentservice/payment-server

run-campaign:
	./campaignservice/campaign-server

run-bg:
	./orderservice/order-server &
	./paymentservice/payment-server &
	./campaignservice/campaign-server &
	@echo "Services started on ports 8080, 8081, 8082"

run:
	@echo "Run services in separate terminals:"
	@echo "  Terminal 1: make run-order"
	@echo "  Terminal 2: make run-payment" 
	@echo "  Terminal 3: make run-campaign"
	@echo "  Terminal 4: make run-ui"
	@echo ""
	@echo "Or run all in background: make run-bg"

run-ui:
	cd adminportal && npm run serve

stop:
	pkill order-server || true
	pkill payment-server || true
	pkill campaign-server || true

clean:
	rm -f orderservice/order-server
	rm -f paymentservice/payment-server
	rm -f campaignservice/campaign-server
	rm -rf adminportal/build

help:
	@echo "Commands:"
	@echo "  make build     - Build Go services"
	@echo "  make build-ui  - Build React UI"
	@echo "  make all       - Build everything"
	@echo "  make run       - Start services"
	@echo "  make run-ui    - Start UI"
	@echo "  make stop      - Stop services"
	@echo "  make clean     - Clean builds"