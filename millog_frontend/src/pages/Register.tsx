import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import './Login.css'

export default function Register() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [fullName, setFullName] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await api.auth.register({ email, password, full_name: fullName })
      const loginRes = await api.auth.login(email, password)
      login(loginRes.token, loginRes.refresh_token, loginRes.user)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Помилка реєстрації')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      <div className="login-card">
        <h1>Реєстрація волонтера</h1>
        <p>Створіть обліковий запис для виконання заявок</p>
        <p className="login-hint">
          <Link to="/login">Маєте обліковий запис — увійти</Link>
        </p>
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            placeholder="ПІБ"
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
            required
          />
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Пароль (мін. 8 символів)"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
          />
          {error && <div className="form-error">{error}</div>}
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Реєстрація...' : 'Зареєструватися'}
          </button>
        </form>
      </div>
    </div>
  )
}
