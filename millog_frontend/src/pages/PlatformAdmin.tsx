import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { api } from '../api/client'
import './PlatformAdmin.css'

type Tenant = {
  id: string
  name: string
  slug: string
  subscription_tier: 'FREE' | 'BASIC' | 'PRO' | 'ENTERPRISE'
  subscription_expires_at?: string | null
  owner_email?: string | null
  is_active: boolean
  created_at: string
  updated_at: string
  user_count: number
}

const TIERS = ['FREE', 'BASIC', 'PRO', 'ENTERPRISE'] as const

export default function PlatformAdmin() {
  const [stats, setStats] = useState<any>(null)
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)

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

  const changeTier = async (t: Tenant, tier: string) => {
    try {
      await api.platform.updateTier(t.id, tier)
      toast.success(`Тариф "${t.name}" → ${tier}`)
      load()
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

  return (
    <div className="platform-admin">
      <div className="pa-header">
        <h1>Платформний адмін</h1>
        <p>Керування організаціями, тарифами та підписками.</p>
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
        <input
          type="text"
          placeholder="Пошук за назвою/slug…"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && load()}
        />
        <button onClick={load}>Пошук</button>
      </div>

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
                <th></th>
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
                  <td className="pa-actions">
                    <button onClick={() => toggleActive(t)}>
                      {t.is_active ? 'Призупинити' : 'Активувати'}
                    </button>
                    <button className="danger" onClick={() => removeTenant(t)}>
                      Видалити
                    </button>
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
