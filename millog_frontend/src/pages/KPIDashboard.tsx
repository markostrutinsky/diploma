import React, { useEffect, useState } from 'react'
import { api } from '../api/client'
import { usePermissions } from '../hooks/usePermissions'
import { PaywallScreen } from '../components/FeatureGate'
import './KPIDashboard.css'

interface KPIDashboard {
  reporting_period: string
  inventory_total_value: number
  sla: {
    on_time_percent: number
    total_requests: number
    on_time_count: number
    late_count: number
    avg_delay_hours: number
  }
  tco: {
    total_fuel_cost: number
    total_units_shipped: number
    cost_per_unit: number
    trend: string
  }
  risk: {
    critical_resources_percent: number
    critical_resources: string[]
    at_risk_count: number
    total_resources: number
  }
  depletion_forecast: {
    within_7_days: string[]
    within_14_days: string[]
    within_30_days: string[]
    action_required: boolean
  }
}

const KPIDashboard: React.FC = () => {
  const perms = usePermissions()
  const hasAccess = perms.hasFeature('advanced_analytics')
  const [kpiData, setKpiData] = useState<KPIDashboard | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (perms.authLoading) return
    if (!hasAccess) {
      setLoading(false)
      return
    }
    fetchKPIData()
  }, [hasAccess, perms.authLoading])

  const fetchKPIData = async () => {
    try {
      setLoading(true)
      const response = await api.analytics.getAdvancedKPIs()
      setKpiData(response)
      setError(null)
    } catch (err: any) {
      const message = err.message || ''
      if (message.includes('402')) {
        // Тариф не дозволяє — показуємо красиву заглушку, а не червону плашку
        setError(null)
      } else {
        setError('Помилка при завантаженні KPI даних')
      }
    } finally {
      setLoading(false)
    }
  }

  if (!hasAccess || error) {
    return (
      <PaywallScreen
        feature="advanced_analytics"
        description="Показники SLA, TCO, ризиків та прогноз дефіциту допомагають керівнику приймати стратегічні рішення на основі даних, а не відчуттів."
      />
    )
  }

  if (loading) {
    return <div className="loading">Завантаження KPI даних...</div>
  }

  if (!kpiData) {
    return <div className="empty-state">Немає даних KPI</div>
  }

  const getMetricColor = (value: number, thresholds: { good: number; warning: number }): 'good' | 'warning' | 'danger' => {
    if (value >= thresholds.good) return 'good'
    if (value >= thresholds.warning) return 'warning'
    return 'danger'
  }

  const slaColor = getMetricColor(kpiData.sla.on_time_percent, { good: 95, warning: 85 })
  const riskColor = getMetricColor(100 - kpiData.risk.critical_resources_percent, { good: 95, warning: 85 })

  return (
    <div className="kpi-dashboard">
      <div className="kpi-header">
        <h1>🎯 Розширена панель ефективності</h1>
        <p className="period">{kpiData.reporting_period}</p>
      </div>

      <div className="kpi-grid">
        {/* SLA Metric */}
        <div className={`kpi-card metric-${slaColor}`}>
          <div className="kpi-icon">📊</div>
          <div className="kpi-content">
            <h3>SLA - Вчасність</h3>
            <div className="kpi-value">{kpiData.sla.on_time_percent.toFixed(1)}%</div>
            <div className="kpi-meta">
              {kpiData.sla.on_time_count} вчасно з {kpiData.sla.total_requests}
            </div>
            <div className="kpi-detail">Затримка: {kpiData.sla.avg_delay_hours.toFixed(1)} годин</div>
          </div>
          <div className="kpi-status-bar">
            <div 
              className="status-fill" 
              style={{ width: `${Math.min(kpiData.sla.on_time_percent, 100)}%` }}
            ></div>
          </div>
        </div>

        {/* Inventory Value Metric */}
        <div className="kpi-card metric-good">
          <div className="kpi-icon">🏦</div>
          <div className="kpi-content">
            <h3>Вартість майна</h3>
            <div className="kpi-value">{(kpiData.inventory_total_value || 0).toLocaleString('uk-UA', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} ₴</div>
            <div className="kpi-meta">Загальна балансова вартість залишків</div>
          </div>
        </div>

        {/* TCO Metric */}
        <div className="kpi-card metric-good">
          <div className="kpi-icon">💰</div>
          <div className="kpi-content">
            <h3>TCO - Вартість на одиницю</h3>
            <div className="kpi-value">{kpiData.tco.cost_per_unit.toFixed(2)} ₴</div>
            <div className="kpi-meta">
              {kpiData.tco.total_units_shipped} доставлено
            </div>
            <div className="kpi-detail">
              Паливо: {kpiData.tco.total_fuel_cost.toFixed(2)} ₴ 
              <span className={kpiData.tco.trend === 'down' ? 'trend-down' : kpiData.tco.trend === 'up' ? 'trend-up' : 'trend-stable'}>
                {kpiData.tco.trend === 'down' ? ' ↓' : kpiData.tco.trend === 'up' ? ' ↑' : ' →'}
              </span>
            </div>
          </div>
        </div>

        {/* Risk Metric */}
        <div className={`kpi-card metric-${riskColor}`}>
          <div className="kpi-icon">⚠️</div>
          <div className="kpi-content">
            <h3>Ризик - Дефіцит</h3>
            <div className="kpi-value">{kpiData.risk.critical_resources_percent.toFixed(1)}%</div>
            <div className="kpi-meta">
              {kpiData.risk.at_risk_count} ресурсів під ризиком
            </div>
            {kpiData.risk.critical_resources.length > 0 && (
              <div className="kpi-alert">
                <strong>Критичні:</strong> {kpiData.risk.critical_resources.slice(0, 2).join(', ')}
              </div>
            )}
          </div>
        </div>

        {/* Depletion Forecast */}
        <div className={`kpi-card metric-${kpiData.depletion_forecast.action_required ? 'danger' : 'good'}`}>
          <div className="kpi-icon">📉</div>
          <div className="kpi-content">
            <h3>Прогноз - Скорочення запасів</h3>
            <div className="kpi-forecast">
              <div className="forecast-item urgent">
                <span className="forecast-days">7 днів</span>
                <span className="forecast-count">{kpiData.depletion_forecast.within_7_days.length} товарів</span>
              </div>
              <div className="forecast-item warning">
                <span className="forecast-days">14 днів</span>
                <span className="forecast-count">{kpiData.depletion_forecast.within_14_days.length} товарів</span>
              </div>
              <div className="forecast-item info">
                <span className="forecast-days">30 днів</span>
                <span className="forecast-count">{kpiData.depletion_forecast.within_30_days.length} товарів</span>
              </div>
            </div>
            {kpiData.depletion_forecast.action_required && (
              <div className="kpi-action">⚡ Потрібна дія!</div>
            )}
          </div>
        </div>
      </div>

      {/* Details Section */}
      {kpiData.risk.critical_resources.length > 0 && (
        <div className="kpi-details-section">
          <h3>🚨 Критичні ресурси</h3>
          <ul className="resource-list">
            {kpiData.risk.critical_resources.map((resource, idx) => (
              <li key={idx} className="resource-item">
                {resource}
              </li>
            ))}
          </ul>
        </div>
      )}

      {(kpiData.depletion_forecast.within_7_days.length > 0 ||
        kpiData.depletion_forecast.within_14_days.length > 0) && (
        <div className="kpi-action-section">
          <span style={{ color: 'var(--text-muted)', fontSize: '0.9rem' }}>⚡ Є товари що закінчуються:</span>
          <button className="btn-primary" onClick={() => window.location.href = '/inventory'}>
            📦 Переглянути склад
          </button>
          <button className="btn-secondary" onClick={() => window.location.href = '/requests'}>
            📋 Створити заявку
          </button>
        </div>
      )}
    </div>
  )
}

export default KPIDashboard
