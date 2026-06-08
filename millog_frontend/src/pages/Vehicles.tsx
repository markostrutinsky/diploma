import React, { useEffect, useRef, useState } from 'react'
import { api, ROLE_NAMES, type Vehicle, type FuelRecordType, type FuelRecord, type MaintenanceRecord, type SystemUser, type DriverHistoryRecord, type VehicleShipmentRecord, type Warehouse, type LogShipmentRefuelPayload} from '../api/client'
import { usePermissions } from '../hooks/usePermissions'
import toast from 'react-hot-toast'
import { useNavigate } from 'react-router-dom'
import Pagination from '../components/Pagination'
import SearchableSelect from '../components/SearchableSelect'
import './Vehicles.css'

export default function Vehicles() {
  const navigate = useNavigate()
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [usersList, setUsersList] = useState<SystemUser[]>([])
  const [warehousesList, setWarehousesList] = useState<Warehouse[]>([])
  const [loading, setLoading] = useState(true)

  const [viewTab, setViewTab] = useState<'ACTIVE' | 'ARCHIVE'>('ACTIVE')

  // ---------------------------------------------------------
  // НОВИЙ СТЕЙТ ДЛЯ ПОШУКУ
  // ---------------------------------------------------------
  const [searchQuery, setSearchQuery] = useState('')

  const [showVehicleForm, setShowVehicleForm] = useState(false)
  const [fuelModalVehicle, setFuelModalVehicle] = useState<Vehicle | null>(null)
  
  const [maintenanceModalVehicle, setMaintenanceModalVehicle] = useState<Vehicle | null>(null)
  const [maintenanceForm, setMaintenanceForm] = useState<{
    odometer_km: number;
    description: string;
    performed_by: string;
    cost_amount: number;
    document: File | null;
  }>({
    odometer_km: 0,
    description: '',
    performed_by: '',
    cost_amount: 0,
    document: null
  })

  const [statusModalVehicle, setStatusModalVehicle] = useState<Vehicle | null>(null)
  const [statusForm, setStatusForm] = useState({
    status: 'IN_REPAIR',
    reason: ''
  })

  const [anomalyAlert, setAnomalyAlert] = useState<FuelRecord | null>(null)
  
  const [historyVehicle, setHistoryVehicle] = useState<Vehicle | null>(null)
  const [historyTab, setHistoryTab] = useState<'FUEL' | 'MAINTENANCE' | 'DRIVERS' | 'TRIPS'>('FUEL')
  const [fuelRecords, setFuelRecords] = useState<FuelRecord[]>([])
  const [maintenanceRecords, setMaintenanceRecords] = useState<MaintenanceRecord[]>([])
  const [historyLoading, setHistoryLoading] = useState(false)
  const [driverRecords, setDriverRecords] = useState<DriverHistoryRecord[]>([])
  const [tripRecords, setTripRecords] = useState<VehicleShipmentRecord[]>([])

  // ⛽ Дозаправка під час рейсу
  const [tripRefuelModal, setTripRefuelModal] = useState<{ shipmentId: string; fromWarehouse: string; toWarehouse: string } | null>(null)
  const [tripRefuelForm, setTripRefuelForm] = useState({ liters: '', station_name: '', odometer_km: '', cost_uah: '' })
  const [tripRefuelProcessing, setTripRefuelProcessing] = useState(false)

  const [driverModalVehicle, setDriverModalVehicle] = useState<Vehicle | null>(null)
  const [driverForm, setDriverForm] = useState({ driver_id: '' })
  const maintenanceFileRef = useRef<HTMLInputElement>(null)
  
  const [newVehicle, setNewVehicle] = useState({
    brand: '',
    model: '',
    plate_number: '',
    type: 'VAN',           
    capacity_kg: 1500,     
    tank_capacity: 0,
    fuel_norm: 0,
    driver_id: '',
    home_warehouse_id: ''
  })

  const [fuelForm, setFuelForm] = useState({
    record_type: 'EXPENSE' as FuelRecordType,
    liters: 0,
    odometer_km: '',
  })

  const [editingVehicle, setEditingVehicle] = useState<Vehicle | null>(null);
  const [vehicleToDelete, setVehicleToDelete] = useState<Vehicle | null>(null);
  const [editForm, setEditForm] = useState({ brand: '', model: '', plate_number: '', capacity_kg: 0 });
  const [isProcessing, setIsProcessing] = useState(false);

  const [activeMenuId, setActiveMenuId] = useState<string | null>(null);
  const perms = usePermissions();
  const hasFuelAntifraud = perms.hasFeature('fuel_antifraud');
  const hasPredictiveMaint = perms.hasFeature('predictive_maintenance');
  // Сумісність з існуючим UI, який посилається на isPro (показ фіч PRO)
  const isPro = hasFuelAntifraud || hasPredictiveMaint;

  useEffect(() => {
    const closeMenu = () => setActiveMenuId(null);
    document.addEventListener('click', closeMenu);
    return () => document.removeEventListener('click', closeMenu);
  }, []);

  const canManageVehicles = perms.can('vehicle_manage');
  const loadData = () => {
    setLoading(true)
    Promise.all([
      api.vehicles.list().catch(() => []),
      api.users.getVisible().catch(() => []),
      api.warehouses.list().catch(() => [])
    ]).then(([vehiclesData, usersData, warehousesData]) => {
      setVehicles(Array.isArray(vehiclesData) ? vehiclesData : [])
      setUsersList(Array.isArray(usersData) ? usersData : [])
      setWarehousesList(Array.isArray(warehousesData) ? warehousesData : [])
    }).finally(() => setLoading(false))
  }

  useEffect(() => {
    loadData()
  }, [])

  const handleCreateVehicle = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsProcessing(true)
    try {
      const payload = { 
        ...newVehicle, 
        driver_id: newVehicle.driver_id || undefined 
      }
      await api.vehicles.create(payload)
      toast.success('Автомобіль успішно додано!')
      setNewVehicle({ brand: '', model: '', plate_number: '', type: 'VAN', capacity_kg: 1500, tank_capacity: 0, fuel_norm: 0, driver_id: '', home_warehouse_id: '' })
      loadData()
      setShowVehicleForm(false)
    } catch (err: any) {
      // Перевіряємо, чи це помилка ліміту (402 Payment Required)
      if (err?.response?.status === 402 || err?.message?.includes('ліміт') || err?.message?.includes('Ліміт')) {
        const errorMsg = err?.response?.data?.error || err?.message || 'Досягнуто ліміт транспортних засобів для вашого тарифу';
        toast.error(
          (t) => (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <strong>🚫 {errorMsg}</strong>
              <button
                onClick={() => {
                  toast.dismiss(t.id);
                  navigate('/billing');
                }}
                style={{
                  padding: '8px 16px',
                  background: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)',
                  color: 'white',
                  border: 'none',
                  borderRadius: '6px',
                  cursor: 'pointer',
                  fontWeight: '600'
                }}
              >
                💎 Оновити тариф
              </button>
            </div>
          ),
          { duration: 8000 }
        );
      } else {
        toast.error(err?.message || 'Помилка створення авто');
      }
    } finally {
      setIsProcessing(false)
    }
  }

  const handleAddFuel = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!fuelModalVehicle) return
    setIsProcessing(true)

    try {
      const dataToSubmit = {
        record_type: fuelForm.record_type,
        liters: fuelForm.liters,
        odometer_km: fuelForm.odometer_km ? parseInt(fuelForm.odometer_km, 10) : undefined
      }

      const record = await api.vehicles.addFuelRecord(fuelModalVehicle.id, dataToSubmit)
      
      if (record.is_anomaly && hasFuelAntifraud) {
        setFuelModalVehicle(null)
        setFuelForm({ record_type: 'EXPENSE', liters: 0, odometer_km: '' })
        setAnomalyAlert(record)
        loadData()
      } else {
        toast.success('Запис про пальне успішно додано!')
        setFuelForm({ record_type: 'EXPENSE', liters: 0, odometer_km: '' })
        loadData()
        setFuelModalVehicle(null)
      }
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка збереження пального')
    } finally {
      setIsProcessing(false)
    }
  }

  // ⛽ Дозаправка під час активного рейсу
  const handleConfirmTripRefuel = async () => {
    if (!tripRefuelModal) return
    const liters = parseFloat(tripRefuelForm.liters)
    if (isNaN(liters) || liters <= 0) return toast.error('Вкажіть кількість літрів')
    setTripRefuelProcessing(true)
    try {
      const payload: LogShipmentRefuelPayload = { liters }
      if (tripRefuelForm.station_name.trim()) payload.station_name = tripRefuelForm.station_name.trim()
      if (tripRefuelForm.odometer_km) payload.odometer_km = parseInt(tripRefuelForm.odometer_km, 10)
      if (tripRefuelForm.cost_uah) payload.cost_uah = parseFloat(tripRefuelForm.cost_uah)
      const result = await api.inventory.logShipmentRefuel(tripRefuelModal.shipmentId, payload)
      toast.success(`⛽ Дозаправка ${result.liters} л зареєстрована!`)
      setTripRefuelModal(null)
      // Оновлюємо список рейсів в панелі авто
      if (historyVehicle) {
        api.vehicles.getShipmentHistory(historyVehicle.id).then(setTripRecords).catch(() => {})
      }
    } catch (err: any) {
      toast.error(err.message || 'Помилка реєстрації дозаправки')
    } finally {
      setTripRefuelProcessing(false)
    }
  }

  const handlePerformMaintenance = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!maintenanceModalVehicle) return
    
    if (!maintenanceForm.document) {
      toast.error('Помилка: Необхідно завантажити скан Акту (PDF або Фото)!')
      return
    }

    setIsProcessing(true)
    try {
      const formData = new FormData()
      formData.append('current_odometer', maintenanceForm.odometer_km.toString())
      formData.append('description', maintenanceForm.description)
      formData.append('performed_by', maintenanceForm.performed_by)
      formData.append('cost_amount', maintenanceForm.cost_amount.toString())
      
      if (maintenanceForm.document) {
        formData.append('document', maintenanceForm.document)
      }

      await api.vehicles.performMaintenance(maintenanceModalVehicle.id, formData)
      
      const msg = maintenanceModalVehicle.status === 'IN_REPAIR' 
        ? 'Ремонт завершено! Машина знову на лінії.' 
        : 'Акт ТО успішно зафіксовано!'
      
      toast.success(msg)
      loadData()
      setMaintenanceModalVehicle(null)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка фіксації акту')
    } finally {
      setIsProcessing(false)
    }
  }

  const handleUpdateStatus = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!statusModalVehicle) return
    setIsProcessing(true)

    try {
      await api.vehicles.updateStatus(statusModalVehicle.id, {
        status: statusForm.status,
        reason: statusForm.reason
      })
      toast.success('Статус машини оновлено!')
      loadData()
      setStatusModalVehicle(null)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка оновлення статусу')
    } finally {
      setIsProcessing(false)
    }
  }

  const handleAssignDriver = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!driverModalVehicle) return
    setIsProcessing(true)

    try {
      const driverIdToSubmit = driverForm.driver_id === '' ? null : driverForm.driver_id
      await api.vehicles.assignDriver(driverModalVehicle.id, driverIdToSubmit)
      toast.success('Водія успішно призначено!')
      loadData()
      setDriverModalVehicle(null)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка оновлення водія')
    } finally {
      setIsProcessing(false)
    }
  }

  const handleViewHistory = async (vehicle: Vehicle) => {
    setHistoryVehicle(vehicle)
    setHistoryTab('FUEL')
    setHistoryLoading(true)
    setFuelRecords([])
    setMaintenanceRecords([])
    setDriverRecords([])
    setTripRecords([])
    
    try {
      const [fRecords, mRecords, dRecords, tRecords] = await Promise.all([
        api.vehicles.getFuelHistory(vehicle.id).catch(() => []),
        api.vehicles.getMaintenanceHistory(vehicle.id).catch(() => []),
        api.vehicles.getDriverHistory(vehicle.id).catch(() => []),
        api.vehicles.getShipmentHistory(vehicle.id).catch(() => [])
      ])
      setFuelRecords(fRecords || [])
      setMaintenanceRecords(mRecords || [])
      setDriverRecords(dRecords || [])
      setTripRecords(tRecords || [])
    } catch (err) {
      toast.error('Не вдалося завантажити історію')
    } finally {
      setHistoryLoading(false)
    }
  }

  const handleOpenEdit = (v: Vehicle) => {
    setEditingVehicle(v);
    setEditForm({
      brand: v.brand,
      model: v.model || '',
      plate_number: v.plate_number,
      capacity_kg: v.capacity_kg
    });
  };

  const handleEditSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingVehicle) return;
    setIsProcessing(true);
    try {
      await api.vehicles.update(editingVehicle.id, editForm as Partial<Vehicle>);
      toast.success('Дані автомобіля оновлено!');
      setEditingVehicle(null);
      loadData();
    } catch (err: any) {
      toast.error(err.message || 'Помилка оновлення автомобіля');
    } finally {
      setIsProcessing(false);
    }
  };

  const handleDelete = async () => {
    if (!vehicleToDelete) return;
    setIsProcessing(true);
    try {
      await api.vehicles.delete(vehicleToDelete.id);
      toast.success('Автомобіль списано!');
      setVehicleToDelete(null);
      loadData();
    } catch (err: any) {
      toast.error(err.message || 'Неможливо списати (можливо, авто виконує рейс)');
    } finally {
      setIsProcessing(false);
    }
  };

  const getDriverName = (driverId?: string) => {
    if (!driverId) return <span className="text-muted unassigned-text-italic">Не призначено</span>;
    const driver = usersList.find(u => u.id === driverId);
    if (!driver) return <span className="text-muted">Невідомий ({driverId.substring(0, 6)})</span>;
    return <span>{driver.full_name}</span>;
  }

  // --- ПАГІНАЦІЯ (хуки мають бути до будь-якого раннього return) ---
  const VEHICLES_PAGE_SIZE = 20;
  const [vehiclesPage, setVehiclesPage] = useState(0);
  React.useEffect(() => { setVehiclesPage(0); }, [searchQuery, viewTab]);

  if (loading) {
    return (
      <div className="page-loading">
        <div className="spinner" />
        <p>Завантаження автопарку...</p>
      </div>
    )
  }

  // ---------------------------------------------------------
  // ЛОГІКА ФІЛЬТРАЦІЇ ПОШУКУ
  // ---------------------------------------------------------
  const filteredVehicles = vehicles.filter(v => {
    if (!searchQuery) return true;
    
    const query = searchQuery.toLowerCase().trim();
    const brandModel = `${v.brand} ${v.model || ''}`.toLowerCase();
    const plate = (v.plate_number || '').toLowerCase();
    
    // Шукаємо ім'я водія
    const driver = usersList.find(u => u.id === v.driver_id);
    const driverName = (driver?.full_name || '').toLowerCase();

    return brandModel.includes(query) || plate.includes(query) || driverName.includes(query);
  });

  const activeVehicles = filteredVehicles.filter(v => v.status !== 'INACTIVE')
  const archivedVehicles = filteredVehicles.filter(v => v.status === 'INACTIVE')
  const displayedVehicles = viewTab === 'ACTIVE' ? activeVehicles : archivedVehicles

  const vehiclesTotalPages = Math.max(1, Math.ceil(displayedVehicles.length / VEHICLES_PAGE_SIZE));
  const safeVehiclesPage = Math.min(vehiclesPage, vehiclesTotalPages - 1);
  const pagedVehicles = displayedVehicles.slice(
    safeVehiclesPage * VEHICLES_PAGE_SIZE,
    (safeVehiclesPage + 1) * VEHICLES_PAGE_SIZE
  );

  const assignedDriverIds = vehicles
    .filter(v => v.status !== 'INACTIVE' && v.driver_id)
    .map(v => v.driver_id);

  const availableDriversForNew = usersList.filter(u => !assignedDriverIds.includes(u.id));

  // Фільтруємо водіїв для присвоєння:
  // 1. Не зайнятий (або вже закріплений за цим авто)
  // 2. Належить до того ж підрозділу (unit_id), що і склад базування авто
  const vehicleWarehouseUnitId = driverModalVehicle
    ? warehousesList.find(w => w.id === driverModalVehicle.home_warehouse_id)?.unit_id
    : undefined;

  const availableDriversForAssign = usersList.filter(u => {
    const notTaken = !assignedDriverIds.includes(u.id) || u.id === driverModalVehicle?.driver_id;
    if (!notTaken) return false;
    // Якщо склад авто має unit_id — показуємо лише водіїв з того ж підрозділу
    if (vehicleWarehouseUnitId != null) {
      return u.unit_id === vehicleWarehouseUnitId;
    }
    return true; // склад без підрозділу — без фільтра
  });

  return (
    <div className="vehicles-page">

      <div className="page-header">
        <h1>Автопарк та ПММ</h1>
        <div className="page-actions">
          {canManageVehicles && (
            <button className="btn btn-primary" onClick={() => setShowVehicleForm(true)}>
              + Автомобіль
            </button>
          )}
        </div>
      </div>

      {/* МОДАЛКА: АНОМАЛІЯ ПАЛЬНОГО */}
      {anomalyAlert && (
        <div className="modal-overlay vehicles-modal">
          <div className="modal confirm-modal anomaly-modal">
            <h3>🚨 Виявлено аномалію!</h3>
            <p className="modal-description">
              Запис збережено в систему, але він позначений як підозрілий і потребує перевірки.
            </p>
            <div className="warning-box">
              <p style={{margin: '0 0 6px', fontWeight: 600}}>Причина від системи:</p>
              <p style={{margin: 0, color: '#dc2626'}}>{anomalyAlert.anomaly_reason}</p>
            </div>
            <div className="modal-actions center-actions">
              <button className="btn btn-primary" onClick={() => setAnomalyAlert(null)}>
                Зрозуміло
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ⛽ Модаль дозаправки в дорозі (з вкладки Автопарку → Рейси) */}
      {tripRefuelModal && (
        <div className="modal-overlay vehicles-modal" onClick={() => setTripRefuelModal(null)}>
          <div className="modal" style={{ maxWidth: '420px', width: '100%' }} onClick={e => e.stopPropagation()}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
              <h3 style={{ margin: 0 }}>⛽ Дозаправка в дорозі</h3>
              <button onClick={() => setTripRefuelModal(null)} style={{ background: 'none', border: 'none', fontSize: '24px', cursor: 'pointer', color: 'var(--text-muted)' }}>&times;</button>
            </div>
            <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginBottom: '16px', padding: '8px 12px', background: 'rgba(251,191,36,0.08)', border: '1px solid rgba(251,191,36,0.25)', borderRadius: '6px' }}>
              🚚 {tripRefuelModal.fromWarehouse} → {tripRefuelModal.toWarehouse}
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
              <div>
                <label className="form-label">Літри <span style={{ color: '#ef4444' }}>*</span></label>
                <input type="number" min="0.1" step="0.1" className="erp-input" style={{ width: '100%', boxSizing: 'border-box' }}
                  placeholder="напр. 35.5" value={tripRefuelForm.liters}
                  onChange={e => setTripRefuelForm(f => ({ ...f, liters: e.target.value }))} autoFocus />
              </div>
              <div>
                <label className="form-label">Назва АЗС (необов'язково)</label>
                <input type="text" className="erp-input" style={{ width: '100%', boxSizing: 'border-box' }}
                  placeholder="ОККО, WOG, ANP, Shell..." value={tripRefuelForm.station_name}
                  onChange={e => setTripRefuelForm(f => ({ ...f, station_name: e.target.value }))} />
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '10px' }}>
                <div>
                  <label className="form-label">Одометр (км)</label>
                  <input type="number" min="0" className="erp-input" style={{ width: '100%', boxSizing: 'border-box' }}
                    placeholder="12450" value={tripRefuelForm.odometer_km}
                    onChange={e => setTripRefuelForm(f => ({ ...f, odometer_km: e.target.value }))} />
                </div>
                <div>
                  <label className="form-label">Вартість (грн)</label>
                  <input type="number" min="0" step="0.01" className="erp-input" style={{ width: '100%', boxSizing: 'border-box' }}
                    placeholder="2100" value={tripRefuelForm.cost_uah}
                    onChange={e => setTripRefuelForm(f => ({ ...f, cost_uah: e.target.value }))} />
                </div>
              </div>
            </div>
            <div className="modal-actions" style={{ marginTop: '20px' }}>
              <button className="btn btn-secondary" onClick={() => setTripRefuelModal(null)}>Скасувати</button>
              <button
                className="btn"
                style={{ background: '#f59e0b', color: '#1a1a2e', fontWeight: 700 }}
                onClick={handleConfirmTripRefuel}
                disabled={tripRefuelProcessing}
              >
                {tripRefuelProcessing ? '⏳ Збереження...' : '⛽ Зареєструвати'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* МОДАЛКИ (Створення, Редагування, ТО, Пальне тощо) */}
      {showVehicleForm && canManageVehicles && (
        <div className="modal-overlay vehicles-modal" onClick={() => setShowVehicleForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Новий транспорт</h3>
            <form onSubmit={handleCreateVehicle}>
              <div className="form-group">
                <label>Марка (напр., Ford, Renault)</label>
                <input value={newVehicle.brand} onChange={(e) => setNewVehicle({ ...newVehicle, brand: e.target.value })} required className="erp-input" />
              </div>
              <div className="form-row-2">
                <div className="form-group">
                  <label>Модель</label>
                  <input value={newVehicle.model} onChange={(e) => setNewVehicle({ ...newVehicle, model: e.target.value })} className="erp-input" />
                </div>
                <div className="form-group">
                  <label>Номерний знак <span style={{color: '#ef4444'}}>*</span></label>
                  <input value={newVehicle.plate_number} onChange={(e) => setNewVehicle({ ...newVehicle, plate_number: e.target.value })} required className="erp-input" />
                </div>
              </div>

              <div className="form-row-2">
                <div className="form-group">
                  <label>Тип кузова <span style={{color: '#ef4444'}}>*</span></label>
                  <select value={newVehicle.type} onChange={(e) => setNewVehicle({...newVehicle, type: e.target.value})} className="erp-input">
                    <option value="PICKUP">🛻 Пікап / Джип</option>
                    <option value="VAN">🚐 Мікроавтобус / Фургон</option>
                    <option value="TRUCK">🚛 Вантажівка</option>
                  </select>
                </div>
                <div className="form-group">
                  <label>Вантажопідйомність (кг) <span style={{color: '#ef4444'}}>*</span></label>
                  <input type="number" min="1" value={newVehicle.capacity_kg || ''} onChange={(e) => setNewVehicle({ ...newVehicle, capacity_kg: parseFloat(e.target.value) })} required className="erp-input" />
                </div>
              </div>

              <div className="form-row-2">
                <div className="form-group">
                  <label>Об'єм бака (л) <span style={{color: '#ef4444'}}>*</span></label>
                  <input type="number" min="1" value={newVehicle.tank_capacity || ''} onChange={(e) => setNewVehicle({ ...newVehicle, tank_capacity: parseFloat(e.target.value) })} required className="erp-input" />
                </div>
                <div className="form-group">
                  <label>Норма (л/100км) <span style={{color: '#ef4444'}}>*</span></label>
                  <input type="number" min="1" step="0.1" value={newVehicle.fuel_norm || ''} onChange={(e) => setNewVehicle({ ...newVehicle, fuel_norm: parseFloat(e.target.value) })} required className="erp-input" />
                </div>
              </div>

              <div className="form-group">
                <label>Закріплений водій</label>
                <SearchableSelect
                  options={availableDriversForNew.map(u => ({ value: u.id, label: `${u.full_name} (${ROLE_NAMES[u.role] || u.role})` }))}
                  value={newVehicle.driver_id}
                  onChange={(val) => setNewVehicle({ ...newVehicle, driver_id: val })}
                  emptyLabel="-- Без закріплення --"
                  searchPlaceholder="Пошук співробітника..."
                  disabled={isProcessing}
                />
                <span style={{ fontSize: '11px', color: '#64748b', marginTop: '4px', display: 'block' }}>
                  У списку відображаються лише вільні співробітники.
                </span>
              </div>

              <div className="form-group">
                <label>Базовий склад <span style={{color: '#ef4444'}}>*</span></label>
                <SearchableSelect
                  options={warehousesList.map(w => ({ value: w.id, label: w.name }))}
                  value={newVehicle.home_warehouse_id}
                  onChange={(val) => setNewVehicle({ ...newVehicle, home_warehouse_id: val })}
                  placeholder="Оберіть склад приписки..."
                  searchPlaceholder="Пошук складу..."
                  disabled={isProcessing}
                />
                <span style={{ fontSize: '11px', color: '#64748b', marginTop: '4px', display: 'block' }}>
                  Де машина базується постійно. Поточна локація може змінюватись після рейсів.
                </span>
              </div>

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowVehicleForm(false)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={isProcessing || !newVehicle.brand?.trim() || !newVehicle.plate_number?.trim() || !newVehicle.capacity_kg || !newVehicle.tank_capacity || !newVehicle.fuel_norm || !newVehicle.home_warehouse_id}>Створити</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {editingVehicle && canManageVehicles && (
        <div className="modal-overlay vehicles-modal" onClick={() => !isProcessing && setEditingVehicle(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Редагувати параметри авто</h3>
            <form onSubmit={handleEditSubmit}>
              <div className="form-group">
                <label>Державний номер</label>
                <input className="erp-input" value={editForm.plate_number} onChange={(e) => setEditForm({ ...editForm, plate_number: e.target.value })} required disabled={isProcessing} />
              </div>
              <div className="form-row-2">
                <div className="form-group">
                  <label>Марка</label>
                  <input className="erp-input" value={editForm.brand} onChange={(e) => setEditForm({ ...editForm, brand: e.target.value })} required disabled={isProcessing} />
                </div>
                <div className="form-group">
                  <label>Модель</label>
                  <input className="erp-input" value={editForm.model} onChange={(e) => setEditForm({ ...editForm, model: e.target.value })} disabled={isProcessing} />
                </div>
              </div>
              <div className="form-group">
                <label>Вантажопідйомність (кг)</label>
                <input className="erp-input" type="number" min="1" value={editForm.capacity_kg} onChange={(e) => setEditForm({ ...editForm, capacity_kg: Number(e.target.value) })} required disabled={isProcessing} />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setEditingVehicle(null)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={isProcessing}>{isProcessing ? 'Збереження...' : 'Зберегти'}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {vehicleToDelete && canManageVehicles && (
        <div className="modal-overlay vehicles-modal" onClick={() => !isProcessing && setVehicleToDelete(null)}>
          <div className="modal confirm-modal" onClick={(e) => e.stopPropagation()}>
            <h3 style={{ color: '#ef4444' }}>⚠️ Списання автомобіля</h3>
            <p className="modal-description" style={{color: '#0f172a'}}>Ви впевнені, що хочете остаточно списати авто <strong>{vehicleToDelete.brand} ({vehicleToDelete.plate_number})</strong>?</p>
            <p className="modal-description">Це видалить його зі списку доступних машин для формування логістичних рейсів.</p>
            <div className="modal-actions">
              <button className="btn btn-secondary" onClick={() => setVehicleToDelete(null)} disabled={isProcessing}>Скасувати</button>
              <button className="btn btn-danger" onClick={handleDelete} disabled={isProcessing}>{isProcessing ? 'Списання...' : 'Списати'}</button>
            </div>
          </div>
        </div>
      )}

      {statusModalVehicle && (
        <div className="modal-overlay vehicles-modal" onClick={() => !isProcessing && setStatusModalVehicle(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>🛠 Відправити в ремонт: {statusModalVehicle.brand} ({statusModalVehicle.plate_number})</h3>
            <form onSubmit={handleUpdateStatus}>
              <p className="modal-description">Заповніть причину — це допоможе відстежити хід ремонту та скласти акт виконаних робіт.</p>
              <div className="form-group">
                <label>Причина направлення в ремонт <span style={{color: '#ef4444'}}>*</span></label>
                <textarea rows={3} placeholder="Напр., Поломка двигуна / ДТП / планове ТО" value={statusForm.reason} onChange={(e) => setStatusForm({...statusForm, reason: e.target.value})} required disabled={isProcessing} className="erp-input" />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setStatusModalVehicle(null)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className="btn btn-danger" disabled={isProcessing || !statusForm.reason?.trim()}>Підтвердити дію</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {driverModalVehicle && (
        <div className="modal-overlay vehicles-modal" onClick={() => !isProcessing && setDriverModalVehicle(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Водій: {driverModalVehicle.brand} ({driverModalVehicle.plate_number})</h3>
            <form onSubmit={handleAssignDriver}>
              <p className="modal-description">
                Виберіть співробітника, за яким буде закріплено даний транспортний засіб.
                {vehicleWarehouseUnitId != null && (
                  <span style={{ display: 'block', marginTop: '6px', color: 'var(--text-muted)', fontSize: '12px' }}>
                    🏠 Показуються лише співробітники підрозділу складу базування авто.
                  </span>
                )}
              </p>
              <div className="form-group">
                <label>Відповідальний водій</label>
                <SearchableSelect
                  options={availableDriversForAssign.map(u => ({ value: u.id, label: `${u.full_name} (${ROLE_NAMES[u.role] || u.role})` }))}
                  value={driverForm.driver_id}
                  onChange={(val) => setDriverForm({...driverForm, driver_id: val})}
                  emptyLabel="-- Зняти закріплення (Без водія) --"
                  searchPlaceholder="Пошук співробітника..."
                  disabled={isProcessing}
                />
                {availableDriversForAssign.length === 0 && vehicleWarehouseUnitId != null && (
                  <div style={{ marginTop: '8px', fontSize: '12px', color: '#f59e0b' }}>
                    ⚠️ Немає доступних співробітників у підрозділі цього складу. Спочатку призначте персонал до підрозділу.
                  </div>
                )}
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setDriverModalVehicle(null)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={isProcessing}>Зберегти зміни</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {fuelModalVehicle && (
        <div className="modal-overlay vehicles-modal" onClick={() => !isProcessing && setFuelModalVehicle(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Пальне: {fuelModalVehicle.brand} ({fuelModalVehicle.plate_number})</h3>
            <form onSubmit={handleAddFuel}>
              <div className="form-group">
                <label>Тип операції</label>
                <select value={fuelForm.record_type} onChange={(e) => setFuelForm({ ...fuelForm, record_type: e.target.value as FuelRecordType })} className={`erp-input fuel-type-select ${fuelForm.record_type === 'REFUEL' ? 'type-refuel' : 'type-expense'}`} disabled={isProcessing}>
                  <option value="EXPENSE">Списання (Витрата)</option>
                  <option value="REFUEL">Заправка (Прихід)</option>
                </select>
              </div>
              <div className="form-row-2">
                <div className="form-group">
                  <label>Літри <span style={{color: '#ef4444'}}>*</span></label>
                  <input type="number" min="0.1" step="0.1" value={fuelForm.liters || ''} onChange={(e) => setFuelForm({ ...fuelForm, liters: parseFloat(e.target.value) })} required disabled={isProcessing} className="erp-input" />
                </div>
                <div className="form-group">
                  <label>Поточний одометр (км) {fuelForm.record_type === 'EXPENSE' && <span style={{color: '#ef4444'}}>*</span>}</label>
                  <input type="number" min={fuelModalVehicle?.current_odometer || 0} placeholder={`Напр. ${fuelModalVehicle?.current_odometer ? fuelModalVehicle.current_odometer + 150 : 150500}`} value={fuelForm.odometer_km} onChange={(e) => setFuelForm({ ...fuelForm, odometer_km: e.target.value })} required={fuelForm.record_type === 'EXPENSE'} disabled={isProcessing} className="erp-input" />
                  <span style={{ display: 'block', fontSize: '11px', color: '#64748b', marginTop: '6px' }}>Останній запис: <strong>{fuelModalVehicle?.current_odometer || 0}</strong> км</span>
                </div>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setFuelModalVehicle(null)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className={fuelForm.record_type === 'REFUEL' ? "btn btn-success" : "btn btn-danger"} disabled={isProcessing || !fuelForm.liters || fuelForm.liters <= 0 || (fuelForm.record_type === 'EXPENSE' && !fuelForm.odometer_km)}>{fuelForm.record_type === 'REFUEL' ? 'Заправити' : 'Списати'}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {maintenanceModalVehicle && (
        <div className="modal-overlay vehicles-modal" onClick={() => !isProcessing && setMaintenanceModalVehicle(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>{maintenanceModalVehicle.status === 'IN_REPAIR' ? '✅ Завершення ремонту: ' : '🛠 Акт виконаних робіт: '} {maintenanceModalVehicle.brand} ({maintenanceModalVehicle.plate_number})</h3>
            <form onSubmit={handlePerformMaintenance}>
              <div className="form-row-2">
                <div className="form-group">
                  <label>Виконавець (Власний сервіс / СТО)</label>
                  <input type="text" placeholder="Напр. Внутрішній сервіс або СТО 'Гараж'" value={maintenanceForm.performed_by} onChange={(e) => setMaintenanceForm({...maintenanceForm, performed_by: e.target.value})} disabled={isProcessing} className="erp-input" />
                </div>
              </div>
              <div className="form-group">
                <label>Опис робіт та замінених запчастин <span style={{color: '#ef4444'}}>*</span></label>
                <textarea rows={3} placeholder="Заміна мастила 5w40 (5л), масляний фільтр..." value={maintenanceForm.description} onChange={(e) => setMaintenanceForm({...maintenanceForm, description: e.target.value})} required disabled={isProcessing} className="erp-input" />
              </div>
              <div className="form-row-2 form-row-bottom-align">
                <div className="form-group" style={{marginBottom: 0}}>
                  <label>Загальна вартість (Грн)</label>
                  <input type="number" min="0" step="0.01" placeholder="Напр. 15000" value={maintenanceForm.cost_amount || ''} onChange={(e) => setMaintenanceForm({...maintenanceForm, cost_amount: parseFloat(e.target.value)})} disabled={isProcessing} className="erp-input" />
                </div>
                <div className="form-group" style={{marginBottom: 0}}>
                  <label>Скан Акту (PDF/Фото) <span style={{color: '#ef4444'}}>*</span></label>
                  <label className="file-upload-custom">
                  <input type="file" ref={maintenanceFileRef} className="file-input-hidden" accept="image/*,application/pdf" onChange={(e) => { if (e.target.files && e.target.files.length > 0) setMaintenanceForm({...maintenanceForm, document: e.target.files[0]}) }} disabled={isProcessing} />
                    <span className="file-upload-text">{maintenanceForm.document ? `📎 ${maintenanceForm.document.name}` : '📁 Натисніть, щоб вибрати...'}</span>
                  </label>
                </div>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setMaintenanceModalVehicle(null)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className={maintenanceModalVehicle.status === 'IN_REPAIR' ? "btn btn-success" : "btn btn-primary"} disabled={isProcessing || !maintenanceForm.description?.trim() || !maintenanceForm.performed_by?.trim()}>{maintenanceModalVehicle.status === 'IN_REPAIR' ? 'Повернути на лінію' : 'Зафіксувати ТО'}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {historyVehicle && (
        <div className="modal-overlay vehicles-modal" onClick={() => setHistoryVehicle(null)}>
          <div className="modal history-modal" onClick={(e) => e.stopPropagation()}>
            <h3>Картка автомобіля: {historyVehicle.brand} ({historyVehicle.plate_number})</h3>
            
            <div className="vehicle-info-cards">
              <div className="v-card">
                <span className="v-card-label">Номерний знак</span>
                <span className="v-card-value plate-badge">{historyVehicle.plate_number}</span>
              </div>
              <div className="v-card">
                <span className="v-card-label">Тип і Вантаж</span>
                <span className="v-card-value">
                  {historyVehicle.type === 'PICKUP' ? '🛻 Пікап' : historyVehicle.type === 'VAN' ? '🚐 Фургон' : '🚛 Вантажівка'}
                  <small>({historyVehicle.capacity_kg} кг)</small>
                </span>
              </div>
              <div className="v-card">
                <span className="v-card-label">Статус</span>
                <span className={`badge ${historyVehicle.status === 'ACTIVE' ? 'badge-success' : historyVehicle.status === 'IN_REPAIR' ? 'badge-warning' : 'badge-critical'}`}>
                  {historyVehicle.status === 'ACTIVE' ? 'На лінії' : historyVehicle.status === 'IN_REPAIR' ? 'В ремонті' : historyVehicle.status === 'ON_MISSION' ? 'У рейсі' : 'Списане'}
                </span>
              </div>
              
              {(() => {
                const currentFuel = fuelRecords.reduce((total, record) => {
                  return record.record_type === 'REFUEL' ? total + record.liters : total - record.liters;
                }, 0);
                
                const displayFuel = Math.max(0, currentFuel); 
                const tankCap = historyVehicle.tank_capacity || 1;
                const fillPercentage = Math.min(100, Math.round((displayFuel / tankCap) * 100));

                return (
                  <div className="v-card fuel-gauge-card">
                    <div className="fuel-gauge-header">
                      <span className="v-card-label">Залишок у баку</span>
                      <span className="v-card-value fuel-gauge-value-text">
                        {displayFuel.toFixed(1)} / {tankCap} <small>л</small>
                      </span>
                    </div>
                    <div className="gauge-bg">
                      <div className={`gauge-fill ${fillPercentage < 20 ? 'gauge-low' : 'gauge-normal'}`} style={{ width: `${fillPercentage}%` }}></div>
                    </div>
                  </div>
                );
              })()}
            </div>
            
            <div className="history-tabs">
              <button className={`history-tab ${historyTab === 'FUEL' ? 'active' : ''}`} onClick={() => setHistoryTab('FUEL')}>⛽ Історія пального</button>
              <button className={`history-tab ${historyTab === 'MAINTENANCE' ? 'active' : ''}`} onClick={() => setHistoryTab('MAINTENANCE')}>🛠 Акти виконаних робіт</button>
              <button className={`history-tab ${historyTab === 'DRIVERS' ? 'active' : ''}`} onClick={() => setHistoryTab('DRIVERS')}>👥 Історія водіїв</button>
              <button className={`history-tab ${historyTab === 'TRIPS' ? 'active' : ''}`} onClick={() => setHistoryTab('TRIPS')}>🚚 Історія рейсів</button>
            </div>

            {historyLoading ? (
              <div className="history-spinner">Завантаження даних...</div>
            ) : (
              <div className="history-table-wrapper">
                
                {historyTab === 'FUEL' && (
                  fuelRecords.length === 0 ? (
                    <p className="history-empty">Записів про пальне ще немає.</p>
                  ) : (
                    <table className="data-table fuel-history-table">
                      <thead>
                        <tr><th>Дата</th><th>Тип</th><th>Літри</th><th>Одометр</th><th>Статус</th></tr>
                      </thead>
                      <tbody>
                        {fuelRecords.map(record => (
                          <tr key={record.id} className={record.is_anomaly ? 'row-inactive' : ''}>
                            <td className="text-muted">
                              {new Date(record.created_at).toLocaleString('uk-UA', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })}
                            </td>
                            <td className={`fuel-type-select ${record.record_type === 'REFUEL' ? 'type-refuel' : 'type-expense'}`}>
                              {record.record_type === 'REFUEL' ? 'Прихід' : 'Списання'}
                            </td>
                            <td style={{fontWeight: 600}}>{record.liters} л</td>
                            <td style={{fontWeight: 500}}>{record.odometer_km ? `${record.odometer_km} км` : '-'}</td>
                            <td>
                              {record.is_anomaly ? (
                                <span className="badge badge-critical" title={record.anomaly_reason || 'Підозрілий запис'}>Аномалія ⚠️</span>
                              ) : (
                                <span className="badge badge-success">ОК</span>
                              )}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  )
                )}

                {historyTab === 'MAINTENANCE' && (
                  maintenanceRecords.length === 0 ? (
                    <p className="history-empty">Актів виконаних робіт ще немає.</p>
                  ) : (
                    <table className="data-table maint-history-table">
                      <thead>
                        <tr>
                          <th>Дата ТО</th>
                          <th>Водій</th>
                          <th>Одометр</th>
                          <th>Опис виконаних робіт</th>
                          <th>Виконавець</th>
                          <th>Вартість</th>
                          <th>Документ</th>
                        </tr>
                      </thead>
                      <tbody>
                        {maintenanceRecords.map(record => (
                          <tr key={record.id}>
                            <td className="text-muted">
                              {new Date(record.created_at).toLocaleString('uk-UA', { day: '2-digit', month: '2-digit', year: 'numeric' })}
                            </td>
                            <td className="text-muted">
                              {record.driver_name ? `👤 ${record.driver_name}` : <span className="unassigned-text-italic">Не призначено</span>}
                            </td>
                            <td style={{fontWeight: 500}}>{record.odometer_km} км</td>
                            <td className="desc-cell">{record.description}</td>
                            <td style={{color: '#475569'}}>{record.performed_by || '-'}</td>
                            <td style={{fontWeight: 600, color: '#0f172a'}}>
                              {record.cost_amount > 0 ? `${record.cost_amount} ₴` : '-'}
                            </td>
                            <td>
                              {record.document_url ? (
                                <a 
                                  href={record.document_url} 
                                  target="_blank" 
                                  rel="noopener noreferrer" 
                                  className="document-link-btn"
                                >
                                  📄 Відкрити
                                </a>
                              ) : '-'}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  )
                )}

                {historyTab === 'DRIVERS' && (
                  driverRecords.length === 0 ? (
                    <p className="history-empty">Історія водіїв порожня.</p>
                  ) : (
                    <table className="data-table">
                      <thead>
                        <tr>
                          <th>Дата призначення</th>
                          <th>Співробітник (Водій)</th>
                          <th>Стан</th>
                        </tr>
                      </thead>
                      <tbody>
                        {driverRecords.map((record, idx) => (
                          <tr key={record.id} className={idx === 0 && record.driver_id ? 'row-current-driver' : ''}>
                            <td className="text-muted">
                              {new Date(record.assigned_at).toLocaleString('uk-UA', { 
                                day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' 
                              })}
                            </td>
                            <td className={record.driver_id ? 'driver-assigned' : 'driver-unassigned'}>
                              {record.driver_id ? `👤 ${record.driver_name}` : `🚫 ${record.driver_name || 'Знято закріплення'}`}
                            </td>
                            <td>
                              {idx === 0 && record.driver_id ? (
                                <span className="badge badge-success">Поточний</span>
                              ) : idx === 0 && !record.driver_id ? (
                                <span className="badge badge-critical">Без водія</span>
                              ) : (
                                <span className="badge" style={{background:'#e2e8f0',color:'#64748b'}}>Архів</span>
                              )}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  )
                )}

                {historyTab === 'TRIPS' && (
                  tripRecords.length === 0 ? (
                    <p className="history-empty">Рейсів ще не було.</p>
                  ) : (
                    <table className="data-table">
                      <thead>
                        <tr>
                          <th>Дата</th>
                          <th>Звідки</th>
                          <th>Куди</th>
                          <th className="text-center">Статус</th>
                          <th className="text-center">Відстань</th>
                          <th className="text-center">Дії</th>
                        </tr>
                      </thead>
                      <tbody>
                        {tripRecords.map(trip => (
                          <tr key={trip.id}>
                            <td className="text-muted">
                              {new Date(trip.created_at).toLocaleString('uk-UA', {
                                day: '2-digit', month: '2-digit', year: 'numeric'
                              })}
                            </td>
                            <td style={{fontWeight: 500}}>{trip.from_warehouse_name}</td>
                            <td style={{fontWeight: 500}}>{trip.to_warehouse_name}</td>
                            <td className="text-center">
                              {trip.status === 'PENDING' ? (
                                <span className="badge badge-neutral">⏳ Очікує</span>
                              ) : trip.status === 'IN_TRANSIT' ? (
                                <span className="badge badge-warning">🚛 В дорозі</span>
                              ) : (
                                <span className="badge badge-success">✅ Доставлено</span>
                              )}
                            </td>
                            <td className="text-center" style={{fontWeight: 600}}>
                              {trip.actual_km > 0
                                ? `${trip.actual_km} км`
                                : trip.distance_km > 0
                                  ? `~${Math.round(trip.distance_km)} км`
                                  : '—'}
                            </td>
                            <td className="text-center">
                              {trip.status === 'IN_TRANSIT' && (
                                <button
                                  className="btn btn-sm"
                                  style={{ background: '#f59e0b', color: '#1a1a2e', fontWeight: 600, border: 'none', padding: '4px 10px', borderRadius: '6px', cursor: 'pointer', fontSize: '12px' }}
                                  onClick={() => {
                                    setTripRefuelModal({ shipmentId: trip.id, fromWarehouse: trip.from_warehouse_name, toWarehouse: trip.to_warehouse_name })
                                    setTripRefuelForm({ liters: '', station_name: '', odometer_km: '', cost_uah: '' })
                                  }}
                                  title="Зареєструвати дозаправку"
                                >⛽ Дозаправка</button>
                              )}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  )
                )}
              </div>
            )}

            <div className="modal-actions">
              <button className="btn btn-secondary" onClick={() => setHistoryVehicle(null)}>Закрити</button>
            </div>
          </div>
        </div>
      )}

      {/* ТАБЛИЦЯ АВТОПАРКУ */}
      <div className="card">
        <div className="table-header-flex" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '16px', marginBottom: '16px' }}>
          <h2 style={{ margin: 0 }}>{viewTab === 'ACTIVE' ? 'Автомобілі на балансі' : 'Архів списаної техніки'}</h2>
          
          {/* НОВЕ ПОЛЕ ПОШУКУ */}
          <div style={{ position: 'relative', width: '260px' }}>
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

          <div className="table-header-actions" style={{ display: 'flex', gap: '8px' }}>
            <button className={`tab-toggle-btn ${viewTab === 'ACTIVE' ? 'active' : ''}`} onClick={() => setViewTab('ACTIVE')}>Активні авто ({activeVehicles.length})</button>
            <button className={`tab-toggle-btn archive ${viewTab === 'ARCHIVE' ? 'active' : ''}`} onClick={() => setViewTab('ARCHIVE')}>Списані ({archivedVehicles.length})</button>
          </div>
        </div>

        {displayedVehicles.length === 0 ? (
          <p className="history-empty" style={{marginTop: '20px'}}>
            {searchQuery ? `За запитом "${searchQuery}" нічого не знайдено` : (viewTab === 'ACTIVE' ? 'Активний автопарк порожній' : 'Немає списаної техніки')}
          </p>
        ) : (
          <>
          <table className="data-table">
            <thead>
              <tr>
                <th>Марка / Модель</th>
                <th>Номерний знак</th>
                <th>Тип та Вантаж</th>
                <th>Закріплений водій</th>
                <th>База / Локація</th>
                <th>Бак (Норма)</th>
                <th>Статус</th>
                {viewTab === 'ACTIVE' && <th>До ТО</th>}
                <th className="col-actions-menu">Дії</th>
              </tr>
            </thead>
            <tbody>
              {pagedVehicles.map((v) => {
                return (
                  <tr key={v.id} className={v.status === 'INACTIVE' ? 'row-inactive' : ''}>
                    <td className="vehicle-brand-cell">
                      {v.brand} {v.model && <span className="vehicle-model-text">{v.model}</span>}
                    </td>
                    <td><span className="plate-badge">{v.plate_number}</span></td>
                    
                    <td>
                      <div>{v.type === 'PICKUP' ? '🛻 Пікап' : v.type === 'VAN' ? '🚐 Фургон' : v.type === 'TRUCK' ? '🚛 Вантажівка' : '🚗 Авто'}</div>
                      <div className="norm-text" style={{ marginLeft: 0 }}>Макс: {v.capacity_kg} кг</div>
                    </td>

                    <td>{getDriverName(v.driver_id)}</td>

                    <td>
                      <div style={{ fontSize: '0.8rem' }}>
                        <div title="Базовий склад (постійна приписка)">
                          🏠 {v.home_warehouse_name || <span style={{color: 'var(--text-muted)'}}>не вказано</span>}
                        </div>
                        {v.current_warehouse_id !== v.home_warehouse_id && (
                          <div style={{ color: 'var(--text-muted)', marginTop: 2 }} title="Поточна локація після рейсу">
                            📍 {v.current_warehouse_name || '—'}
                          </div>
                        )}
                      </div>
                    </td>

                    <td>{v.tank_capacity} л <span className="norm-text">({v.fuel_norm} л/100км)</span></td>
                    
                    <td>
                      <span className={`badge ${v.status === 'ACTIVE' ? 'badge-success' : v.status === 'IN_REPAIR' ? 'badge-warning' : v.status === 'ON_MISSION' ? 'badge-primary' : 'badge-critical'}`}>
                        {v.status === 'ACTIVE' ? 'На лінії' : v.status === 'IN_REPAIR' ? 'В ремонті' : v.status === 'ON_MISSION' ? 'У рейсі' : 'Списане'}
                      </span>
                    </td>
                    
                    {viewTab === 'ACTIVE' && (
  <td>
    {v.maintenance_status === 'OVERDUE' ? (
      <span className="badge badge-critical">🚨 Прострочено</span>
    ) : (
      <div className="maint-cell">
        {/* Завжди показуємо км залишку */}
        <div className={isPro && v.predicted_maint_date ? "pro-prediction" : "basic-maint"} title={isPro ? `Середній пробіг: ${Math.round(v.avg_km_per_day || 0)} км/день` : undefined}>
          <span className={v.maintenance_status === 'WARNING' ? 'text-warning' : 'text-muted'} style={{fontWeight: 600}}>
            Залишок: {v.km_to_next_maintenance} км
          </span>
          {/* PRO: дата прогнозу як додаткова підказка */}
          {isPro && v.predicted_maint_date ? (
            <>
              <div style={{fontSize: '11px', color: 'var(--text-muted)', marginTop: '2px'}}>
                📅 {new Date(v.predicted_maint_date).toLocaleDateString('uk-UA')}
              </div>
              <div className="pro-badge">PRO прогноз</div>
            </>
          ) : !isPro ? (
            <div className="upsell-link" onClick={() => toast.error("Прогноз дати доступний лише у PRO тарифі")}>
              🔒 Дізнатись дату
            </div>
          ) : null}
        </div>
      </div>
    )}
  </td>
)}

                    <td className="col-actions-menu">
                      <div className="dropdown-container" onClick={(e) => e.stopPropagation()}>
                        <button 
                          className={`btn-kebab ${activeMenuId === v.id ? 'active' : ''}`} 
                          onClick={() => setActiveMenuId(activeMenuId === v.id ? null : v.id)}
                        >
                          ⋮
                        </button>
                        
                        {activeMenuId === v.id && (
                          <div className="actions-dropdown-menu">
                            <button onClick={() => { handleViewHistory(v); setActiveMenuId(null); }}>
                              📊 Картка авто
                            </button>

                            {v.status === 'ON_MISSION' && (
                              <>
                                <button
                                  style={{ color: '#f59e0b', fontWeight: 600 }}
                                  onClick={async () => {
                                    setActiveMenuId(null)
                                    try {
                                      const trips = await api.vehicles.getShipmentHistory(v.id)
                                      const active = trips.find(t => t.status === 'IN_TRANSIT')
                                      if (active) {
                                        setTripRefuelForm({ liters: '', station_name: '', odometer_km: '', cost_uah: '' })
                                        setTripRefuelModal({ shipmentId: active.id, fromWarehouse: active.from_warehouse_name, toWarehouse: active.to_warehouse_name })
                                      } else {
                                        toast.error('Активний рейс не знайдено')
                                      }
                                    } catch {
                                      toast.error('Не вдалося отримати рейс')
                                    }
                                  }}
                                >
                                  ⛽ Дозаправка в дорозі
                                </button>
                                <div className="dropdown-divider"></div>
                              </>
                            )}
                            
                            {canManageVehicles && viewTab === 'ACTIVE' && v.status !== 'ON_MISSION' && (
                              <>
                                <button onClick={() => { setDriverModalVehicle(v); setDriverForm({ driver_id: v.driver_id || '' }); setActiveMenuId(null); }}>
                                  👤 Призначити водія
                                </button>
                                
                                {v.status === 'ACTIVE' && (
                                  <button style={{color: '#2563eb'}} onClick={() => { setFuelModalVehicle(v); setFuelForm({ record_type: 'EXPENSE', liters: 0, odometer_km: '' }); setActiveMenuId(null); }}>
                                    ⛽ Додати пальне
                                  </button>
                                )}

                                {v.status === 'IN_REPAIR' ? (
                                  <button className="text-success" onClick={() => { setMaintenanceModalVehicle(v); setMaintenanceForm({ odometer_km: v.current_odometer || 0, description: '', performed_by: '', cost_amount: 0, document: null }); if (maintenanceFileRef.current) maintenanceFileRef.current.value = ''; setActiveMenuId(null); }}>
                                    ✅ Завершити ремонт
                                  </button>
                                ) : (
                                  <button onClick={() => { setMaintenanceModalVehicle(v); setMaintenanceForm({ odometer_km: v.current_odometer || 0, description: '', performed_by: '', cost_amount: 0, document: null }); if (maintenanceFileRef.current) maintenanceFileRef.current.value = ''; setActiveMenuId(null); }}>
                                    🛠 Зафіксувати ТО
                                  </button>
                                )}
                                
                                {v.status !== 'IN_REPAIR' && (
                                  <button onClick={() => { setStatusModalVehicle(v); setStatusForm({ status: 'IN_REPAIR', reason: '' }); setActiveMenuId(null); }}>
                                    � Відправити в ремонт
                                  </button>
                                )}

                                <div className="dropdown-divider"></div>

                                <button onClick={() => { handleOpenEdit(v); setActiveMenuId(null); }}>
                                  ✏️ Редагувати параметри
                                </button>
                                <button className="text-danger" onClick={() => { setVehicleToDelete(v); setActiveMenuId(null); }}>
                                  🗑️ Списати авто
                                </button>
                              </>
                            )}
                          </div>
                        )}
                      </div>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
          <Pagination
            currentPage={safeVehiclesPage}
            totalPages={vehiclesTotalPages}
            onPageChange={setVehiclesPage}
            totalItems={displayedVehicles.length}
            itemsPerPage={VEHICLES_PAGE_SIZE}
          />
          </>
        )}
      </div>
    </div>
  )
}