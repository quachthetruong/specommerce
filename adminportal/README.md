# SPE Commerce Admin UI

A React-based admin dashboard for managing SPE Commerce services.

## Features

- **Orders Management**: Search and view orders with pagination and sorting
- **Payments Management**: Search and view payments with pagination and sorting  
- **Campaign Management**: Create iPhone campaigns and view winners
- **Responsive Design**: Works on desktop and mobile devices
- **Real-time Data**: Connects to microservices APIs

## Tech Stack

- React 18 with TypeScript
- Tailwind CSS for styling
- React Router for navigation
- Axios for API calls
- Responsive design with mobile support

## API Endpoints

### Order Service (Port 8080)
- `GET /api/admin/v1/orders/search` - Search orders with pagination

### Payment Service (Port 8081)  
- `GET /api/admin/v1/payments/search` - Search payments with pagination

### Campaign Service (Port 8082)
- `POST /api/admin/v1/campaigns/iphones` - Create iPhone campaign
- `GET /api/admin/v1/campaigns/iphones/winners` - Get campaign winners

## Getting Started

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm start
   ```

3. Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

## Environment Variables

Create a `.env` file in the root directory:

```
REACT_APP_API_URL=http://localhost
```

## Build

To build the app for production:

```bash
npm run build
```

## Pages

- `/orders` - Orders search and management
- `/payments` - Payments search and management  
- `/campaigns` - Campaign creation and winner viewing

## Components

- `Layout` - Main layout with navigation
- `Pagination` - Reusable pagination component
- `SearchFilters` - Sort and page size filters
- `LoadingSpinner` - Loading indicator