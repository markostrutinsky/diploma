import React, { useEffect, useState } from 'react'
import { api } from '../api/client'
import './GPSTracking.css'

interface Vehicle {
  vehicle_id: number
  plate_number?: string
  latitude: number
  longitude: number
  speed: number
  heading?: number | null
  timestamp: string
  updated_seconds_ago: number
}

interface FleetMapData {
  timestamp: string
  count: number
  vehicles: Vehicle[]
}

const GPSTracking: React.FC = () => {
  const [fleetData, setFleetData] = useState<FleetMapData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [autoRefresh, setAutoRefresh] = useState(true)
  const [refreshInterval, setRefreshInterval] = useState(30) // seconds

  // Модалка з детальною інформацією по конкретній машині + історія маршруту
  const [detailVehicle, setDetailVehicle] = useState<Vehicle | null>(null)
  const [trajectory, setTrajectory] = useState<any[]>([])
  const [trajectoryLoading, setTrajectoryLoading] = useState(false)

  const openDetail = async (vehicle: Vehicle) => {
    setDetailVehicle(vehicle)
    setTrajectory([])
    setTrajectoryLoading(true)
    try {
      const endTime = new Date().toISOString()
      const startTime = new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString() // 24 год тому
      const resp = await api.gps.getVehicleTrajectory(String(vehicle.vehicle_id), startTime, endTime)
      const points = Array.isArray(resp?.locations) ? resp.locations : Array.isArray(resp) ? resp : []
      setTrajectory(points)
    } catch (err) {
      setTrajectory([])
    } finally {
      setTrajectoryLoading(false)
    }
  }

  const closeDetail = () => {
    setDetailVehicle(null)
    setTrajectory([])
  }

  useEffect(() => {
    fetchFleetData()
    
    if (autoRefresh) {
      const interval = setInterval(fetchFleetData, refreshInterval * 1000)
      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval])

  const fetchFleetData = async () => {
    try {
      const response = await api.gps.getFleetMap()
      setFleetData(response)
      setError(null)
    } catch (err: any) {
      const message = err.message || ''
      if (message.includes('402')) {
        setError('GPS трекінг доступний лише для PRO та ENTERPRISE')
      } else {
        setError('Помилка при завантаженні GPS даних')
      }
    } finally {
      setLoading(false)
    }
  }

  const getSpeedColor = (speed: number): string => {
    if (speed === 0) return '#999' // Parked
    if (speed < 20) return '#ffc107' // Slow
    if (speed < 60) return '#17a2b8' // Normal
    return '#28a745' // Fast
  }

  const getHeadingRotation = (heading?: number | null): number => {
    return heading ?? 0
  }

  if (loading && !fleetData) {
    return <div className="loading">Завантаження GPS карти...</div>
  }

  if (error) {
    return (
      <div className="error-message">
        <p>{error}</p>
        <a href="/billing" className="btn-upgrade">Оновити план</a>
      </div>
    )
  }

  if (!fleetData) {
    return <div className="empty-state">Немає даних GPS трекінгу</div>
  }

  return (
    <div className="gps-tracking-container">
      <div className="gps-header">
        <h1>🗺️ Real-Time Fleet Tracking</h1>
        <div className="header-stats">
          <span className="stat">
            Машин онлайн: <strong>{fleetData.count}</strong>
          </span>
          <span className="stat">
            Оновлено: {new Date(fleetData.timestamp).toLocaleTimeString('uk-UA')}
          </span>
        </div>
      </div>

      <div className="gps-controls">
        <label className="control-label">
          <input 
            type="checkbox" 
            checked={autoRefresh}
            onChange={(e) => setAutoRefresh(e.target.checked)}
          />
          Auto-refresh
        </label>
        {autoRefresh && (
          <select 
            value={refreshInterval} 
            onChange={(e) => setRefreshInterval(parseInt(e.target.value))}
            className="refresh-interval"
          >
            <option value={10}>10 сек</option>
            <option value={20}>20 сек</option>
            <option value={30}>30 сек</option>
            <option value={60}>1 хв</option>
          </select>
        )}
        <button className="btn-refresh" onClick={fetchFleetData}>
          ↻ Оновити зараз
        </button>
      </div>

      <div className="vehicles-grid">
        {(fleetData.vehicles?.length ?? 0) > 0 ? (
          (fleetData.vehicles ?? []).map((vehicle) => (
            <div key={vehicle.vehicle_id} className="vehicle-card">
              <div 
                className="vehicle-marker"
                style={{
                  backgroundColor: getSpeedColor(vehicle.speed),
                  transform: `rotate(${getHeadingRotation(vehicle.heading)}deg)`
                }}
              >
                {vehicle.speed > 0 ? '▶' : '◼'}
              </div>
              
              <div className="vehicle-info">
                <div className="vehicle-plate">
                  {vehicle.plate_number || `Авто №${vehicle.vehicle_id}`}
                </div>
                
                <div className="vehicle-location">
                  <span className="coord lat">
                    {Number(vehicle.latitude).toFixed(4)}°N
                  </span>
                  <span className="coord lon">
                    {Number(vehicle.longitude).toFixed(4)}°E
                  </span>
                </div>

                <div className="vehicle-status">
                  <div className="speed-box">
                    <span className="speed-value">{Number(vehicle.speed).toFixed(1)}</span>
                    <span className="speed-unit">км/год</span>
                  </div>
                  
                  <div className="status-box">
                    <span className="status-label">
                      {vehicle.speed === 0 ? '🛑 Припинена' : '🚗 Рухається'}
                    </span>
                    <span className="updated-time">
                      {vehicle.updated_seconds_ago ?? 0}с тому
                    </span>
                  </div>
                </div>

                <button 
                  className="btn-detail"
                  onClick={() => openDetail(vehicle)}
                >
                  Деталі
                </button>
              </div>
            </div>
          ))
        ) : (
          <div className="no-vehicles">Немає машин онлайн</div>
        )}
      </div>

      <div className="gps-legend">
        <div className="legend-item">
          <div className="legend-color" style={{ backgroundColor: '#999' }}>◼</div>
          <span>Припинена (0 км/год)</span>
        </div>
        <div className="legend-item">
          <div className="legend-color" style={{ backgroundColor: '#ffc107' }}>▶</div>
          <span>Низька швидкість (&lt;20 км/год)</span>
        </div>
        <div className="legend-item">
          <div className="legend-color" style={{ backgroundColor: '#17a2b8' }}>▶</div>
          <span>Нормальна швидкість (20-60 км/год)</span>
        </div>
        <div className="legend-item">
          <div className="legend-color" style={{ backgroundColor: '#28a745' }}>▶</div>
          <span>Висока швидкість (&gt;60 км/год)</span>
        </div>
      </div>

      {detailVehicle && (
        <div className="gps-modal-overlay" onClick={closeDetail}>
          <div className="gps-modal" onClick={(e) => e.stopPropagation()}>
            <div className="gps-modal-header">
              <h2>🚗 {detailVehicle.plate_number || `Машина #${detailVehicle.vehicle_id}`}</h2>
              <button className="gps-modal-close" onClick={closeDetail}>✕</button>
            </div>
            <div className="gps-modal-body">
              <div className="gps-detail-grid">
                <div className="gps-detail-item">
                  <span className="gps-detail-label">Широта</span>
                  <span className="gps-detail-value">{Number(detailVehicle.latitude).toFixed(6)}°</span>
                </div>
                <div className="gps-detail-item">
                  <span className="gps-detail-label">Довгота</span>
                  <span className="gps-detail-value">{Number(detailVehicle.longitude).toFixed(6)}°</span>
                </div>
                <div className="gps-detail-item">
                  <span className="gps-detail-label">Швидкість</span>
                  <span className="gps-detail-value">{Number(detailVehicle.speed).toFixed(1)} км/год</span>
                </div>
                <div className="gps-detail-item">
                  <span className="gps-detail-label">Курс</span>
                  <span className="gps-detail-value">{Math.round(detailVehicle.heading ?? 0)}°</span>
                </div>
                <div className="gps-detail-item">
                  <span className="gps-detail-label">Останній пінг</span>
                  <span className="gps-detail-value">{new Date(detailVehicle.timestamp).toLocaleString('uk-UA')}</span>
                </div>
                <div className="gps-detail-item">
                  <span className="gps-detail-label">Статус</span>
                  <span className="gps-detail-value">
                    {detailVehicle.speed === 0 ? '🛑 Стоїть' : '🚗 Рухається'}
                  </span>
                </div>
              </div>

              <h3 className="gps-traj-title">📍 Маршрут за останні 24 години</h3>
              {trajectoryLoading ? (
                <p className="gps-traj-empty">Завантаження маршруту...</p>
              ) : trajectory.length === 0 ? (
                <p className="gps-traj-empty">Немає даних маршруту за вказаний період.</p>
              ) : (
                <div className="gps-traj-stats">
                  <div className="gps-detail-item">
                    <span className="gps-detail-label">Точок зафіксовано</span>
                    <span className="gps-detail-value">{trajectory.length}</span>
                  </div>
                  <div className="gps-detail-item">
                    <span className="gps-detail-label">Перша точка</span>
                    <span className="gps-detail-value">
                      {new Date(trajectory[0]?.timestamp).toLocaleString('uk-UA')}
                    </span>
                  </div>
                  <div className="gps-detail-item">
                    <span className="gps-detail-label">Остання точка</span>
                    <span className="gps-detail-value">
                      {new Date(trajectory[trajectory.length - 1]?.timestamp).toLocaleString('uk-UA')}
                    </span>
                  </div>
                </div>
              )}
            </div>
            <div className="gps-modal-footer">
              <button className="btn-refresh" onClick={closeDetail}>Закрити</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default GPSTracking
