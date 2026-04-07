import React, { useEffect, useState } from 'react';
import { api, type Warehouse, type Unit } from '../api/client';
import { useAuth } from '../contexts/AuthContext';
import './Inventory.css'; // Перевикористовуємо стилі таблиць та модалок

export default function Warehouses() {
  const { user } = useAuth();
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [units, setUnits] = useState<Unit[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);

  const [newWarehouse, setNewWarehouse] = useState({
    unit_id: '' as number | '',
    name: '',
    location_type: 'STATIONARY' as 'STATIONARY' | 'MOBILE',
    latitude: '',
    longitude: ''
  });

  // Хто може створювати склади
  const managerRoles = [
    'ADMIN', 'BRIGADE_CMDR', 'BRIGADE_LOGIST', 
    'BATTALION_CMDR', 'BATTALION_LOGIST', 
    'COMPANY_CMDR', 'PLATOON_CMDR'
  ];
  const canManageWarehouses = managerRoles.includes(user?.role || '');

  const loadData = async () => {
    try {
      const [wRes, uRes] = await Promise.all([
        api.warehouses.list().catch(() => []),
        api.units.list().catch(() => [])
      ]);
      setWarehouses(Array.isArray(wRes) ? wRes : []);
      setUnits(Array.isArray(uRes) ? uRes : []);
    } catch (error) {
      console.error('Помилка завантаження даних:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newWarehouse.unit_id) return;

    try {
      await api.warehouses.create({
        unit_id: Number(newWarehouse.unit_id),
        name: newWarehouse.name,
        location_type: newWarehouse.location_type,
        latitude: newWarehouse.latitude ? parseFloat(newWarehouse.latitude) : undefined,
        longitude: newWarehouse.longitude ? parseFloat(newWarehouse.longitude) : undefined,
      });
      setShowForm(false);
      setNewWarehouse({ unit_id: '', name: '', location_type: 'STATIONARY', latitude: '', longitude: '' });
      loadData(); // Оновлюємо список
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Помилка при створенні складу');
    }
  };

  if (loading) {
    return (
      <div className="page-loading">
        <div className="spinner" />
        <p>Завантаження інфраструктури...</p>
      </div>
    );
  }

  return (
    <div className="inventory-page">
      <div className="page-header">
        <h1>Склади та Інфраструктура</h1>
        <div className="page-actions">
          {canManageWarehouses && (
            <button className="btn btn-primary" onClick={() => setShowForm(true)}>
              + Створити склад
            </button>
          )}
        </div>
      </div>

      {showForm && canManageWarehouses && (
        <div className="modal-overlay" onClick={() => setShowForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Новий склад / пункт забезпечення</h3>
            <form onSubmit={handleCreate}>
              
              <div className="form-group">
                <label>Підрозділ-власник <span className="required">*</span></label>
                <select
  value={newWarehouse.unit_id}
  onChange={(e) => setNewWarehouse({ ...newWarehouse, unit_id: e.target.value ? Number(e.target.value) : '' })}
  required
>
  <option value="" disabled>Оберіть підрозділ</option>
  {units.map((u) => (
    <option key={u.id} value={u.id}>{u.name}</option> // Переконайся, що тут u.id
  ))}
</select>
              </div>

              <div className="form-group">
                <label>Назва складу <span className="required">*</span></label>
                <input
                  placeholder="Напр. БПЗ 1-го Батальйону"
                  value={newWarehouse.name}
                  onChange={(e) => setNewWarehouse({ ...newWarehouse, name: e.target.value })}
                  required
                />
              </div>

              <div className="form-group">
                <label>Тип локації <span className="required">*</span></label>
                <select
                  value={newWarehouse.location_type}
                  onChange={(e) => setNewWarehouse({ ...newWarehouse, location_type: e.target.value as 'STATIONARY' | 'MOBILE' })}
                  required
                >
                  <option value="STATIONARY">Стаціонарний (Будівля, ангар, бліндаж)</option>
                  <option value="MOBILE">Мобільний (Вантажівка, фура)</option>
                </select>
              </div>

              <div className="form-row-2">
                <div className="form-group">
                  <label>Широта (Latitude)</label>
                  <input
                    type="number"
                    step="0.000001"
                    placeholder="Напр. 48.4500"
                    value={newWarehouse.latitude}
                    onChange={(e) => setNewWarehouse({ ...newWarehouse, latitude: e.target.value })}
                  />
                </div>
                <div className="form-group">
                  <label>Довгота (Longitude)</label>
                  <input
                    type="number"
                    step="0.000001"
                    placeholder="Напр. 34.9833"
                    value={newWarehouse.longitude}
                    onChange={(e) => setNewWarehouse({ ...newWarehouse, longitude: e.target.value })}
                  />
                </div>
              </div>

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => setShowForm(false)}>
                  Скасувати
                </button>
                <button type="submit" className="btn btn-primary" disabled={!newWarehouse.unit_id}>
                  Створити
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className="card card-table">
        {warehouses.length === 0 ? (
          <p className="empty-state">Склади ще не створені. Додайте перший склад, щоб почати облік майна.</p>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>Назва складу</th>
                <th>Підрозділ</th>
                <th>Тип</th>
                <th>Координати</th>
              </tr>
            </thead>
            <tbody>
              {warehouses.map((w) => {
                const unitName = units.find(u => u.id === w.unit_id)?.name || 'Невідомий підрозділ';
                return (
                  <tr key={w.id}>
                    <td style={{ fontWeight: 'bold' }}>{w.name}</td>
                    <td>{unitName}</td>
                    <td>
                      <span className={`badge ${w.location_type === 'MOBILE' ? 'badge-warning' : 'badge-success'}`}>
                        {w.location_type === 'MOBILE' ? 'Мобільний 🚛' : 'Стаціонарний 🏢'}
                      </span>
                    </td>
                    <td className="text-muted">
                      {w.latitude && w.longitude 
                        ? `${w.latitude}, ${w.longitude}` 
                        : 'Не вказано'}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}