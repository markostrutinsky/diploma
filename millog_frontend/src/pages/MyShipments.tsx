import { useEffect, useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { useGPS } from '../contexts/GPSContext';
import { getInMemoryToken } from '../api/client';
import toast from 'react-hot-toast';
import Pagination from '../components/Pagination';
import ModalPortal from '../components/ModalPortal';
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
  distance_km?: number;
  items?: Array<{ resource_name: string; quantity: number; unit: string }>;
}

export default function MyShipments() {
  const { token } = useAuth();
  const { gpsStatus, lastCoords, hasActiveShipment, refreshActiveShipment } = useGPS();

  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'pending' | 'in_transit' | 'delivered'>('pending');
  const [shipmentsPage, setShipmentsPage] = useState(0);
  const SHIPMENTS_PAGE_SIZE = 20;

  // Модалка підтвердження доставки з фактичним пробігом
  const [completeModal, setCompleteModal] = useState<{ shipmentId: string; plannedKm: number } | null>(null);
  const [actualKm, setActualKm] = useState('');
  const [isCompleting, setIsCompleting] = useState(false);

  useEffect(() => { setShipmentsPage(0); }, [activeTab]);

  useEffect(() => {
    loadMyShipments();
  }, [token]);

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
      // Негайно оновлюємо GPS-контекст — не чекаємо 30-секундного polling
      await refreshActiveShipment();
    } catch (err: any) {
      toast.error(err.message || 'Помилка початку рейсу', { id: toastId });
    }
  };

  const handleCompleteShipment = async (shipmentId: string, plannedKm: number) => {
    setCompleteModal({ shipmentId, plannedKm });
    setActualKm(plannedKm > 0 ? String(Math.round(plannedKm)) : '');
  };

  const handleConfirmComplete = async () => {
    if (!completeModal) return;
    const authToken = token || getInMemoryToken();
    if (!authToken) return;
    setIsCompleting(true);
    const toastId = toast.loading('Завершуємо рейс...');
    try {
      const km = parseInt(actualKm, 10) || 0;
      const res = await fetch(`/api/inventory/shipments/${completeModal.shipmentId}/receive`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${authToken}`, 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ actual_km: km }),
      });
      if (!res.ok) { const e = await res.json(); throw new Error(e.error || 'Помилка'); }
      toast.success(`Рейс завершено! Одометр оновлено на +${km} км ✅`, { id: toastId });
      setCompleteModal(null);
      await loadMyShipments();
    } catch (err: any) {
      toast.error(err.message || 'Помилка завершення рейсу', { id: toastId });
    } finally {
      setIsCompleting(false);
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
      hour: '2-digit', minute: '2-digit',
    });
  };

  const getStatusBadge = (status: string) => {
    const badges: Record<string, { text: string; class: string }> = {
      PENDING:    { text: '⏳ Очікує',    class: 'status-pending' },
      IN_TRANSIT: { text: '🚚 В дорозі',  class: 'status-in-transit' },
      DELIVERED:  { text: '✅ Доставлено', class: 'status-delivered' },
    };
    const badge = badges[status] || { text: status, class: '' };
    return <span className={`status-badge ${badge.class}`}>{badge.text}</span>;
  };

  if (loading) {
    return (
      <div className="my-shipments-page">
        <div className="loading-state"><div className="spinner" /><p>Завантаження рейсів...</p></div>
      </div>
    );
  }

  const shipmentsTotalPages = Math.max(1, Math.ceil(filteredShipments.length / SHIPMENTS_PAGE_SIZE));
  const safeShipmentsPage = Math.min(shipmentsPage, shipmentsTotalPages - 1);
  const pagedShipments = filteredShipments.slice(safeShipmentsPage * SHIPMENTS_PAGE_SIZE, (safeShipmentsPage + 1) * SHIPMENTS_PAGE_SIZE);

  return (
    <div className="my-shipments-page">
      <div className="page-header">
        <div className="header-content">
          <h1>🚚 Мої Рейси</h1>
          <p className="subtitle">Призначені вам доставки</p>
        </div>
        <button className="btn-refresh" onClick={loadMyShipments}>🔄 Оновити</button>
      </div>

      {/* ── GPS статус-банер (дані з глобального GPSContext) ── */}
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
            <span>🔒 Геолокацію заблоковано. Дозвольте доступ у браузері та оновіть сторінку.</span>
          )}
          {gpsStatus === 'unavailable' && (
            <span>📡 Геолокація недоступна на цьому пристрої. GPS трекінг вимкнено.</span>
          )}
          {gpsStatus === 'error' && <span>❌ Геолокація не підтримується цим браузером</span>}
          {(gpsStatus === 'idle' || gpsStatus === 'no_shipment') && <span>⏳ Очікуємо GPS сигнал...</span>}
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
          <>
          <div className="shipments-grid">
            {pagedShipments.map(shipment => (
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
                    <button className="btn-complete" onClick={() => handleCompleteShipment(shipment.id, shipment.distance_km || 0)}>✅ Підтвердити доставку</button>
                  )}
                </div>
              </div>
            ))}
          </div>
          <Pagination
            currentPage={safeShipmentsPage}
            totalPages={shipmentsTotalPages}
            onPageChange={setShipmentsPage}
            totalItems={filteredShipments.length}
            itemsPerPage={SHIPMENTS_PAGE_SIZE}
          />
          </>
        )}
      </div>

      {/* Модалка підтвердження доставки з фактичним пробігом */}
      {completeModal && (
        <ModalPortal>
          <div className="modal-overlay" style={{ zIndex: 1000 }} onClick={() => !isCompleting && setCompleteModal(null)}>
            <div className="modal" style={{ maxWidth: '420px' }} onClick={e => e.stopPropagation()}>
              <h3>✅ Підтвердити доставку</h3>
              <p style={{ color: 'var(--text-muted)', fontSize: '14px', marginBottom: '20px' }}>
                Вкажіть фактичний пробіг цього рейсу. Система автоматично запише витрату пального та оновить одометр авто.
              </p>
              <div style={{ marginBottom: '8px' }}>
                <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--text-muted)', marginBottom: '6px' }}>
                  Фактичний пробіг (км) <span style={{ color: '#ef4444' }}>*</span>
                </label>
                <input
                  type="number"
                  min="1"
                  className="erp-input"
                  style={{ width: '100%', boxSizing: 'border-box' }}
                  value={actualKm}
                  onChange={e => setActualKm(e.target.value)}
                  disabled={isCompleting}
                  autoFocus
                />
                {completeModal.plannedKm > 0 && (
                  <span style={{ display: 'block', fontSize: '11px', color: '#64748b', marginTop: '5px' }}>
                    Планова відстань маршруту: <strong>{completeModal.plannedKm} км</strong>. Якщо маршрут відрізнявся — введіть реальний пробіг.
                  </span>
                )}
              </div>
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px', marginTop: '24px', paddingTop: '16px', borderTop: '1px solid var(--border)' }}>
                <button className="btn btn-secondary" onClick={() => setCompleteModal(null)} disabled={isCompleting}>Скасувати</button>
                <button className="btn btn-primary" onClick={handleConfirmComplete} disabled={isCompleting || !actualKm || parseInt(actualKm) < 1}>
                  {isCompleting ? 'Зберігаємо...' : 'Завершити рейс'}
                </button>
              </div>
            </div>
          </div>
        </ModalPortal>
      )}
    </div>
  );
}
