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
      <div className="login-card" style={isForgotMode ? { 
        maxWidth: '500px', 
        gridTemplateColumns: '1fr',
        gap: '0'
      } : undefined}>
        
        {isForgotMode ? (
          <div className="login-form-section">
            <h1>Відновлення пароля</h1>
            <p style={{ marginBottom: '1.5rem' }}>
              Введіть ваш email, і ми надішлемо інструкції для встановлення нового пароля.
            </p>
            <form onSubmit={handleResetSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
              <input
                type="email"
                placeholder="Ваш Email"
                value={resetEmail}
                onChange={(e) => setResetEmail(e.target.value)}
                required
                style={{ width: '100%', margin: 0, padding: '0.875rem 1rem' }}
              />
              <button type="submit" className="btn btn-primary" disabled={isResetting} style={{ width: '100%', margin: 0 }}>
                {isResetting ? 'Відправка...' : 'Надіслати інструкції'}
              </button>
              <button 
                type="button" 
                className="btn btn-secondary" 
                onClick={() => setIsForgotMode(false)}
                style={{ width: '100%', margin: 0 }}
              >
                Повернутися до входу
              </button>
            </form>
          </div>
        ) : (
          <>
            <div className="login-form-section">
              <h1>Omnilog</h1>
              <p>Увійдіть до системи</p>
              
              {/* Demo hint для системного адміна */}
              <details className="demo-hint" style={{ 
                marginBottom: '1rem', 
                padding: '0.75rem', 
                background: 'rgba(59, 130, 246, 0.1)', 
                border: '1px solid #0ea5e9',
                borderRadius: '6px',
                fontSize: '0.875rem'
              }}>
                <summary style={{ cursor: 'pointer', fontWeight: 600, color: '#0369a1' }}>
                  🔐 Демо-доступи
                </summary>
                <div style={{ marginTop: '0.5rem', color: '#0c4a6e', lineHeight: '1.6' }}>
                  <strong>Системний адміністратор:</strong><br/>
                  Email: <code>platform@omnilog.system</code><br/>
                  Password: <code>AdminSystem2024!</code>
                  <hr style={{ margin: '0.5rem 0', border: 'none', borderTop: '1px solid #bae6fd' }}/>
                  <small style={{ color: '#075985' }}>
                    SYSTEM_ADMIN має доступ до всіх організацій через <code>/platform</code>
                  </small>
                </div>
              </details>
              
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
                
                <div style={{ 
                  textAlign: 'right',
                  marginBottom: '1.25rem',
                  background: 'transparent',
                  padding: 0
                }}>
                  <span 
                    onClick={() => {
                      setResetEmail(email) 
                      setIsForgotMode(true)
                      setError(null)
                    }}
                    style={{
                      color: 'var(--accent-light)',
                      fontSize: '0.875rem',
                      fontWeight: 500,
                      cursor: 'pointer',
                      textDecoration: 'none',
                      transition: 'all 0.2s ease'
                    }}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.color = 'var(--accent)';
                      e.currentTarget.style.textDecoration = 'underline';
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.color = 'var(--accent-light)';
                      e.currentTarget.style.textDecoration = 'none';
                    }}
                  >
                    Забули пароль?
                  </span>
                </div>

                {error && <div className="form-error">{error}</div>}
                
                <button type="submit" className="btn btn-primary" disabled={loading}>
                  {loading ? 'Вхід...' : 'Увійти'}
                </button>
              </form>
            </div>

            {/* Варіанти справа */}
            <div className="login-options-section">
              <h2>Вперше у нас?</h2>
              
              <div className="login-action-section" style={{ 
                background: 'transparent', 
                backgroundColor: 'transparent',
                border: 'none',
                boxShadow: 'none'
              }}>
                <p className="action-label">🆕 Новий бізнес</p>
                <Link to="/signup" className="btn btn-success btn-large">
                  ✨ Створити організацію
                </Link>
                <p className="action-hint">Перший запуск для нового бізнесу</p>
              </div>

              <div className="login-action-divider" style={{ 
                background: 'transparent', 
                backgroundColor: 'transparent',
                border: 'none',
                boxShadow: 'none'
              }}>або</div>

              <div className="login-action-section" style={{ 
                background: 'transparent', 
                backgroundColor: 'transparent',
                border: 'none',
                boxShadow: 'none'
              }}>
                <p className="action-label">👷 Підрядник</p>
                <Link to="/register" className="btn btn-warning btn-large">
                  📝 Реєстрація підрядника
                </Link>
                <p className="action-hint">Приєднайтесь як виконавець</p>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  )
}