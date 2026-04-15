import { useEffect, useState } from 'react'
import { api, type Unit, type User } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import toast, { Toaster } from 'react-hot-toast'
import './Units.css'

const UNIT_TYPES: Record<string, string> = {
  BRIGADE: 'Бригада',
  BATTALION: 'Батальйон',
  COMPANY: 'Рота',
  PLATOON: 'Взвод',
}

const ROLE_UNIT_CREATION_MAP: Record<string, string[]> = {
  'ADMIN': ['BRIGADE', 'BATTALION', 'COMPANY', 'PLATOON'],
  'BRIGADE_CMDR': ['BATTALION', 'COMPANY', 'PLATOON'],
  'BRIGADE_LOGIST': ['BATTALION', 'COMPANY', 'PLATOON'],
  'BRIGADE_STOREKEEPER': ['BATTALION', 'COMPANY', 'PLATOON'],
  'BATTALION_CMDR': ['COMPANY', 'PLATOON'],
  'BATTALION_LOGIST': ['COMPANY', 'PLATOON'],
  'BATTALION_STOREKEEPER': ['COMPANY', 'PLATOON'],
  'COMPANY_CMDR': ['PLATOON'],
  'COMPANY_SERGEANT': ['PLATOON'],
}

const REQUIRED_PARENT_TYPE: Record<string, string> = {
  BATTALION: 'BRIGADE',
  COMPANY: 'BATTALION',
  PLATOON: 'COMPANY',
}

