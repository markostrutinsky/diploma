import { useEffect, useState } from 'react'
import { api, type Unit, type User, UNIT_TYPE_NAMES } from '../api/client' // 🔥 Додали UNIT_TYPE_NAMES
import { useAuth } from '../contexts/AuthContext'
import { usePermissions } from '../hooks/usePermissions'
import toast, { Toaster } from 'react-hot-toast'
import './Units.css'

const ROLE_UNIT_CREATION_MAP: Record<string, string[]> = {
  'ADMIN': ['REGION', 'BRANCH', 'DEPARTMENT', 'TEAM'],
  'REGION_DIRECTOR': ['BRANCH', 'DEPARTMENT', 'TEAM'],
  'REGION_LOGISTICIAN': ['BRANCH', 'DEPARTMENT', 'TEAM'],
  'REGION_STOREKEEPER': ['BRANCH', 'DEPARTMENT', 'TEAM'],
  'BRANCH_MANAGER': ['DEPARTMENT', 'TEAM'],
  'BRANCH_LOGISTICIAN': ['DEPARTMENT', 'TEAM'],
  'BRANCH_STOREKEEPER': ['DEPARTMENT', 'TEAM'],
  'DEPT_MANAGER': ['TEAM'],
  'DEPT_SUPERVISOR': ['TEAM'],
}

const REQUIRED_PARENT_TYPE: Record<string, string> = {
  BRANCH: 'REGION',
  DEPARTMENT: 'BRANCH',
  TEAM: 'DEPARTMENT',
}

