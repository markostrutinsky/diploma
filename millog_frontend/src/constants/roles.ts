// Централізовані групи ролей для керування доступом.
// Використовується як у бічній панелі (Layout), так і в маршрутах (ProtectedRoute),
// щоб правила доступу залишались синхронізованими.

export const ROLES = {
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

// --- Групи доступу до розділів системи ---

export const ROLE_GROUPS = {
  // Аналітика (керівники + логісти регіонального рівня + регіональний комірник)
  analytics: [
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.REGION_STOREKEEPER,
  ],

  // Складський блок: Ресурси, Склади
  inventory: [
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.REGION_STOREKEEPER,
    ROLES.BRANCH_MANAGER,
    ROLES.BRANCH_STOREKEEPER,
    ROLES.DEPT_MANAGER,
    ROLES.DEPT_SUPERVISOR,
  ],

  // Автопарк
  transport: [
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_MANAGER,
    ROLES.BRANCH_LOGISTICIAN,
  ],

  // Оргструктура (підрозділи)
  units: [
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.DEPT_MANAGER,
  ],

  // Керування користувачами
  users: [
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.DEPT_MANAGER,
    ROLES.TEAM_LEAD,
  ],

  // Термінал/Каса (видача зі складу)
  kiosk: [
    ROLES.ADMIN,
    ROLES.REGION_STOREKEEPER,
    ROLES.BRANCH_STOREKEEPER,
    ROLES.DEPT_SUPERVISOR,
  ],

  // Управління контрактами (заявки підрядникам)
  contracts: [
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.DEPT_MANAGER,
  ],

  // Хто може погоджувати внутрішні заявки
  approvers: [
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.DEPT_MANAGER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
  ],

  // Внутрішні заявки (усі, крім підрядників)
  requests: [
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

  // Заявки підрядникам видно самим підрядникам + менеджерам з групи contracts
  contractorRequestsView: [
    ROLES.ADMIN,
    ROLES.REGION_DIRECTOR,
    ROLES.BRANCH_MANAGER,
    ROLES.REGION_LOGISTICIAN,
    ROLES.BRANCH_LOGISTICIAN,
    ROLES.DEPT_MANAGER,
    ROLES.CONTRACTOR,
  ],

  // Тільки суперадмін
  superAdmin: [ROLES.ADMIN],
} as const

export function hasRole(userRole: string | undefined, allowed: readonly string[]): boolean {
  if (!userRole) return false
  return allowed.includes(userRole)
}
