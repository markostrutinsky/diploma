import { Navigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import './Layout.css'

interface ProtectedRouteProps {
  children: React.ReactNode
  forbidRoles?: string[]
  allowedRoles?: string[] // 🔥 Додали нову властивість
  requireTenantContext?: boolean
}

export default function ProtectedRoute({ children, forbidRoles, allowedRoles, requireTenantContext }: ProtectedRouteProps) {
  const { user, token, loading, supportTenant } = useAuth()

  if (loading) {
    return <div className="loading-full">Завантаження...</div>
  }

  if (!token) {
    return <Navigate to="/login" replace />
  }

  // Перевірка 1: Чи не заборонена ця роль (Чорний список)
  if (user && forbidRoles && forbidRoles.includes(user.role)) {
    console.warn(`Доступ заборонено (forbidRoles) для ролі: ${user.role}`)
    return <Navigate to="/" replace />
  }

  // Перевірка 2: Чи дозволена ця роль (Білий список)
  if (user && allowedRoles && !allowedRoles.includes(user.role)) {
    console.warn(`Доступ заборонено (allowedRoles) для ролі: ${user.role}`)
    return <Navigate to="/" replace />
  }

  if (user?.role === 'SYSTEM_ADMIN' && requireTenantContext && !supportTenant) {
    return <Navigate to="/platform" replace />
  }

  return <>{children}</>
}
