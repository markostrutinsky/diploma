import React, { useEffect, useState, useRef } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { api, ROLE_NAMES } from '../api/client'
import toast from 'react-hot-toast'
import { ROLE_GROUPS, ROLES, hasRole } from '../constants/roles'
import NotificationCenter from './NotificationCenter'
import { SubscriptionExpiryBanner } from './FeatureGate'
import './Layout.css'

interface LayoutProps {
  children: React.ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const location = useLocation()
  const { user, token, loading, logout, supportTenant, exitSupportTenant } = useAuth()

  // --- Логіка сповіщень ---
  const [pendingCount, setPendingCount] = useState(0)
  const prevCountRef = useRef(0)

  const userRole = user?.role || ''
  const isCONTRACTOR = userRole === ROLES.CONTRACTOR
  const isAdmin = userRole === ROLES.ADMIN
  const isSystemAdmin = userRole === ROLES.SYSTEM_ADMIN
  const hasTenantContext = !isSystemAdmin || !!supportTenant
  
  // Перевіряємо, чи має поточний користувач права погоджувати заявки
  const canApproveRequests = hasTenantContext && hasRole(userRole, ROLE_GROUPS.approvers)

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

  const PUBLIC_ROUTES = ['/login', '/signup', '/register', '/bootstrap', '/setup-password']
  const showSidebar = !!token && !PUBLIC_ROUTES.includes(location.pathname)

  // --- ГРУПИ ДОСТУПУ (синхронізовано з маршрутами у App.tsx) ---
  const canSeeAnalytics = hasTenantContext && hasRole(userRole, ROLE_GROUPS.analytics)
  const canSeeInventory = hasTenantContext && hasRole(userRole, ROLE_GROUPS.inventory)
  const canSeeTransport = hasTenantContext && hasRole(userRole, ROLE_GROUPS.transport)
  const canSeeRequests = hasTenantContext && hasRole(userRole, ROLE_GROUPS.requests)

  const canManageUnits = hasTenantContext && hasRole(userRole, ROLE_GROUPS.units)
  const canManageUsers = hasTenantContext && hasRole(userRole, ROLE_GROUPS.users)

  const canSeeKiosk = hasTenantContext && hasRole(userRole, ROLE_GROUPS.kiosk)
  const canManageContracts = hasTenantContext && hasRole(userRole, ROLE_GROUPS.contracts)

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
          {isSystemAdmin && (
            <div className={`support-mode-banner ${supportTenant ? 'active' : ''}`}>
              <span>{supportTenant ? `Support mode: ${supportTenant.name}` : 'Platform mode'}</span>
              {supportTenant && (
                <button type="button" onClick={exitSupportTenant}>
                  Вийти
                </button>
              )}
            </div>
          )}
          <nav className="sidebar-nav">
            
            {/* 1. Спільні розділи */}
            {hasTenantContext && (
              <Link to="/" className={location.pathname === '/' ? 'active' : ''}>
                Головна
              </Link>
            )}

            {!isCONTRACTOR && hasTenantContext && (
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
                  Показники ефективності
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
                  GPS Трекінг
                </Link>
                <Link to="/maintenance" className={location.pathname === '/maintenance' ? 'active' : ''}>
                  Графік ТО
                </Link>
                <Link to="/fuel-anomalies" className={location.pathname === '/fuel-anomalies' ? 'active' : ''}>
                  Контроль пального
                </Link>
              </>
            )}


            {/* 5. Заявки */}
            {canSeeRequests && (
              <Link 
                to="/requests" 
                className={`sidebar-nav-link-with-badge ${location.pathname === '/requests' ? 'active' : ''}`}
              >
                <span className="sidebar-link-label">Заявки</span>
                {pendingCount > 0 && (
                  <span className="sidebar-link-badge">
                    {pendingCount}
                  </span>
                )}
              </Link>
            )}

            {/* 5.1 Мої Рейси - тільки для співробітників (водіїв та складських) */}
            {userRole === 'EMPLOYEE' && (
              <Link 
                to="/my-shipments" 
                className={location.pathname === '/my-shipments' ? 'active' : ''}
              >
                Мої Рейси
              </Link>
            )}

            {/* Заявки підрядникам бачать або самі підрядники, або менеджери */}
            {(isCONTRACTOR || canManageContracts) && (
              <Link to="/contractor-requests" className={location.pathname === '/contractor-requests' ? 'active' : ''}>
                {isCONTRACTOR ? 'Відкриті завдання' : 'Заявки підрядникам'}
              </Link>
            )}

            {/* Схвалення підрядників (тільки менеджери, не самі підрядники) */}
            {canManageContracts && (
              <Link to="/admin/contractors" className={location.pathname === '/admin/contractors' ? 'active' : ''}>
                Підрядники
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
                Термінал видачі майна
              </Link>
            )}
            
            {/* 8. Права Суперадміна */}
            {(isAdmin || isSystemAdmin || hasRole(userRole, ROLE_GROUPS.superAdmin)) && (
              <>
                <Link to="/audit" className={location.pathname === '/audit' ? 'active' : ''}>
                  Журнал аудиту
                </Link>
              </>
            )}

            {/* 9. Платформа (тільки власник SaaS) */}
            {isSystemAdmin && (
              <Link to="/platform" className={location.pathname === '/platform' ? 'active' : ''}>
                Адміністрування платформи
              </Link>
            )}

            {/* Тарифні плани — бачать ті, хто бачить аналітику (керівники) */}
            {canSeeAnalytics && (
              <Link to="/billing" className={location.pathname === '/billing' ? 'active' : ''}>
                Тарифні плани
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

      <div className={`main-content ${showSidebar ? 'main-content--app' : 'main-content--public'}`}>
        {!showSidebar && (
          <header className="top-bar">
            <Link to="/" className="logo-text">Omnilog</Link>
            <Link to="/login" className="btn btn-primary">Увійти</Link>
          </header>
        )}
        <main className={`page-content ${showSidebar ? 'page-content--app' : 'page-content--public'}`}>
          <SubscriptionExpiryBanner />
          <div
            key={location.pathname}
            className={`route-transition ${showSidebar ? 'route-transition--app' : 'route-transition--public'}`}
          >
            {children}
          </div>
        </main>
      </div>
    </div>
  )
}
