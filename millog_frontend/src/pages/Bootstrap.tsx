import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import './Bootstrap.css'

export default function Bootstrap() {
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
      await api.auth.bootstrap(email, password, fullName)
      const res = await api.auth.login(email, password)
      login(res.token, res.refresh_token, res.user)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Помилка')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="bootstrap-page">
      <div className="bootstrap-card">
        <h1>Перший запуск</h1>
        <p>Створіть обліковий запис адміністратора</p>
        <form onSubmit={handleSubmit}>
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            type="text"
            placeholder="ПІБ"
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
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
            {loading ? 'Створення...' : 'Створити'}
          </button>
        </form>
      </div>
    </div>
  )
}
