import React, { useEffect, useState, useMemo } from 'react'
import { api, type SupplyRequest, type Resource, type Vehicle, type Warehouse, type User, type Unit } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import toast from 'react-hot-toast'
import './Requests.css'

export default function Requests() {
  const { user } = useAuth()
  const [requests, setRequests] = useState<SupplyRequest[]>([])
  const [resources, setResources] = useState<Resource[]>([])
  const [warehouses, setWarehouses] = useState<Warehouse[]>([])
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [units, setUnits] = useState<Unit[]>([]) // НОВЕ: Список підрозділів для ієрархії
  
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  
  const [newReq, setNewReq] = useState({ resource_id: '', quantity: 1, target_warehouse_id: '' })

  const [filterStatus, setFilterStatus] = useState<string>('ALL')
  const [filterWarehouseId, setFilterWarehouseId] = useState<string>('ALL')

  const [selectedReqIds, setSelectedReqIds] = useState<Set<string>>(new Set())
  const [showDispatchModal, setShowDispatchModal] = useState(false)
  const [dispatchForm, setDispatchForm] = useState({
    from_warehouse_id: '',
    to_warehouse_id: '',
    vehicle_id: '',
    priority: 'NORMAL'
  })

  const canCreate = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'PLATOON_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST', 'COMPANY_SERGEANT'].includes(user?.role || '')
  const canApprove = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST'].includes(user?.role || '')

  const loadData = async () => {
    setLoading(true)
    try {
      const token = localStorage.getItem('token')
      const [reqs, resRes, whs, vehs, usersRes, unitsRes] = await Promise.all([
        api.requests.list().catch(() => []),
        api.inventory.listResources(undefined).catch(() => []),
        api.warehouses.list().catch(() => []),
        (api as any).vehicles?.list().catch(() => fetch('/api/vehicles', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json()).catch(() => [])) || [],
        api.users.getVisible().catch(() => fetch('/api/users', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json()).catch(() => [])),
        api.units.list().catch(() => []) // Завантажуємо підрозділи
      ])
      
      setRequests(Array.isArray(reqs) ? reqs : [])
      setResources(Array.isArray(resRes) ? resRes : [])
      const whsArray = Array.isArray(whs) ? whs : []
      setWarehouses(whsArray)
      setVehicles(Array.isArray(vehs) ? vehs : [])
      setUsers(Array.isArray(usersRes) ? usersRes : [])
      setUnits(Array.isArray(unitsRes) ? unitsRes : [])
      
      if (Array.isArray(resRes) && resRes.length > 0 && !newReq.resource_id) {
        setNewReq(prev => ({ ...prev, resource_id: resRes[0].id }))
      }
      if (whsArray.length > 0 && !newReq.target_warehouse_id) {
        setNewReq(prev => ({ ...prev, target_warehouse_id: whsArray[0].id }))
      }
    } catch (error) { console.error(error) } finally { setLoading(false) }
  }

  useEffect(() => { loadData() }, [showForm])

  const activeTargetWarehouseId = selectedReqIds.size > 0 
    ? requests.find(r => r.id === Array.from(selectedReqIds)[0])?.target_warehouse_id 
    : null;

  const toggleSelection = (id: string) => {
    const newSet = new Set(selectedReqIds)
    if (newSet.has(id)) newSet.delete(id)
    else newSet.add(id)
    setSelectedReqIds(newSet)
  }

  const filteredRequests = requests.filter(r => {
    const matchStatus = filterStatus === 'ALL' || r.status === filterStatus;
    const matchWarehouse = filterWarehouseId === 'ALL' || r.target_warehouse_id === filterWarehouseId;
    return matchStatus && matchWarehouse;
  })

  const selectedRequestsDetails = requests.filter(r => selectedReqIds.has(r.id))
  const currentTotalWeight = useMemo(() => {
    return selectedRequestsDetails.reduce((sum, req) => {
      const resource = resources.find(res => res.id === req.resource_id)
      return sum + ((resource?.weight_kg || 1) * req.quantity)
    }, 0)
  }, [selectedRequestsDetails, resources])

  const selectedVehicle = vehicles.find(v => v.id === dispatchForm.vehicle_id)
  const isOverweight = selectedVehicle ? currentTotalWeight > selectedVehicle.capacity_kg : false
  const fillPercentage = selectedVehicle ? Math.min(100, (currentTotalWeight / selectedVehicle.capacity_kg) * 100) : 0
  let barStatusClass = fillPercentage >= 100 ? 'bar-critical' : fillPercentage > 80 ? 'bar-warning' : 'bar-safe' 

  // --- ЛОГІКА ІЄРАРХІЇ: Доступні склади для відправки ---
  // --- ЛОГІКА ІЄРАРХІЇ ТА НАЯВНОСТІ: Доступні склади для відправки ---
  const allowedSourceWarehouses = useMemo(() => {
    if (!dispatchForm.to_warehouse_id || units.length === 0) return [];

    const targetWarehouse = warehouses.find(w => w.id === dispatchForm.to_warehouse_id);
    if (!targetWarehouse) return [];

    const targetUnit = units.find(u => u.id === targetWarehouse.unit_id);
    if (!targetUnit) return [];

    // 1. ФІЛЬТРАЦІЯ ПО ІЄРАРХІЇ
    const allowedUnitIds = new Set<number>();
    allowedUnitIds.add(targetUnit.id); 

    let currentParentId = targetUnit.parent_id;
    let depth = 0; 
    while (currentParentId && depth < 20) {
      allowedUnitIds.add(currentParentId);
      const parentNode = units.find(u => u.id === currentParentId);
      currentParentId = parentNode?.parent_id;
      depth++;
    }

    const hierarchicallyAllowed = warehouses.filter(w => allowedUnitIds.has(w.unit_id) && w.id !== dispatchForm.to_warehouse_id);

    // 2. НОВА ЛОГІКА: ФІЛЬТРАЦІЯ ПО НАЯВНОСТІ МАЙНА
    
    // Крок А: Рахуємо, скільки і чого (за назвою) нам треба для цього збірного рейсу
    const requiredItems = selectedRequestsDetails.reduce((acc, req) => {
      const res = resources.find(r => r.id === req.resource_id);
      if (res && res.name) {
        acc[res.name] = (acc[res.name] || 0) + req.quantity;
      }
      return acc;
    }, {} as Record<string, number>);

    // Крок Б: Залишаємо тільки ті склади, де є ВСІ потрібні ресурси в достатній кількості
    return hierarchicallyAllowed.filter(w => {
      for (const [name, neededQty] of Object.entries(requiredItems)) {
        // Рахуємо сумарний залишок ресурсу з такою назвою саме на ЦЬОМУ складі
        const availableQty = resources
          .filter(r => r.warehouse_id === w.id && r.name === name)
          .reduce((sum, r) => sum + r.quantity, 0);

        // Якщо хоча б одного ресурсу не вистачає — відкидаємо цей склад
        if (availableQty < neededQty) {
          return false;
        }
      }
      return true; // Склад має все необхідне!
    });

  }, [dispatchForm.to_warehouse_id, warehouses, units, selectedRequestsDetails, resources]);

  const handleOpenDispatchModal = () => {
    setDispatchForm(prev => ({
      ...prev,
      to_warehouse_id: activeTargetWarehouseId || '',
      from_warehouse_id: '' // Скидаємо попередній вибір
    }))
    setShowDispatchModal(true)
  }

  const handleDispatchSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (selectedReqIds.size === 0) return toast.error('Не вибрано жодної заявки!')
    if (!selectedVehicle) return toast.error('Оберіть транспорт!')
    if (!dispatchForm.from_warehouse_id) return toast.error('Оберіть склад відправник!')
    if (isOverweight) return toast.error(`Перевантаження! Максимум ${selectedVehicle.capacity_kg} кг`)

    const toastId = 'dispatch_toast'
    toast.loading('Формуємо збірний рейс...', { id: toastId })

    try {
      const token = localStorage.getItem('token')
      const payloadItems = selectedRequestsDetails.map(req => ({
        resource_id: req.resource_id,
        quantity: req.quantity,
        request_id: req.id 
      }))

      const payload = {
        from_warehouse_id: dispatchForm.from_warehouse_id,
        to_warehouse_id: dispatchForm.to_warehouse_id,
        vehicle_id: dispatchForm.vehicle_id,
        priority: dispatchForm.priority,
        items: payloadItems 
      }

      const response = await fetch('/api/inventory/shipments', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify(payload)
      })

      if (!response.ok) throw new Error((await response.json().catch(() => ({}))).error || 'Помилка сервера')

      toast.success(`🚚 Збірний рейс відправлено (${payloadItems.length} позицій)!`, { id: toastId, duration: 4000 })
      setShowDispatchModal(false)
      setSelectedReqIds(new Set()) 
      loadData() 
    } catch (error: any) { 
      toast.error(error.message || 'Не вдалося створити рейс', { id: toastId, duration: 5000 }) 
    }
  }

  const handleCreate = async (e: React.FormEvent) => { 
    e.preventDefault(); 
    if (!newReq.target_warehouse_id) {
      toast.error("❌ Оберіть цільовий склад!", { duration: 5000 })
      return
    }

    try { 
      await api.requests.create(newReq as any); 
      setShowForm(false); 
      setNewReq({ resource_id: resources[0]?.id || '', quantity: 1, target_warehouse_id: warehouses[0]?.id || '' }); 
      loadData(); 
      toast.success('Заявку створено!') 
    } catch (err) { 
      toast.error(err instanceof Error ? err.message : 'Помилка') 
    } 
  }

  const handleApprove = async (id: string, approved: boolean) => { 
    let comment = ""; 
    if (!approved) { 
      const reason = window.prompt("Вкажіть причину відхилення:"); 
      if (reason === null) return; comment = reason; 
    } 
    try { 
      await (api.requests.approve as any)(id, approved, comment); 
      toast.success(approved ? 'Заявку погоджено!' : 'Заявку відхилено.'); 
      loadData(); 
    } catch (err) { toast.error(err instanceof Error ? err.message : 'Помилка') } 
  }

  const statusLabel: Record<string, string> = { PENDING: 'Очікує', APPROVED: 'Затверджено', DISPATCHED: 'В дорозі', REJECTED: 'Відхилено', COMPLETED: 'Виконано', OPEN: 'Відкрито' }
  const availableVehicles = vehicles.filter(v => v.status === 'ACTIVE' && (v.type === 'VAN' || v.type === 'TRUCK' || v.type === 'PICKUP'))

  if (loading) return <div className="page-loading"><div className="spinner" /></div>

  return (
    <div className="requests-page">
      <div className="page-header header-flex">
        <div className="header-titles">
          <h1>Заявки на постачання</h1>
          <p className="header-subtitle">Управління потребами складів</p>
        </div>
        <div className="header-actions">
          {selectedReqIds.size > 0 && canApprove && (
            <button className="btn btn-dispatch" onClick={handleOpenDispatchModal}>
              🚚 Сформувати рейс на цей склад ({selectedReqIds.size})
            </button>
          )}
          {canCreate && <button className="btn btn-primary" onClick={() => setShowForm(true)}>+ Нова заявка</button>}
        </div>
      </div>

      <div className="filters-bar" style={{ display: 'flex', gap: '16px', marginBottom: '24px', backgroundColor: '#f8fafc', padding: '16px', borderRadius: '8px', border: '1px solid #e2e8f0' }}>
        <div style={{ flex: 1 }}>
          <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: '#64748b', marginBottom: '4px' }}>Статус заявки</label>
          <select className="erp-input" value={filterStatus} onChange={e => setFilterStatus(e.target.value)}>
            <option value="ALL">Всі статуси</option>
            <option value="PENDING">⏳ Очікують погодження</option>
            <option value="APPROVED">📦 Затверджені (Очікують логістику)</option>
            <option value="DISPATCHED">🚛 В дорозі (Прямують на склад)</option> {/* <--- ДОДАНО */}
            <option value="COMPLETED">✅ Доставлені на склад</option>
          </select>
        </div>
        <div style={{ flex: 2 }}>
          <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: '#64748b', marginBottom: '4px' }}>Цільовий Склад</label>
          <select className="erp-input" value={filterWarehouseId} onChange={e => setFilterWarehouseId(e.target.value)}>
            <option value="ALL">Всі склади</option>
            {warehouses.map(w => <option key={w.id} value={w.id}>{w.name}</option>)}
          </select>
        </div>
      </div>

      {showDispatchModal && (
        <div className="modal-overlay" onClick={() => setShowDispatchModal(false)}>
          <div className="modal modal-wide" onClick={(e) => e.stopPropagation()}>
            <h3 className="modal-title">🚚 Логістика: Рейс на склад</h3>
            
            <div className="dispatch-summary">
              <div className="summary-title">Вантаж до відправки ({selectedReqIds.size} позицій):</div>
              <ul className="summary-list">
                {selectedRequestsDetails.map(req => {
                  const res = resources.find(r => r.id === req.resource_id)
                  const itemWeight = (res?.weight_kg || 1) * req.quantity
                  return (
                    <li key={req.id} className="summary-item">
                      <span className="item-name">📦 {res?.name || 'Невідомий ресурс'} — {req.quantity} шт.</span>
                      <span className="item-weight">~{itemWeight.toFixed(1)} кг</span>
                    </li>
                  )
                })}
              </ul>
              <div className="summary-total">Загальна розрахункова вага: {currentTotalWeight.toFixed(1)} кг</div>
            </div>

            <form onSubmit={handleDispatchSubmit}>
              <div className="form-row-2 gap-16 mb-16">
                <div className="form-group flex-1 mb-0">
                  <label>Звідки (Склад відправник)</label>
                  <select className="erp-input" value={dispatchForm.from_warehouse_id} onChange={e => setDispatchForm({...dispatchForm, from_warehouse_id: e.target.value})} required>
                    <option value="" disabled>Оберіть склад...</option>
                    {allowedSourceWarehouses.map(w => {
                      const u = units.find(unit => unit.id === w.unit_id);
                      return <option key={w.id} value={w.id}>{w.name} ({u?.name})</option>
                    })}
                  </select>
                  {allowedSourceWarehouses.length === 0 && <span className="error-text" style={{marginTop: '4px'}}>Немає доступних складів вище по ієрархії!</span>}
                </div>
                <div className="form-group flex-1 mb-0">
                  <label>Куди (Заблоковано системою)</label>
                  <select className="erp-input" value={dispatchForm.to_warehouse_id} disabled style={{ backgroundColor: '#e2e8f0', color: '#475569', cursor: 'not-allowed' }}>
                    {warehouses.map(w => <option key={w.id} value={w.id}>{w.name}</option>)}
                  </select>
                  <span style={{ fontSize: '11px', color: '#64748b', display: 'block', marginTop: '4px' }}>Склад отримувач визначено згідно з заявкою</span>
                </div>
              </div>

              <div className="form-group mb-8">
                <label>Вільний Транспорт</label>
                <select className="erp-input" value={dispatchForm.vehicle_id} onChange={e => setDispatchForm({...dispatchForm, vehicle_id: e.target.value})} required>
                  <option value="" disabled>Оберіть транспорт...</option>
                  {availableVehicles.map(v => <option key={v.id} value={v.id}>{v.brand} ({v.plate_number}) - Макс {v.capacity_kg} кг</option>)}
                </select>
                {availableVehicles.length === 0 && <span className="error-text">Немає вільних авто для логістики!</span>}
              </div>

              {selectedVehicle && (
                <div className="capacity-indicator">
                  <div className="capacity-header"><span className="capacity-label">Завантаженість кузова</span><span className={`capacity-value ${isOverweight ? 'text-critical' : 'text-normal'}`}>{fillPercentage.toFixed(1)}% ({currentTotalWeight.toFixed(1)} / {selectedVehicle.capacity_kg} кг)</span></div>
                  <div className="progress-bg"><div className={`progress-fill ${barStatusClass}`} style={{ width: `${fillPercentage}%` }} /></div>
                </div>
              )}

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowDispatchModal(false)}>Скасувати</button>
                <button type="submit" className="btn btn-dispatch" disabled={availableVehicles.length === 0 || isOverweight || allowedSourceWarehouses.length === 0}>Відправити рейс 🚀</button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className="card">
        {filteredRequests.length === 0 ? (
          <p className="empty-state">Заявок не знайдено</p>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                {canApprove && <th className="th-checkbox"></th>}
                <th>Ресурс</th>
                <th>Кількість</th>
                <th>Склад отримувач</th>
                <th>Статус</th>
                <th>Дата</th>
                {canApprove && <th>Дії</th>}
              </tr>
            </thead>
            <tbody>
              {filteredRequests.map((r) => {
                const isLocked = activeTargetWarehouseId !== null && activeTargetWarehouseId !== r.target_warehouse_id;
                const isSelected = selectedReqIds.has(r.id);
                const authorUser = users.find(u => u.id === r.created_by);
                const targetWarehouse = warehouses.find(w => w.id === r.target_warehouse_id);
                
                return (
                <tr key={r.id} className={`${isSelected ? 'row-selected' : ''} ${isLocked ? 'row-locked' : ''}`}>
                  {canApprove && (
                    <td className="td-checkbox">
                      {r.status === 'APPROVED' ? (
                        <input 
                          type="checkbox" 
                          className="custom-checkbox"
                          checked={isSelected} 
                          onChange={() => toggleSelection(r.id)}
                          disabled={isLocked}
                        />
                      ) : null}
                    </td>
                  )}
                  <td className="font-medium">{resources.find((res) => res.id === r.resource_id)?.name || r.resource_id}</td>
                  <td>{r.quantity} шт</td>
                  <td className="text-muted" style={{ fontSize: '13px' }}>
                    <div style={{ fontWeight: 600, color: '#334155' }}>
                      📍 {targetWarehouse?.name || 'Не вказано'}
                    </div>
                    <div style={{ fontSize: '11px', color: '#94a3b8' }}>Замовив: {authorUser?.full_name}</div>
                  </td>
                  <td>
                    <span className={`badge badge-${r.status === 'PENDING' ? 'warning' : r.status === 'APPROVED' ? 'success' : r.status === 'DISPATCHED' ? 'warning' : 'neutral'}`}>
                      {statusLabel[r.status] || r.status}
                    </span>
                  </td>
                  <td className="text-muted">{new Date(r.created_at).toLocaleDateString('uk-UA')}</td>
                  
                  {canApprove && (
                    <td>
                      {r.status === 'PENDING' ? (
                        <div className="action-buttons-flex">
                          <button className="btn btn-sm btn-primary" onClick={() => handleApprove(r.id, true)}>Затвердити</button>
                          <button className="btn btn-sm btn-danger-outline" onClick={() => handleApprove(r.id, false)}>Відхилити</button>
                        </div>
                      ) : r.status === 'APPROVED' ? (
                        <span className="status-text-waiting">
                          {isLocked ? '⛔ Інший напрямок' : 'Очікує логістику'}
                        </span>
                      ) : r.status === 'DISPATCHED' ? (
                        <span className="status-text-waiting" style={{ color: '#d97706' }}>
                          🚛 В дорозі
                        </span>
                      ) : (
                        <span className="status-text-closed">
                          {r.status === 'COMPLETED' ? '✅ Доставлено' : '🔒 Закрито'}
                        </span>
                      )}
                    </td>
                  )}
                </tr>
              )})}
            </tbody>
          </table>
        )}
      </div>

      {showForm && (
        <div className="modal-overlay" onClick={() => setShowForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3 className="modal-title">Нова заявка на постачання</h3>
            <form onSubmit={handleCreate}>
              <div className="form-group">
                <label>Ресурс</label>
                <select className="erp-input" value={newReq.resource_id} onChange={(e) => setNewReq({ ...newReq, resource_id: e.target.value })} required>
                  <option value="" disabled>Оберіть ресурс</option>
                  {resources.map((r) => <option key={r.id} value={r.id}>{r.name} (залишок: {r.quantity})</option>)}
                </select>
              </div>
              <div className="form-group">
                <label>Кількість</label>
                <input className="erp-input" type="number" min={1} value={newReq.quantity} onChange={(e) => setNewReq({ ...newReq, quantity: parseInt(e.target.value) || 1 })} required />
              </div>
              <div className="form-group">
                <label>На який склад доставити?</label>
                <select className="erp-input" value={newReq.target_warehouse_id} onChange={(e) => setNewReq({ ...newReq, target_warehouse_id: e.target.value })} required>
                  <option value="" disabled>Оберіть ваш склад...</option>
                  {warehouses.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
                </select>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowForm(false)}>Скасувати</button>
                <button type="submit" className="btn btn-primary">Створити</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}