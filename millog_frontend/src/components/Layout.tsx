import React, { useEffect, useState, useRef } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { api, ROLE_NAMES } from '../api/client'
import toast from 'react-hot-toast'
import { ROLE_GROUPS, ROLES, hasRole } from '../constants/roles'
import NotificationCenter from './NotificationCenter'
import './Layout.css'

interface LayoutProps {
  children: React.ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const location = useLocation()
  const { user, token, loading, logout } = useAuth()

  // --- Логіка сповіщень ---
  const [pendingCount, setPendingCount] = useState(0)
  const prevCountRef = useRef(0)

  const userRole = user?.role || ''
  const isCONTRACTOR = userRole === ROLES.CONTRACTOR
  
  // Перевіряємо, чи має поточний користувач права погоджувати заявки
  const canApproveRequests = hasRole(userRole, ROLE_GROUPS.approvers)

  useEffect(() => {
    if (!token || !canApproveRequests) return

    const checkNewRequests = async () => {
      try {
        const reqs = await api.requests.list()
        // Безпечна перевірка через any на випадок нюансів TypeScript
        if (Array.isArray(reqs)) {
          const currentPending = (reqs as any[]).filter(r => r?.status === 'PENDING').length
          
          if (currentPending > prevCountRef.current && prevCountRef.current !== 0) {
            toast('Надійшла нова заявка на постачання!', { 
              icon: '🔔', 
              duration: 5000,
              style: { border: '1px solid #3b82f6', padding: '16px', color: '#1e3a8a' }
            })
          }
          
          prevCountRef.current = currentPending
          setPendingCount(currentPending)
        }
      } catch (err) {
        // Тихо ігноруємо помилки мережі в фоні
      }
    }

    checkNewRequests()
    const interval = setInterval(checkNewRequests, 15000) 
    
    return () => clearInterval(interval)
  }, [token, canApproveRequests])

  if (loading) {
    return (
      <div className="layout-loading">
        <div className="spinner" />
        <p>Завантаження Omnilog...</p>
      </div>
    )
  }

  const showSidebar = !!token

  // --- ГРУПИ ДОСТУПУ (синхронізовано з маршрутами у App.tsx) ---
  const canSeeAnalytics = hasRole(userRole, ROLE_GROUPS.analytics)
  const canSeeInventory = hasRole(userRole, ROLE_GROUPS.inventory)
  const canSeeTransport = hasRole(userRole, ROLE_GROUPS.transport)
  const canSeeRequests = hasRole(userRole, ROLE_GROUPS.requests)

  const canManageUnits = hasRole(userRole, ROLE_GROUPS.units)
  const canManageUsers = hasRole(userRole, ROLE_GROUPS.users)

  const canSeeKiosk = hasRole(userRole, ROLE_GROUPS.kiosk)
  const canManageContracts = hasRole(userRole, ROLE_GROUPS.contracts)

  const isAdmin = userRole === ROLES.ADMIN
  const isSystemAdmin = userRole === ROLES.SYSTEM_ADMIN

  return (
    <div className={`layout ${showSidebar ? 'with-sidebar' : ''}`}>
      {showSidebar && (
        <aside className="sidebar">
          <div className="sidebar-header">
            <Link to="/" className="sidebar-logo">
              Omnilog
            </Link>
            <NotificationCenter />
          </div>
          <nav className="sidebar-nav">
            
            {/* 1. Спільні розділи */}
            <Link to="/" className={location.pathname === '/' ? 'active' : ''}>
              Головна
            </Link>

            {!isCONTRACTOR && (
              <Link to="/my-equipment" className={location.pathname === '/my-equipment' ? 'active' : ''}>
                Профіль
              </Link>
            )}

            {/* 2. Аналітика */}
            {canSeeAnalytics && (
              <>
                <Link to="/analytics" className={location.pathname === '/analytics' ? 'active' : ''}>
                  Аналітика
                </Link>
                <Link to="/kpi" className={location.pathname === '/kpi' ? 'active' : ''}>
                  📊 KPI (SLA / TCO / Ризики)
                </Link>
              </>
            )}

            {/* 3. Складський Блок (Майно та Склади) */}
            {canSeeInventory && (
              <>
                <Link to="/inventory" className={location.pathname === '/inventory' ? 'active' : ''}>
                  Ресурси
                </Link>
                <Link to="/warehouses" className={location.pathname === '/warehouses' ? 'active' : ''}>
                  Склади
                </Link>
              </>
            )}

            {/* 4. Транспортний блок */}
            {canSeeTransport && (
              <>
                <Link to="/vehicles" className={location.pathname === '/vehicles' ? 'active' : ''}>
                  Автопарк
                </Link>
                <Link to="/gps" className={location.pathname === '/gps' ? 'active' : ''}>
                  🌍 GPS Трекінг
                </Link>
                <Link to="/maintenance" className={location.pathname === '/maintenance' ? 'active' : ''}>
                  🔮 Графік ТО
                </Link>
                <Link to="/fuel-anomalies" className={location.pathname === '/fuel-anomalies' ? 'active' : ''}>
                  🛡️ Антифрод пального
                </Link>
              </>
            )}

            {/* 5. Заявки */}
            {canSeeRequests && (
              <Link 
                to="/requests" 
                className={location.pathname === '/requests' ? 'active' : ''}
                style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}
              >
                <span>Заявки</span>
                {pendingCount > 0 && (
                  <span style={{
                    backgroundColor: '#ef4444', color: 'white', fontSize: '11px', fontWeight: 'bold',
                    padding: '2px 8px', borderRadius: '12px', boxShadow: '0 0 8px rgba(239, 68, 68, 0.4)'
                  }}>
                    {pendingCount}
                  </span>
                )}
              </Link>
            )}

            {/* Заявки підрядникам бачать або самі підрядники, або менеджери */}
            {(isCONTRACTOR || canManageContracts) && (
              <Link to="/contractor-requests" className={location.pathname === '/contractor-requests' ? 'active' : ''}>
                {isCONTRACTOR ? 'Відкриті завдання' : 'Заявки підрядникам'}
              </Link>
            )}

            {/* 6. Управління Організацією */}
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

            {/* 7. Термінал / Каса */}
            {canSeeKiosk && (
              <Link to="/kiosk" className={location.pathname === '/kiosk' ? 'active' : ''}>
                📦 Термінал (Каса)
              </Link>
            )}
            
            {/* 8. Права Суперадміна */}
            {isAdmin && (
              <>
                <Link to="/audit" className={location.pathname === '/audit' ? 'active' : ''}>
                  🛡️ Журнал аудиту
                </Link>
              </>
            )}

            {/* 9. Платформа (тільки власник SaaS) */}
            {isSystemAdmin && (
              <Link to="/platform" className={location.pathname === '/platform' ? 'active' : ''}>
                🌐 Платформа (Tenants)
              </Link>
            )}

            {/* Тарифні плани — бачать ті, хто бачить аналітику (керівники) */}
            {canSeeAnalytics && (
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