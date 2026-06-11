import React, { useEffect, useState, useRef, useMemo } from 'react';
import { api, getInMemoryToken, type Warehouse, type Unit, type ShipmentRefuel } from '../api/client';
import { usePermissions } from '../hooks/usePermissions';
import { MapContainer, TileLayer, Marker, Popup, Polyline, Tooltip } from 'react-leaflet';
import L from 'leaflet';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';
import Pagination from '../components/Pagination';
import 'leaflet/dist/leaflet.css';
import './Warehouses.css';
import InventoryAuditModal from '../components/InventoryAuditModal';
import SearchableSelect from '../components/SearchableSelect';

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
    const polyline = polyRef.current; 
    if (!polyline) return; 
    if (typeof polyline.getElement !== 'function') return;
    const el = polyline.getElement(); 
    if (!el) return;
    
    el.style.strokeDasharray = pathOptions?.dashArray || '12, 12'; 
    let offset = 0;
    
    const animate = () => { 
      offset -= 0.5; 
      el.style.strokeDashoffset = offset.toString(); 
      requestAnimationFrame(animate); 
    };
    
    const frame = requestAnimationFrame(animate); 
    return () => cancelAnimationFrame(frame);
  }, [pathOptions?.dashArray]);

  return (
    // @ts-ignore
    <Polyline ref={polyRef} positions={positions} pathOptions={pathOptions}>
      {children}
    </Polyline>
  );
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

  const closePopup = () => { markerRef.current?.closePopup(); };

  return (
    <Marker draggable={true} eventHandlers={eventHandlers} position={[warehouse.latitude, warehouse.longitude]} ref={markerRef} icon={icon}>
      <Popup minWidth={270}>
        <div style={{ fontFamily: 'Inter, system-ui, sans-serif', padding: '4px 2px' }}>
          <h4 style={{ margin: '0 0 6px 0', fontSize: '16px', color: 'var(--text-bright)', fontWeight: 700 }}>{warehouse.name}</h4>
          <div style={{ fontSize: '13px', color: 'var(--text-muted)', marginBottom: '14px' }}>🛡️ {unitName}</div>
          
          <div style={{ display: 'flex', gap: '8px', marginBottom: '20px' }}>
            <span className="badge badge-neutral">
              {warehouse.location_type === 'MOBILE' ? '🚛 Мобільний' : '🏢 Стаціонарний'}
            </span>
            <span className="badge badge-success">🟢 На зв'язку</span>
          </div>
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
            <button className="btn btn-primary" style={{ width: '100%', padding: '10px', fontSize: '13px', fontWeight: 600, display: 'flex', justifyContent: 'center', borderRadius: '6px', border: 'none', cursor: 'pointer' }} onClick={() => { closePopup(); onViewInventory(warehouse); }}>📦 Переглянути залишки</button>
            <button className="btn btn-secondary" style={{ width: '100%', padding: '10px', fontSize: '13px', fontWeight: 600, display: 'flex', justifyContent: 'center', borderRadius: '6px', border: '1px solid var(--border)', cursor: 'pointer' }} onClick={() => { closePopup(); onDispatchTrip(warehouse); }}>🚚 Сформувати рейс сюди</button>
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
  const navigate = useNavigate();
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

  const [editingWarehouse, setEditingWarehouse] = useState<Warehouse | null>(null);
  const [warehouseToDelete, setWarehouseToDelete] = useState<Warehouse | null>(null);
  const [editForm, setEditForm] = useState({ name: '', capacity_level: '', zone_type: '' });
  const [isProcessing, setIsProcessing] = useState(false);

  const [auditWarehouse, setAuditWarehouse] = useState<Warehouse | null>(null);
  const [receiveShipmentModal, setReceiveShipmentModal] = useState<{ id: string; distance_km: number } | null>(null);
  const [receiveActualKm, setReceiveActualKm] = useState<string>('');
  const [receiveGPSSource, setReceiveGPSSource] = useState<{
    has_gps: boolean;
    points: number;
    gps_km: number;
    planned_km: number;
    planned_rt_km: number;
    route_status: 'on_route' | 'deviated' | 'unknown';
    deviation_pct: number;
    min_allowed_km: number;
  } | null>(null);

  // ⛽ Дозаправка в дорозі
  const [refuelModal, setRefuelModal] = useState<{ shipmentId: string; vehiclePlate: string } | null>(null);
  const [refuelForm, setRefuelForm] = useState({ liters: '', station_name: '', odometer_km: '', cost_uah: '' });
  const [refuelProcessing, setRefuelProcessing] = useState(false);
  const [shipmentRefuels, setShipmentRefuels] = useState<ShipmentRefuel[]>([]); // для модалі прийому

  const [viewInventoryItems, setViewInventoryItems] = useState<InventoryItem[]>([]);
  const [viewInventoryLoading, setViewInventoryLoading] = useState(false);

  const [warehousesPage, setWarehousesPage] = useState(0);
  const [shipmentsPage, setShipmentsPage] = useState(0);
  const WH_PAGE_SIZE = 10;
  const SHIP_PAGE_SIZE = 10;

  const perms = usePermissions();
  const canManageWarehouses = perms.can('warehouse_manage');
  const canAuditWarehouse = perms.can('warehouse_audit');
  const canUseWarehouseActions = canManageWarehouses || canAuditWarehouse;

  const loadData = async () => {
    try {
      const token = getInMemoryToken();
      const [wRes, uRes, vRes, sRes] = await Promise.all([
        api.warehouses.list().catch(() => []),
        api.units.list().catch(() => []),
        api.vehicles.list().catch(() => []) || [],
        fetch('/api/inventory/shipments', { headers: { 'Authorization': `Bearer ${token}` }, credentials: 'include' }).then(res => res.ok ? res.json() : [])
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
    // Збираємо нащадків (для UPSTREAM рейсів)
    const collectDescendants = (parentId: number, maxDepth = 20, currentDepth = 0) => {
      if (currentDepth >= maxDepth) return;
      const children = units.filter(u => u.parent_id === parentId);
      children.forEach(child => {
        allowedUnitIds.add(child.id);
        collectDescendants(child.id, maxDepth, currentDepth + 1);
      });
    };
    collectDescendants(targetUnit.id);
    return warehouses.filter(w => allowedUnitIds.has(w.unit_id) && w.id !== dispatchTargetWarehouse.id);
  }, [dispatchTargetWarehouse, warehouses, units]);

  const handleViewInventory = async (warehouse: Warehouse) => {
    setViewInventoryWarehouse(warehouse);
    setViewInventoryItems([]);
    setViewInventoryLoading(true);
    try {
      const items = await api.inventory.getByWarehouse(warehouse.id);
      setViewInventoryItems(Array.isArray(items) ? items : []);
    } catch (err) {
      console.error('Помилка завантаження залишків:', err);
    } finally {
      setViewInventoryLoading(false);
    }
  };

  const handleOpenDispatch = (targetWarehouse: Warehouse) => {
    setDispatchTargetWarehouse(targetWarehouse);
    setDispatchParentWarehouse(null); 
    setWarehouseInventory([]);
    setManifest([]); 
    setQtyToAdd(''); 
    setItemToAdd('');
    setActiveRoadRoute(null);
    setVehicles([]); // Скидаємо транспорт - завантажиться після вибору складу-відправника
  };

  const handleSourceChange = async (sourceId: string) => {
    const sourceW = warehouses.find(w => w.id === sourceId);
    setDispatchParentWarehouse(sourceW || null);
    setManifest([]);
    setSelectedVehicleId('');
    if (sourceW && dispatchTargetWarehouse) {
      try {
        const [invRes, vehiclesRes] = await Promise.all([
          api.inventory.getByWarehouse(sourceW.id),
          api.vehicles.getAvailableForRoute(sourceW.id, dispatchTargetWarehouse.id).catch(() => [])
        ]);
        const rawInventory: InventoryItem[] = Array.isArray(invRes) ? invRes : [];
        // Merge duplicate entries with the same name (different DB rows for same resource)
        const mergedMap = new Map<string, InventoryItem>();
        rawInventory.forEach(item => {
          const existing = mergedMap.get(item.name);
          if (existing) {
            mergedMap.set(item.name, {
              ...existing,
              available: (existing.available ?? 0) + (item.available ?? 0),
            });
          } else {
            mergedMap.set(item.name, { ...item });
          }
        });
        const fetchedInventory = Array.from(mergedMap.values());
        setWarehouseInventory(fetchedInventory);
        setItemToAdd(fetchedInventory.length > 0 ? fetchedInventory[0].id : '');
        setVehicles(Array.isArray(vehiclesRes) ? vehiclesRes : []);
      } catch (err) { console.error("Помилка завантаження:", err); }
      await buildRoute(sourceW, dispatchTargetWarehouse);
    }
  };

  const handleCloseDispatch = () => { setDispatchTargetWarehouse(null); setDispatchParentWarehouse(null); setActiveRoadRoute(null); setSelectedVehicleId(''); loadData(); };

  const availableVehicles = vehicles.filter(v => v.type === 'VAN' || v.type === 'TRUCK' || v.type === 'PICKUP');
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

    // Перевірка пального
    const MIN_FUEL_LITERS = 5;
    const vehicleFuel = (selectedVehicle as any).current_fuel_liters ?? 0;
    if (vehicleFuel < MIN_FUEL_LITERS) {
      return toast.error(
        `⛽ Неможливо відправити рейс! Машина "${selectedVehicle.brand} (${selectedVehicle.plate_number})" має лише ${vehicleFuel.toFixed(1)} л пального. Заправте мінімум ${MIN_FUEL_LITERS} л перед рейсом.`,
        { duration: 7000 }
      );
    }

    const toastId = 'dispatch_toast';
    toast.loading('Відправка рейсу в систему...', { id: toastId });

    try {
      const payload = {
        from_warehouse_id: dispatchParentWarehouse.id,
        to_warehouse_id: dispatchTargetWarehouse.id,
        vehicle_id: selectedVehicle.id,
        priority: dispatchPriority,
        items: manifest.map(m => ({ resource_id: m.item.id, quantity: m.quantity })),
        distance_km: parseFloat(activeRoadRoute.distance) || 0,
      };
      
      const token = getInMemoryToken();
      const response = await fetch('/api/inventory/shipments', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        credentials: 'include',
        body: JSON.stringify(payload)
      });
      
      if (!response.ok) {
        const errData = await response.json().catch(() => ({}));
        throw new Error(errData.error || 'Помилка сервера');
      }
      
      toast.success('Рейс успішно сформовано!', { id: toastId, duration: 4000 });
      handleCloseDispatch(); // викликає loadData() всередині
    } catch (error: any) { 
      toast.error(error.message || 'Помилка при створенні рейсу', { id: toastId, duration: 5000 }); 
    }
  };

  const handleReceiveShipment = async (shipmentId: string) => {
    const shipment = shipmentsList.find((s: any) => s.id === shipmentId);
    const distanceKm = parseFloat(shipment?.distance_km) || 0;
    setReceiveShipmentModal({ id: shipmentId, distance_km: distanceKm });
    setReceiveActualKm('');
    setReceiveGPSSource(null);
    setShipmentRefuels([]);

    // Паралельно: GPS-трек + список дозаправок в дорозі
    try {
      const token = getInMemoryToken();
      const [gpsRes, refuelsRes] = await Promise.allSettled([
        fetch(`/api/inventory/shipments/${shipmentId}/gps-distance`, {
          headers: { 'Authorization': `Bearer ${token}` }, credentials: 'include',
        }),
        fetch(`/api/inventory/shipments/${shipmentId}/refuels`, {
          headers: { 'Authorization': `Bearer ${token}` }, credentials: 'include',
        }),
      ]);

      if (gpsRes.status === 'fulfilled' && gpsRes.value.ok) {
        const data = await gpsRes.value.json();
        setReceiveGPSSource(data);
        setReceiveActualKm(data.suggested_km > 0 ? String(data.suggested_km) : (distanceKm > 0 ? String(parseFloat((distanceKm * 2).toFixed(1))) : ''));
      } else {
        setReceiveActualKm(distanceKm > 0 ? String(parseFloat((distanceKm * 2).toFixed(1))) : '');
      }

      if (refuelsRes.status === 'fulfilled' && refuelsRes.value.ok) {
        const data = await refuelsRes.value.json();
        setShipmentRefuels(Array.isArray(data) ? data : []);
      }
    } catch {
      setReceiveActualKm(distanceKm > 0 ? String(parseFloat((distanceKm * 2).toFixed(1))) : '');
    }
  };

  const handleReceiveShipmentConfirm = async () => {
    if (!receiveShipmentModal) return;
    const shipmentId = receiveShipmentModal.id;
    const actualKmValue = parseInt(receiveActualKm, 10);
    if (isNaN(actualKmValue) || actualKmValue < 1) {
      return toast.error('Вкажіть фактичний пробіг (мінімум 1 км)');
    }
    try {
      toast.loading('Приймаємо вантаж на склад...', { id: 'receive' });
      const token = getInMemoryToken();
      const payload: any = {
        actual_km: actualKmValue,
        gps_km: receiveGPSSource?.gps_km ?? 0,
        route_status: receiveGPSSource?.route_status ?? 'unknown',
        deviation_pct: receiveGPSSource?.deviation_pct ?? 0,
      };
      const res = await fetch(`/api/inventory/shipments/${shipmentId}/receive`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const errData = await res.json().catch(() => ({}));
        throw new Error(errData.error || 'Помилка сервера (500)');
      }
      toast.success('Вантаж успішно прийнято! Машину звільнено.', { id: 'receive' });
      setReceiveShipmentModal(null);
      loadData();
    } catch (err: any) {
      toast.error(err.message || 'Не вдалося прийняти вантаж', { id: 'receive' });
    }
  };

  const handleStartShipment = async (shipmentId: string) => {
    try {
      toast.loading('Підтверджуємо відправку...', { id: 'start' });
      const token = getInMemoryToken();
      const res = await fetch(`/api/inventory/shipments/${shipmentId}/start`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` }, credentials: 'include' });
      if (!res.ok) {
        const errData = await res.json().catch(() => ({}));
        throw new Error(errData.error || 'Помилка сервера');
      }
      toast.success('Рейс відправлено! 🚚', { id: 'start' });
      loadData(); 
    } catch (err: any) {
      toast.error(err.message || 'Не вдалося відправити рейс', { id: 'start' });
    }
  };

  // ⛽ Підтвердити дозаправку
  const handleConfirmRefuel = async () => {
    if (!refuelModal) return;
    const liters = parseFloat(refuelForm.liters);
    if (isNaN(liters) || liters <= 0) return toast.error('Вкажіть кількість літрів');
    setRefuelProcessing(true);
    try {
      const payload: any = { liters };
      if (refuelForm.station_name.trim()) payload.station_name = refuelForm.station_name.trim();
      if (refuelForm.odometer_km) payload.odometer_km = parseInt(refuelForm.odometer_km, 10);
      if (refuelForm.cost_uah) payload.cost_uah = parseFloat(refuelForm.cost_uah);
      const result = await api.inventory.logShipmentRefuel(refuelModal.shipmentId, payload);
      toast.success(`✅ Дозаправка ${result.liters} л зареєстрована!`);
      setRefuelModal(null);
    } catch (err: any) {
      toast.error(err.message || 'Помилка реєстрації дозаправки');
    } finally {
      setRefuelProcessing(false);
    }
  };

  const handleDownloadPDF = async (shipmentId: string) => {
    const toastId = toast.loading('Формування накладної...');
    try {
      const { blob, filename } = await api.inventory.downloadShipmentPDF(shipmentId);
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
    } catch (err: any) { 
      // Перевіряємо, чи це помилка ліміту (402 Payment Required)
      if (err?.response?.status === 402 || err?.message?.includes('ліміт') || err?.message?.includes('Ліміт')) {
        const errorMsg = err?.response?.data?.error || err?.message || 'Досягнуто ліміт складів для вашого тарифу';
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
        toast.error(err?.response?.data?.error || err?.message || 'Помилка при створенні складу');
      }
    } 
  };
  
  const handleMarkerDragEnd = (warehouse: Warehouse, lat: number, lng: number) => { setConfirmMove({ warehouse, lat: lat.toFixed(6), lng: lng.toFixed(6) }); };
  
  const saveNewLocation = async () => { 
    if (!confirmMove) return; 
    try { 
      toast.loading('Збереження локації...', { id: 'move' }); 
      await api.warehouses.updateLocation(confirmMove.warehouse.id, parseFloat(confirmMove.lat), parseFloat(confirmMove.lng)); 
      toast.success('Локацію оновлено!', { id: 'move', duration: 3000 }); 
      if (dispatchTargetWarehouse && dispatchParentWarehouse) {
        const updatedTarget = confirmMove.warehouse.id === dispatchTargetWarehouse.id ? { ...dispatchTargetWarehouse, latitude: parseFloat(confirmMove.lat), longitude: parseFloat(confirmMove.lng) } : dispatchTargetWarehouse;
        const updatedParent = confirmMove.warehouse.id === dispatchParentWarehouse.id ? { ...dispatchParentWarehouse, latitude: parseFloat(confirmMove.lat), longitude: parseFloat(confirmMove.lng) } : dispatchParentWarehouse;
        buildRoute(updatedParent, updatedTarget);
      }
      setConfirmMove(null); loadData(); 
    } catch (error) { toast.error('Помилка збереження', { id: 'move' }); } 
  };

  const handleOpenEdit = (w: Warehouse) => {
    setEditingWarehouse(w);
    setEditForm({
      name: w.name,
      capacity_level: w.capacity_level || 'MEDIUM',
      zone_type: w.zone_type || 'REAR'
    });
  };

  const handleEditSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingWarehouse) return;
    setIsProcessing(true);
    try {
      await api.warehouses.update(editingWarehouse.id, editForm as Partial<Warehouse>);
      toast.success('Дані складу оновлено');
      setEditingWarehouse(null);
      loadData();
    } catch (err: any) {
      toast.error(err.message || 'Помилка оновлення');
    } finally {
      setIsProcessing(false);
    }
  };

  const handleDelete = async () => {
    if (!warehouseToDelete) return;
    setIsProcessing(true);
    try {
      await api.warehouses.delete(warehouseToDelete.id);
      toast.success('Склад видалено');
      setWarehouseToDelete(null);
      loadData();
    } catch (err: any) {
      toast.error(err.message || 'Помилка видалення');
    } finally {
      setIsProcessing(false);
    }
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
    <div className="warehouses-page">
      <div className="page-header">
        <h1>Склади та Інфраструктура</h1>
        <div className="page-actions">{canManageWarehouses && <button className="btn btn-primary" onClick={() => setShowForm(true)}>+ Створити склад</button>}</div>
      </div>
      
      <div className="erp-tabs">
        <button className={`tab-btn ${activeTab === 'list' ? 'active' : ''}`} onClick={() => setActiveTab('list')}>📋 Список</button>
        <button className={`tab-btn ${activeTab === 'map' ? 'active' : ''}`} onClick={() => setActiveTab('map')}>🗺️ Карта локацій</button>
        <button className={`tab-btn ${activeTab === 'shipments' ? 'active' : ''}`} onClick={() => setActiveTab('shipments')}>🚚 Рейси та Накладні</button>
      </div>

      {/* Модалка прийому рейсу з введенням фактичного пробігу */}
      {receiveShipmentModal && (() => {
        const gps = receiveGPSSource;
        const actualKmNum = parseFloat(receiveActualKm) || 0;
        const isFraudWarning = gps?.has_gps && gps.gps_km > 0 && actualKmNum > 0 && actualKmNum < gps.gps_km * 0.6;
        const isSubmitDisabled = !receiveActualKm || actualKmNum <= 0 || isFraudWarning;
        return (
          <div className="modal-overlay" onClick={() => { setReceiveShipmentModal(null); setReceiveActualKm(''); setReceiveGPSSource(null); }}>
            <div className="modal" style={{ maxWidth: '460px' }} onClick={e => e.stopPropagation()}>
              <h3 className="modal-title">📦 Прийом вантажу</h3>
              {/* GPS status block */}
              {gps ? (
                gps.has_gps ? (
                  <div style={{ marginBottom: '16px' }}>
                    <div style={{ padding: '12px', background: gps.route_status === 'on_route' ? 'rgba(34,197,94,0.07)' : 'rgba(251,191,36,0.07)', border: `1px solid ${gps.route_status === 'on_route' ? 'rgba(34,197,94,0.3)' : 'rgba(251,191,36,0.3)'}`, borderRadius: '8px', marginBottom: '10px' }}>
                      <div style={{ fontWeight: 600, color: gps.route_status === 'on_route' ? '#22c55e' : '#f59e0b', marginBottom: '6px' }}>
                        {gps.route_status === 'on_route' ? '✅ Маршрут виконано' : '⚠️ Відхилення від маршруту'}
                        {gps.route_status !== 'on_route' && gps.deviation_pct > 0 && <span style={{ fontSize: '12px', marginLeft: '8px', fontWeight: 400 }}>({gps.deviation_pct.toFixed(0)}%)</span>}
                      </div>
                      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px', fontSize: '12px' }}>
                        <div style={{ padding: '8px', background: 'rgba(99,102,241,0.1)', borderRadius: '6px', textAlign: 'center' }}>
                          <div style={{ color: 'var(--text-muted)', marginBottom: '2px' }}>📡 GPS-трек</div>
                          <div style={{ fontWeight: 700, fontSize: '16px' }}>{gps.gps_km} км</div>
                          <div style={{ color: 'var(--text-muted)', fontSize: '10px' }}>{gps.points} точок</div>
                        </div>
                        <div style={{ padding: '8px', background: 'rgba(59,130,246,0.1)', borderRadius: '6px', textAlign: 'center' }}>
                          <div style={{ color: 'var(--text-muted)', marginBottom: '2px' }}>🗺️ OSRM (туди+назад)</div>
                          <div style={{ fontWeight: 700, fontSize: '16px' }}>{gps.planned_rt_km > 0 ? `${gps.planned_rt_km} км` : '—'}</div>
                          <div style={{ color: 'var(--text-muted)', fontSize: '10px' }}>{gps.planned_km > 0 ? `${Math.round(gps.planned_km)} × 2` : 'не розраховувався'}</div>
                        </div>
                      </div>
                    </div>
                  </div>
                ) : receiveShipmentModal.distance_km > 0 ? (
                  <div style={{ padding: '12px', background: 'rgba(59,130,246,0.07)', border: '1px solid rgba(59,130,246,0.25)', borderRadius: '8px', marginBottom: '16px' }}>
                    <div style={{ fontWeight: 600 }}>🗺️ Маршрут з карти (OSRM)</div>
                    <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginTop: '4px' }}>GPS не записувався. Планова відстань: {receiveShipmentModal.distance_km.toFixed(1)} × 2 = <strong>{(receiveShipmentModal.distance_km * 2).toFixed(1)} км</strong>.</div>
                  </div>
                ) : (
                  <div style={{ padding: '12px', background: 'rgba(239,68,68,0.07)', border: '1px solid rgba(239,68,68,0.25)', borderRadius: '8px', marginBottom: '16px' }}>
                    <div style={{ fontWeight: 600, color: '#ef4444' }}>⚠️ Пробіг невідомий</div>
                    <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginTop: '4px' }}>GPS не записувався, маршрут не розраховувався.</div>
                  </div>
                )
              ) : (
                <div style={{ padding: '12px', background: 'rgba(148,163,184,0.07)', border: '1px solid rgba(148,163,184,0.2)', borderRadius: '8px', marginBottom: '16px', fontSize: '13px', color: 'var(--text-muted)' }}>⏳ Завантаження GPS-даних...</div>
              )}
              <div style={{ marginBottom: isFraudWarning ? '8px' : '16px' }}>
                <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--text-muted)', marginBottom: '6px' }}>Фактичний пробіг (км) <span style={{ color: '#ef4444' }}>*</span></label>
                <input type="number" min="0.1" step="0.1" className="erp-input" style={{ width: '100%', boxSizing: 'border-box', borderColor: isFraudWarning ? '#ef4444' : undefined }} value={receiveActualKm} onChange={e => setReceiveActualKm(e.target.value)} autoFocus />
                <span style={{ display: 'block', fontSize: '11px', color: '#64748b', marginTop: '4px' }}>
                  {gps?.has_gps ? `📡 З GPS-треку. Мінімально: ${gps.min_allowed_km} км.` : receiveShipmentModal.distance_km > 0 ? `🗺️ З OSRM (×2 = ${(receiveShipmentModal.distance_km * 2).toFixed(1)} км). Якщо їхали іншим шляхом — введіть реальний пробіг.` : 'Введіть пробіг туди і назад.'}
                </span>
              </div>
              {isFraudWarning && (
                <div style={{ padding: '10px 12px', background: 'rgba(239,68,68,0.1)', border: '1px solid rgba(239,68,68,0.4)', borderRadius: '8px', marginBottom: '16px' }}>
                  <div style={{ fontWeight: 600, color: '#ef4444', fontSize: '13px' }}>🚫 Введений пробіг підозріло малий</div>
                  <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginTop: '4px' }}>GPS показує {gps!.gps_km} км, ви ввели {actualKmNum} км. Мінімум: <strong>{gps!.min_allowed_km} км</strong>.</div>
                </div>
              )}

              {/* ⛽ Дозаправки в дорозі */}
              {shipmentRefuels.length > 0 && (
                <div style={{ marginBottom: '16px', padding: '12px', background: 'rgba(251,191,36,0.08)', border: '1px solid rgba(251,191,36,0.3)', borderRadius: '8px' }}>
                  <div style={{ fontWeight: 600, fontSize: '13px', marginBottom: '8px' }}>⛽ Дозаправки в дорозі</div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                    {shipmentRefuels.map(r => (
                      <div key={r.id} style={{ display: 'flex', justifyContent: 'space-between', fontSize: '12px', padding: '4px 8px', background: 'rgba(255,255,255,0.05)', borderRadius: '4px' }}>
                        <span>{r.station_name || 'Заправка'}</span>
                        <span style={{ fontWeight: 700, color: '#fbbf24' }}>+{r.liters} л{r.cost_uah ? ` · ${r.cost_uah} грн` : ''}</span>
                      </div>
                    ))}
                    <div style={{ textAlign: 'right', fontWeight: 700, fontSize: '12px', marginTop: '4px', color: '#fbbf24' }}>
                      Разом: +{shipmentRefuels.reduce((s, r) => s + r.liters, 0).toFixed(1)} л
                    </div>
                  </div>
                </div>
              )}

              <div className="modal-actions">
                <button className="btn btn-secondary" onClick={() => { setReceiveShipmentModal(null); setReceiveActualKm(''); setReceiveGPSSource(null); setShipmentRefuels([]); }}>Скасувати</button>
                <button className="btn btn-primary" onClick={handleReceiveShipmentConfirm} disabled={isSubmitDisabled}>✅ Підтвердити прийом</button>
              </div>
            </div>
          </div>
        );
      })()}

      {/* ⛽ Модаль реєстрації дозаправки в дорозі */}
      {refuelModal && (
        <div className="modal-overlay" onClick={() => setRefuelModal(null)}>
          <div className="modal" style={{ maxWidth: '420px', width: '100%' }} onClick={e => e.stopPropagation()}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
              <h3 style={{ margin: 0 }}>⛽ Дозаправка в дорозі</h3>
              <button onClick={() => setRefuelModal(null)} style={{ background: 'none', border: 'none', fontSize: '24px', cursor: 'pointer', color: 'var(--text-muted)' }}>&times;</button>
            </div>
            <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginBottom: '16px', padding: '8px 12px', background: 'rgba(251,191,36,0.08)', border: '1px solid rgba(251,191,36,0.25)', borderRadius: '6px' }}>
              🚛 Транспорт: <strong>{refuelModal.vehiclePlate}</strong>
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
              <div>
                <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--text-muted)', marginBottom: '6px' }}>
                  Літри <span style={{ color: '#ef4444' }}>*</span>
                </label>
                <input type="number" min="0.1" step="0.1" className="erp-input" style={{ width: '100%', boxSizing: 'border-box' }}
                  placeholder="напр. 35.5" value={refuelForm.liters}
                  onChange={e => setRefuelForm(f => ({ ...f, liters: e.target.value }))} autoFocus />
              </div>
              <div>
                <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--text-muted)', marginBottom: '6px' }}>
                  Назва АЗС (необов'язково)
                </label>
                <input type="text" className="erp-input" style={{ width: '100%', boxSizing: 'border-box' }}
                  placeholder="напр. ОККО, WOG, ANP..." value={refuelForm.station_name}
                  onChange={e => setRefuelForm(f => ({ ...f, station_name: e.target.value }))} />
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '10px' }}>
                <div>
                  <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--text-muted)', marginBottom: '6px' }}>
                    Одометр (км)
                  </label>
                  <input type="number" min="0" className="erp-input" style={{ width: '100%', boxSizing: 'border-box' }}
                    placeholder="напр. 12450" value={refuelForm.odometer_km}
                    onChange={e => setRefuelForm(f => ({ ...f, odometer_km: e.target.value }))} />
                </div>
                <div>
                  <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--text-muted)', marginBottom: '6px' }}>
                    Вартість (грн)
                  </label>
                  <input type="number" min="0" step="0.01" className="erp-input" style={{ width: '100%', boxSizing: 'border-box' }}
                    placeholder="напр. 2100" value={refuelForm.cost_uah}
                    onChange={e => setRefuelForm(f => ({ ...f, cost_uah: e.target.value }))} />
                </div>
              </div>
            </div>
            <div className="modal-actions" style={{ marginTop: '20px' }}>
              <button className="btn btn-secondary" onClick={() => setRefuelModal(null)}>Скасувати</button>
              <button className="btn btn-warning" onClick={handleConfirmRefuel} disabled={refuelProcessing}>
                {refuelProcessing ? '⏳ Збереження...' : '⛽ Зареєструвати'}
              </button>
            </div>
          </div>
        </div>
      )}

      {viewInventoryWarehouse && (
        <div className="modal-overlay" onClick={() => setViewInventoryWarehouse(null)}>
          <div className="modal" style={{ maxWidth: '600px', width: '100%' }} onClick={e => e.stopPropagation()}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
              <h3 style={{ margin: 0, color: 'var(--text-bright)' }}>📦 Залишки: {viewInventoryWarehouse.name}</h3>
              <button onClick={() => setViewInventoryWarehouse(null)} style={{ background: 'none', border: 'none', fontSize: '24px', cursor: 'pointer', color: 'var(--text-muted)' }}>&times;</button>
            </div>
            <div style={{ padding: '4px 0', color: 'var(--text-muted)', minHeight: '80px' }}>
              {viewInventoryLoading ? (
                <div style={{ textAlign: 'center', padding: '32px' }}>
                  <div className="spinner" style={{ margin: '0 auto' }} />
                  <p style={{ marginTop: '12px' }}>Завантаження залишків...</p>
                </div>
              ) : viewInventoryItems.length === 0 ? (
                <div style={{ textAlign: 'center', padding: '24px', color: 'var(--text-muted)' }}>
                  <p>На цьому складі немає майна.</p>
                </div>
              ) : (
                <table className="data-table" style={{ width: '100%' }}>
                  <thead>
                    <tr>
                      <th>Назва</th>
                      <th>Категорія</th>
                      <th style={{ textAlign: 'right' }}>Кількість</th>
                      <th style={{ textAlign: 'right' }}>Вага, кг</th>
                    </tr>
                  </thead>
                  <tbody>
                    {viewInventoryItems.map(item => (
                      <tr key={item.id}>
                        <td className="font-medium">{item.name}</td>
                        <td style={{ color: 'var(--text-muted)', fontSize: '13px' }}>{(item as any).category || '—'}</td>
                        <td style={{ textAlign: 'right' }}><strong>{item.available ?? (item as any).quantity ?? 0}</strong></td>
                        <td style={{ textAlign: 'right', color: 'var(--text-muted)' }}>{item.weight_kg ?? '—'}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
            <div className="modal-actions">
              <button className="btn btn-primary" onClick={() => { handleOpenDispatch(viewInventoryWarehouse); setViewInventoryWarehouse(null); }}>🚚 Сформувати рейс сюди</button>
            </div>
          </div>
        </div>
      )}

      {showForm && canManageWarehouses && (
        <div className="modal-overlay" onClick={() => setShowForm(false)}>
           <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Новий склад</h3>
            <form onSubmit={handleCreate}>
              <div className="form-group"><label>Орг. одиниця <span className="required">*</span></label><SearchableSelect options={units.map(u => ({ value: String(u.id), label: u.name }))} value={String(newWarehouse.unit_id ?? '')} onChange={(val) => setNewWarehouse({ ...newWarehouse, unit_id: val ? Number(val) : '' })} placeholder="Оберіть одиницю..." /></div>
              <div className="form-group"><label>Назва <span className="required">*</span></label><input className="erp-input" value={newWarehouse.name} onChange={(e) => setNewWarehouse({ ...newWarehouse, name: e.target.value })} required /></div>
              <div className="form-group"><label>Тип <span className="required">*</span></label><select className="erp-input" value={newWarehouse.location_type} onChange={(e) => setNewWarehouse({ ...newWarehouse, location_type: e.target.value as 'STATIONARY' | 'MOBILE' })} required><option value="STATIONARY">Стаціонарний</option><option value="MOBILE">Мобільний</option></select></div>
              <div className="form-row-2">
                <div className="form-group"><label>Широта</label><input className="erp-input" type="number" step="0.000001" value={newWarehouse.latitude} onChange={(e) => setNewWarehouse({ ...newWarehouse, latitude: e.target.value })} /></div>
                <div className="form-group"><label>Довгота</label><input className="erp-input" type="number" step="0.000001" value={newWarehouse.longitude} onChange={(e) => setNewWarehouse({ ...newWarehouse, longitude: e.target.value })} /></div>
              </div>
              <div className="modal-actions"><button type="button" className="btn btn-secondary cancel-margin" onClick={() => setShowForm(false)}>Скасувати</button><button type="submit" className="btn btn-primary" disabled={!newWarehouse.unit_id || !newWarehouse.name?.trim()}>Створити</button></div>
            </form>
          </div>
        </div>
      )}

      {editingWarehouse && (
        <div className="modal-overlay" onClick={() => !isProcessing && setEditingWarehouse(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Редагувати параметри складу</h3>
            <form onSubmit={handleEditSubmit}>
              <div className="form-group">
                <label>Назва складу</label>
                <input className="erp-input" value={editForm.name} onChange={(e) => setEditForm({ ...editForm, name: e.target.value })} required disabled={isProcessing} />
              </div>
              <div className="form-row-2">
                <div className="form-group flex-1">
                  <label>Місткість</label>
                  <select className="erp-input" value={editForm.capacity_level} onChange={(e) => setEditForm({ ...editForm, capacity_level: e.target.value })} disabled={isProcessing}>
                    <option value="LARGE">Великий (Центральний)</option>
                    <option value="MEDIUM">Середній (Регіональний)</option>
                    <option value="SMALL">Малий (Мобільний/Локальний)</option>
                  </select>
                </div>
                <div className="form-group flex-1">
                  <label>Зона розташування</label>
                  <select className="erp-input" value={editForm.zone_type} onChange={(e) => setEditForm({ ...editForm, zone_type: e.target.value })} disabled={isProcessing}>
                    <option value="REAR">Центральний хаб</option>
                    <option value="TACTICAL">Регіональний вузол</option>
                    <option value="FORWARD">Точка видачі (Остання миля)</option>
                  </select>
                </div>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setEditingWarehouse(null)} disabled={isProcessing}>Скасувати</button>
                <button type="submit" className="btn btn-primary" disabled={isProcessing}>{isProcessing ? 'Збереження...' : 'Зберегти'}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {warehouseToDelete && (
        <div className="modal-overlay" onClick={() => !isProcessing && setWarehouseToDelete(null)}>
          <div className="modal confirm-modal" onClick={(e) => e.stopPropagation()}>
            <h3 style={{ color: '#ef4444' }}>⚠️ Видалення складу</h3>
            <p>Ви впевнені, що хочете ліквідувати склад <strong>{warehouseToDelete.name}</strong>?</p>
            <p style={{ fontSize: '12px', color: 'var(--text-muted)' }}>Система дозволить це зробити лише якщо на балансі складу нуль одиниць майна.</p>
            <div className="modal-actions">
              <button className="btn btn-secondary" onClick={() => setWarehouseToDelete(null)} disabled={isProcessing}>Скасувати</button>
              <button className="btn btn-danger" onClick={handleDelete} disabled={isProcessing}>{isProcessing ? 'Видалення...' : 'Ліквідувати'}</button>
            </div>
          </div>
        </div>
      )}

      {confirmMove && (
        <div className="modal-overlay" onClick={() => { setConfirmMove(null); loadData(); }}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h3 className="modal-title-blue">📍 Зміна локації</h3>
            <div className="form-row-2">
              <div className="form-group"><label>Широта (Lat)</label><input type="number" step="0.000001" className="erp-input" value={confirmMove.lat} onChange={(e) => setConfirmMove({...confirmMove, lat: e.target.value})} /></div>
              <div className="form-group"><label>Довгота (Lng)</label><input type="number" step="0.000001" className="erp-input" value={confirmMove.lng} onChange={(e) => setConfirmMove({...confirmMove, lng: e.target.value})} /></div>
            </div>
            <div className="modal-actions"><button className="btn btn-secondary cancel-margin" onClick={() => { setConfirmMove(null); loadData(); }}>Скасувати</button><button className="btn btn-primary" onClick={saveNewLocation}>Зберегти</button></div>
          </div>
        </div>
      )}

      {activeTab === 'list' && (
        <div className="card card-table">
          {warehouses.length === 0 ? <p className="empty-state">Склади ще не створені.</p> : (() => {
            const totalWhPages = Math.max(1, Math.ceil(warehouses.length / WH_PAGE_SIZE));
            const safeWhPage = Math.min(warehousesPage, totalWhPages - 1);
            const pagedWarehouses = warehouses.slice(safeWhPage * WH_PAGE_SIZE, (safeWhPage + 1) * WH_PAGE_SIZE);
            return (
              <>
            <table className="data-table">
              <thead>
                <tr>
                  <th>Назва складу</th>
                  <th>Орг. одиниця</th>
                  <th>Тип</th>
                  <th>Координати</th>
                  {canUseWarehouseActions && <th>Дії</th>}
                </tr>
              </thead>
              <tbody>
                {pagedWarehouses.map((w) => (
                  <tr key={w.id}>
                    <td className="font-bold">{w.name}</td>
                    <td>{units.find(u => u.id === w.unit_id)?.name || 'Невідомо'}</td>
                    <td><span className={`badge ${w.location_type === 'MOBILE' ? 'badge-warning' : 'badge-success'}`}>{w.location_type === 'MOBILE' ? 'Мобільний' : 'Стаціонарний'}</span></td>
                    <td className="text-muted">{w.latitude && w.longitude ? `${w.latitude}, ${w.longitude}` : 'Не вказано'}</td>
                    {canUseWarehouseActions && (
                      <td>
                        <div className="warehouse-action-buttons">
                          {canAuditWarehouse && (
                            <button 
                              className="wh-btn" 
                              style={{ backgroundColor: 'var(--bg-input)', border: '1px solid var(--border)', color: 'var(--text)' }}
                              onClick={() => setAuditWarehouse(w)}
                            >
                              📋 Переоблік
                            </button>
                          )}
                          
                          {canManageWarehouses && (
                            <>
                              <button className="wh-btn wh-edit" onClick={() => handleOpenEdit(w)}>
                                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                                Редагувати
                              </button>
                              <button className="wh-btn wh-delete" onClick={() => setWarehouseToDelete(w)}>
                                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                                Видалити
                              </button>
                            </>
                          )}
                        </div>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
            <Pagination
              currentPage={safeWhPage}
              totalPages={totalWhPages}
              onPageChange={setWarehousesPage}
              totalItems={warehouses.length}
              itemsPerPage={WH_PAGE_SIZE}
            />
              </>
            );
          })()}
        </div>
      )}

      {activeTab === 'shipments' && (
        <div className="card card-table">
          {shipmentsList.length === 0 ? (
            <div className="empty-state"><h3>Немає активних або минулих рейсів</h3></div>
          ) : (() => {
            const totalShipPages = Math.max(1, Math.ceil(shipmentsList.length / SHIP_PAGE_SIZE));
            const safeShipPage = Math.min(shipmentsPage, totalShipPages - 1);
            const pagedShipments = shipmentsList.slice(safeShipPage * SHIP_PAGE_SIZE, (safeShipPage + 1) * SHIP_PAGE_SIZE);
            return (
              <>
            <table className="data-table">
              <thead>
                <tr>
                  <th>Дата відправки</th>
                  <th>Звідки (Відправник)</th>
                  <th>Куди (Одержувач)</th>
                  <th>Транспорт</th>
                  <th className="text-center">Пріоритет</th>
                  <th className="text-center">Статус</th>
                  {canManageWarehouses && <th style={{ textAlign: 'right' }}>Дії</th>}
                </tr>
              </thead>
              <tbody>
                {pagedShipments.map(s => (
                  <tr key={s.id}>
                    <td>{new Date(s.created_at).toLocaleString('uk-UA')}</td>
                    <td className="font-bold">{s.from_warehouse}</td>
                    <td className="font-bold">{s.to_warehouse}</td>
                    <td>{s.vehicle}</td>
                    <td className="text-center">
                      {s.priority === 'URGENT' ? <span className="badge badge-critical">🔴 Терміново</span> : <span className="badge badge-neutral">🟢 Плановий</span>}
                    </td>
                    <td className="text-center">
                      {s.status === 'PENDING' ? <span className="badge badge-neutral">⏳ Очікує</span> :
                       s.status === 'IN_TRANSIT' ? <span className="badge badge-warning">🚛 В дорозі</span> : 
                       <span className="badge badge-success">✅ Доставлено</span>}
                    </td>
                    {canManageWarehouses && (
                      <td className="shipment-actions-cell">
                        <div className="shipment-actions-inner">
                          {s.status === 'PENDING' && (
                            <button className="btn btn-info btn-sm" onClick={() => handleStartShipment(s.id)}>🚀 Відправити</button>
                          )}
                          {s.status === 'IN_TRANSIT' && (
                            <button className="btn btn-primary btn-sm" onClick={() => handleReceiveShipment(s.id)}>📦 Прийняти</button>
                          )}
                          <button className="btn btn-secondary btn-sm" onClick={() => handleDownloadPDF(s.id)} title="Завантажити накладну">📄 Друк ТТН</button>
                        </div>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
            <Pagination
              currentPage={safeShipPage}
              totalPages={totalShipPages}
              onPageChange={setShipmentsPage}
              totalItems={shipmentsList.length}
              itemsPerPage={SHIP_PAGE_SIZE}
            />
              </>
            );
          })()}
        </div>
      )}

      {activeTab === 'map' && (
        <div className={`map-outer-wrapper${dispatchTargetWarehouse ? ' has-sidebar' : ''}`}>
          <div className="map-container-wrapper">
          {warehousesWithCoords.length === 0 ? (
            <div className="empty-state"><h3>Немає об'єктів на карті</h3></div>
          ) : (
            <>
              <MapContainer center={[48.3794, 31.1656]} zoom={6} style={{ height: '100%', width: '100%', zIndex: 0 }}>
                <TileLayer attribution='&copy; OpenStreetMap' url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                
                {activeRoadRoute && activeRoadRoute.positions.length > 0 && !activeRoadRoute.error && (
                  <>
                    {/* @ts-ignore */}
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
                          <div style={{ fontSize: '11px', color: 'var(--text-muted)', fontWeight: 600, textTransform: 'uppercase', marginBottom: '6px' }}>🔗 Підпорядкованість</div>
                          <div style={{ fontSize: '13px', color: 'var(--text-bright)', marginBottom: '2px' }}><strong>Головний:</strong> {line.parentName}</div>
                          <div style={{ fontSize: '13px', color: 'var(--text-bright)', marginBottom: '8px' }}><strong>Підпорядкований:</strong> {line.childName}</div>
                          <div style={{ fontSize: '12px', color: '#2563eb', fontWeight: 600, background: 'rgba(59, 130, 246, 0.12)', padding: '4px 6px', borderRadius: '4px', display: 'inline-block' }}>📏 Пряма відстань: ~{line.distance.toFixed(1)} км</div>
                        </div>
                      </Tooltip>
                    </AnimatedPolyline>
                  );
                })}

                {warehousesWithCoords.map((w) => (
                  <DraggableMarker key={w.id} warehouse={w} icon={w.location_type === 'MOBILE' ? mobileIcon : stationaryIcon} unitName={units.find(u => u.id === w.unit_id)?.name || 'Невідомо'} onDragEnd={handleMarkerDragEnd} onViewInventory={handleViewInventory} onDispatchTrip={handleOpenDispatch} onHover={setHoveredUnitId} />
                ))}
              </MapContainer>
            </>
          )}
          </div>

          {dispatchTargetWarehouse && (
            <div className="dispatch-sidebar">
              <div className="dispatch-header">
                <h3>🚚 Формування рейсу</h3>
                <button onClick={handleCloseDispatch} className="close-btn">&times;</button>
              </div>

              <div className="dispatch-content">
                <p style={{marginBottom: '20px', fontSize: '13px', color: 'var(--text-muted)'}}>Одержувач: <strong style={{color:'var(--text-bright)'}}>{dispatchTargetWarehouse.name}</strong></p>

                <div className="form-group">
                  <label style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text)', marginBottom: '6px', display: 'block' }}>Звідки (Склад відправник)</label>
                  <SearchableSelect
                    options={allowedSourceWarehouses.map(w => ({ value: String(w.id), label: w.name }))}
                    value={String(dispatchParentWarehouse?.id ?? '')}
                    onChange={val => handleSourceChange(val)}
                    placeholder="Оберіть склад вище по ієрархії..."
                    inlineDropdown
                  />
                </div>

                {activeRoadRoute && (
                  <div className={`route-info-box ${activeRoadRoute.error ? 'route-error' : 'route-success'}`}>
                    {activeRoadRoute.isLoading ? '⏳ Прокладання маршруту...' : activeRoadRoute.error ? <span>{activeRoadRoute.error}</span> : <div style={{display: 'flex', justifyContent: 'space-between'}}><span>Відстань: <strong>{activeRoadRoute.distance} км</strong></span><span>Час: <strong>~{activeRoadRoute.duration} хв</strong></span></div>}
                  </div>
                )}

                <div className="form-group">
                  <label style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text)', marginBottom: '6px', display: 'block' }}>Пріоритет</label>
                  <select className="erp-input" value={dispatchPriority} onChange={e => setDispatchPriority(e.target.value)}>
                    <option value="NORMAL">🟢 Звичайний (Плановий)</option>
                    <option value="URGENT">🔴 Терміновий</option>
                  </select>
                </div>

                <div className="form-group">
                  <label style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text)', marginBottom: '6px', display: 'block' }}>Транспорт та Водій</label>
                  <SearchableSelect
                    options={availableVehicles.map(v => ({ value: v.id, label: `${v.brand} (${v.plate_number}) • ${v.capacity_kg} кг` }))}
                    value={selectedVehicleId}
                    onChange={val => setSelectedVehicleId(val)}
                    placeholder={dispatchParentWarehouse ? 'Оберіть вільний ТЗ...' : 'Спочатку оберіть склад-відправника'}
                    disabled={!dispatchParentWarehouse}
                    inlineDropdown
                  />
                  {dispatchParentWarehouse && availableVehicles.length === 0 && (
                    <span style={{ fontSize: '11px', color: 'var(--warning)', marginTop: '4px', display: 'block' }}>⚠️ Немає доступного транспорту для цього маршруту</span>
                  )}
                </div>

                {selectedVehicle && activeRoadRoute && !activeRoadRoute.error && (
                  <div className="fuel-forecast-box">
                    <strong>⛽ Прогноз пального:</strong><br/>
                    {selectedVehicle.fuel_norm ? `${((parseFloat(activeRoadRoute.distance) * 2 / 100) * selectedVehicle.fuel_norm).toFixed(1)} л` : 'Не вказана норма'}
                  </div>
                )}

                <h4 style={{marginTop: '24px', fontSize: '14px', color: 'var(--text-bright)'}}>📦 Вантажний маніфест</h4>
                <div className="manifest-container">
                  <SearchableSelect
                    options={warehouseInventory.map(item => ({ value: item.id, label: item.name + (getRemainingAvailable(item.id) <= 0 ? ' (вичерп.)' : '') }))}
                    value={itemToAdd}
                    onChange={val => setItemToAdd(val)}
                    placeholder="Оберіть майно..."
                    inlineDropdown
                  />
                  <div style={{ display: 'flex', gap: '8px' }}>
                    <input type="number" className="erp-input" placeholder="К-сть" value={qtyToAdd} onChange={e => setQtyToAdd(e.target.value)} style={{flex: 1}}/>
                    <button className="btn btn-secondary" onClick={handleAddToManifest}>+ Додати</button>
                  </div>
                </div>

                <div className="manifest-table-wrapper">
                  <table className="manifest-table">
                    <thead><tr><th>Товар</th><th>К-сть</th><th>Вага</th><th></th></tr></thead>
                    <tbody>
                      {manifest.map(m => (
                        <tr key={m.item.id}>
                          <td>{m.item.name}</td><td>{m.quantity}</td><td>{(getSafeWeight(m.item) * m.quantity).toFixed(1)} кг</td>
                          <td style={{textAlign: 'right'}}><button style={{ color: '#ef4444', background: 'none', border: 'none', cursor: 'pointer', fontSize: '16px' }} onClick={() => handleRemoveFromManifest(m.item.id)}>&times;</button></td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <div className={`weight-summary ${isOverweight ? 'weight-error' : 'weight-ok'}`}>
                  <strong>Вага: {currentTotalWeight.toFixed(1)} кг</strong>
                  <span>Ліміт: {selectedVehicle?.capacity_kg || 0} кг</span>
                </div>
              </div>

              <div className="dispatch-footer">
                <button className="btn btn-secondary" style={{flex: 1}} onClick={handleCloseDispatch}>Скасувати</button>
                <button className="btn btn-primary" style={{flex: 1}} onClick={handleDispatchSubmit} disabled={isOverweight || manifest.length === 0 || !selectedVehicleId}>Відправити 🚀</button>
              </div>
            </div>
          )}
        </div>
      )}

      {auditWarehouse && (
        <InventoryAuditModal 
          warehouseId={auditWarehouse.id} 
          warehouseName={auditWarehouse.name}
          onClose={() => setAuditWarehouse(null)} 
        />
      )}
    </div>
  );
}
