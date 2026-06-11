import { useState, useEffect } from 'react';
import toast from 'react-hot-toast';
import { api } from '../api/client';
import { usePermissions } from '../hooks/usePermissions';
import { PaywallScreen } from '../components/FeatureGate';
import Pagination from '../components/Pagination';
import './MaintenanceSchedule.css';

interface MaintenanceItem {
  id: number;
  vehicle_id: string;
  vehicle_plate: string;
  scheduled_record_id?: string;
  service_type: string; // OIL_CHANGE, TIRE_ROTATION, FILTER_REPLACEMENT, INSPECTION
  last_service_date: string;
  next_service_date: string;
  days_remaining: number;
  mileage_since_service: number;
  recommended_mileage: number;
  current_odometer: number;
  priority: 'LOW' | 'MEDIUM' | 'HIGH';
  status: 'COMPLETED' | 'DUE' | 'SCHEDULED' | 'OVERDUE';
}

interface MaintenanceScheduleData {
  schedules: MaintenanceItem[];
  total_overdue: number;
  total_due_soon: number;
  average_compliance: number;
}

export function MaintenanceSchedule() {
  const perms = usePermissions();
  const hasAccess = perms.hasFeature('predictive_maintenance');
  const [scheduleData, setScheduleData] = useState<MaintenanceScheduleData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filterPriority, setFilterPriority] = useState<'ALL' | 'LOW' | 'MEDIUM' | 'HIGH'>('ALL');
  const [detailItem, setDetailItem] = useState<MaintenanceItem | null>(null);
  const [schedulePage, setSchedulePage] = useState(0);
  const [schedulingId, setSchedulingId] = useState<number | null>(null);
  const SCHEDULE_PAGE_SIZE = 12;

  useEffect(() => { setSchedulePage(0); }, [filterPriority]);

  useEffect(() => {
    if (perms.authLoading) return;
    if (!hasAccess) {
      setLoading(false);
      return;
    }
    fetchMaintenanceSchedule();
  }, [hasAccess, perms.authLoading]);

  const fetchMaintenanceSchedule = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.analytics.getPredictiveMaintenanceSchedule();
      setScheduleData(response);
    } catch (err: any) {
      if (err.response?.status === 402 || String(err.message || '').includes('402')) {
        // Тариф не дозволяє — показуємо paywall
        setError(null);
      } else {
        setError('Помилка при завантаженні графіка ТО');
      }
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string, priority: string): string => {
    if (status === 'OVERDUE') return 'danger';
    if (status === 'DUE') {
      if (priority === 'HIGH') return 'warning';
      if (priority === 'MEDIUM') return 'info';
      return 'secondary';
    }
    if (status === 'SCHEDULED') {
      if (priority === 'HIGH') return 'warning';
      if (priority === 'MEDIUM') return 'info';
      return 'success';
    }
    return 'secondary';
  };

  const getStatusLabel = (status: string): string => {
    const labels: { [key: string]: string } = {
      COMPLETED: '✅ Завершено',
      DUE: '📝 Потребує планування',
      SCHEDULED: '📅 Заплановано',
      OVERDUE: '⚠️ ПРОСТРОЧЕНО',
    };
    return labels[status] || status;
  };

  const getPriorityLabel = (priority: string): string => {
    const labels: { [key: string]: string } = {
      LOW: '🟢 Низький',
      MEDIUM: '🟡 Середній',
      HIGH: '🔴 Високий',
    };
    return labels[priority] || priority;
  };

  const getServiceLabel = (type: string): string => {
    const labels: { [key: string]: string } = {
      OIL_CHANGE: '🛢️ Заміна масла',
      TIRE_ROTATION: '🛞 Ротація шин',
      FILTER_REPLACEMENT: '🔧 Заміна фільтрів',
      INSPECTION: '🔍 Планове ТО',
    };
    return labels[type] || type;
  };

  const handleScheduleMaintenance = async (item: MaintenanceItem) => {
    setSchedulingId(item.id);
    try {
      await api.vehicles.scheduleMaintenance(item.vehicle_id, {
        odometer_km: Math.max(0, Math.round(item.current_odometer ?? item.mileage_since_service)),
        service_type: item.service_type,
        scheduled_for: item.next_service_date,
        description: `${getServiceLabel(item.service_type)}: заплановано з графіка предиктивного обслуговування`,
      });
      toast.success(`📅 ТО для ${item.vehicle_plate} збережено в історії обслуговування`);
      await fetchMaintenanceSchedule();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Помилка планування ТО');
    } finally {
      setSchedulingId(null);
    }
  };

  const filteredSchedules = (scheduleData?.schedules ?? []).filter(item =>
    filterPriority === 'ALL' || item.priority === filterPriority
  );

  if (!hasAccess || error) {
    return (
      <PaywallScreen
        feature="predictive_maintenance"
        description="Автоматичне планування ТО на основі пробігу, історії обслуговування та інтервалів виробника. Менше непланових простоїв — більше вчасних рейсів."
      />
    );
  }

  if (loading) {
    return <div className="maintenance-container"><div className="loading">Завантаження графіка ТО...</div></div>;
  }

  const scheduleTotalPages = Math.max(1, Math.ceil(filteredSchedules.length / SCHEDULE_PAGE_SIZE));
  const safeSchedulePage = Math.min(schedulePage, scheduleTotalPages - 1);
  const pagedSchedules = filteredSchedules.slice(safeSchedulePage * SCHEDULE_PAGE_SIZE, (safeSchedulePage + 1) * SCHEDULE_PAGE_SIZE);

  return (
    <div className="maintenance-container">
      <div className="maintenance-header">
        <h1>📅 Графік Предиктивного Обслуговування</h1>
        <p>Планування ТО на основі рівня використання та історії обслуговування</p>
      </div>

      <div className="maintenance-stats">
        <div className="stat-card">
          <div className="stat-value" style={{ color: '#ef4444' }}>
            {scheduleData?.total_overdue || 0}
          </div>
          <div className="stat-label">⚠️ Прострочених</div>
        </div>
        <div className="stat-card">
          <div className="stat-value" style={{ color: '#f59e0b' }}>
            {scheduleData?.total_due_soon || 0}
          </div>
          <div className="stat-label">📅 Найближчі 30 днів</div>
        </div>
        <div className="stat-card">
          <div className="stat-value" style={{ color: '#10b981' }}>
            {Math.round(scheduleData?.average_compliance || 0)}%
          </div>
          <div className="stat-label">✅ Середнє дотримання</div>
        </div>
      </div>

      <div className="maintenance-filters">
        <label>Фільтр по пріоритету:</label>
        <div className="filter-buttons">
          {['ALL', 'LOW', 'MEDIUM', 'HIGH'].map(priority => (
            <button
              key={priority}
              className={`filter-btn ${filterPriority === priority ? 'active' : ''}`}
              onClick={() => setFilterPriority(priority as any)}
            >
              {priority === 'ALL' ? '📊 Всі' : getPriorityLabel(priority).split(' ')[1]}
            </button>
          ))}
        </div>
      </div>

      <div className="maintenance-grid">
        {filteredSchedules.length === 0 ? (
          <div className="no-data">
            <p>📭 Немає послуг технічного обслуговування для фільтра</p>
          </div>
        ) : (
          pagedSchedules.map(item => (
            <div key={item.id} className={`maintenance-card ${getStatusColor(item.status, item.priority)}`}>
              <div className="card-header">
                <div className="vehicle-info">
                  <h3>{item.vehicle_plate}</h3>
                  <span className="service-type">{getServiceLabel(item.service_type)}</span>
                </div>
                <div className="status-badge">{getStatusLabel(item.status)}</div>
              </div>

              <div className="card-content">
                <div className="info-row">
                  <span className="label">Пріоритет:</span>
                  <span className="value">{getPriorityLabel(item.priority)}</span>
                </div>
                <div className="info-row">
                  <span className="label">Останнє обслуговування:</span>
                  <span className="value">{new Date(item.last_service_date).toLocaleDateString('uk-UA')}</span>
                </div>
                <div className="info-row">
                  <span className="label">Наступне обслуговування:</span>
                  <span className="value">{new Date(item.next_service_date).toLocaleDateString('uk-UA')}</span>
                </div>

                <div className="days-remaining">
                  <div className="days-label">Днів до обслуговування</div>
                  <div className={`days-value ${item.days_remaining < 0 ? 'overdue' : item.days_remaining < 7 ? 'warning' : 'normal'}`}>
                    {item.days_remaining < 0 ? `Прострочено на ${Math.abs(item.days_remaining)} днів` : `${item.days_remaining} днів`}
                  </div>
                </div>

                <div className="mileage-bar">
                  <div className="bar-label">
                    <span>Пробіг з останнього ТО</span>
                    <span className="percentage">{Math.round((item.mileage_since_service / item.recommended_mileage) * 100)}%</span>
                  </div>
                  <div className="progress-bar">
                    <div
                      className="progress-fill"
                      style={{
                        width: `${Math.min((item.mileage_since_service / item.recommended_mileage) * 100, 100)}%`,
                        backgroundColor: item.mileage_since_service > item.recommended_mileage ? '#ef4444' : '#3b82f6'
                      }}
                    />
                  </div>
                  <div className="mileage-info">
                    <span>{Math.round(item.mileage_since_service)} км</span>
                    <span>/</span>
                    <span>{Math.round(item.recommended_mileage)} км</span>
                  </div>
                </div>
              </div>

              <div className="card-actions">
                <button
                  className="btn-schedule"
                  disabled={item.status === 'SCHEDULED' || schedulingId === item.id}
                  onClick={() => handleScheduleMaintenance(item)}
                >
                  {item.status === 'SCHEDULED' ? '📅 Вже заплановано' : schedulingId === item.id ? 'Збереження...' : '📅 Запланувати ТО'}
                </button>
                <button className="btn-details" onClick={() => setDetailItem(item)}>
                  📋 Деталі
                </button>
              </div>
            </div>
          ))
        )}
      </div>
      <Pagination
        currentPage={safeSchedulePage}
        totalPages={scheduleTotalPages}
        onPageChange={setSchedulePage}
        totalItems={filteredSchedules.length}
        itemsPerPage={SCHEDULE_PAGE_SIZE}
      />

      {detailItem && (
        <div className="maint-modal-overlay" onClick={() => setDetailItem(null)}>
          <div className="maint-modal" onClick={(e) => e.stopPropagation()}>
            <div className="maint-modal-header">
              <h2>📋 {detailItem.vehicle_plate}</h2>
              <button className="maint-modal-close" onClick={() => setDetailItem(null)}>✕</button>
            </div>
            <div className="maint-modal-body">
              <div className="maint-modal-block">
                <span className="label">Вид обслуговування:</span>
                <span className="value">{getServiceLabel(detailItem.service_type)}</span>
              </div>
              <div className="maint-modal-block">
                <span className="label">Пріоритет:</span>
                <span className="value">{getPriorityLabel(detailItem.priority)}</span>
              </div>
              <div className="maint-modal-block">
                <span className="label">Статус:</span>
                <span className="value">{getStatusLabel(detailItem.status)}</span>
              </div>
              <div className="maint-modal-block">
                <span className="label">Останнє ТО:</span>
                <span className="value">{new Date(detailItem.last_service_date).toLocaleDateString('uk-UA')}</span>
              </div>
              <div className="maint-modal-block">
                <span className="label">Наступне ТО:</span>
                <span className="value">{new Date(detailItem.next_service_date).toLocaleDateString('uk-UA')}</span>
              </div>
              <div className="maint-modal-block">
                <span className="label">Пробіг з моменту ТО:</span>
                <span className="value">
                  {Math.round(detailItem.mileage_since_service)} / {Math.round(detailItem.recommended_mileage)} км
                  {' '}({Math.round((detailItem.mileage_since_service / detailItem.recommended_mileage) * 100)}%)
                </span>
              </div>
              <div className="maint-modal-block">
                <span className="label">Залишилось:</span>
                <span className="value">
                  {detailItem.days_remaining < 0
                    ? `Прострочено на ${Math.abs(detailItem.days_remaining)} днів`
                    : `${detailItem.days_remaining} днів`}
                </span>
              </div>
              <div className="maint-modal-note">
                ℹ️ Графік розраховується на основі пробігу з останнього ТО та рекомендованого інтервалу обслуговування.
              </div>
            </div>
            <div className="maint-modal-footer">
              <button className="btn-details" onClick={() => setDetailItem(null)}>Закрити</button>
            </div>
          </div>
        </div>
      )}


    </div>
  );
}

export default MaintenanceSchedule;
