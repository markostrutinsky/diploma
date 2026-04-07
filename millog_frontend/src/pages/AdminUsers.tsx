import { useEffect, useState } from 'react'
import { api, type CreateUserRequest, type Unit, type User } from '../api/client'
import './AdminUsers.css'

const ROLES = [
  { value: 'ADMIN', label: 'Адмін' },
  { value: 'BRIGADE_CMDR', label: 'Командир бригади' },
  { value: 'BATTALION_CMDR', label: 'Командир батальйону' },
  { value: 'COMPANY_CMDR', label: 'Командир роти' },
  { value: 'PLATOON_CMDR', label: 'Командир взводу' },
  { value: 'BRIGADE_LOGIST', label: 'Логіст бригади' },
  { value: 'BRIGADE_STOREKEEPER', label: 'Комірник бригади' },
  { value: 'BATTALION_LOGIST', label: 'Логіст батальйону' },
  { value: 'BATTALION_STOREKEEPER', label: 'Комірник батальйону' },
  { value: 'COMPANY_SERGEANT', label: 'Старшина роти' },
] as const

const ROLE_CREATION_MAP: Record<string, string[]> = {
  'ADMIN': [
    'ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'PLATOON_CMDR',
    'BRIGADE_LOGIST', 'BRIGADE_STOREKEEPER', 'BATTALION_LOGIST', 'BATTALION_STOREKEEPER', 'COMPANY_SERGEANT'
  ],
  'BRIGADE_CMDR': [
    'BATTALION_CMDR', 'COMPANY_CMDR', 'PLATOON_CMDR',
    'BRIGADE_LOGIST', 'BRIGADE_STOREKEEPER', 'BATTALION_LOGIST', 'BATTALION_STOREKEEPER', 'COMPANY_SERGEANT'
  ],
  'BATTALION_CMDR': [
    'COMPANY_CMDR', 'PLATOON_CMDR',
    'BATTALION_LOGIST', 'BATTALION_STOREKEEPER', 'COMPANY_SERGEANT'
  ],
  'COMPANY_CMDR': [
    'PLATOON_CMDR', 'COMPANY_SERGEANT'
  ],
  'BRIGADE_LOGIST': [
    'BRIGADE_STOREKEEPER', 'BATTALION_LOGIST', 'BATTALION_STOREKEEPER'
  ],
  'BATTALION_LOGIST': [
    'BATTALION_STOREKEEPER'
  ],
}

