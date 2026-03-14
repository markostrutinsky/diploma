import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { api, type User } from '../api/client'

interface AuthContextType {
  user: User | null
  token: string | null
  loading: boolean
  login: (token: string, refreshToken: string, user: User) => void
  logout: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'))
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!token) {
      setUser(null)
      setLoading(false)
      return
    }
    const loadUser = async () => {
      try {
        const u = await api.auth.me()
        setUser(u)
      } catch {
        const refreshToken = localStorage.getItem('refresh_token')
        if (refreshToken) {
          try {
            const res = await api.auth.refresh(refreshToken)
            localStorage.setItem('token', res.token)
            localStorage.setItem('refresh_token', res.refresh_token)
            setToken(res.token)
            setUser(res.user)
            return
          } catch {
            // refresh failed
          }
        }
        localStorage.removeItem('token')
        localStorage.removeItem('refresh_token')
        setToken(null)
        setUser(null)
      } finally {
        setLoading(false)
      }
    }
    loadUser()
  }, [token])

  const login = (t: string, rt: string, u: User) => {
    localStorage.setItem('token', t)
    localStorage.setItem('refresh_token', rt)
    setToken(t)
    setUser(u)
  }

  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    setToken(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, token, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
