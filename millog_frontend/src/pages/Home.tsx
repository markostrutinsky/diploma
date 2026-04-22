import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { usePermissions } from '../hooks/usePermissions'
import { api } from '../api/client'
import './Home.css'

export default function Home() {
  const { user } = useAuth()
  const perms = usePermissions()
  const isCONTRACTOR = perms.isContractor
  const isAdmin = perms.isAdmin
  const isManager = perms.isManager
  
  const [stats, setStats] = useState({ 
    resources: 0, 
    categories: 0, 
    pendingRequests: 0,
    openContractorRequests: 0,
    myActiveTasks: 0 
  })

  // Стейт для серйозного віджета
  const [unitMetrics, setUnitMetrics] = useState<any[]>([])
  const [recentLogs, setRecentLogs] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [loadingLogs, setLoadingLogs] = useState(true)

  useEffect(() => {
    setLoading(true)
    
    if (isCONTRACTOR) {
      api.contractorRequests.list() 
        .then((data) => {
          const allReqs = Array.isArray(data) ? data : []
          setStats(prev => ({
            ...prev,
            openContractorRequests: allReqs.filter(r => r.status === 'OPEN').length,
            myActiveTasks: allReqs.filter(r => r.status === 'TAKEN' && r.taken_by === String(user?.id)).length
          }))
        })
        .catch(console.error)
        .finally(() => setLoading(false))
    } else {
      Promise.all([
        api.inventory.listResources(undefined).catch(() => []),
        api.inventory.listCategories().catch(() => []),
        api.requests.list().catch(() => []),
        api.units.list().catch(() => [])
      ])
        .then(([res, cats, reqs, units]) => {
          const safeRes = Array.isArray(res) ? res : []
          const safeCats = Array.isArray(cats) ? cats : []
          const safeReqs = Array.isArray(reqs) ? reqs : []
          const safeUnits = Array.isArray(units) ? units : []

          setStats(prev => ({
            ...prev,
            resources: safeRes.length,
            categories: safeCats.length,
            pendingRequests: safeReqs.filter((r) => r.status === 'PENDING').length,
          }))

          // ---------------------------------------------------------
          // 🧠 СЕРЙОЗНА БІЗНЕС-ЛОГІКА: Розрахунок Unit Health Score
          // ---------------------------------------------------------
          const metricsData = safeUnits.slice(0, 4).map(unit => {
            // Використовуємо (r: any) та (req: any), щоб TypeScript не сварився на відсутні поля
            const unitResources = safeRes.filter((r: any) => r.unit_id === unit.id)
            const unitRequests = safeReqs.filter((req: any) => req.unit_id === unit.id && req.status === 'PENDING')
            
            const totalItems = unitResources.length
            const pendingCount = unitRequests.length

            // Стан активів. Якщо status немає, вважаємо, що майно справне (ACTIVE)
            const activeItems = unitResources.filter((r: any) => !r.status || r.status === 'ACTIVE').length
            const healthScore = totalItems > 0 ? (activeItems / totalItems) * 100 : 100

            // Операційне навантаження
            const opsScore = Math.max(0, 100 - (pendingCount * 5))

            // Фінансова вартість (якщо поля price немає, воно просто додасть 0)
            const totalValue = unitResources.reduce((sum: number, r: any) => sum + (Number(r.price) || 0), 0)

            // Композитний індекс
            const overallScore = Math.round((healthScore * 0.6) + (opsScore * 0.4))

            return { 
              name: unit.name, 
              overallScore, 
              healthScore: Math.round(healthScore),
              pendingCount,
              totalItems,
              totalValue
            }
          })
          
          // Сортуємо відділи від найбільш проблемних до найефективніших
          metricsData.sort((a, b) => a.overallScore - b.overallScore)
          
          setUnitMetrics(metricsData)
        })
        .catch(console.error)
        .finally(() => setLoading(false))
    }

    if (isAdmin) {
      setLoadingLogs(true)
      api.admin.getAuditLogs()
        .then((logs: any) => { 
          setRecentLogs(Array.isArray(logs) ? logs.slice(0, 5) : []) 
        })
        .catch(console.error)
        .finally(() => setLoadingLogs(false))
    }
  }, [isCONTRACTOR, isAdmin, user?.id])

  const formatTime = (dateString: string) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    return date.toLocaleTimeString('uk-UA', { hour: '2-digit', minute: '2-digit' })
  }

  // Функція для форматування валюти
  const formatCurrency = (val: number) => {
    if (val === 0) return 'Н/Д'
    return new Intl.NumberFormat('uk-UA', { style: 'currency', currency: 'UAH', maximumFractionDigits: 0 }).format(val)
  }

  return (
    <div className="dashboard-container">
      <div className="page-header">
        <h1>Головна</h1>
        <p className="page-subtitle">{user?.full_name} ({user?.role})</p>
      </div>

      <div className="stats-grid">
        {/* ... (код верхніх карток залишається без змін) ... */}
        {isCONTRACTOR ? (
          <>
            <Link to="/contractor-requests" className="stat-card">
              <span className="stat-value">{stats.openContractorRequests}</span>
              <span className="stat-label">Доступних завдань</span>
            </Link>
            <Link to="/contractor-requests" className="stat-card stat-warning">
              <span className="stat-value">{stats.myActiveTasks}</span>
              <span className="stat-label">Моїх завдань</span>
            </Link>
          </>
        ) : (
          <>
            <Link to="/inventory" className="stat-card">
              <span className="stat-value">{stats.resources}</span>
              <span className="stat-label">Ресурсів на балансі</span>
            </Link>
            <Link to="/inventory" className="stat-card">
              <span className="stat-value">{stats.categories}</span>
              <span className="stat-label">Категорій</span>
            </Link>
            <Link to="/requests" className="stat-card stat-warning">
              <span className="stat-value">{stats.pendingRequests}</span>
              <span className="stat-label">Заявок очікує</span>
            </Link>
          </>
        )}
      </div>

      <div className="quick-actions-section">
        {/* ... (код кнопок швидких дій залишається без змін) ... */}
        <h2>Швидкі дії</h2>
        <div className="actions-row">
          {!isCONTRACTOR && (
            <>
              <Link to="/inventory" className="action-button">📦 Ресурси</Link>
              <Link to="/warehouses" className="action-button">🏢 Склади</Link>
              <Link to="/requests" className="action-button">📝 Внутрішні заявки</Link>
              <Link to="/fleet" className="action-button">🚙 Автопарк</Link>
            </>
          )}
          <Link to="/contractor-requests" className="action-button">
            🤝 {isCONTRACTOR ? 'Доступні завдання' : 'Заявки підрядникам'}
          </Link>
          {isManager && !isCONTRACTOR && (
            <>
              <Link to="/analytics" className="action-button">📊 Аналітика</Link>
              <Link to="/admin/units" className="action-button">🏢 Оргструктура</Link>
              <Link to="/admin/users" className="action-button">👥 Користувачі</Link>
            </>
          )}
        </div>
      </div>

      <div className="widgets-grid">
        <div className="widgets-main-column">
          
          {/* НОВИЙ ВІДЖЕТ: Аналітика відділів */}
          {!isCONTRACTOR && (
            <div className="erp-widget">
              <div className="widget-header">
                <h3>Зведення: Індекс здоров'я відділів (OEE)</h3>
                <span className="widget-badge">Аналітика</span>
              </div>
              <div className="metrics-list">
                {loading ? <p className="empty-text">Розрахунок показників...</p> : unitMetrics.length > 0 ? (
                  unitMetrics.map((unit, i) => (
                    <div key={i} className="metric-card">
                      <div className="metric-header">
                        <span className="metric-name">{unit.name}</span>
                        <div className="metric-score-wrap">
                           <span className="score-label">Оцінка:</span>
                           <span className={`score-value ${unit.overallScore > 80 ? 'text-success' : unit.overallScore > 50 ? 'text-warning' : 'text-danger'}`}>
                             {unit.overallScore}%
                           </span>
                        </div>
                      </div>
                      
                      <div className="metric-details">
                        <div className="detail-col">
                          <span className="detail-val">{unit.healthScore}%</span>
                          <span className="detail-lbl">Справність майна</span>
                        </div>
                        <div className="detail-col">
                          <span className={`detail-val ${unit.pendingCount > 0 ? 'text-warning' : ''}`}>{unit.pendingCount}</span>
                          <span className="detail-lbl">Заявок в черзі</span>
                        </div>
                        <div className="detail-col">
                          <span className="detail-val">{unit.totalItems}</span>
                          <span className="detail-lbl">К-сть активів</span>
                        </div>
                        <div className="detail-col right-align">
                          <span className="detail-val">{formatCurrency(unit.totalValue)}</span>
                          <span className="detail-lbl">Оціночна вартість</span>
                        </div>
                      </div>
                    </div>
                  ))
                ) : <p className="empty-text">Дані про відділи відсутні.</p>}
              </div>
            </div>
          )}

          {/* ... (Журнал аудиту залишається без змін) ... */}
          {isAdmin && (
            <div className="erp-widget">
              <div className="widget-header">
                <h3>Останні події в системі</h3>
                <Link to="/audit" className="widget-link">Всі логи</Link>
              </div>
              <div className="feed-list">
                {loadingLogs ? (
                  <p className="empty-text">Завантаження подій...</p>
                ) : recentLogs.length > 0 ? (
                  recentLogs.map((log, idx) => (
                    <div key={log.id || idx} className="feed-item">
                      <div className="feed-time">{formatTime(log.created_at)}</div>
                      <div className="feed-content">
                        {/* Використовуємо user_email замість user_id */}
                        <strong>{log.user_email || 'Система'}</strong>
                        <span>
                          <span style={{fontWeight: 600, color: '#3b82f6', marginRight: '6px'}}>
                            {/* Використовуємо action_type */}
                            [{log.action_type || 'SYS'}]
                          </span> 
                          {/* Використовуємо details */}
                          {log.details || 'Деталі події відсутні'}
                        </span>
                      </div>
                    </div>
                  ))
                ) : (
                  <p className="empty-text">Немає останніх подій.</p>
                )}
              </div>
            </div>
          )}

        </div>

        {/* ... (Права колонка зі сповіщеннями залишається без змін) ... */}
        <div className="widgets-side-column">
          <div className="erp-widget">
             <div className="widget-header">
               <h3>Потребує уваги</h3>
             </div>
             <div className="alerts-container">
               {!isCONTRACTOR && stats.pendingRequests > 0 && (
                 <Link to="/requests" className="alert-card warning-alert">
                   <div className="alert-icon">⚠️</div>
                   <div className="alert-content">
                     <strong>{stats.pendingRequests} нових заявок</strong>
                     <span>Чекають вашого погодження</span>
                   </div>
                 </Link>
               )}
               {isCONTRACTOR && stats.openContractorRequests > 0 && (
                 <Link to="/contractor-requests" className="alert-card info-alert">
                   <div className="alert-icon">📦</div>
                   <div className="alert-content">
                     <strong>Нові завдання</strong>
                     <span>Доступно {stats.openContractorRequests} завдань для виконання</span>
                   </div>
                 </Link>
               )}
               {((!isCONTRACTOR && stats.pendingRequests === 0) || (isCONTRACTOR && stats.openContractorRequests === 0)) && (
                 <div className="alert-card success-alert">
                   <div className="alert-icon">✅</div>
                   <div className="alert-content">
                     <strong>Все опрацьовано</strong>
                     <span>Нових критичних запитів немає</span>
                   </div>
                 </div>
               )}
             </div>
          </div>
        </div>

      </div>
    </div>
  )
}