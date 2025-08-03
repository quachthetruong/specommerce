# SPgroup E-Commerce Platform
## Interview Coding Assignment - Empowering the Future of Energy

### Overview
This project implements a high-throughput, low-latency e-commerce platform designed to handle hot-sale campaigns with concurrent order processing. The system supports promotional campaigns where the first 50 customers with orders above SGD 200 can win limited quantities of free iPhones (20 total gifts).

### Architecture
The system follows a microservices architecture with event-driven communication:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Admin Portal  │    │   Order Service  │    │ Campaign Service│
│   (React/TS)    │    │   (Go/Postgres)  │    │  (Go/Postgres) │
│   Port: 3000    │◄──►│   Port: 8080     │◄──►│   Port: 8082   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │                        │
                                │       Kafka            │
                                │   ┌─────────────┐      │
                                └──►│   Topics:   │◄─────┘
                                    │ - payment_* │
                       ┌────────────┤ - order_*   │
                       │            └─────────────┘
                       ▼
           ┌─────────────────────┐
           │  Payment Service    │
           │   (Go/Postgres)     │
           │   Port: 8081        │
           └─────────────────────┘
```

### Core Features

#### Business Requirements
- **Normal Order Processing**: Customers can place orders and make payments
- **Hot Sale Campaign**: 
  - First 50 customers with successful transactions
  - Order amount > SGD 200 (configurable)
  - One gift per customer (no duplicate wins)
  - Limited to 20 free iPhones (first-come-first-serve)

#### Technical Implementation
- **Microservices**: Order, Payment, and Campaign services with separate databases
- **Event-Driven**: Kafka for inter-service communication
- **Concurrency Safe**: Atomic operations for campaign prize allocation
- **Admin Portal**: React/TypeScript frontend for order management
- **High Throughput**: Optimized for >10,000 transactions/minute

### Services

#### 1. Order Service (Port: 8080)
- **Database**: PostgreSQL (Port: 5432)
- **Responsibilities**:
  - Order creation and management
  - Payment request orchestration
  - Order status updates
- **API Endpoints**:
  - `POST /api/orders` - Create new order
  - `GET /api/admin/orders` - Admin order listing with pagination/sorting

#### 2. Payment Service (Port: 8081)
- **Database**: PostgreSQL (Port: 5433)
- **Responsibilities**:
  - Payment processing (mocked for success)
  - Payment status tracking
  - Payment response events
- **API Endpoints**:
  - `GET /api/admin/payments` - Admin payment records

#### 3. Campaign Service (Port: 8082)
- **Database**: PostgreSQL (Port: 5434)
- **Responsibilities**:
  - Campaign configuration management
  - Prize eligibility validation
  - Winner tracking and limits
- **API Endpoints**:
  - `GET /api/admin/campaigns` - Campaign status
  - `POST /api/admin/campaigns` - Create/update campaigns

#### 4. Admin Portal (Port: 3000)
- **Technology**: React + TypeScript + Tailwind CSS
- **Features**:
  - Order search and filtering
  - Payment record viewing
  - Campaign management
  - Pagination (10 records per page)
  - Sorting by creation time and amount
  - iPhone inventory tracking

### Database Design

#### Orders Table
```sql
- id (UUID, Primary Key)
- customer_name (VARCHAR)
- customer_email (VARCHAR)
- amount (DECIMAL)
- status (VARCHAR)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

#### Payments Table
```sql
- id (UUID, Primary Key)
- order_id (UUID, Foreign Key)
- amount (DECIMAL)
- status (VARCHAR)
- processed_at (TIMESTAMP)
```

#### Campaigns Table
```sql
- id (UUID, Primary Key)
- name (VARCHAR)
- min_amount (DECIMAL)
- max_winners (INTEGER)
- current_winners (INTEGER)
- is_active (BOOLEAN)
- created_at (TIMESTAMP)
```

### Event Flow
1. **Order Creation**: Customer places order → Order Service creates order
2. **Payment Processing**: Order Service → Payment Request Event → Payment Service
3. **Payment Success**: Payment Service → Payment Response Event → Order Service
4. **Campaign Check**: Order Service → Order Success Event → Campaign Service
5. **Prize Award**: Campaign Service validates eligibility and awards prizes

