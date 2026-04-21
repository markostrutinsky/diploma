import React, { useEffect, useState, useMemo } from 'react'
import { api, type SupplyRequest, type Resource, type Vehicle, type Warehouse, type User, type Unit, type VehicleBin, type RequestItem } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import { usePermissions } from '../hooks/usePermissions'
import { PaywallBadge } from '../components/FeatureGate'
import toast from 'react-hot-toast'
import './Requests.css'

const APPROVAL_MATRIX: Record<string, string[]> = {
  // Заявки звичайних виконавців погоджують їхні ліди, АБО логісти, АБО вище керівництво
  'EMPLOYEE': ['TEAM_LEAD', 'DEPT_MANAGER', 'DEPT_SUPERVISOR', 'BRANCH_MANAGER', 'BRANCH_LOGISTICIAN', 'REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'ADMIN'],
  
  // Тімліда погоджує менеджер відділу, логіст філії/регіону або керівництво
  'TEAM_LEAD': ['DEPT_MANAGER', 'DEPT_SUPERVISOR', 'BRANCH_MANAGER', 'BRANCH_LOGISTICIAN', 'REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'ADMIN'],
  
  'DEPT_SUPERVISOR': ['DEPT_MANAGER', 'BRANCH_MANAGER', 'BRANCH_LOGISTICIAN', 'REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'ADMIN'],
  'DEPT_MANAGER': ['BRANCH_MANAGER', 'BRANCH_LOGISTICIAN', 'REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'ADMIN'],
  
  // Керівника філії погоджує директор регіону або логіст регіону
  'BRANCH_MANAGER': ['REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'ADMIN'],
  
  'BRANCH_LOGISTICIAN': ['BRANCH_MANAGER', 'REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'ADMIN'],
  'BRANCH_STOREKEEPER': ['BRANCH_MANAGER', 'BRANCH_LOGISTICIAN', 'REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'ADMIN'],
  
  'REGION_STOREKEEPER': ['REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'ADMIN'],
  
  // А от верхівка погоджується тільки Адміном
  'REGION_DIRECTOR': ['ADMIN'],
  'REGION_LOGISTICIAN': ['REGION_DIRECTOR', 'ADMIN']
};

export default function Requests() {
  const { user } = useAuth()
  const [requests, setRequests] = useState<SupplyRequest[]>([])
  const [resources, setResources] = useState<Resource[]>([])
  const [warehouses, setWarehouses] = useState<Warehouse[]>([])
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [units, setUnits] = useState<Unit[]>([]) 
  
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  
  const [newReq, setNewReq] = useState({ resource_id: '', quantity: 1, target_warehouse_id: '' })

  const [filterStatus, setFilterStatus] = useState<string>('ALL')
  const [filterWarehouseId, setFilterWarehouseId] = useState<string>('ALL')
  const [searchQuery, setSearchQuery] = useState('')

  const [selectedReqIds, setSelectedReqIds] = useState<Set<string>>(new Set())
  const [showDispatchModal, setShowDispatchModal] = useState(false)
  const [dispatchForm, setDispatchForm] = useState({
    from_warehouse_id: '',
    to_warehouse_id: '',
    vehicle_id: '',
    priority: 'NORMAL'
  })

  useEffect(() => {
    if (dispatchForm.from_warehouse_id && dispatchForm.to_warehouse_id) {
      
      const fromWh = warehouses.find(w => String(w.id) === String(dispatchForm.from_warehouse_id));
      const toWh = warehouses.find(w => String(w.id) === String(dispatchForm.to_warehouse_id));

      if (fromWh?.unit_id && toWh?.unit_id) {
        api.vehicles.getAvailableForRoute(fromWh.unit_id, toWh.unit_id)
          .then(data => {
            // 🔥 ГОЛОВНИЙ ФІКС: Якщо бекенд прислав null, робимо з нього порожній масив
            const safeData = Array.isArray(data) ? data : [];
            
            setVehicles(safeData);
            
            if (dispatchForm.vehicle_id && !safeData.find(v => String(v.id) === String(dispatchForm.vehicle_id))) {
              setDispatchForm(prev => ({ ...prev, vehicle_id: '' }));
            }
          })
          .catch(err => {
            console.error("Помилка завантаження авто:", err);
          });
      } else {
        setVehicles([]);
      }
    } else {
      setVehicles([]);
    }
  }, [dispatchForm.from_warehouse_id, dispatchForm.to_warehouse_id, warehouses]);
  // --- СТЕЙТИ ДЛЯ SMART РОЗПОДІЛУ (AI) ---
  const [showSmartPreview, setShowSmartPreview] = useState(false)
  const [smartRoutes, setSmartRoutes] = useState<VehicleBin[]>([])
  const [unassignedItems, setUnassignedItems] = useState<RequestItem[]>([])

  const [rejectModalData, setRejectModalData] = useState<SupplyRequest | null>(null)
  const [rejectComment, setRejectComment] = useState('')
  const [cancelModalData, setCancelModalData] = useState<SupplyRequest | null>(null)
  const [isProcessing, setIsProcessing] = useState(false)

  const perms = usePermissions()
  const canCreate = perms.can('request_create')
  const canApprove = perms.can('request_approve')
  const hasSmartDispatch = perms.hasFeature('smart_dispatch')

  const loadData = async () => {
    setLoading(true)
    try {
      const token = localStorage.getItem('token')
      const [reqs, resRes, whs, vehs, usersRes, unitsRes] = await Promise.all([
        api.requests.list().catch(() => []),
        api.inventory.listResources(undefined).catch(() => []),
        api.warehouses.list().catch(() => []),
        // Залишаємо фолбеки для тих методів, які ти ще не переніс в клієнт (якщо такі є)
        (api as any).vehicles?.list().catch(() => fetch('/api/vehicles', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json()).catch(() => [])) || [],
        api.users.getVisible().catch(() => fetch('/api/users', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json()).catch(() => [])),
        api.units.list().catch(() => []) 
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
    
    const query = searchQuery.toLowerCase().trim();
    let matchSearch = true;
    
    if (query !== '') {
      const resourceName = (resources.find(res => res.id === r.resource_id)?.name || '').toLowerCase();
      const authorName = (users.find(u => u.id === r.created_by)?.full_name || '').toLowerCase();
      const warehouseName = (warehouses.find(w => w.id === r.target_warehouse_id)?.name || '').toLowerCase();
      const reqId = (r.id || '').toLowerCase();

      matchSearch = resourceName.includes(query) || authorName.includes(query) || warehouseName.includes(query) || reqId.includes(query);
    }

    return matchStatus && matchWarehouse && matchSearch;
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

  const allowedSourceWarehouses = useMemo(() => {
    if (!dispatchForm.to_warehouse_id || units.length === 0) return [];
    const targetWarehouse = warehouses.find(w => w.id === dispatchForm.to_warehouse_id);
    if (!targetWarehouse) return [];
    const targetUnit = units.find(u => u.id === targetWarehouse.unit_id);
    if (!targetUnit) return [];

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

    const requiredItems = selectedRequestsDetails.reduce((acc, req) => {
      const res = resources.find(r => r.id === req.resource_id);
      if (res && res.name) {
        acc[res.name] = (acc[res.name] || 0) + req.quantity;
      }
      return acc;
    }, {} as Record<string, number>);

    return hierarchicallyAllowed.filter(w => {
      for (const [name, neededQty] of Object.entries(requiredItems)) {
        const availableQty = resources
          .filter(r => r.warehouse_id === w.id && r.name === name)
          .reduce((sum, r) => sum + r.quantity, 0);
        if (availableQty < neededQty) return false;
      }
      return true; 
    });
  }, [dispatchForm.to_warehouse_id, warehouses, units, selectedRequestsDetails, resources]);

  const localAlternatives = useMemo(() => {
    if (!newReq.resource_id || !newReq.target_warehouse_id) return [];
    const targetW = warehouses.find(w => w.id === newReq.target_warehouse_id);
    const resName = resources.find(r => r.id === newReq.resource_id)?.name;
    if (!targetW || !resName) return [];

    return warehouses
      .filter(w => w.unit_id === targetW.unit_id && w.id !== targetW.id)
      .map(w => {
        const availableQty = resources
          .filter(r => r.warehouse_id === w.id && r.name === resName)
          .reduce((sum, r) => sum + r.quantity, 0);
        return { warehouse: w, availableQty };
      })
      .filter(info => info.availableQty > 0);
  }, [newReq.resource_id, newReq.target_warehouse_id, warehouses, resources]);

  const canApproveThis = useMemo(() => (r: SupplyRequest) => {
  if (!user) return false;
  
  // 1. Самопогодження заборонено жорстко
  if (r.created_by === user.id) return false;
  
  // 2. Знаходимо, хто створив заявку
  const creator = users.find(u => u.id === r.created_by);
  if (!creator) return false;
  
  // 3. Дістаємо масив ролей, які мають право погоджувати заявки від цього творця
  const allowedApprovers = APPROVAL_MATRIX[creator.role] || [];
  
  // 4. Перевіряємо, чи є поточна роль юзера в цьому списку
  return allowedApprovers.includes(user.role);
}, [user, users]);

  const handleOpenDispatchModal = () => {
    setDispatchForm(prev => ({
      ...prev,
      to_warehouse_id: activeTargetWarehouseId || '',
      from_warehouse_id: '' 
    }))
    setShowDispatchModal(true)
  }

  // --- 1. ВІДПРАВКА РУЧНОГО РЕЙСУ ЧЕРЕЗ КЛІЄНТ ---
  const handleDispatchSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (selectedReqIds.size === 0) return toast.error('Не вибрано жодної заявки!')
    if (!selectedVehicle) return toast.error('Оберіть транспорт!')
    if (!dispatchForm.from_warehouse_id) return toast.error('Оберіть склад відправник!')
    if (isOverweight) return toast.error(`Перевантаження! Максимум ${selectedVehicle.capacity_kg} кг`)

    const toastId = 'dispatch_toast'
    toast.loading('Формуємо збірний рейс...', { id: toastId })

    try {
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

      await (api as any).inventory.createShipment(payload)

      toast.success(`🚚 Збірний рейс відправлено!`, { id: toastId, duration: 4000 })
      setShowDispatchModal(false)
      setSelectedReqIds(new Set()) 
      loadData() 
    } catch (error: any) { 
      toast.error(error.message || 'Не вдалося створити рейс', { id: toastId, duration: 5000 }) 
    }
  }

  // --- 2. ВИКЛИК SMART РОЗПОДІЛУ ЧЕРЕЗ КЛІЄНТ ---
  const handleSmartDispatchPreview = async () => {
    if (selectedReqIds.size === 0) return toast.error('Виберіть хоча б одну заявку!');
    
    const toastId = toast.loading('🧠 Алгоритм First-Fit Decreasing аналізує вантаж...');
    
    try {
      const data = await (api as any).inventory.smartDispatchPreview(Array.from(selectedReqIds));
      
      setSmartRoutes(data.routes || []);
      setUnassignedItems(data.unassigned || []);
      setShowSmartPreview(true);
      
      toast.success('Оптимальний розподіл знайдено!', { id: toastId });
    } catch (error: any) {
      toast.error(error.message || 'Не вдалося прорахувати маршрути', { id: toastId });
    }
  };

  // --- 3. ПІДТВЕРДЖЕННЯ SMART РОЗПОДІЛУ ЧЕРЕЗ КЛІЄНТ ---
  const confirmSmartRoutes = async () => {
    const toastId = toast.loading('🚀 Формуємо серію рейсів...');
    
    try {
      await (api as any).inventory.smartDispatchConfirm(smartRoutes);

      toast.success('Всі рейси успішно відправлено!', { id: toastId });
      setShowSmartPreview(false);
      setSelectedReqIds(new Set());
      loadData();
    } catch (error: any) {
      toast.error(error.message || 'Помилка при збереженні рейсів', { id: toastId });
    }
  };

  const handleCreate = async (e: React.FormEvent) => { 
    e.preventDefault(); 
    if (!newReq.target_warehouse_id) return toast.error("❌ Оберіть цільовий склад!", { duration: 5000 })
    try { 
      await api.requests.create(newReq as any); 
      setShowForm(false); 
      setNewReq({ resource_id: resources[0]?.id || '', quantity: 1, target_warehouse_id: warehouses[0]?.id || '' }); 
      loadData(); 
      toast.success('Заявку створено!') 
    } catch (err) { toast.error(err instanceof Error ? err.message : 'Помилка') } 
  }

  const handleApprove = async (id: string) => { 
    try { 
      await (api.requests.approve as any)(id, true); 
      toast.success('Заявку погоджено!'); 
      loadData(); 
    } catch (err) { toast.error(err instanceof Error ? err.message : 'Помилка погодження') } 
  }

  const handleRejectSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!rejectModalData) return;
    setIsProcessing(true);
    try {
      await api.requests.reject(rejectModalData.id, rejectComment);
      toast.success('Заявку відхилено');
      setRejectModalData(null);
      setRejectComment('');
      loadData();
    } catch (err: any) {
      toast.error(err.message || 'Помилка відхилення');
    } finally { setIsProcessing(false); }
  }

  const handleCancel = async () => {
    if (!cancelModalData) return;
    setIsProcessing(true);
    try {
      await api.requests.cancel(cancelModalData.id);
      toast.success('Вашу заявку скасовано');
      setCancelModalData(null);
      loadData();
    } catch (err: any) {
      toast.error(err.message || 'Помилка скасування');
    } finally { setIsProcessing(false); }
  }

  const statusLabel: Record<string, string> = { 
    PENDING: 'Очікує', 
    APPROVED: 'Затверджено', 
    DISPATCHED: 'В дорозі', 
    REJECTED: 'Відхилено', 
    COMPLETED: 'Виконано', 
    OPEN: 'Відкрито',
    CANCELLED: 'Скасовано'
  }
  
  const availableVehicles = vehicles.filter(v => v.status === 'ACTIVE' && (v.type === 'VAN' || v.type === 'TRUCK' || v.type === 'PICKUP'))

  if (loading) return <div className="page-loading"><div className="spinner" /></div>

  const showActionsColumn = canApprove || canCreate;

  return (
    <div className="requests-page">
      
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px', paddingTop: '16px', flexWrap: 'wrap', gap: '16px' }}>
        <div>
          <h1 style={{ margin: '0 0 6px 0', fontSize: '1.75rem', fontWeight: 'bold', color: '#0f172a' }}>
            Заявки на постачання
          </h1>
          <p style={{ margin: 0, color: '#64748b', fontSize: '14px' }}>
            Управління потребами складів
          </p>
        </div>
        <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
          {selectedReqIds.size > 0 && canApprove && (
            <>
              <button className="btn btn-secondary" onClick={handleOpenDispatchModal}>
                🚚 Ручний рейс ({selectedReqIds.size})
              </button>
              
              {hasSmartDispatch ? (
                <button 
                  className="btn btn-smart" 
                  onClick={handleSmartDispatchPreview}
                >
                  ✨ Smart Розподіл (AI)
                </button>
              ) : (
                <button
                  className="btn btn-smart"
                  onClick={() => toast(
                    'Smart Розподіл — платна фіча. Перегляньте тариф PRO на сторінці «Тарифні плани».',
                    { icon: '🔒', duration: 5000 }
                  )}
                  style={{ opacity: 0.6, cursor: 'not-allowed' }}
                  title="Доступно на тарифі PRO"
                >
                  🔒 Smart Розподіл <PaywallBadge feature="smart_dispatch" compact />
                </button>
              )}
            </>
          )}
          {canCreate && (
            <button className="btn btn-primary" onClick={() => setShowForm(true)}>
              + Нова заявка
            </button>
          )}
        </div>
      </div>

      <div className="filters-bar" style={{ display: 'flex', gap: '16px', marginBottom: '24px', backgroundColor: '#f8fafc', padding: '16px', borderRadius: '8px', border: '1px solid #e2e8f0', flexWrap: 'wrap' }}>
        
        <div style={{ flex: '1 1 250px' }}>
          <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: '#64748b', marginBottom: '4px' }}>Пошук</label>
          <div style={{ position: 'relative' }}>
            <span style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: '#94a3b8', fontSize: '14px' }}>
              🔍
            </span>
            <input
              type="text"
              className="erp-input"
              placeholder="ID, Ресурс, Автор..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              style={{ paddingLeft: '35px' }}
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

        <div style={{ flex: '1 1 200px' }}>
          <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: '#64748b', marginBottom: '4px' }}>Статус заявки</label>
          <select className="erp-input" value={filterStatus} onChange={e => setFilterStatus(e.target.value)}>
            <option value="ALL">Всі статуси</option>
            <option value="PENDING">⏳ Очікують погодження</option>
            <option value="APPROVED">📦 Затверджені (Очікують логістику)</option>
            <option value="DISPATCHED">🚛 В дорозі (Прямують на склад)</option>
            <option value="COMPLETED">✅ Доставлені на склад</option>
            <option value="REJECTED">❌ Відхилені логістом</option>
            <option value="CANCELLED">🚫 Скасовані ініціатором</option>
          </select>
        </div>
        
        <div style={{ flex: '1 1 200px' }}>
          <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: '#64748b', marginBottom: '4px' }}>Цільовий Склад</label>
          <select className="erp-input" value={filterWarehouseId} onChange={e => setFilterWarehouseId(e.target.value)}>
            <option value="ALL">Всі склади</option>
            {warehouses.map(w => <option key={w.id} value={w.id}>{w.name}</option>)}
          </select>
        </div>
      </div>

      {/* --- МОДАЛКА SMART РОЗПОДІЛУ --- */}
      {showSmartPreview && (
        <div className="modal-overlay" onClick={() => setShowSmartPreview(false)}>
          <div className="modal modal-wide" onClick={(e) => e.stopPropagation()}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '15px', marginBottom: '24px' }}>
              <div style={{ fontSize: '2.5rem', background: '#f5f3ff', padding: '10px', borderRadius: '12px' }}>🧠</div>
              <div>
                <h3 className="modal-title" style={{ margin: 0, fontSize: '1.4rem' }}>Інтелектуальна маршрутизація</h3>
                <p style={{ margin: '4px 0 0', color: '#64748b', fontSize: '13px' }}>
                  Алгоритм оптимального пакування (First-Fit Decreasing) розподілив вантаж
                </p>
              </div>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', marginBottom: '30px', maxHeight: '50vh', overflowY: 'auto', paddingRight: '5px' }}>
              {smartRoutes.map((route, idx) => {
                const fillPercentage = Math.min(100, (route.used_weight / route.max_weight) * 100);
                const isOverweight = fillPercentage >= 100;

                return (
                  <div key={idx} className="smart-route-card">
                    <div className="smart-route-header">
                      <strong style={{ color: '#0f172a', fontSize: '1.1rem' }}>🚛 {route.name}</strong>
                      <span className={`badge ${isOverweight ? 'badge-critical' : 'badge-success'}`}>
                        Завантажено: {Math.round(fillPercentage)}%
                      </span>
                    </div>
                    
                    {/* Прогрес-бар */}
                    <div className="smart-route-progress-bg">
                      <div 
                        className={`smart-route-progress-fill ${isOverweight ? 'bg-red' : 'bg-purple'}`} 
                        style={{ width: `${fillPercentage}%` }}
                      ></div>
                    </div>
                    <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '12px', color: '#64748b', marginBottom: '12px' }}>
                      <span>0 кг</span>
                      <span>{route.used_weight.toFixed(1)} / {route.max_weight} кг</span>
                    </div>

                    <ul className="smart-route-items">
                      {route.items.map((item: any) => (
                        <li key={item.id}>📦 {item.name} <span className="smart-item-weight">— {item.weight_kg} кг</span></li>
                      ))}
                    </ul>
                  </div>
                )
              })}

              {unassignedItems.length > 0 && (
                <div className="smart-unassigned-card">
                  <strong style={{ color: '#b91c1c', display: 'flex', alignItems: 'center', gap: '8px' }}>
                    ⚠️ Не вистачило вільних машин для:
                  </strong>
                  <ul className="smart-route-items mt-2">
                    {unassignedItems.map((item: any) => (
                      <li key={item.id} style={{ color: '#991b1b' }}>
                        {item.name} <span className="smart-item-weight" style={{ color: '#b91c1c' }}>({item.weight_kg} кг)</span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>

            <div className="modal-actions">
              <button type="button" className="btn btn-secondary" onClick={() => setShowSmartPreview(false)}>Скасувати</button>
              <button 
                type="button" 
                className="btn btn-success" 
                onClick={confirmSmartRoutes}
                disabled={smartRoutes.length === 0}
              >
                ✅ Затвердити та відправити рейси
              </button>
            </div>
          </div>
        </div>
      )}

      {rejectModalData && (
        <div className="modal-overlay" onClick={() => !isProcessing && setRejectModalData(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3 style={{ color: '#ef4444' }}>Відхилення заявки</h3>
            <p className="text-muted">
              Відхилити заявку на <strong>{resources.find(r => r.id === rejectModalData.resource_id)?.name}</strong> ({rejectModalData.quantity} шт.)?
            </p>
            <form onSubmit={handleRejectSubmit}>
              <div className="form-group" style={{ textAlign: 'left', marginTop: '15px' }}>
                <label>Причина відмови <span className="required">*</span></label>
                <textarea 
                  className="erp-input" 
                  rows={3}
                  placeholder="Наприклад: Відсутнє на центральному складі"
                  value={rejectComment} 
                  onChange={(e) => setRejectComment(e.target.value)} 
                  required 
                  disabled={isProcessing}
                />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setRejectModalData(null)} disabled={isProcessing}>Назад</button>
                <button type="submit" className="btn btn-danger" disabled={!rejectComment || isProcessing}>{isProcessing ? 'Обробка...' : 'Підтвердити відмову'}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {cancelModalData && (
        <div className="modal-overlay" onClick={() => !isProcessing && setCancelModalData(null)}>
          <div className="modal confirm-modal" onClick={(e) => e.stopPropagation()}>
            <h3 style={{ color: '#64748b' }}>Скасування заявки</h3>
            <p>Ви впевнені, що хочете відкликати свою заявку на <strong>{resources.find(r => r.id === cancelModalData.resource_id)?.name}</strong>?</p>
            <div className="modal-actions">
              <button className="btn btn-secondary" onClick={() => setCancelModalData(null)} disabled={isProcessing}>Ні, залишити</button>
              <button className="btn" style={{ backgroundColor: '#64748b', color: 'white' }} onClick={handleCancel} disabled={isProcessing}>{isProcessing ? 'Обробка...' : 'Так, скасувати'}</button>
            </div>
          </div>
        </div>
      )}

      {showDispatchModal && (
        <div className="modal-overlay" onClick={() => setShowDispatchModal(false)}>
          <div className="modal modal-wide" onClick={(e) => e.stopPropagation()}>
            <h3 className="modal-title">🚚 Логістика: Ручний рейс на склад</h3>
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
                  {allowedSourceWarehouses.length === 0 && <span className="error-text" style={{marginTop: '4px'}}>Немає доступних складів вище по ієрархії або на них недостатньо майна!</span>}
                </div>
                <div className="form-group flex-1 mb-0">
                  <label>Куди (Заблоковано системою)</label>
                  <select className="erp-input" value={dispatchForm.to_warehouse_id} disabled style={{ backgroundColor: '#e2e8f0', color: '#475569', cursor: 'not-allowed' }}>
                    {warehouses.map(w => <option key={w.id} value={w.id}>{w.name}</option>)}
                  </select>
                </div>
              </div>
              <div className="form-group mb-8">
                <label>Вільний Транспорт</label>
                <select 
                  className="erp-input" 
                  value={dispatchForm.vehicle_id} 
                  onChange={e => setDispatchForm({...dispatchForm, vehicle_id: e.target.value})} 
                  required
                  disabled={!dispatchForm.from_warehouse_id || !dispatchForm.to_warehouse_id}
                >
                  <option value="" disabled>
                    {!dispatchForm.from_warehouse_id || !dispatchForm.to_warehouse_id 
                      ? "Спочатку оберіть склади відправки та отримання" 
                      : "Оберіть транспорт..."}
                  </option>
                  {availableVehicles.map(v => <option key={v.id} value={v.id}>{v.brand} ({v.plate_number}) - Макс {v.capacity_kg} кг</option>)}
                </select>
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

      {/* Таблиця Заявок */}
      <div className="card">
        {filteredRequests.length === 0 ? (
          <div className="empty-state">
            {searchQuery ? `За запитом "${searchQuery}" нічого не знайдено` : 'Заявок не знайдено'}
          </div>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                {canApprove && <th className="th-checkbox"></th>}
                <th>Ресурс</th>
                <th>Кількість</th>
                <th>Склад отримувач & Коментарі</th>
                <th>Статус</th>
                <th>Дата</th>
                {showActionsColumn && <th>Дії</th>}
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
                  <td className="font-medium">
                    {resources.find((res) => res.id === r.resource_id)?.name || r.resource_id}
                    <div style={{fontSize: '0.75rem', color: '#94a3b8', marginTop: '2px', fontWeight: 'normal'}}>
                      ID: {r.id.split('-')[0].toUpperCase()}
                    </div>
                  </td>
                  <td style={{ fontWeight: 600 }}>{r.quantity} шт</td>
                  <td className="text-muted" style={{ fontSize: '13px' }}>
                    <div style={{ fontWeight: 600, color: '#334155' }}>📍 {targetWarehouse?.name || 'Не вказано'}</div>
                    <div style={{ fontSize: '11px', color: '#94a3b8' }}>Замовив: {authorUser?.full_name}</div>
                    
                    {r.comment && (
                      <div className={`comment-box ${r.status === 'REJECTED' ? 'rejected' : ''}`}>
                        💬 {r.comment}
                      </div>
                    )}
                  </td>
                  <td>
                    <span className={`badge badge-${r.status === 'PENDING' ? 'warning' : r.status === 'APPROVED' ? 'success' : r.status === 'DISPATCHED' ? 'warning' : r.status === 'REJECTED' ? 'danger' : 'neutral'}`}>
                      {statusLabel[r.status] || r.status}
                    </span>
                  </td>
                  <td className="text-muted">{new Date(r.created_at).toLocaleDateString('uk-UA')}</td>
                  
                  {showActionsColumn && (
                    <td>
                      {r.status === 'PENDING' ? (
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', alignItems: 'flex-start' }}>
                          
                          {/* Розумна перевірка по матриці ієрархії */}
                          {canApproveThis(r) && (
                            <div className="action-buttons-flex">
                              <button className="btn btn-sm btn-primary" onClick={() => handleApprove(r.id)}>Затвердити</button>
                              <button className="btn btn-sm btn-danger-outline" onClick={() => setRejectModalData(r)}>Відхилити</button>
                            </div>
                          )}
                          
                          {/* Скасування тільки власних заявок */}
                          {r.created_by === user?.id && (
                             <button 
                               className="btn btn-sm" 
                               style={{ backgroundColor: '#f1f5f9', color: '#64748b', border: '1px dashed #cbd5e1', fontSize: '12px' }} 
                               onClick={() => setCancelModalData(r)}
                             >
                               Скасувати власну
                             </button>
                          )}

                        </div>
                      ) : r.status === 'APPROVED' ? (
                        <span className="status-text-waiting">{isLocked ? '⛔ Інший напрямок' : 'Очікує логістику'}</span>
                      ) : r.status === 'DISPATCHED' ? (
                        <span className="status-text-waiting" style={{ color: '#d97706' }}>🚛 В дорозі</span>
                      ) : (
                        <span className="status-text-closed">
                          {r.status === 'COMPLETED' ? '✅ Доставлено' : r.status === 'REJECTED' ? '❌ Відхилено' : '🔒 Закрито'}
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
                <input 
                  className="erp-input" 
                  type="number" 
                  min={1} 
                  value={newReq.quantity} 
                  onChange={(e) => {
                    const val = e.target.value;
                    setNewReq({ ...newReq, quantity: val === '' ? ('' as any) : parseInt(val) });
                  }} 
                  required 
                />
              </div>
              <div className="form-group">
                <label>На який склад доставити?</label>
                <select className="erp-input" value={newReq.target_warehouse_id} onChange={(e) => setNewReq({ ...newReq, target_warehouse_id: e.target.value })} required>
                  <option value="" disabled>Оберіть ваш склад...</option>
                  {warehouses.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
                </select>
              </div>

              {localAlternatives.length > 0 && (
                <div style={{ marginBottom: '16px', padding: '12px', backgroundColor: '#fffbeb', border: '1px solid #fde68a', borderRadius: '8px', fontSize: '13px', color: '#b45309' }}>
                  <strong style={{ display: 'block', marginBottom: '6px' }}>💡 Знайдено внутрішні резерви!</strong>
                  У вашій орг. структурі вже є цей ресурс на сусідніх складах:
                  <ul style={{ margin: '6px 0 0 20px', padding: 0 }}>
                    {localAlternatives.map(alt => (
                      <li key={alt.warehouse.id} style={{ marginBottom: '4px' }}>
                        {alt.warehouse.name} — <strong>{alt.availableQty} шт.</strong>
                      </li>
                    ))}
                  </ul>
                  <div style={{ marginTop: '8px', fontSize: '11px', fontStyle: 'italic', color: '#92400e', lineHeight: '1.4' }}>
                    * Замість того, щоб замовляти майно з центрального складу, можливо, варто попросити завідувача цього складу просто перемістити його вам.
                  </div>
                </div>
              )}

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