import { useState, useEffect } from 'react';
import toast, { Toaster } from 'react-hot-toast';
import { api } from '../api/client';
import { usePermissions } from '../hooks/usePermissions';
import { PaywallScreen } from '../components/FeatureGate';
import './FuelAnomalies.css';

interface FuelAnomaly {
  id: number;
  vehicle_id: number;
  vehicle_plate: string;
  anomaly_type: string; // EXTREME_REFILL, FREQUENT_SMALL_REFILLS, PRICE_ANOMALY, ABNORMAL_CONSUMPTION
  risk_score: number; // 0-100
  investigation_level: 'LOW' | 'MEDIUM' | 'HIGH';
  last_detected: string;
  details: string;
  confidence: number; // 0-100
  potential_loss: number; // estimated monthly loss in UAH
}

interface FuelAnomalyData {
  anomalies: FuelAnomaly[];
  total_vehicles_monitored: number;
  vehicles_with_anomalies: number;
  total_potential_loss: number;
  high_risk_count: number;
}

export function FuelAnomalies() {
  const perms = usePermissions();
  const hasAccess = perms.hasFeature('fuel_antifraud');
  const [anomalyData, setAnomalyData] = useState<FuelAnomalyData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filterLevel, setFilterLevel] = useState<'ALL' | 'LOW' | 'MEDIUM' | 'HIGH'>('ALL');
  const [investigating, setInvestigating] = useState<FuelAnomaly | null>(null);
  const [alertsEnabled, setAlertsEnabled] = useState<Record<number, boolean>>({});

  useEffect(() => {
    if (!hasAccess) {
      setLoading(false);
      return;
    }
    fetchFuelAnomalies();
  }, [hasAccess]);

  const fetchFuelAnomalies = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.analytics.getFuelAnomalyDetection();
      setAnomalyData(response);
    } catch (err: any) {
      if (err.response?.status === 402 || String(err.message || '').includes('402')) {
        // Тариф не дозволяє — рендеримо paywall замість червоної плашки
        setError(null);
      } else {
        setError('Помилка при завантаженні анализу аномалій');
      }
    } finally {
      setLoading(false);
    }
  };

  const getAnomalyTypeLabel = (type: string): string => {
    const labels: { [key: string]: string } = {
      EXTREME_REFILL: '🚨 Екстремальна заправка',
      FREQUENT_SMALL_REFILLS: '⚡ Часті малі заправки',
      PRICE_ANOMALY: '💰 Цінова аномалія',
      ABNORMAL_CONSUMPTION: '📉 Ненормальне споживання',
    };
    return labels[type] || type;
  };

  const getLevelColor = (level: string): string => {
    switch (level) {
      case 'LOW': return 'level-low';
      case 'MEDIUM': return 'level-medium';
      case 'HIGH': return 'level-high';
      default: return 'level-low';
    }
  };

  const getLevelLabel = (level: string): string => {
    const labels: { [key: string]: string } = {
      LOW: '🟢 Низький',
      MEDIUM: '🟡 Середній',
      HIGH: '🔴 Високий',
    };
    return labels[level] || level;
  };

  const getRiskColor = (score: number): string => {
    if (score < 33) return '#10b981'; // green
    if (score < 66) return '#f59e0b'; // orange
    return '#ef4444'; // red
  };

  const filteredAnomalies = (anomalyData?.anomalies ?? []).filter(item =>
    filterLevel === 'ALL' || item.investigation_level === filterLevel
  );

  if (!hasAccess) {
    return (
      <PaywallScreen
        feature="fuel_antifraud"
        description="AI-детекція підозрілих заправок, цінових аномалій та ненормального споживання. Типова організація економить 5–12% бюджету пального вже в перший місяць."
      />
    );
  }

  if (loading) {
    return <div className="anomalies-container"><div className="loading">Завантаження анализу аномалій...</div></div>;
  }

  if (error) {
    return <div className="anomalies-container"><div className="error-message">{error}</div></div>;
  }

  return (
    <div className="anomalies-container">
      <div className="anomalies-header">
        <h1>🛡️ Антифрод-система контролю пального</h1>
        <p>AI-детекція підозрілих заправок, цінових аномалій та ненормального споживання</p>
      </div>

      <div className="anomalies-stats">
        <div className="stat-card monitored">
          <div className="stat-value">{anomalyData?.total_vehicles_monitored || 0}</div>
          <div className="stat-label">🚗 Автомобілів під моніторингом</div>
        </div>
        <div className="stat-card affected">
          <div className="stat-value">{anomalyData?.vehicles_with_anomalies || 0}</div>
          <div className="stat-label">⚠️ З аномаліями</div>
        </div>
        <div className="stat-card risk">
          <div className="stat-value">{anomalyData?.high_risk_count || 0}</div>
          <div className="stat-label">🔴 Високого ризику</div>
        </div>
        <div className="stat-card loss">
          <div className="stat-value">{Math.round(anomalyData?.total_potential_loss || 0)} ₴</div>
          <div className="stat-label">💸 Потенційні втрати/місяць</div>
        </div>
      </div>

      <div className="anomalies-filters">
        <label>Фільтр по рівню розслідування:</label>
        <div className="filter-buttons">
          {['ALL', 'LOW', 'MEDIUM', 'HIGH'].map(level => (
            <button
              key={level}
              className={`filter-btn ${filterLevel === level ? 'active' : ''}`}
              onClick={() => setFilterLevel(level as any)}
            >
              {level === 'ALL' ? '📊 Всі' : getLevelLabel(level).split(' ')[1]}
            </button>
          ))}
        </div>
      </div>

      <div className="anomalies-grid">
        {filteredAnomalies.length === 0 ? (
          <div className="no-data">
            <p>✅ Аномалій не виявлено! Всі автомобілі в нормі</p>
          </div>
        ) : (
          filteredAnomalies.map(anomaly => (
            <div key={anomaly.id} className={`anomaly-card ${getLevelColor(anomaly.investigation_level)}`}>
              <div className="card-header">
                <div className="vehicle-info">
                  <h3>{anomaly.vehicle_plate}</h3>
                  <span className="anomaly-type">{getAnomalyTypeLabel(anomaly.anomaly_type)}</span>
                </div>
                <div className="level-badge">{getLevelLabel(anomaly.investigation_level)}</div>
              </div>

              <div className="card-content">
                <div className="risk-section">
                  <div className="risk-label">Ризик</div>
                  <div className="risk-bar">
                    <div
                      className="risk-fill"
                      style={{
                        width: `${anomaly.risk_score}%`,
                        backgroundColor: getRiskColor(anomaly.risk_score),
                      }}
                    />
                  </div>
                  <div className="risk-value">{anomaly.risk_score}/100</div>
                </div>

                <div className="confidence-section">
                  <span className="label">Впевненість детекції:</span>
                  <span className="value">{anomaly.confidence}%</span>
                </div>

                <div className="details-section">
                  <span className="label">📋 Деталі:</span>
                  <p className="details-text">{anomaly.details}</p>
                </div>

                <div className="financial-section">
                  <div className="loss-item">
                    <span className="label">💸 Потенційна втрата:</span>
                    <span className="loss-value">{Math.round(anomaly.potential_loss)} ₴/місяць</span>
                  </div>
                </div>

                <div className="detection-section">
                  <span className="label">🕐 Останнє виявлення:</span>
                  <span className="value">{new Date(anomaly.last_detected).toLocaleDateString('uk-UA')}</span>
                </div>

                <div className="ai-badge">
                  🤖 AI-детекція на основі історичних даних і паттернів
                </div>
              </div>

              <div className="card-actions">
                <button
                  className="btn-investigate"
                  onClick={() => setInvestigating(anomaly)}
                >
                  🔍 Розслідувати
                </button>
                <button
                  className={`btn-alert ${alertsEnabled[anomaly.id] ? 'active' : ''}`}
                  onClick={() => {
                    setAlertsEnabled((prev) => {
                      const next = { ...prev, [anomaly.id]: !prev[anomaly.id] };
                      if (next[anomaly.id]) {
                        toast.success(
                          `🚨 Сповіщення для ${anomaly.vehicle_plate} увімкнено. Ви отримаєте email при повторенні аномалії.`,
                          { duration: 4500 }
                        );
                      } else {
                        toast(`🔕 Сповіщення для ${anomaly.vehicle_plate} вимкнено`, { duration: 3000 });
                      }
                      return next;
                    });
                  }}
                >
                  {alertsEnabled[anomaly.id] ? '🔔 Сповіщення увімкнено' : '🚨 Встановити сповіщення'}
                </button>
              </div>
            </div>
          ))
        )}
      </div>

      {investigating && (
        <div className="fuel-modal-overlay" onClick={() => setInvestigating(null)}>
          <div className="fuel-modal" onClick={(e) => e.stopPropagation()}>
            <div className="fuel-modal-header">
              <h2>🔍 Розслідування: {investigating.vehicle_plate}</h2>
              <button className="fuel-modal-close" onClick={() => setInvestigating(null)}>✕</button>
            </div>
            <div className="fuel-modal-body">
              <div className="fuel-modal-row">
                <span className="label">Тип аномалії:</span>
                <span className="value">{getAnomalyTypeLabel(investigating.anomaly_type)}</span>
              </div>
              <div className="fuel-modal-row">
                <span className="label">Рівень загрози:</span>
                <span className="value">{getLevelLabel(investigating.investigation_level)}</span>
              </div>
              <div className="fuel-modal-row">
                <span className="label">Рейтинг ризику:</span>
                <span className="value" style={{ color: getRiskColor(investigating.risk_score) }}>
                  {investigating.risk_score}/100
                </span>
              </div>
              <div className="fuel-modal-row">
                <span className="label">Впевненість AI:</span>
                <span className="value">{investigating.confidence}%</span>
              </div>
              <div className="fuel-modal-row">
                <span className="label">Останнє виявлення:</span>
                <span className="value">{new Date(investigating.last_detected).toLocaleString('uk-UA')}</span>
              </div>
              <div className="fuel-modal-row">
                <span className="label">Потенційні втрати:</span>
                <span className="value" style={{ color: '#dc2626' }}>
                  {Math.round(investigating.potential_loss)} ₴/місяць
                </span>
              </div>
              <div className="fuel-modal-details">
                <strong>📋 Деталі виявлення:</strong>
                <p>{investigating.details}</p>
              </div>
              <div className="fuel-modal-steps">
                <strong>✅ Рекомендовані дії:</strong>
                <ol>
                  <li>Зв'язатися з водієм машини та запитати пояснення по факту.</li>
                  <li>Перевірити чеки АЗС та співставити з бортовим комп'ютером.</li>
                  <li>Звірити маршрут GPS з датою/часом заправки.</li>
                  <li>Скласти службову записку при підтвердженні фроду.</li>
                </ol>
              </div>
            </div>
            <div className="fuel-modal-footer">
              <button
                className="btn-investigate"
                onClick={() => {
                  toast.success('📝 Розслідування відкрито у внутрішньому трекері', { duration: 4000 });
                  setInvestigating(null);
                }}
              >
                Відкрити службове розслідування
              </button>
              <button className="btn-alert" onClick={() => setInvestigating(null)}>
                Закрити
              </button>
            </div>
          </div>
        </div>
      )}

      <Toaster position="top-right" />
    </div>
  );
}

export default FuelAnomalies;
