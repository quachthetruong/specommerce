import React from 'react';

interface SearchFiltersProps {
  onSortChange: (sort: string) => void;
  onPageSizeChange: (size: number) => void;
  currentSort: string;
  currentPageSize: number;
  sortOptions: { value: string; label: string }[];
}

const SearchFilters: React.FC<SearchFiltersProps> = ({
  onSortChange,
  onPageSizeChange,
  currentSort,
  currentPageSize,
  sortOptions,
}) => {
  return (
    <div className="flex flex-col sm:flex-row gap-4 mb-6">
      <div className="flex-1">
        <label htmlFor="sort" className="block text-sm font-medium text-gray-700 mb-2">
          Sort By
        </label>
        <select
          id="sort"
          value={currentSort}
          onChange={(e) => onSortChange(e.target.value)}
          className="form-input"
        >
          {sortOptions.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      </div>
      
      <div className="flex-1">
        <label htmlFor="pageSize" className="block text-sm font-medium text-gray-700 mb-2">
          Items per page
        </label>
        <select
          id="pageSize"
          value={currentPageSize}
          onChange={(e) => onPageSizeChange(Number(e.target.value))}
          className="form-input"
        >
          <option value={10}>10</option>
          <option value={20}>20</option>
          <option value={50}>50</option>
          <option value={100}>100</option>
        </select>
      </div>
    </div>
  );
};

export default SearchFilters;