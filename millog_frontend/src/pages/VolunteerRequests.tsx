import { useEffect, useState } from 'react'
import { api, type VolunteerRequest } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import './VolunteerRequests.css'

const STATUS_LABELS: Record<string, string> = {
  OPEN: 'Відкрита',
  TAKEN: 'В роботі',
  COMPLETED: 'Виконана',
}

export default function VolunteerRequests() {
  const { user } = useAuth()
  const [requests, setRequests] = useState<VolunteerRequest[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ title: '', description: '' })

  const canCreate = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST', 'BRIGADE_STOREKEEPER', 'BATTALION_STOREKEEPER', 'COMPANY_SERGEANT'].includes(user?.role || '')
  const isVolunteer = user?.role === 'VOLUNTEER'

  const loadData = () => {
    setLoading(true)
    api.volunteerRequests.list()
      .then((data) => setRequests(Array.isArray(data) ? data : []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadData()
  }, [])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.volunteerRequests.create(form)
      setShowForm(false)
      setForm({ title: '', description: '' })
      loadData()
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const handleTake = async (id: string) => {
    try {
      await api.volunteerRequests.take(id)
      loadData()
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const handleComplete = async (id: string) => {
    try {
      await api.volunteerRequests.complete(id)
      loadData()
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const myRequests = requests.filter((r) => r.taken_by === user?.id)

  if (loading) {
    return (
      <div className="page-loading">
        <div className="spinner" />
        <p>Завантаження...</p>
      </div>
    )
  }

  return (
    <div className="volunteer-requests-page">
      <div className="page-header">
        <h1>Заявки для волонтерів</h1>
        {canCreate && (
          <button className="btn btn-primary" onClick={() => setShowForm(true)}>
            + Створити заявку
          </button>
        )}
      </div>

      {showForm && (
        <div className="modal-overlay" onClick={() => setShowForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Нова заявка для волонтерів</h3>
            <form onSubmit={handleCreate}>
              <div className="form-group">
                <label>Назва</label>
                <input
                  value={form.title}
                  onChange={(e) => setForm({ ...form, title: e.target.value })}
                  required
                />
              </div>
              <div className="form-group">
                <label>Опис</label>
                <textarea
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  rows={4}
                />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowForm(false)}>
                  Скасувати
                </button>
                <button type="submit" className="btn btn-primary">Створити</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {isVolunteer && myRequests.length > 0 && (
        <div className="card">
          <h2>Мої заявки</h2>
          <ul className="request-list">
            {myRequests.map((r) => (
              <li key={r.id}>
                <div>
                  <strong>{r.title}</strong>
                  {r.description && <p>{r.description}</p>}
                  <span className="badge badge-warning">{STATUS_LABELS[r.status]}</span>
                </div>
                {r.status === 'TAKEN' && (
                  <button className="btn btn-primary btn-sm" onClick={() => handleComplete(r.id)}>
                    Відмітити виконаною
                  </button>
                )}
              </li>
            ))}
          </ul>
        </div>
      )}

      <div className="card">
        <h2>{isVolunteer ? 'Відкриті заявки' : 'Всі заявки'}</h2>
        {requests.length === 0 ? (
          <p className="empty-state">Немає заявок</p>
        ) : (
          <ul className="request-list">
            {requests.map((r) => (
              <li key={r.id}>
                <div>
                  <strong>{r.title}</strong>
                  {r.description && <p>{r.description}</p>}
                  <span className={`badge badge-${r.status === 'COMPLETED' ? 'success' : 'warning'}`}>
                    {STATUS_LABELS[r.status] || r.status}
                  </span>
                  <span className="request-date">{new Date(r.created_at).toLocaleDateString('uk-UA')}</span>
                </div>
                {isVolunteer && r.status === 'OPEN' && (
                  <button className="btn btn-primary btn-sm" onClick={() => handleTake(r.id)}>
                    Взяти в роботу
                  </button>
                )}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}
