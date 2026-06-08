import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { api, type ContractorMembership } from '../api/client'
import './AdminContractors.css'

const STATUS_LABELS: Record<string, string> = {
  PENDING: 'Очікує',
  APPROVED: 'Схвалено',
  REJECTED: 'Відхилено',
}

const STATUS_BADGE: Record<string, string> = {
  PENDING: 'warning',
  APPROVED: 'success',
  REJECTED: 'danger',
}

type FilterKey = 'ALL' | 'PENDING' | 'APPROVED' | 'REJECTED'

export default function AdminContractors() {
  const [items, setItems] = useState<ContractorMembership[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<FilterKey>('PENDING')
  const [actingId, setActingId] = useState<string | null>(null)

  const load = () => {
    setLoading(true)
    api.contractorMemberships
      .list()
      .then((data) => setItems(Array.isArray(data) ? data : []))
      .catch((e) => toast.error(e instanceof Error ? e.message : 'Помилка завантаження'))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    load()
  }, [])

  const decide = async (id: string, action: 'approve' | 'reject') => {
    setActingId(id)
    try {
      if (action === 'approve') await api.contractorMemberships.approve(id)
      else await api.contractorMemberships.reject(id)
      toast.success(action === 'approve' ? 'Підрядника схвалено' : 'Підрядника відхилено')
      load()
    } catch (e) {
      toast.error(e instanceof Error ? e.message : 'Помилка')
    } finally {
      setActingId(null)
    }
  }

  const counts = {
    ALL: items.length,
    PENDING: items.filter((i) => i.status === 'PENDING').length,
    APPROVED: items.filter((i) => i.status === 'APPROVED').length,
    REJECTED: items.filter((i) => i.status === 'REJECTED').length,
  }

  const filtered = filter === 'ALL' ? items : items.filter((i) => i.status === filter)

  const FILTERS: { key: FilterKey; label: string }[] = [
    { key: 'PENDING', label: 'Очікують' },
    { key: 'APPROVED', label: 'Схвалені' },
    { key: 'REJECTED', label: 'Відхилені' },
    { key: 'ALL', label: 'Усі' },
  ]

  return (
    <div className="admin-contractors">
      <div className="page-header">
        <h1>Підрядники</h1>
        <p className="page-subtitle">
          Зовнішні підрядники реєструються самостійно та бачать спільну дошку завдань.
          Узяти ваше завдання в роботу зможе лише той, кого ви тут схвалите.
        </p>
      </div>

      <div className="contractors-info">
        ℹ️ Коли підрядник уперше намагається взяти ваше завдання, він автоматично з'являється тут
        зі статусом <strong>«Очікує»</strong>. Схвалений підрядник може брати ваші завдання і надалі
        (зокрема співпрацювати з кількома організаціями одночасно).
      </div>

      <div className="contractors-filters">
        {FILTERS.map((f) => (
          <button
            key={f.key}
            className={`contractors-filter-btn ${filter === f.key ? 'active' : ''}`}
            onClick={() => setFilter(f.key)}
          >
            {f.label}
            <span className="filter-count">{counts[f.key]}</span>
          </button>
        ))}
      </div>

      <div className="card">
        {loading ? (
          <div className="page-loading">
            <div className="spinner" />
            <p>Завантаження…</p>
          </div>
        ) : filtered.length === 0 ? (
          <p className="empty-state">
            {filter === 'PENDING'
              ? 'Немає підрядників, що очікують на рішення.'
              : 'Записів не знайдено.'}
          </p>
        ) : (
          <div className="table-scroll">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Підрядник</th>
                  <th>Email</th>
                  <th>Телефон</th>
                  <th>Подано</th>
                  <th>Статус</th>
                  <th style={{ textAlign: 'right' }}>Дії</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((m) => (
                  <tr key={m.id}>
                    <td>{m.contractor_name || '—'}</td>
                    <td>{m.contractor_email || '—'}</td>
                    <td>{m.contractor_phone || '—'}</td>
                    <td>{new Date(m.requested_at).toLocaleDateString('uk-UA')}</td>
                    <td>
                      <span className={`badge badge-${STATUS_BADGE[m.status] || 'neutral'}`}>
                        {STATUS_LABELS[m.status] || m.status}
                      </span>
                    </td>
                    <td>
                      <div className="contractors-actions">
                        {m.status !== 'APPROVED' && (
                          <button
                            className="btn btn-success btn-sm"
                            disabled={actingId === m.id}
                            onClick={() => decide(m.id, 'approve')}
                          >
                            Схвалити
                          </button>
                        )}
                        {m.status !== 'REJECTED' && (
                          <button
                            className="btn btn-danger btn-sm"
                            disabled={actingId === m.id}
                            onClick={() => decide(m.id, 'reject')}
                          >
                            {m.status === 'APPROVED' ? 'Відкликати' : 'Відхилити'}
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
