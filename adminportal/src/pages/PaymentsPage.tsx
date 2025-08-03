import React, { useState, useEffect } from 'react';
import { paymentService } from '../services/api';
import { Payment, PaginatedResponse } from '../types';

const PaymentsPage: React.FC = () => {
  const [payments, setPayments] = useState<PaginatedResponse<Payment> | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [sortBy, setSortBy] = useState('-created_at');

  const fetchPayments = async () => {
    try {
      setLoading(true);
      const response = await paymentService.searchPayments({
        page: currentPage,
        size: pageSize,
        sort: sortBy,
      });
      setPayments(response);
    } catch (err) {
      setError('Failed to fetch payments');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPayments();
  }, [currentPage, pageSize, sortBy]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div className="text-red-600">{error}</div>;

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Payments</h1>
      
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
              <th className="px-4 py-2 text-left">Payment ID</th>
              <th className="px-4 py-2 text-left">Order ID</th>
              <th className="px-4 py-2 text-left">Customer</th>
              <th className="px-4 py-2 text-left">Amount</th>
              <th className="px-4 py-2 text-left">Status</th>
              <th className="px-4 py-2 text-left">Date</th>
            </tr>
          </thead>
          <tbody>
            {payments?.data?.map((payment) => (
              <tr key={payment.id} className="border-t">
                <td className="px-4 py-2">{payment.id}</td>
                <td className="px-4 py-2">{payment.order_id}</td>
                <td className="px-4 py-2">{payment.customer_id}</td>
                <td className="px-4 py-2">${payment.total_amount}</td>
                <td className="px-4 py-2">{payment.status}</td>
                <td className="px-4 py-2">{new Date(payment.created_at).toLocaleDateString()}</td>
              </tr>
            )) || []}
          </tbody>
        </table>
        
        {(!payments || !payments.data || payments.data.length === 0) && (
          <div className="p-8 text-center text-gray-500">No payments found</div>
        )}
      </div>

      {/* Pagination */}
      {payments && payments.metadata && payments.metadata.total_pages > 1 && (
        <div className="flex justify-between items-center mt-4 p-4 bg-white rounded-lg shadow">
          <div className="text-sm text-gray-600">
            Page {payments.metadata.page_number} of {payments.metadata.total_pages} 
            ({payments.metadata.total} total payments)
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
              disabled={currentPage >= payments.metadata.total_pages}
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

export default PaymentsPage;