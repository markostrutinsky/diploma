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

  // Стейти: Редагування та видалення/списання
  const [editingVehicle, setEditingVehicle] = useState<Vehicle | null>(null);
  const [vehicleToDelete, setVehicleToDelete] = useState<Vehicle | null>(null);
  const [editForm, setEditForm] = useState({ brand: '', model: '', plate_number: '', capacity_kg: 0 });
  const [isProcessing, setIsProcessing] = useState(false);

  const canManageVehicles = ['ADMIN', 'BRIGADE_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST', 'COMPANY_SERGEANT'].includes(user?.role || '')

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
        ? 'Ремонт завершено! Машина знову в строю.' 
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
      toast.success('Екіпаж успішно оновлено!')
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
    <div className="inventory-page">
      <Toaster position="top-right" />
      <div className="page-header">
        <h1>Автопарк та ГСМ</h1>
        <div className="page-actions">
          {canManageVehicles && (
            <button className="btn btn-primary" onClick={() => setShowVehicleForm(true)}>
              + Автомобіль
            </button>
          )}
        </div>
      </div>

      {anomalyAlert && (
        <div className="modal-overlay">
          <div className="modal confirm-modal anomaly-modal">
            <h3 className="text-danger">🚨 Виявлено аномалію!</h3>
            <p className="confirm-text">
              Запис збережено в систему, але він позначений як підозрілий і потребує перевірки.
            </p>
            <div className="warning-box">
              <p><strong>Причина від бекенду:</strong></p>
              <p className="text-danger">{anomalyAlert.anomaly_reason}</p>
            </div>
            <div className="modal-actions center-actions">
              <button className="btn btn-primary" onClick={() => setAnomalyAlert(null)}>
                Зрозуміло
              </button>
            </div>
          </div>
        </div>
      )}

      {/* МОДАЛКА НОВОГО АВТО */}
      {showVehicleForm && canManageVehicles && (
        <div className="modal-overlay" onClick={() => setShowVehicleForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Новий транспорт</h3>
            <form onSubmit={handleCreateVehicle}>
              <div className="form-group">
                <label>Марка (напр., Nissan, КрАЗ)</label>
                <input
                  value={newVehicle.brand}
                  onChange={(e) => setNewVehicle({ ...newVehicle, brand: e.target.value })}
                  required
                />
              </div>
              <div className="form-row-2">
                <div className="form-group">
                  <label>Модель</label>
                  <input
                    value={newVehicle.model}
                    onChange={(e) => setNewVehicle({ ...newVehicle, model: e.target.value })}
                  />
                </div>
                <div className="form-group">
                  <label>Номерний знак <span className="required">*</span></label>
                  <input
                    value={newVehicle.plate_number}
                    onChange={(e) => setNewVehicle({ ...newVehicle, plate_number: e.target.value })}
                    required
                  />
                </div>
              </div>

              <div className="form-row-2">
                <div className="form-group">
                  <label>Тип кузова <span className="required">*</span></label>
                  <select 
                    value={newVehicle.type} 
                    onChange={(e) => setNewVehicle({...newVehicle, type: e.target.value})}
                    className="fuel-type-select select-normal-weight"
                  >
                    <option value="PICKUP">🛻 Пікап / Джип</option>
                    <option value="VAN">🚐 Мікроавтобус / Фургон</option>
                    <option value="TRUCK">🚛 Вантажівка</option>
                  </select>
                </div>
                <div className="form-group">
                  <label>Вантажопідйомність (кг) <span className="required">*</span></label>
                  <input
                    type="number"
                    min="1"
                    value={newVehicle.capacity_kg || ''}
                    onChange={(e) => setNewVehicle({ ...newVehicle, capacity_kg: parseFloat(e.target.value) })}
                    required
                  />
                </div>
              </div>

              <div className="form-row-2">
                <div className="form-group">
                  <label>Об'єм бака (л) <span className="required">*</span></label>
                  <input
                    type="number"
                    min="1"
                    value={newVehicle.tank_capacity || ''}
                    onChange={(e) => setNewVehicle({ ...newVehicle, tank_capacity: parseFloat(e.target.value) })}
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Норма (л/100км) <span className="required">*</span></label>
                  <input
                    type="number"
                    min="1"
                    step="0.1"
                    value={newVehicle.fuel_norm || ''}
                    onChange={(e) => setNewVehicle({ ...newVehicle, fuel_norm: parseFloat(e.target.value) })}
                    required
                  />
                </div>
              </div>

              <div className="form-group">
                <label>Закріплений водій (Екіпаж)</label>
                <select 
                  value={newVehicle.driver_id} 
                  onChange={(e) => setNewVehicle({ ...newVehicle, driver_id: e.target.value })}
                  className="fuel-type-select select-normal-weight"
                >
                  <option value="">-- Без закріплення --</option>
                  {usersList.map(u => (
                    <option key={u.id} value={u.id}>
                      {u.full_name} ({u.role})
                    </option>
                  ))}
                </select>
              </div>

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => setShowVehicleForm(false)} disabled={isProcessing}>
                  Скасувати
                </button>
                <button type="submit" className="btn btn-primary" disabled={isProcessing}>
                  Створити
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* МОДАЛКА РЕДАГУВАННЯ АВТО */}
      {editingVehicle && canManageVehicles && (
        <div className="modal-overlay" onClick={() => !isProcessing && setEditingVehicle(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Редагувати параметри авто</h3>
            <form onSubmit={handleEditSubmit}>
              <div className="form-group">
                <label>Державний номер</label>
                <input className="erp-input" value={editForm.plate_number} onChange={(e) => setEditForm({ ...editForm, plate_number: e.target.value })} required disabled={isProcessing} />
              </div>
              <div className="form-row-2">
                <div className="form-group flex-1">
                  <label>Марка</label>
                  <input className="erp-input" value={editForm.brand} onChange={(e) => setEditForm({ ...editForm, brand: e.target.value })} required disabled={isProcessing} />
                </div>
                <div className="form-group flex-1">
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

      {/* МОДАЛКА ПІДТВЕРДЖЕННЯ СПИСАННЯ АВТО */}
      {vehicleToDelete && canManageVehicles && (
        <div className="modal-overlay" onClick={() => !isProcessing && setVehicleToDelete(null)}>
          <div className="modal confirm-modal" onClick={(e) => e.stopPropagation()}>
            <h3 style={{ color: '#ef4444' }}>⚠️ Списання автомобіля</h3>
            <p>Ви впевнені, що хочете остаточно списати авто <strong>{vehicleToDelete.brand} ({vehicleToDelete.plate_number})</strong>?</p>
            <p style={{ fontSize: '12px', color: '#64748b' }}>Це видалить його зі списку доступних машин для формування логістичних рейсів.</p>
            <div className="modal-actions">
              <button className="btn btn-secondary" onClick={() => setVehicleToDelete(null)} disabled={isProcessing}>Скасувати</button>
              <button className="btn btn-danger" onClick={handleDelete} disabled={isProcessing}>{isProcessing ? 'Списання...' : 'Списати'}</button>
            </div>
          </div>
        </div>
      )}

      {/* МОДАЛКА ЗМІНИ СТАТУСУ */}
      {statusModalVehicle && (
        <div className="modal-overlay" onClick={() => !isProcessing && setStatusModalVehicle(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Зміна статусу: {statusModalVehicle.brand} ({statusModalVehicle.plate_number})</h3>
            <form onSubmit={handleUpdateStatus}>
              <p className="modal-description">
                Виберіть новий стан техніки та обов'язково вкажіть причину.
              </p>
              
              <div className="form-group">
                <label>Новий статус</label>
                <select 
                  value={statusForm.status} 
                  onChange={(e) => setStatusForm({...statusForm, status: e.target.value})}
                  className="fuel-type-select"
                  disabled={isProcessing}
                >
                  {statusModalVehicle.status !== 'IN_REPAIR' && (
                    <option value="IN_REPAIR">🛠 Відправити в ремонт</option>
                  )}
                  <option value="INACTIVE">🔥 Списати (Безповоротна втрата)</option>
                </select>
              </div>

              <div className="form-group">
                <label>Причина <span className="required">*</span></label>
                <textarea
                  rows={3}
                  placeholder={statusForm.status === 'IN_REPAIR' ? "Напр., Поломка двигуна / ДТП" : "Напр., Знищено внаслідок бойових дій"}
                  value={statusForm.reason}
                  onChange={(e) => setStatusForm({...statusForm, reason: e.target.value})}
                  required
                  disabled={isProcessing}
                />
              </div>

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => setStatusModalVehicle(null)} disabled={isProcessing}>
                  Скасувати
                </button>
                <button type="submit" className="btn btn-danger" disabled={isProcessing}>
                  Підтвердити дію
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* МОДАЛКА ЗМІНИ ВОДІЯ */}
      {driverModalVehicle && (
        <div className="modal-overlay" onClick={() => !isProcessing && setDriverModalVehicle(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Екіпаж: {driverModalVehicle.brand} ({driverModalVehicle.plate_number})</h3>
            <form onSubmit={handleAssignDriver}>
              <p className="modal-description">
                Виберіть військовослужбовця, за яким буде закріплено даний транспортний засіб.
              </p>
              
              <div className="form-group">
                <label>Відповідальний водій</label>
                <select 
                  value={driverForm.driver_id} 
                  onChange={(e) => setDriverForm({...driverForm, driver_id: e.target.value})}
                  className="fuel-type-select select-normal-weight"
                  disabled={isProcessing}
                >
                  <option value="">-- Зняти закріплення (Без водія) --</option>
                  {usersList.map(u => (
                    <option key={u.id} value={u.id}>
                      {u.full_name} ({u.role})
                    </option>
                  ))}
                </select>
              </div>

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => setDriverModalVehicle(null)} disabled={isProcessing}>
                  Скасувати
                </button>
                <button type="submit" className="btn btn-primary" disabled={isProcessing}>
                  Зберегти зміни
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* МОДАЛКА ПАЛЬНОГО */}
      {fuelModalVehicle && (
        <div className="modal-overlay" onClick={() => !isProcessing && setFuelModalVehicle(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Пальне: {fuelModalVehicle.brand} ({fuelModalVehicle.plate_number})</h3>
            <form onSubmit={handleAddFuel}>
              <div className="form-group">
                <label>Тип операції</label>
                <select
                  value={fuelForm.record_type}
                  onChange={(e) => setFuelForm({ ...fuelForm, record_type: e.target.value as FuelRecordType })}
                  className={`fuel-type-select ${fuelForm.record_type === 'REFUEL' ? 'type-refuel' : 'type-expense'}`}
                  disabled={isProcessing}
                >
                  <option value="EXPENSE">Списання (Витрата)</option>
                  <option value="REFUEL">Заправка (Прихід)</option>
                </select>
              </div>
              
              <div className="form-row-2">
                <div className="form-group">
                  <label>Літри <span className="required">*</span></label>
                  <input
                    type="number"
                    min="0.1"
                    step="0.1"
                    value={fuelForm.liters || ''}
                    onChange={(e) => setFuelForm({ ...fuelForm, liters: parseFloat(e.target.value) })}
                    required
                    disabled={isProcessing}
                  />
                </div>
                <div className="form-group">
                  <label>
                    Поточний одометр (км)
                    {fuelForm.record_type === 'EXPENSE' && <span className="required">*</span>}
                  </label>
                  <input
                    type="number"
                    min={(fuelModalVehicle as any).current_odometer || 0}
                    placeholder={`Напр. ${(fuelModalVehicle as any).current_odometer ? (fuelModalVehicle as any).current_odometer + 150 : 150500}`}
                    value={fuelForm.odometer_km}
                    onChange={(e) => setFuelForm({ ...fuelForm, odometer_km: e.target.value })}
                    required={fuelForm.record_type === 'EXPENSE'}
                    disabled={isProcessing}
                  />
                  <span style={{ display: 'block', fontSize: '11px', color: '#64748b', marginTop: '6px' }}>
                    Останній запис: <strong>{(fuelModalVehicle as any).current_odometer || 0}</strong> км
                  </span>
                </div>
              </div>

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => setFuelModalVehicle(null)} disabled={isProcessing}>
                  Скасувати
                </button>
                <button 
                  type="submit" 
                  className={fuelForm.record_type === 'REFUEL' ? "btn btn-success" : "btn btn-danger"}
                  disabled={isProcessing}
                >
                  {fuelForm.record_type === 'REFUEL' ? 'Заправити' : 'Списати'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* МОДАЛКА ТО */}
      {maintenanceModalVehicle && (
        <div className="modal-overlay" onClick={() => !isProcessing && setMaintenanceModalVehicle(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>
              {maintenanceModalVehicle.status === 'IN_REPAIR' 
                ? '✅ Завершення ремонту: ' 
                : '🛠 Акт виконаних робіт: '} 
              {maintenanceModalVehicle.brand} ({maintenanceModalVehicle.plate_number})
            </h3>
            <form onSubmit={handlePerformMaintenance}>
              
              <div className="form-row-2">
                <div className="form-group">
                  <label>Одометр після ремонту (км) <span className="required">*</span></label>
                  <input
                    type="number"
                    min={maintenanceModalVehicle.last_maintenance_odometer}
                    value={maintenanceForm.odometer_km || ''}
                    onChange={(e) => setMaintenanceForm({...maintenanceForm, odometer_km: parseInt(e.target.value, 10)})}
                    required
                    disabled={isProcessing}
                  />
                </div>
                <div className="form-group">
                  <label>Виконавець (Свої сили / СТО)</label>
                  <input
                    type="text"
                    placeholder="Напр. Ремрота або СТО 'Гараж'"
                    value={maintenanceForm.performed_by}
                    onChange={(e) => setMaintenanceForm({...maintenanceForm, performed_by: e.target.value})}
                    disabled={isProcessing}
                  />
                </div>
              </div>

              <div className="form-group">
                <label>Опис робіт та замінених запчастин <span className="required">*</span></label>
                <textarea
                  rows={3}
                  placeholder="Заміна мастила 5w40 (5л), масляний фільтр..."
                  value={maintenanceForm.description}
                  onChange={(e) => setMaintenanceForm({...maintenanceForm, description: e.target.value})}
                  required
                  disabled={isProcessing}
                />
              </div>

              <div className="form-row-2 form-row-bottom-align">
                <div className="form-group form-group-no-margin">
                  <label>Загальна вартість (Грн)</label>
                  <input
                    type="number"
                    min="0"
                    step="0.01"
                    placeholder="Напр. 15000"
                    value={maintenanceForm.cost_amount || ''}
                    onChange={(e) => setMaintenanceForm({...maintenanceForm, cost_amount: parseFloat(e.target.value)})}
                    disabled={isProcessing}
                  />
                </div>
                
                <div className="form-group form-group-no-margin">
                  <label>Скан Акту (PDF/Фото) <span className="required">*</span></label>
                  <label className="file-upload-custom">
                    <input 
                      type="file" 
                      className="file-input-hidden"
                      accept="image/*,application/pdf"
                      onChange={(e) => {
                        if (e.target.files && e.target.files.length > 0) {
                          setMaintenanceForm({...maintenanceForm, document: e.target.files[0]})
                        }
                      }}
                      disabled={isProcessing}
                    />
                    <span className="file-upload-text">
                      {maintenanceForm.document 
                        ? `📎 ${maintenanceForm.document.name}` 
                        : '📁 Натисніть, щоб вибрати...'}
                    </span>
                  </label>
                </div>
              </div>

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => setMaintenanceModalVehicle(null)} disabled={isProcessing}>
                  Скасувати
                </button>
                <button type="submit" className="btn btn-primary" disabled={isProcessing}>
                  {maintenanceModalVehicle.status === 'IN_REPAIR' ? 'Повернути в стрій' : 'Зафіксувати ТО'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* МОДАЛКА ІСТОРІЇ */}
      {historyVehicle && (
        <div className="modal-overlay" onClick={() => setHistoryVehicle(null)}>
          <div className="modal history-modal" onClick={(e) => e.stopPropagation()}>
            <h3>Паспорт машини: {historyVehicle.brand} ({historyVehicle.plate_number})</h3>
            
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
                  {historyVehicle.status === 'ACTIVE' ? 'В строю' : historyVehicle.status === 'IN_REPAIR' ? 'В ремонті' : historyVehicle.status === 'ON_MISSION' ? 'У рейсі' : 'Списане'}
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
              <button className={`history-tab ${historyTab === 'DRIVERS' ? 'active' : ''}`} onClick={() => setHistoryTab('DRIVERS')}>👥 Історія екіпажів</button>
            </div>

            {historyLoading ? (
              <div className="spinner history-spinner" />
            ) : (
              <div className="history-table-wrapper">
                
                {historyTab === 'FUEL' && (
                  fuelRecords.length === 0 ? (
                    <p className="history-empty">Записів про пальне ще немає.</p>
                  ) : (
                    <table className="data-table">
                      <thead>
                        <tr><th>Дата</th><th>Тип</th><th>Літри</th><th>Одометр</th><th>Статус</th></tr>
                      </thead>
                      <tbody>
                        {fuelRecords.map(record => (
                          <tr key={record.id} className={record.is_anomaly ? 'row-critical' : ''}>
                            <td className="date-cell">
                              {new Date(record.created_at).toLocaleString('uk-UA', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })}
                            </td>
                            <td className={`fuel-type-select ${record.record_type === 'REFUEL' ? 'type-refuel' : 'type-expense'}`}>
                              {record.record_type === 'REFUEL' ? 'Прихід' : 'Списання'}
                            </td>
                            <td className="liters-cell">{record.liters} л</td>
                            <td>{record.odometer_km ? `${record.odometer_km} км` : '-'}</td>
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
                    <table className="data-table">
                      <thead>
                        <tr>
                          <th>Дата ТО</th>
                          <th>Екіпаж</th>
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
                            <td className="date-cell">
                              {new Date(record.created_at).toLocaleString('uk-UA', { day: '2-digit', month: '2-digit', year: 'numeric' })}
                            </td>
                            <td className="driver-history-cell">
                              {record.driver_name ? `👤 ${record.driver_name}` : <span className="unassigned-text">Не призначено</span>}
                            </td>
                            <td className="odometer-cell">{record.odometer_km} км</td>
                            <td className="desc-cell">{record.description}</td>
                            <td className="performer-cell">{record.performed_by || '-'}</td>
                            <td className="cost-cell">
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
                    <p className="history-empty">Історія екіпажів порожня.</p>
                  ) : (
                    <table className="data-table">
                      <thead>
                        <tr>
                          <th>Дата призначення</th>
                          <th>Військовослужбовець (Екіпаж)</th>
                        </tr>
                      </thead>
                      <tbody>
                        {driverRecords.map(record => (
                          <tr key={record.id}>
                            <td className="date-cell">
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

            <div className="modal-actions history-actions">
              <button className="btn btn-secondary" onClick={() => setHistoryVehicle(null)}>
                Закрити
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ТАБЛИЦЯ АВТОПАРКУ */}
      <div className="card card-table">
        <div className="table-header-flex">
          <h2>{viewTab === 'ACTIVE' ? 'Автомобілі на балансі' : 'Архів списаної техніки'}</h2>
          
          <div className="table-header-actions">
            <button 
              className={`tab-toggle-btn ${viewTab === 'ACTIVE' ? 'active' : ''}`}
              onClick={() => setViewTab('ACTIVE')}
            >
              Активні авто ({activeVehicles.length})
            </button>
            <button 
              className={`tab-toggle-btn archive ${viewTab === 'ARCHIVE' ? 'active' : ''}`}
              onClick={() => setViewTab('ARCHIVE')}
            >
              Списані ({archivedVehicles.length})
            </button>
          </div>
        </div>

        {displayedVehicles.length === 0 ? (
          <p className="empty-state">{viewTab === 'ACTIVE' ? 'Активний автопарк порожній' : 'Немає списаної техніки'}</p>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>Марка / Модель</th>
                <th>Номерний знак</th>
                <th>Тип та Вантаж</th>
                <th>Екіпаж / Водій</th>
                <th>Бак (Норма)</th>
                <th>Статус</th>
                {viewTab === 'ACTIVE' && <th>До ТО</th>}
                <th>Дії</th>
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
                      <div>
                        {v.type === 'PICKUP' ? '🛻 Пікап' : v.type === 'VAN' ? '🚐 Фургон' : v.type === 'TRUCK' ? '🚛 Вантажівка' : '🚗 Авто'}
                      </div>
                      <div className="norm-text" style={{ marginLeft: 0 }}>Макс: {v.capacity_kg} кг</div>
                    </td>

                    <td>{getDriverName(v.driver_id)}</td>

                    <td>{v.tank_capacity} л <span className="norm-text">({v.fuel_norm} л/100км)</span></td>
                    
                    <td>
                      <span className={`badge ${v.status === 'ACTIVE' ? 'badge-success' : v.status === 'IN_REPAIR' ? 'badge-warning' : v.status === 'ON_MISSION' ? 'badge-primary' : 'badge-critical'}`}>
                        {v.status === 'ACTIVE' ? 'В строю' : v.status === 'IN_REPAIR' ? 'В ремонті' : v.status === 'ON_MISSION' ? 'У рейсі' : 'Списане'}
                      </span>
                    </td>
                    
                    {viewTab === 'ACTIVE' && (
                      <td>
                        {v.maintenance_status === 'OVERDUE' && <span className="badge badge-critical" title={`Інтервал: ${v.maintenance_interval_km} км`}>Прострочено на {Math.abs(v.km_to_next_maintenance)} км</span>}
                        {v.maintenance_status === 'WARNING' && <span className="badge badge-warning" title={`Інтервал: ${v.maintenance_interval_km} км`}>Залишилось {v.km_to_next_maintenance} км</span>}
                        {v.maintenance_status === 'OK' && <span className="text-muted" title={`Інтервал: ${v.maintenance_interval_km} км`}>Ще {v.km_to_next_maintenance} км</span>}
                      </td>
                    )}

                    <td>
                      <div className="action-buttons-group">
                        <button className="btn btn-secondary btn-sm btn-fuel-action" onClick={() => handleViewHistory(v)}>
                          📊 Паспорт
                        </button>
                        
                        {canManageVehicles && viewTab === 'ACTIVE' && (
                          <>
                            <button 
                              className="btn btn-secondary btn-sm btn-fuel-action" 
                              onClick={() => {
                                setDriverModalVehicle(v)
                                setDriverForm({ driver_id: v.driver_id || '' })
                              }}
                            >
                              👤 Водій
                            </button>
                            {v.status === 'ACTIVE' && (
                              <button 
                                className="btn btn-secondary btn-sm btn-fuel-action" 
                                onClick={() => {
                                  setFuelModalVehicle(v)
                                  setFuelForm({ record_type: 'EXPENSE', liters: 0, odometer_km: '' })
                                }}
                              >
                                ⛽ Пальне
                              </button>
                            )}

                            {v.status === 'IN_REPAIR' ? (
                              <button 
                                className="btn btn-finish-repair btn-sm btn-fuel-action" 
                                onClick={() => {
                                  setMaintenanceModalVehicle(v)
                                  setMaintenanceForm({ odometer_km: v.current_odometer || 0, description: '', performed_by: '', cost_amount: 0, document: null })
                                }}
                              >
                                ✅ Завершити ремонт
                              </button>
                            ) : (
                              <button 
                                className="btn btn-secondary btn-sm btn-fuel-action" 
                                onClick={() => {
                                  setMaintenanceModalVehicle(v)
                                  setMaintenanceForm({ odometer_km: v.current_odometer || 0, description: '', performed_by: '', cost_amount: 0, document: null })
                                }}
                              >
                                🛠 Зафіксувати ТО
                              </button>
                            )}
                            
                            <button 
                              className="btn btn-secondary btn-sm btn-fuel-action" 
                              onClick={() => {
                                setStatusModalVehicle(v)
                                setStatusForm({ status: v.status === 'IN_REPAIR' ? 'INACTIVE' : 'IN_REPAIR', reason: '' })
                              }}
                            >
                              🚦 Статус
                            </button>

                            {/* НОВІ КНОПКИ РЕДАГУВАННЯ І ВИДАЛЕННЯ */}
                            <button 
                              className="btn-unit-action btn-unit-edit" 
                              onClick={() => handleOpenEdit(v)}
                              title="Редагувати параметри авто"
                            >
                              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                            </button>
                            <button 
                              className="btn-unit-action btn-unit-delete" 
                              onClick={() => setVehicleToDelete(v)}
                              title="Остаточно списати автомобіль"
                            >
                              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                            </button>

                          </>
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