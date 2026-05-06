// Централізовані правила доступу системи Omnilog.
// Тут живуть: перелік ролей, групи (white-lists), дії (permissions),
// а також матриця платних фіч (feature gating за підпискою).
// Використовується через хук `usePermissions` і компонент `FeatureGate`.

export const ROLES = {
  SYSTEM_ADMIN: 'SYSTEM_ADMIN',
  TENANT_ADMIN: 'TENANT_ADMIN',
  ADMIN: 'ADMIN',
  REGION_DIRECTOR: 'REGION_DIRECTOR',
  REGION_LOGISTICIAN: 'REGION_LOGISTICIAN',
  REGION_STOREKEEPER: 'REGION_STOREKEEPER',
  BRANCH_MANAGER: 'BRANCH_MANAGER',
  BRANCH_LOGISTICIAN: 'BRANCH_LOGISTICIAN',
  BRANCH_STOREKEEPER: 'BRANCH_STOREKEEPER',
  DEPT_MANAGER: 'DEPT_MANAGER',
  DEPT_SUPERVISOR: 'DEPT_SUPERVISOR',
  TEAM_LEAD: 'TEAM_LEAD',
  EMPLOYEE: 'EMPLOYEE',
  CONTRACTOR: 'CONTRACTOR',
} as const

export type Role = typeof ROLES[keyof typeof ROLES]

export type Tier = 'BASIC' | 'PRO' | 'ENTERPRISE'

// ================================================================
// 1. ГРУПИ ДОСТУПУ ДО РОЗДІЛІВ (роути у App.tsx + меню у Layout)
// ================================================================

export const ROLE_GROUPS = {
  analytics: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.REGION_STOREKEEPER,
  ],
  inventory: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.REGION_STOREKEEPER,
    ROLES.BRANCH_MANAGER,
    ROLES.BRANCH_STOREKEEPER,
    ROLES.DEPT_MANAGER,
    ROLES.DEPT_SUPERVISOR,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
  ],
  transport: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_MANAGER,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.DEPT_SUPERVISOR,
  ],
  units: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.DEPT_MANAGER,
    ROLES.REGION_LOGISTICIAN,
  ],
  users: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.DEPT_MANAGER,
    ROLES.TEAM_LEAD,
  ],
  kiosk: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.REGION_LOGISTICIAN,
    ROLES.REGION_STOREKEEPER,
    ROLES.BRANCH_MANAGER,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.BRANCH_STOREKEEPER,
    ROLES.DEPT_MANAGER,
    ROLES.DEPT_SUPERVISOR,
    ROLES.TEAM_LEAD,
    ROLES.EMPLOYEE,
  ],
  contracts: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.DEPT_MANAGER,
  ],
  approvers: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.DEPT_MANAGER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
  ],
  requests: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.REGION_LOGISTICIAN,
    ROLES.REGION_STOREKEEPER,
    ROLES.BRANCH_MANAGER,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.BRANCH_STOREKEEPER,
    ROLES.DEPT_MANAGER,
    ROLES.DEPT_SUPERVISOR,
    ROLES.TEAM_LEAD,
    // EMPLOYEE прибрано - вони працюють з рейсами, а не заявками
  ],
  contractorRequestsView: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.DEPT_MANAGER,
    ROLES.CONTRACTOR,
  ],
  superAdmin: [ROLES.SYSTEM_ADMIN, ROLES.TENANT_ADMIN, ROLES.ADMIN],
  platform: [ROLES.SYSTEM_ADMIN],
} as const

// ================================================================
// 2. МАТРИЦЯ ДІЙ (ACTIONS) — хто що може робити
// ================================================================

