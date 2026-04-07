import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { api } from '../api/client'
import './Home.css'

export default function Home() {
  const { user } = useAuth()
  const isVolunteer = user?.role === 'VOLUNTEER'
  
  // Додаємо окремі стейти для волонтерської статистики
  const [stats, setStats] = useState({ 
    resources: 0, 
    categories: 0, 
    pendingRequests: 0,
    // Нові поля для волонтера
    openVolunteerRequests: 0,
    myActiveTasks: 0 
  })

  useEffect(() => {
    if (isVolunteer) {
      // Логіка для волонтера: тягнемо тільки волонтерські запити
      api.volunteerRequests.list() // Припустимо, цей метод повертає всі VolunteerRequest
        .then((data) => {
          const allReqs = Array.isArray(data) ? data : []
          setStats(prev => ({
            ...prev,
            openVolunteerRequests: allReqs.filter(r => r.status === 'OPEN').length,
            myActiveTasks: allReqs.filter(r => r.status === 'TAKEN' && r.taken_by === String(user?.id)).length
          }))
        })
        .catch(console.error)
    } else {
      // Логіка для військових: стандартна складська статистика
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
  }, [isVolunteer, user?.id])

  const canManageUsers = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'PLATOON_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST', 'BRIGADE_STOREKEEPER', 'BATTALION_STOREKEEPER', 'COMPANY_SERGEANT'].includes(user?.role || '')

  return (
    <div className="dashboard">
      <div className="page-header">
        <h1>Головна</h1>
        <p className="page-subtitle">{user?.full_name}</p>
      </div>

      {/* --- СТАТИСТИКА --- */}
      <div className="stats-grid">
        {isVolunteer ? (
          <>
            <Link to="/volunteer-requests" className="stat-card">
              <span className="stat-value">{stats.openVolunteerRequests}</span>
              <span className="stat-label">Відкритих потреб</span>
            </Link>
            <Link to="/volunteer-requests" className="stat-card stat-warning">
              <span className="stat-value">{stats.myActiveTasks}</span>
              <span className="stat-label">Моїх завдань</span>
            </Link>
          </>
        ) : (
          <>
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
          </>
        )}
      </div>

      {/* --- ШВИДКІ ДІЇ --- */}
      <div className="quick-actions">
        <h2>Швидкі дії</h2>
        <div className="actions-grid">
          {/* Військові дії (ховаємо від волонтера) */}
          {!isVolunteer && (
            <>
              <Link to="/inventory" className="action-card">Ресурси</Link>
              <Link to="/requests" className="action-card">Заявки</Link>
            </>
          )}

          {/* Спільна дія, але з різним текстом */}
          <Link to="/volunteer-requests" className="action-card">
            {isVolunteer ? 'Переглянути потреби' : 'Для волонтерів'}
          </Link>

          {/* Управління користувачами (вже має свою перевірку) */}
          {canManageUsers && !isVolunteer && (
            <Link to="/admin/users" className="action-card">
              Користувачі
            </Link>
          )}
        </div>
      </div>
    </div>
  )
}