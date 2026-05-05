import { useEffect, useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
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

export default function MyShipments() {
  const { token } = useAuth();
  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'pending' | 'in_transit' | 'delivered'>('pending');

  useEffect(() => {
    loadMyShipments();
  }, [token]);

  const loadMyShipments = async () => {
    if (!token) return;
    try {
      setLoading(true);
      const res = await fetch('/api/inventory/shipments/my', {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      
      if (!res.ok) {
        throw new Error('Не вдалося завантажити рейси');
      }
      
      const data = await res.json();
      console.log('📦 My shipments:', data);
      setShipments(Array.isArray(data) ? data : []);
    } catch (err: any) {
      console.error('Error loading shipments:', err);
      toast.error(err.message || 'Помилка завантаження рейсів');
    } finally {
      setLoading(false);
    }
  };

  const handleStartShipment = async (shipmentId: string) => {
    if (!token) return;
    const toastId = toast.loading('Починаємо рейс...');
    try {
      const res = await fetch(`/api/inventory/shipments/${shipmentId}/start`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` }
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Не вдалося почати рейс');
      }

      // Позначаємо всі повідомлення як прочитані (повідомлення про призначення рейсу вже не актуальне)
      try {
        await fetch('/api/notifications/mark-all-read', {
          method: 'POST',
          headers: { 'Authorization': `Bearer ${token}` }
        });
      } catch (notifErr) {
        console.log('Could not mark notifications as read:', notifErr);
      }

      toast.success('Рейс розпочато! 🚚', { id: toastId });
      loadMyShipments();
    } catch (err: any) {
      console.error('Error starting shipment:', err);
      toast.error(err.message || 'Помилка початку рейсу', { id: toastId });
    }
  };

  const handleCompleteShipment = async (shipmentId: string) => {
    if (!token) return;
    const toastId = toast.loading('Завершуємо рейс...');
    try {
      const res = await fetch(`/api/inventory/shipments/${shipmentId}/receive`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` }
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Не вдалося завершити рейс');
      }

      toast.success('Рейс завершено! ✅', { id: toastId });
      loadMyShipments();
    } catch (err: any) {
      console.error('Error completing shipment:', err);
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
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getStatusBadge = (status: string) => {
    const badges = {
      'PENDING': { text: '⏳ Очікує', class: 'status-pending' },
      'IN_TRANSIT': { text: '🚚 В дорозі', class: 'status-in-transit' },
      'DELIVERED': { text: '✅ Доставлено', class: 'status-delivered' }
    };
    const badge = badges[status as keyof typeof badges] || { text: status, class: '' };
    return <span className={`status-badge ${badge.class}`}>{badge.text}</span>;
  };

  if (loading) {
    return (
      <div className="my-shipments-page">
        <div className="loading-state">
          <div className="spinner"></div>
          <p>Завантаження рейсів...</p>
        </div>
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
        <button className="btn-refresh" onClick={loadMyShipments}>
          🔄 Оновити
        </button>
      </div>

      <div className="tabs-container">
        <button
          className={`tab-btn ${activeTab === 'pending' ? 'active' : ''}`}
          onClick={() => setActiveTab('pending')}
        >
          ⏳ Очікують відправки ({shipments.filter(s => s.status === 'PENDING').length})
        </button>
        <button
          className={`tab-btn ${activeTab === 'in_transit' ? 'active' : ''}`}
          onClick={() => setActiveTab('in_transit')}
        >
          🚚 В дорозі ({shipments.filter(s => s.status === 'IN_TRANSIT').length})
        </button>
        <button
          className={`tab-btn ${activeTab === 'delivered' ? 'active' : ''}`}
          onClick={() => setActiveTab('delivered')}
        >
          ✅ Доставлено ({shipments.filter(s => s.status === 'DELIVERED').length})
        </button>
      </div>

      <div className="shipments-container">
        {filteredShipments.length === 0 ? (
          <div className="empty-state">
            <div className="empty-icon">📦</div>
            <h3>Рейсів не знайдено</h3>
            <p>
              {activeTab === 'pending' && 'У вас поки немає рейсів, які очікують відправки'}
              {activeTab === 'in_transit' && 'У вас поки немає активних рейсів'}
              {activeTab === 'delivered' && 'У вас поки немає завершених рейсів'}
            </p>
          </div>
        ) : (
          <div className="shipments-grid">
            {filteredShipments.map(shipment => (
              <div key={shipment.id} className="shipment-card">
                <div className="shipment-header">
                  <div className="shipment-id">
                    <strong>Рейс #{shipment.id.substring(0, 8)}</strong>
                    {getStatusBadge(shipment.status)}
                  </div>
                  <div className="vehicle-info">
                    🚗 {shipment.vehicle_plate || shipment.vehicle_id.substring(0, 8)}
                  </div>
                </div>

                <div className="route-info">
                  <div className="route-point">
                    <span className="route-label">📍 Від:</span>
                    <strong>{shipment.from_warehouse_name || 'Склад'}</strong>
                  </div>
                  <div className="route-arrow">→</div>
                  <div className="route-point">
                    <span className="route-label">📍 До:</span>
                    <strong>{shipment.to_warehouse_name || 'Склад'}</strong>
                  </div>
                </div>

                {shipment.items && shipment.items.length > 0 && (
                  <div className="items-list">
                    <strong>Вантаж:</strong>
                    <ul>
                      {shipment.items.map((item, idx) => (
                        <li key={idx}>
                          {item.resource_name}: {item.quantity} {item.unit}
                        </li>
                      ))}
                    </ul>
                  </div>
                )}

                <div className="shipment-times">
                  <div className="time-item">
                    <span className="time-label">Створено:</span>
                    <span className="time-value">{formatDate(shipment.created_at)}</span>
                  </div>
                  {shipment.started_at && (
                    <div className="time-item">
                      <span className="time-label">Відправлено:</span>
                      <span className="time-value">{formatDate(shipment.started_at)}</span>
                    </div>
                  )}
                  {shipment.delivered_at && (
                    <div className="time-item">
                      <span className="time-label">Доставлено:</span>
                      <span className="time-value">{formatDate(shipment.delivered_at)}</span>
                    </div>
                  )}
                </div>

                <div className="shipment-actions">
                  {shipment.status === 'PENDING' && (
                    <button
                      className="btn-start"
                      onClick={() => handleStartShipment(shipment.id)}
                    >
                      🚀 Почати рейс
                    </button>
                  )}
                  {shipment.status === 'IN_TRANSIT' && (
                    <button
                      className="btn-complete"
                      onClick={() => handleCompleteShipment(shipment.id)}
                    >
                      ✅ Підтвердити доставку
                    </button>
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