export const ACTIONS = {
  // Ресурси
  resource_manage: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_STOREKEEPER,
    ROLES.BRANCH_STOREKEEPER,
    ROLES.DEPT_SUPERVISOR,
    ROLES.REGION_LOGISTICIAN,
  ],
  resource_assign: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_STOREKEEPER,
    ROLES.BRANCH_STOREKEEPER,
    ROLES.DEPT_SUPERVISOR,
  ],
  resource_writeoff: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_STOREKEEPER,
    ROLES.BRANCH_STOREKEEPER,
    ROLES.REGION_LOGISTICIAN,
  ],
  resource_transfer: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_STOREKEEPER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_STOREKEEPER,
  ],

  // Категорії
  category_manage: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_LOGISTICIAN,
    ROLES.REGION_DIRECTOR,
  ],

  // Заявки на постачання
  request_create: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.DEPT_MANAGER,
    ROLES.TEAM_LEAD,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.DEPT_SUPERVISOR,
    ROLES.EMPLOYEE,
  ],
  request_approve: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.DEPT_MANAGER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
  ],
  request_dispatch: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_STOREKEEPER,
    ROLES.BRANCH_STOREKEEPER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
  ],

  // Склади
  warehouse_manage: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_MANAGER,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.DEPT_MANAGER,
  ],
  warehouse_audit: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_STOREKEEPER,
    ROLES.BRANCH_STOREKEEPER,
    ROLES.REGION_LOGISTICIAN,
  ],

  // Автопарк
  vehicle_manage: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.DEPT_SUPERVISOR,
  ],
  vehicle_fuel_log: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.DEPT_SUPERVISOR,
  ],
  vehicle_maintenance: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
  ],

  // Організаційна структура
  unit_manage: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.DEPT_MANAGER,
    ROLES.REGION_LOGISTICIAN,
  ],

  // Персонал
  user_invite: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.DEPT_MANAGER,
    ROLES.TEAM_LEAD,
  ],
  user_block: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
  ],
  user_change_role: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
  ],

  // Заявки підрядникам
  contractor_request_create: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.DEPT_MANAGER,
  ],
  contractor_request_take: [ROLES.CONTRACTOR],

  // Адмін-функції
  audit_view: [ROLES.SYSTEM_ADMIN, ROLES.TENANT_ADMIN, ROLES.ADMIN],
  billing_manage: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
  ],
  sla_trigger: [ROLES.SYSTEM_ADMIN, ROLES.TENANT_ADMIN, ROLES.ADMIN],

  // Термінал
  kiosk_operate: [
    ROLES.SYSTEM_ADMIN,
    ROLES.TENANT_ADMIN,
    ROLES.ADMIN,
    ROLES.REGION_STOREKEEPER,
    ROLES.BRANCH_STOREKEEPER,
    ROLES.DEPT_SUPERVISOR,
  ],
} as const

export type ActionKey = keyof typeof ACTIONS

// ================================================================
// 3. ПЛАТНІ ФІЧІ (FEATURE GATING)
// ================================================================

export const FEATURES = {
  smart_dispatch: { minTier: 'PRO' as Tier, label: 'Smart Розподіл рейсів' },
  smart_replenish: { minTier: 'PRO' as Tier, label: 'Smart Поповнення складу' },
  advanced_analytics: { minTier: 'PRO' as Tier, label: 'Розширена аналітика (SLA, TCO, ризики)' },
  fuel_antifraud: { minTier: 'PRO' as Tier, label: 'Антифрод-система пального' },
  predictive_maintenance: { minTier: 'PRO' as Tier, label: 'Прогноз ТО за пробігом' },
  gps_tracking: { minTier: 'PRO' as Tier, label: 'GPS-трекінг автопарку (real-time мапа)' },
  excel_import: { minTier: 'PRO' as Tier, label: 'Масовий імпорт з Excel' },
  audit_log_extended: { minTier: 'PRO' as Tier, label: 'Журнал аудиту понад 7 днів' },
  // Примітка: «Пріоритетна підтримка» та «Мульти-регіональна консолідація»
  // лишились тільки як маркетинговий текст на сторінці Billing — у коді
  // жодного окремого gate'а під них немає, тому тут ми їх не декларуємо,
  // аби hasFeature() не зміг повернути false на фічу, якої реально немає.
} as const

export type FeatureKey = keyof typeof FEATURES

const TIER_WEIGHT: Record<Tier, number> = {
  BASIC: 0,
  PRO: 1,
  ENTERPRISE: 2,
}

export const TIER_NAMES: Record<Tier, string> = {
  BASIC: 'Standard',
  PRO: 'Enterprise PRO',
  ENTERPRISE: 'Enterprise+',
}

// ================================================================
// 4. HELPERS
// ================================================================

export function hasRole(userRole: string | undefined, allowed: readonly string[]): boolean {
  if (!userRole) return false
  return allowed.includes(userRole)
}

export function can(userRole: string | undefined, action: ActionKey): boolean {
  if (!userRole) return false
  const allowed = ACTIONS[action] as readonly string[]
  return allowed.includes(userRole)
}

export function tierAtLeast(userTier: Tier | undefined | null, required: Tier): boolean {
  if (!userTier) return false
  return TIER_WEIGHT[userTier] >= TIER_WEIGHT[required]
}

export function hasFeature(
  userTier: Tier | undefined | null,
  userRole: string | undefined,
  feature: FeatureKey,
): boolean {
  // SYSTEM_ADMIN, TENANT_ADMIN та legacy ADMIN завжди мають доступ до всіх фіч
  if (
    userRole === ROLES.SYSTEM_ADMIN ||
    userRole === ROLES.TENANT_ADMIN ||
    userRole === ROLES.ADMIN
  ) return true
  const { minTier } = FEATURES[feature]
  return tierAtLeast(userTier, minTier)
}
