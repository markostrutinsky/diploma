import { useEffect, useState } from 'react'
import { api, type CreateUserRequest, type Unit, type User, ROLE_NAMES, UNIT_TYPE_NAMES, type UserRole } from '../api/client'
import { usePermissions } from '../hooks/usePermissions'
import { useNavigate } from 'react-router-dom'
import './AdminUsers.css'

// Ролі власника організації — мають однакові права в UI (повний доступ в межах tenant).
// TENANT_ADMIN — нова назва, ADMIN — застаріле ім'я, SYSTEM_ADMIN — платформний адмін.
const OWNER_ROLES: readonly string[] = ['ADMIN', 'TENANT_ADMIN', 'SYSTEM_ADMIN']
const isOwnerRole = (role: string | null | undefined): boolean =>
  !!role && OWNER_ROLES.includes(role)

const ROLES: { value: UserRole, label: string }[] = [
  { value: 'TENANT_ADMIN', label: ROLE_NAMES['TENANT_ADMIN'] },
  { value: 'REGION_DIRECTOR', label: ROLE_NAMES['REGION_DIRECTOR'] },
  { value: 'BRANCH_MANAGER', label: ROLE_NAMES['BRANCH_MANAGER'] },
  { value: 'DEPT_MANAGER', label: ROLE_NAMES['DEPT_MANAGER'] },
  { value: 'TEAM_LEAD', label: ROLE_NAMES['TEAM_LEAD'] },
  { value: 'REGION_LOGISTICIAN', label: ROLE_NAMES['REGION_LOGISTICIAN'] },
  { value: 'REGION_STOREKEEPER', label: ROLE_NAMES['REGION_STOREKEEPER'] },
  { value: 'BRANCH_LOGISTICIAN', label: ROLE_NAMES['BRANCH_LOGISTICIAN'] },
  { value: 'BRANCH_STOREKEEPER', label: ROLE_NAMES['BRANCH_STOREKEEPER'] },
  { value: 'DEPT_SUPERVISOR', label: ROLE_NAMES['DEPT_SUPERVISOR'] },
  { value: 'EMPLOYEE', label: ROLE_NAMES['EMPLOYEE'] },
]

const ROLE_UNIT_TYPE_MAP: Record<string, string[]> = {
  'REGION_DIRECTOR': ['REGION'],
  'REGION_LOGISTICIAN': ['REGION'],
  'REGION_STOREKEEPER': ['REGION'],
  'BRANCH_MANAGER': ['BRANCH'],
  'BRANCH_LOGISTICIAN': ['BRANCH'],
  'BRANCH_STOREKEEPER': ['BRANCH'],
  'DEPT_MANAGER': ['DEPARTMENT'],
  'DEPT_SUPERVISOR': ['DEPARTMENT'],
  'EMPLOYEE': ['REGION', 'BRANCH', 'DEPARTMENT', 'TEAM'], // Співробітник може бути в будь-якому підрозділі
  'TEAM_LEAD': ['TEAM'], // Керівник групи керує командою
  'ADMIN': [], // Адмін не прив'язаний до конкретного підрозділу
  'TENANT_ADMIN': [],
  'SYSTEM_ADMIN': [],
}

const OWNER_CREATABLE_ROLES = [
  'TENANT_ADMIN', 'REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 'TEAM_LEAD',
  'REGION_LOGISTICIAN', 'REGION_STOREKEEPER', 'BRANCH_LOGISTICIAN', 'BRANCH_STOREKEEPER', 'DEPT_SUPERVISOR', 'EMPLOYEE'
]

