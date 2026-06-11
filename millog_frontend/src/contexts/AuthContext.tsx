import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { api, getSupportTenantId, hasAuthSessionMarker, setAuthSessionMarker, setInMemoryToken, setSupportTenantId, type User } from '../api/client'

export type SupportTenant = {
  id: string
  name: string
  slug: string
}

interface AuthContextType {
  user: User | null
  token: string | null
  loading: boolean
  login: (token: string, refreshToken: string, user: User) => void
  logout: () => void
  refreshUser: () => Promise<void>
  supportTenant: SupportTenant | null
  enterSupportTenant: (tenant: SupportTenant) => void
  exitSupportTenant: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [supportTenant, setSupportTenant] = useState<SupportTenant | null>(() => {
    try {
      const raw = localStorage.getItem('omnilog_support_tenant_meta')
      const tenant = raw ? JSON.parse(raw) as SupportTenant : null
      return tenant?.id === getSupportTenantId() ? tenant : null
    } catch {
      return null
    }
  })

  useEffect(() => {
    // При запуску намагаємось відновити сесію через httpOnly cookie (refresh_token)
    const restoreSession = async () => {
      if (!hasAuthSessionMarker()) {
        setLoading(false)
        return
      }

      try {
        const res = await api.auth.refresh()
        setInMemoryToken(res.token)
        setToken(res.token)
        setUser(res.user)
        if (res.user.role !== 'SYSTEM_ADMIN') {
          setSupportTenantId(null)
          setSupportTenant(null)
          localStorage.removeItem('omnilog_support_tenant_meta')
        }
      } catch {
        // Cookie відсутній або протермінований — користувач не авторизований
        setInMemoryToken(null)
        setToken(null)
        setUser(null)
        setAuthSessionMarker(false)
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
    setAuthSessionMarker(true)
    if (u.role !== 'SYSTEM_ADMIN') {
      setSupportTenantId(null)
      setSupportTenant(null)
      localStorage.removeItem('omnilog_support_tenant_meta')
    }
  }

  const logout = async () => {
    try { await api.auth.logout() } catch { /* ігноруємо помилку */ }
    setInMemoryToken(null)
    setToken(null)
    setUser(null)
    setAuthSessionMarker(false)
    setSupportTenantId(null)
    setSupportTenant(null)
    localStorage.removeItem('omnilog_support_tenant_meta')
  }

  const refreshUser = async () => {
    try {
      // Просто оновлюємо user об'єкт з /auth/me без ротації токена
      const freshUser = await api.auth.me()
      setUser(freshUser)
      if (freshUser.role !== 'SYSTEM_ADMIN') {
        exitSupportTenant()
      }
    } catch {
      // Якщо /auth/me недоступний — пробуємо через refresh
      if (!hasAuthSessionMarker()) return

      try {
        const res = await api.auth.refresh()
        setInMemoryToken(res.token)
        setToken(res.token)
        setUser(res.user)
        if (res.user.role !== 'SYSTEM_ADMIN') {
          setSupportTenantId(null)
          setSupportTenant(null)
          localStorage.removeItem('omnilog_support_tenant_meta')
        }
      } catch {
        setAuthSessionMarker(false)
      }
    }
  }

  const enterSupportTenant = (tenant: SupportTenant) => {
    setSupportTenantId(tenant.id)
    setSupportTenant(tenant)
    localStorage.setItem('omnilog_support_tenant_meta', JSON.stringify(tenant))
  }

  const exitSupportTenant = () => {
    setSupportTenantId(null)
    setSupportTenant(null)
    localStorage.removeItem('omnilog_support_tenant_meta')
  }

  return (
    <AuthContext.Provider value={{ user, token, loading, login, logout, refreshUser, supportTenant, enterSupportTenant, exitSupportTenant }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
