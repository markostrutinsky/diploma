import { useEffect, useState } from 'react'
import { api, type ContractorRequest, type Unit } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
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
  const [resources, setResources] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [createForm, setCreateForm] = useState({ title: '', description: '' })
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

  const inventoryRoles = [
    'ADMIN', 'REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 
    'REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN', 'REGION_STOREKEEPER', 
    'BRANCH_STOREKEEPER', 'DEPT_SUPERVISOR'
  ]
  
  const canCreateRequest = inventoryRoles.includes(user?.role || '')
  const isCONTRACTOR = user?.role === 'CONTRACTOR'

  const canManageThisRequest = (requestUnitId?: number | null) => {
    if (!user || !inventoryRoles.includes(user.role)) return false;
    if (['ADMIN', 'REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'REGION_STOREKEEPER'].includes(user.role)) {
      return true;
    }
    return requestUnitId === user.unit_id;
  }

  const getPlannedQuantity = (title: string) => {
    const match = title.match(/\d+/);
    return match ? parseInt(match[0], 10) : null;
  }

  // --- РОЗУМНЕ ОЧИЩЕННЯ НАЗВИ ---
  const getCleanResourceName = (title: string) => {
    return title.replace(/^\s*\d+\s*(шт\.?|штук|пар|пари|-)?\s*/i, '').trim();
  }
  // --------------------------------------------

  const loadData = () => {
    setLoading(true)
    api.contractorRequests.list()
      .then((data) => setRequests(Array.isArray(data) ? data : []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  const loadCategories = () => {
    api.inventory.listCategories()
      .then((data) => setCategories(data || []))
      .catch(console.error)
  }

  const loadResources = () => {
    if (!isCONTRACTOR) {
      api.inventory.listResources()
        .then((data) => setResources(data || []))
        .catch(console.error)
    }
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
    loadResources()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const body = {
        ...createForm,
        unit_id: selectedUnitId ? Number(selectedUnitId) : undefined
      }

      await api.contractorRequests.create(body as any)
      setShowCreateForm(false)
      setCreateForm({ title: '', description: '' })
      setSelectedUnitId('') 
      loadData()
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка')
    }
  }

  const handleTake = async (id: string) => {
    try { await api.contractorRequests.take(id); loadData() } catch (err) { alert(err instanceof Error ? err.message : 'Помилка') }
  }

  const handleDeliver = async (id: string) => {
    try { await api.contractorRequests.deliver(id); loadData() } catch (err) { alert(err instanceof Error ? err.message : 'Помилка') }
  }

  const handleReject = async (id: string) => {
    if (!window.confirm("Ви впевнені, що хочете відхилити цю доставку (наприклад, через брак або невідповідність)?")) return;
    try { await api.contractorRequests.reject(id); loadData() } catch (err) { alert(err instanceof Error ? err.message : 'Помилка') }
  }

  const handleCancel = async (id: string) => {
    if (!window.confirm("Скасувати заявку? Вона зникне зі списку доступних завдань для підрядників.")) return;
    try { await api.contractorRequests.cancel(id); loadData() } catch (err) { alert(err instanceof Error ? err.message : 'Помилка') }
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
      loadData()
      loadResources() 
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка прийомки на склад')
    }
  }

  const myActiveTasks = requests.filter((r) => r.taken_by === user?.id && r.status === 'TAKEN')
  const displayedRequests = isCONTRACTOR ? requests.filter((r) => r.status === 'OPEN') : requests

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
                <button type="submit" className="btn btn-primary">Опублікувати</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {acceptModalId && (
        <div className="modal-overlay" onClick={() => { setAcceptModalId(null); setNameMismatchWarning(false); }}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Оприбуткування майна</h3>
            
            <div style={{ display: 'flex', gap: '15px', marginBottom: '20px', padding: '10px', background: '#f8f9fa', borderRadius: '8px' }}>
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
                <div style={{ backgroundColor: '#fff3cd', color: '#856404', padding: '12px', borderRadius: '6px', border: '1px solid #ffeeba', marginBottom: '15px', fontSize: '14px' }}>
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
                  disabled={acceptMode === 'EXISTING' && !acceptForm.resource_id}
                >
                  {nameMismatchWarning ? 'Підтвердити і прийняти' : 'Прийняти на баланс'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {isCONTRACTOR && myActiveTasks.length > 0 && (
        <div className="card my-tasks-card" style={{ marginBottom: '24px' }}>
          <h2>Мої завдання в роботі</h2>
          <ul className="request-list">
            {myActiveTasks.map((r) => (
              <li key={r.id}>
                <div className="request-info">
                  <strong>{r.title}</strong>
                  {r.unit_name && <span className="unit-badge">🏢 {r.unit_name}</span>}
                  {r.description && <p className="request-desc">{r.description}</p>}
                </div>
                <button className="btn btn-success btn-sm" onClick={() => handleDeliver(r.id)}>
                  Позначити доставленим
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}

      <div className="card">
        <h2>{isCONTRACTOR ? 'Відкриті завдання' : 'Історія та статус заявок'}</h2>
        {displayedRequests.length === 0 ? (
          <p className="empty-state">Наразі немає записів</p>
        ) : (
          <ul className="request-list">
            {displayedRequests.map((r) => (
              <li key={r.id}>
                <div className="request-info">
                  <strong>{r.title}</strong>
                  {r.unit_name && <span className="unit-badge">🏢 {r.unit_name}</span>}
                  {r.description && <p className="request-desc">{r.description}</p>}
                  
                  <div className="request-meta">
                    <span className={`badge badge-${getBadgeClass(r.status)}`}>
                      {STATUS_LABELS[r.status?.toUpperCase()] || r.status}
                    </span>
                    <span className="request-date">{new Date(r.created_at).toLocaleDateString('uk-UA')}</span>
                  </div>
                </div>

                <div className="action-buttons-row">
                  {isCONTRACTOR && r.status === 'OPEN' && (
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
        )}
      </div>
    </div>
  )
}