// SYSTEM_ADMIN свідомо відсутній як ключ: платформний адмін не створює користувачів
// через адмінку організації (для цього є TENANT_ADMIN у межах конкретного tenant).
const ROLE_CREATION_MAP: Record<string, string[]> = {
  'TENANT_ADMIN': OWNER_CREATABLE_ROLES,
  'ADMIN': OWNER_CREATABLE_ROLES,
  'REGION_DIRECTOR': [
    'BRANCH_MANAGER', 'DEPT_MANAGER', 'TEAM_LEAD',
    'REGION_LOGISTICIAN', 'REGION_STOREKEEPER', 'BRANCH_LOGISTICIAN', 'BRANCH_STOREKEEPER', 'DEPT_SUPERVISOR', 'EMPLOYEE'
  ],
  'BRANCH_MANAGER': [
    'DEPT_MANAGER', 'TEAM_LEAD',
    'BRANCH_LOGISTICIAN', 'BRANCH_STOREKEEPER', 'DEPT_SUPERVISOR', 'EMPLOYEE'
  ],
  'DEPT_MANAGER': [
    'TEAM_LEAD', 'DEPT_SUPERVISOR', 'EMPLOYEE'
  ],
  'REGION_LOGISTICIAN': [
    'REGION_STOREKEEPER', 'BRANCH_LOGISTICIAN', 'BRANCH_STOREKEEPER'
  ],
  'BRANCH_LOGISTICIAN': [
    'BRANCH_STOREKEEPER'
  ],
}

