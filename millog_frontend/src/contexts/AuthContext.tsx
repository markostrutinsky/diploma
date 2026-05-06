import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { api, setInMemoryToken, type User } from '../api/client'

interface AuthContextType {
  user: User | null
  token: string | null
  loading: boolean
  login: (token: string, refreshToken: string, user: User) => void
  logout: () => void
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // При запуску намагаємось відновити сесію через httpOnly cookie (refresh_token)
    const restoreSession = async () => {
      try {
        const res = await api.auth.refresh()
        setInMemoryToken(res.token)
        setToken(res.token)
        setUser(res.user)
      } catch {
        // Cookie відсутній або протермінований — користувач не авторизований
        setInMemoryToken(null)
        setToken(null)
        setUser(null)
      } finally {
        setLoading(false)
      }
    }
    restoreSession()
  }, [])

  const login = (t: string, _rt: string, u: User) => {
    // refresh_token зберігається у httpOnly cookie на бекенді
    // access_token зберігаємо лише в пам'яті
    setInMemoryToken(t)
    setToken(t)
    setUser(u)
  }

  const logout = async () => {
    try { await api.auth.logout() } catch { /* ігноруємо помилку */ }
    setInMemoryToken(null)
    setToken(null)
    setUser(null)
  }

  const refreshUser = async () => {
    try {
      // Просто оновлюємо user об'єкт з /auth/me без ротації токена
      const freshUser = await api.auth.me()
      setUser(freshUser)
    } catch {
      // Якщо /auth/me недоступний — пробуємо через refresh
      try {
        const res = await api.auth.refresh()
        setInMemoryToken(res.token)
        setToken(res.token)
        setUser(res.user)
      } catch { /* ігноруємо */ }
    }
  }

  return (
    <AuthContext.Provider value={{ user, token, loading, login, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
