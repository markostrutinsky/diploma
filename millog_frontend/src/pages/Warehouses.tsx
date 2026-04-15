import React, { useEffect, useState, useRef, useMemo } from 'react';
import { api, type Warehouse, type Unit } from '../api/client';
import { useAuth } from '../contexts/AuthContext';
import { MapContainer, TileLayer, Marker, Popup, Polyline, Tooltip } from 'react-leaflet';
import L from 'leaflet';
import toast from 'react-hot-toast';
import 'leaflet/dist/leaflet.css';
import './Inventory.css';

const stationaryIcon = L.divIcon({ html: '<div style="font-size: 24px; text-shadow: 0 2px 4px rgba(0,0,0,0.3);">🏢</div>', className: 'custom-map-marker', iconSize: [30, 30], iconAnchor: [15, 30], popupAnchor: [0, -30] });
const mobileIcon = L.divIcon({ html: '<div style="font-size: 24px; text-shadow: 0 2px 4px rgba(0,0,0,0.3);">🚛</div>', className: 'custom-map-marker', iconSize: [30, 30], iconAnchor: [15, 30], popupAnchor: [0, -30] });

type VehicleType = 'PICKUP' | 'VAN' | 'TRUCK';
type VehicleStatus = 'ACTIVE' | 'INACTIVE' | 'IN_REPAIR' | 'ON_MISSION';

export interface Vehicle { id: string; type: VehicleType; brand: string; model?: string; plate_number: string; capacity_kg: number; driver_name?: string; status: VehicleStatus; fuel_norm?: number; tank_capacity?: number; }
export interface InventoryItem { id: string; name: string; available: number; quantity?: number; weight_kg: number; }
interface ManifestItem { item: InventoryItem; quantity: number; }

const AnimatedPolyline = ({ positions, pathOptions, children }: any) => {
  const polyRef = useRef<any>(null);
  useEffect(() => {
    const polyline = polyRef.current; if (!polyline) return; const el = polyline.getElement(); if (!el) return;
    el.style.strokeDasharray = pathOptions.dashArray || '12, 12'; let offset = 0;
    const animate = () => { offset -= 0.5; el.style.strokeDashoffset = offset.toString(); requestAnimationFrame(animate); };
    const frame = requestAnimationFrame(animate); return () => cancelAnimationFrame(frame);
  }, [pathOptions.dashArray]);
  return <Polyline ref={polyRef} positions={positions} pathOptions={pathOptions}>{children}</Polyline>;
};

const calculateDistance = (lat1: number, lon1: number, lat2: number, lon2: number) => {
  const R = 6371; const dLat = (lat2 - lat1) * (Math.PI / 180); const dLon = (lon2 - lon1) * (Math.PI / 180);
  const a = Math.sin(dLat / 2) * Math.sin(dLat / 2) + Math.cos(lat1 * (Math.PI / 180)) * Math.cos(lat2 * (Math.PI / 180)) * Math.sin(dLon / 2) * Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a)); return R * c;
};

const DraggableMarker = ({ warehouse, icon, unitName, onDragEnd, onViewInventory, onDispatchTrip, onHover }: any) => {
  const markerRef = useRef<L.Marker>(null);
  const eventHandlers = useMemo(() => ({
    dragend() { const marker = markerRef.current; if (marker != null) { const { lat, lng } = marker.getLatLng(); onDragEnd(warehouse, lat, lng); } },
    mouseover() { if (onHover) onHover(warehouse.unit_id); }, mouseout() { if (onHover) onHover(null); }
  }), [warehouse, onDragEnd, onHover]);

  return (
    <Marker draggable={true} eventHandlers={eventHandlers} position={[warehouse.latitude, warehouse.longitude]} ref={markerRef} icon={icon}>
      <Popup minWidth={270}>
        <div style={{ fontFamily: 'Inter, system-ui, sans-serif', padding: '4px 2px' }}>
          <h4 style={{ margin: '0 0 6px 0', fontSize: '16px', color: '#0f172a', fontWeight: 700 }}>{warehouse.name}</h4>
          <div style={{ fontSize: '13px', color: '#64748b', marginBottom: '14px' }}>🛡️ {unitName}</div>
          
          <div style={{ display: 'flex', gap: '8px', marginBottom: '20px' }}>
            <span className="badge badge-neutral">
              {warehouse.location_type === 'MOBILE' ? '🚛 Мобільний' : '🏢 Стаціонарний'}
            </span>
            <span className="badge badge-success">🟢 На зв'язку</span>
          </div>
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
            <button className="btn btn-primary" style={{ width: '100%', padding: '10px', fontSize: '13px', fontWeight: 600, display: 'flex', justifyContent: 'center', borderRadius: '6px', border: 'none', cursor: 'pointer' }} onClick={() => onViewInventory(warehouse)}>📦 Переглянути залишки</button>
            <button className="btn btn-secondary" style={{ width: '100%', padding: '10px', fontSize: '13px', fontWeight: 600, display: 'flex', justifyContent: 'center', borderRadius: '6px', border: '1px solid #cbd5e1', cursor: 'pointer' }} onClick={() => onDispatchTrip(warehouse)}>🚚 Сформувати рейс сюди</button>
          </div>
          
          <div style={{ fontSize: '11px', color: '#94a3b8', marginTop: '16px', paddingTop: '12px', textAlign: 'center', borderTop: '1px dashed #cbd5e1' }}>
            🖱️ Потягніть іконку для переміщення
          </div>
        </div>
      </Popup>
    </Marker>
  );
};

