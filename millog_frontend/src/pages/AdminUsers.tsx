import { useEffect, useState } from 'react'
import { api, type CreateUserRequest, type Unit } from '../api/client'
import './AdminUsers.css'

const ROLES = [
  { value: 'ADMIN', label: 'Адмін' },
  { value: 'BRIGADE_CMDR', label: 'Командир бригади' },
  { value: 'BATTALION_CMDR', label: 'Командир батальйону' },
  { value: 'COMPANY_CMDR', label: 'Командир роти' },
  { value: 'PLATOON_CMDR', label: 'Командир взводу' },
  { value: 'BRIGADE_LOGIST', label: 'Логіст бригади' },
  { value: 'BRIGADE_STOREKEEPER', label: 'Комірник бригади' },
  { value: 'BATTALION_LOGIST', label: 'Логіст батальйону' },
  { value: 'BATTALION_STOREKEEPER', label: 'Комірник батальйону' },
  { value: 'COMPANY_SERGEANT', label: 'Старшина роти' },
] as const

export default function AdminUsers() {
  const [units, setUnits] = useState<Unit[]>([])
  const [form, setForm] = useState<Omit<CreateUserRequest, 'role'> & { role: (typeof ROLES)[number]['value'] }>({
    email: '',
    full_name: '',
    role: 'COMPANY_SERGEANT',
  })
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    api.units.list().then((data) => setUnits(Array.isArray(data) ? data : [])).catch(() => {})
  }, [])
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setSuccess(null)
    setLoading(true)

    try {
      const body: CreateUserRequest = {
        email: form.email.trim(),
        full_name: form.full_name.trim(),
        role: form.role,
      }
      if (form.username?.trim()) body.username = form.username.trim()
      if (form.phone?.trim()) body.phone = form.phone.trim()
      if (form.unit_id) body.unit_id = form.unit_id

      const res = await api.admin.createUser(body)
      setSuccess(res.message)
      setForm({ email: '', full_name: '', role: 'COMPANY_SERGEANT', username: '', phone: '' })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Помилка')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="admin-users">
      <div className="page-header">
        <h1>Користувачі</h1>
        <p className="page-subtitle">
          Створіть обліковий запис. Користувач отримає лист з посиланням для встановлення паролю. Волонтери реєструються самостійно.
        </p>
      </div>

      <div className="card">
        <form className="admin-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Email <span className="required">*</span></label>
            <input
              type="email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              required
              placeholder="user@example.com"
            />
          </div>
          <div className="form-group">
            <label>ПІБ <span className="required">*</span></label>
            <input
              type="text"
              value={form.full_name}
              onChange={(e) => setForm({ ...form, full_name: e.target.value })}
              required
              placeholder="Іван Петренко"
            />
          </div>
          <div className="form-group">
            <label>Логін</label>
            <input
              type="text"
              value={form.username ?? ''}
              onChange={(e) => setForm({ ...form, username: e.target.value || undefined })}
              placeholder="Опціонально"
            />
          </div>
          <div className="form-group">
            <label>Телефон</label>
            <input
              type="tel"
              value={form.phone ?? ''}
              onChange={(e) => setForm({ ...form, phone: e.target.value || undefined })}
              placeholder="+380..."
            />
          </div>
          <div className="form-group">
            <label>Роль</label>
            <select
              value={form.role}
              onChange={(e) => setForm({ ...form, role: e.target.value as (typeof ROLES)[number]['value'] })}
            >
              {ROLES.map((r) => (
                <option key={r.value} value={r.value}>{r.label}</option>
              ))}
            </select>
          </div>
          <div className="form-group">
            <label>Підрозділ</label>
            <select
              value={form.unit_id ?? ''}
              onChange={(e) => {
                const val = e.target.value
                setForm({ ...form, unit_id: val ? parseInt(val, 10) : undefined })
              }}
            >
              <option value="">Не обрано</option>
              {units.map((u) => (
                <option key={u.id} value={u.id}>{u.name}</option>
              ))}
            </select>
          </div>

          {error && <div className="form-error">{error}</div>}
          {success && <div className="form-success">{success}</div>}

          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Створення...' : 'Створити'}
          </button>
        </form>
      </div>
    </div>
  )
}
