import React, { useEffect, useRef, useState } from 'react'
import { MapContainer, TileLayer, Marker, Popup, Polyline, useMap } from 'react-leaflet'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { api } from '../api/client'
import { usePermissions } from '../hooks/usePermissions'
import { PaywallScreen } from '../components/FeatureGate'
import './GPSTracking.css'

// Leaflet default icon fix (Vite/webpack asset issue)
delete (L.Icon.Default.prototype as any)._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  iconUrl:       'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  shadowUrl:     'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
})

const makeVehicleIcon = (speed: number, heading: number) => {
  const color = speed === 0 ? '#6b7280' : speed < 20 ? '#f59e0b' : speed < 60 ? '#3b82f6' : '#10b981'
  const arrow = speed > 0 ? `<div style="transform:rotate(${heading}deg);font-size:14px;line-height:1;">➤</div>` : `<div style="font-size:14px;">◼</div>`
  return L.divIcon({
    className: '',
    html: `<div style="
      background:${color};color:#fff;border-radius:50%;
      width:36px;height:36px;display:flex;align-items:center;
      justify-content:center;border:2px solid #fff;
      box-shadow:0 2px 8px rgba(0,0,0,0.35);">${arrow}</div>`,
    iconSize: [36, 36],
    iconAnchor: [18, 18],
    popupAnchor: [0, -20],
  })
}

