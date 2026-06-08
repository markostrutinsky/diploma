import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { api, type ContractorRequest, type Unit, type ContractorMembership, type ContractorMembershipStatus } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import { usePermissions } from '../hooks/usePermissions'
import Pagination from '../components/Pagination'
import './VolunteerRequests.css'

const STATUS_LABELS: Record<string, string> = {
  OPEN: 'Відкрита',
  TAKEN: 'В роботі',
  DELIVERED: 'Очікує прийомки',
  ACCEPTED: 'Прийнята на баланс',
  REJECTED: 'Відхилена',
  CANCELED: 'Скасована',
  COMPLETED: 'Виконана (старий статус)', 
}

export default function ContractorRequests() {
  const { user } = useAuth()
  const [requests, setRequests] = useState<ContractorRequest[]>([])
  const [categories, setCategories] = useState<{ id: string; name: string }[]>([])
  const [units, setUnits] = useState<Unit[]>([])
  const [warehouses, setWarehouses] = useState<{ id: string; name: string }[]>([])
  const [resources, setResources] = useState<any[]>([])
  const [memberships, setMemberships] = useState<ContractorMembership[]>([])
  const [loading, setLoading] = useState(true)
  
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [createForm, setCreateForm] = useState({ title: '', description: '', target_warehouse_id: '' })
  const [selectedUnitId, setSelectedUnitId] = useState<number | ''>('')

  const [acceptModalId, setAcceptModalId] = useState<string | null>(null)
  const [acceptMode, setAcceptMode] = useState<'NEW' | 'EXISTING'>('NEW')
  const [acceptForm, setAcceptForm] = useState({
    resource_id: '', 
    category_id: '', 
    name: '',
    quantity: 1,
    unit_type: 'PCS', 
  })
  
  const [nameMismatchWarning, setNameMismatchWarning] = useState(false)

  // ---------------------------------------------------------
  // НОВИЙ СТЕЙТ ДЛЯ ПОШУКУ
  // ---------------------------------------------------------
  const [searchQuery, setSearchQuery] = useState('')
  const [volPage, setVolPage] = useState(0)
  const VOL_PAGE_SIZE = 20

  // Ролі, що можуть приймати завдання підрядника на баланс.
  // Має збігатися з backend models.ContractorRequestCreatorRoles.
  const inventoryRoles = [
    'SYSTEM_ADMIN', 'TENANT_ADMIN', 'ADMIN', 'REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER',
    'REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN', 'REGION_STOREKEEPER',
    'BRANCH_STOREKEEPER', 'DEPT_SUPERVISOR'
  ]

  const perms = usePermissions()
  const canCreateRequest = perms.can('contractor_request_create')
  const canTakeRequest = perms.can('contractor_request_take')
  const isCONTRACTOR = perms.isContractor

  const canManageThisRequest = (requestUnitId?: number | null) => {
    if (!user || !inventoryRoles.includes(user.role)) return false;
    // Адміни та регіональні ролі приймають завдання будь-якого підрозділу.
    if (['SYSTEM_ADMIN', 'TENANT_ADMIN', 'ADMIN', 'REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'REGION_STOREKEEPER'].includes(user.role)) {
      return true;
    }
    // Решта (рівень філії/відділу) — лише завдання свого підрозділу.
    // Якщо у завдання не вказано підрозділ — теж дозволяємо (його приймуть на головний склад).
    return requestUnitId == null || requestUnitId === user.unit_id;
  }

  const getPlannedQuantity = (title: string) => {
    const match = title.match(/\d+/);
    return match ? parseInt(match[0], 10) : null;
  }

  const getCleanResourceName = (title: string) => {
    return title.replace(/^\s*\d+\s*(шт\.?|штук|пар|пари|-)?\s*/i, '').trim();
  }

  const loadData = () => {
    api.contractorRequests.list()
      .then((data) => setRequests(Array.isArray(data) ? data : []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  const loadCategories = () => {
    // Підрядник не має доступу до категорій складу (та й не потребує їх) → не смикаємо 403.
    if (isCONTRACTOR) return
    api.inventory.listCategories()
      .then((data) => setCategories(data || []))
      .catch(console.error)
  }

  const loadWarehouses = () => {
    // Підрядник не має доступу до складів організації → пропускаємо, щоб не ловити 403.
    if (isCONTRACTOR) return
    api.warehouses.list()
      .then((data) => setWarehouses(data || []))
      .catch(console.error)
  }

  const loadResources = () => {
    if (!isCONTRACTOR) {
      api.inventory.listResources()
        .then((data) => setResources(data || []))
        .catch(console.error)
    }
  }

  // Членства підрядника: з якими організаціями він уже співпрацює / подав заявку.
  // Потрібно, щоб згрупувати дошку по організаціях і показати правильну дію
  // (Співпрацювати / очікує / можна брати).
  const loadMemberships = () => {
    if (!isCONTRACTOR) return
    api.contractorMemberships.mine()
      .then((data) => setMemberships(Array.isArray(data) ? data : []))
      .catch(console.error)
  }

  useEffect(() => {
    if (user?.role === 'ADMIN') {
      api.units.list()
        .then(data => setUnits(Array.isArray(data) ? data : []))
        .catch(console.error)
    }
  }, [user?.role])

  useEffect(() => {
    loadData()
    loadCategories()
    loadWarehouses()
    loadResources()
    loadMemberships()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => { setVolPage(0) }, [searchQuery])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const body = {
        ...createForm,
        unit_id: selectedUnitId ? Number(selectedUnitId) : undefined
      }

      await api.contractorRequests.create(body)
      setShowCreateForm(false)
      setCreateForm({ title: '', description: '', target_warehouse_id: '' })
      setSelectedUnitId('') 
      toast.success('Заявку опубліковано')
      loadData()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const handleTake = async (id: string) => {
    try {
      await api.contractorRequests.take(id)
      toast.success('Завдання взято в роботу')
      loadData()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка')
    }
  }

  // Підрядник надсилає заявку на співпрацю з організацією (не беручи завдання).
  const handleApply = async (tenantId: string, tenantName?: string) => {
    const org = tenantName ? `«${tenantName}»` : 'організацією'
    try {
      const res = await api.contractorMemberships.apply(tenantId)
      if (res.status === 'APPROVED') {
        toast.success(`Співпрацю з ${org} підтверджено — можете брати завдання.`)
      } else if (res.status === 'REJECTED') {
        toast.error(`Організація ${org} відхилила вашу заявку на співпрацю.`)
      } else {
        toast.success(`Заявку на співпрацю з ${org} надіслано. Очікуйте підтвердження адміністратора.`)
      }
      loadMemberships()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Не вдалося надіслати заявку')
    }
  }

  const handleDeliver = async (id: string) => {
    try {
      await api.contractorRequests.deliver(id)
      toast.success('Завдання позначено доставленим')
      loadData()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const handleReject = async (id: string) => {
    if (!window.confirm("Ви впевнені, що хочете відхилити цю доставку (наприклад, через брак або невідповідність)?")) return;
    try {
      await api.contractorRequests.reject(id)
      toast.success('Доставку відхилено')
      loadData()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const handleCancel = async (id: string) => {
    if (!window.confirm("Скасувати заявку? Вона зникне зі списку доступних завдань для підрядників.")) return;
    try {
      await api.contractorRequests.cancel(id)
      toast.success('Заявку скасовано')
      loadData()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const submitAccept = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!acceptModalId) return

    const originalRequest = requests.find(r => r.id === acceptModalId);
    if (originalRequest) {
      const planned = getPlannedQuantity(originalRequest.title);
      if (planned !== null && planned !== acceptForm.quantity) {
        const confirmText = acceptForm.quantity > planned 
          ? `Ви вводите більше майна (${acceptForm.quantity}), ніж було в заявці (${planned}). Поставити надлишок на баланс?`
          : `Ви вводите менше майна (${acceptForm.quantity}), ніж було в заявці (${planned}). Продовжити?`;
        
        if (!window.confirm(confirmText)) return;
      }

      if (acceptMode === 'EXISTING' && !nameMismatchWarning) {
        const reqTitle = originalRequest.title.toLowerCase(); 
        const selectedName = acceptForm.name.toLowerCase();    
        
        const reqWords = reqTitle.split(/\s+/).filter(w => w.length > 3);
        const hasCommonWord = reqWords.some(w => selectedName.includes(w));

        if (!hasCommonWord) {
          setNameMismatchWarning(true); 
          return; 
        }
      }
    }

    try {
      const payload = {
        ...acceptForm,
        resource_id: acceptMode === 'EXISTING' && acceptForm.resource_id ? String(acceptForm.resource_id) : undefined,
      }

      await api.contractorRequests.accept(acceptModalId, payload)

      setAcceptModalId(null)
      setNameMismatchWarning(false) 
      toast.success('Майно прийнято на баланс')
      loadData()
      loadResources() 
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка прийомки на склад')
    }
  }

  // ---------------------------------------------------------
  // ЛОГІКА ФІЛЬТРАЦІЇ ПОШУКУ
  // ---------------------------------------------------------
  const filteredRequests = requests.filter(r => {
    if (!searchQuery) return true;
    const query = searchQuery.toLowerCase().trim();
    const title = (r.title || '').toLowerCase();
    const desc = (r.description || '').toLowerCase();
    const unitName = (r.unit_name || '').toLowerCase();
    
    return title.includes(query) || desc.includes(query) || unitName.includes(query);
  });

  const myActiveTasks = filteredRequests.filter((r) => r.taken_by === user?.id && r.status === 'TAKEN');
  const displayedRequests = isCONTRACTOR ? filteredRequests.filter((r) => r.status === 'OPEN') : filteredRequests;

  const volTotalPages = Math.max(1, Math.ceil(displayedRequests.length / VOL_PAGE_SIZE));
  const safeVolPage = Math.min(volPage, volTotalPages - 1);
  const pagedVolRequests = displayedRequests.slice(safeVolPage * VOL_PAGE_SIZE, (safeVolPage + 1) * VOL_PAGE_SIZE);

  const getBadgeClass = (status: string) => {
    switch(status?.toUpperCase()) {
      case 'ACCEPTED': 
      case 'COMPLETED': return 'success'
      case 'DELIVERED': return 'info'
      case 'OPEN': return 'warning'
      case 'REJECTED': 
      case 'CANCELED': return 'danger'
      default: return 'neutral'
    }
  }

  // Статус співпраці підрядника з кожною організацією (tenant_id → статус).
  const membershipByTenant: Record<string, ContractorMembershipStatus> = {}
  for (const m of memberships) membershipByTenant[m.tenant_id] = m.status

  // Групуємо відкриті завдання за організацією-замовником. Так підрядник бачить дошку
  // згрупованою, а кнопка «Співпрацювати» з'являється ОДИН раз на організацію (а не біля
  // кожної заявки), навіть якщо організація опублікувала десятки потреб.
  const openGroups: { tenantId: string; tenantName: string; items: ContractorRequest[] }[] = []
  if (isCONTRACTOR) {
    const indexByTenant: Record<string, number> = {}
    for (const r of displayedRequests) {
      const tid = r.tenant_id || 'unknown'
      if (indexByTenant[tid] === undefined) {
        indexByTenant[tid] = openGroups.length
        openGroups.push({ tenantId: tid, tenantName: r.tenant_name || 'Організація', items: [] })
      }
      openGroups[indexByTenant[tid]].items.push(r)
    }
  }

  // Дія/статус співпраці на рівні організації (показуємо в заголовку групи).
  const renderCollabAction = (tenantId: string, tenantName: string) => {
    const status = membershipByTenant[tenantId]
    if (status === 'APPROVED') {
      return <span className="collab-badge collab-approved">✓ Ви співпрацюєте</span>
    }
    if (status === 'PENDING') {
      return <span className="collab-badge collab-pending">⏳ Заявку надіслано · очікує підтвердження</span>
    }
    if (status === 'REJECTED') {
      return <span className="collab-badge collab-rejected">✕ Співпрацю відхилено</span>
    }
    return (
      <button className="btn btn-primary btn-sm collab-btn" onClick={() => handleApply(tenantId, tenantName)}>
        🤝 Співпрацювати
      </button>
    )
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
    <div className="contractor-requests-page">
      <div className="page-header">
        <h1>{isCONTRACTOR ? 'Доступні завдання' : 'Заявки підрядникам'}</h1>
        {canCreateRequest && (
          <button className="btn btn-primary" onClick={() => setShowCreateForm(true)}>
            + Створити заявку
          </button>
        )}
      </div>

      {/* НОВИЙ БЛОК ПОШУКУ */}
      <div className="filters-bar" style={{ marginBottom: '24px', backgroundColor: 'var(--bg-input)', padding: '16px', borderRadius: '8px', border: '1px solid var(--border)' }}>
        <div style={{ position: 'relative', width: '100%', maxWidth: '400px' }}>
          <span style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: '#94a3b8', fontSize: '14px' }}>
            🔍
          </span>
          <input
            type="text"
            className="erp-input"
            placeholder="Пошук"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            style={{ paddingLeft: '35px', width: '100%' }}
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
      </div>

      {showCreateForm && (
        <div className="modal-overlay" onClick={() => setShowCreateForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Нове завдання для підрядника</h3>
            <form onSubmit={handleCreate}>
              <div className="form-group">
                <label>Назва потреби</label>
                <input
                  value={createForm.title}
                  onChange={(e) => setCreateForm({ ...createForm, title: e.target.value })}
                  placeholder="Наприклад: 5 ноутбуків Dell, 10 офісних крісел..."
                  required
                />
              </div>
              <div className="form-group">
                <label>Опис (Деталі)</label>
                <textarea
                  value={createForm.description}
                  onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
                  placeholder="Вкажіть точну модель, бажані характеристики або терміни постачання..."
                  rows={4}
                />
              </div>

              <div className="form-group" style={{ marginTop: '15px' }}>
                <label>📍 Склад призначення (куди доставити)</label>
                <select 
                  value={createForm.target_warehouse_id} 
                  onChange={(e) => setCreateForm({ ...createForm, target_warehouse_id: e.target.value })}
                  style={{ width: '100%', padding: '8px', borderRadius: '4px', border: '1px solid #ced4da' }}
                >
                  <option value="">-- Оберіть склад --</option>
                  {warehouses.map(w => (
                    <option key={w.id} value={w.id}>{w.name}</option>
                  ))}
                </select>
                <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginTop: '6px', marginBottom: 0 }}>
                  Підрядник побачить цей склад у деталях завдання
                </p>
              </div>
              
              {user?.role === 'ADMIN' && (
                <div className="form-group" style={{ marginTop: '15px' }}>
                  <label>Для якої орг. одиниці замовляємо?</label>
                  <select 
                    value={selectedUnitId} 
                    onChange={(e) => setSelectedUnitId(e.target.value ? Number(e.target.value) : '')}
                    style={{ width: '100%', padding: '8px', borderRadius: '4px', border: '1px solid #ced4da' }}
                  >
                    <option value="">-- Центральний офіс (без прив'язки) --</option>
                    {units.map(u => (
                      <option key={u.id} value={u.id}>{u.name}</option>
                    ))}
                  </select>
                </div>
              )}

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowCreateForm(false)}>
                  Скасувати
                </button>
                <button type="submit" className="btn btn-primary" disabled={!createForm.title?.trim()}>Опублікувати</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {acceptModalId && (
        <div className="modal-overlay" onClick={() => { setAcceptModalId(null); setNameMismatchWarning(false); }}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Оприбуткування майна</h3>
            
            <div style={{ display: 'flex', gap: '15px', marginBottom: '20px', padding: '10px', background: 'var(--bg-input)', borderRadius: '8px' }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer', margin: 0 }}>
                <input 
                  type="radio" 
                  checked={acceptMode === 'NEW'} 
                  onChange={() => { setAcceptMode('NEW'); setNameMismatchWarning(false); }} 
                />
                Створити нову картку товару
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer', margin: 0 }}>
                <input 
                  type="radio" 
                  checked={acceptMode === 'EXISTING'} 
                  onChange={() => setAcceptMode('EXISTING')} 
                />
                Додати до існуючого майна
              </label>
            </div>

            <form onSubmit={submitAccept}>
              {acceptMode === 'EXISTING' ? (
                <div className="form-group">
                  <label>Оберіть майно з вашого складу</label>
                  <select 
                    value={acceptForm.resource_id} 
                    onChange={(e) => {
                      setNameMismatchWarning(false); 
                      const res = resources.find(r => String(r.id) === e.target.value);
                      setAcceptForm({
                        ...acceptForm, 
                        resource_id: e.target.value,
                        name: res?.name || '',
                        category_id: res?.category_id || '',
                        unit_type: res?.unit_type || 'PCS'
                      });
                    }}
                    required
                  >
                    <option value="" disabled>-- Оберіть ресурс зі складу --</option>
                    {resources.map(r => (
                      <option key={r.id} value={r.id}>
                        {r.name} (Зараз на складі: {r.quantity} {r.unit_type})
                      </option>
                    ))}
                  </select>
                  {resources.length === 0 && (
                    <p style={{ color: '#dc3545', fontSize: '12px', marginTop: '5px' }}>Ваш склад порожній.</p>
                  )}
                </div>
              ) : (
                <>
                  <div className="form-group">
                    <label>Категорія майна</label>
                    <select 
                      value={acceptForm.category_id} 
                      onChange={(e) => setAcceptForm({...acceptForm, category_id: e.target.value})} 
                      required
                    >
                      <option value="" disabled>Оберіть категорію</option>
                      {categories.map(c => (
                        <option key={c.id} value={c.id}>{c.name}</option>
                      ))}
                    </select>
                  </div>
                  <div className="form-group">
                    <label>Точна номенклатурна назва</label>
                    <input 
                      value={acceptForm.name} 
                      onChange={(e) => setAcceptForm({ ...acceptForm, name: e.target.value })} 
                      required 
                    />
                    <small style={{ color: '#6c757d', fontSize: '12px', display: 'block', marginTop: '4px' }}>
                      ⚠️ Вказуйте назву без кількості (правильно: "Ноутбук Dell", неправильно: "5 ноутбуків").
                    </small>
                  </div>
                </>
              )}

              <div className="form-row-2">
                <div className="form-group">
                  <label>
                    Кількість
                    <span style={{ color: '#6c757d', fontWeight: 400, marginLeft: '6px', fontSize: '12px' }}>
                      (у заявці: {getPlannedQuantity(requests.find(r => r.id === acceptModalId)?.title || "") || '?'})
                    </span>
                  </label>
                  <input 
                    type="number" 
                    min="1" 
                    value={acceptForm.quantity} 
                    onChange={(e) => setAcceptForm({ ...acceptForm, quantity: Number(e.target.value) })} 
                    required 
                  />
                </div>
                {acceptMode === 'NEW' && (
                  <div className="form-group">
                    <label>Од. виміру</label>
                    <select 
                      value={acceptForm.unit_type} 
                      onChange={(e) => setAcceptForm({ ...acceptForm, unit_type: e.target.value })}
                    >
                      <option value="PCS">шт</option>
                      <option value="KIT">комплект</option>
                      <option value="KG">кг</option>
                      <option value="L">літр</option>
                    </select>
                  </div>
                )}
              </div>

              {nameMismatchWarning && (
                <div style={{ backgroundColor: 'rgba(245, 158, 11, 0.15)', color: '#856404', padding: '12px', borderRadius: '6px', border: '1px solid rgba(245, 158, 11, 0.3)', marginBottom: '15px', fontSize: '14px' }}>
                  ⚠️ <strong>Увага, можлива помилка!</strong><br/>
                  Ви приймаєте заявку на <strong>"{requests.find(r => r.id === acceptModalId)?.title}"</strong>, 
                  але обрали на складі <strong>"{acceptForm.name}"</strong>. <br/>
                  Якщо ви впевнені, натисніть "Прийняти на баланс" ще раз.
                </div>
              )}

              <div className="modal-actions">
                <button 
                  type="button" 
                  className="btn btn-secondary" 
                  onClick={() => { setAcceptModalId(null); setNameMismatchWarning(false); }}
                >
                  Скасувати
                </button>
                <button 
                  type="submit" 
                  className="btn btn-primary" 
                  disabled={
                    (acceptMode === 'EXISTING' && !acceptForm.resource_id) ||
                    (acceptMode === 'NEW' && (!acceptForm.category_id || !acceptForm.name?.trim()))
                  }
                >
                  {nameMismatchWarning ? 'Підтвердити і прийняти' : 'Прийняти на баланс'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {isCONTRACTOR && (
        <div className="card my-tasks-card" style={{ marginBottom: '24px' }}>
          <h2>Мої завдання в роботі</h2>
          {myActiveTasks.length === 0 ? (
            <p className="empty-state">
              {searchQuery ? `За запитом "${searchQuery}" активних завдань не знайдено` : 'У вас немає завдань в роботі'}
            </p>
          ) : (
            <ul className="request-list">
              {myActiveTasks.map((r) => (
                <li key={r.id}>
                  <div className="request-info">
                    <strong>{r.title}</strong>
                    <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap', marginTop: '4px' }}>
                      {isCONTRACTOR && r.tenant_name && <span className="unit-badge org-badge">🏛️ {r.tenant_name}</span>}
                      {r.unit_name && <span className="unit-badge">🏢 {r.unit_name}</span>}
                      {r.warehouse_name && <span className="unit-badge" style={{ background: 'var(--info-bg)', color: 'var(--info)' }}>📍 {r.warehouse_name}</span>}
                    </div>
                    {r.description && <p className="request-desc">{r.description}</p>}
                  </div>
                  <button className="btn btn-success btn-sm" onClick={() => handleDeliver(r.id)}>
                    Позначити доставленим
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}

      {isCONTRACTOR ? (
        <div className="card">
          <h2>Відкриті завдання</h2>
          {displayedRequests.length === 0 ? (
            <p className="empty-state">
              {searchQuery ? `За запитом "${searchQuery}" нічого не знайдено` : 'Наразі немає відкритих завдань'}
            </p>
          ) : (
            <div className="org-groups">
              {openGroups.map((g) => {
                const status = membershipByTenant[g.tenantId]
                const approved = status === 'APPROVED'
                return (
                  <div className="org-group" key={g.tenantId}>
                    <div className="org-group-header">
                      <div className="org-group-title">
                        <span className="org-group-name">🏛️ {g.tenantName}</span>
                        <span className="org-group-count">завдань: {g.items.length}</span>
                      </div>
                      {renderCollabAction(g.tenantId, g.tenantName)}
                    </div>

                    {!approved && (
                      <p className="org-group-hint">
                        {status === 'PENDING'
                          ? '⏳ Щойно адміністратор підтвердить співпрацю — ви зможете брати завдання цієї організації.'
                          : status === 'REJECTED'
                            ? '✕ Організація відхилила вашу заявку, тож її завдання поки недоступні.'
                            : 'Натисніть «Співпрацювати», щоб подати заявку й отримати доступ до завдань цієї організації.'}
                      </p>
                    )}

                    <ul className="request-list">
                      {g.items.map((r) => (
                        <li key={r.id} className={approved ? '' : 'request-locked'}>
                          <div className="request-info">
                            <strong>{r.title}</strong>
                            <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap', marginTop: '4px' }}>
                              {r.unit_name && <span className="unit-badge">🏢 {r.unit_name}</span>}
                              {r.warehouse_name && <span className="unit-badge" style={{ background: 'var(--info-bg)', color: 'var(--info)' }}>📍 {r.warehouse_name}</span>}
                            </div>
                            {r.description && <p className="request-desc">{r.description}</p>}
                            <div className="request-meta">
                              <span className="badge badge-warning">{STATUS_LABELS['OPEN']}</span>
                              <span className="request-date">{new Date(r.created_at).toLocaleDateString('uk-UA')}</span>
                            </div>
                          </div>
                          <div className="action-buttons-row">
                            {approved ? (
                              <button className="btn btn-primary btn-sm take-btn" onClick={() => handleTake(r.id)}>
                                Взяти в роботу
                              </button>
                            ) : (
                              <span className="lock-hint">🔒 Потрібне підтвердження</span>
                            )}
                          </div>
                        </li>
                      ))}
                    </ul>
                  </div>
                )
              })}
            </div>
          )}
        </div>
      ) : (
      <div className="card">
        <h2>Історія та статус заявок</h2>
        {displayedRequests.length === 0 ? (
          <p className="empty-state">
            {searchQuery ? `За запитом "${searchQuery}" нічого не знайдено` : 'Наразі немає записів'}
          </p>
        ) : (
          <>
          <ul className="request-list">
            {pagedVolRequests.map((r) => (
              <li key={r.id}>
                <div className="request-info">
                  <strong>{r.title}</strong>
                  <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap', marginTop: '4px' }}>
                    {isCONTRACTOR && r.tenant_name && <span className="unit-badge org-badge">🏛️ {r.tenant_name}</span>}
                    {r.unit_name && <span className="unit-badge">🏢 {r.unit_name}</span>}
                    {r.warehouse_name && <span className="unit-badge" style={{ background: 'var(--info-bg)', color: 'var(--info)' }}>📍 {r.warehouse_name}</span>}
                  </div>
                  {r.description && <p className="request-desc">{r.description}</p>}
                  
                  <div className="request-meta">
                    <span className={`badge badge-${getBadgeClass(r.status)}`}>
                      {STATUS_LABELS[r.status?.toUpperCase()] || r.status}
                    </span>
                    <span className="request-date">{new Date(r.created_at).toLocaleDateString('uk-UA')}</span>
                  </div>
                </div>

                <div className="action-buttons-row">
                  {canTakeRequest && r.status === 'OPEN' && (
                    <button className="btn btn-primary btn-sm take-btn" onClick={() => handleTake(r.id)}>
                      Взяти в роботу
                    </button>
                  )}

                  {canManageThisRequest(r.unit_id) && r.status === 'DELIVERED' && (
                    <>
                      <button 
                        className="btn-action-small btn-success-small" 
                        onClick={() => {
                          setNameMismatchWarning(false)
                          
                          const cleanTitle = getCleanResourceName(r.title);
                          const reqWords = cleanTitle.toLowerCase().split(/\s+/).filter(w => w.length > 3);
                          
                          const matchedResource = resources.find(res => {
                            const resNameLower = res.name.toLowerCase();
                            return reqWords.some(w => resNameLower.includes(w));
                          });

                          if (matchedResource) {
                            setAcceptMode('EXISTING')
                            setAcceptForm(prev => ({ 
                              ...prev, 
                              resource_id: String(matchedResource.id),
                              name: matchedResource.name,
                              quantity: getPlannedQuantity(r.title) || 1,
                              category_id: matchedResource.category_id || '',
                              unit_type: matchedResource.unit_type || 'PCS'
                            }))
                          } else {
                            setAcceptMode('NEW')
                            setAcceptForm(prev => ({ 
                              ...prev, 
                              resource_id: '',
                              name: cleanTitle, 
                              quantity: getPlannedQuantity(r.title) || 1,
                              category_id: categories.length > 0 ? categories[0].id : ''
                            }))
                          }

                          setAcceptModalId(r.id)
                        }}
                      >
                        Прийняти на баланс
                      </button>
                      <button className="btn-action-small btn-danger-small" onClick={() => handleReject(r.id)}>
                        Відхилити
                      </button>
                    </>
                  )}

                  {canManageThisRequest(r.unit_id) && r.status === 'OPEN' && (
                    <button className="btn-action-small btn-secondary" onClick={() => handleCancel(r.id)}>
                      Скасувати запит
                    </button>
                  )}
                </div>
              </li>
            ))}
          </ul>
          <Pagination
            currentPage={safeVolPage}
            totalPages={volTotalPages}
            onPageChange={setVolPage}
            totalItems={displayedRequests.length}
            itemsPerPage={VOL_PAGE_SIZE}
          />
          </>
        )}
      </div>
      )}
    </div>
  )
}