export default function AdminUsers() {
  const [formUnits, setFormUnits] = useState<Unit[]>([])
  const [allUnits, setAllUnits] = useState<Unit[]>([])
  const [usersList, setUsersList] = useState<User[]>([])
  const [activeTab, setActiveTab] = useState<'all' | 'in_unit' | 'reserve'>('all')

  const [currentUserRole, setCurrentUserRole] = useState<string | null>(null)
  const [currentUserId, setCurrentUserId] = useState<string | null>(null)
  const [availableRoles, setAvailableRoles] = useState<typeof ROLES[number][]>([])
  
  const [showForm, setShowForm] = useState(false)
  const [editingUser, setEditingUser] = useState<User | null>(null)
  
  const [confirmModal, setConfirmModal] = useState<{ type: 'block' | 'unblock', id: string, name: string } | null>(null)
  
  const [form, setForm] = useState<Omit<CreateUserRequest, 'role'> & { role: (typeof ROLES)[number]['value'] }>({
    email: '',
    full_name: '',
    role: 'COMPANY_SERGEANT',
  })
  
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  const canManagePersonnel = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'PLATOON_CMDR'].includes(currentUserRole || '')

  const loadUsersAndUnits = () => {
    api.units.list()
      .then(data => setAllUnits(Array.isArray(data) ? data : []))
      .catch(console.error)

    api.users.getVisible() 
      .then(data => setUsersList(Array.isArray(data) ? data : []))
      .catch(console.error)
  }

  useEffect(() => {
    loadUsersAndUnits()
    api.auth.me()
      .then((user) => {
        setCurrentUserRole(user.role)
        setCurrentUserId(String(user.id))
      })
      .catch((err) => console.error("Не вдалося отримати користувача", err))
  }, [])

  useEffect(() => {
    if (currentUserRole) {
      const allowedRoles = ROLE_CREATION_MAP[currentUserRole] || []
      const filtered = ROLES.filter(r => allowedRoles.includes(r.value))
      setAvailableRoles(filtered)

      if (filtered.length > 0 && !allowedRoles.includes(form.role) && !editingUser) {
        setForm(prev => ({ ...prev, role: filtered[0].value as any, unit_id: undefined }))
      }
    }
  }, [currentUserRole, editingUser])

  useEffect(() => {
    if (!form.role || !currentUserRole) return

    if (currentUserRole === 'ADMIN') {
      api.units.getAvailableForRole(form.role)
        .then((data) => setFormUnits(Array.isArray(data) ? data : []))
        .catch(() => setFormUnits([]))
    } else {
      api.units.getMyHierarchyForRole(form.role)
        .then((data) => setFormUnits(Array.isArray(data) ? data : []))
        .catch(() => setFormUnits([]))
    }
  }, [form.role, currentUserRole])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null); setSuccess(null); setLoading(true)

    try {
      const body: CreateUserRequest = {
        email: form.email.trim(),
        full_name: form.full_name.trim(),
        role: form.role,
      }
      if (form.username?.trim()) body.username = form.username.trim()
      if (form.phone?.trim()) body.phone = form.phone.trim()
      if (form.unit_id) body.unit_id = form.unit_id

      const res = await api.admin.createUser(body)
      setSuccess(res.message)
      
      const defaultRole = availableRoles.length > 0 ? availableRoles[0].value : 'COMPANY_SERGEANT'
      setForm({ email: '', full_name: '', role: defaultRole as any, username: '', phone: '', unit_id: undefined })
      
      loadUsersAndUnits()
      setTimeout(() => { setShowForm(false); setSuccess(null) }, 1500)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Помилка')
    } finally {
      setLoading(false)
    }
  }

  const handleUpdateRole = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null); setSuccess(null); setLoading(true)

    if (!editingUser) return

    try {
      await api.users.updateRole(editingUser.id, form.role, form.unit_id || null)
      setSuccess("Кадрове переміщення виконано успішно")
      loadUsersAndUnits()
      
      setTimeout(() => { setEditingUser(null); setSuccess(null) }, 1500)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Помилка')
    } finally {
      setLoading(false)
    }
  }

  const handleConfirmStatusChange = async () => {
    if (!confirmModal) return

    setLoading(true)
    setError(null)
    setSuccess(null)

    try {
      if (confirmModal.type === 'block') {
        await api.users.block(confirmModal.id)
        setSuccess(`Користувача ${confirmModal.name} заблоковано`)
      } else {
        await api.users.unblock(confirmModal.id)
        setSuccess(`Користувача ${confirmModal.name} розблоковано`)
      }
      
      loadUsersAndUnits() 
      setTimeout(() => setSuccess(null), 3000)
      setConfirmModal(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Помилка виконання дії')
    } finally {
      setLoading(false)
    }
  }

  const openEditModal = (user: User) => {
    setEditingUser(user)
    setForm({
      email: user.email,
      full_name: user.full_name || '',
      role: user.role as any,
      unit_id: user.unit_id || undefined,
    })
    setSuccess(null)
    setError(null)
  }

  const filteredUsers = usersList.filter(user => {
    if (currentUserRole !== 'ADMIN' && user.role === 'ADMIN') {
      return false;
    }
    // Волонтери не є військовим резервом, тому приховуємо їх з цієї вкладки
    if (activeTab === 'reserve') return !user.unit_id && user.role !== 'ADMIN' && user.role !== 'VOLUNTEER';
    if (activeTab === 'in_unit') return !!user.unit_id;
    return true;
  })

  const getUnitName = (unitId?: number | null) => {
    if (!unitId) return null;
    const unit = allUnits.find(u => u.id === unitId);
    return unit ? unit.name : `ID: ${unitId}`;
  }

  return (
    <div className="admin-users">
      <div className="page-header page-header-flex">
        <div>
          <h1>Кадровий облік</h1>
          <p className="page-subtitle">
            Керування особовим складом та призначення на посади.
          </p>
        </div>
        {canManagePersonnel && (
          <button 
            className="btn btn-primary" 
            onClick={() => {
              setSuccess(null); setError(null);
              setEditingUser(null);
              const defaultRole = availableRoles.length > 0 ? availableRoles[0].value : 'COMPANY_SERGEANT'
              setForm({ email: '', full_name: '', role: defaultRole as any, unit_id: undefined })
              setShowForm(true);
            }}
            disabled={availableRoles.length === 0}
          >
            + Створити
          </button>
        )}
      </div>

      {(showForm || editingUser) && (
        <div className="modal-overlay" onClick={() => !loading && (setShowForm(false), setEditingUser(null))}>
           <div className="modal modal-form" onClick={(e) => e.stopPropagation()}>
            <h3>{editingUser ? `Переміщення: ${editingUser.full_name}` : 'Створити обліковий запис'}</h3>
            {!editingUser && (
              <p className="modal-description">Користувач отримає лист з посиланням для встановлення паролю.</p>
            )}
            
            <form className="admin-form" onSubmit={editingUser ? handleUpdateRole : handleSubmit}>
              
              {!editingUser && (
                <div className="form-group">
                  <label>Email <span className="required">*</span></label>
                  <input type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} required />
                </div>
              )}
              
              <div className="form-group">
                <label>ПІБ {editingUser && <span className="text-muted">(тільки для читання)</span>}</label>
                <input 
                  type="text" 
                  value={form.full_name} 
                  onChange={(e) => !editingUser && setForm({ ...form, full_name: e.target.value })} 
                  required 
                  readOnly={!!editingUser} 
                  className={editingUser ? "input-readonly" : ""}
                />
              </div>
              
              {!editingUser && (
                <div className="form-row-2">
                  <div className="form-group">
                    <label>Логін</label>
                    <input type="text" value={form.username ?? ''} onChange={(e) => setForm({ ...form, username: e.target.value || undefined })} />
                  </div>
                  <div className="form-group">
                    <label>Телефон</label>
                    <input type="tel" value={form.phone ?? ''} onChange={(e) => setForm({ ...form, phone: e.target.value || undefined })} />
                  </div>
                </div>
              )}

              <div className="form-group">
                <label>{editingUser ? 'Нова посада / Звання' : 'Посада / Звання'}</label>
                <select
                  value={form.role}
                  onChange={(e) => {
                    const newRole = e.target.value as (typeof ROLES)[number]['value']
                    setForm({ ...form, role: newRole, unit_id: undefined }) 
                  }}
                  disabled={availableRoles.length === 0}
                >
                  {availableRoles.map((r) => (
                    <option key={r.value} value={r.value}>{r.label}</option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>{editingUser ? 'Новий підрозділ' : 'Підрозділ (залиште пустим для Резерву)'}</label>
                <select
                  value={form.unit_id ?? ''}
                  onChange={(e) => {
                    const val = e.target.value
                    setForm({ ...form, unit_id: val ? parseInt(val, 10) : undefined })
                  }}
                  disabled={formUnits.length === 0}
                >
                  <option value="">
                    {formUnits.length === 0 ? 'Немає вільних підрозділів (буде в резерві)' : '-- В Резерв --'}
                  </option>
                  {formUnits.map((u) => (
                    <option key={u.id} value={u.id}>{u.name}</option>
                  ))}
                </select>
              </div>

              {error && <div className="form-error modal-alert">{error}</div>}
              {success && <div className="form-success modal-alert">{success}</div>}

              <div className="modal-actions">
                <button 
                  type="button" 
                  className="btn btn-secondary" 
                  onClick={() => { setShowForm(false); setEditingUser(null); }}
                  disabled={loading}
                >
                  Скасувати
                </button>
                <button 
                  type="submit" 
                  className="btn btn-primary" 
                  disabled={loading || availableRoles.length === 0}
                >
                  {loading ? 'Обробка...' : (editingUser ? 'Підтвердити переміщення' : 'Створити')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {confirmModal && (
        <div className="modal-overlay" onClick={() => !loading && setConfirmModal(null)}>
          <div className="modal modal-form" onClick={(e) => e.stopPropagation()}>
            <h3 style={{ marginBottom: '16px' }}>
              {confirmModal.type === 'block' ? 'Підтвердження блокування' : 'Підтвердження розблокування'}
            </h3>
            
            <p className="modal-description" style={{ fontSize: '15px', color: '#495057', marginBottom: '24px' }}>
              {confirmModal.type === 'block' 
                ? `Ви впевнені, що хочете перевести користувача ${confirmModal.name} в резерв (заблокувати)? Він втратить доступ до системи.`
                : `Повернути користувача ${confirmModal.name} до активного складу? Він опиниться в Кадровому резерві.`
              }
            </p>

            {error && <div className="form-error modal-alert">{error}</div>}

            <div className="modal-actions">
              <button 
                type="button" 
                className="btn btn-secondary" 
                onClick={() => setConfirmModal(null)}
                disabled={loading}
              >
                Скасувати
              </button>
              <button 
                type="button" 
                style={{ 
                  backgroundColor: confirmModal.type === 'block' ? '#dc3545' : '#198754', 
                  borderColor: confirmModal.type === 'block' ? '#dc3545' : '#198754',
                  color: 'white'
                }}
                className="btn" 
                onClick={handleConfirmStatusChange}
                disabled={loading}
              >
                {loading ? 'Обробка...' : 'Підтвердити'}
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="card card-table card-table-spaced">
        <div className="table-header-with-tabs">
          <h2>Особовий склад</h2>
          <div className="inventory-tabs">
            <button className={`tab-btn ${activeTab === 'all' ? 'active' : ''}`} onClick={() => setActiveTab('all')}>Всі</button>
            <button className={`tab-btn ${activeTab === 'in_unit' ? 'active' : ''}`} onClick={() => setActiveTab('in_unit')}>У підрозділах</button>
            <button className={`tab-btn ${activeTab === 'reserve' ? 'active' : ''}`} onClick={() => setActiveTab('reserve')}>Кадровий резерв</button>
          </div>
        </div>

        {filteredUsers.length === 0 ? (
          <p className="empty-state">
            {activeTab === 'reserve' ? 'Кадровий резерв порожній' : 'Немає користувачів'}
          </p>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>ПІБ</th>
                <th>Посада</th>
                <th>Підрозділ</th>
                <th>Статус</th>
                {canManagePersonnel && <th>Дії</th>}
              </tr>
            </thead>
            <tbody>
              {filteredUsers.map(u => {
                const roleLabel = u.role === 'VOLUNTEER' 
                  ? 'Волонтер' 
                  : ROLES.find(r => r.value === u.role)?.label || u.role;
                  
                const unitName = getUnitName(u.unit_id);
                
                const isPending = u.status === 'PENDING';
                const isBlocked = u.status === 'BLOCKED';

                const canManageThisUser = currentUserRole === 'ADMIN' || (ROLE_CREATION_MAP[currentUserRole || '']?.includes(u.role));
                
                // ЛОГІКА БЛОКУВАННЯ:
                // Якщо юзер у загальному резерві (unit_id == null) і ти не Адмін — блокувати заборонено!
                const isGeneralReserve = u.unit_id == null;
                const isAdmin = currentUserRole === 'ADMIN';
                const canChangeStatus = isAdmin || !isGeneralReserve;

                return (
                  <tr key={u.id}>
                    <td>
                      <div className="user-name-cell">{u.full_name || '-'}</div>
                    </td>
                    <td>{roleLabel}</td>
                    <td>
                      {unitName ? (
                        <span className="unit-name-cell">{unitName}</span>
                      ) : u.role === 'ADMIN' ? (
                        <span style={{ color: '#6c757d', fontSize: '13px' }}>Системний персонал</span>
                      ) : u.role === 'VOLUNTEER' ? (
                        <span style={{ color: '#6c757d', fontSize: '13px' }}>Цивільний</span>
                      ) : (
                        <span className="badge badge-neutral">Резерв</span>
                      )}
                    </td>
                    <td>
                      <span className={`badge badge-${isPending ? 'warning' : isBlocked ? 'danger' : 'success'}`}>
                        {isPending ? 'Очікує реєстрації' : isBlocked ? 'Заблокований' : 'Активний'}
                      </span>
                    </td>
                    
                    {canManagePersonnel && (
                      <td>
                        {u.id !== currentUserId && canManageThisUser && (
                          <div className="action-buttons-row">
                            <button 
                              className="btn-action-small" 
                              onClick={() => openEditModal(u)}
                              disabled={isBlocked} 
                              title={isBlocked ? "Неможливо змінити звання заблокованому користувачу" : "Перемістити або змінити посаду"}
                            >
                              Змінити звання
                            </button>
                            
                            {isBlocked ? (
                              <button 
                                className="btn-action-small btn-success-small" 
                                onClick={() => setConfirmModal({ type: 'unblock', id: u.id, name: u.full_name || u.email })}
                                disabled={!canChangeStatus}
                                title={!canChangeStatus ? "Немає прав: користувач у загальному резерві" : "Повернути до активного складу"}
                              >
                                Розблокувати
                              </button>
                            ) : (
                              <button 
                                className="btn-action-small btn-danger-small" 
                                onClick={() => setConfirmModal({ type: 'block', id: u.id, name: u.full_name || u.email })}
                                disabled={!canChangeStatus}
                                title={!canChangeStatus ? "Немає прав: користувач у загальному резерві" : "Заблокувати користувача"}
                              >
                                Заблокувати
                              </button>
                            )}
                          </div>
                        )}
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