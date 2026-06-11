import { type ChangeEvent, type FormEvent, useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { api } from '../api/client'
import { useAuth } from '../contexts/AuthContext'
import './PlatformAdmin.css'

type Tenant = {
  id: string
  name: string
  slug: string
  subscription_tier: 'BASIC' | 'PRO' | 'ENTERPRISE'
  subscription_expires_at?: string | null
  owner_email?: string | null
  is_active: boolean
  created_at: string
  updated_at: string
  user_count: number
}

const TIERS = ['BASIC', 'PRO', 'ENTERPRISE'] as const

export default function PlatformAdmin() {
  const { refreshUser, supportTenant, enterSupportTenant, exitSupportTenant } = useAuth()
  const [stats, setStats] = useState<any>(null)
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)
  const [createOpen, setCreateOpen] = useState(false)
  const [activeMenuId, setActiveMenuId] = useState<string | null>(null)
  const [createForm, setCreateForm] = useState({
    organization_name: '',
    slug: '',
    owner_full_name: '',
    owner_email: '',
    owner_password: '',
  })

  const load = async () => {
    try {
      setLoading(true)
      const [s, list] = await Promise.all([api.platform.stats(), api.platform.listTenants(search)])
      setStats(s)
      setTenants(list as Tenant[])
    } catch (e: any) {
      toast.error(e?.message || 'Не вдалось завантажити дані')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    if (!activeMenuId) return

    const closeMenu = () => setActiveMenuId(null)
    window.addEventListener('click', closeMenu)
    return () => window.removeEventListener('click', closeMenu)
  }, [activeMenuId])

  const changeTier = async (t: Tenant, tier: string) => {
    try {
      await api.platform.updateTier(t.id, tier)
      toast.success(`Тариф "${t.name}" → ${tier}`)
      load()
      // Оновлюємо user в AuthContext — щоб платні фічи відразу відкрились
      await refreshUser()
    } catch (e: any) {
      toast.error(e?.message || 'Не вдалось змінити тариф')
    }
  }

  const toggleActive = async (t: Tenant) => {
    try {
      await api.platform.setActive(t.id, !t.is_active)
      toast.success(!t.is_active ? 'Організацію активовано' : 'Організацію призупинено')
      load()
    } catch (e: any) {
      toast.error(e?.message || 'Помилка')
    }
  }

  const removeTenant = async (t: Tenant) => {
    if (!confirm(`ВИДАЛИТИ організацію "${t.name}" разом з усіма даними? Цю дію неможливо відкотити!`)) return
    if (prompt(`Щоб підтвердити, введіть slug: ${t.slug}`) !== t.slug) {
      toast.error('Slug не співпадає — операцію скасовано')
      return
    }
    try {
      await api.platform.deleteTenant(t.id)
      toast.success('Організацію видалено')
      load()
    } catch (e: any) {
      toast.error(e?.message || 'Помилка видалення')
    }
  }

  const enterSupportMode = (t: Tenant) => {
    enterSupportTenant({ id: t.id, name: t.name, slug: t.slug })
    toast.success(`Support mode: ${t.name}`)
  }

  const updateCreateForm = (key: keyof typeof createForm) => (e: ChangeEvent<HTMLInputElement>) => {
    let value = e.target.value
    if (key === 'slug') value = value.toLowerCase().replace(/[^a-z0-9-]/g, '-').slice(0, 60)
    setCreateForm((prev) => ({ ...prev, [key]: value }))
  }

  const createTenant = async (e: FormEvent) => {
    e.preventDefault()
    try {
      await api.platform.createTenant({
        organization_name: createForm.organization_name.trim(),
        slug: createForm.slug.trim(),
        owner_email: createForm.owner_email.trim(),
        owner_full_name: createForm.owner_full_name.trim(),
        owner_password: createForm.owner_password,
      })
      toast.success('Організацію створено')
      setCreateForm({ organization_name: '', slug: '', owner_full_name: '', owner_email: '', owner_password: '' })
      setCreateOpen(false)
      load()
    } catch (err: any) {
      toast.error(err?.message || 'Не вдалось створити організацію')
    }
  }

  return (
    <div className="platform-admin">
      <div className="pa-header">
        <div className="pa-header-main">
          <div className="pa-title-row">
            <h1>Платформний адмін</h1>
            <span className={`pa-mode-badge ${supportTenant ? 'support' : ''}`}>
              {supportTenant ? `Support mode: ${supportTenant.name}` : 'Platform mode'}
            </span>
          </div>
          <p>Керування організаціями, тарифами та підписками.</p>
          {supportTenant && (
            <button type="button" className="pa-mode-exit" onClick={exitSupportTenant}>
              Вийти з організації
            </button>
          )}
        </div>
      </div>

      {stats && (
        <div className="pa-stats">
          <div className="pa-stat-card">
            <div className="pa-stat-label">Усього організацій</div>
            <div className="pa-stat-value">{stats.total_tenants}</div>
          </div>
          <div className="pa-stat-card">
            <div className="pa-stat-label">Активних</div>
            <div className="pa-stat-value">{stats.active_tenants}</div>
          </div>
          <div className="pa-stat-card">
            <div className="pa-stat-label">Усього користувачів</div>
            <div className="pa-stat-value">{stats.total_users}</div>
          </div>
          <div className="pa-stat-card">
            <div className="pa-stat-label">Нові за 30 днів</div>
            <div className="pa-stat-value">{stats.new_tenants_30_days}</div>
          </div>
          {TIERS.map((tier) => (
            <div className="pa-stat-card" key={tier}>
              <div className="pa-stat-label">{tier}</div>
              <div className="pa-stat-value">{stats.tenants_by_tier?.[tier] ?? 0}</div>
            </div>
          ))}
        </div>
      )}

      <div className="pa-toolbar">
        <div className="pa-toolbar-search">
          <input
            type="text"
            placeholder="Пошук за назвою/slug…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && load()}
          />
        </div>
        <div className="pa-toolbar-actions">
          <button onClick={() => setCreateOpen((value) => !value)}>
            {createOpen ? 'Закрити' : 'Створити організацію'}
          </button>
        </div>
      </div>

      {createOpen && (
        <form className="pa-create-form" onSubmit={createTenant}>
          <input required value={createForm.organization_name} onChange={updateCreateForm('organization_name')} placeholder="Назва організації" minLength={2} />
          <input required value={createForm.slug} onChange={updateCreateForm('slug')} placeholder="slug" minLength={2} pattern="[a-z0-9-]+" />
          <input required value={createForm.owner_full_name} onChange={updateCreateForm('owner_full_name')} placeholder="ПІБ власника" />
          <input required type="email" value={createForm.owner_email} onChange={updateCreateForm('owner_email')} placeholder="Email власника" />
          <input required type="password" value={createForm.owner_password} onChange={updateCreateForm('owner_password')} placeholder="Тимчасовий пароль" minLength={8} />
          <button type="submit">Створити</button>
        </form>
      )}

      <div className="pa-table-wrap">
        {loading ? (
          <div className="pa-loading">Завантаження…</div>
        ) : (
          <table className="pa-table">
            <thead>
              <tr>
                <th>Організація</th>
                <th>Slug</th>
                <th>Тариф</th>
                <th>Юзерів</th>
                <th>Власник</th>
                <th>Створено</th>
                <th>Статус</th>
                <th>Дії</th>
              </tr>
            </thead>
            <tbody>
              {tenants.map((t) => (
                <tr key={t.id} className={t.is_active ? '' : 'pa-row-inactive'}>
                  <td>
                    <strong>{t.name}</strong>
                  </td>
                  <td>
                    <code>{t.slug}</code>
                  </td>
                  <td>
                    <select value={t.subscription_tier} onChange={(e) => changeTier(t, e.target.value)}>
                      {TIERS.map((tier) => (
                        <option key={tier} value={tier}>
                          {tier}
                        </option>
                      ))}
                    </select>
                  </td>
                  <td>{t.user_count}</td>
                  <td>{t.owner_email || '—'}</td>
                  <td>{new Date(t.created_at).toLocaleDateString('uk-UA')}</td>
                  <td>
                    <span className={`pa-badge ${t.is_active ? 'ok' : 'off'}`}>
                      {t.is_active ? 'АКТИВНА' : 'ПРИЗУПИНЕНА'}
                    </span>
                  </td>
                  <td className="pa-actions-cell">
                    <div className="pa-actions">
                      <button className="primary" onClick={() => enterSupportMode(t)} disabled={!t.is_active}>
                        Support
                      </button>
                      <div className="pa-dropdown" onClick={(e) => e.stopPropagation()}>
                        <button
                          type="button"
                          className={`pa-kebab ${activeMenuId === t.id ? 'active' : ''}`}
                          aria-label={`Додаткові дії для ${t.name}`}
                          aria-expanded={activeMenuId === t.id}
                          onClick={(e) => {
                            e.stopPropagation()
                            setActiveMenuId(activeMenuId === t.id ? null : t.id)
                          }}
                        >
                          ⋮
                        </button>
                        {activeMenuId === t.id && (
                          <div className="pa-dropdown-menu">
                            <button
                              type="button"
                              onClick={() => {
                                toggleActive(t)
                                setActiveMenuId(null)
                              }}
                            >
                              {t.is_active ? 'Pause' : 'Resume'}
                            </button>
                            <div className="pa-dropdown-divider" />
                            <button
                              type="button"
                              className="danger"
                              onClick={() => {
                                removeTenant(t)
                                setActiveMenuId(null)
                              }}
                            >
                              Видалити
                            </button>
                          </div>
                        )}
                      </div>
                    </div>
                  </td>
                </tr>
              ))}
              {tenants.length === 0 && (
                <tr>
                  <td colSpan={8} className="pa-empty">
                    Нічого не знайдено
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
