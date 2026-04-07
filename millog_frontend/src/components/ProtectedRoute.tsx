import { Navigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import './Layout.css'

interface ProtectedRouteProps {
  children: React.ReactNode
  forbidRoles?: string[]
}

export default function ProtectedRoute({ children, forbidRoles }: ProtectedRouteProps) {
  const { user, token, loading } = useAuth()

  if (loading) {
    return <div className="loading-full">Завантаження...</div>
  }

  if (!token) {
    return <Navigate to="/login" replace />
  }

  if (user && forbidRoles && forbidRoles.includes(user.role)) {
    console.warn(`Доступ заборонено для ролі: ${user.role}`)
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}