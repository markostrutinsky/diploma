import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { api } from '../api/client'
import './Home.css'

export default function Home() {
  const { user } = useAuth()
  const [stats, setStats] = useState({ resources: 0, categories: 0, pendingRequests: 0 })

  useEffect(() => {
    Promise.all([api.inventory.listResources(undefined), api.inventory.listCategories(), api.requests.list()])
      .then(([res, cats, reqs]) => {
        const safeRes = Array.isArray(res) ? res : []
        const safeCats = Array.isArray(cats) ? cats : []
        const safeReqs = Array.isArray(reqs) ? reqs : []
        setStats({
          resources: safeRes.length,
          categories: safeCats.length,
          pendingRequests: safeReqs.filter((r) => r.status === 'PENDING').length,
        })
      })
      .catch(() => {})
  }, [])

  return (
    <div className="dashboard">
      <div className="page-header">
        <h1>Головна</h1>
        <p className="page-subtitle">{user?.full_name}</p>
      </div>

      <div className="stats-grid">
        <Link to="/inventory" className="stat-card">
          <span className="stat-value">{stats.resources}</span>
          <span className="stat-label">Ресурсів</span>
        </Link>
        <Link to="/inventory" className="stat-card">
          <span className="stat-value">{stats.categories}</span>
          <span className="stat-label">Категорій</span>
        </Link>
        <Link to="/requests" className="stat-card stat-warning">
          <span className="stat-value">{stats.pendingRequests}</span>
          <span className="stat-label">Заявок очікує</span>
        </Link>
      </div>

      <div className="quick-actions">
        <h2>Швидкі дії</h2>
        <div className="actions-grid">
          <Link to="/inventory" className="action-card">
            Ресурси
          </Link>
          <Link to="/requests" className="action-card">
            Заявки
          </Link>
          <Link to="/volunteer-requests" className="action-card">
            Для волонтерів
          </Link>
          {['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'PLATOON_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST', 'BRIGADE_STOREKEEPER', 'BATTALION_STOREKEEPER', 'COMPANY_SERGEANT'].includes(user?.role || '') && (
            <Link to="/admin/users" className="action-card">
              Користувачі
            </Link>
          )}
        </div>
      </div>
    </div>
  )
}
