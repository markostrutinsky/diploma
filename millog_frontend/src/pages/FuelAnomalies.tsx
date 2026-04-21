import { useState, useEffect } from 'react';
import { api } from '../api/client';
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
  const [anomalyData, setAnomalyData] = useState<FuelAnomalyData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filterLevel, setFilterLevel] = useState<'ALL' | 'LOW' | 'MEDIUM' | 'HIGH'>('ALL');

  useEffect(() => {
    fetchFuelAnomalies();
  }, []);

  const fetchFuelAnomalies = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.analytics.getFuelAnomalyDetection();
      setAnomalyData(response);
    } catch (err: any) {
      if (err.response?.status === 402) {
        setError('Антифрод-система доступна тільки для PRO підписки');
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

  const filteredAnomalies = anomalyData?.anomalies.filter(item =>
    filterLevel === 'ALL' || item.investigation_level === filterLevel
  ) || [];

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
                <button className="btn-investigate">🔍 Розслідувати</button>
                <button className="btn-alert">🚨 Встановити оповіщення</button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

export default FuelAnomalies;
