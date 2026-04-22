import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { api } from '../api/client'
import './SignupTenant.css'

export default function SignupTenant() {
  const navigate = useNavigate()
  const [form, setForm] = useState({
    organization_name: '',
    slug: '',
    owner_full_name: '',
    owner_email: '',
    owner_password: '',
    confirm_password: '',
  })
  const [loading, setLoading] = useState(false)

  const onChange = (k: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement>) => {
    let v = e.target.value
    if (k === 'slug') v = v.toLowerCase().replace(/[^a-z0-9-]/g, '-').slice(0, 60)
    setForm((f) => ({ ...f, [k]: v }))
  }

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (form.owner_password !== form.confirm_password) {
      toast.error('Паролі не співпадають')
      return
    }
    if (form.owner_password.length < 8) {
      toast.error('Пароль має бути щонайменше 8 символів')
      return
    }
    try {
      setLoading(true)
      await api.auth.signupTenant({
        organization_name: form.organization_name.trim(),
        slug: form.slug.trim(),
        owner_email: form.owner_email.trim(),
        owner_full_name: form.owner_full_name.trim(),
        owner_password: form.owner_password,
      })
      toast.success('Організацію створено! Тепер увійдіть.')
      navigate('/login')
    } catch (err: any) {
      toast.error(err?.message || 'Помилка реєстрації')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="signup-tenant-page">
      <div className="signup-tenant-card">
        <h1>Створити організацію</h1>
        <p className="subtitle">
          Реєстрація нової організації на платформі. Ви станете її власником (TENANT_ADMIN).
        </p>

        <form onSubmit={submit}>
          <label>
            Назва організації*
            <input
              required
              type="text"
              value={form.organization_name}
              onChange={onChange('organization_name')}
              placeholder="ТОВ «Нова Логістика»"
              minLength={2}
            />
          </label>

          <label>
            Короткий ідентифікатор (slug)*
            <input
              required
              type="text"
              value={form.slug}
              onChange={onChange('slug')}
              placeholder="nova-log"
              minLength={2}
              pattern="[a-z0-9-]+"
            />
            <small>Унікальний, латиницею, без пробілів. Використовується в URL.</small>
          </label>

          <label>
            Ваше ПІБ*
            <input
              required
              type="text"
              value={form.owner_full_name}
              onChange={onChange('owner_full_name')}
              placeholder="Іван Петренко"
            />
          </label>

          <label>
            Ваш email*
            <input
              required
              type="email"
              value={form.owner_email}
              onChange={onChange('owner_email')}
              placeholder="owner@company.ua"
            />
          </label>

          <label>
            Пароль* (мін. 8 символів)
            <input
              required
              type="password"
              value={form.owner_password}
              onChange={onChange('owner_password')}
              minLength={8}
            />
          </label>

          <label>
            Повторіть пароль*
            <input
              required
              type="password"
              value={form.confirm_password}
              onChange={onChange('confirm_password')}
              minLength={8}
            />
          </label>

          <button type="submit" className="primary-btn" disabled={loading}>
            {loading ? 'Створення…' : 'Створити організацію'}
          </button>
        </form>

        <div className="signup-tenant-footer">
          Вже є акаунт? <Link to="/login">Увійти</Link>
        </div>
      </div>
    </div>
  )
}
