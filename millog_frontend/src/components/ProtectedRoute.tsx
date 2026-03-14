import { Navigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import './Layout.css'

export default function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { token, loading } = useAuth()
  if (loading) return <div className="loading-full">Завантаження...</div>
  if (!token) return <Navigate to="/login" replace />
  return <>{children}</>
}
