import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { api } from '../api/client'
import './Home.css'

export default function Home() {
  const { user } = useAuth()
  const isCONTRACTOR = user?.role === 'CONTRACTOR'
  
  // Додаємо окремі стейти для статистики підрядника
  const [stats, setStats] = useState({ 
    resources: 0, 
    categories: 0, 
    pendingRequests: 0,
    // Нові поля для підрядника
    openContractorRequests: 0,
    myActiveTasks: 0 
  })

  useEffect(() => {
    if (isCONTRACTOR) {
      // Логіка для підрядника: тягнемо тільки доступні зовнішні завдання
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
    } else {
      // Логіка для персоналу компанії: стандартна складська статистика
      Promise.all([
        api.inventory.listResources(undefined), 
        api.inventory.listCategories(), 
        api.requests.list()
      ])
        .then(([res, cats, reqs]) => {
          const safeRes = Array.isArray(res) ? res : []
          const safeCats = Array.isArray(cats) ? cats : []
          const safeReqs = Array.isArray(reqs) ? reqs : []
          setStats(prev => ({
            ...prev,
            resources: safeRes.length,
            categories: safeCats.length,
            pendingRequests: safeReqs.filter((r) => r.status === 'PENDING').length,
          }))
        })
        .catch(() => {})
    }
  }, [isCONTRACTOR, user?.id])

  const canManageUsers = ['ADMIN', 'REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 'TEAM_LEAD', 'REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN', 'REGION_STOREKEEPER', 'BRANCH_STOREKEEPER', 'DEPT_SUPERVISOR'].includes(user?.role || '')

  return (
    <div className="dashboard">
      <div className="page-header">
        <h1>Головна</h1>
        <p className="page-subtitle">{user?.full_name}</p>
      </div>

      {/* --- СТАТИСТИКА --- */}
      <div className="stats-grid">
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

      {/* --- ШВИДКІ ДІЇ --- */}
      <div className="quick-actions">
        <h2>Швидкі дії</h2>
        <div className="actions-grid">
          {/* Внутрішні дії (ховаємо від зовнішнього підрядника) */}
          {!isCONTRACTOR && (
            <>
              <Link to="/inventory" className="action-card">Ресурси</Link>
              <Link to="/requests" className="action-card">Заявки</Link>
            </>
          )}

          {/* Спільна дія, але з різним текстом */}
          <Link to="/contractor-requests" className="action-card">
            {isCONTRACTOR ? 'Відкриті завдання' : 'Заявки підрядникам'}
          </Link>

          {/* Управління користувачами (вже має свою перевірку) */}
          {canManageUsers && !isCONTRACTOR && (
            <Link to="/admin/users" className="action-card">
              Користувачі
            </Link>
          )}
        </div>
      </div>
    </div>
  )
}