### Prerequisites
- **Go**: Version 1.19+
- **Node.js**: Version 16+
- **Docker & Docker Compose**: Latest version
- **Make**: Build automation tool

### Quick Start

#### 1. Start Infrastructure
```bash
# Start databases and Kafka
docker-compose up -d

# Verify services are running
docker-compose ps
```

#### 2. Build Services
```bash
# Build all Go services
make build

# Build React admin portal
make build-ui
```

#### 3. Run Services
Open 4 separate terminals and run:

**Terminal 1 - Order Service:**
```bash
make run-order
```

**Terminal 2 - Payment Service:**
```bash
make run-payment
```

**Terminal 3 - Campaign Service:**
```bash
make run-campaign
```

**Terminal 4 - Admin Portal:**
```bash
make run-ui
```

#### 4. Access Applications
- **Admin Portal**: http://localhost:3000
- **Order Service API**: http://localhost:8080
- **Payment Service API**: http://localhost:8081
- **Campaign Service API**: http://localhost:8082
- **Kafka UI**: http://localhost:3030

### Testing the System

#### 1. Create a Campaign
```bash
curl -X POST http://localhost:8082/api/admin/campaigns \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 13 Pro Campaign",
    "min_amount": 200.00,
    "max_winners": 20,
    "is_active": true
  }'
```

#### 2. Simulate Orders
Use the provided order simulator:
```bash
go run simple_mock_orders.go
```

#### 3. Monitor via Admin Portal
- View orders at http://localhost:3000/orders
- Check payments at http://localhost:3000/payments  
- Monitor campaigns at http://localhost:3000/campaigns

### API Documentation
Swagger documentation is available for each service:
- Order Service: `./orderservice/docs/openapi/`
- Payment Service: `./paymentservice/docs/openapi/`
- Campaign Service: `./campaignservice/docs/openapi/`

### Concurrency & Atomicity
The system handles concurrent requests through:
- **Database Transactions**: ACID compliance for critical operations
- **Optimistic Locking**: Version-based conflict resolution
- **Event Ordering**: Kafka ensures message ordering within partitions
- **Idempotency**: Duplicate request protection

### Monitoring & Observability
- **Kafka UI**: Real-time message flow monitoring
- **Application Logs**: Structured logging in each service
- **Health Checks**: Service availability endpoints

### Project Structure
```
specommerce/
├── adminportal/          # React admin interface
├── orderservice/         # Order management microservice
├── paymentservice/       # Payment processing microservice  
├── campaignservice/      # Campaign management microservice
├── docker-compose.yml    # Infrastructure setup
├── Makefile             # Build and run automation
└── simple_mock_orders.go # Order simulation tool
```

### Assumptions & Design Decisions
1. **Payment Success**: All payments are mocked to succeed for demonstration
2. **Single Campaign**: System supports one active campaign at a time
3. **Customer Identification**: Based on email address for duplicate prevention
4. **Prize Allocation**: First-come-first-serve based on successful payment timestamp
5. **Data Consistency**: Eventually consistent across services via event sourcing

### Future Enhancements
- Authentication and authorization
- Real payment gateway integration
- Advanced campaign rules engine
- Metrics and alerting system
- Load balancing and auto-scaling
- Data analytics and reporting

### Troubleshooting

#### Common Issues
1. **Port Conflicts**: Ensure ports 3000, 5432-5434, 8080-8082, 9093 are available
2. **Database Connection**: Verify PostgreSQL containers are running
3. **Kafka Connectivity**: Check Kafka container health
4. **Build Failures**: Ensure Go modules are properly downloaded

#### Cleanup
```bash
# Stop all services
make stop

# Clean build artifacts
make clean

# Stop and remove containers
docker-compose down -v
```

---

**Author**: SPgroup Interview Candidate  
**Date**: August 2025  
**Technology Stack**: Go, React/TypeScript, PostgreSQL, Kafka, Docker