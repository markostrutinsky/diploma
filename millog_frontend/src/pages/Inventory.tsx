import { useEffect, useState } from 'react'
import { api, type Resource, type ResourceCategory, type Unit } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import './Inventory.css'

export default function Inventory() {
  const { user } = useAuth()
  const [categories, setCategories] = useState<ResourceCategory[]>([])
  const [resources, setResources] = useState<Resource[]>([])
  const [units, setUnits] = useState<Unit[]>([])
  const [filterUnitId, setFilterUnitId] = useState<number | ''>('')
  const [loading, setLoading] = useState(true)
  const [showCategoryForm, setShowCategoryForm] = useState(false)
  const [showResourceForm, setShowResourceForm] = useState(false)
  const [newCat, setNewCat] = useState({ name: '', description: '' })
  const [newRes, setNewRes] = useState({ category_id: '', unit_id: undefined as number | undefined, name: '', quantity: 0, min_quantity: 0 })

  const canManage = ['ADMIN', 'BRIGADE_STOREKEEPER', 'BATTALION_STOREKEEPER', 'COMPANY_SERGEANT'].includes(user?.role || '')

  const loadData = () => {
    const unitId = filterUnitId || undefined
    Promise.all([
      api.inventory.listCategories(),
      api.inventory.listResources(unitId),
      api.units.list(),
    ])
      .then(([cats, res, u]) => {
        const safeCats = Array.isArray(cats) ? cats : []
        const safeRes = Array.isArray(res) ? res : []
        const safeUnits = Array.isArray(u) ? u : []
        setCategories(safeCats)
        setResources(safeRes)
        setUnits(safeUnits)
        setNewRes((r) => (safeCats.length && !r.category_id ? { ...r, category_id: safeCats[0].id } : r))
      })
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    setLoading(true)
    loadData()
  }, [showCategoryForm, showResourceForm, filterUnitId])

  const handleCreateCategory = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.inventory.createCategory(newCat)
      setShowCategoryForm(false)
      setNewCat({ name: '', description: '' })
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const handleCreateResource = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.inventory.createResource({
        ...newRes,
        unit_id: newRes.unit_id || undefined,
      })
      setShowResourceForm(false)
      setNewRes({ category_id: categories[0]?.id || '', unit_id: undefined, name: '', quantity: 0, min_quantity: 0 })
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка')
    }
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
    <div className="inventory-page">
      <div className="page-header">
        <h1>Облік ресурсів</h1>
        <div className="page-actions">
          {units.length > 0 && (
            <select
              value={filterUnitId}
              onChange={(e) => setFilterUnitId(e.target.value ? parseInt(e.target.value, 10) : '')}
              className="filter-select"
            >
              <option value="">Всі склади</option>
              {units.map((u) => (
                <option key={u.id} value={u.id}>{u.name}</option>
              ))}
            </select>
          )}
          {canManage && (
            <>
              <button className="btn btn-secondary" onClick={() => setShowCategoryForm(true)}>
                + Категорія
              </button>
              <button className="btn btn-primary" onClick={() => setShowResourceForm(true)}>
                + Ресурс
              </button>
            </>
          )}
        </div>
      </div>

      {showCategoryForm && (
        <div className="modal-overlay" onClick={() => setShowCategoryForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Нова категорія</h3>
            <form onSubmit={handleCreateCategory}>
              <div className="form-group">
                <label>Назва</label>
                <input
                  value={newCat.name}
                  onChange={(e) => setNewCat({ ...newCat, name: e.target.value })}
                  required
                />
              </div>
              <div className="form-group">
                <label>Опис</label>
                <input
                  value={newCat.description}
                  onChange={(e) => setNewCat({ ...newCat, description: e.target.value })}
                />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowCategoryForm(false)}>
                  Скасувати
                </button>
                <button type="submit" className="btn btn-primary">Створити</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showResourceForm && (
        <div className="modal-overlay" onClick={() => setShowResourceForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Новий ресурс</h3>
            <form onSubmit={handleCreateResource}>
              <div className="form-group">
                <label>Категорія</label>
                <select
                  value={newRes.category_id}
                  onChange={(e) => setNewRes({ ...newRes, category_id: e.target.value })}
                  required
                >
                  <option value="">Оберіть категорію</option>
                  {categories.map((c) => (
                    <option key={c.id} value={c.id}>{c.name}</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>Склад (підрозділ)</label>
                <select
                  value={newRes.unit_id ?? ''}
                  onChange={(e) => setNewRes({ ...newRes, unit_id: e.target.value ? parseInt(e.target.value, 10) : undefined })}
                >
                  <option value="">Не обрано</option>
                  {units.map((u) => (
                    <option key={u.id} value={u.id}>{u.name}</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>Назва</label>
                <input
                  value={newRes.name}
                  onChange={(e) => setNewRes({ ...newRes, name: e.target.value })}
                  required
                />
              </div>
              <div className="form-row-2">
                <div className="form-group">
                  <label>Кількість</label>
                  <input
                    type="number"
                    value={newRes.quantity || ''}
                    onChange={(e) => setNewRes({ ...newRes, quantity: parseInt(e.target.value) || 0 })}
                  />
                </div>
                <div className="form-group">
                  <label>Мін. залишок</label>
                  <input
                    type="number"
                    value={newRes.min_quantity || ''}
                    onChange={(e) => setNewRes({ ...newRes, min_quantity: parseInt(e.target.value) || 0 })}
                  />
                </div>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowResourceForm(false)}>
                  Скасувати
                </button>
                <button type="submit" className="btn btn-primary">Створити</button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className="content-grid">
        <div className="card">
          <h2>Категорії</h2>
          {categories.length === 0 ? (
            <p className="empty-state">Немає категорій</p>
          ) : (
            <ul className="category-list">
              {categories.map((c) => (
                <li key={c.id}>{c.name}</li>
              ))}
            </ul>
          )}
        </div>
        <div className="card card-table">
          <h2>Ресурси</h2>
          {resources.length === 0 ? (
            <p className="empty-state">Немає ресурсів</p>
          ) : (
            <table className="data-table">
              <thead>
                <tr>
                  <th>Назва</th>
                  <th>Склад</th>
                  <th>Кількість</th>
                  <th>Мін.</th>
                  <th>Стан</th>
                </tr>
              </thead>
              <tbody>
                {resources.map((r) => (
                  <tr key={r.id} className={r.quantity <= r.min_quantity ? 'row-warning' : ''}>
                    <td>{r.name}</td>
                    <td>{r.unit_id ? units.find((u) => u.id === r.unit_id)?.name : '-'}</td>
                    <td>{r.quantity}</td>
                    <td>{r.min_quantity}</td>
                    <td>
                      <span className={r.quantity <= r.min_quantity ? 'badge badge-warning' : 'badge badge-success'}>
                        {r.quantity <= r.min_quantity ? 'Нестача' : 'OK'}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  )
}