export default function Warehouses() {
  const { user } = useAuth();
  
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [units, setUnits] = useState<Unit[]>([]);
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [warehouseInventory, setWarehouseInventory] = useState<InventoryItem[]>([]);
  const [shipmentsList, setShipmentsList] = useState<any[]>([]); 
  
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'list' | 'map' | 'shipments'>('map'); 
  
  const [showForm, setShowForm] = useState(false);
  const [confirmMove, setConfirmMove] = useState<{ warehouse: Warehouse, lat: string, lng: string } | null>(null);
  const [viewInventoryWarehouse, setViewInventoryWarehouse] = useState<Warehouse | null>(null);
  const [dispatchTargetWarehouse, setDispatchTargetWarehouse] = useState<Warehouse | null>(null);
  const [dispatchParentWarehouse, setDispatchParentWarehouse] = useState<Warehouse | null>(null);
  const [hoveredUnitId, setHoveredUnitId] = useState<number | null>(null);
  const [activeRoadRoute, setActiveRoadRoute] = useState<{ positions: [number, number][]; distance: string; duration: number; isLoading: boolean; error: string | null; } | null>(null);

  const [selectedVehicleId, setSelectedVehicleId] = useState<string>('');
  const [manifest, setManifest] = useState<ManifestItem[]>([]);
  const [itemToAdd, setItemToAdd] = useState<string>('');
  const [qtyToAdd, setQtyToAdd] = useState<string>('');
  const [dispatchPriority, setDispatchPriority] = useState('NORMAL');

  const [newWarehouse, setNewWarehouse] = useState({ unit_id: '' as number | '', name: '', location_type: 'STATIONARY' as 'STATIONARY' | 'MOBILE', latitude: '', longitude: '' });

  const canManageWarehouses = ['ADMIN', 'BRIGADE_CMDR', 'BRIGADE_LOGIST', 'BATTALION_CMDR', 'BATTALION_LOGIST', 'COMPANY_CMDR', 'PLATOON_CMDR'].includes(user?.role || '');

  const loadData = async () => {
    try {
      setLoading(true);
      const token = localStorage.getItem('token'); 
      
      const [wRes, uRes, vRes, sRes] = await Promise.all([
        api.warehouses.list().catch(() => []), 
        api.units.list().catch(() => []),
        (api as any).vehicles?.list().catch(() => []) || [],
        fetch('/api/inventory/shipments', { headers: { 'Authorization': `Bearer ${token}` } }).then(res => res.ok ? res.json() : [])
      ]);
      setWarehouses(Array.isArray(wRes) ? wRes : []); 
      setUnits(Array.isArray(uRes) ? uRes : []);
      setVehicles(Array.isArray(vRes) ? vRes : []);
      setShipmentsList(Array.isArray(sRes) ? sRes : []);
    } catch (error) { console.error(error); } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const buildRoute = async (parentW: Warehouse, targetW: Warehouse) => {
    setActiveRoadRoute({ positions: [], distance: '', duration: 0, isLoading: true, error: null });
    if (parentW.latitude && parentW.longitude && targetW.latitude && targetW.longitude) {
      try {
        const url = `https://router.project-osrm.org/route/v1/driving/${parentW.longitude},${parentW.latitude};${targetW.longitude},${targetW.latitude}?overview=full&geometries=geojson`;
        const response = await fetch(url); const data = await response.json();
        if (data.routes && data.routes.length > 0) {
          const route = data.routes[0];
          const roadPositions = route.geometry.coordinates.map((coord: any[]) => [coord[1], coord[0]]);
          setActiveRoadRoute({ positions: roadPositions, distance: (route.distance / 1000).toFixed(1), duration: Math.round(route.duration / 60), isLoading: false, error: null });
        }
      } catch (error) { setActiveRoadRoute({ positions: [], distance: '', duration: 0, isLoading: false, error: '❌ Помилка з\'єднання з сервісом маршрутів.' }); }
    } else { setActiveRoadRoute({ positions: [], distance: '', duration: 0, isLoading: false, error: '❌ Не вказані точні координати на карті.' }); }
  };

  const allowedSourceWarehouses = useMemo(() => {
    if (!dispatchTargetWarehouse || units.length === 0) return [];
    const targetUnit = units.find(u => u.id === dispatchTargetWarehouse.unit_id);
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
    return warehouses.filter(w => allowedUnitIds.has(w.unit_id) && w.id !== dispatchTargetWarehouse.id);
  }, [dispatchTargetWarehouse, warehouses, units]);

  const handleOpenDispatch = (targetWarehouse: Warehouse) => {
    setDispatchTargetWarehouse(targetWarehouse);
    setDispatchParentWarehouse(null); 
    setWarehouseInventory([]);
    setManifest([]); 
    setQtyToAdd(''); 
    setItemToAdd('');
    setActiveRoadRoute(null);
  };

  const handleSourceChange = async (sourceId: string) => {
    const sourceW = warehouses.find(w => w.id === sourceId);
    setDispatchParentWarehouse(sourceW || null);
    setManifest([]); 
    if (sourceW && dispatchTargetWarehouse) {
      try {
        const invRes = await (api as any).inventory?.getByWarehouse(sourceW.id);
        const fetchedInventory = Array.isArray(invRes) ? invRes : [];
        setWarehouseInventory(fetchedInventory);
        if (fetchedInventory.length > 0) setItemToAdd(fetchedInventory[0].id);
      } catch (err) { console.error("Помилка завантаження залишків:", err); }
      await buildRoute(sourceW, dispatchTargetWarehouse);
    }
  };

  const handleCloseDispatch = () => { setDispatchTargetWarehouse(null); setDispatchParentWarehouse(null); setActiveRoadRoute(null); setSelectedVehicleId(''); };

  const availableVehicles = vehicles.filter(v => v.status === 'ACTIVE' && (v.type === 'VAN' || v.type === 'TRUCK'));
  const selectedVehicle = availableVehicles.find(v => v.id === selectedVehicleId);
  
  const getSafeAvailable = (item: any) => item.available ?? item.quantity ?? 0;
  const getSafeWeight = (item: any) => item.weight_kg ?? 1;
  const currentTotalWeight = manifest.reduce((sum, item) => sum + (getSafeWeight(item.item) * item.quantity), 0);
  const isOverweight = selectedVehicle ? currentTotalWeight > selectedVehicle.capacity_kg : false;

  const getRemainingAvailable = (itemId: string) => {
    const invItem = warehouseInventory.find(i => i.id === itemId);
    if (!invItem) return 0;
    const inManifest = manifest.find(m => m.item.id === itemId)?.quantity || 0;
    return getSafeAvailable(invItem) - inManifest;
  };

  const currentAvailableToSelect = getRemainingAvailable(itemToAdd);

  const handleAddToManifest = () => {
    const parsedQty = parseInt(qtyToAdd, 10);
    if (isNaN(parsedQty) || parsedQty <= 0) return toast.error('Введіть коректну кількість');
    if (parsedQty > currentAvailableToSelect) return toast.error(`На складі залишилось лише ${currentAvailableToSelect} шт.`);
    const invItem = warehouseInventory.find(i => i.id === itemToAdd);
    if (!invItem) return;

    setManifest(prev => {
      const existing = prev.find(p => p.item.id === itemToAdd);
      if (existing) { return prev.map(p => p.item.id === itemToAdd ? { ...p, quantity: p.quantity + parsedQty } : p); }
      return [...prev, { item: invItem, quantity: parsedQty }];
    });
    setQtyToAdd('');
  };

  const handleRemoveFromManifest = (itemId: string) => { setManifest(prev => prev.filter(m => m.item.id !== itemId)); };

  const handleDispatchSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!dispatchParentWarehouse || !dispatchTargetWarehouse || !activeRoadRoute || activeRoadRoute.error) return toast.error('Помилка логістики.');
    if (!selectedVehicle) return toast.error('Оберіть транспортний засіб!');
    if (manifest.length === 0) return toast.error('Маніфест порожній! Додайте вантаж.');
    if (isOverweight) return toast.error(`Перевантаження! Максимум: ${selectedVehicle.capacity_kg} кг`);

    const toastId = 'dispatch_toast';
    toast.loading('Відправка рейсу в систему...', { id: toastId });

    try {
      const payload = {
        from_warehouse_id: dispatchParentWarehouse.id,
        to_warehouse_id: dispatchTargetWarehouse.id,
        vehicle_id: selectedVehicle.id,
        priority: dispatchPriority,
        items: manifest.map(m => ({ resource_id: m.item.id, quantity: m.quantity }))
      };
      
      const token = localStorage.getItem('token');
      const response = await fetch('/api/inventory/shipments', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify(payload)
      });
      
      if (!response.ok) {
        const errData = await response.json().catch(() => ({}));
        throw new Error(errData.error || 'Помилка сервера');
      }
      
      toast.success('Рейс успішно сформовано!', { id: toastId, duration: 4000 });
      handleCloseDispatch();
      loadData(); 
    } catch (error: any) { 
      toast.error(error.message || 'Помилка при створенні рейсу', { id: toastId, duration: 5000 }); 
    }
  };

  const handleReceiveShipment = async (shipmentId: string) => {
    try {
      toast.loading('Приймаємо вантаж на склад...', { id: 'receive' });
      const token = localStorage.getItem('token');
      const res = await fetch(`/api/inventory/shipments/${shipmentId}/receive`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } });
      if (!res.ok) {
        const errData = await res.json().catch(() => ({}));
        throw new Error(errData.error || 'Помилка сервера (500)');
      }
      toast.success('Вантаж успішно прийнято! Машину звільнено.', { id: 'receive' });
      loadData(); 
    } catch (err: any) {
      toast.error(err.message || 'Не вдалося прийняти вантаж', { id: 'receive' });
    }
  };

  const handleDownloadPDF = async (shipmentId: string) => {
    const toastId = toast.loading('Формування накладної...');
    try {
      const { blob, filename } = await (api as any).inventory.downloadShipmentPDF(shipmentId);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', filename);
      document.body.appendChild(link);
      link.click();
      link.parentNode?.removeChild(link);
      window.URL.revokeObjectURL(url);
      toast.success('Накладну завантажено!', { id: toastId });
    } catch (error) {
      toast.error('Не вдалося завантажити документ', { id: toastId });
    }
  };

  const handleCreate = async (e: React.FormEvent) => { 
    e.preventDefault(); 
    if (!newWarehouse.unit_id) return; 
    try { 
      await api.warehouses.create({ unit_id: Number(newWarehouse.unit_id), name: newWarehouse.name, location_type: newWarehouse.location_type, latitude: newWarehouse.latitude ? parseFloat(newWarehouse.latitude) : undefined, longitude: newWarehouse.longitude ? parseFloat(newWarehouse.longitude) : undefined, }); 
      setShowForm(false); setNewWarehouse({ unit_id: '', name: '', location_type: 'STATIONARY', latitude: '', longitude: '' }); 
      loadData(); toast.success('Склад успішно створено'); 
    } catch (err) { toast.error('Помилка при створенні складу'); } 
  };
  
  const handleMarkerDragEnd = (warehouse: Warehouse, lat: number, lng: number) => { setConfirmMove({ warehouse, lat: lat.toFixed(6), lng: lng.toFixed(6) }); };
  
  const saveNewLocation = async () => { 
    if (!confirmMove) return; 
    try { 
      toast.loading('Збереження дислокації...', { id: 'move' }); 
      await api.warehouses.updateLocation(confirmMove.warehouse.id, parseFloat(confirmMove.lat), parseFloat(confirmMove.lng)); 
      toast.success('Дислокацію оновлено!', { id: 'move', duration: 3000 }); 
      if (dispatchTargetWarehouse && dispatchParentWarehouse) {
        const updatedTarget = confirmMove.warehouse.id === dispatchTargetWarehouse.id ? { ...dispatchTargetWarehouse, latitude: parseFloat(confirmMove.lat), longitude: parseFloat(confirmMove.lng) } : dispatchTargetWarehouse;
        const updatedParent = confirmMove.warehouse.id === dispatchParentWarehouse.id ? { ...dispatchParentWarehouse, latitude: parseFloat(confirmMove.lat), longitude: parseFloat(confirmMove.lng) } : dispatchParentWarehouse;
        buildRoute(updatedParent, updatedTarget);
      }
      setConfirmMove(null); loadData(); 
    } catch (error) { toast.error('Помилка збереження', { id: 'move' }); } 
  };
  
  const warehousesWithCoords = warehouses.filter(w => w.latitude && w.longitude);
  
  const hierarchyLines = useMemo(() => {
    const lines: any[] = [];
    warehousesWithCoords.forEach(childW => {
      const childUnit = units.find(u => u.id === childW.unit_id);
      if (childUnit && childUnit.parent_id) {
        const parentWarehouses = warehousesWithCoords.filter(w => w.unit_id === childUnit.parent_id);
        parentWarehouses.forEach(parentW => {
          if (parentW.latitude && parentW.longitude) {
            lines.push({
              id: `line-${parentW.id}-${childW.id}`,
              positions: [[parentW.latitude, parentW.longitude], [childW.latitude as number, childW.longitude as number]],
              parentName: parentW.name, childName: childW.name,
              distance: calculateDistance(parentW.latitude, parentW.longitude, childW.latitude as number, childW.longitude as number),
              parentId: parentW.unit_id, childId: childW.unit_id
            });
          }
        });
      }
    }); 
    return lines;
  }, [warehousesWithCoords, units]);

  if (loading) return <div className="page-loading"><div className="spinner" /></div>;

  return (
    <div className="inventory-page">
      <div className="page-header">
        <h1>Склади та Інфраструктура</h1>
        <div className="page-actions">{canManageWarehouses && <button className="btn btn-primary" onClick={() => setShowForm(true)}>+ Створити склад</button>}</div>
      </div>
      
      <div className="erp-tabs">
        <button className={`tab-btn ${activeTab === 'list' ? 'active' : ''}`} onClick={() => setActiveTab('list')}>📋 Список</button>
        <button className={`tab-btn ${activeTab === 'map' ? 'active' : ''}`} onClick={() => setActiveTab('map')}>🗺️ Карта дислокації</button>
        <button className={`tab-btn ${activeTab === 'shipments' ? 'active' : ''}`} onClick={() => setActiveTab('shipments')}>🚚 Рейси та Накладні</button>
      </div>

      {/* 🔥 ПОВЕРНУВ ЗНИКЛЕ ВІКНО ПЕРЕГЛЯДУ ЗАЛИШКІВ */}
      {viewInventoryWarehouse && (
        <div className="modal-overlay" onClick={() => setViewInventoryWarehouse(null)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
              <h3 style={{ margin: 0, color: '#1e293b' }}>📦 Залишки: {viewInventoryWarehouse.name}</h3>
              <button onClick={() => setViewInventoryWarehouse(null)} style={{ background: 'none', border: 'none', fontSize: '24px', cursor: 'pointer', color: '#64748b' }}>&times;</button>
            </div>
            <div style={{ padding: '20px 0', color: '#64748b' }}>
              <p>Функціонал перегляду в розробці...</p>
            </div>
            <div className="modal-actions">
              <button className="btn btn-primary" onClick={() => { handleOpenDispatch(viewInventoryWarehouse); setViewInventoryWarehouse(null); }}>🚚 Сформувати рейс сюди</button>
            </div>
          </div>
        </div>
      )}

      {/* МОДАЛКА: СТВОРЕННЯ СКЛАДУ */}
      {showForm && canManageWarehouses && (
        <div className="modal-overlay" onClick={() => setShowForm(false)}>
           <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Новий склад</h3>
            <form onSubmit={handleCreate}>
              <div className="form-group"><label>Підрозділ <span className="required">*</span></label><select className="erp-input" value={newWarehouse.unit_id} onChange={(e) => setNewWarehouse({ ...newWarehouse, unit_id: e.target.value ? Number(e.target.value) : '' })} required><option value="" disabled>Оберіть підрозділ</option>{units.map((u) => <option key={u.id} value={u.id}>{u.name}</option>)}</select></div>
              <div className="form-group"><label>Назва <span className="required">*</span></label><input className="erp-input" value={newWarehouse.name} onChange={(e) => setNewWarehouse({ ...newWarehouse, name: e.target.value })} required /></div>
              <div className="form-group"><label>Тип <span className="required">*</span></label><select className="erp-input" value={newWarehouse.location_type} onChange={(e) => setNewWarehouse({ ...newWarehouse, location_type: e.target.value as 'STATIONARY' | 'MOBILE' })} required><option value="STATIONARY">Стаціонарний</option><option value="MOBILE">Мобільний</option></select></div>
              <div className="form-row-2">
                <div className="form-group"><label>Lat</label><input className="erp-input" type="number" step="0.000001" value={newWarehouse.latitude} onChange={(e) => setNewWarehouse({ ...newWarehouse, latitude: e.target.value })} /></div>
                <div className="form-group"><label>Lng</label><input className="erp-input" type="number" step="0.000001" value={newWarehouse.longitude} onChange={(e) => setNewWarehouse({ ...newWarehouse, longitude: e.target.value })} /></div>
              </div>
              <div className="modal-actions"><button type="button" className="btn btn-secondary cancel-margin" onClick={() => setShowForm(false)}>Скасувати</button><button type="submit" className="btn btn-primary" disabled={!newWarehouse.unit_id}>Створити</button></div>
            </form>
          </div>
        </div>
      )}

      {/* МОДАЛКА: ПІДТВЕРДЖЕННЯ ЗМІНИ ЛОКАЦІЇ НА КАРТІ */}
      {confirmMove && (
        <div className="modal-overlay" onClick={() => { setConfirmMove(null); loadData(); }}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h3 className="modal-title-blue">📍 Зміна дислокації</h3>
            <div className="form-row-2">
              <div className="form-group"><label>Широта (Lat)</label><input type="number" step="0.000001" className="erp-input" value={confirmMove.lat} onChange={(e) => setConfirmMove({...confirmMove, lat: e.target.value})} /></div>
              <div className="form-group"><label>Довгота (Lng)</label><input type="number" step="0.000001" className="erp-input" value={confirmMove.lng} onChange={(e) => setConfirmMove({...confirmMove, lng: e.target.value})} /></div>
            </div>
            <div className="modal-actions"><button className="btn btn-secondary cancel-margin" onClick={() => { setConfirmMove(null); loadData(); }}>Скасувати</button><button className="btn btn-primary" onClick={saveNewLocation}>Зберегти</button></div>
          </div>
        </div>
      )}

      {/* СПИСОК СКЛАДІВ */}
      {activeTab === 'list' && (
        <div className="card card-table">
          {warehouses.length === 0 ? <p className="empty-state">Склади ще не створені.</p> : (
            <table className="data-table">
              <thead><tr><th>Назва складу</th><th>Підрозділ</th><th>Тип</th><th>Координати</th></tr></thead>
              <tbody>
                {warehouses.map((w) => (
                  <tr key={w.id}>
                    <td style={{ fontWeight: 'bold' }}>{w.name}</td>
                    <td>{units.find(u => u.id === w.unit_id)?.name || 'Невідомо'}</td>
                    <td><span className={`badge ${w.location_type === 'MOBILE' ? 'badge-warning' : 'badge-success'}`}>{w.location_type === 'MOBILE' ? 'Мобільний' : 'Стаціонарний'}</span></td>
                    <td className="text-muted">{w.latitude && w.longitude ? `${w.latitude}, ${w.longitude}` : 'Не вказано'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}

      {/* ТАБЛИЦЯ РЕЙСІВ ТА ТТН */}
      {activeTab === 'shipments' && (
        <div className="card card-table">
          {shipmentsList.length === 0 ? (
            <div className="empty-state"><h3>Немає активних або минулих рейсів</h3></div>
          ) : (
            <table className="data-table">
              <thead>
                <tr>
                  <th>Дата відправки</th>
                  <th>Звідки (Відправник)</th>
                  <th>Куди (Одержувач)</th>
                  <th>Транспорт</th>
                  <th className="text-center">Пріоритет</th>
                  <th className="text-center">Статус</th>
                  {canManageWarehouses && <th className="actions-col">Дії</th>}
                </tr>
              </thead>
              <tbody>
                {shipmentsList.map(s => (
                  <tr key={s.id}>
                    <td>{new Date(s.created_at).toLocaleString('uk-UA')}</td>
                    <td className="font-bold">{s.from_warehouse}</td>
                    <td className="font-bold">{s.to_warehouse}</td>
                    <td>{s.vehicle}</td>
                    <td className="text-center">
                      {s.priority === 'URGENT' ? <span className="badge badge-critical">🔴 Терміново</span> : <span className="badge badge-neutral">🟢 Плановий</span>}
                    </td>
                    <td className="text-center">
                      {s.status === 'DISPATCHED' ? <span className="badge badge-warning">🚛 В дорозі</span> : <span className="badge badge-success">✅ Доставлено</span>}
                    </td>
                    {canManageWarehouses && (
                      <td className="actions-col">
                        <div className="actions-flex">
                          <button className="btn btn-secondary btn-sm" onClick={() => handleDownloadPDF(s.id)} title="Завантажити накладну">📄 Друк ТТН</button>
                          {s.status === 'DISPATCHED' && (
                            <button className="btn btn-primary btn-sm" onClick={() => handleReceiveShipment(s.id)}>📦 Прийняти</button>
                          )}
                        </div>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}

      {/* КАРТА ДИСЛОКАЦІЇ */}
      {activeTab === 'map' && (
        <div style={{ height: '700px', width: '100%', borderRadius: '12px', border: '1px solid #e2e8f0', overflow: 'hidden', backgroundColor: '#f8fafc', position: 'relative' }}>
          {warehousesWithCoords.length === 0 ? (
            <div className="empty-state"><h3>Немає об'єктів на карті</h3></div>
          ) : (
            <>
              <MapContainer center={[48.3794, 31.1656]} zoom={6} style={{ height: '100%', width: '100%', zIndex: 0 }}>
                <TileLayer attribution='&copy; OpenStreetMap' url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                
                {activeRoadRoute && activeRoadRoute.positions.length > 0 && !activeRoadRoute.error && (
                  <>
                    <Polyline positions={activeRoadRoute.positions} pathOptions={{ color: '#c4b5fd', weight: 8, opacity: 0.8 }} />
                    <AnimatedPolyline positions={activeRoadRoute.positions} pathOptions={{ color: '#7c3aed', weight: 8, dashArray: '15, 25' }} />
                  </>
                )}

                {hierarchyLines.map((line) => {
                  const isHoveringAny = hoveredUnitId !== null;
                  const isConnectedToHovered = isHoveringAny && (line.parentId === hoveredUnitId || line.childId === hoveredUnitId);
                  const lineColor = isHoveringAny ? (isConnectedToHovered ? '#f59e0b' : '#cbd5e1') : '#3b82f6';
                  const lineOpacity = isHoveringAny ? (isConnectedToHovered ? 1 : 0.3) : 0.6;
                  return (
                    <AnimatedPolyline key={`${line.id}-${lineColor}`} positions={line.positions} pathOptions={{ color: lineColor, weight: 4, opacity: lineOpacity }}>
                      <Tooltip sticky>
                        <div style={{ padding: '4px', fontFamily: 'Inter, sans-serif' }}>
                          <div style={{ fontSize: '11px', color: '#64748b', fontWeight: 600, textTransform: 'uppercase', marginBottom: '6px' }}>🔗 Підпорядкованість</div>
                          <div style={{ fontSize: '13px', color: '#0f172a', marginBottom: '2px' }}><strong>Головний:</strong> {line.parentName}</div>
                          <div style={{ fontSize: '13px', color: '#0f172a', marginBottom: '8px' }}><strong>Підпорядкований:</strong> {line.childName}</div>
                          <div style={{ fontSize: '12px', color: '#2563eb', fontWeight: 600, background: '#eff6ff', padding: '4px 6px', borderRadius: '4px', display: 'inline-block' }}>📏 Пряма відстань: ~{line.distance.toFixed(1)} км</div>
                        </div>
                      </Tooltip>
                    </AnimatedPolyline>
                  );
                })}

                {warehousesWithCoords.map((w) => (
                  <DraggableMarker key={w.id} warehouse={w} icon={w.location_type === 'MOBILE' ? mobileIcon : stationaryIcon} unitName={units.find(u => u.id === w.unit_id)?.name || 'Невідомо'} onDragEnd={handleMarkerDragEnd} onViewInventory={setViewInventoryWarehouse} onDispatchTrip={handleOpenDispatch} onHover={setHoveredUnitId} />
                ))}
              </MapContainer>

              {/* БОКОВА ПАНЕЛЬ: ФОРМУВАННЯ РЕЙСУ */}
              {dispatchTargetWarehouse && (
                <div style={{ position: 'absolute', top: '16px', right: '16px', bottom: '16px', width: '420px', backgroundColor: '#ffffff', borderRadius: '12px', boxShadow: '0 10px 30px rgba(0,0,0,0.15)', zIndex: 1000, display: 'flex', flexDirection: 'column', overflow: 'hidden', border: '1px solid #e2e8f0', fontFamily: 'Inter, sans-serif' }}>
                  <div style={{ padding: '16px 20px', borderBottom: '1px solid #e2e8f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center', backgroundColor: '#f8fafc' }}>
                    <h3 style={{ margin: 0, fontSize: '16px', color: '#1e293b', fontWeight: 600 }}>🚚 Формування рейсу</h3>
                    <button onClick={handleCloseDispatch} style={{ background: 'none', border: 'none', fontSize: '28px', lineHeight: '1', cursor: 'pointer', color: '#94a3b8', padding: 0, marginTop: '-4px' }}>&times;</button>
                  </div>

                  <div style={{ padding: '20px', overflowY: 'auto', flex: 1 }}>
                    <p style={{marginBottom: '20px', fontSize: '13px', color: '#64748b'}}>Одержувач: <strong style={{color:'#0f172a'}}>{dispatchTargetWarehouse.name}</strong></p>

                    <div className="form-group">
                      <label style={{ fontSize: '12px', fontWeight: 600, color: '#475569', marginBottom: '6px', display: 'block' }}>Звідки (Склад відправник)</label>
                      <select className="erp-input" value={dispatchParentWarehouse?.id || ''} onChange={e => handleSourceChange(e.target.value)}>
                        <option value="" disabled>Оберіть склад вище по ієрархії...</option>
                        {allowedSourceWarehouses.map(w => <option key={w.id} value={w.id}>{w.name}</option>)}
                      </select>
                    </div>

                    {activeRoadRoute && (
                      <div style={{ backgroundColor: activeRoadRoute.error ? '#fef2f2' : '#f3e8ff', padding: '12px 16px', borderRadius: '8px', marginBottom: '20px', borderLeft: `4px solid ${activeRoadRoute.error ? '#ef4444' : '#8b5cf6'}`, fontSize: '13px' }}>
                        {activeRoadRoute.isLoading ? '⏳ Прокладання маршруту...' : activeRoadRoute.error ? <span style={{color: '#ef4444'}}>{activeRoadRoute.error}</span> : <div style={{display: 'flex', justifyContent: 'space-between'}}><span>Відстань: <strong>{activeRoadRoute.distance} км</strong></span><span>Час: <strong>~{activeRoadRoute.duration} хв</strong></span></div>}
                      </div>
                    )}

                    <div className="form-group">
                      <label style={{ fontSize: '12px', fontWeight: 600, color: '#475569', marginBottom: '6px', display: 'block' }}>Пріоритет</label>
                      <select className="erp-input" value={dispatchPriority} onChange={e => setDispatchPriority(e.target.value)}>
                        <option value="NORMAL">🟢 Звичайний (Плановий)</option>
                        <option value="URGENT">🔴 Терміновий</option>
                      </select>
                    </div>

                    <div className="form-group">
                      <label style={{ fontSize: '12px', fontWeight: 600, color: '#475569', marginBottom: '6px', display: 'block' }}>Транспорт та Водій</label>
                      <select className="erp-input" value={selectedVehicleId} onChange={e => setSelectedVehicleId(e.target.value)}>
                        <option value="" disabled>Оберіть вільний ТЗ...</option>
                        {availableVehicles.map(v => <option key={v.id} value={v.id}>{v.brand} ({v.plate_number}) • {v.capacity_kg} кг</option>)}
                      </select>
                    </div>

                    {selectedVehicle && activeRoadRoute && !activeRoadRoute.error && (
                      <div style={{ marginTop: '12px', padding: '12px', backgroundColor: '#fffbeb', border: '1px solid #fde68a', borderRadius: '8px', fontSize: '13px', color: '#78350f' }}>
                        <strong>⛽ Прогноз пального:</strong><br/>
                        {selectedVehicle.fuel_norm ? `${((parseFloat(activeRoadRoute.distance) * 2 / 100) * selectedVehicle.fuel_norm).toFixed(1)} л` : 'Не вказана норма'}
                      </div>
                    )}

                    <h4 style={{marginTop: '24px', fontSize: '14px', color: '#1e293b'}}>📦 Вантажний маніфест</h4>
                    <div style={{ backgroundColor: '#f8fafc', padding: '12px', borderRadius: '8px', border: '1px solid #e2e8f0', marginBottom: '16px' }}>
                      <select className="erp-input" style={{ marginBottom: '8px', backgroundColor: '#fff' }} value={itemToAdd} onChange={e => setItemToAdd(e.target.value)}>
                        <option value="" disabled>Оберіть майно...</option>
                        {warehouseInventory.map(item => <option key={item.id} value={item.id} disabled={getRemainingAvailable(item.id) <= 0}>{item.name}</option>)}
                      </select>
                      <div style={{ display: 'flex', gap: '8px' }}>
                        <input type="number" className="erp-input" placeholder="К-сть" value={qtyToAdd} onChange={e => setQtyToAdd(e.target.value)} style={{flex: 1, backgroundColor: '#fff'}}/>
                        <button className="btn btn-secondary" onClick={handleAddToManifest}>+ Додати</button>
                      </div>
                    </div>

                    <div style={{ border: '1px solid #e2e8f0', borderRadius: '8px', overflow: 'hidden' }}>
                      <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '12px', textAlign: 'left' }}>
                        <thead style={{ backgroundColor: '#f1f5f9', borderBottom: '1px solid #e2e8f0' }}><tr><th style={{ padding: '8px 12px' }}>Товар</th><th style={{ padding: '8px 12px' }}>К-сть</th><th style={{ padding: '8px 12px' }}>Вага</th><th></th></tr></thead>
                        <tbody>
                          {manifest.map(m => (
                            <tr key={m.item.id} style={{ borderBottom: '1px solid #e2e8f0' }}>
                              <td style={{ padding: '8px 12px' }}>{m.item.name}</td><td style={{ padding: '8px 12px' }}>{m.quantity}</td><td style={{ padding: '8px 12px' }}>{(getSafeWeight(m.item) * m.quantity).toFixed(1)} кг</td>
                              <td style={{textAlign: 'right'}}><button style={{ color: '#ef4444', background: 'none', border: 'none', cursor: 'pointer', fontSize: '16px' }} onClick={() => handleRemoveFromManifest(m.item.id)}>&times;</button></td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>

                    <div style={{ marginTop: '16px', padding: '12px 16px', borderRadius: '8px', backgroundColor: isOverweight ? '#fef2f2' : '#f0fdf4', border: `1px solid ${isOverweight ? '#fca5a5' : '#bbf7d0'}`, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <strong style={{ color: isOverweight ? '#dc2626' : '#16a34a' }}>Вага: {currentTotalWeight.toFixed(1)} кг</strong>
                      <span style={{ fontSize: '12px', color: isOverweight ? '#ef4444' : '#64748b' }}>Ліміт: {selectedVehicle?.capacity_kg || 0} кг</span>
                    </div>
                  </div>

                  <div style={{ padding: '16px 20px', borderTop: '1px solid #e2e8f0', backgroundColor: '#f8fafc', display: 'flex', gap: '12px' }}>
                    <button className="btn btn-secondary" style={{flex: 1}} onClick={handleCloseDispatch}>Скасувати</button>
                    <button className="btn btn-primary" style={{flex: 1}} onClick={handleDispatchSubmit} disabled={isOverweight || manifest.length === 0}>Відправити 🚀</button>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      )}
    </div>
  );
}