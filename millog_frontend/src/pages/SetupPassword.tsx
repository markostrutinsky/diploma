import { useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { api } from '../api/client'
import './SetupPassword.css'

export default function SetupPassword() {
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token')

  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!token) return
    setError(null)

    if (password.length < 8) {
      setError('Пароль має бути не менше 8 символів')
      return
    }
    if (!/[A-Z]/.test(password)) {
      setError('Пароль має містити хоча б одну велику літеру')
      return
    }
    if (!/[0-9]/.test(password)) {
      setError('Пароль має містити хоча б одну цифру')
      return
    }
    if (!/[^A-Za-z0-9]/.test(password)) {
      setError('Пароль має містити хоча б один спецсимвол (!@#$... тощо)')
      return
    }
    if (password !== confirm) {
      setError('Паролі не збігаються')
      return
    }

    setLoading(true)
    try {
      await api.auth.setupPassword(token, password)
      setSuccess(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Помилка встановлення пароля')
    } finally {
      setLoading(false)
    }
  }

  if (!token) {
    return (
      <div className="setup-password">
        <div className="setup-card">
          <h1>Недійсне посилання</h1>
          <p>Посилання для встановлення пароля відсутнє або застаріло.</p>
        </div>
      </div>
    )
  }

  if (success) {
    return (
      <div className="setup-password">
        <div className="setup-card success">
          <h1>Пароль встановлено</h1>
          <p>Ваш обліковий запис активовано. Можете увійти в систему.</p>
          <Link to="/login" className="btn btn-primary" style={{ marginTop: '1rem', display: 'inline-block' }}>
            Увійти
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="setup-password">
      <div className="setup-card">
        <h1>Встановити пароль</h1>
        <p>Оберіть надійний пароль для вашого облікового запису.</p>
        <ul style={{ fontSize: '13px', color: 'var(--text-muted)', marginBottom: '1rem', paddingLeft: '1.2rem' }}>
          <li style={{ color: password.length >= 8 ? 'var(--success)' : undefined }}>Мінімум 8 символів</li>
          <li style={{ color: /[A-Z]/.test(password) ? 'var(--success)' : undefined }}>Хоча б одна велика літера</li>
          <li style={{ color: /[0-9]/.test(password) ? 'var(--success)' : undefined }}>Хоча б одна цифра</li>
          <li style={{ color: /[^A-Za-z0-9]/.test(password) ? 'var(--success)' : undefined }}>Хоча б один спецсимвол (!@#$...)</li>
        </ul>

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Пароль</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Мінімум 8 символів"
              required
              minLength={8}
            />
          </div>
          <div className="form-group">
            <label>Підтвердження</label>
            <input
              type="password"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              placeholder="Повторіть пароль"
              required
            />
          </div>
          {error && <div className="form-error">{error}</div>}
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Збереження...' : 'Зберегти'}
          </button>
        </form>
      </div>
    </div>
  )
}
