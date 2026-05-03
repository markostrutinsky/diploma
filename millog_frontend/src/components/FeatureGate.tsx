import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import { FEATURES, type FeatureKey, TIER_NAMES } from '../constants/roles'
import { usePermissions } from '../hooks/usePermissions'
import './FeatureGate.css'

interface FeatureGateProps {
  feature: FeatureKey
  children: ReactNode
  /** Якщо true — PRO-контент відображається, але поверх нього напівпрозора заглушка з CTA. */
  showTeaser?: boolean
  /** Якщо задано — рендериться замість заглушки, коли фіча недоступна. */
  fallback?: ReactNode
}

/**
 * Обгортка для платної функціональності.
 * - Якщо користувач має доступ → рендериться `children`.
 * - Інакше: або `fallback`, або `<PaywallCard />` з кнопкою на /billing.
 */
export function FeatureGate({ feature, children, showTeaser, fallback }: FeatureGateProps) {
  const perms = usePermissions()

  if (perms.hasFeature(feature)) {
    return <>{children}</>
  }

  if (showTeaser) {
    return (
      <div className="feature-gate-teaser">
        <div className="feature-gate-teaser__content" aria-hidden="true">
          {children}
        </div>
        <div className="feature-gate-teaser__overlay">
          <PaywallCard feature={feature} />
        </div>
      </div>
    )
  }

  return <>{fallback ?? <PaywallCard feature={feature} />}</>
}

interface PaywallCardProps {
  feature: FeatureKey
}

/** Карточка-заглушка з CTA на сторінку тарифів. */
export function PaywallCard({ feature }: PaywallCardProps) {
  const meta = FEATURES[feature]
  return (
    <div className="paywall-card">
      <div className="paywall-card__icon">🔒</div>
      <div className="paywall-card__body">
        <div className="paywall-card__title">{meta.label}</div>
        <div className="paywall-card__hint">
          Доступно на тарифі <strong>{TIER_NAMES[meta.minTier]}</strong> і вище.
        </div>
      </div>
      <Link to="/billing" className="paywall-card__cta">
        Переглянути тарифи
      </Link>
    </div>
  )
}

interface PaywallScreenProps {
  feature: FeatureKey
  /** Необов'язковий опис чому саме ця фіча корисна — зʼявиться під основним текстом. */
  description?: string
}

/**
 * Повноекранна заглушка для платних сторінок (GPS, KPI, Analytics, ТО, Антифрод).
 * Замість червоної плашки "error" показує елегантну картку з CTA.
 */
export function PaywallScreen({ feature, description }: PaywallScreenProps) {
  const meta = FEATURES[feature]
  return (
    <div className="paywall-screen">
      <div className="paywall-screen__card">
        <div className="paywall-screen__icon">🔒</div>
        <div className="paywall-screen__title">{meta.label}</div>
        <div className="paywall-screen__tier">
          Доступно на тарифі <strong>{TIER_NAMES[meta.minTier]}</strong> і вище
        </div>
        {description && (
          <p className="paywall-screen__description">{description}</p>
        )}
        <Link to="/billing" className="paywall-screen__cta">
          Переглянути тарифи
        </Link>
      </div>
    </div>
  )
}

interface PaywallBadgeProps {
  feature: FeatureKey
  compact?: boolean
}

/** Маленький значок «PRO» поруч з назвою пункту. */
export function PaywallBadge({ feature, compact }: PaywallBadgeProps) {
  const meta = FEATURES[feature]
  return (
    <span
      className={`paywall-badge paywall-badge--${meta.minTier.toLowerCase()}${compact ? ' paywall-badge--compact' : ''}`}
      title={`Потрібен тариф ${TIER_NAMES[meta.minTier]}`}
    >
      {meta.minTier === 'ENTERPRISE' ? '💎' : '✨'} {meta.minTier}
    </span>
  )
}
