import { useEffect, useState } from 'react'
import { api, type SupplyRequest, type Resource } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import './Requests.css'

export default function Requests() {
  const { user } = useAuth()
  const [requests, setRequests] = useState<SupplyRequest[]>([])
  const [resources, setResources] = useState<Resource[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [newReq, setNewReq] = useState({ resource_id: '', quantity: 1 })

  const canCreate = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'PLATOON_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST', 'COMPANY_SERGEANT'].includes(user?.role || '')
  const canApprove = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST'].includes(user?.role || '')

  const loadData = () => {
    setLoading(true)
    Promise.all([api.requests.list(), api.inventory.listResources(undefined)])
      .then(([reqs, res]) => {
        const safeReqs = Array.isArray(reqs) ? reqs : []
        const safeRes = Array.isArray(res) ? res : []
        setRequests(safeReqs)
        setResources(safeRes)
        setNewReq((r) => (safeRes.length && !r.resource_id ? { ...r, resource_id: safeRes[0].id } : r))
      })
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    setLoading(true)
    loadData()
  }, [showForm])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.requests.create(newReq)
      setShowForm(false)
      setNewReq({ resource_id: resources[0]?.id || '', quantity: 1 })
      loadData()
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const handleApprove = async (id: string, approved: boolean) => {
    let comment = ""
    
    if (!approved) {
      const reason = window.prompt("Вкажіть причину відхилення (необов'язково):")
      if (reason === null) return 
      comment = reason
    }

    try {
      await (api.requests.approve as any)(id, approved, comment) 
      loadData()
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const statusLabel: Record<string, string> = {
    PENDING: 'Очікує',
    APPROVED: 'Затверджено',
    REJECTED: 'Відхилено',
  }

  if (loading) {
    return (
      <div className="page-loading">
        <div className="spinner" />
        <p>Завантаження...</p>
      </div>
    )
  }

  return (
    <div className="requests-page">
      <div className="page-header">
        <h1>Заявки на постачання</h1>
        {canCreate && (
          <button className="btn btn-primary" onClick={() => setShowForm(true)}>
            + Нова заявка
          </button>
        )}
      </div>

      {showForm && (
        <div className="modal-overlay" onClick={() => setShowForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Нова заявка</h3>
            <form onSubmit={handleCreate}>
              <div className="form-group">
                <label>Ресурс</label>
                <select
                  value={newReq.resource_id}
                  onChange={(e) => setNewReq({ ...newReq, resource_id: e.target.value })}
                  required
                >
                  <option value="">Оберіть ресурс</option>
                  {resources.map((r) => (
                    <option key={r.id} value={r.id}>{r.name} (залишок: {r.quantity})</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>Кількість</label>
                <input
                  type="number"
                  min={1}
                  value={newReq.quantity}
                  onChange={(e) => setNewReq({ ...newReq, quantity: parseInt(e.target.value) || 1 })}
                  required
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

      <div className="card">
        {requests.length === 0 ? (
          <p className="empty-state">Немає заявок</p>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>Ресурс</th>
                <th>Кількість</th>
                <th>Статус</th>
                <th>Дата</th>
                {canApprove && <th>Дії</th>}
              </tr>
            </thead>
            <tbody>
              {requests.map((r) => (
                <tr key={r.id}>
                  <td>{resources.find((res) => res.id === r.resource_id)?.name || r.resource_id}</td>
                  <td>{r.quantity}</td>
                  <td>
                    <span className={`badge badge-${r.status === 'PENDING' ? 'warning' : r.status === 'APPROVED' ? 'success' : 'danger'}`}>
                      {statusLabel[r.status] || r.status}
                    </span>
                  </td>
                  <td>{new Date(r.created_at).toLocaleDateString('uk-UA')}</td>
                  
                  {canApprove && (
                    <td>
                      {r.status === 'PENDING' ? (
                        r.created_by !== String(user?.id) ? (
                          <div className="action-buttons">
                            <button className="btn btn-sm btn-primary" onClick={() => handleApprove(r.id, true)}>
                              Затвердити
                            </button>
                            <button className="btn btn-sm btn-secondary" onClick={() => handleApprove(r.id, false)}>
                              Відхилити
                            </button>
                          </div>
                        ) : (
                          <span className="status-text-own">
                            🔒 Власна заявка
                          </span>
                        )
                      ) : (
                        <span className="status-text-processed">
                          ✓ Оброблено
                        </span>
                      )}
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}