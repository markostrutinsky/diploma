import React, { useEffect, useState, useMemo, useRef } from 'react'
import { api, type SupplyRequest, type Resource, type Vehicle, type Warehouse, type User, type Unit, type VehicleBin, type RequestItem } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import { usePermissions } from '../hooks/usePermissions'
import { PaywallBadge } from '../components/FeatureGate'
import toast from 'react-hot-toast'
import Pagination from '../components/Pagination'
import SearchableSelect from '../components/SearchableSelect'
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
  const [users, setUsers] = useState<User[]>([])
  const [units, setUnits] = useState<Unit[]>([]) 
  
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  
  const [newReq, setNewReq] = useState({ 
    resource_name: '', 
    resource_category_id: '', 
    quantity: 1, 
    target_warehouse_id: '' 
  })
  
  const [uniqueResources, setUniqueResources] = useState<Array<{ name: string; category_id: string }>>([])

  const [filterStatus, setFilterStatus] = useState<string>('ALL')
  const [filterWarehouseId, setFilterWarehouseId] = useState<string>('ALL')
  const [searchQuery, setSearchQuery] = useState('')
  const [requestsPage, setRequestsPage] = useState(0)
  const REQUESTS_PAGE_SIZE = 25

  const [selectedReqIds, setSelectedReqIds] = useState<Set<string>>(new Set())
  const [showDispatchModal, setShowDispatchModal] = useState(false)
  const [dispatchForm, setDispatchForm] = useState({
    from_warehouse_id: '',
    to_warehouse_id: '',
    vehicle_id: '',
    priority: 'NORMAL'
  })

  // Окремий стейт для авто у модалці ручного рейсу — щоб оновлення списку авто
  // не перерендерювало всю велику сторінку із заявками.
  const [dispatchVehicles, setDispatchVehicles] = useState<Vehicle[]>([])
  const [dispatchVehiclesLoading, setDispatchVehiclesLoading] = useState(false)
  // Ref на актуальний масив warehouses, щоб useEffect нижче міг читати
  // свіжі дані без warehouses у deps (інакше loadData() щоразу міняє посилання
  // і ефект перезапускається, що й спричиняє «миготіння» модалки).
  const warehousesRef = useRef<Warehouse[]>([])

  useEffect(() => {
    if (dispatchForm.from_warehouse_id && dispatchForm.to_warehouse_id) {
      setDispatchVehiclesLoading(true)
      api.vehicles.getAvailableForRoute(dispatchForm.from_warehouse_id, dispatchForm.to_warehouse_id)
        .then(data => {
            const safeData = Array.isArray(data) ? data : [];
            setDispatchVehicles(safeData);
            if (dispatchForm.vehicle_id && !safeData.find(v => String(v.id) === String(dispatchForm.vehicle_id))) {
              setDispatchForm(prev => ({ ...prev, vehicle_id: '' }));
            }
          })
          .catch(err => {
            console.error("Помилка завантаження авто:", err);
          })
          .finally(() => setDispatchVehiclesLoading(false));
    } else {
      setDispatchVehicles([]);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [dispatchForm.from_warehouse_id, dispatchForm.to_warehouse_id]);
  // --- СТЕЙТИ ДЛЯ SMART РОЗПОДІЛУ (AI) ---
  const [showSmartPreview, setShowSmartPreview] = useState(false)
  const [smartRoutes, setSmartRoutes] = useState<VehicleBin[]>([])
  const [unassignedItems, setUnassignedItems] = useState<RequestItem[]>([])
  // Склад-відправник для Smart Розподілу (обирається ДО запуску алгоритму).
  const [smartFromWarehouseId, setSmartFromWarehouseId] = useState('')
  // Чи вже запустили розрахунок після вибору складу
  const [smartPreviewCalculated, setSmartPreviewCalculated] = useState(false)
  const [smartPreviewLoading, setSmartPreviewLoading] = useState(false)
  const [smartNoVehiclesMsg, setSmartNoVehiclesMsg] = useState('')

  const [rejectModalData, setRejectModalData] = useState<SupplyRequest | null>(null)
  const [rejectComment, setRejectComment] = useState('')
  const [cancelModalData, setCancelModalData] = useState<SupplyRequest | null>(null)
  const [isProcessing, setIsProcessing] = useState(false)

  const perms = usePermissions()
  const canCreate = perms.can('request_create')
  const canApprove = perms.can('request_approve')
  const hasSmartDispatch = perms.hasFeature('smart_dispatch')

  const loadData = async () => {
    try {
      const [reqs, resRes, uniqueRes, whs, usersRes, unitsRes] = await Promise.all([
        api.requests.list().catch(() => []),
        api.inventory.listResources(undefined).catch(() => []),
        api.inventory.getUniqueResourceNames(undefined).catch(() => []), // Завантажуємо унікальні назви
        api.warehouses.list().catch(() => []),
        api.users.getVisible().catch(() => []),
        api.units.list().catch(() => []) 
      ])
      
      setRequests(Array.isArray(reqs) ? reqs : [])
      console.log('📊 Loaded requests:', reqs)
      console.log('📊 Requests count:', Array.isArray(reqs) ? reqs.length : 0)
      setResources(Array.isArray(resRes) ? resRes : [])
      setUniqueResources(Array.isArray(uniqueRes) ? uniqueRes : []) // Зберігаємо унікальні назви
      console.log('🔍 Loaded unique resources:', uniqueRes)
      console.log('🔍 Unique resources count:', Array.isArray(uniqueRes) ? uniqueRes.length : 0)
      const whsArray = Array.isArray(whs) ? whs : []
      setWarehouses(whsArray)
      warehousesRef.current = whsArray  // оновлюємо ref синхронно зі стейтом
      setUsers(Array.isArray(usersRes) ? usersRes : [])
      setUnits(Array.isArray(unitsRes) ? unitsRes : [])
      
      // Встановлюємо перші значення при завантаженні даних
      if (Array.isArray(uniqueRes) && uniqueRes.length > 0) {
        setNewReq(prev => ({ 
          ...prev, 
          resource_name: prev.resource_name || uniqueRes[0].name,
          resource_category_id: prev.resource_category_id || uniqueRes[0].category_id
        }))
      }
      if (whsArray.length > 0) {
        setNewReq(prev => ({ 
          ...prev, 
          target_warehouse_id: prev.target_warehouse_id || whsArray[0].id 
        }))
      }
    } catch (error) { console.error(error) } finally { setLoading(false) }
  }

  useEffect(() => { loadData() }, [])
  useEffect(() => { setRequestsPage(0) }, [filterStatus, filterWarehouseId, searchQuery])

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
      const resourceName = (r.resource_name || '').toLowerCase();
      const authorName = (users.find(u => u.id === r.created_by)?.full_name || '').toLowerCase();
      const warehouseName = (warehouses.find(w => w.id === r.target_warehouse_id)?.name || '').toLowerCase();
      const reqId = (r.id || '').toLowerCase();

      matchSearch = resourceName.includes(query) || authorName.includes(query) || warehouseName.includes(query) || reqId.includes(query);
    }

    return matchStatus && matchWarehouse && matchSearch;
  })

  console.log('📊 Filtered requests count:', filteredRequests.length, 'Total:', requests.length, 'Filter status:', filterStatus)

  const selectedRequestsDetails = requests.filter(r => selectedReqIds.has(r.id))
  const currentTotalWeight = useMemo(() => {
    return selectedRequestsDetails.reduce((sum, req) => {
      // Шукаємо ресурс по назві та категорії
      const resource = resources.find(res => 
        res.name === req.resource_name && 
        (!req.resource_category_id || res.category_id === req.resource_category_id)
      )
      return sum + ((resource?.weight_kg || 1) * req.quantity)
    }, 0)
  }, [selectedRequestsDetails, resources])

  const selectedVehicle = dispatchVehicles.find(v => v.id === dispatchForm.vehicle_id)
  const isOverweight = selectedVehicle ? currentTotalWeight > selectedVehicle.capacity_kg : false
  const fillPercentage = selectedVehicle ? Math.min(100, (currentTotalWeight / selectedVehicle.capacity_kg) * 100) : 0
  let barStatusClass = fillPercentage >= 100 ? 'bar-critical' : fillPercentage > 80 ? 'bar-warning' : 'bar-safe' 

  const allowedSourceWarehouses = useMemo(() => {
    if (!dispatchForm.to_warehouse_id) return [];

    // Збираємо список необхідних ресурсів з обраних заявок
    const requiredItems = selectedRequestsDetails.reduce((acc, req) => {
      const name = req.resource_name;
      if (name) {
        acc[name] = (acc[name] || 0) + req.quantity;
      }
      return acc;
    }, {} as Record<string, number>);

    // Показуємо будь-який склад (крім складу-призначення),
    // де є ВСІ потрібні ресурси у достатній кількості.
    // Ієрархічного обмеження немає: ресурс може фізично знаходитись
    // в будь-якому складі незалежно від гілки організаційного дерева.
    return warehouses.filter(w => {
      if (w.id === dispatchForm.to_warehouse_id) return false;
      for (const [name, neededQty] of Object.entries(requiredItems)) {
        const availableQty = resources
          .filter(r => r.warehouse_id === w.id && r.name === name)
          .reduce((sum, r) => sum + r.quantity, 0);
        if (availableQty < neededQty) return false;
      }
      return true;
    });
  }, [dispatchForm.to_warehouse_id, warehouses, selectedRequestsDetails, resources]);

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
    // 🎯 Тепер користувач обирає склад-відправника вручну,
    // тому що заявки не прив'язані до конкретного ресурсу на складі
    
    // Перевіряємо, чи всі заявки мають однаковий target_warehouse
    const targetWarehouseIds = new Set(selectedRequestsDetails.map(r => r.target_warehouse_id));
    
    if (targetWarehouseIds.size > 1) {
      return toast.error(
        '❌ Обрані заявки мають різні склади призначення! Оберіть заявки з одним цільовим складом.',
        { duration: 6000 }
      );
    }
    
    setDispatchForm(prev => ({
      ...prev,
      to_warehouse_id: activeTargetWarehouseId || '',
      from_warehouse_id: '' // Користувач обере вручну
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

    // Перевірка пального
    const MIN_FUEL_LITERS = 5;
    const vehicleFuel = selectedVehicle.current_fuel_liters ?? 0;
    if (vehicleFuel < MIN_FUEL_LITERS) {
      return toast.error(
        `⛽ Неможливо відправити рейс! Машина "${selectedVehicle.brand} (${selectedVehicle.plate_number})" має лише ${vehicleFuel.toFixed(1)} л пального. Заправте мінімум ${MIN_FUEL_LITERS} л перед рейсом.`,
        { duration: 7000 }
      );
    }

    const toastId = 'dispatch_toast'
    toast.loading('Формуємо збірний рейс...', { id: toastId })

    try {
      // Отримуємо ресурси з складу-відправника
      const warehouseResources = resources.filter(r => r.warehouse_id === dispatchForm.from_warehouse_id)
      
      const payloadItems = selectedRequestsDetails.map(req => {
        // Знаходимо ресурс на складі-відправнику по назві та категорії
        const matchingResource = warehouseResources.find(r => 
          r.name === req.resource_name && 
          (!req.resource_category_id || r.category_id === req.resource_category_id)
        )
        
        if (!matchingResource) {
          throw new Error(`Ресурс "${req.resource_name}" не знайдено на обраному складі-відправнику!`)
        }
        
        return {
          resource_id: matchingResource.id,
          quantity: req.quantity,
          request_id: req.id 
        }
      })

      const payload = {
        from_warehouse_id: dispatchForm.from_warehouse_id,
        to_warehouse_id: dispatchForm.to_warehouse_id,
        vehicle_id: dispatchForm.vehicle_id,
        priority: dispatchForm.priority,
        items: payloadItems 
      }

      await api.inventory.createShipment(payload)

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

    // 🔥 Перевірка: всі обрані заявки мають йти на один і той самий склад-отримувач,
    // інакше "розумний" розподіл по машинах змішає різні маршрути в одну фуру.
    const targets = new Set(
      selectedRequestsDetails.map(r => r.target_warehouse_id || '')
    );
    if (targets.size > 1) {
      return toast.error(
        'Обрані заявки йдуть на різні склади-отримувачі. Оберіть заявки з одним цільовим складом для Smart Розподілу.',
        { duration: 6000 }
      );
    }

    // Синхронізуємо to_warehouse_id щоб allowedSourceWarehouses працював у модалці.
    const commonTarget = Array.from(targets)[0] || '';
    setDispatchForm(prev => ({ ...prev, to_warehouse_id: commonTarget }));
    setSmartFromWarehouseId('');
    setSmartRoutes([]);
    setUnassignedItems([]);
    setSmartPreviewCalculated(false);
    setSmartNoVehiclesMsg('');
    setShowSmartPreview(true);
  };

  // Запускається після того як користувач обрав склад-відправник у Smart модалці
  const runSmartCalculation = async (overrideFromId?: string) => {
    const fromId = overrideFromId ?? smartFromWarehouseId;
    if (!fromId) return toast.error('Оберіть склад відправник!');

    setSmartPreviewLoading(true);
    setSmartNoVehiclesMsg('');
    const toastId = toast.loading('🧠 Алгоритм First-Fit Decreasing аналізує вантаж...');

    try {
      const data = await api.inventory.smartDispatchPreview(Array.from(selectedReqIds), fromId);

      setSmartRoutes(data.routes || []);
      setUnassignedItems(data.unassigned || []);
      setSmartPreviewCalculated(true);

      toast.success('Оптимальний розподіл знайдено!', { id: toastId });
    } catch (error: any) {
      const msg = error.message || 'Не вдалося прорахувати маршрути';
      setSmartNoVehiclesMsg(msg);
      setSmartRoutes([]);
      setUnassignedItems([]);
      setSmartPreviewCalculated(true);
      toast.dismiss(toastId);
    } finally {
      setSmartPreviewLoading(false);
    }
  };

  // --- 3. ПІДТВЕРДЖЕННЯ SMART РОЗПОДІЛУ ЧЕРЕЗ КЛІЄНТ ---
  const confirmSmartRoutes = async () => {
    if (!smartFromWarehouseId) {
      return toast.error('Оберіть склад відправник!');
    }
    if (smartRoutes.length === 0) {
      return toast.error('Немає маршрутів для відправки.');
    }

    // Перевірка пального: жодна машина не повинна мати < 5л
    const MIN_FUEL_LITERS = 5;
    const emptyVehicles = smartRoutes.filter(r => (r.fuel_liters ?? 0) < MIN_FUEL_LITERS);
    if (emptyVehicles.length > 0) {
      return toast.error(
        `⛽ Неможливо відправити рейси! Машини не мають достатньо пального (потрібно мінімум ${MIN_FUEL_LITERS} л): ${emptyVehicles.map(r => r.name).join(', ')}`,
        { duration: 7000 }
      );
    }

    const toastId = toast.loading('🚀 Формуємо серію рейсів...');

    try {
      const payload = {
        from_warehouse_id: smartFromWarehouseId,
        priority: dispatchForm.priority || 'NORMAL',
        routes: smartRoutes.map(r => ({
          vehicle_id: r.id,
          // В preview r.items.id — це supply_request_id (див. GetRequestsForDispatch).
          request_ids: r.items.map((i: any) => i.id),
        })),
      };

      const res = await api.inventory.smartDispatchConfirm(payload);

      toast.success(res.message || 'Всі рейси успішно відправлено!', { id: toastId });
      setShowSmartPreview(false);
      setSelectedReqIds(new Set());
      setSmartFromWarehouseId('');
      loadData();
    } catch (error: any) {
      toast.error(error.message || 'Помилка при збереженні рейсів', { id: toastId });
    }
  };

  const handleCreate = async (e: React.FormEvent) => { 
    e.preventDefault(); 
    if (!newReq.target_warehouse_id) return toast.error("❌ Оберіть цільовий склад!", { duration: 5000 })
    if (!newReq.resource_name) return toast.error("❌ Оберіть ресурс!", { duration: 5000 })
    try { 
      await api.requests.create(newReq); 
      setShowForm(false); 
      setNewReq({ 
        resource_name: uniqueResources[0]?.name || '', 
        resource_category_id: uniqueResources[0]?.category_id || '',
        quantity: 1, 
        target_warehouse_id: warehouses[0]?.id || '' 
      }); 
      loadData(); 
      toast.success('Заявку створено!') 
    } catch (err) { toast.error(err instanceof Error ? err.message : 'Помилка') } 
  }

  const handleApprove = async (id: string) => { 
    try { 
      await api.requests.approve(id, true, undefined); 
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
    LOADING: 'Завантажується',
    DISPATCHED: 'В дорозі',
    REJECTED: 'Відхилено', 
    COMPLETED: 'Виконано', 
    OPEN: 'Відкрито',
    CANCELLED: 'Скасовано'
  }
  
  const availableVehicles = useMemo(() => {
    // Фільтруємо авто зі списку маршруту (dispatchVehicles), що повернув бекенд.
    // Бекенд вже повертає тільки ACTIVE авто цього маршруту, але фільтруємо
    // ще раз по типу кузова для надійності.
    return dispatchVehicles.filter(v => v.status === 'ACTIVE' && (v.type === 'VAN' || v.type === 'TRUCK' || v.type === 'PICKUP'))
  }, [dispatchVehicles])

  if (loading) return <div className="page-loading"><div className="spinner" /></div>

  const reqTotalPages = Math.max(1, Math.ceil(filteredRequests.length / REQUESTS_PAGE_SIZE));
  const safeReqPage = Math.min(requestsPage, reqTotalPages - 1);
  const pagedRequests = filteredRequests.slice(safeReqPage * REQUESTS_PAGE_SIZE, (safeReqPage + 1) * REQUESTS_PAGE_SIZE);

  const showActionsColumn = canApprove || canCreate;

  return (
    <div className="requests-page">
      
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        flexWrap: 'wrap',
        gap: '16px',
        position: 'sticky',
        top: '-2rem',
        zIndex: 100,
        backgroundColor: 'var(--bg-body, #1a1a2e)',
        padding: 'calc(2rem + 16px) 0 12px 0',
        margin: '-2rem 0 8px 0',
        borderBottom: '1px solid var(--border-color, rgba(255,255,255,0.08))',
        boxShadow: '0 2px 8px rgba(0,0,0,0.6)',
      }}>
        <div>
          <h1 style={{ margin: '0 0 6px 0', fontSize: '1.75rem', fontWeight: 'bold', color: 'var(--text-bright)' }}>
            Заявки на постачання
          </h1>
          <p style={{ margin: 0, color: 'var(--text-muted)', fontSize: '14px' }}>
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

      <div className="filters-bar" style={{ display: 'flex', gap: '16px', marginBottom: '24px', backgroundColor: 'var(--bg-input)', padding: '16px', borderRadius: '8px', border: '1px solid var(--border)', flexWrap: 'wrap' }}>
        
        <div style={{ flex: '1 1 250px' }}>
          <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '4px' }}>Пошук</label>
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
          <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '4px' }}>Статус заявки</label>
          <select className="erp-input" value={filterStatus} onChange={e => setFilterStatus(e.target.value)}>
            <option value="ALL">Всі статуси</option>
            <option value="PENDING">⏳ Очікують погодження</option>
            <option value="APPROVED">📦 Затверджені (Очікують логістику)</option>
            <option value="LOADING">🔄 Завантажуються (Рейс сформовано)</option>
            <option value="DISPATCHED">🚛 В дорозі (Прямують на склад)</option>
            <option value="COMPLETED">✅ Доставлені на склад</option>
            <option value="REJECTED">❌ Відхилені логістом</option>
            <option value="CANCELLED">🚫 Скасовані ініціатором</option>
          </select>
        </div>
        
        <div style={{ flex: '1 1 200px' }}>
          <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: 'var(--text-muted)', marginBottom: '4px' }}>Цільовий Склад</label>
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
            <div style={{ marginBottom: '20px' }}>
              <h3 className="modal-title" style={{ margin: '0 0 4px 0' }}>✨ Інтелектуальна маршрутизація</h3>
              <p style={{ margin: 0, color: 'var(--text-muted)', fontSize: '13px' }}>
                Алгоритм First-Fit Decreasing підбере машини зі складу відправника або отримувача
              </p>
            </div>

            {/* КРОК 1: Оберіть склад-відправник — розрахунок запускається автоматично */}
            <div className="form-group" style={{ marginBottom: '16px' }}>
              <label style={{ display: 'block', fontSize: '0.8125rem', fontWeight: 600, color: 'var(--text)', marginBottom: '6px' }}>
                📦 Склад-відправник <span className="required">*</span>
              </label>
              <SearchableSelect
                options={allowedSourceWarehouses.map(w => {
                  const u = units.find(unit => unit.id === w.unit_id);
                  return { value: w.id, label: `${w.name}${u ? ` (${u.name})` : ''}` };
                })}
                value={smartFromWarehouseId}
                onChange={val => {
                  setSmartFromWarehouseId(val);
                  setSmartPreviewCalculated(false);
                  setSmartRoutes([]);
                  setUnassignedItems([]);
                  setSmartNoVehiclesMsg('');
                  // Автоматично запускаємо розрахунок після вибору складу
                  if (val) {
                    setTimeout(() => runSmartCalculation(val), 0);
                  }
                }}
                placeholder="Оберіть склад відправника..."
                searchPlaceholder="Пошук складу..."
              />
              {allowedSourceWarehouses.length === 0 && (
                <span className="error-text" style={{ marginTop: '4px' }}>
                  Немає доступних складів для цих заявок!
                </span>
              )}
              {smartPreviewLoading && (
                <div style={{ marginTop: '8px', fontSize: '13px', color: 'var(--text-muted)' }}>
                  ⏳ Алгоритм First-Fit Decreasing аналізує вантаж...
                </div>
              )}
            </div>

            {/* КРОК 2: Результати розрахунку */}
            {smartPreviewCalculated && (
              <>
                {smartNoVehiclesMsg ? (
                  <div style={{ padding: '16px', backgroundColor: 'rgba(239, 68, 68, 0.1)', border: '1px solid rgba(239, 68, 68, 0.3)', borderRadius: '8px', marginBottom: '16px' }}>
                    <strong style={{ color: '#ef4444' }}>⚠️ Немає доступного транспорту</strong>
                    <div style={{ marginTop: '6px', fontSize: '13px', color: 'var(--text-muted)', lineHeight: 1.5 }}>
                      {smartNoVehiclesMsg}. Перевірте наявність транспорту на сторінці "Транспорт" або оберіть інший склад відправника.
                    </div>
                  </div>
                ) : (
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', marginBottom: '24px', maxHeight: '45vh', overflowY: 'auto', paddingRight: '5px' }}>
                    {smartRoutes.map((route, idx) => {
                      const fillPercentage = Math.min(100, (route.used_weight / route.max_weight) * 100);
                      const isOverweight = fillPercentage >= 100;

                      return (
                        <div key={idx} className="smart-route-card">
                          <div className="smart-route-header">
                            <strong style={{ color: 'var(--text-bright)', fontSize: '1.1rem' }}>🚛 {route.name}</strong>
                            <span className={`badge ${isOverweight ? 'badge-critical' : 'badge-success'}`}>
                              Завантажено: {Math.round(fillPercentage)}%
                            </span>
                          </div>
                          <div className="smart-route-progress-bg">
                            <div
                              className={`smart-route-progress-fill ${isOverweight ? 'bg-red' : 'bg-purple'}`}
                              style={{ width: `${fillPercentage}%` }}
                            />
                          </div>
                          <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '12px', color: 'var(--text-muted)', marginBottom: '8px' }}>
                            <span>0 кг</span>
                            <span>{route.used_weight.toFixed(1)} / {route.max_weight} кг</span>
                          </div>
                          {/* ⛽ Рядок пального */}
                          {route.fuel_norm != null && route.tank_capacity != null && (
                            (() => {
                              const fuel = route.fuel_liters ?? 0;
                              const tank = route.tank_capacity ?? 1;
                              const norm = route.fuel_norm ?? 0;
                              const fuelPct = Math.min(100, Math.round((fuel / tank) * 100));
                              const maxRange = norm > 0 ? Math.floor(fuel / norm * 100) : null;
                              const isFuelLow = fuelPct < 20;
                              return (
                                <div style={{ marginBottom: '10px', padding: '8px', background: isFuelLow ? 'rgba(239,68,68,0.08)' : 'rgba(34,197,94,0.07)', borderRadius: '6px', fontSize: '12px', border: `1px solid ${isFuelLow ? 'rgba(239,68,68,0.25)' : 'rgba(34,197,94,0.2)'}` }}>
                                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '4px' }}>
                                    <span style={{ color: isFuelLow ? '#ef4444' : 'var(--text)', fontWeight: 600 }}>⛽ {fuel.toFixed(1)} / {tank} л</span>
                                    {maxRange !== null && <span style={{ color: 'var(--text-muted)' }}>~{maxRange} км запас</span>}
                                  </div>
                                  <div style={{ height: '5px', background: 'var(--border)', borderRadius: '3px', overflow: 'hidden' }}>
                                    <div style={{ height: '100%', width: `${fuelPct}%`, background: isFuelLow ? '#ef4444' : '#22c55e', borderRadius: '3px' }} />
                                  </div>
                                  {isFuelLow && <div style={{ color: '#ef4444', marginTop: '4px' }}>⚠️ Потребує заправки перед рейсом!</div>}
                                </div>
                              );
                            })()
                          )}
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
                )}
              </>
            )}

            <div className="modal-actions" style={{ justifyContent: 'flex-end' }}>
              <button type="button" className="btn btn-secondary" onClick={() => setShowSmartPreview(false)}>Скасувати</button>
              <button
                type="button"
                className="btn btn-success"
                onClick={confirmSmartRoutes}
                disabled={smartRoutes.length === 0 || !smartFromWarehouseId}
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
              Відхилити заявку на <strong>{rejectModalData.resource_name}</strong> ({rejectModalData.quantity} шт.)?
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
            <h3 style={{ color: 'var(--text-muted)' }}>Скасування заявки</h3>
            <p>Ви впевнені, що хочете відкликати свою заявку на <strong>{cancelModalData.resource_name}</strong>?</p>
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
                  // Шукаємо ресурс по назві та категорії для отримання ваги
                  const res = resources.find(r => 
                    r.name === req.resource_name && 
                    (!req.resource_category_id || r.category_id === req.resource_category_id)
                  )
                  const itemWeight = (res?.weight_kg || 1) * req.quantity
                  return (
                    <li key={req.id} className="summary-item">
                      <span className="item-name">📦 {req.resource_name} — {req.quantity} шт.</span>
                      <span className="item-weight">~{itemWeight.toFixed(1)} кг</span>
                    </li>
                  )
                })}
              </ul>
              <div className="summary-total">Загальна розрахункова вага: {currentTotalWeight.toFixed(1)} кг</div>
            </div>
            <form onSubmit={handleDispatchSubmit}>
              <div className="dispatch-row">
                <div className="dispatch-col">
                  <label>Звідки відправляємо <span className="required">*</span></label>
                  <SearchableSelect
                    options={allowedSourceWarehouses.map(w => {
                      const u = units.find(unit => unit.id === w.unit_id);
                      return { value: w.id, label: `${w.name}${u ? ` (${u.name})` : ''}` };
                    })}
                    value={dispatchForm.from_warehouse_id}
                    onChange={(val) => {
                      setDispatchVehiclesLoading(true);
                      setDispatchVehicles([]);
                      setDispatchForm(prev => ({ ...prev, from_warehouse_id: val, vehicle_id: '' }));
                    }}
                    placeholder="Оберіть склад-відправника"
                    searchPlaceholder="Пошук складу..."
                  />
                </div>
                <div className="dispatch-col">
                  <label>🔒 Куди (фіксований системою)</label>
                  <input
                    type="text"
                    className="dispatch-locked-field"
                    value={warehouses.find(w => w.id === dispatchForm.to_warehouse_id)?.name || '—'}
                    disabled
                    readOnly
                  />
                </div>
              </div>
              <div className="form-group mb-8">
                <label>Пріоритет рейсу</label>
                <select
                  className="erp-input"
                  value={dispatchForm.priority}
                  onChange={e => setDispatchForm({...dispatchForm, priority: e.target.value})}
                >
                  <option value="NORMAL">🟢 Звичайний (Плановий)</option>
                  <option value="URGENT">🔴 Терміновий</option>
                </select>
              </div>
              <div className="form-group mb-8">
                <label>Вільний Транспорт</label>
                <select 
                  className="erp-input" 
                  value={dispatchForm.vehicle_id} 
                  onChange={e => setDispatchForm({...dispatchForm, vehicle_id: e.target.value})} 
                  required
                  disabled={!dispatchForm.from_warehouse_id || !dispatchForm.to_warehouse_id || dispatchVehiclesLoading}
                >
                  <option value="" disabled>
                    {dispatchVehiclesLoading
                      ? "Завантаження транспорту..."
                      : !dispatchForm.from_warehouse_id || !dispatchForm.to_warehouse_id 
                        ? "Спочатку оберіть склади відправки та отримання" 
                        : "Оберіть транспорт..."}
                  </option>
                  {availableVehicles.map(v => <option key={v.id} value={v.id}>{v.brand} ({v.plate_number}) - Макс {v.capacity_kg} кг</option>)}
                </select>
                {!dispatchVehiclesLoading && dispatchForm.from_warehouse_id && dispatchForm.to_warehouse_id && availableVehicles.length === 0 && (
                  <div style={{ marginTop: '8px', padding: '12px', backgroundColor: 'rgba(239, 68, 68, 0.1)', border: '1px solid rgba(239, 68, 68, 0.3)', borderRadius: '6px', fontSize: '13px', color: '#ef4444' }}>
                    <strong>⚠️ Немає доступного транспорту</strong>
                    <div style={{ marginTop: '4px', fontSize: '12px', lineHeight: '1.4' }}>
                      На маршруті між обраними складами немає вільних автомобілів зі статусом ACTIVE. Перевірте наявність транспорту на сторінці "Транспорт" або оберіть інші склади.
                    </div>
                  </div>
                )}
              </div>
              {selectedVehicle && (
                <div className="capacity-indicator">
                  <div className="capacity-header"><span className="capacity-label">Завантаженість кузова</span><span className={`capacity-value ${isOverweight ? 'text-critical' : 'text-normal'}`}>{fillPercentage.toFixed(1)}% ({currentTotalWeight.toFixed(1)} / {selectedVehicle.capacity_kg} кг)</span></div>
                  <div className="progress-bg"><div className={`progress-fill ${barStatusClass}`} style={{ width: `${fillPercentage}%` }} /></div>
                </div>
              )}
              {/* ⛽ Індикатор пального */}
              {selectedVehicle && (
                (() => {
                  const fuel = selectedVehicle.current_fuel_liters ?? 0;
                  const norm = selectedVehicle.fuel_norm ?? 0;
                  const tank = selectedVehicle.tank_capacity ?? 1;
                  const maxRangeKm = norm > 0 ? Math.floor(fuel / norm * 100) : null;
                  const fuelPct = Math.min(100, Math.round((fuel / tank) * 100));
                  const isFuelLow = fuelPct < 20;
                  return (
                    <div style={{ marginTop: '12px', padding: '12px', backgroundColor: isFuelLow ? 'rgba(239,68,68,0.08)' : 'rgba(34,197,94,0.07)', border: `1px solid ${isFuelLow ? 'rgba(239,68,68,0.3)' : 'rgba(34,197,94,0.25)'}`, borderRadius: '8px', fontSize: '13px' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '6px' }}>
                        <span style={{ fontWeight: 600, color: isFuelLow ? '#ef4444' : 'var(--text)' }}>⛽ Пальне: {fuel.toFixed(1)} / {tank} л</span>
                        <span style={{ color: 'var(--text-muted)', fontSize: '12px' }}>{norm} л/100 км</span>
                      </div>
                      <div style={{ height: '6px', background: 'var(--border)', borderRadius: '4px', overflow: 'hidden', marginBottom: '6px' }}>
                        <div style={{ height: '100%', width: `${fuelPct}%`, background: isFuelLow ? '#ef4444' : '#22c55e', borderRadius: '4px', transition: 'width 0.3s' }} />
                      </div>
                      {maxRangeKm !== null && (
                        <div style={{ color: isFuelLow ? '#ef4444' : 'var(--text-muted)', fontSize: '12px' }}>
                          {isFuelLow
                            ? `⚠️ Мало пального! Запас ходу ~${maxRangeKm} км. Заправте перед рейсом.`
                            : `🗺 Запас ходу ~${maxRangeKm} км`}
                        </div>
                      )}
                      {fuel === 0 && (
                        <div style={{ color: '#ef4444', fontWeight: 600, fontSize: '12px', marginTop: '4px' }}>
                          🚫 Бак порожній! Рейс неможливий без заправки.
                        </div>
                      )}
                    </div>
                  );
                })()
              )}
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowDispatchModal(false)}>Скасувати</button>
                <button 
                  type="submit" 
                  className="btn btn-dispatch" 
                  disabled={
                    !dispatchForm.vehicle_id || 
                    !dispatchForm.from_warehouse_id || 
                    isOverweight || 
                    allowedSourceWarehouses.length === 0
                  }
                >
                  Відправити рейс 🚀
                </button>
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
          <>
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
              {pagedRequests.map((r) => {
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
                    {r.resource_name || 'Невідомий ресурс'}
                    <div style={{fontSize: '0.75rem', color: '#94a3b8', marginTop: '2px', fontWeight: 'normal'}}>
                      ID: {r.id.split('-')[0].toUpperCase()}
                    </div>
                  </td>
                  <td style={{ fontWeight: 600 }}>{r.quantity} шт</td>
                  <td className="text-muted" style={{ fontSize: '13px' }}>
                    <div style={{ fontWeight: 600, color: 'var(--text)' }}>📍 {targetWarehouse?.name || 'Не вказано'}</div>
                    <div style={{ fontSize: '11px', color: '#94a3b8' }}>Замовив: {authorUser?.full_name}</div>
                    
                    {r.comment && (
                      <div className={`comment-box ${r.status === 'REJECTED' ? 'rejected' : ''}`}>
                        💬 {r.comment}
                      </div>
                    )}
                  </td>
                  <td>
                    <span className={`badge badge-${r.status === 'PENDING' ? 'warning' : r.status === 'APPROVED' ? 'success' : r.status === 'LOADING' ? 'info' : r.status === 'DISPATCHED' ? 'warning' : r.status === 'REJECTED' ? 'danger' : 'neutral'}`}>
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
                               style={{ backgroundColor: 'var(--bg-input)', color: 'var(--text-muted)', border: '1px dashed var(--border-light)', fontSize: '12px' }} 
                               onClick={() => setCancelModalData(r)}
                             >
                               Скасувати власну
                             </button>
                          )}

                        </div>
                      ) : r.status === 'APPROVED' ? (
                        <span className="status-text-waiting">{isLocked ? '⛔ Інший напрямок' : 'Очікує логістику'}</span>
                      ) : r.status === 'LOADING' ? (
                        <span className="status-text-waiting" style={{ color: '#6366f1' }}>📦 Завантажується</span>
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
          <Pagination
            currentPage={safeReqPage}
            totalPages={reqTotalPages}
            onPageChange={setRequestsPage}
            totalItems={filteredRequests.length}
            itemsPerPage={REQUESTS_PAGE_SIZE}
          />
          </>
        )}
      </div>

      {showForm && (
        <div className="modal-overlay" onClick={() => setShowForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3 className="modal-title">Нова заявка на постачання</h3>
            <form onSubmit={handleCreate}>
              <div className="form-group">
                <label>Ресурс <span className="required">*</span></label>
                <SearchableSelect
                  options={uniqueResources.map(r => ({ value: r.name, label: r.name }))}
                  value={newReq.resource_name}
                  onChange={(val) => {
                    const selectedResource = uniqueResources.find(r => r.name === val);
                    setNewReq({
                      ...newReq,
                      resource_name: val,
                      resource_category_id: selectedResource?.category_id || '',
                      target_warehouse_id: ''
                    });
                  }}
                  placeholder="Оберіть ресурс"
                  searchPlaceholder="Пошук ресурсу..."
                />
              </div>
              <div className="form-group">
                <label>Кількість <span className="required">*</span></label>
                <input 
                  className="erp-input" 
                  type="number" 
                  min={1} 
                  value={newReq.quantity} 
                  onChange={(e) => {
                    const val = e.target.value;
                    setNewReq({ ...newReq, quantity: val === '' ? 0 : parseInt(val) });
                  }} 
                  required 
                />
              </div>
              <div className="form-group">
                <label>Куди доставити? <span className="required">*</span></label>
                <SearchableSelect
                  options={warehouses.map(w => {
                    const u = units.find(unit => unit.id === w.unit_id);
                    return { value: w.id, label: `${w.name}${u ? ` (${u.name})` : ''}` };
                  })}
                  value={newReq.target_warehouse_id}
                  onChange={(val) => setNewReq({ ...newReq, target_warehouse_id: val })}
                  placeholder={newReq.resource_name ? 'Оберіть склад призначення...' : 'Спочатку оберіть ресурс'}
                  searchPlaceholder="Пошук складу..."
                  disabled={!newReq.resource_name}
                />
              </div>

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowForm(false)}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={!newReq.resource_name || !newReq.target_warehouse_id || newReq.quantity < 1}>Створити</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}