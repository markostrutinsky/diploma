import { useMemo } from 'react'
import { useAuth } from '../contexts/AuthContext'
import {
  ACTIONS,
  FEATURES,
  ROLE_GROUPS,
  ROLES,
  type ActionKey,
  type FeatureKey,
  type Tier,
  can as canFn,
  hasFeature as hasFeatureFn,
  hasRole as hasRoleFn,
  tierAtLeast as tierAtLeastFn,
  isSubscriptionActive,
} from '../constants/roles'

/**
 * Єдиний хук для перевірки прав доступу в UI.
 *
 * Приклад:
 * ```tsx
 * const perms = usePermissions()
 * if (perms.can('resource_manage')) { ... }
 * if (perms.hasFeature('smart_dispatch')) { ... }
 * ```
 */
export function usePermissions() {
  const { user, loading: authLoading } = useAuth()

  return useMemo(() => {
    const role = user?.role
    const tier = (user?.effective_subscription_tier || user?.unit?.subscription_tier) as
      | Tier
      | undefined
    const expiresAt = user?.subscription_expires_at

    const isAuthenticated = !!user
    const isAdmin = role === ROLES.ADMIN
    const isContractor = role === ROLES.CONTRACTOR
    const isManager = hasRoleFn(role, [
      ROLES.ADMIN,
      ROLES.REGION_DIRECTOR,
      ROLES.BRANCH_MANAGER,
      ROLES.DEPT_MANAGER,
      ROLES.TEAM_LEAD,
    ])

    /** true якщо підписка ще діє (null = безстрокова) */
    const subscriptionActive = isSubscriptionActive(expiresAt)

    return {
      user,
      role,
      tier,
      expiresAt,
      subscriptionActive,
      /** true поки AuthContext ще відновлює сесію — використовуй для затримки feature-gate */
      authLoading,
      isAuthenticated,
      isAdmin,
      isContractor,
      isManager,

      /** Чи належить користувач до однієї з груп з ROLE_GROUPS */
      inGroup: (group: keyof typeof ROLE_GROUPS) =>
        hasRoleFn(role, ROLE_GROUPS[group]),

      /** Чи дозволено користувачеві виконати дію (дивись ACTIONS) */
      can: (action: ActionKey) => canFn(role, action),

      /** Чи можна виконати хоч одну з переданих дій */
      canAny: (...actions: ActionKey[]) => actions.some((a) => canFn(role, a)),

      /** Чи відкрита платна фіча поточній підписці (враховує expiry) */
      hasFeature: (feature: FeatureKey) => hasFeatureFn(tier, role, feature, expiresAt),

      /** Чи тариф не нижче переданого */
      tierAtLeast: (t: Tier) => tierAtLeastFn(tier, t),

      /** Метадані про платну фічу — для Paywall-банерів */
      featureMeta: (feature: FeatureKey) => FEATURES[feature],

      /** Явний перелік — може знадобитись для серверних payload-ів */
      allActions: ACTIONS,
    }
  }, [user, authLoading])
}

export type PermissionsHook = ReturnType<typeof usePermissions>