export default function AdminUsers() {
  const navigate = useNavigate()
  const [formUnits, setFormUnits] = useState<Unit[]>([])
  const [allUnits, setAllUnits] = useState<Unit[]>([])
  const [usersList, setUsersList] = useState<User[]>([])
  const [activeTab, setActiveTab] = useState<'all' | 'in_unit' | 'reserve'>('all')

  const [currentUserRole, setCurrentUserRole] = useState<string | null>(null)
  const [currentUserId, setCurrentUserId] = useState<string | null>(null)
  const [availableRoles, setAvailableRoles] = useState<typeof ROLES>([])
  
  // ---------------------------------------------------------
  // НОВИЙ СТЕЙТ ДЛЯ ПОШУКУ
  // ---------------------------------------------------------
  const [searchQuery, setSearchQuery] = useState('')

  const [showForm, setShowForm] = useState(false)
  const [editingUser, setEditingUser] = useState<User | null>(null)
  
  const [confirmModal, setConfirmModal] = useState<{ type: 'block' | 'unblock', id: string, name: string } | null>(null)
  
  const [form, setForm] = useState<Omit<CreateUserRequest, 'role'> & { role: UserRole }>({
    email: '',
    full_name: '',
    role: 'DEPT_SUPERVISOR',
  })
  
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  const perms = usePermissions()
  const canManagePersonnel = perms.can('user_invite')

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
        setForm(prev => ({ ...prev, role: filtered[0].value, unit_id: undefined }))
      }
    }
  }, [currentUserRole, editingUser])

  useEffect(() => {
    if (!form.role || !currentUserRole) return

    if (isOwnerRole(currentUserRole)) {
      const allowedTypes = ROLE_UNIT_TYPE_MAP[form.role] || [];
      
      // Фільтруємо за типом підрозділу, якщо для ролі є обмеження
      if (allowedTypes.length > 0) {
        // Примітка: переконайтеся, що ваш інтерфейс Unit має поле type (або unit_type)
        const filteredUnits = allUnits.filter(unit => allowedTypes.includes(unit.unit_type));
        setFormUnits(filteredUnits);
      } else {
        // Для ролей без конкретних обмежень (наприклад, ADMIN) або очищаємо, або віддаємо все
        setFormUnits(isOwnerRole(form.role) ? [] : allUnits);
      }
    } else {
      api.units.getMyHierarchyForRole(form.role)
        .then((data) => setFormUnits(Array.isArray(data) ? data : []))
        .catch(() => setFormUnits([]))
    }
  }, [form.role, currentUserRole, allUnits])

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
      
      const defaultRole = availableRoles.length > 0 ? availableRoles[0].value : 'DEPT_SUPERVISOR'
      setForm({ email: '', full_name: '', role: defaultRole, username: '', phone: '', unit_id: undefined })
      
      loadUsersAndUnits()
      setTimeout(() => { setShowForm(false); setSuccess(null) }, 1500)
    } catch (err: any) {
      // Перевіряємо, чи це помилка ліміту (402 Payment Required)
      if (err?.response?.status === 402 || err?.message?.includes('ліміт') || err?.message?.includes('Ліміт')) {
        const errorMsg = err?.response?.data?.error || err?.message || 'Досягнуто ліміт користувачів для вашого тарифу';
        setError(errorMsg);
        // Показуємо toast з кнопкою переходу на білінг
        setTimeout(() => {
          const shouldNavigate = window.confirm(`${errorMsg}\n\nПерейти до сторінки тарифних планів?`);
          if (shouldNavigate) {
            navigate('/billing');
          }
        }, 100);
      } else {
        setError(err?.message || 'Помилка');
      }
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
        setSuccess(`Співробітника ${confirmModal.name} переведено в резерв`)
      } else {
        await api.users.unblock(confirmModal.id)
        setSuccess(`Співробітника ${confirmModal.name} розблоковано`)
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
      role: user.role,
      unit_id: user.unit_id || undefined,
    })
    setSuccess(null)
    setError(null)
  }

  const getUnitName = (unitId?: number | null) => {
    if (!unitId) return null;
    const unit = allUnits.find(u => u.id === unitId);
    return unit ? unit.name : `ID: ${unitId}`;
  }

  // ---------------------------------------------------------
  // ЛОГІКА ФІЛЬТРАЦІЇ ПОШУКУ
  // ---------------------------------------------------------
  const filteredUsers = usersList.filter(user => {
    // Базові правила вкладок і ролей
    if (!isOwnerRole(currentUserRole) && isOwnerRole(user.role)) {
      return false;
    }
    if (activeTab === 'reserve') {
      if (user.unit_id || isOwnerRole(user.role) || user.role === 'CONTRACTOR') return false;
    }
    if (activeTab === 'in_unit' && !user.unit_id) {
      return false;
    }

    // Логіка пошуку
    if (searchQuery) {
      const query = searchQuery.toLowerCase().trim();
      const fullName = (user.full_name || '').toLowerCase();
      const email = (user.email || '').toLowerCase();
      const roleLabel = (ROLE_NAMES[user.role] || user.role).toLowerCase();
      const unitName = (getUnitName(user.unit_id) || '').toLowerCase();

      if (!fullName.includes(query) && !email.includes(query) && !roleLabel.includes(query) && !unitName.includes(query)) {
        return false;
      }
    }

    return true;
  })

  return (
    <div className="admin-users">
      <div className="page-header page-header-flex">
        <div>
          <h1>Управління персоналом</h1>
          <p className="page-subtitle">
            Керування співробітниками та призначення на посади.
          </p>
        </div>
        {canManagePersonnel && (
          <button 
            className="btn btn-primary" 
            onClick={() => {
              setSuccess(null); setError(null);
              setEditingUser(null);
              const defaultRole = availableRoles.length > 0 ? availableRoles[0].value : 'DEPT_SUPERVISOR'
              setForm({ email: '', full_name: '', role: defaultRole, unit_id: undefined })
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
              <p className="modal-description">Співробіттник отримає лист з посиланням для встановлення паролю.</p>
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
                <label>{editingUser ? 'Нова посада' : 'Посада'}</label>
                <select
                  value={form.role}
                  onChange={(e) => {
                    const newRole = e.target.value as UserRole
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
                <label>{editingUser ? 'Нова орг. одиниця' : 'Орг. одиниця (залиште пустим для Резерву)'}</label>
                <select
                  value={form.unit_id ?? ''}
                  onChange={(e) => {
                    const val = e.target.value
                    setForm({ ...form, unit_id: val ? parseInt(val, 10) : undefined })
                  }}
                  disabled={formUnits.length === 0}
                >
                  <option value="">
                    {formUnits.length === 0 ? 'Немає доступних орг. одиниць (буде в резерві)' : '-- В Кадровий резерв --'}
                  </option>
                  {formUnits.map((u) => (
                    <option key={u.id} value={u.id}>
                      {u.name} ({UNIT_TYPE_NAMES[u.unit_type] || u.unit_type})
                    </option>
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
                  disabled={loading || availableRoles.length === 0 || (!editingUser && (!form.email || !form.full_name?.trim()))}
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
                ? `Ви впевнені, що хочете перевести співробітника ${confirmModal.name} в резерв (заблокувати)? Він втратить доступ до системи.`
                : `Повернути співробітника ${confirmModal.name} до активного персоналу? Він опиниться в Кадровому резерві.`
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
        <div className="table-header-with-tabs" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '16px', marginBottom: '16px' }}>
          <h2 style={{ margin: 0 }}>Персонал</h2>
          
          {/* НОВЕ ПОЛЕ ПОШУКУ */}
          <div style={{ position: 'relative', width: '280px' }}>
            <span style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: '#94a3b8', fontSize: '14px' }}>
              🔍
            </span>
            <input
              type="text"
              className="erp-input"
              placeholder="Пошук"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              style={{ 
                paddingLeft: '35px', 
                borderRadius: '20px', 
                paddingBottom: '6px', 
                paddingTop: '6px',
                border: '1px solid #cbd5e1'
              }}
            />
            {searchQuery && (
              <button 
                onClick={() => setSearchQuery('')}
                style={{ position: 'absolute', right: '12px', top: '50%', transform: 'translateY(-50%)', background: 'none', border: 'none', cursor: 'pointer', color: '#94a3b8', padding: 0 }}
              >
                ✖
              </button>
            )}
          </div>

          <div className="inventory-tabs">
            <button className={`tab-btn ${activeTab === 'all' ? 'active' : ''}`} onClick={() => setActiveTab('all')}>Всі</button>
            <button className={`tab-btn ${activeTab === 'in_unit' ? 'active' : ''}`} onClick={() => setActiveTab('in_unit')}>В орг. одиницях</button>
            <button className={`tab-btn ${activeTab === 'reserve' ? 'active' : ''}`} onClick={() => setActiveTab('reserve')}>Кадровий резерв</button>
          </div>
        </div>

        {filteredUsers.length === 0 ? (
          <p className="empty-state">
            {searchQuery ? `За запитом "${searchQuery}" нічого не знайдено` : (activeTab === 'reserve' ? 'Кадровий резерв порожній' : 'Немає користувачів')}
          </p>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>ПІБ</th>
                <th>Посада</th>
                <th>Орг. одиниця</th>
                <th>Статус</th>
                {canManagePersonnel && <th>Дії</th>}
              </tr>
            </thead>
            <tbody>
              {filteredUsers.map(u => {
                const roleLabel = u.role === 'CONTRACTOR' 
                  ? 'Підрядник' 
                  : ROLE_NAMES[u.role] || u.role;
                  
                const unitName = getUnitName(u.unit_id);
                
                const isPending = u.status === 'PENDING';
                const isBlocked = u.status === 'BLOCKED';

                const canManageThisUser = isOwnerRole(currentUserRole) || (ROLE_CREATION_MAP[currentUserRole || '']?.includes(u.role));
                
                const isGeneralReserve = u.unit_id == null;
                const isAdmin = isOwnerRole(currentUserRole);
                const canChangeStatus = isAdmin || !isGeneralReserve;

                return (
                  <tr key={u.id}>
                    <td>
                      <div className="user-name-cell">{u.full_name || '-'}</div>
                      <div style={{fontSize: '0.8rem', color: '#64748b', marginTop: '2px'}}>{u.email}</div>
                    </td>
                    <td>{roleLabel}</td>
                    <td>
                      {unitName ? (
                        <span className="unit-name-cell">{unitName}</span>
                      ) : isOwnerRole(u.role) ? (
                        <span style={{ color: '#6c757d', fontSize: '13px' }}>Системний персонал</span>
                      ) : u.role === 'CONTRACTOR' ? (
                        <span style={{ color: '#6c757d', fontSize: '13px' }}>Зовнішній</span>
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
                              title={isBlocked ? "Неможливо змінити посаду заблокованому користувачу" : "Перемістити або змінити посаду"}
                            >
                              Змінити посаду
                            </button>
                            
                            {isBlocked ? (
                              <button 
                                className="btn-action-small btn-success-small" 
                                onClick={() => setConfirmModal({ type: 'unblock', id: u.id, name: u.full_name || u.email })}
                                disabled={!canChangeStatus}
                                title={!canChangeStatus ? "Немає прав: користувач у загальному резерві" : "Повернути до активного персоналу"}
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