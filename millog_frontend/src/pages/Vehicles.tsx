import React, { useEffect, useState } from 'react'
import { api, type Vehicle, type FuelRecordType, type FuelRecord, type MaintenanceRecord, type SystemUser, type DriverHistoryRecord} from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import toast, { Toaster } from 'react-hot-toast'
import './Vehicles.css' 

export default function Vehicles() {
  const { user } = useAuth()
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [usersList, setUsersList] = useState<SystemUser[]>([])
  const [loading, setLoading] = useState(true)

  const [viewTab, setViewTab] = useState<'ACTIVE' | 'ARCHIVE'>('ACTIVE')

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
  const [historyTab, setHistoryTab] = useState<'FUEL' | 'MAINTENANCE' | 'DRIVERS'>('FUEL')
  const [fuelRecords, setFuelRecords] = useState<FuelRecord[]>([])
  const [maintenanceRecords, setMaintenanceRecords] = useState<MaintenanceRecord[]>([])
  const [historyLoading, setHistoryLoading] = useState(false)
  const [driverRecords, setDriverRecords] = useState<DriverHistoryRecord[]>([])
  
  const [driverModalVehicle, setDriverModalVehicle] = useState<Vehicle | null>(null)
  const [driverForm, setDriverForm] = useState({ driver_id: '' })
  
  const [newVehicle, setNewVehicle] = useState({
    brand: '',
    model: '',
    plate_number: '',
    type: 'VAN',           
    capacity_kg: 1500,     
    tank_capacity: 0,
    fuel_norm: 0,
    driver_id: ''
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

  useEffect(() => {
    const closeMenu = () => setActiveMenuId(null);
    document.addEventListener('click', closeMenu);
    return () => document.removeEventListener('click', closeMenu);
  }, []);

  const canManageVehicles = ['ADMIN', 'REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN', 'DEPT_SUPERVISOR'].includes(user?.role || '')

  const loadData = () => {
    setLoading(true)
    Promise.all([
      api.vehicles.list().catch(() => []),
      api.users.getVisible().catch(() => [])
    ]).then(([vehiclesData, usersData]) => {
      setVehicles(Array.isArray(vehiclesData) ? vehiclesData : [])
      setUsersList(Array.isArray(usersData) ? usersData : [])
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
      await api.vehicles.create(payload as any) 
      toast.success('Автомобіль успішно додано!')
      setNewVehicle({ brand: '', model: '', plate_number: '', type: 'VAN', capacity_kg: 1500, tank_capacity: 0, fuel_norm: 0, driver_id: '' })
      loadData()
      setShowVehicleForm(false)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка створення авто')
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
      
      if (record.is_anomaly) {
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

      await api.vehicles.performMaintenance(maintenanceModalVehicle.id, formData as any)
      
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
        status: statusForm.status as any,
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
    
    try {
      const [fRecords, mRecords, dRecords] = await Promise.all([
        api.vehicles.getFuelHistory(vehicle.id).catch(() => []),
        api.vehicles.getMaintenanceHistory(vehicle.id).catch(() => []),
        api.vehicles.getDriverHistory(vehicle.id).catch(() => [])
      ])
      setFuelRecords(fRecords || [])
      setMaintenanceRecords(mRecords || [])
      setDriverRecords(dRecords || [])
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

  const activeVehicles = vehicles.filter(v => v.status !== 'INACTIVE')
  const archivedVehicles = vehicles.filter(v => v.status === 'INACTIVE')
  const displayedVehicles = viewTab === 'ACTIVE' ? activeVehicles : archivedVehicles

  if (loading) {
    return (
      <div className="page-loading">
        <div className="spinner" />
        <p>Завантаження автопарку...</p>
      </div>
    )
  }

  return (
    <div className="vehicles-page">
      <Toaster position="top-right" />
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
                <select value={newVehicle.driver_id} onChange={(e) => setNewVehicle({ ...newVehicle, driver_id: e.target.value })} className="erp-input">
                  <option value="">-- Без закріплення --</option>
                  {usersList.map(u => <option key={u.id} value={u.id}>{u.full_name} ({u.role})</option>)}
                </select>
              </div>

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowVehicleForm(false)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={isProcessing}>Створити</button>
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
            <h3>Зміна статусу: {statusModalVehicle.brand} ({statusModalVehicle.plate_number})</h3>
            <form onSubmit={handleUpdateStatus}>
              <p className="modal-description">Виберіть новий стан техніки та обов'язково вкажіть причину.</p>
              <div className="form-group">
                <label>Новий статус</label>
                <select value={statusForm.status} onChange={(e) => setStatusForm({...statusForm, status: e.target.value})} className="erp-input" disabled={isProcessing}>
                  {statusModalVehicle.status !== 'IN_REPAIR' && <option value="IN_REPAIR">🛠 Відправити в ремонт</option>}
                  <option value="INACTIVE">🔥 Списати (Тотал / Повна втрата)</option>
                </select>
              </div>
              <div className="form-group">
                <label>Причина <span style={{color: '#ef4444'}}>*</span></label>
                <textarea rows={3} placeholder={statusForm.status === 'IN_REPAIR' ? "Напр., Поломка двигуна / ДТП" : "Напр., Авто не підлягає відновленню після ДТП"} value={statusForm.reason} onChange={(e) => setStatusForm({...statusForm, reason: e.target.value})} required disabled={isProcessing} className="erp-input" />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setStatusModalVehicle(null)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className="btn btn-danger" disabled={isProcessing}>Підтвердити дію</button>
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
              <p className="modal-description">Виберіть співробітника, за яким буде закріплено даний транспортний засіб.</p>
              <div className="form-group">
                <label>Відповідальний водій</label>
                <select value={driverForm.driver_id} onChange={(e) => setDriverForm({...driverForm, driver_id: e.target.value})} className="erp-input" disabled={isProcessing}>
                  <option value="">-- Зняти закріплення (Без водія) --</option>
                  {usersList.map(u => <option key={u.id} value={u.id}>{u.full_name} ({u.role})</option>)}
                </select>
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
                  <input type="number" min={(fuelModalVehicle as any).current_odometer || 0} placeholder={`Напр. ${(fuelModalVehicle as any).current_odometer ? (fuelModalVehicle as any).current_odometer + 150 : 150500}`} value={fuelForm.odometer_km} onChange={(e) => setFuelForm({ ...fuelForm, odometer_km: e.target.value })} required={fuelForm.record_type === 'EXPENSE'} disabled={isProcessing} className="erp-input" />
                  <span style={{ display: 'block', fontSize: '11px', color: '#64748b', marginTop: '6px' }}>Останній запис: <strong>{(fuelModalVehicle as any).current_odometer || 0}</strong> км</span>
                </div>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setFuelModalVehicle(null)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className={fuelForm.record_type === 'REFUEL' ? "btn btn-success" : "btn btn-danger"} disabled={isProcessing}>{fuelForm.record_type === 'REFUEL' ? 'Заправити' : 'Списати'}</button>
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
                  <label>Одометр після ремонту (км) <span style={{color: '#ef4444'}}>*</span></label>
                  <input type="number" min={maintenanceModalVehicle.last_maintenance_odometer} value={maintenanceForm.odometer_km || ''} onChange={(e) => setMaintenanceForm({...maintenanceForm, odometer_km: parseInt(e.target.value, 10)})} required disabled={isProcessing} className="erp-input" />
                </div>
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
                    <input type="file" className="file-input-hidden" accept="image/*,application/pdf" onChange={(e) => { if (e.target.files && e.target.files.length > 0) setMaintenanceForm({...maintenanceForm, document: e.target.files[0]}) }} disabled={isProcessing} />
                    <span className="file-upload-text">{maintenanceForm.document ? `📎 ${maintenanceForm.document.name}` : '📁 Натисніть, щоб вибрати...'}</span>
                  </label>
                </div>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setMaintenanceModalVehicle(null)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className={maintenanceModalVehicle.status === 'IN_REPAIR' ? "btn btn-success" : "btn btn-primary"} disabled={isProcessing}>{maintenanceModalVehicle.status === 'IN_REPAIR' ? 'Повернути на лінію' : 'Зафіксувати ТО'}</button>
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
                                  href={`http://localhost:8080${record.document_url}`} 
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
                        </tr>
                      </thead>
                      <tbody>
                        {driverRecords.map(record => (
                          <tr key={record.id}>
                            <td className="text-muted">
                              {new Date(record.assigned_at).toLocaleString('uk-UA', { 
                                day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' 
                              })}
                            </td>
                            <td className={record.driver_id ? 'driver-assigned' : 'driver-unassigned'}>
                              {record.driver_id ? `👤 ${record.driver_name}` : `🚫 ${record.driver_name}`}
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
        <div className="table-header-flex">
          <h2>{viewTab === 'ACTIVE' ? 'Автомобілі на балансі' : 'Архів списаної техніки'}</h2>
          
          <div className="table-header-actions">
            <button className={`tab-toggle-btn ${viewTab === 'ACTIVE' ? 'active' : ''}`} onClick={() => setViewTab('ACTIVE')}>Активні авто ({activeVehicles.length})</button>
            <button className={`tab-toggle-btn archive ${viewTab === 'ARCHIVE' ? 'active' : ''}`} onClick={() => setViewTab('ARCHIVE')}>Списані ({archivedVehicles.length})</button>
          </div>
        </div>

        {displayedVehicles.length === 0 ? (
          <p className="history-empty" style={{marginTop: '20px'}}>{viewTab === 'ACTIVE' ? 'Активний автопарк порожній' : 'Немає списаної техніки'}</p>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>Марка / Модель</th>
                <th>Номерний знак</th>
                <th>Тип та Вантаж</th>
                <th>Закріплений водій</th>
                <th>Бак (Норма)</th>
                <th>Статус</th>
                {viewTab === 'ACTIVE' && <th>До ТО</th>}
                <th className="col-actions-menu">Дії</th>
              </tr>
            </thead>
            <tbody>
              {displayedVehicles.map((v) => {
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

                    <td>{v.tank_capacity} л <span className="norm-text">({v.fuel_norm} л/100км)</span></td>
                    
                    <td>
                      <span className={`badge ${v.status === 'ACTIVE' ? 'badge-success' : v.status === 'IN_REPAIR' ? 'badge-warning' : v.status === 'ON_MISSION' ? 'badge-primary' : 'badge-critical'}`}>
                        {v.status === 'ACTIVE' ? 'На лінії' : v.status === 'IN_REPAIR' ? 'В ремонті' : v.status === 'ON_MISSION' ? 'У рейсі' : 'Списане'}
                      </span>
                    </td>
                    
                    {viewTab === 'ACTIVE' && (
                      <td>
                        {v.maintenance_status === 'OVERDUE' && <span className="badge badge-critical" title={`Інтервал: ${v.maintenance_interval_km} км`}>Прострочено на {Math.abs(v.km_to_next_maintenance)} км</span>}
                        {v.maintenance_status === 'WARNING' && <span className="badge badge-warning" title={`Інтервал: ${v.maintenance_interval_km} км`}>Залишилось {v.km_to_next_maintenance} км</span>}
                        {v.maintenance_status === 'OK' && <span className="text-muted" title={`Інтервал: ${v.maintenance_interval_km} км`}>Ще {v.km_to_next_maintenance} км</span>}
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
                            
                            {canManageVehicles && viewTab === 'ACTIVE' && (
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
                                  <button className="text-success" onClick={() => { setMaintenanceModalVehicle(v); setMaintenanceForm({ odometer_km: v.current_odometer || 0, description: '', performed_by: '', cost_amount: 0, document: null }); setActiveMenuId(null); }}>
                                    ✅ Завершити ремонт
                                  </button>
                                ) : (
                                  <button onClick={() => { setMaintenanceModalVehicle(v); setMaintenanceForm({ odometer_km: v.current_odometer || 0, description: '', performed_by: '', cost_amount: 0, document: null }); setActiveMenuId(null); }}>
                                    🛠 Зафіксувати ТО
                                  </button>
                                )}
                                
                                <button onClick={() => { setStatusModalVehicle(v); setStatusForm({ status: v.status === 'IN_REPAIR' ? 'INACTIVE' : 'IN_REPAIR', reason: '' }); setActiveMenuId(null); }}>
                                  🚦 Змінити статус
                                </button>

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
        )}
      </div>
    </div>
  )
}