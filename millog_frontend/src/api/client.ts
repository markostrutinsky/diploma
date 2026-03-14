const API_BASE = '/api'

function getToken(): string | null {
  return localStorage.getItem('token')
}

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE}${path}`
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  const token = getToken()
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(url, { ...options, headers })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error(data.error || `HTTP ${res.status}`)
  }
  return data as T
}

export const api = {
  auth: {
    login: (email: string, password: string) =>
      request<LoginResponse>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      }),
    refresh: (refreshToken: string) =>
      request<LoginResponse>('/auth/refresh', {
        method: 'POST',
        body: JSON.stringify({ refresh_token: refreshToken }),
      }),
    register: (body: { email: string; password: string; full_name: string }) =>
      request<{ id: string; message: string }>('/auth/register', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    setupPassword: (token: string, password: string) =>
      request<{ message: string }>('/auth/setup-password', {
        method: 'POST',
        body: JSON.stringify({ token, password }),
      }),
    me: () => request<User>('/auth/me'),
    bootstrap: (email: string, password: string, full_name: string) =>
      request<{ message: string }>('/bootstrap', {
        method: 'POST',
        body: JSON.stringify({ email, password, full_name }),
      }),
  },
  admin: {
    createUser: (body: CreateUserRequest) =>
      request<CreateUserResponse>('/admin/users', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
  },
  inventory: {
    listCategories: () => request<ResourceCategory[]>('/inventory/categories'),
    createCategory: (body: { name: string; description?: string }) =>
      request<ResourceCategory>('/inventory/categories', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    listResources: (unitId?: number) =>
      request<Resource[]>(unitId ? `/inventory/resources?unit_id=${unitId}` : '/inventory/resources'),
    createResource: (body: CreateResourceRequest) =>
      request<Resource>('/inventory/resources', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
  },
  requests: {
    list: () => request<SupplyRequest[]>('/requests'),
    create: (body: { resource_id: string; quantity: number }) =>
      request<SupplyRequest>('/requests', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    approve: (id: string, approved: boolean, comment?: string) =>
      request<{ message: string }>(`/requests/${id}/approve`, {
        method: 'POST',
        body: JSON.stringify({ approved, comment }),
      }),
  },
  units: {
    list: () => request<Unit[]>('/units'),
    create: (body: { parent_id?: number; name: string; unit_type: string }) =>
      request<Unit>('/units', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
  },
  volunteerRequests: {
    list: () => request<VolunteerRequest[]>('/volunteer-requests'),
    create: (body: { title: string; description?: string }) =>
      request<VolunteerRequest>('/volunteer-requests', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    take: (id: string) =>
      request<{ message: string }>(`/volunteer-requests/${id}/take`, { method: 'POST' }),
    complete: (id: string) =>
      request<{ message: string }>(`/volunteer-requests/${id}/complete`, { method: 'POST' }),
  },
}

export interface Unit {
  id: number
  parent_id?: number
  name: string
  unit_type: string
}

export interface VolunteerRequest {
  id: string
  created_by: string
  title: string
  description: string
  status: string
  taken_by?: string
  taken_at?: string
  completed_at?: string
  created_at: string
}

export interface User {
  id: string
  email: string
  full_name: string
  role: string
  status: string
}

export interface LoginResponse {
  token: string
  refresh_token: string
  expires_at: number
  user: User
}

export type UserRole =
  | 'ADMIN'
  | 'BRIGADE_CMDR'
  | 'BATTALION_CMDR'
  | 'COMPANY_CMDR'
  | 'PLATOON_CMDR'
  | 'BRIGADE_LOGIST'
  | 'BRIGADE_STOREKEEPER'
  | 'BATTALION_LOGIST'
  | 'BATTALION_STOREKEEPER'
  | 'COMPANY_SERGEANT'
  | 'VOLUNTEER'

export interface CreateUserRequest {
  username?: string
  email: string
  full_name: string
  phone?: string
  role: UserRole
  unit_id?: number
}

export interface CreateUserResponse {
  id: string
  message: string
}

export interface ResourceCategory {
  id: string
  name: string
  description: string
}

export interface Resource {
  id: string
  category_id: string
  unit_id?: number
  name: string
  description: string
  quantity: number
  serial_number: string
  location: string
  condition: string
  min_quantity: number
}

export interface CreateResourceRequest {
  category_id: string
  unit_id?: number
  name: string
  description?: string
  quantity?: number
  serial_number?: string
  location?: string
  condition?: string
  min_quantity?: number
}

export interface SupplyRequest {
  id: string
  created_by: string
  resource_id: string
  quantity: number
  status: string
  approved_by?: string
  comment: string
  created_at: string
}
