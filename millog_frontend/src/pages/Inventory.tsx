import React, { useEffect, useState } from 'react';
import { api, type Resource, type ResourceCategory, type Unit, type Warehouse } from '../api/client';
import { useAuth } from '../contexts/AuthContext';
import { usePermissions } from '../hooks/usePermissions';
import { PaywallBadge } from '../components/FeatureGate';
import toast, { Toaster } from 'react-hot-toast';
import { Html5QrcodeScanner } from 'html5-qrcode';
import './Inventory.css';

export default function Inventory() {
  const { user } = useAuth();
  const perms = usePermissions();
  
  const [categories, setCategories] = useState<ResourceCategory[]>([]);
  const [resources, setResources] = useState<Resource[]>([]);
  const [units, setUnits] = useState<Unit[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [filterUnitId, setFilterUnitId] = useState<number | ''>('');
  
  const [loading, setLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);
  
  const [showCategoryForm, setShowCategoryForm] = useState(false);
  const [showResourceForm, setShowResourceForm] = useState(false);
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
  
  const [editingCategory, setEditingCategory] = useState<ResourceCategory | null>(null);
  const [editCategoryForm, setEditCategoryForm] = useState({ name: '', description: '' });
  const [deletingCategory, setDeletingCategory] = useState<ResourceCategory | null>(null);

  const [isScannerOpen, setIsScannerOpen] = useState(false);
  const [scannedResource, setScannedResource] = useState<any | null>(null);

  const [newCat, setNewCat] = useState({ name: '', description: '' });
  
  const [newRes, setNewRes] = useState({ 
    category_id: '', 
    unit_id: undefined as number | undefined, 
    warehouse_id: '',
    name: '',
    barcode: '', 
    quantity: 0, 
    unit_type: 'PCS' as 'PCS' | 'KIT' | 'KG' | 'L',
    min_quantity: 0,
    weight_kg: 1
  });

  const [selectedCategoryId, setSelectedCategoryId] = useState<string | null>(null);
  const [activeMenuId, setActiveMenuId] = useState<string | null>(null);
  
  // ---------------------------------------------------------
  // НОВИЙ СТЕЙТ ДЛЯ ПОШУКУ
  // ---------------------------------------------------------
  const [searchQuery, setSearchQuery] = useState('');
  // --- СТЕЙТ ДЛЯ ІМПОРТУ ЕКСЕЛЬ ---
  const [showImportModal, setShowImportModal] = useState(false);
  const [importFile, setImportFile] = useState<File | null>(null);
  const [importUnitId, setImportUnitId] = useState<number | ''>('');
  const [importWarehouseId, setImportWarehouseId] = useState<string>('');
  const [isImporting, setIsImporting] = useState(false);

  useEffect(() => {
    const closeMenu = () => setActiveMenuId(null);
    document.addEventListener('click', closeMenu);
    return () => document.removeEventListener('click', closeMenu);
  }, []);

  const canManageResources = perms.can('resource_manage');
  const canManageCategories = perms.can('category_manage');
  const hasExcelImport = perms.hasFeature('excel_import');
  
  const [qrPreviewData, setQrPreviewData] = useState<{ id: string; name: string; url: string } | null>(null);

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

  useEffect(() => {
    if (isScannerOpen) {
      const scanner = new Html5QrcodeScanner(
        "qr-reader",
        { fps: 10, qrbox: { width: 250, height: 250 } },
        false
      );

      scanner.render(
        async (decodedText) => {
          scanner.clear();
          setIsScannerOpen(false);
          
          if (decodedText.startsWith('millog-resource:')) {
            const resourceId = decodedText.split(':')[1];
            const toastId = toast.loading('Шукаємо майно в базі...');
            try {
              const resourceData = await api.inventory.getById(resourceId);
              setScannedResource(resourceData);
              toast.success('Майно знайдено!', { id: toastId });
            } catch (err) {
              toast.error('Майно не знайдено або видалено', { id: toastId });
            }
          } else {
            toast.error('Невідомий формат QR-коду');
          }
        },
        (error) => {
          console.warn('Помилка сканування QR-коду:', error);
        }
      );

      return () => {
        scanner.clear().catch(e => console.error("Помилка зупинки сканера", e));
      };
    }
  }, [isScannerOpen]);

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

  const handleEditCategorySubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingCategory) return;
    try {
      await api.inventory.updateCategory(editingCategory.id, editCategoryForm);
      toast.success('Категорію успішно оновлено!');
      setEditingCategory(null);
      loadData();
    } catch (err: any) {
      toast.error(err.message || 'Помилка при оновленні категорії');
    }
  };

  const confirmDeleteCategory = async () => {
    if (!deletingCategory) return;
    try {
      await api.inventory.deleteCategory(deletingCategory.id);
      toast.success('Категорію видалено!');
      setDeletingCategory(null);
      if (selectedCategoryId === deletingCategory.id) setSelectedCategoryId(null);
      loadData();
    } catch (err: any) {
      toast.error(err.message || 'Неможливо видалити: у цій категорії ще є майно');
      setDeletingCategory(null);
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
        barcode: '', 
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

  const handleDownloadQR = async (resourceId: string, resourceName: string) => {
    const toastId = toast.loading('Генерація QR-коду...');
    try {
      const blob = await api.inventory.downloadResourceQR(resourceId);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      const safeName = resourceName.replace(/\s+/g, '_');
      link.setAttribute('download', `QR_${safeName}.png`);
      document.body.appendChild(link);
      link.click();
      link.parentNode?.removeChild(link);
      window.URL.revokeObjectURL(url);
      toast.success('Наклейку завантажено!', { id: toastId });
    } catch (error) {
      toast.error('Помилка генерації QR', { id: toastId });
    }
  };

  const handleShowQR = async (resourceId: string, resourceName: string) => {
    const toastId = toast.loading('Генерація QR-коду...');
    try {
      const blob = await api.inventory.downloadResourceQR(resourceId);
      const url = window.URL.createObjectURL(blob);
      setQrPreviewData({ id: resourceId, name: resourceName, url });
      toast.dismiss(toastId);
    } catch (error) {
      toast.error('Помилка генерації QR', { id: toastId });
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

  const handleDownloadTemplate = async () => {
    try {
      await api.inventory.downloadImportTemplate();
      toast.success('Шаблон завантажено');
    } catch (err) {
      toast.error('Не вдалося завантажити шаблон');
    }
  };

  const handleImportSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!importFile || !importUnitId || !importWarehouseId) {
      toast.error('Заповніть усі поля та оберіть файл');
      return;
    }

    setIsImporting(true);
    const toastId = toast.loading('Імпортуємо дані з Excel...');
    
    try {
      const result = await api.inventory.importExcel(
        Number(importUnitId), 
        importWarehouseId, 
        importFile
      );
      toast.success(`Успіх! Додано ${result.success_count} позицій`, { id: toastId });
      setShowImportModal(false);
      setImportFile(null);
      loadData(); // Оновлюємо таблицю, щоб побачити нові товари
    } catch (err: any) {
      toast.error(err.message || 'Помилка при імпорті', { id: toastId });
    } finally {
      setIsImporting(false);
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

  // ---------------------------------------------------------
  // ОНОВЛЕНА ЛОГІКА ФІЛЬТРАЦІЇ (Додано пошук)
  // ---------------------------------------------------------
  const filteredResources = resources.filter(r => {
    const isDeleted = r.condition === 'WRITTEN_OFF';
    const matchesTab = activeTab === 'active' ? !isDeleted : isDeleted;
    const matchesCategory = selectedCategoryId ? r.category_id === selectedCategoryId : true;
    const isOrphanEmpty = (r.unit_id === 0 || r.unit_id == null) && r.quantity === 0 && !isDeleted;
    
    // Перевірка пошукового запиту
    const query = searchQuery.toLowerCase().trim();
    const matchesSearch = query === '' || 
                          r.name.toLowerCase().includes(query) || 
                          r.id.toLowerCase().includes(query); // Дозволяє знайти по шматку UUID

    return matchesTab && matchesCategory && !isOrphanEmpty && matchesSearch;
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
  
  const allowedUsersForAssignment = assignModalData 
    ? usersList.filter(u => u.role !== 'CONTRACTOR' && u.unit_id === assignModalData.resource.unit_id)
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
              <option value="">Всі орг. одиниці</option>
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
          
          <button className="btn btn-secondary" onClick={() => setIsScannerOpen(true)}>
            📷 Сканувати QR
          </button>
          
          {canManageResources && (
            hasExcelImport ? (
              <button className="btn btn-secondary" onClick={() => setShowImportModal(true)}>
                📥 Імпорт Excel
              </button>
            ) : (
              <button
                className="btn btn-secondary"
                onClick={() => toast('Імпорт Excel доступний у тарифі PRO', { icon: '🔒' })}
                title="Доступно в PRO"
                style={{ opacity: 0.7 }}
              >
                🔒 Імпорт Excel <PaywallBadge feature="excel_import" compact />
              </button>
            )
          )}

          {canManageResources && (
            <button className="btn btn-primary" onClick={() => setShowResourceForm(true)}>
              + Ресурс
            </button>
          )}
        </div>
      </div>

      {/* ============================== МОДАЛКИ КАТЕГОРІЙ ============================== */}
      {/* ... (Модалки категорій залишені без змін) ... */}
      {showCategoryForm && canManageCategories && (
        <div className="modal-overlay inventory-modal" onClick={() => setShowCategoryForm(false)}>
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
                <button type="button" className="btn btn-secondary" onClick={() => setShowCategoryForm(false)}>Скасувати</button>
                <button type="submit" className="btn btn-primary">Створити</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {editingCategory && canManageCategories && (
        <div className="modal-overlay inventory-modal" onClick={() => setEditingCategory(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Редагувати категорію</h3>
            <form onSubmit={handleEditCategorySubmit}>
              <div className="form-group">
                <label>Назва</label>
                <input className="erp-input" value={editCategoryForm.name} onChange={(e) => setEditCategoryForm({ ...editCategoryForm, name: e.target.value })} required />
              </div>
              <div className="form-group">
                <label>Опис</label>
                <input className="erp-input" value={editCategoryForm.description} onChange={(e) => setEditCategoryForm({ ...editCategoryForm, description: e.target.value })} />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setEditingCategory(null)}>Скасувати</button>
                <button type="submit" className="btn btn-primary">Зберегти</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {deletingCategory && canManageCategories && (
        <div className="modal-overlay inventory-modal" onClick={() => setDeletingCategory(null)}>
          <div className="modal confirm-modal" onClick={(e) => e.stopPropagation()}>
            <h3 style={{ color: '#ef4444' }}>⚠️ Видалення категорії</h3>
            <p className="confirm-text text-left">
              Видалити категорію <strong>{deletingCategory.name}</strong>?<br/>
              <small style={{color: '#64748b'}}>Це можливо лише якщо в категорії немає майна.</small>
            </p>
            <div className="modal-actions">
              <button type="button" className="btn btn-secondary" onClick={() => setDeletingCategory(null)}>Скасувати</button>
              <button type="button" className="btn btn-danger" onClick={confirmDeleteCategory}>🗑️ Видалити</button>
            </div>
          </div>
        </div>
      )}

      {/* ============================== ІНШІ МОДАЛКИ (Майно) ============================== */}
      {/* ... (Всі інші модалки залишені без змін) ... */}
      {showResourceForm && canManageResources && (
        <div className="modal-overlay inventory-modal" onClick={() => setShowResourceForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Новий ресурс</h3>
            <form onSubmit={handleCreateResource}>
              <div className="form-group">
                <label>Категорія <span style={{color: '#ef4444'}}>*</span></label>
                <select className="erp-input" value={newRes.category_id} onChange={(e) => setNewRes({ ...newRes, category_id: e.target.value })} required>
                  <option value="">Оберіть категорію</option>
                  {categories.map((c) => (<option key={c.id} value={c.id}>{c.name}</option>))}
                </select>
              </div>
              
              <div className="form-group">
                <label>Назва майна <span style={{color: '#ef4444'}}>*</span></label>
                <input className="erp-input" placeholder="Напр. Ноутбук Dell Latitude" value={newRes.name} onChange={(e) => setNewRes({ ...newRes, name: e.target.value })} required />
              </div>

              <div className="form-group">
                <label>Заводський штрих-код (необов'язково)</label>
                <input 
                  className="erp-input" 
                  placeholder="Відскануйте код або введіть вручну..." 
                  value={newRes.barcode} 
                  onChange={(e) => setNewRes({ ...newRes, barcode: e.target.value })} 
                />
              </div>
              
              <div className="form-row-2">
                <div className="form-group">
                  <label>Власник (Орг. одиниця) <span style={{color: '#ef4444'}}>*</span></label>
                  <select className="erp-input" value={newRes.unit_id ?? ''} onChange={(e) => { setNewRes({ ...newRes, unit_id: e.target.value ? parseInt(e.target.value, 10) : undefined, warehouse_id: '' }) }} required>
                    <option value="" disabled>Оберіть орг. одиницю</option>
                    {units.map((u) => (<option key={u.id} value={u.id}>{u.name}</option>))}
                  </select>
                </div>
                
                <div className="form-group">
                  <label>Склад (Локація) <span style={{color: '#ef4444'}}>*</span></label>
                  <select className="erp-input" value={newRes.warehouse_id} onChange={(e) => setNewRes({ ...newRes, warehouse_id: e.target.value })} required>
                    <option value="" disabled>-- Оберіть конкретний склад --</option>
                    {availableWarehousesForNew.map((w) => (<option key={w.id} value={w.id}>{w.name}</option>))}
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
                  <label>Кількість (На склад)</label>
                  <input className="erp-input" type="number" min="0" value={newRes.quantity.toString()} onChange={(e) => { const val = parseInt(e.target.value, 10); setNewRes({ ...newRes, quantity: isNaN(val) ? 0 : val }) }} required />
                </div>
              </div>
              
              <div className="form-row-2">
                <div className="form-group">
                  <label>Мін. залишок</label>
                  <input className="erp-input" type="number" min="0" value={newRes.min_quantity.toString()} onChange={(e) => { const val = parseInt(e.target.value, 10); setNewRes({ ...newRes, min_quantity: isNaN(val) ? 0 : val }) }} required />
                </div>
                <div className="form-group">
                  <label>Вага 1 од. (кг) <span style={{color: '#ef4444'}}>*</span></label>
                  <input className="erp-input" type="number" min="0.1" step="0.1" value={newRes.weight_kg || ''} onChange={(e) => setNewRes({ ...newRes, weight_kg: parseFloat(e.target.value) })} required />
                </div>
              </div>
              
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowResourceForm(false)}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={!newRes.unit_id}>Створити запис</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {assignModalData && canManageResources && (
        <div className="modal-overlay inventory-modal" onClick={() => { setAssignModalData(null); setAssignError(null); }}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Видача майна персоналу</h3>
            <p style={{color: '#64748b', textAlign: 'left'}}>
              Видаємо: <strong style={{color: '#0f172a'}}>{assignModalData.resource.name}</strong><br />
              Доступно: {assignModalData.resource.quantity} {formatUnitType(assignModalData.resource.unit_type)}
            </p>
            {assignError && (<div style={{color: '#dc2626', padding: '10px', background: '#fef2f2', borderRadius: '6px', marginBottom: '16px'}}>❌ {assignError}</div>)}
            <form onSubmit={handleAssignSubmit}>
              <div className="form-group text-left">
                <label>Кількість до видачі</label>
                <input className="erp-input" type="number" min="1" max={assignModalData.resource.quantity} value={assignModalData.quantity.toString()} onChange={(e) => { let val = parseInt(e.target.value, 10); if (isNaN(val)) val = 1; if (val > assignModalData.resource.quantity) val = assignModalData.resource.quantity; setAssignModalData({ ...assignModalData, quantity: val }) }} required />
              </div>
              <div className="form-group text-left">
                <label>Отримувач</label>
                <select className="erp-input" value={assignModalData.user_id} onChange={(e) => setAssignModalData({ ...assignModalData, user_id: e.target.value })} required>
                  <option value="" disabled>-- Оберіть особу --</option>
                  {allowedUsersForAssignment.map((u) => (
                    <option key={u.id} value={u.id}>{u.full_name || u.email}</option>
                  ))}
                </select>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => { setAssignModalData(null); setAssignError(null); }}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={allowedUsersForAssignment.length === 0}>Видати</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {resourceToDelete && canManageResources && (
        <div className="modal-overlay inventory-modal" onClick={() => { setResourceToDelete(null); setDeleteError(null); }}>
          <div className="modal confirm-modal" onClick={(e) => e.stopPropagation()}>
            <h3 style={{color: '#ef4444'}}>Підтвердження видалення</h3>
            {deleteError && (<div style={{color: '#dc2626', padding: '10px', background: '#fef2f2', borderRadius: '6px', marginBottom: '16px'}}>❌ {deleteError}</div>)}
            <p className="confirm-text text-left">Видалити <strong>{resourceToDelete.name}</strong>?</p>
            <div className="modal-actions">
              <button type="button" className="btn btn-secondary" onClick={() => { setResourceToDelete(null); setDeleteError(null); }}>Скасувати</button>
              <button type="button" className="btn btn-danger" onClick={confirmDelete}>🗑️ Видалити</button>
            </div>
          </div>
        </div>
      )}

      {editModalId && canManageResources && (
        <div className="modal-overlay inventory-modal" onClick={() => setEditModalId(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Редагування ресурсу</h3>
            <form onSubmit={handleEditSubmit}>
              <div className="form-group">
                <label>Назва ресурсу</label>
                <input className="erp-input" value={editForm.name} onChange={(e) => setEditForm({ ...editForm, name: e.target.value })} required />
              </div>
              <div className="form-group">
                <label>Мінімальний залишок</label>
                <input className="erp-input" type="number" min="0" value={editForm.min_quantity.toString()} onChange={(e) => { const val = parseInt(e.target.value, 10); setEditForm({ ...editForm, min_quantity: isNaN(val) ? 0 : val }) }} required />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setEditModalId(null)}>Скасувати</button>
                <button type="submit" className="btn btn-primary">Зберегти зміни</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {writeOffModalData && canManageResources && (
        <div className="modal-overlay inventory-modal" onClick={() => { setWriteOffModalData(null); setWriteOffError(null); }}>
          <div className="modal confirm-modal" onClick={(e) => e.stopPropagation()}>
            <h3 style={{color: '#ef4444'}}>Списання майна</h3>
            {writeOffError && (<div style={{color: '#dc2626', padding: '10px', background: '#fef2f2', borderRadius: '6px', marginBottom: '16px'}}>❌ {writeOffError}</div>)}
            <form onSubmit={handleWriteOffSubmit}>
              <div className="form-group" style={{ textAlign: 'left' }}>
                <label>Кількість до списання</label>
                <input className="erp-input" type="number" min="1" max={writeOffModalData.resource.quantity} value={writeOffModalData.quantity.toString()} onChange={(e) => { let val = parseInt(e.target.value, 10); if (isNaN(val)) val = 1; if (val > writeOffModalData.resource.quantity) val = writeOffModalData.resource.quantity; setWriteOffModalData({ ...writeOffModalData, quantity: val }) }} required />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => { setWriteOffModalData(null); setWriteOffError(null); }}>Скасувати</button>
                <button type="submit" className="btn btn-danger">Списати</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {qrPreviewData && (
        <div className="modal-overlay inventory-modal" onClick={() => {
          window.URL.revokeObjectURL(qrPreviewData.url);
          setQrPreviewData(null);
        }}>
          <div className="modal qr-modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Наклейка майна</h3>
              <button className="close-btn" onClick={() => setQrPreviewData(null)} style={{background: 'none', border: 'none', fontSize: '24px', cursor: 'pointer', color: '#64748b'}}>&times;</button>
            </div>
            <div className="qr-modal-body">
              <p className="qr-resource-name">{qrPreviewData.name}</p>
              <div className="qr-image-wrapper">
                <img src={qrPreviewData.url} alt="Resource QR Code" className="qr-image" />
              </div>
              <p className="qr-id-info">
                ID: <code className="qr-id-code">{qrPreviewData.id.split('-')[0].toUpperCase()}</code>
              </p>
            </div>
            <div className="modal-actions qr-modal-actions">
              <button className="btn btn-primary btn-block" onClick={() => {
                const link = document.createElement('a');
                link.href = qrPreviewData.url;
                link.setAttribute('download', `QR_${qrPreviewData.name.replace(/\s+/g, '_')}.png`);
                document.body.appendChild(link);
                link.click();
                link.remove();
              }}>📥 Завантажити для друку</button>
              <button className="btn btn-secondary btn-block" onClick={() => setQrPreviewData(null)}>Закрити</button>
            </div>
          </div>
        </div>
      )}

      {isScannerOpen && (
        <div className="modal-overlay inventory-modal" onClick={() => setIsScannerOpen(false)}>
          <div className="modal scanner-modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Наведіть камеру на QR-код</h3>
              <button className="close-btn" onClick={() => setIsScannerOpen(false)} style={{background: 'none', border: 'none', fontSize: '24px', cursor: 'pointer', color: '#64748b'}}>&times;</button>
            </div>
            <div className="scanner-modal-body">
              <div id="qr-reader" className="qr-reader-container"></div>
              <p className="scanner-permission-text">
                Дозвольте браузеру доступ до камери
              </p>
            </div>
          </div>
        </div>
      )}

      {scannedResource && (
        <div className="modal-overlay inventory-modal" onClick={() => setScannedResource(null)}>
          <div className="modal scan-result-modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Картка майна (Скан)</h3>
              <button className="close-btn" onClick={() => setScannedResource(null)} style={{background: 'none', border: 'none', fontSize: '24px', cursor: 'pointer', color: '#64748b'}}>&times;</button>
            </div>
            
            <div className="scanner-modal-body">
              <div className="scan-result-icon">
                <span>📦</span>
                <h2>{scannedResource.name}</h2>
              </div>
              
              <table className="scan-result-table">
                <tbody>
                  <tr>
                    <td className="scan-result-label">Залишок на складі:</td>
                    <td className="scan-result-value">
                      {scannedResource.quantity} {formatUnitType(scannedResource.unit_type)}
                    </td>
                  </tr>
                  <tr>
                    <td className="scan-result-label">Стан:</td>
                    <td className="scan-result-value">
                      <span className="badge badge-success">Нове / Справне</span>
                    </td>
                  </tr>
                  <tr>
                    <td className="scan-result-label">ID в системі:</td>
                    <td className="scan-result-id">
                      {scannedResource.id.split('-')[0].toUpperCase()}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div className="modal-actions">
              <button className="btn btn-primary btn-block" style={{width: '100%', justifyContent: 'center'}} onClick={() => setScannedResource(null)}>
                Закрити картку
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ============================== ОСНОВНИЙ КОНТЕНТ (СІТКА) ============================== */}
      <div className="content-grid">
        
        {/* КАТЕГОРІЇ (Сайдбар) */}
        <div className="card">
          <h2>Категорії</h2>
          {categories.length === 0 ? (
            <p className="empty-state">Немає категорій</p>
          ) : (
            <ul className="category-list clean-list">
              <li onClick={() => setSelectedCategoryId(null)} className={`category-item ${selectedCategoryId === null ? 'active-category' : ''}`}>
                <span className="category-name-text">Всі категорії</span>
              </li>
              
              {categories.map((c) => (
                <li 
                  key={c.id} 
                  onClick={() => setSelectedCategoryId(c.id)} 
                  className={`category-item ${selectedCategoryId === c.id ? 'active-category' : ''}`}
                >
                  <span className="category-name-text">{c.name}</span>
                  
                  {canManageCategories && (
                    <div className="category-actions">
                      <button
                        type="button"
                        className="btn-icon-small"
                        onClick={(e) => { 
                          e.stopPropagation(); 
                          setEditingCategory(c); 
                          setEditCategoryForm({ name: c.name, description: c.description || '' }); 
                        }}
                        title="Редагувати"
                      >
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                      </button>
                      
                      <button
                        type="button"
                        className="btn-icon-small delete"
                        onClick={(e) => { 
                          e.stopPropagation(); 
                          setDeletingCategory(c); 
                        }}
                        title="Видалити"
                      >
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                      </button>
                    </div>
                  )}
                </li>
              ))}
            </ul>
          )}
        </div>
        
        {/* ТАБЛИЦЯ МАЙНА */}
        <div className={`card card-table ${isRefreshing ? 'refreshing-fade' : ''}`}>
          
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px', flexWrap: 'wrap', gap: '15px' }}>
            
            {/* Ліва частина: Заголовок + Пошук */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '20px', flexWrap: 'wrap' }}>
              <h2 style={{ margin: 0 }}>
                Ресурси {selectedCategoryId && `(${categories.find(c => c.id === selectedCategoryId)?.name})`}
              </h2>
              
              {/* ПОЛЕ ПОШУКУ */}
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
                    paddingTop: '6px', 
                    paddingBottom: '6px', 
                    height: 'auto',
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
            </div>

            {/* Права частина: Вкладки (Таби) */}
            <div style={{ display: 'flex', gap: '8px', borderBottom: '2px solid #e2e8f0' }}>
              <button 
                style={{ background: 'none', border: 'none', padding: '8px 16px', fontWeight: 600, cursor: 'pointer', color: activeTab === 'active' ? '#2563eb' : '#64748b', borderBottom: activeTab === 'active' ? '2px solid #2563eb' : '2px solid transparent', marginBottom: '-2px' }} 
                onClick={() => setActiveTab('active')}
              >
                На балансі
              </button>
              <button 
                style={{ background: 'none', border: 'none', padding: '8px 16px', fontWeight: 600, cursor: 'pointer', color: activeTab === 'written_off' ? '#2563eb' : '#64748b', borderBottom: activeTab === 'written_off' ? '2px solid #2563eb' : '2px solid transparent', marginBottom: '-2px' }} 
                onClick={() => setActiveTab('written_off')}
              >
                Списані
              </button>
            </div>
          </div>

          {/* ============================== МОДАЛКА ІМПОРТУ EXCEL ============================== */}
      {showImportModal && (
        <div className="modal-overlay inventory-modal" onClick={() => !isImporting && setShowImportModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Масове завантаження майна</h3>
              <button className="close-btn" onClick={() => setShowImportModal(false)}>&times;</button>
            </div>
            
            <form onSubmit={handleImportSubmit}>
              <div className="info-box" style={{ background: '#f0f9ff', padding: '12px', borderRadius: '8px', marginBottom: '16px', fontSize: '0.9rem', border: '1px solid #bae6fd' }}>
                💡 Завантажте шаблон, заповніть його даними та перетягніть файл сюди. 
                Система автоматично створить нові записи в базі.
                <button type="button" onClick={handleDownloadTemplate} style={{ display: 'block', marginTop: '8px', color: '#0284c7', background: 'none', border: 'none', textDecoration: 'underline', cursor: 'pointer', padding: 0 }}>
                  📥 Завантажити шаблон .xlsx
                </button>
              </div>

              <div className="form-row-2">
                <div className="form-group">
                  {/* Замінили "Підрозділ-власник" */}
                  <label>Власник (Орг. одиниця)</label> 
                  <select className="erp-input" value={importUnitId} onChange={(e) => { setImportUnitId(Number(e.target.value)); setImportWarehouseId(''); }} required>
                    {/* Замінили "-- Оберіть підрозділ --" */}
                    <option value="">-- Оберіть орг. одиницю --</option> 
                    {units.map((u) => (<option key={u.id} value={u.id}>{u.name}</option>))}
                  </select>
                </div>
                
                <div className="form-group">
                  <label>Цільовий склад</label>
                  <select className="erp-input" value={importWarehouseId} onChange={(e) => setImportWarehouseId(e.target.value)} required disabled={!importUnitId}>
                    <option value="">-- Оберіть склад --</option>
                    {warehouses.filter(w => Number(w.unit_id) === Number(importUnitId)).map((w) => (
                      <option key={w.id} value={w.id}>{w.name}</option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="form-group">
                <label>Файл Excel (.xlsx)</label>
                <div 
                  className={`file-drop-zone ${importFile ? 'has-file' : ''}`}
                  onDragOver={(e) => e.preventDefault()}
                  onDrop={(e) => {
                    e.preventDefault();
                    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
                      setImportFile(e.dataTransfer.files[0]);
                    }
                  }}
                >
                  <input 
                    type="file" 
                    accept=".xlsx" 
                    id="excel-upload"
                    hidden 
                    onChange={(e) => setImportFile(e.target.files ? e.target.files[0] : null)} 
                  />
                  <label htmlFor="excel-upload" style={{ cursor: 'pointer', display: 'block', padding: '30px' }}>
                    {importFile ? `📄 ${importFile.name}` : 'Перетягніть файл сюди або натисніть для вибору'}
                  </label>
                </div>
              </div>

              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowImportModal(false)} disabled={isImporting}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={isImporting || !importFile}>
                  {isImporting ? 'Завантаження...' : '🚀 Розпочати імпорт'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

          {/* Відображення результатів */}
          {filteredResources.length === 0 ? (
            <div className="empty-state">
              {searchQuery 
                ? `За запитом "${searchQuery}" нічого не знайдено` 
                : (activeTab === 'active' ? 'Немає ресурсів' : 'Списаних ресурсів немає')
              }
            </div>
          ) : (
          <table className="data-table table-inventory">
            <thead>
              <tr>
                <th>Назва</th>
                <th>Склад (Локація)</th>
                <th style={{textAlign: 'center'}}>Загальна кількість</th>
                <th style={{textAlign: 'center'}}>Мін. залишок</th>
                <th style={{textAlign: 'center'}}>Стан</th>
                {canManageResources && <th className="col-actions-menu">Дії</th>}
              </tr>
            </thead>
            <tbody>
                {sortedUnitIds.map(unitId => {
                  const unitResources = groupedResources[unitId];
                  const isOrphan = unitId === 0;
                  const unitName = isOrphan ? (activeTab === 'active' ? '⚠️ НЕРОЗПОДІЛЕНИЙ ЗАЛИШОК' : '🗄️ Архів') : units.find((u) => u.id === unitId)?.name || 'Невідома орг. одиниця';
                  const isMyUnit = user?.unit_id === unitId;
                  const isDangerRow = isOrphan && activeTab === 'active';

                  return (
                    <React.Fragment key={unitId}>
                      <tr className={`unit-header-row ${isMyUnit ? 'my-unit-header' : ''}`}>
                        <td colSpan={canManageResources ? 6 : 5} style={{color: isDangerRow ? '#ef4444' : ''}}>
                          {isDangerRow ? '🚨 ' : (isOrphan ? '' : '🏢 ')} {unitName} {isMyUnit && <span className="my-unit-badge">(Ваш)</span>}
                        </td>
                      </tr>
                      {unitResources.map((r) => {
                        const isWrittenOff = r.condition === 'WRITTEN_OFF';
                        const issuedQty = 0; // issued_quantity не в типі Resource
                        const totalQuantity = r.quantity + issuedQty;
                        
                        let status = 'success'; let statusText = 'OK';
                        if (isWrittenOff) { status = 'neutral'; statusText = 'Списано'; } 
                        else if (r.quantity === 0 && issuedQty === 0) { status = 'critical'; statusText = 'Відсутньо'; } 
                        else if (r.quantity <= r.min_quantity) { status = 'warning'; statusText = 'Нестача'; }
                        
                        const warehouseNameStr = r.warehouse_id ? warehouses.find(w => w.id === r.warehouse_id)?.name || 'Невідомий склад' : 'В дорозі';

                        return (
                          <tr key={r.id} style={{opacity: isWrittenOff ? 0.7 : 1}}>
                            <td style={{fontWeight: 500}}>
                              {r.name}
                              <div style={{fontSize: '0.75rem', color: '#94a3b8', marginTop: '2px', fontWeight: 'normal'}}>
                                ID: {r.id.split('-')[0].toUpperCase()}
                              </div>
                            </td>
                            <td>
                              <div className="location-stack">
                                <span className="stock-info">🏢 Склад: <strong>{r.quantity}</strong> ({warehouseNameStr})</span>
                                {issuedQty > 0 && !isWrittenOff && (<span className="issued-info">👤 На руках: <strong>{issuedQty}</strong></span>)}
                              </div>
                            </td>
                            <td style={{textAlign: 'center', fontWeight: 'bold'}}>
                              {totalQuantity} <small style={{color: '#64748b', fontWeight: 'normal', marginLeft: '4px'}}>{formatUnitType(r.unit_type)}</small>
                            </td>
                            <td style={{textAlign: 'center'}}>{r.min_quantity}</td>
                            <td style={{textAlign: 'center'}}><span className={`badge badge-${status}`}>{statusText}</span></td>
                            
                            {canManageResources && (
                              <td className="col-actions-menu">
                                <div className="dropdown-container" onClick={(e) => e.stopPropagation()}>
                                  <button className={`btn-kebab ${activeMenuId === r.id ? 'active' : ''}`} onClick={() => setActiveMenuId(activeMenuId === r.id ? null : r.id)}>⋮</button>
                                  {activeMenuId === r.id && (
                                    <div className="actions-dropdown-menu">
                                      {!isWrittenOff && (
                                        <>
                                          <button onClick={() => { handleShowQR(r.id, r.name); setActiveMenuId(null); }}>🔍 Переглянути QR-код</button>
                                          <button onClick={() => { handleDownloadQR(r.id, r.name); setActiveMenuId(null); }}>🖨️ Друк наклейки (QR)</button>
                                          {r.quantity > 0 && (
                                            <button style={{color: '#2563eb'}} onClick={() => { setAssignModalData({ resource: r, quantity: 1, user_id: '' }); setActiveMenuId(null); }}>👤 Видати співробітнику</button>
                                          )}
                                          <button onClick={() => { setEditForm({ name: r.name, min_quantity: r.min_quantity }); setEditModalId(r.id); setActiveMenuId(null); }}>✏️ Редагувати</button>
                                          <button onClick={() => { setWriteOffModalData({ resource: r, quantity: r.quantity }); setActiveMenuId(null); }}>📦 Списати зі складу</button>
                                          <div className="dropdown-divider"></div>
                                        </>
                                      )}
                                      <button className="text-danger" onClick={() => { setResourceToDelete(r); setActiveMenuId(null); }}>🗑️ Видалити запис</button>
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