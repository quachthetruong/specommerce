import React, { useState, useEffect } from 'react';
import { orderService } from '../services/api';
import { Order, PaginatedResponse } from '../types';

const OrdersPage: React.FC = () => {
  const [orders, setOrders] = useState<PaginatedResponse<Order> | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [sortBy, setSortBy] = useState('-created_at');

  const fetchOrders = async () => {
    try {
      setLoading(true);
      const response = await orderService.searchOrders({
        page: currentPage,
        size: pageSize,
        sort: sortBy,
      });
      setOrders(response);
    } catch (err) {
      setError('Failed to fetch orders');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, [currentPage, pageSize, sortBy]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div className="text-red-600">{error}</div>;

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Orders</h1>
      
      {/* Controls */}
      <div className="bg-white rounded-lg shadow p-4 mb-4">
        <div className="flex gap-4 items-end">
          <div>
            <label className="block text-sm font-medium mb-1">Sort By</label>
            <select 
              value={sortBy} 
              onChange={(e) => {
                setSortBy(e.target.value);
                setCurrentPage(1);
              }}
              className="border rounded px-3 py-2"
            >
              <option value="-created_at">Created Date (Newest)</option>
              <option value="created_at">Created Date (Oldest)</option>
              <option value="-total_amount">Amount (Highest)</option>
              <option value="total_amount">Amount (Lowest)</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Page Size</label>
            <select 
              value={pageSize} 
              onChange={(e) => {
                setPageSize(Number(e.target.value));
                setCurrentPage(1);
              }}
              className="border rounded px-3 py-2"
            >
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
              <option value={100}>100</option>
            </select>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-2 text-left">ID</th>
              <th className="px-4 py-2 text-left">Customer</th>
              <th className="px-4 py-2 text-left">Amount</th>
              <th className="px-4 py-2 text-left">Status</th>
              <th className="px-4 py-2 text-left">Date</th>
            </tr>
          </thead>
          <tbody>
            {orders?.data?.map((order) => (
              <tr key={order.id} className="border-t">
                <td className="px-4 py-2">{order.id}</td>
                <td className="px-4 py-2">{order.customer_name}</td>
                <td className="px-4 py-2">${order.total_amount}</td>
                <td className="px-4 py-2">{order.status}</td>
                <td className="px-4 py-2">{new Date(order.created_at).toLocaleDateString()}</td>
              </tr>
            )) || []}
          </tbody>
        </table>
        
        {(!orders || !orders.data || orders.data.length === 0) && (
          <div className="p-8 text-center text-gray-500">No orders found</div>
        )}
      </div>

      {/* Pagination */}
      {orders && orders.metadata && orders.metadata.total_pages > 1 && (
        <div className="flex justify-between items-center mt-4 p-4 bg-white rounded-lg shadow">
          <div className="text-sm text-gray-600">
            Page {orders.metadata.page_number} of {orders.metadata.total_pages} 
            ({orders.metadata.total} total orders)
          </div>
          <div className="flex gap-2">
            <button
              onClick={() => setCurrentPage(currentPage - 1)}
              disabled={currentPage <= 1}
              className="px-3 py-1 border rounded disabled:opacity-50"
            >
              Previous
            </button>
            <button
              onClick={() => setCurrentPage(currentPage + 1)}
              disabled={currentPage >= orders.metadata.total_pages}
              className="px-3 py-1 border rounded disabled:opacity-50"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default OrdersPage;