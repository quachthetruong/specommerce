import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/Layout';
import OrdersPage from './pages/OrdersPage';
import PaymentsPage from './pages/PaymentsPage';
import CampaignsPage from './pages/CampaignsPage';

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Navigate to="/orders" replace />} />
          <Route path="/orders" element={<OrdersPage />} />
          <Route path="/payments" element={<PaymentsPage />} />
          <Route path="/campaigns" element={<CampaignsPage />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;