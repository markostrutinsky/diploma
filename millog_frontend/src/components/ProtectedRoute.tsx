import { Navigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import './Layout.css'

interface ProtectedRouteProps {
  children: React.ReactNode
  forbidRoles?: string[]
  allowedRoles?: string[] // 🔥 Додали нову властивість
}

export default function ProtectedRoute({ children, forbidRoles, allowedRoles }: ProtectedRouteProps) {
  const { user, token, loading } = useAuth()

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

  return <>{children}</>
}