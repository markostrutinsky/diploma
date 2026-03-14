import { useEffect, useState } from 'react'
import { api, type Unit } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import './Units.css'

const UNIT_TYPES: Record<string, string> = {
  BRIGADE: 'Бригада',
  BATTALION: 'Батальйон',
  COMPANY: 'Рота',
  PLATOON: 'Взвод',
}

export default function Units() {
  const { user } = useAuth()
  const [units, setUnits] = useState<Unit[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ parent_id: '', name: '', unit_type: 'BRIGADE' })

  const canManage = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST', 'BRIGADE_STOREKEEPER', 'BATTALION_STOREKEEPER'].includes(user?.role || '')

  const loadData = () => {
    api.units.list()
      .then((data) => setUnits(Array.isArray(data) ? data : []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    setLoading(true)
    loadData()
  }, [])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.units.create({
        parent_id: form.parent_id ? parseInt(form.parent_id) : undefined,
        name: form.name,
        unit_type: form.unit_type,
      })
      setShowForm(false)
      setForm({ parent_id: '', name: '', unit_type: 'BRIGADE' })
      loadData()
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const getUnitName = (id: number) => units.find((u) => u.id === id)?.name || '-'

  if (loading) {
    return (
      <div className="page-loading">
        <div className="spinner" />
        <p>Завантаження...</p>
      </div>
    )
  }

  return (
    <div className="units-page">
      <div className="page-header">
        <h1>Підрозділи</h1>
        {canManage && (
          <button className="btn btn-primary" onClick={() => setShowForm(true)}>
            + Додати
          </button>
        )}
      </div>

      {showForm && (
        <div className="modal-overlay" onClick={() => setShowForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Новий підрозділ</h3>
            <form onSubmit={handleCreate}>
              <div className="form-group">
                <label>Тип</label>
                <select
                  value={form.unit_type}
                  onChange={(e) => setForm({ ...form, unit_type: e.target.value })}
                  required
                >
                  <option value="BRIGADE">Бригада</option>
                  <option value="BATTALION">Батальйон</option>
                  <option value="COMPANY">Рота</option>
                  <option value="PLATOON">Взвод</option>
                </select>
              </div>
              <div className="form-group">
                <label>Батьківський підрозділ</label>
                <select
                  value={form.parent_id}
                  onChange={(e) => setForm({ ...form, parent_id: e.target.value })}
                >
                  <option value="">Немає</option>
                  {units.map((u) => (
                    <option key={u.id} value={u.id}>{u.name} ({UNIT_TYPES[u.unit_type]})</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>Назва</label>
                <input
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
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
        {units.length === 0 ? (
          <p className="empty-state">Немає підрозділів</p>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>Назва</th>
                <th>Тип</th>
                <th>Батьківський</th>
              </tr>
            </thead>
            <tbody>
              {units.map((u) => (
                <tr key={u.id}>
                  <td>{u.name}</td>
                  <td>{UNIT_TYPES[u.unit_type] || u.unit_type}</td>
                  <td>{u.parent_id ? getUnitName(u.parent_id) : '-'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
