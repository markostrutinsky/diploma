import { useEffect, useState } from 'react'
import { api, type Unit, type User } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
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
  const [form, setForm] = useState({ parent_id: '', name: '', unit_type: '' })

  const [unitToChangeCommander, setUnitToChangeCommander] = useState<Unit | null>(null)
  const [availableUsers, setAvailableUsers] = useState<User[]>([])
  const [newCommanderId, setNewCommanderId] = useState<string>('')
  const [isChangingCommander, setIsChangingCommander] = useState(false)

  const currentUserRole = user?.role || ''
  const allowedUnitTypes = ROLE_UNIT_CREATION_MAP[currentUserRole] || []
  
  const canManageUnits = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR'].includes(currentUserRole)

  const expectedParentType = REQUIRED_PARENT_TYPE[form.unit_type]
  const availableParents = units.filter(u => u.unit_type === expectedParentType)

  const loadData = () => {
    setLoading(true)
    Promise.all([
      api.units.list().catch(() => []),
      api.users.listCommanders().catch((err) => {
        console.warn('Не вдалося завантажити командирів:', err)
        return []
      })
    ])
      .then(([unitsData, commandersData]) => {
        setUnits(Array.isArray(unitsData) ? unitsData : [])
        setCommanders(Array.isArray(commandersData) ? commandersData : [])
      })
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadData()
  }, [])

  const handleOpenForm = () => {
    setForm({ 
      parent_id: '', 
      name: '', 
      unit_type: allowedUnitTypes[0] || '' 
    })
    setShowForm(true)
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (expectedParentType && !form.parent_id) {
      alert(`Для цього підрозділу обов'язково треба обрати батьківський підрозділ типу "${UNIT_TYPES[expectedParentType]}"`)
      return
    }

    try {
      await api.units.create({
        parent_id: form.parent_id ? parseInt(form.parent_id) : undefined,
        name: form.name,
        unit_type: form.unit_type,
      })
      setShowForm(false)
      loadData()
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка')
    }
  }

const openChangeCommanderModal = (unit: Unit) => {
  setUnitToChangeCommander(unit)
  setNewCommanderId('')
  
  const expectedRole = `${unit.unit_type}_CMDR` 
  const currentCmdr = getCommander(unit.id)

  const candidates = commanders.filter(c => 
    c.id !== currentCmdr?.id &&                 
    c.role === expectedRole &&                  
    ['ACTIVE', 'PENDING'].includes(c.status)    
  )
  
  setAvailableUsers(candidates)
}

  const handleConfirmCommanderChange = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!unitToChangeCommander || !newCommanderId) return

    setIsChangingCommander(true)
    try {
      await api.units.changeCommander(unitToChangeCommander.id, newCommanderId)
      setUnitToChangeCommander(null)
      loadData() 
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка при зміні командира')
    } finally {
      setIsChangingCommander(false)
    }
  }

  const getUnitName = (id: number) => units.find((u) => u.id === id)?.name || '-'

  const getCommander = (unitId: number) => {
    const commanderRoles = ['BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'PLATOON_CMDR']
    return commanders.find(c => c.unit_id === unitId && commanderRoles.includes(c.role))
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
    <div className="units-page">
      <div className="page-header">
        <h1>Підрозділи</h1>
        {canManageUnits && (
          <button className="btn btn-primary" onClick={handleOpenForm}>
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
                  onChange={(e) => setForm({ ...form, unit_type: e.target.value, parent_id: '' })}
                  required
                >
                  {allowedUnitTypes.map(type => (
                    <option key={type} value={type}>
                      {UNIT_TYPES[type]}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>Батьківський підрозділ</label>
                <select
                  value={form.parent_id}
                  onChange={(e) => setForm({ ...form, parent_id: e.target.value })}
                  required={!!expectedParentType}
                  disabled={!!expectedParentType && availableParents.length === 0}
                >
                  <option value="">Немає</option>
                  {availableParents.map((u) => (
                    <option key={u.id} value={u.id}>{u.name} ({UNIT_TYPES[u.unit_type]})</option>
                  ))}
                </select>
                
                {expectedParentType && availableParents.length === 0 && (
                  <span className="form-help-error">
                    Немає доступних підрозділів типу "{UNIT_TYPES[expectedParentType]}" для підпорядкування. Створіть його спочатку.
                  </span>
                )}
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
                <button 
                  type="submit" 
                  className="btn btn-primary"
                  disabled={!!expectedParentType && availableParents.length === 0}
                >
                  Створити
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {unitToChangeCommander && (
        <div className="modal-overlay" onClick={() => !isChangingCommander && setUnitToChangeCommander(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Зміна командира</h3>
            <p className="modal-subtitle">
              Підрозділ: <strong>{unitToChangeCommander.name}</strong>
            </p>
            
            <form onSubmit={handleConfirmCommanderChange}>
              <div className="form-group mb-4">
                <label>Новий командир <span className="required">*</span></label>
                <select
                  value={newCommanderId}
                  onChange={(e) => setNewCommanderId(e.target.value)}
                  required
                  disabled={isChangingCommander}
                >
                  <option value="" disabled>-- Оберіть кандидата --</option>
                  {availableUsers.length === 0 ? (
                    <option value="" disabled>Немає доступних кандидатів</option>
                  ) : (
                    availableUsers.map(user => (
                      <option key={user.id} value={user.id}>
                        {user.full_name || user.email} {user.unit_id ? `(зараз у ${getUnitName(user.unit_id)})` : '(без підрозділу)'}
                      </option>
                    ))
                  )}
                </select>
              </div>

              <div className="modal-actions">
                <button 
                  type="button" 
                  className="btn btn-secondary" 
                  onClick={() => setUnitToChangeCommander(null)}
                  disabled={isChangingCommander}
                >
                  Скасувати
                </button>
                <button 
                  type="submit" 
                  className="btn btn-primary"
                  disabled={!newCommanderId || isChangingCommander}
                >
                  {isChangingCommander ? 'Збереження...' : 'Зберегти'}
                </button>
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
                <th>Командир</th>
                {canManageUnits && <th>Дії</th>}
              </tr>
            </thead>
            <tbody>
              {units.map((u) => {
                const commander = getCommander(u.id);
                const isPending = commander?.status === 'PENDING';

                return (
                  <tr key={u.id}>
                    <td>{u.name}</td>
                    <td>{UNIT_TYPES[u.unit_type] || u.unit_type}</td>
                    <td>{u.parent_id ? getUnitName(u.parent_id) : '-'}</td>
                    <td>
                      {commander ? (
                        <div className="commander-cell">
                          <span className={isPending ? 'commander-name pending' : 'commander-name'}>
                            {commander.full_name || commander.email}
                          </span>
                          {isPending && (
                            <span className="commander-status-badge">
                              ⏳ Очікує
                            </span>
                          )}
                        </div>
                      ) : (
                        <span className="commander-empty">Не призначено</span>
                      )}
                    </td>
                    {canManageUnits && (
                      <td>
                        <button 
                          className="btn btn-secondary btn-action-small"
                          onClick={() => openChangeCommanderModal(u)}
                        >
                          Змінити
                        </button>
                      </td>
                    )}
                  </tr>
                )
              })}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}