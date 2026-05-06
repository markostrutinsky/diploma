import { useEffect, useRef, useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { getInMemoryToken } from '../api/client';
import toast from 'react-hot-toast';
import './MyShipments.css';

interface Shipment {
  id: string;
  vehicle_id: string;
  vehicle_plate?: string;
  from_warehouse_id: string;
  from_warehouse_name?: string;
  to_warehouse_id: string;
  to_warehouse_name?: string;
  status: 'PENDING' | 'IN_TRANSIT' | 'DELIVERED';
  created_at: string;
  started_at?: string;
  delivered_at?: string;
  items?: Array<{
    resource_name: string;
    quantity: number;
    unit: string;
  }>;
}

// ─── GPS стан ───────────────────────────────────────────────
type GpsStatus = 'idle' | 'active' | 'error' | 'no_permission' | 'unavailable' | 'no_shipment';

export default function MyShipments() {
  const { token } = useAuth();
  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'pending' | 'in_transit' | 'delivered'>('pending');

  // GPS
  const [gpsStatus, setGpsStatus] = useState<GpsStatus>('idle');
  const [lastCoords, setLastCoords] = useState<{ lat: number; lng: number; speed: number | null } | null>(null);
  const watchIdRef = useRef<number | null>(null);
  const pingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const lastPositionRef = useRef<GeolocationPosition | null>(null);

  const hasActiveShipment = shipments.some(s => s.status === 'IN_TRANSIT');

  useEffect(() => {
    loadMyShipments();
  }, [token]);

  // Автоматично запускаємо/зупиняємо GPS залежно від наявності IN_TRANSIT рейсу
  useEffect(() => {
    if (hasActiveShipment) {
      startGpsTracking();
    } else {
      stopGpsTracking();
    }
    return () => stopGpsTracking();
  }, [hasActiveShipment]);

  const loadMyShipments = async () => {
    const authToken = token || getInMemoryToken();
    if (!authToken) return;
    try {
      setLoading(true);
      const res = await fetch('/api/inventory/shipments/my', {
        headers: { 'Authorization': `Bearer ${authToken}` },
        credentials: 'include',
      });
      if (!res.ok) throw new Error('Не вдалося завантажити рейси');
      const data = await res.json();
      setShipments(Array.isArray(data) ? data : []);
    } catch (err: any) {
      toast.error(err.message || 'Помилка завантаження рейсів');
    } finally {
      setLoading(false);
    }
  };

  // ── GPS tracking ─────────────────────────────────────────
  const sendPing = async (pos: GeolocationPosition) => {
    const authToken = getInMemoryToken() || token;
    if (!authToken) return;
    try {
      const res = await fetch('/api/gps/driver/ping', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authToken}` },
        credentials: 'include',
        body: JSON.stringify({
          latitude: pos.coords.latitude,
          longitude: pos.coords.longitude,
          altitude: pos.coords.altitude ?? 0,
          speed: pos.coords.speed ?? 0,
          heading: pos.coords.heading ?? 0,
          accuracy: pos.coords.accuracy,
        }),
      });
      if (res.ok) {
        setGpsStatus('active');
        setLastCoords({
          lat: pos.coords.latitude,
          lng: pos.coords.longitude,
          speed: pos.coords.speed !== null ? Math.round((pos.coords.speed ?? 0) * 3.6) : null, // m/s → km/h
        });
      } else if (res.status === 403) {
        setGpsStatus('no_shipment');
        stopGpsTracking();
      }
    } catch {
      // мовчки ігноруємо мережеву помилку пінгу
    }
  };

  const startGpsTracking = () => {
    if (!navigator.geolocation) {
      setGpsStatus('error');
      return;
    }
    setGpsStatus('idle');
    // Перший пінг одразу — перевіряємо дозвіл
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        lastPositionRef.current = pos;
        sendPing(pos);
        // watchPosition — браузер сам слідкує і кличе callback при русі
        watchIdRef.current = navigator.geolocation.watchPosition(
          (p) => { lastPositionRef.current = p; },
          (e) => setGpsStatus(e.code === 1 ? 'no_permission' : 'unavailable'),
          { enableHighAccuracy: true, timeout: 15000, maximumAge: 5000 }
        );
        // Надсилаємо пінг кожні 10 секунд
        pingIntervalRef.current = setInterval(() => {
          if (lastPositionRef.current) sendPing(lastPositionRef.current);
        }, 10_000);
      },
      (e) => setGpsStatus(e.code === 1 ? 'no_permission' : 'unavailable'),
      { enableHighAccuracy: true, timeout: 15000 }
    );
  };

  const stopGpsTracking = () => {
    if (watchIdRef.current !== null) {
      navigator.geolocation.clearWatch(watchIdRef.current);
      watchIdRef.current = null;
    }
    if (pingIntervalRef.current !== null) {
      clearInterval(pingIntervalRef.current);
      pingIntervalRef.current = null;
    }
    setGpsStatus('idle');
  };
  // ─────────────────────────────────────────────────────────

  const handleStartShipment = async (shipmentId: string) => {
    const authToken = token || getInMemoryToken();
    if (!authToken) return;
    const toastId = toast.loading('Починаємо рейс...');
    try {
      const res = await fetch(`/api/inventory/shipments/${shipmentId}/start`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${authToken}` },
        credentials: 'include',
      });
      if (!res.ok) { const e = await res.json(); throw new Error(e.error || 'Помилка'); }
      toast.success('Рейс розпочато! 🚚', { id: toastId });
      await loadMyShipments();
    } catch (err: any) {
      toast.error(err.message || 'Помилка початку рейсу', { id: toastId });
    }
  };

  const handleCompleteShipment = async (shipmentId: string) => {
    const authToken = token || getInMemoryToken();
    if (!authToken) return;
    const toastId = toast.loading('Завершуємо рейс...');
    try {
      const res = await fetch(`/api/inventory/shipments/${shipmentId}/receive`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${authToken}` },
        credentials: 'include',
      });
      if (!res.ok) { const e = await res.json(); throw new Error(e.error || 'Помилка'); }
      toast.success('Рейс завершено! ✅', { id: toastId });
      await loadMyShipments();
    } catch (err: any) {
      toast.error(err.message || 'Помилка завершення рейсу', { id: toastId });
    }
  };

  const filteredShipments = shipments.filter(s => {
    if (activeTab === 'pending') return s.status === 'PENDING';
    if (activeTab === 'in_transit') return s.status === 'IN_TRANSIT';
    if (activeTab === 'delivered') return s.status === 'DELIVERED';
    return true;
  });

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleString('uk-UA', {
      day: '2-digit', month: '2-digit', year: 'numeric',
      hour: '2-digit', minute: '2-digit'
    });
  };

  const getStatusBadge = (status: string) => {
    const badges: Record<string, { text: string; class: string }> = {
      'PENDING':    { text: '⏳ Очікує',    class: 'status-pending' },
      'IN_TRANSIT': { text: '🚚 В дорозі',  class: 'status-in-transit' },
      'DELIVERED':  { text: '✅ Доставлено', class: 'status-delivered' },
    };
    const badge = badges[status] || { text: status, class: '' };
    return <span className={`status-badge ${badge.class}`}>{badge.text}</span>;
  };

  if (loading) {
    return (
      <div className="my-shipments-page">
        <div className="loading-state"><div className="spinner"></div><p>Завантаження рейсів...</p></div>
      </div>
    );
  }

  return (
    <div className="my-shipments-page">
      <div className="page-header">
        <div className="header-content">
          <h1>🚚 Мої Рейси</h1>
          <p className="subtitle">Призначені вам доставки</p>
        </div>
        <button className="btn-refresh" onClick={loadMyShipments}>🔄 Оновити</button>
      </div>

      {/* ── GPS статус-банер ── */}
      {hasActiveShipment && (
        <div className={`gps-banner gps-banner--${gpsStatus}`}>
          {gpsStatus === 'active' && (
            <>
              <span className="gps-dot gps-dot--active" />
              <span>
                📡 GPS активний — координати надсилаються кожні 10 сек
                {lastCoords && (
                  <span className="gps-coords">
                    &nbsp;| {lastCoords.lat.toFixed(4)}°, {lastCoords.lng.toFixed(4)}°
                    {lastCoords.speed !== null && ` | ${lastCoords.speed} км/год`}
                  </span>
                )}
              </span>
            </>
          )}
          {gpsStatus === 'no_permission' && (
            <span>🔒 Геолокацію заблоковано. Дозвольте доступ у браузері та <button className="gps-retry-btn" onClick={startGpsTracking}>спробуйте знову</button></span>
          )}
          {gpsStatus === 'unavailable' && (
            <span>📡 Геолокація недоступна на цьому пристрої (немає GPS/WiFi-позиціонування). GPS трекінг вимкнено. <button className="gps-retry-btn" onClick={startGpsTracking}>повторити</button></span>
          )}
          {gpsStatus === 'error' && <span>❌ Геолокація не підтримується цим браузером</span>}
          {gpsStatus === 'idle' && <span>⏳ Очікуємо GPS сигнал...</span>}
        </div>
      )}

      <div className="tabs-container">
        <button className={`tab-btn ${activeTab === 'pending' ? 'active' : ''}`} onClick={() => setActiveTab('pending')}>
          ⏳ Очікують ({shipments.filter(s => s.status === 'PENDING').length})
        </button>
        <button className={`tab-btn ${activeTab === 'in_transit' ? 'active' : ''}`} onClick={() => setActiveTab('in_transit')}>
          🚚 В дорозі ({shipments.filter(s => s.status === 'IN_TRANSIT').length})
        </button>
        <button className={`tab-btn ${activeTab === 'delivered' ? 'active' : ''}`} onClick={() => setActiveTab('delivered')}>
          ✅ Доставлено ({shipments.filter(s => s.status === 'DELIVERED').length})
        </button>
      </div>

      <div className="shipments-container">
        {filteredShipments.length === 0 ? (
          <div className="empty-state">
            <div className="empty-icon">📦</div>
            <h3>Рейсів не знайдено</h3>
            <p>
              {activeTab === 'pending' && 'У вас поки немає рейсів, що очікують відправки'}
              {activeTab === 'in_transit' && 'У вас поки немає активних рейсів'}
              {activeTab === 'delivered' && 'У вас поки немає завершених рейсів'}
            </p>
          </div>
        ) : (
          <div className="shipments-grid">
            {filteredShipments.map(shipment => (
              <div key={shipment.id} className={`shipment-card ${shipment.status === 'IN_TRANSIT' ? 'shipment-card--active' : ''}`}>
                <div className="shipment-header">
                  <div className="shipment-id">
                    <strong>Рейс #{shipment.id.substring(0, 8)}</strong>
                    {getStatusBadge(shipment.status)}
                  </div>
                  <div className="vehicle-info">🚗 {shipment.vehicle_plate || shipment.vehicle_id.substring(0, 8)}</div>
                </div>

                <div className="route-info">
                  <div className="route-point"><span className="route-label">📍 Від:</span><strong>{shipment.from_warehouse_name || 'Склад'}</strong></div>
                  <div className="route-arrow">→</div>
                  <div className="route-point"><span className="route-label">📍 До:</span><strong>{shipment.to_warehouse_name || 'Склад'}</strong></div>
                </div>

                {shipment.items && shipment.items.length > 0 && (
                  <div className="items-list">
                    <strong>Вантаж:</strong>
                    <ul>{shipment.items.map((item, idx) => (<li key={idx}>{item.resource_name}: {item.quantity} {item.unit}</li>))}</ul>
                  </div>
                )}

                <div className="shipment-times">
                  <div className="time-item"><span className="time-label">Створено:</span><span className="time-value">{formatDate(shipment.created_at)}</span></div>
                  {shipment.started_at && <div className="time-item"><span className="time-label">Відправлено:</span><span className="time-value">{formatDate(shipment.started_at)}</span></div>}
                  {shipment.delivered_at && <div className="time-item"><span className="time-label">Доставлено:</span><span className="time-value">{formatDate(shipment.delivered_at)}</span></div>}
                </div>

                <div className="shipment-actions">
                  {shipment.status === 'PENDING' && (
                    <button className="btn-start" onClick={() => handleStartShipment(shipment.id)}>🚀 Почати рейс</button>
                  )}
                  {shipment.status === 'IN_TRANSIT' && (
                    <button className="btn-complete" onClick={() => handleCompleteShipment(shipment.id)}>✅ Підтвердити доставку</button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}


interface Shipment {
  id: string;
  vehicle_id: string;
  vehicle_plate?: string;
  from_warehouse_id: string;
  from_warehouse_name?: string;
  to_warehouse_id: string;
  to_warehouse_name?: string;
  status: 'PENDING' | 'IN_TRANSIT' | 'DELIVERED';
  created_at: string;
  started_at?: string;
  delivered_at?: string;
  items?: Array<{
    resource_name: string;
    quantity: number;
    unit: string;
  }>;
}