export default function Units() {
  const { user } = useAuth()
  const [units, setUnits] = useState<Unit[]>([])
  const [managers, setManagers] = useState<User[]>([]) // Колишні commanders
  const [loading, setLoading] = useState(true)
  
  const [showForm, setShowForm] = useState(false)
  const [editingUnit, setEditingUnit] = useState<Unit | null>(null)
  const [form, setForm] = useState({ parent_id: '', name: '', unit_type: '' })

  const [unitToChangeManager, setUnitToChangeManager] = useState<Unit | null>(null)
  const [unitToDelete, setUnitToDelete] = useState<Unit | null>(null)
  const [availableUsers, setAvailableUsers] = useState<User[]>([])
  const [newManagerId, setNewManagerId] = useState<string>('')
  const [isProcessing, setIsProcessing] = useState(false)

  const currentUserRole = user?.role || ''
  const perms = usePermissions()
  const allowedUnitTypes = ROLE_UNIT_CREATION_MAP[currentUserRole] || []
  const canManageUnits = perms.can('unit_manage')

  const expectedParentType = REQUIRED_PARENT_TYPE[form.unit_type]
  const availableParents = units.filter(u => u.unit_type === expectedParentType && u.id !== editingUnit?.id)

  const loadData = () => {
    setLoading(true)
    Promise.all([
      api.units.list().catch(() => []),
      api.users.listManagers().catch(() => []) // Колишній listCommanders
    ])
      .then(([unitsData, managersData]) => {
        setUnits(Array.isArray(unitsData) ? unitsData : [])
        setManagers(Array.isArray(managersData) ? managersData : [])
      })
      .finally(() => setLoading(false))
  }

  useEffect(() => { loadData() }, [])

  const handleOpenCreate = () => {
    setEditingUnit(null)
    setForm({ parent_id: '', name: '', unit_type: allowedUnitTypes[0] || '' })
    setShowForm(true)
  }

  const handleOpenEdit = (unit: Unit) => {
    setEditingUnit(unit)
    setForm({ 
      parent_id: unit.parent_id?.toString() || '', 
      name: unit.name, 
      unit_type: unit.unit_type 
    })
    setShowForm(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (expectedParentType && !form.parent_id) {
      toast.error(`Оберіть батьківську одиницю типу "${UNIT_TYPE_NAMES[expectedParentType]}"`)
      return
    }

    setIsProcessing(true)
    try {
      const payload = {
        parent_id: form.parent_id ? parseInt(form.parent_id) : undefined,
        name: form.name,
        unit_type: form.unit_type,
      }

      if (editingUnit) {
        await api.units.update(editingUnit.id, payload)
        toast.success('Орг. одиницю оновлено')
      } else {
        await api.units.create(payload)
        toast.success('Орг. одиницю створено')
      }
      setShowForm(false)
      loadData()
    } catch (err: any) {
      toast.error(err.message || 'Помилка збереження')
    } finally {
      setIsProcessing(false)
    }
  }

  const handleDelete = async () => {
    if (!unitToDelete) return
    setIsProcessing(true)
    try {
      await api.units.delete(unitToDelete.id)
      toast.success('Одиницю видалено')
      setUnitToDelete(null)
      loadData()
    } catch (err: any) {
      toast.error(err.message || 'Неможливо видалити одиницю, на балансі якої є майно або люди')
    } finally {
      setIsProcessing(false)
    }
  }

  const openChangeManagerModal = (unit: Unit) => {
    setUnitToChangeManager(unit)
    setNewManagerId('')
    
    // Мапінг корпоративних типів на ролі керівників
    const MANAGER_ROLES_MAP: Record<string, string> = {
      'REGION': 'REGION_DIRECTOR',
      'BRANCH': 'BRANCH_MANAGER',
      'DEPARTMENT': 'DEPT_MANAGER',
      'TEAM': 'TEAM_LEAD'
    }
    
    const expectedRole = MANAGER_ROLES_MAP[unit.unit_type] 
    const currentMgr = getManager(unit.id)
    const candidates = managers.filter(c => 
      c.id !== currentMgr?.id && c.role === expectedRole && ['ACTIVE', 'PENDING'].includes(c.status)
    )
    setAvailableUsers(candidates)
  }

  const handleConfirmManagerChange = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!unitToChangeManager || !newManagerId) return
    setIsProcessing(true)
    try {
      await api.units.changeManager(unitToChangeManager.id, newManagerId)
      setUnitToChangeManager(null)
      toast.success('Керівника успішно змінено')
      loadData() 
    } catch (err: any) {
      toast.error(err.message || 'Помилка при зміні керівника')
    } finally {
      setIsProcessing(false)
    }
  }

  const getUnitName = (id: number) => units.find((u) => u.id === id)?.name || '-'
  const getManager = (unitId: number) => {
    const roles = ['REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 'TEAM_LEAD']
    return managers.find(c => c.unit_id === unitId && roles.includes(c.role))
  }

  if (loading) return <div className="page-loading"><div className="spinner" /><p>Завантаження...</p></div>

  return (
    <div className="units-page">
      <Toaster position="top-right" />
      <div className="page-header">
        <h1>Організаційна структура</h1>
        {canManageUnits && <button className="btn btn-primary" onClick={handleOpenCreate}>+ Додати одиницю</button>}
      </div>

      {showForm && (
        <div className="modal-overlay" onClick={() => !isProcessing && setShowForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>{editingUnit ? 'Редагувати орг. одиницю' : 'Нова орг. одиниця'}</h3>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Тип</label>
                <select value={form.unit_type} onChange={(e) => setForm({ ...form, unit_type: e.target.value, parent_id: '' })} required disabled={isProcessing}>
                  {allowedUnitTypes.map(type => <option key={type} value={type}>{UNIT_TYPE_NAMES[type] || type}</option>)}
                </select>
              </div>
              <div className="form-group">
                <label>Батьківська ланка (Підпорядкування)</label>
                <select value={form.parent_id} onChange={(e) => setForm({ ...form, parent_id: e.target.value })} required={!!expectedParentType} disabled={isProcessing}>
                  <option value="">Без підпорядкування</option>
                  {availableParents.map((u) => <option key={u.id} value={u.id}>{u.name} ({UNIT_TYPE_NAMES[u.unit_type] || u.unit_type})</option>)}
                </select>
              </div>
              <div className="form-group">
                <label>Назва</label>
                <input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required disabled={isProcessing} />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowForm(false)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={isProcessing}>{isProcessing ? 'Збереження...' : (editingUnit ? 'Оновити' : 'Створити')}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {unitToDelete && (
        <div className="modal-overlay" onClick={() => !isProcessing && setUnitToDelete(null)}>
          <div className="modal confirm-modal" onClick={(e) => e.stopPropagation()}>
            <h3 style={{ color: '#ef4444' }}>⚠️ Видалення</h3>
            <p>Ви впевнені, що хочете ліквідувати <strong>{unitToDelete.name}</strong>?</p>
            <div className="modal-actions">
              <button className="btn btn-secondary" onClick={() => setUnitToDelete(null)} disabled={isProcessing}>Скасувати</button>
              <button className="btn btn-danger" onClick={handleDelete} disabled={isProcessing}>{isProcessing ? 'Видалення...' : 'Видалити'}</button>
            </div>
          </div>
        </div>
      )}

      {unitToChangeManager && (
        <div className="modal-overlay" onClick={() => !isProcessing && setUnitToChangeManager(null)}>
           <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Зміна керівника</h3>
            <form onSubmit={handleConfirmManagerChange}>
              <div className="form-group">
                <label>Новий керівник</label>
                <select value={newManagerId} onChange={(e) => setNewManagerId(e.target.value)} required disabled={isProcessing}>
                  <option value="" disabled>-- Оберіть кандидата --</option>
                  {availableUsers.map(u => (
                    <option key={u.id} value={u.id}>{u.full_name || u.email}</option>
                  ))}
                </select>
                {availableUsers.length === 0 && (
                  <p style={{ color: '#64748b', fontSize: '12px', marginTop: '8px' }}>
                    Немає вільних користувачів з відповідною роллю. Створіть нового користувача в розділі "Користувачі".
                  </p>
                )}
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setUnitToChangeManager(null)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={!newManagerId || isProcessing}>Зберегти зміни</button>
              </div>
            </form>
           </div>
        </div>
      )}

      <div className="card">
        <table className="data-table">
          <thead>
            <tr>
              <th>Назва</th>
              <th>Тип</th>
              <th>Батьківська ланка</th>
              <th>Керівник</th>
              {canManageUnits && <th>Дії</th>}
            </tr>
          </thead>
          <tbody>
            {units.map((u) => {
              const manager = getManager(u.id);
              return (
                <tr key={u.id}>
                  <td className="font-bold">{u.name}</td>
                  <td><span className="unit-type-tag">{UNIT_TYPE_NAMES[u.unit_type] || u.unit_type}</span></td>
                  <td>{u.parent_id ? getUnitName(u.parent_id) : '-'}</td>
                  <td>
                    {manager ? (
                      <div className="unit-cmdr-wrapper">
                        <span className="cmdr-name-text">
                          {manager.full_name} {manager.status === 'PENDING' && '⏳'}
                        </span>
                        <button 
                          className="btn-cmdr-change" 
                          onClick={() => openChangeManagerModal(u)}
                        >
                          Змінити
                        </button>
                      </div>
                    ) : (
                      <button className="btn-cmdr-assign" onClick={() => openChangeManagerModal(u)}>
                        + Призначити керівника
                      </button>
                    )}
                  </td>

                  {canManageUnits && (
                    <td style={{ width: '1%', whiteSpace: 'nowrap' }}>
                      <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                        <button 
                          style={{ display: 'inline-flex', alignItems: 'center', gap: '6px', padding: '6px 12px', background: '#eff6ff', color: '#2563eb', border: '1px solid #dbeafe', borderRadius: '6px', cursor: 'pointer', fontSize: '13px', fontWeight: 500 }}
                          onClick={() => handleOpenEdit(u)}
                          title="Редагувати"
                        >
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                          Редагувати
                        </button>
                        
                        <button 
                          style={{ display: 'inline-flex', alignItems: 'center', gap: '6px', padding: '6px 12px', background: '#fff1f2', color: '#e11d48', border: '1px solid #ffe4e6', borderRadius: '6px', cursor: 'pointer', fontSize: '13px', fontWeight: 500 }}
                          onClick={() => setUnitToDelete(u)}
                          title="Видалити"
                        >
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                          Видалити
                        </button>
                      </div>
                    </td>
                  )}
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}