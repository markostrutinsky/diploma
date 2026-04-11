import React, { useEffect, useState } from 'react';
import { api, type Resource, type ResourceCategory, type Unit, type Warehouse } from '../api/client';
import { useAuth } from '../contexts/AuthContext';
import toast, { Toaster } from 'react-hot-toast';
import './Inventory.css';

export default function Inventory() {
  const { user } = useAuth();
  
  const [categories, setCategories] = useState<ResourceCategory[]>([]);
  const [resources, setResources] = useState<Resource[]>([]);
  const [units, setUnits] = useState<Unit[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [filterUnitId, setFilterUnitId] = useState<number | ''>('');
  
  const [loading, setLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);
  
  const [showCategoryForm, setShowCategoryForm] = useState(false);
  const [showResourceForm, setShowResourceForm] = useState(false);
  const [transferError, setTransferError] = useState<string | null>(null);
  const [writeOffModalData, setWriteOffModalData] = useState<{ resource: Resource; quantity: number; } | null>(null);
  const [usersList, setUsersList] = useState<any[]>([]);
  const [assignModalData, setAssignModalData] = useState<{ resource: Resource; quantity: number; user_id: string; } | null>(null);
  const [assignError, setAssignError] = useState<string | null>(null);
  const [writeOffError, setWriteOffError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'active' | 'written_off'>('active');
  const [editModalId, setEditModalId] = useState<string | null>(null);
  const [editForm, setEditForm] = useState({ name: '', min_quantity: 0 });
  const [resourceToDelete, setResourceToDelete] = useState<Resource | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const [transferModalData, setTransferModalData] = useState<{ resource: Resource; quantity: number; target_unit_id: number | ''; target_warehouse_id: string; } | null>(null);

  const [newCat, setNewCat] = useState({ name: '', description: '' });
  
  const [newRes, setNewRes] = useState({ 
    category_id: '', 
    unit_id: undefined as number | undefined, 
    warehouse_id: '',
    name: '', 
    quantity: 0, 
    unit_type: 'PCS' as 'PCS' | 'KIT' | 'KG' | 'L',
    min_quantity: 0,
    weight_kg: 1
  });

  const [selectedCategoryId, setSelectedCategoryId] = useState<string | null>(null);
  const [activeMenuId, setActiveMenuId] = useState<string | null>(null);

  useEffect(() => {
    const closeMenu = () => setActiveMenuId(null);
    document.addEventListener('click', closeMenu);
    return () => document.removeEventListener('click', closeMenu);
  }, []);

  const canManageResources = ['ADMIN', 'BRIGADE_STOREKEEPER', 'BATTALION_STOREKEEPER', 'COMPANY_SERGEANT', 'BRIGADE_LOGIST'].includes(user?.role || '');
  const canManageCategories = ['ADMIN', 'BRIGADE_LOGIST', 'BRIGADE_CMDR'].includes(user?.role || '');

  const loadData = () => {
    const unitId = filterUnitId || undefined;
    setIsRefreshing(true);

    Promise.all([
      api.inventory.listCategories().catch(() => []),
      api.inventory.listResources(unitId).catch(() => []),
      api.units.list().catch(() => []),
      api.warehouses.list().catch(() => []),
      api.users.getVisible().catch(() => [])
    ])
      .then(([cats, res, u, w, users]) => {
        const safeCats = Array.isArray(cats) ? cats : [];
        setCategories(safeCats);
        setResources(Array.isArray(res) ? res : []);
        setUnits(Array.isArray(u) ? u : []);
        setWarehouses(Array.isArray(w) ? w : []);
        setUsersList(Array.isArray(users) ? users : []);
        setNewRes((r) => (safeCats.length && !r.category_id ? { ...r, category_id: safeCats[0].id } : r));
      })
      .catch(console.error)
      .finally(() => { 
        setLoading(false); 
        setIsRefreshing(false); 
      });
  };

  useEffect(() => { 
    loadData(); 
  }, [filterUnitId]);

  const handleCreateCategory = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.inventory.createCategory(newCat);
      setShowCategoryForm(false);
      setNewCat({ name: '', description: '' });
      toast.success('Категорію успішно створено');
      loadData();
    } catch (err) { 
      toast.error(err instanceof Error ? err.message : 'Помилка створення категорії'); 
    }
  };

  const handleAssignSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!assignModalData) return;
    setAssignError(null);
    try {
      await api.inventory.assignResource(assignModalData.resource.id, { 
        quantity: assignModalData.quantity, 
        user_id: assignModalData.user_id 
      });
      setAssignModalData(null);
      toast.success('Майно успішно видано');
      loadData();
    } catch (err) { 
      setAssignError(err instanceof Error ? err.message : 'Помилка при видачі'); 
    }
  };

  const handleCreateResource = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.inventory.createResource({ 
        ...newRes, 
        unit_id: newRes.unit_id || undefined, 
        warehouse_id: newRes.warehouse_id || undefined 
      });
      setShowResourceForm(false);
      setNewRes({ 
        category_id: categories[0]?.id || '', 
        unit_id: undefined, 
        warehouse_id: '', 
        name: '', 
        quantity: 0, 
        unit_type: 'PCS', 
        min_quantity: 0, 
        weight_kg: 1 
      });
      toast.success('Ресурс успішно додано на склад');
      loadData();
    } catch (err) { 
      toast.error(err instanceof Error ? err.message : 'Помилка збереження ресурсу'); 
    }
  };

  const handleWriteOffSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!writeOffModalData) return;
    setWriteOffError(null);
    try {
      await api.inventory.writeOffResource(writeOffModalData.resource.id, writeOffModalData.quantity);
      setWriteOffModalData(null);
      toast.success('Майно успішно списано');
      loadData(); 
    } catch (err) { 
      setWriteOffError(err instanceof Error ? err.message : 'Помилка при списанні'); 
    }
  };

  const handleEditSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editModalId) return;
    try {
      await api.inventory.updateResource(editModalId, editForm);
      setEditModalId(null);
      toast.success('Дані ресурсу оновлено');
      loadData();
    } catch (err) { 
      toast.error(err instanceof Error ? err.message : 'Помилка оновлення ресурсу'); 
    }
  };

  const handleTransferSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!transferModalData) return;
    setTransferError(null);
    try {
      await api.inventory.transferResource(transferModalData.resource.id, { 
        quantity: transferModalData.quantity, 
        target_unit_id: transferModalData.target_unit_id === '' ? undefined : Number(transferModalData.target_unit_id), 
        target_warehouse_id: transferModalData.target_warehouse_id === '' ? undefined : transferModalData.target_warehouse_id 
      });
      setTransferModalData(null);
      toast.success('Майно переміщено в інший підрозділ');
      loadData();
    } catch (err) { 
      setTransferError(err instanceof Error ? err.message : 'Помилка при переміщенні'); 
    }
  };

  const confirmDelete = async () => {
    if (!resourceToDelete) return;
    setDeleteError(null);
    try {
      await api.inventory.deleteResource(resourceToDelete.id);
      setResourceToDelete(null);
      toast.success('Запис безповоротно видалено');
      loadData(); 
    } catch (err) { 
      setDeleteError(err instanceof Error ? err.message : 'Помилка при видаленні'); 
    }
  };

  const formatUnitType = (type: string) => {
    switch(type) { 
      case 'PCS': return 'шт'; 
      case 'KIT': return 'компл'; 
      case 'KG': return 'кг'; 
      case 'L': return 'л'; 
      default: return 'шт'; 
    }
  };

  if (loading) {
    return (
      <div className="page-loading">
        <div className="spinner" />
        <p>Завантаження інвентарю...</p>
      </div>
    );
  }

  const filteredResources = resources.filter(r => {
    const isDeleted = r.condition === 'WRITTEN_OFF';
    const matchesTab = activeTab === 'active' ? !isDeleted : isDeleted;
    const matchesCategory = selectedCategoryId ? r.category_id === selectedCategoryId : true;
    const isOrphanEmpty = (r.unit_id === 0 || r.unit_id == null) && r.quantity === 0 && !isDeleted;
    return matchesTab && matchesCategory && !isOrphanEmpty;
  });

  const groupedResources = filteredResources.reduce((acc, resource) => {
    const uId = resource.unit_id || 0; 
    if (!acc[uId]) acc[uId] = [];
    acc[uId].push(resource);
    return acc;
  }, {} as Record<number, Resource[]>);

  const sortedUnitIds = Object.keys(groupedResources).map(Number).sort((a, b) => {
    if (a === user?.unit_id) return -1;
    if (b === user?.unit_id) return 1;
    const indexA = units.findIndex(u => u.id === a);
    const indexB = units.findIndex(u => u.id === b);
    return indexA - indexB;
  });

  const availableWarehousesForNew = warehouses.filter(w => Number(w.unit_id) === Number(newRes.unit_id));
  const availableWarehousesForTransfer = transferModalData ? warehouses.filter(w => Number(w.unit_id) === Number(transferModalData.target_unit_id)) : [];
  
  // --- ЛОГІКА ІЄРАРХІЇ: Фільтрація людей для видачі ---
  const allowedUsersForAssignment = assignModalData 
    ? usersList.filter(u => u.role !== 'VOLUNTEER' && u.unit_id === assignModalData.resource.unit_id)
    : [];

  return (
    <div className="inventory-page">
      <Toaster position="top-right" />
      
      <div className="page-header">
        <h1>Облік ресурсу</h1>
        <div className="page-actions">
          {units.length > 0 && (
            <select 
              value={filterUnitId} 
              onChange={(e) => setFilterUnitId(e.target.value ? parseInt(e.target.value, 10) : '')} 
              className="filter-select erp-input"
            >
              <option value="">Всі підрозділи</option>
              {units.map((u) => (
                <option key={u.id} value={u.id}>{u.name}</option>
              ))}
            </select>
          )}
          
          {canManageCategories && (
            <button className="btn btn-secondary" onClick={() => setShowCategoryForm(true)}>
              + Категорія
            </button>
          )}
          
          {canManageResources && (
            <button className="btn btn-primary" onClick={() => setShowResourceForm(true)}>
              + Ресурс
            </button>
          )}
        </div>
      </div>

      {showCategoryForm && canManageCategories && (
        <div className="modal-overlay" onClick={() => setShowCategoryForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Нова категорія</h3>
            <form onSubmit={handleCreateCategory}>
              <div className="form-group">
                <label>Назва</label>
                <input className="erp-input" value={newCat.name} onChange={(e) => setNewCat({ ...newCat, name: e.target.value })} required />
              </div>
              <div className="form-group">
                <label>Опис</label>
                <input className="erp-input" value={newCat.description} onChange={(e) => setNewCat({ ...newCat, description: e.target.value })} />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => setShowCategoryForm(false)}>Скасувати</button>
                <button type="submit" className="btn btn-primary">Створити</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showResourceForm && canManageResources && (
        <div className="modal-overlay" onClick={() => setShowResourceForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Новий ресурс</h3>
            <form onSubmit={handleCreateResource}>
              <div className="form-group">
                <label>Категорія <span className="required">*</span></label>
                <select className="erp-input" value={newRes.category_id} onChange={(e) => setNewRes({ ...newRes, category_id: e.target.value })} required>
                  <option value="">Оберіть категорію</option>
                  {categories.map((c) => (<option key={c.id} value={c.id}>{c.name}</option>))}
                </select>
              </div>
              
              <div className="form-group">
                <label>Назва майна <span className="required">*</span></label>
                <input className="erp-input" placeholder="Напр. Бронежилет Корсар М3" value={newRes.name} onChange={(e) => setNewRes({ ...newRes, name: e.target.value })} required />
              </div>
              
              <div className="form-row-2">
                <div className="form-group">
                  <label>Власник (Підрозділ) <span className="required">*</span></label>
                  <select className="erp-input" value={newRes.unit_id ?? ''} onChange={(e) => { setNewRes({ ...newRes, unit_id: e.target.value ? parseInt(e.target.value, 10) : undefined, warehouse_id: '' }) }} required>
                    <option value="" disabled>Оберіть підрозділ</option>
                    {units.map((u) => (<option key={u.id} value={u.id}>{u.name}</option>))}
                  </select>
                </div>
                
                <div className="form-group">
                  <label>Склад (Фізична локація)</label>
                  <select className="erp-input" value={newRes.warehouse_id} onChange={(e) => setNewRes({ ...newRes, warehouse_id: e.target.value })}>
                    {!newRes.unit_id ? (
                      <option value="">-- Спершу оберіть підрозділ --</option>
                    ) : (
                      <>
                        <option value="">-- В дорозі / Не вказано --</option>
                        {availableWarehousesForNew.map((w) => (<option key={w.id} value={w.id}>{w.name}</option>))}
                      </>
                    )}
                  </select>
                </div>
              </div>
              
              <div className="form-row-2">
                <div className="form-group">
                  <label>Одиниці виміру</label>
                  <select className="erp-input" value={newRes.unit_type} onChange={(e) => setNewRes({ ...newRes, unit_type: e.target.value as 'PCS'|'KIT'|'KG'|'L' })}>
                    <option value="PCS">Штуки (шт)</option>
                    <option value="KIT">Комплекти (компл)</option>
                    <option value="KG">Кілограми (кг)</option>
                    <option value="L">Літри (л)</option>
                  </select>
                </div>
                <div className="form-group">
                  <label>Початкова кількість (На склад)</label>
                  <input className="erp-input" type="number" min="0" value={newRes.quantity.toString()} onChange={(e) => { const val = parseInt(e.target.value, 10); setNewRes({ ...newRes, quantity: isNaN(val) ? 0 : val }) }} required />
                </div>
              </div>
              
              <div className="form-row-2">
                <div className="form-group">
                  <label>Мін. залишок (для сповіщень)</label>
                  <input className="erp-input" type="number" min="0" value={newRes.min_quantity.toString()} onChange={(e) => { const val = parseInt(e.target.value, 10); setNewRes({ ...newRes, min_quantity: isNaN(val) ? 0 : val }) }} required />
                </div>
                <div className="form-group">
                  <label>Вага 1 од. (кг) <span className="required">*</span></label>
                  <input className="erp-input" type="number" min="0.1" step="0.1" value={newRes.weight_kg || ''} onChange={(e) => setNewRes({ ...newRes, weight_kg: parseFloat(e.target.value) })} required />
                </div>
              </div>
              
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => setShowResourceForm(false)}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={!newRes.unit_id}>Створити запис</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* МОДАЛКА: ВИДАЧА МАЙНА (ОНОВЛЕНО) */}
      {assignModalData && canManageResources && (
        <div className="modal-overlay" onClick={() => { setAssignModalData(null); setAssignError(null); }}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Видача майна персоналу</h3>
            <p className="text-muted text-left" style={{ marginBottom: '15px' }}>
              Видаємо: <strong>{assignModalData.resource.name}</strong><br />
              Власник: {units.find(u => u.id === assignModalData.resource.unit_id)?.name || 'Невідомо'}<br />
              Доступно на складі: {assignModalData.resource.quantity} {formatUnitType(assignModalData.resource.unit_type)}
            </p>
            {assignError && (<div className="modal-error-box">❌ {assignError}</div>)}
            <form onSubmit={handleAssignSubmit}>
              <div className="form-group text-left">
                <label>Кількість до видачі</label>
                <input 
                  className="erp-input" 
                  type="number" 
                  min="1" 
                  max={assignModalData.resource.quantity} 
                  value={assignModalData.quantity.toString()} 
                  onChange={(e) => { 
                    let val = parseInt(e.target.value, 10); 
                    if (isNaN(val)) val = 1; 
                    if (val > assignModalData.resource.quantity) val = assignModalData.resource.quantity; 
                    setAssignModalData({ ...assignModalData, quantity: val }) 
                  }} 
                  required 
                />
              </div>
              <div className="form-group text-left">
                <label>Співробітник (Отримувач)</label>
                <select 
                  className="erp-input" 
                  value={assignModalData.user_id} 
                  onChange={(e) => setAssignModalData({ ...assignModalData, user_id: e.target.value })} 
                  required
                >
                  <option value="" disabled>-- Оберіть особу --</option>
                  {allowedUsersForAssignment.length === 0 ? (
                    <option value="" disabled>У цьому підрозділі ще немає доданих співробітників!</option>
                  ) : (
                    allowedUsersForAssignment.map((u) => { 
                      const userUnit = units.find(unit => unit.id === u.unit_id); 
                      const unitName = userUnit ? userUnit.name : 'Без підрозділу'; 
                      return (<option key={u.id} value={u.id}>{u.full_name || u.email} — 🏢 {unitName}</option>);
                    })
                  )}
                </select>
                {allowedUsersForAssignment.length === 0 && (
                  <span className="error-text" style={{marginTop: '6px', display: 'block', fontSize: '12px'}}>
                    Майно можна видати ТІЛЬКИ співробітникам того підрозділу, якому воно належить.
                  </span>
                )}
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => { setAssignModalData(null); setAssignError(null); }}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={allowedUsersForAssignment.length === 0}>Підтвердити видачу</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {resourceToDelete && canManageResources && (
        <div className="modal-overlay" onClick={() => { setResourceToDelete(null); setDeleteError(null); }}>
          <div className="modal confirm-modal" onClick={(e) => e.stopPropagation()}>
            <h3>Підтвердження видалення</h3>
            {deleteError && (<div className="modal-error-box">❌ {deleteError}</div>)}
            <p className="confirm-text text-left">Ти дійсно хочеш <strong>назавжди видалити</strong> картку майна <strong>{resourceToDelete.name}</strong>?</p>
            <div className="warning-box">
              <p>Цей запис буде повністю стертий з бази даних разом з історією видач.</p>
              <p className="critical-text">Цю дію не можна скасувати.</p>
            </div>
            <div className="modal-actions">
              <button type="button" className="btn btn-secondary cancel-margin" onClick={() => { setResourceToDelete(null); setDeleteError(null); }}>Скасувати</button>
              <button type="button" className="btn btn-danger" onClick={confirmDelete}>🗑️ Видалити</button>
            </div>
          </div>
        </div>
      )}

      {editModalId && canManageResources && (
        <div className="modal-overlay" onClick={() => setEditModalId(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Редагування ресурсу</h3>
            <form onSubmit={handleEditSubmit}>
              <div className="form-group">
                <label>Назва ресурсу</label>
                <input className="erp-input" value={editForm.name} onChange={(e) => setEditForm({ ...editForm, name: e.target.value })} required />
              </div>
              <div className="form-group">
                <label>Мінімальний залишок (для складу)</label>
                <input className="erp-input" type="number" min="0" value={editForm.min_quantity.toString()} onChange={(e) => { const val = parseInt(e.target.value, 10); setEditForm({ ...editForm, min_quantity: isNaN(val) ? 0 : val }) }} required />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => setEditModalId(null)}>Скасувати</button>
                <button type="submit" className="btn btn-primary">Зберегти зміни</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {transferModalData && canManageResources && (
        <div className="modal-overlay" onClick={() => { setTransferModalData(null); setTransferError(null); }}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Переміщення ресурсу</h3>
            <p className="text-muted text-left" style={{ marginBottom: '15px' }}>
              Переміщення: <strong>{transferModalData.resource.name}</strong><br />
              Доступно на поточному складі: {transferModalData.resource.quantity} {formatUnitType(transferModalData.resource.unit_type)}
            </p>
            {transferError && (<div className="modal-error-box">❌ {transferError}</div>)}
            <form onSubmit={handleTransferSubmit}>
              <div className="form-group text-left">
                <label>Кількість для переміщення</label>
                <input className="erp-input" type="number" min="1" max={transferModalData.resource.quantity} value={transferModalData.quantity.toString()} onChange={(e) => { let val = parseInt(e.target.value, 10); if (isNaN(val)) val = 1; if (val > transferModalData.resource.quantity) val = transferModalData.resource.quantity; setTransferModalData({ ...transferModalData, quantity: val }) }} required />
              </div>
              <div className="form-group text-left">
                <label>Кому передаємо (Підрозділ)</label>
                <select className="erp-input" value={transferModalData.target_unit_id} onChange={(e) => { setTransferModalData({ ...transferModalData, target_unit_id: e.target.value ? parseInt(e.target.value, 10) : '', target_warehouse_id: '' }) }} required>
                  {units.map((u) => (<option key={u.id} value={u.id}>{u.name}</option>))}
                </select>
              </div>
              <div className="form-group text-left">
                <label>На який склад</label>
                <select className="erp-input" value={transferModalData.target_warehouse_id} onChange={(e) => setTransferModalData({ ...transferModalData, target_warehouse_id: e.target.value })}>
                  {transferModalData.target_unit_id === '' ? (
                    <option value="">-- Спершу оберіть підрозділ --</option>
                  ) : (
                    <>
                      <option value="">-- В дорозі / Не вказано --</option>
                      {availableWarehousesForTransfer.map((w) => (<option key={w.id} value={w.id}>{w.name}</option>))}
                    </>
                  )}
                </select>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => { setTransferModalData(null); setTransferError(null); }}>Скасувати</button>
                <button type="submit" className="btn btn-primary">Підтвердити переміщення</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {writeOffModalData && canManageResources && (
        <div className="modal-overlay" onClick={() => { setWriteOffModalData(null); setWriteOffError(null); }}>
          <div className="modal confirm-modal" onClick={(e) => e.stopPropagation()}>
            <h3>Списання майна (зі складу)</h3>
            <p className="text-muted" style={{ marginBottom: '15px', textAlign: 'left' }}>
              Списання: <strong>{writeOffModalData.resource.name}</strong><br />
              Доступно на складі: {writeOffModalData.resource.quantity} {formatUnitType(writeOffModalData.resource.unit_type)}
            </p>
            {writeOffError && (<div className="modal-error-box">❌ {writeOffError}</div>)}
            <form onSubmit={handleWriteOffSubmit}>
              <div className="form-group" style={{ textAlign: 'left' }}>
                <label>Скільки штук списати?</label>
                <input className="erp-input" type="number" min="1" max={writeOffModalData.resource.quantity} value={writeOffModalData.quantity.toString()} onChange={(e) => { let val = parseInt(e.target.value, 10); if (isNaN(val)) val = 1; if (val > writeOffModalData.resource.quantity) val = writeOffModalData.resource.quantity; setWriteOffModalData({ ...writeOffModalData, quantity: val }) }} required />
              </div>
              <div className="warning-box">
                <p>Вказана кількість майна отримає статус <strong>«Списано»</strong> і перейде у відповідну вкладку.</p>
                <p className="critical-text">Цю дію не можна скасувати.</p>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => { setWriteOffModalData(null); setWriteOffError(null); }}>Скасувати</button>
                <button type="submit" className="btn btn-danger">Підтвердити списання</button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className="content-grid">
        <div className="card">
          <h2>Категорії</h2>
          {categories.length === 0 ? (
            <p className="empty-state">Немає категорій</p>
          ) : (
            <ul className="category-list clean-list">
              <li onClick={() => setSelectedCategoryId(null)} className={`category-item ${selectedCategoryId === null ? 'active-category' : ''}`}>Всі категорії</li>
              {categories.map((c) => (
                <li key={c.id} onClick={() => setSelectedCategoryId(c.id)} className={`category-item ${selectedCategoryId === c.id ? 'active-category' : ''}`}>{c.name}</li>
              ))}
            </ul>
          )}
        </div>
        
        <div className={`card card-table ${isRefreshing ? 'refreshing-fade' : ''}`}>
          <div className="table-header-with-tabs">
            <h2>Ресурси {selectedCategoryId && `(${categories.find(c => c.id === selectedCategoryId)?.name})`}</h2>
            <div className="inventory-tabs">
              <button className={`tab-btn ${activeTab === 'active' ? 'active' : ''}`} onClick={() => setActiveTab('active')}>На балансі</button>
              <button className={`tab-btn ${activeTab === 'written_off' ? 'active' : ''}`} onClick={() => setActiveTab('written_off')}>Списані</button>
            </div>
          </div>

          {filteredResources.length === 0 ? (
            <p className="empty-state">{activeTab === 'active' ? 'Немає ресурсів' : 'Списаних ресурсів немає'}</p>
          ) : (
          <table className="data-table table-inventory">
            <thead>
              <tr>
                <th>Назва</th>
                <th>Склад (Локація)</th>
                <th className="text-center">Загальна кількість</th>
                <th className="text-center">Мін. залишок</th>
                <th className="text-center">Стан</th>
                {canManageResources && <th>Дії</th>}
              </tr>
            </thead>
            <tbody>
                {sortedUnitIds.map(unitId => {
                  const unitResources = groupedResources[unitId];
                  const isOrphan = unitId === 0;
                  const unitName = isOrphan ? (activeTab === 'active' ? '⚠️ НЕРОЗПОДІЛЕНИЙ ЗАЛИШОК (Помилка обліку: призначте підрозділ)' : '🗄️ Старий архів (Без підрозділу)') : units.find((u) => u.id === unitId)?.name || 'Невідомий підрозділ';
                  const isMyUnit = user?.unit_id === unitId;
                  const isDangerRow = isOrphan && activeTab === 'active';

                  return (
                    <React.Fragment key={unitId}>
                      <tr className={`unit-header-row ${isMyUnit ? 'my-unit-header' : ''} ${isDangerRow ? 'bg-orphan' : ''}`}>
                        <td colSpan={canManageResources ? 6 : 5} className={isDangerRow ? 'text-danger fw-bold' : ''}>
                          {isDangerRow ? '🚨 ' : (isOrphan ? '' : '🏢 На балансі: ')} {unitName} {isMyUnit && <span className="my-unit-badge">(Ваш підрозділ)</span>}
                        </td>
                      </tr>
                      {unitResources.map((r) => {
                        const isWrittenOff = r.condition === 'WRITTEN_OFF';
                        const issuedQty = (r as any).issued_quantity || 0;
                        const totalQuantity = r.quantity + issuedQty;
                        
                        let status = 'success'; let statusText = 'OK';
                        if (isWrittenOff) { status = 'neutral'; statusText = 'Списано'; } 
                        else if (r.quantity === 0 && issuedQty === 0) { status = 'critical'; statusText = 'Відсутньо'; } 
                        else if (r.quantity <= r.min_quantity) { status = 'warning'; statusText = 'Нестача'; }
                        
                        const warehouseNameStr = r.warehouse_id ? warehouses.find(w => w.id === r.warehouse_id)?.name || 'Невідомий склад' : 'В дорозі / Не вказано';

                        return (
                          <tr key={r.id} className={`row-${status} ${isWrittenOff ? 'written-off-row' : ''}`}>
                            <td className="resource-name-cell">{r.name}</td>
                            <td className="location-cell">
                              <div className="location-stack">
                                <span className="stock-info">🏢 Склад: <strong>{r.quantity} {formatUnitType(r.unit_type)}</strong> ({warehouseNameStr})</span>
                                {issuedQty > 0 && !isWrittenOff && (<span className="issued-info">👤 На руках у о/с: <strong>{issuedQty} {formatUnitType(r.unit_type)}</strong></span>)}
                              </div>
                            </td>
                            <td className={isWrittenOff ? 'written-off-qty text-center font-bold' : 'text-center font-bold'}>
                              {totalQuantity} <small className="text-muted">{formatUnitType(r.unit_type)}</small>
                            </td>
                            <td className="text-center">{r.min_quantity}</td>
                            <td className="text-center"><span className={`badge badge-${status}`}>{statusText}</span></td>
                            
                            {canManageResources && (
                              <td>
                                <div className="dropdown-container" onClick={(e) => e.stopPropagation()}>
                                  <button className={`btn-kebab ${activeMenuId === r.id ? 'active' : ''}`} onClick={() => setActiveMenuId(activeMenuId === r.id ? null : r.id)}>⋮</button>
                                  {activeMenuId === r.id && (
                                    <div className="actions-dropdown-menu">
                                      {!isWrittenOff && (
                                        <>
                                          {r.quantity > 0 && (
                                            <button className="text-primary-action" onClick={() => { setAssignModalData({ resource: r, quantity: 1, user_id: '' }); setActiveMenuId(null); }}>
                                              👤 Видати співробітнику
                                            </button>
                                          )}
                                          <button onClick={() => { setEditForm({ name: r.name, min_quantity: r.min_quantity }); setEditModalId(r.id); setActiveMenuId(null); }}>
                                            ✏️ Редагувати
                                          </button>
                                          <button onClick={() => { setTransferModalData({ resource: r, quantity: 1, target_unit_id: r.unit_id || '', target_warehouse_id: '' }); setActiveMenuId(null); }}>
                                            🔄 Передати (Трансфер)
                                          </button>
                                          <button onClick={() => { setWriteOffModalData({ resource: r, quantity: r.quantity }); setActiveMenuId(null); }}>
                                            📦 Списати зі складу
                                          </button>
                                          <div className="dropdown-divider"></div>
                                        </>
                                      )}
                                      <button className="text-danger" onClick={() => { setResourceToDelete(r); setActiveMenuId(null); }}>
                                        🗑️ Видалити запис
                                      </button>
                                    </div>
                                  )}
                                </div>
                              </td>
                            )}
                          </tr>
                        )
                      })}
                    </React.Fragment>
                  )
                })}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}