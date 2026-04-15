import React, { useEffect, useState, useRef } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { api } from '../api/client' // ДОДАНО
import toast from 'react-hot-toast' // ДОДАНО
import './Layout.css'

const USER_CREATOR_ROLES = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'PLATOON_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST', 'BRIGADE_STOREKEEPER', 'BATTALION_STOREKEEPER', 'COMPANY_SERGEANT']
const UNIT_MANAGER_ROLES = ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST', 'BRIGADE_STOREKEEPER', 'BATTALION_STOREKEEPER']

const ROLE_LABELS: Record<string, string> = {
  ADMIN: 'Адміністратор',
  BRIGADE_CMDR: 'Командир бригади',
  BATTALION_CMDR: 'Командир батальйону',
  COMPANY_CMDR: 'Командир роти',
  PLATOON_CMDR: 'Командир взводу',
  BRIGADE_LOGIST: 'Логіст бригади',
  BRIGADE_STOREKEEPER: 'Комірник бригади',
  BATTALION_LOGIST: 'Логіст батальйону',
  BATTALION_STOREKEEPER: 'Комірник батальйону',
  COMPANY_SERGEANT: 'Старшина роти',
  VOLUNTEER: 'Волонтер',
}

interface LayoutProps {
  children: React.ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const location = useLocation()
  const { user, token, loading, logout } = useAuth()

  // --- НОВЕ: Логіка сповіщень ---
  const [pendingCount, setPendingCount] = useState(0)
  const prevCountRef = useRef(0)

  // Перевіряємо, чи має поточний користувач права погоджувати заявки (щоб не спамити звичайних солдатів)
  const canApproveRequests = user?.role && ['ADMIN', 'BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'BRIGADE_LOGIST', 'BATTALION_LOGIST'].includes(user.role)

  useEffect(() => {
    // Запускаємо опитування тільки для авторизованих логістів/командирів
    if (!token || !canApproveRequests) return;

    const checkNewRequests = async () => {
      try {
        const reqs = await api.requests.list();
        if (Array.isArray(reqs)) {
          const currentPending = reqs.filter(r => r.status === 'PENDING').length;
          
          // Якщо кількість Очікуючих заявок зросла — викидаємо Пуш-сповіщення (Toast)
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

    checkNewRequests(); // Перша перевірка при завантаженні
    const interval = setInterval(checkNewRequests, 15000); // Опитування кожні 15 секунд
    
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
  const isVolunteer = user?.role === 'VOLUNTEER'

  return (
    <div className={`layout ${showSidebar ? 'with-sidebar' : ''}`}>
      {showSidebar && (
        <aside className="sidebar">
          <div className="sidebar-header">
            <Link to="/" className="sidebar-logo">
              Millog
            </Link>
          </div>
          <nav className="sidebar-nav">
            <Link to="/" className={location.pathname === '/' ? 'active' : ''}>
              Головна
            </Link>

            {!isVolunteer && (
              <>
                <Link to="/my-equipment" className={location.pathname === '/profile' ? 'active' : ''}>
                  Профіль
                </Link>
                <Link to="/analytics" className={location.pathname === '/analytics' ? 'active' : ''}>
                  Аналітика
                </Link>
                <Link to="/inventory" className={location.pathname === '/inventory' ? 'active' : ''}>
                  Ресурси
                </Link>
                
                {/* НОВЕ: Бейдж з кількістю заявок */}
                <Link 
                  to="/requests" 
                  className={location.pathname === '/requests' ? 'active' : ''}
                  style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}
                >
                  <span>Заявки</span>
                  {pendingCount > 0 && (
                    <span style={{
                      backgroundColor: '#ef4444', // Червоний колір
                      color: 'white',
                      fontSize: '11px',
                      fontWeight: 'bold',
                      padding: '2px 8px',
                      borderRadius: '12px',
                      boxShadow: '0 0 8px rgba(239, 68, 68, 0.4)' // Легке світіння
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

            <Link to="/volunteer-requests" className={location.pathname === '/volunteer-requests' ? 'active' : ''}>
              {isVolunteer ? 'Відкриті потреби' : 'Для волонтерів'}
            </Link>

            {canManageUnits && (
              <Link to="/units" className={location.pathname === '/units' ? 'active' : ''}>
                Підрозділи
              </Link>
            )}
            {canManageUsers && (
              <Link to="/admin/users" className={location.pathname === '/admin/users' ? 'active' : ''}>
                Користувачі
              </Link>
            )}
          </nav>
          <div className="sidebar-footer">
            <div className="user-badge">
              <span className="user-name">{user?.full_name}</span>
              <span className="user-role">{user?.role ? ROLE_LABELS[user.role] || user.role : ''}</span>
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
            <Link to="/" className="logo-text">Millog</Link>
            <Link to="/login" className="btn btn-primary">Увійти</Link>
          </header>
        )}
        <main className="page-content">{children}</main>
      </div>
    </div>
  )
}