import { useEffect, useMemo, useRef, useState, type CSSProperties } from 'react'
import { createPortal } from 'react-dom'
import { Link } from 'react-router-dom'
import { api } from '../api/client'
import { usePermissions } from '../hooks/usePermissions'
import './NotificationCenter.css'

type Severity = 'info' | 'warning' | 'danger'

interface NotificationItem {
  id: string
  icon: string
  title: string
  message: string
  link?: string
  severity: Severity
  timestamp: number
}

const POLL_INTERVAL_MS = 45_000

export default function NotificationCenter() {
  const perms = usePermissions()
  const [open, setOpen] = useState(false)
  const [items, setItems] = useState<NotificationItem[]>([])
  const [dropdownStyle, setDropdownStyle] = useState<CSSProperties>()
  const [dismissed, setDismissed] = useState<Set<string>>(() => {
    try {
      const raw = localStorage.getItem('notif:dismissed')
      return new Set(raw ? JSON.parse(raw) : [])
    } catch { return new Set() }
  })
  const ref = useRef<HTMLDivElement>(null)
  const dropdownRef = useRef<HTMLDivElement>(null)

  const canApproveRequests = perms.can('request_approve')
  const canSeeInventory = perms.canAny('resource_manage', 'warehouse_manage')
  const canSeeVehicles = perms.canAny('vehicle_manage', 'vehicle_maintenance', 'vehicle_fuel_log')
  const hasFuelAntifraud = perms.hasFeature('fuel_antifraud')

  useEffect(() => {
    if (!perms.isAuthenticated) return
    let cancelled = false

    const load = async () => {
      const next: NotificationItem[] = []

      // 0. Personal notifications (shipment assignments, etc.) - для всіх користувачів
      try {
        const notifData = await api.notifications.list(10)
        if (notifData?.notifications && Array.isArray(notifData.notifications)) {
          notifData.notifications
            .filter((n: any) => !n.is_read) // Тільки непрочитані
            .forEach((n: any) => {
              let icon = '🔔'
              let severity: Severity = 'info'
              
              switch (n.type) {
                case 'SHIPMENT_ASSIGNED':
                  icon = '🚚'
                  severity = 'warning'
                  break
                case 'REQUEST_APPROVED':
                  icon = '✅'
                  severity = 'info'
                  break
                case 'REQUEST_REJECTED':
                  icon = '❌'
                  severity = 'danger'
                  break
                case 'SHIPMENT_DELIVERED':
                  icon = '📦'
                  severity = 'info'
                  break
                case 'LOW_STOCK':
                  icon = '⚠️'
                  severity = 'warning'
                  break
              }
              
              next.push({
                id: `personal:${n.id}`,
                icon,
                title: n.title || 'Сповіщення',
                message: n.message || '',
                severity,
                timestamp: new Date(n.created_at).getTime(),
              })
            })
        }
      } catch { /* ignore */ }

      // 1. Pending approvals
      if (canApproveRequests) {
        try {
          const reqs: any[] = await api.requests.list()
          const pending = (reqs || []).filter(r => r?.status === 'PENDING')
          if (pending.length > 0) {
            next.push({
              id: `pending-requests:${pending.length}`,
              icon: '📝',
              title: 'Заявки на погодження',
              message: `${pending.length} заявок очікують вашого рішення`,
              link: '/requests',
              severity: 'warning',
              timestamp: Date.now(),
            })
          }
        } catch { /* ignore */ }
      }

      // 2. Low stock
      if (canSeeInventory) {
        try {
          const resources: any[] = await api.inventory.listResources()
          const low = (resources || []).filter(r =>
            typeof r.min_quantity === 'number' && r.min_quantity > 0 &&
            typeof r.quantity === 'number' && r.quantity <= r.min_quantity &&
            r.status !== 'WRITTEN_OFF'
          )
          if (low.length > 0) {
            next.push({
              id: `low-stock:${low.length}`,
              icon: '📉',
              title: 'Низький залишок ресурсів',
              message: `${low.length} позицій нижче мінімуму. Перевірте склад.`,
              link: '/inventory',
              severity: 'danger',
              timestamp: Date.now(),
            })
          }
        } catch { /* ignore */ }
      }

      // 3. Maintenance alerts
      if (canSeeVehicles) {
        try {
          const vehicles: any[] = await api.vehicles.list()
          const overdue = (vehicles || []).filter(v =>
            v?.maintenance_status === 'OVERDUE' || v?.maintenance_status === 'SOON'
          )
          if (overdue.length > 0) {
            const soonCount = overdue.filter(v => v.maintenance_status === 'SOON').length
            const overdueCount = overdue.filter(v => v.maintenance_status === 'OVERDUE').length
            next.push({
              id: `maint:${overdueCount}:${soonCount}`,
              icon: '🔧',
              title: 'Техобслуговування',
              message: overdueCount
                ? `${overdueCount} ТЗ прострочено ТО${soonCount ? `, ${soonCount} — скоро` : ''}`
                : `${soonCount} ТЗ скоро потребують ТО`,
              link: '/vehicles',
              severity: overdueCount > 0 ? 'danger' : 'warning',
              timestamp: Date.now(),
            })
          }

          // Fuel anomalies — PRO only
          if (hasFuelAntifraud) {
            const anomalous = (vehicles || []).filter(v => v?.has_fuel_anomaly)
            if (anomalous.length > 0) {
              next.push({
                id: `fuel-anomaly:${anomalous.length}`,
                icon: '🚨',
                title: 'Підозріле пальне',
                message: `Виявлено аномалії у ${anomalous.length} ТЗ`,
                link: '/vehicles',
                severity: 'danger',
                timestamp: Date.now(),
              })
            }
          }
        } catch { /* ignore */ }
      }

      if (!cancelled) setItems(next)
    }

    load()
    const t = setInterval(load, POLL_INTERVAL_MS)
    return () => { cancelled = true; clearInterval(t) }
  }, [perms.isAuthenticated, canApproveRequests, canSeeInventory, canSeeVehicles, hasFuelAntifraud])

  useEffect(() => {
    const onClick = (e: MouseEvent) => {
      const target = e.target as Node
      if (ref.current?.contains(target) || dropdownRef.current?.contains(target)) return
      setOpen(false)
    }
    document.addEventListener('mousedown', onClick)
    return () => document.removeEventListener('mousedown', onClick)
  }, [])

  useEffect(() => {
    if (!open) return

    const updateDropdownPosition = () => {
      const bellRect = ref.current?.getBoundingClientRect()
      if (!bellRect) return

      const DROPDOWN_WIDTH = 360
      const GAP = 10
      const VIEWPORT_PADDING = 8
      const top = Math.round(bellRect.bottom + GAP)
      const left = Math.min(
        Math.max(VIEWPORT_PADDING, Math.round(bellRect.left)),
        Math.round(window.innerWidth - DROPDOWN_WIDTH - VIEWPORT_PADDING)
      )

      setDropdownStyle({
        top,
        left,
        width: DROPDOWN_WIDTH,
      })
    }

    updateDropdownPosition()
    window.addEventListener('resize', updateDropdownPosition)
    window.addEventListener('scroll', updateDropdownPosition, true)

    return () => {
      window.removeEventListener('resize', updateDropdownPosition)
      window.removeEventListener('scroll', updateDropdownPosition, true)
    }
  }, [open])

  const visibleItems = useMemo(
    () => items.filter(i => !dismissed.has(i.id)),
    [items, dismissed]
  )

  const dismiss = (id: string) => {
    const next = new Set(dismissed)
    next.add(id)
    setDismissed(next)
    try { localStorage.setItem('notif:dismissed', JSON.stringify([...next])) } catch { /* ignore */ }
    
    // Якщо це персональне сповіщення, позначити як прочитане на сервері
    if (id.startsWith('personal:')) {
      const notifId = id.replace('personal:', '')
      api.notifications.markAsRead(notifId).catch(() => { /* ignore */ })
    }
  }

  const clearDismissed = () => {
    setDismissed(new Set())
    try { localStorage.removeItem('notif:dismissed') } catch { /* ignore */ }
  }

  if (!perms.isAuthenticated) return null

  const count = visibleItems.length
  const hasDanger = visibleItems.some(i => i.severity === 'danger')

  return (
    <div className="notif-center" ref={ref}>
      <button
        className={`notif-bell ${count ? 'has-items' : ''} ${hasDanger ? 'has-danger' : ''}`}
        onClick={() => setOpen(o => !o)}
        aria-label="Сповіщення"
        title="Сповіщення"
      >
        <span className="notif-icon">🔔</span>
        {count > 0 && <span className="notif-badge">{count}</span>}
      </button>

      {open && createPortal(
        <div ref={dropdownRef} className="notif-dropdown" style={dropdownStyle}>
          <div className="notif-header">
            <strong>Сповіщення</strong>
            {dismissed.size > 0 && (
              <button className="notif-link-btn" onClick={clearDismissed}>
                Показати приховані
              </button>
            )}
          </div>

          {visibleItems.length === 0 ? (
            <div className="notif-empty">
              <div className="notif-empty-icon">✨</div>
              <div>Усе спокійно — активних подій немає</div>
            </div>
          ) : (
            <ul className="notif-list">
              {visibleItems.map(item => (
                <li key={item.id} className={`notif-item sev-${item.severity}`}>
                  <span className="notif-item-icon">{item.icon}</span>
                  <div className="notif-item-body">
                    <div className="notif-item-title">{item.title}</div>
                    <div className="notif-item-message">{item.message}</div>
                    {item.link && (
                      <Link to={item.link} className="notif-item-link" onClick={() => setOpen(false)}>
                        Перейти →
                      </Link>
                    )}
                  </div>
                  <button
                    className="notif-item-close"
                    onClick={() => dismiss(item.id)}
                    title="Приховати"
                  >×</button>
                </li>
              ))}
            </ul>
          )}
        </div>,
        document.body
      )}
    </div>
  )
}
