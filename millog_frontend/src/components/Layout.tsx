import React, { useEffect, useState, useRef } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { api, ROLE_NAMES } from '../api/client' // 🔥 Імпортували ROLE_NAMES
import toast from 'react-hot-toast'
import './Layout.css'

const USER_CREATOR_ROLES = ['ADMIN', 'REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 'TEAM_LEAD', 'REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN', 'REGION_STOREKEEPER', 'BRANCH_STOREKEEPER', 'DEPT_SUPERVISOR']
const UNIT_MANAGER_ROLES = ['ADMIN', 'REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 'REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN', 'REGION_STOREKEEPER', 'BRANCH_STOREKEEPER']

interface LayoutProps {
  children: React.ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const location = useLocation()
  const { user, token, loading, logout } = useAuth()

  // --- Логіка сповіщень ---
  const [pendingCount, setPendingCount] = useState(0)
  const prevCountRef = useRef(0)

  // Перевіряємо, чи має поточний користувач права погоджувати заявки
  const canApproveRequests = user?.role && ['ADMIN', 'REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 'REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN'].includes(user.role)

  useEffect(() => {
    // Запускаємо опитування тільки для авторизованих логістів/менеджерів
    if (!token || !canApproveRequests) return;

    const checkNewRequests = async () => {
      try {
        const reqs = await api.requests.list();
        if (Array.isArray(reqs)) {
          const currentPending = reqs.filter(r => r.status === 'PENDING').length;
          
          if (currentPending > prevCountRef.current && prevCountRef.current !== 0) {
            toast('Надійшла нова заявка на постачання!', { 
              icon: '🔔', 
              duration: 5000,
              style: { border: '1px solid #3b82f6', padding: '16px', color: '#1e3a8a' }
            });
          }
          
          prevCountRef.current = currentPending;
          setPendingCount(currentPending);
        }
      } catch (err) {
        // Тихо ігноруємо помилки мережі в фоні
      }
    };

    checkNewRequests();
    const interval = setInterval(checkNewRequests, 15000); 
    
    return () => clearInterval(interval);
  }, [token, canApproveRequests]);
  // --------------------------------

  if (loading) {
    return (
      <div className="layout-loading">
        <div className="spinner" />
        <p>Завантаження...</p>
      </div>
    )
  }

  const showSidebar = !!token
  const canManageUsers = user?.role && USER_CREATOR_ROLES.includes(user.role)
  const canManageUnits = user?.role && UNIT_MANAGER_ROLES.includes(user.role)
  
  const isCONTRACTOR = user?.role === 'CONTRACTOR'
  // 🔥 НОВЕ: Звичайний співробітник і підрядник не мають доступу до управління логістикою
  const canManageLogistics = !['CONTRACTOR', 'EMPLOYEE'].includes(user?.role || '')

  return (
    <div className={`layout ${showSidebar ? 'with-sidebar' : ''}`}>
      {showSidebar && (
        <aside className="sidebar">
          <div className="sidebar-header">
            <Link to="/" className="sidebar-logo">
              Omnilog
            </Link>
          </div>
          <nav className="sidebar-nav">
            <Link to="/" className={location.pathname === '/' ? 'active' : ''}>
              Головна
            </Link>

            {!isCONTRACTOR && (
              <Link to="/my-equipment" className={location.pathname === '/my-equipment' ? 'active' : ''}>
                Профіль
              </Link>
            )}

            {/* 🔥 Показуємо ці пункти тільки тим, хто має доступ до логістики */}
            {canManageLogistics && (
              <>
                <Link to="/analytics" className={location.pathname === '/analytics' ? 'active' : ''}>
                  Аналітика
                </Link>
                <Link to="/inventory" className={location.pathname === '/inventory' ? 'active' : ''}>
                  Ресурси
                </Link>
                
                <Link 
                  to="/requests" 
                  className={location.pathname === '/requests' ? 'active' : ''}
                  style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}
                >
                  <span>Заявки</span>
                  {pendingCount > 0 && (
                    <span style={{
                      backgroundColor: '#ef4444',
                      color: 'white',
                      fontSize: '11px',
                      fontWeight: 'bold',
                      padding: '2px 8px',
                      borderRadius: '12px',
                      boxShadow: '0 0 8px rgba(239, 68, 68, 0.4)'
                    }}>
                      {pendingCount}
                    </span>
                  )}
                </Link>

                <Link to="/vehicles" className={location.pathname === '/vehicles' ? 'active' : ''}>
                  Автопарк
                </Link>
                <Link to="/warehouses" className={location.pathname === '/warehouses' ? 'active' : ''}>
                  Склади
                </Link>
              </>
            )}

            <Link to="/contractor-requests" className={location.pathname === '/contractor-requests' ? 'active' : ''}>
              {isCONTRACTOR ? 'Відкриті завдання' : 'Заявки підрядникам'}
            </Link>

            {canManageUnits && (
              <Link to="/units" className={location.pathname === '/units' ? 'active' : ''}>
                Оргструктура
              </Link>
            )}
            {canManageUsers && (
              <Link to="/admin/users" className={location.pathname === '/admin/users' ? 'active' : ''}>
                Користувачі
              </Link>
            )}
            {!isCONTRACTOR && (
              <Link to="/kiosk" className={location.pathname === '/kiosk' ? 'active' : ''}>
                📦 Термінал (Каса)
              </Link>
            )}
            
            {user?.role === 'ADMIN' && (
              <Link to="/audit" className={location.pathname === '/audit' ? 'active' : ''}>
                🛡️ Журнал аудиту
              </Link>
            )}

            {!isCONTRACTOR && (
              <Link to="/billing" className={location.pathname === '/billing' ? 'active' : ''}>
                💎 Тарифні плани
              </Link>
            )}

          </nav>
          <div className="sidebar-footer">
            <div className="user-badge">
              <span className="user-name">{user?.full_name}</span>
              <span className="user-role">{user?.role ? ROLE_NAMES[user.role] || user.role : ''}</span>
            </div>
            <button className="btn btn-secondary btn-sm btn-block" onClick={logout}>
              Вийти
            </button>
          </div>
        </aside>
      )}

      <div className="main-content">
        {!showSidebar && (
          <header className="top-bar">
            <Link to="/" className="logo-text">Omnilog</Link>
            <Link to="/login" className="btn btn-primary">Увійти</Link>
          </header>
        )}
        <main className="page-content">{children}</main>
      </div>
    </div>
  )
}