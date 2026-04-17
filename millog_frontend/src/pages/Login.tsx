import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast' 
import { api } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import './Login.css'

export default function Login() {
  // --- СТЕЙТИ ЛОГІНУ ---
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  
  // --- СТЕЙТИ ВІДНОВЛЕННЯ ПАРОЛЯ ---
  const [isForgotMode, setIsForgotMode] = useState(false)
  const [resetEmail, setResetEmail] = useState('')
  const [isResetting, setIsResetting] = useState(false)

  const { login } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const res = await api.auth.login(email, password)
      login(res.token, res.refresh_token, res.user)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Помилка входу')
    } finally {
      setLoading(false)
    }
  }

  const handleResetSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!resetEmail) return
    
    try {
      setIsResetting(true)
      const res = await api.auth.requestPasswordReset(resetEmail)
      toast.success(res.message || 'Інструкції надіслано на пошту!')
      setIsForgotMode(false) 
    } catch (error: any) {
      toast.error(error.message || 'Помилка при відправці запиту')
    } finally {
      setIsResetting(false)
    }
  }

  return (
    <div className="login-page">
      <div className="login-card">
        
        {isForgotMode ? (
          <>
            <h1>Відновлення пароля</h1>
            <p style={{ marginBottom: '1.5rem' }}>
              Введіть ваш email, і ми надішлемо інструкції для встановлення нового пароля.
            </p>
            <form onSubmit={handleResetSubmit}>
              <input
                type="email"
                placeholder="Ваш Email"
                value={resetEmail}
                onChange={(e) => setResetEmail(e.target.value)}
                required
              />
              <button type="submit" className="btn btn-primary" disabled={isResetting}>
                {isResetting ? 'Відправка...' : 'Надіслати інструкції'}
              </button>
              <button 
                type="button" 
                className="btn btn-secondary" 
                onClick={() => setIsForgotMode(false)}
              >
                Повернутися до входу
              </button>
            </form>
          </>
        ) : (
          <>
            <h1>Omnilog</h1>
            <p>Увійдіть до системи</p>
            <p className="login-hint">
              <Link to="/bootstrap">Перший запуск — створіть адміна</Link>
              {' · '}
              {/* 🔥 Змінено з "Реєстрація волонтера" на "Реєстрація підрядника" */}
              <Link to="/register">Реєстрація підрядника</Link>
            </p>
            <form onSubmit={handleSubmit}>
              <input
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
              <input
                type="password"
                placeholder="Пароль"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                style={{ marginBottom: '0.5rem' }} 
              />
              
              <div className="forgot-password-container">
                <button 
                  type="button" 
                  className="forgot-password-link"
                  onClick={() => {
                    setResetEmail(email) 
                    setIsForgotMode(true)
                    setError(null)
                  }}
                >
                  Забули пароль?
                </button>
              </div>

              {error && <div className="form-error">{error}</div>}
              
              <button type="submit" className="btn btn-primary" disabled={loading}>
                {loading ? 'Вхід...' : 'Увійти'}
              </button>
            </form>
          </>
        )}
      </div>
    </div>
  )
}