export default function Units() {
  const { user } = useAuth()
  const [units, setUnits] = useState<Unit[]>([])
  const [commanders, setCommanders] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  
  const [showForm, setShowForm] = useState(false)
  const [editingUnit, setEditingUnit] = useState<Unit | null>(null)
  const [form, setForm] = useState({ parent_id: '', name: '', unit_type: '' })

  const [unitToChangeCommander, setUnitToChangeCommander] = useState<Unit | null>(null)
  const [unitToDelete, setUnitToDelete] = useState<Unit | null>(null)
  const [availableUsers, setAvailableUsers] = useState<User[]>([])
  const [newCommanderId, setNewCommanderId] = useState<string>('')
  const [isProcessing, setIsProcessing] = useState(false)

  const currentUserRole = user?.role || ''
  const allowedUnitTypes = ROLE_UNIT_CREATION_MAP[currentUserRole] || []
  const canManageUnits = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'BRIGADE_LOGIST'].includes(currentUserRole)

  const expectedParentType = REQUIRED_PARENT_TYPE[form.unit_type]
  const availableParents = units.filter(u => u.unit_type === expectedParentType && u.id !== editingUnit?.id)

  const loadData = () => {
    setLoading(true)
    Promise.all([
      api.units.list().catch(() => []),
      api.users.listCommanders().catch(() => [])
    ])
      .then(([unitsData, commandersData]) => {
        setUnits(Array.isArray(unitsData) ? unitsData : [])
        setCommanders(Array.isArray(commandersData) ? commandersData : [])
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
      toast.error(`Оберіть батьківський підрозділ типу "${UNIT_TYPES[expectedParentType]}"`)
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
        toast.success('Підрозділ оновлено')
      } else {
        await api.units.create(payload)
        toast.success('Підрозділ створено')
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
      toast.success('Підрозділ видалено')
      setUnitToDelete(null)
      loadData()
    } catch (err: any) {
      toast.error(err.message || 'Неможливо видалити підрозділ, у якому є майно або люди')
    } finally {
      setIsProcessing(false)
    }
  }

  const openChangeCommanderModal = (unit: Unit) => {
    setUnitToChangeCommander(unit)
    setNewCommanderId('')
    const expectedRole = `${unit.unit_type}_CMDR` 
    const currentCmdr = getCommander(unit.id)
    const candidates = commanders.filter(c => 
      c.id !== currentCmdr?.id && c.role === expectedRole && ['ACTIVE', 'PENDING'].includes(c.status)
    )
    setAvailableUsers(candidates)
  }

  const handleConfirmCommanderChange = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!unitToChangeCommander || !newCommanderId) return
    setIsProcessing(true)
    try {
      await api.units.changeCommander(unitToChangeCommander.id, newCommanderId)
      setUnitToChangeCommander(null)
      toast.success('Командира змінено')
      loadData() 
    } catch (err: any) {
      toast.error(err.message || 'Помилка при зміні командира')
    } finally {
      setIsProcessing(false)
    }
  }

  const getUnitName = (id: number) => units.find((u) => u.id === id)?.name || '-'
  const getCommander = (unitId: number) => {
    const roles = ['BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'PLATOON_CMDR']
    return commanders.find(c => c.unit_id === unitId && roles.includes(c.role))
  }

  if (loading) return <div className="page-loading"><div className="spinner" /><p>Завантаження...</p></div>

  return (
    <div className="units-page">
      <Toaster position="top-right" />
      <div className="page-header">
        <h1>Підрозділи</h1>
        {canManageUnits && <button className="btn btn-primary" onClick={handleOpenCreate}>+ Додати</button>}
      </div>

      {showForm && (
        <div className="modal-overlay" onClick={() => !isProcessing && setShowForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>{editingUnit ? 'Редагувати підрозділ' : 'Новий підрозділ'}</h3>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Тип</label>
                <select value={form.unit_type} onChange={(e) => setForm({ ...form, unit_type: e.target.value, parent_id: '' })} required disabled={isProcessing}>
                  {allowedUnitTypes.map(type => <option key={type} value={type}>{UNIT_TYPES[type]}</option>)}
                </select>
              </div>
              <div className="form-group">
                <label>Батьківський підрозділ</label>
                <select value={form.parent_id} onChange={(e) => setForm({ ...form, parent_id: e.target.value })} required={!!expectedParentType} disabled={isProcessing}>
                  <option value="">Немає</option>
                  {availableParents.map((u) => <option key={u.id} value={u.id}>{u.name} ({UNIT_TYPES[u.unit_type]})</option>)}
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
            <p>Ви впевнені, що хочете видалити <strong>{unitToDelete.name}</strong>?</p>
            <div className="modal-actions">
              <button className="btn btn-secondary" onClick={() => setUnitToDelete(null)} disabled={isProcessing}>Скасувати</button>
              <button className="btn btn-danger" onClick={handleDelete} disabled={isProcessing}>{isProcessing ? 'Видалення...' : 'Видалити'}</button>
            </div>
          </div>
        </div>
      )}

      {/* Модалка зміни командира залишається такою ж, як була... */}
      {unitToChangeCommander && (
        <div className="modal-overlay" onClick={() => !isProcessing && setUnitToChangeCommander(null)}>
           <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Зміна командира</h3>
            <form onSubmit={handleConfirmCommanderChange}>
              <div className="form-group">
                <label>Новий командир</label>
                <select value={newCommanderId} onChange={(e) => setNewCommanderId(e.target.value)} required disabled={isProcessing}>
                  <option value="" disabled>-- Оберіть кандидата --</option>
                  {availableUsers.map(u => (
                    <option key={u.id} value={u.id}>{u.full_name || u.email}</option>
                  ))}
                </select>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setUnitToChangeCommander(null)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={!newCommanderId || isProcessing}>Зберегти</button>
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
              <th>Підпорядкування</th>
              <th>Командир</th>
              {canManageUnits && <th>Дії</th>}
            </tr>
          </thead>
          <tbody>
            {units.map((u) => {
              const commander = getCommander(u.id);
              return (
                <tr key={u.id}>
                  <td className="font-bold">{u.name}</td>
                  <td><span className="unit-type-tag">{UNIT_TYPES[u.unit_type]}</span></td>
                  <td>{u.parent_id ? getUnitName(u.parent_id) : '-'}</td>
                  {/* Колонка: Командир */}
<td>
  {commander ? (
    <div className="unit-cmdr-wrapper">
      <span className="cmdr-name-text">
        {commander.full_name} {commander.status === 'PENDING' && '⏳'}
      </span>
      <button 
        className="btn-cmdr-change" 
        onClick={() => openChangeCommanderModal(u)}
      >
        Змінити
      </button>
    </div>
  ) : (
    <button className="btn-cmdr-assign" onClick={() => openChangeCommanderModal(u)}>
      + Призначити командира
    </button>
  )}
</td>

{/* Колонка: Дії */}
{canManageUnits && (
  <td>
    <div className="unit-actions-container">
      <button 
        className="btn-unit-action btn-unit-edit" 
        onClick={() => handleOpenEdit(u)}
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
        Редагувати
      </button>
      
      <button 
        className="btn-unit-action btn-unit-delete" 
        onClick={() => setUnitToDelete(u)}
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