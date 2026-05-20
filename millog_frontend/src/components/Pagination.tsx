import './Pagination.css';

interface PaginationProps {
  currentPage: number;       // 0-based
  totalPages: number;
  onPageChange: (page: number) => void;
  totalItems: number;
  itemsPerPage: number;
  /** Кількість елементів на поточній сторінці (може відрізнятись для групової пагінації) */
  currentPageItems?: number;
}

export default function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  totalItems,
  itemsPerPage,
  currentPageItems,
}: PaginationProps) {
  if (totalPages <= 1) return null;

  const from = currentPage * itemsPerPage + 1;
  const to = currentPageItems != null
    ? currentPage * itemsPerPage + currentPageItems
    : Math.min((currentPage + 1) * itemsPerPage, totalItems);

  // Будуємо масив кнопок сторінок з «…»
  const pageButtons: (number | 'ellipsis')[] = [];
  if (totalPages <= 7) {
    for (let i = 0; i < totalPages; i++) pageButtons.push(i);
  } else {
    pageButtons.push(0);
    if (currentPage > 2) pageButtons.push('ellipsis');
    for (let i = Math.max(1, currentPage - 1); i <= Math.min(totalPages - 2, currentPage + 1); i++) {
      pageButtons.push(i);
    }
    if (currentPage < totalPages - 3) pageButtons.push('ellipsis');
    pageButtons.push(totalPages - 1);
  }

  return (
    <div className="pagination-wrap">
      <span className="pagination-info">
        {from}–{to} з {totalItems}
      </span>

      <div className="pagination-controls">
        <button
          className="page-btn page-btn--nav"
          onClick={() => onPageChange(currentPage - 1)}
          disabled={currentPage === 0}
          aria-label="Попередня сторінка"
        >
          ‹
        </button>

        {pageButtons.map((item, idx) =>
          item === 'ellipsis' ? (
            <span key={`ell-${idx}`} className="page-ellipsis">…</span>
          ) : (
            <button
              key={item}
              className={`page-btn ${item === currentPage ? 'page-btn--active' : ''}`}
              onClick={() => onPageChange(item)}
            >
              {item + 1}
            </button>
          )
        )}

        <button
          className="page-btn page-btn--nav"
          onClick={() => onPageChange(currentPage + 1)}
          disabled={currentPage === totalPages - 1}
          aria-label="Наступна сторінка"
        >
          ›
        </button>
      </div>
    </div>
  );
}