// Компонент що переміщує камеру до всіх маркерів
const FitBounds: React.FC<{ vehicles: Vehicle[] }> = ({ vehicles }) => {
  const map = useMap()
  const fitted = useRef(false)
  useEffect(() => {
    if (!fitted.current && vehicles.length > 0) {
      const bounds = L.latLngBounds(vehicles.map(v => [v.latitude, v.longitude]))
      map.fitBounds(bounds, { padding: [50, 50], maxZoom: 13 })
      fitted.current = true
    }
  }, [vehicles])
  return null
}

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
  const perms = usePermissions()
  const hasAccess = perms.hasFeature('gps_tracking')
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
    if (!hasAccess) {
      setLoading(false)
      return
    }
    fetchFleetData()

    if (autoRefresh) {
      const interval = setInterval(fetchFleetData, refreshInterval * 1000)
      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval, hasAccess])

  const fetchFleetData = async () => {
    try {
      const response = await api.gps.getFleetMap()
      setFleetData(response)
      setError(null)
    } catch (err: any) {
      const message = err.message || ''
      if (message.includes('402')) {
        // Тариф не підтримує — показуємо paywall
        setError(null)
      } else {
        setError('Помилка при завантаженні GPS даних')
      }
    } finally {
      setLoading(false)
    }
  }

  if (!hasAccess) {
    return (
      <PaywallScreen
        feature="gps_tracking"
        description="Real-time карта, історія маршрутів за 24 години, геозони та контроль швидкості — щоб диспетчер бачив усі машини в одному вікні."
      />
    )
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

      {/* ── Головна карта флоту ── */}
      <div className="fleet-map-wrapper">
        <MapContainer
          center={[50.4501, 30.5234]}
          zoom={10}
          style={{ height: '100%', width: '100%', borderRadius: '12px' }}
        >
          <TileLayer
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
          />
          {(fleetData.vehicles ?? []).map(vehicle => (
            <Marker
              key={vehicle.vehicle_id}
              position={[vehicle.latitude, vehicle.longitude]}
              icon={makeVehicleIcon(vehicle.speed, vehicle.heading ?? 0)}
              eventHandlers={{ click: () => openDetail(vehicle) }}
            >
              <Popup>
                <div className="fleet-popup">
                  <strong>{vehicle.plate_number || `Авто №${vehicle.vehicle_id}`}</strong>
                  <div>🏎 {Number(vehicle.speed).toFixed(1)} км/год</div>
                  <div>📍 {Number(vehicle.latitude).toFixed(4)}°, {Number(vehicle.longitude).toFixed(4)}°</div>
                  <div>🕐 {new Date(vehicle.timestamp).toLocaleTimeString('uk-UA')}</div>
                  <button className="popup-detail-btn" onClick={() => openDetail(vehicle)}>
                    Детальніше →
                  </button>
                </div>
              </Popup>
            </Marker>
          ))}
          {(fleetData.vehicles?.length ?? 0) > 0 && (
            <FitBounds vehicles={fleetData.vehicles} />
          )}
        </MapContainer>
      </div>

      {(fleetData.vehicles?.length ?? 0) === 0 && (
        <div className="no-vehicles">🚫 Немає машин у дорозі зараз</div>
      )}

      <div className="gps-legend">
        <div className="legend-item"><span className="legend-dot" style={{ background: '#6b7280' }} />Стоїть</div>
        <div className="legend-item"><span className="legend-dot" style={{ background: '#f59e0b' }} />&lt;20 км/год</div>
        <div className="legend-item"><span className="legend-dot" style={{ background: '#3b82f6' }} />20–60 км/год</div>
        <div className="legend-item"><span className="legend-dot" style={{ background: '#10b981' }} />&gt;60 км/год</div>
      </div>

      {/* ── Модалка з деталями + траса ── */}
      {detailVehicle && (
        <div className="gps-modal-overlay" onClick={closeDetail}>
          <div className="gps-modal gps-modal--wide" onClick={(e) => e.stopPropagation()}>
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
                  <span className="gps-detail-value">{detailVehicle.speed === 0 ? '🛑 Стоїть' : '🚗 Рухається'}</span>
                </div>
              </div>

              <h3 className="gps-traj-title">📍 Маршрут за останні 24 години</h3>

              {trajectoryLoading ? (
                <p className="gps-traj-empty">Завантаження маршруту...</p>
              ) : trajectory.length === 0 ? (
                <p className="gps-traj-empty">Немає даних маршруту за вказаний період.</p>
              ) : (
                <>
                  <div className="gps-traj-stats">
                    <div className="gps-detail-item">
                      <span className="gps-detail-label">Точок зафіксовано</span>
                      <span className="gps-detail-value">{trajectory.length}</span>
                    </div>
                    <div className="gps-detail-item">
                      <span className="gps-detail-label">Перша точка</span>
                      <span className="gps-detail-value">{new Date(trajectory[0]?.timestamp).toLocaleString('uk-UA')}</span>
                    </div>
                    <div className="gps-detail-item">
                      <span className="gps-detail-label">Остання точка</span>
                      <span className="gps-detail-value">{new Date(trajectory[trajectory.length - 1]?.timestamp).toLocaleString('uk-UA')}</span>
                    </div>
                  </div>

                  {/* Міні-карта маршруту */}
                  <div className="traj-map-wrapper">
                    <MapContainer
                      center={[detailVehicle.latitude, detailVehicle.longitude]}
                      zoom={12}
                      style={{ height: '100%', width: '100%', borderRadius: '8px' }}
                    >
                      <TileLayer
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                        attribution='&copy; OpenStreetMap'
                      />
                      <Polyline
                        positions={trajectory.map((p: any) => [p.latitude, p.longitude])}
                        color="#3b82f6"
                        weight={3}
                        opacity={0.8}
                      />
                      <Marker position={[trajectory[0]?.latitude, trajectory[0]?.longitude]}>
                        <Popup>🚀 Старт: {new Date(trajectory[0]?.timestamp).toLocaleTimeString('uk-UA')}</Popup>
                      </Marker>
                      <Marker position={[detailVehicle.latitude, detailVehicle.longitude]}
                        icon={makeVehicleIcon(detailVehicle.speed, detailVehicle.heading ?? 0)}>
                        <Popup>📍 Поточна позиція</Popup>
                      </Marker>
                    </MapContainer>
                  </div>
                </>
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
