const API_BASE =
  import.meta.env.VITE_API_URL
    ? `${import.meta.env.VITE_API_URL}/api`
    : '/api'
export const AUTH_SESSION_MARKER = 'omnilog_has_session'
export const SUPPORT_TENANT_KEY = 'omnilog_support_tenant'

// In-memory token store — не зберігається в localStorage, недоступний для XSS
let _inMemoryToken: string | null = null
let _supportTenantId: string | null = localStorage.getItem(SUPPORT_TENANT_KEY)

export function setInMemoryToken(token: string | null) {
  _inMemoryToken = token
}

export function getInMemoryToken(): string | null {
  return _inMemoryToken
}

export function setSupportTenantId(tenantId: string | null) {
  _supportTenantId = tenantId
  if (tenantId) {
    localStorage.setItem(SUPPORT_TENANT_KEY, tenantId)
  } else {
    localStorage.removeItem(SUPPORT_TENANT_KEY)
  }
}

export function getSupportTenantId(): string | null {
  return _supportTenantId
}

export function hasAuthSessionMarker(): boolean {
  return localStorage.getItem(AUTH_SESSION_MARKER) === '1'
}

export function setAuthSessionMarker(enabled: boolean) {
  if (enabled) {
    localStorage.setItem(AUTH_SESSION_MARKER, '1')
  } else {
    localStorage.removeItem(AUTH_SESSION_MARKER)
  }
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
  if (_inMemoryToken) {
    headers['Authorization'] = `Bearer ${_inMemoryToken}`
  }
  if (_supportTenantId && !path.startsWith('/platform') && !path.startsWith('/auth')) {
    headers['X-Support-Tenant-ID'] = _supportTenantId
  }

  const res = await fetch(url, { ...options, headers, credentials: 'include' })
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
    // Refresh використовує httpOnly cookie автоматично (credentials: 'include')
    refresh: () =>
      request<LoginResponse>('/auth/refresh', { method: 'POST', body: '{}' }),
    logout: () =>
      request<{ message: string }>('/auth/logout', { method: 'POST', body: '{}' }),
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
    updatePassword: (data: { old_password: string; new_password: string }) =>
      request<{ message: string }>('/users/password', {
        method: 'PATCH',
        body: JSON.stringify(data),
      }),
    requestPasswordReset: (email: string) =>
      request<{ message: string }>('/auth/forgot-password', {
        method: 'POST',
        body: JSON.stringify({ email }),
      }),
    signupTenant: (body: {
      organization_name: string
      slug: string
      owner_email: string
      owner_full_name: string
      owner_password: string
    }) =>
      request<{ message: string; tenant_id?: string; user_id?: string }>(
        '/auth/tenants/signup',
        { method: 'POST', body: JSON.stringify(body) }
      ),
  },
  platform: {
    stats: () =>
      request<{
        total_tenants: number
        active_tenants: number
        total_users: number
        tenants_by_tier: Record<string, number>
        new_tenants_30_days: number
      }>('/platform/stats'),
    listTenants: (search = '') =>
      request<any[]>(`/platform/tenants${search ? `?search=${encodeURIComponent(search)}` : ''}`),
    getTenant: (id: string) => request<{ tenant: any; user_count: number }>(`/platform/tenants/${id}`),
    getAuditLogs: () => request<any[]>('/platform/audit-logs'),
    createTenant: (body: {
      organization_name: string
      slug: string
      owner_email: string
      owner_full_name: string
      owner_password: string
    }) =>
      request<{ tenant_id: string; owner_id: string; message: string }>('/platform/tenants', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    updateTier: (id: string, tier: string, expires_at?: string | null) =>
      request<{ ok: boolean }>(`/platform/tenants/${id}/tier`, {
        method: 'PATCH',
        body: JSON.stringify({ tier, expires_at: expires_at ?? null }),
      }),
    setActive: (id: string, active: boolean) =>
      request<{ ok: boolean }>(`/platform/tenants/${id}/active`, {
        method: 'PATCH',
        body: JSON.stringify({ active }),
      }),
    deleteTenant: (id: string) =>
      request<{ ok: boolean }>(`/platform/tenants/${id}`, { method: 'DELETE' }),
  },
  admin: {
    createUser: (body: CreateUserRequest) =>
      request<CreateUserResponse>('/admin/users', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    getAuditLogs: () => request<any[]>('/admin/audit-logs'),
    triggerSLACheck: () => request<{
      message: string;
      escalated_count: number;
      existing_escalated?: number;
      pending_total?: number;
      pending_near_sla?: number;
    }>('/admin/sla/trigger', {
      method: 'POST',
    }),
  },
  users: {
    // Змінено з commanders на managers
    listManagers: () => request<User[]>('/users/commanders'),
    getVisible: () => request<User[]>('/users/visible'),
    updateRole: (id: string, role: string, unitId: number | null) => 
      request<{ message: string }>(`/users/${id}/role`, {
        method: 'PUT',
        body: JSON.stringify({ role, unit_id: unitId }),
      }),
    block: (id: string) => request<{ message: string }>(`/users/${id}/block`, { method: 'PUT' }),
    unblock: (id: string) => request<{ message: string }>(`/users/${id}/unblock`, { method: 'PUT' }),
    updateProfile: (data: UpdateProfileData) =>
      request<{ message: string }>('/users/profile', {
        method: 'PATCH',
        body: JSON.stringify(data),
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
    getUniqueResourceNames: (unitId?: number) => 
      request<Array<{ name: string; category_id: string }>>(
        unitId ? `/inventory/resources/unique-names?unit_id=${unitId}` : '/inventory/resources/unique-names'
      ),
    createResource: (body: CreateResourceRequest) =>
      request<Resource>('/inventory/resources', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    getById: (id: string) => 
      request<Resource>(`/inventory/resources/${id}`),
    writeOffResource: (id: string, quantity: number) =>
      request<{ message: string }>(`/inventory/resources/${id}/write-off`, {
        method: 'POST',
        body: JSON.stringify({ quantity }),
      }),
    issueResource: (data: { resource_id: string; user_id: string; quantity: number; notes?: string; warehouse_id?: string }) =>
      request<{ message: string }>('/inventory/issue', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    updateResource: (id: string, data: Partial<Resource>) => 
      request<Resource>(`/inventory/resources/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(data),
      }),
    transferResource: (id: string, body: TransferResourceRequest) =>
      request<{ message: string }>(`/inventory/resources/${id}/transfer`, {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    deleteResource: (id: string) =>
      request<{ message: string }>(`/inventory/resources/${id}`, {
        method: 'DELETE',
      }),

    updateCategory: (id: string, body: { name: string; description?: string }) =>
      request<ResourceCategory>(`/inventory/categories/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(body),
      }),
    deleteCategory: (id: string) =>
      request<{ message: string }>(`/inventory/categories/${id}`, {
        method: 'DELETE',
      }),
    assignResource: (id: string, data: AssignResourceRequest) =>
      request<{ message: string }>(`/inventory/resources/${id}/assign`, {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    getMyEquipment: () => request<MyEquipmentItem[]>('/inventory/my-equipment'),
    reportEquipment: (assignmentId: string, reason: string) =>
      request<{ message: string }>(`/inventory/assignments/${assignmentId}/report`, {
        method: 'POST',
        body: JSON.stringify({ reason }),
      }),
    getByWarehouse: (warehouseId: string) => 
      request<InventoryItem[]>(`/inventory/warehouse/${warehouseId}`),

    downloadShipmentPDF: async (shipmentId: string) => {
      const token = getInMemoryToken();
      const response = await fetch(`${API_BASE}/inventory/shipments/${shipmentId}/pdf`, {
        method: 'GET',
        headers: {
          'Authorization': token ? `Bearer ${token}` : '',
        },
      });

      if (!response.ok) {
        const errData = await response.json().catch(() => ({}));
        throw new Error(errData.error || `Внутрішня помилка сервера: HTTP ${response.status}`);
      }

      const blob = await response.blob();
      
      const contentDisposition = response.headers.get('Content-Disposition');
      let filename = `Waybill_${shipmentId}.pdf`;
      if (contentDisposition) {
        const filenameMatch = contentDisposition.match(/filename="?([^"]+)"?/);
        if (filenameMatch && filenameMatch.length >= 2) {
          filename = filenameMatch[1];
        }
      }

      return { blob, filename };
    },

    downloadResourceQR: async (resourceId: string) => {
      const token = getInMemoryToken();
      const response = await fetch(`${API_BASE}/inventory/resources/${resourceId}/qr`, {
        method: 'GET',
        headers: { 'Authorization': token ? `Bearer ${token}` : '' },
      });

      if (!response.ok) {
        throw new Error('Помилка при генерації QR-коду');
      }

      const blob = await response.blob();
      return blob;
    },

    submitAudit: async (warehouseId: string, discrepancies: AuditDiscrepancy[]) => {
      const token = getInMemoryToken();
      const response = await fetch(`${API_BASE}/inventory/audit`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token ? `Bearer ${token}` : '',
        },
        body: JSON.stringify({
          warehouse_id: warehouseId,
          discrepancies: discrepancies
        })
      });

      if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw new Error(error.error || 'Помилка при збереженні результатів інвентаризації');
      }

      return response.json();
    },

    createShipment: (body: any) =>
      // Зверни увагу: URL на бекенді залишається /inventory/shipments (або можеш змінити його і там)
      // Але на фронті метод викликатиметься через api.requests.createShipment
      request<{ message: string }>('/inventory/shipments', {
        method: 'POST',
        body: JSON.stringify(body),
      }),

    logShipmentRefuel: (shipmentId: string, payload: LogShipmentRefuelPayload) =>
      request<ShipmentRefuel>(`/inventory/shipments/${shipmentId}/refuel`, {
        method: 'POST',
        body: JSON.stringify(payload),
      }),

    getShipmentRefuels: (shipmentId: string) =>
      request<ShipmentRefuel[]>(`/inventory/shipments/${shipmentId}/refuels`),

    smartDispatchPreview: (requestIds: string[], fromWarehouseId?: string) =>
      request<SmartDispatchResult>('/requests/smart-dispatch-preview', {
        method: 'POST',
        body: JSON.stringify({ request_ids: requestIds, from_warehouse_id: fromWarehouseId || '' }),
      }),

    smartDispatchConfirm: (payload: {
      from_warehouse_id: string;
      priority?: string;
      routes: { vehicle_id: string; request_ids: string[] }[];
    }) =>
      request<{ message: string; count: number }>('/requests/smart-dispatch-confirm', {
        method: 'POST',
        body: JSON.stringify(payload),
      }),

    downloadImportTemplate: async () => {
      const token = getInMemoryToken();
      const response = await fetch(`${API_BASE}/inventory/resources/import/template`, {
        method: 'GET',
        headers: { 'Authorization': token ? `Bearer ${token}` : '' },
      });
      if (!response.ok) throw new Error('Помилка завантаження шаблону');
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'OmniLog_Import_Template.xlsx');
      document.body.appendChild(link);
      link.click();
      link.remove();
    },

    importExcel: async (unitId: number, warehouseId: string, file: File) => {
      const token = getInMemoryToken();
      const formData = new FormData();
      formData.append('unit_id', unitId.toString());
      formData.append('warehouse_id', warehouseId);
      formData.append('file', file);

      const response = await fetch(`${API_BASE}/inventory/resources/import`, {
        method: 'POST',
        headers: { 'Authorization': token ? `Bearer ${token}` : '' },
        body: formData,
      });
      const data = await response.json().catch(() => ({}));
      if (!response.ok) throw new Error(data.error || 'Помилка імпорту');
      return data;
    },
  },
  requests: {
    list: () => request<SupplyRequest[]>('/requests'),
    create: (body: { 
      resource_name: string; 
      resource_category_id?: string; 
      quantity: number; 
      target_warehouse_id: string 
    }) =>
      request<SupplyRequest>('/requests', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    approve: (id: string, approved: boolean, comment?: string) =>
      request<{ message: string }>(`/requests/${id}/approve`, {
        method: 'POST',
        body: JSON.stringify({ approved, comment }),
      }),
    reject: (id: string, comment: string) =>
      request<{ message: string }>(`/requests/${id}/reject`, {
        method: 'POST',
        body: JSON.stringify({ comment }),
      }),
      
    cancel: (id: string) =>
      request<{ message: string }>(`/requests/${id}/cancel`, {
        method: 'POST',
      }),
  },
  units: {
    list: () => {
      console.log('🔧 Calling api.units.list()...')
      return request<Unit[]>('/units').then(res => {
        console.log('✅ api.units.list() response:', res)
        return res
      }).catch(err => {
        console.error('❌ api.units.list() failed:', err)
        throw err
      })
    },
    getAvailableForRole: (role: string) => 
      request<Unit[]>(`/units/available?role=${encodeURIComponent(role)}`),

    getMyHierarchyForRole: (role: string) => 
      request<Unit[]>(`/units/my-hierarchy?role=${encodeURIComponent(role)}`),

    create: (body: { parent_id?: number; name: string; unit_type: string }) =>
      request<Unit>('/units', {
        method: 'POST',
        body: JSON.stringify(body),
      }),

    update: (id: number, body: { name: string; unit_type: string; parent_id?: number }) =>
      request<Unit>(`/units/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(body),
      }),
    delete: (id: number) =>
      request<{ message: string }>(`/units/${id}`, {
        method: 'DELETE',
      }),
    // 🔥 Оновлено назву та тіло запиту
    changeManager: (unitId: number, newManagerId: string) => 
      request<{ message: string }>(`/units/${unitId}/change-commander`, {
        method: 'POST',
        body: JSON.stringify({ new_manager_id: newManagerId }),
      }),
  },
  contractorRequests: {
    list: (status?: string) =>
      request<ContractorRequest[]>(`/contractor-requests${status ? `?status=${status}` : ''}`),

    create: (body: { title: string; description: string; unit_id?: number; target_warehouse_id?: string }) =>
      request<ContractorRequest>('/contractor-requests', {
        method: 'POST',
        body: JSON.stringify(body),
      }),

    take: (id: string) =>
      request<{ message: string }>(`/contractor-requests/${id}/take`, { method: 'POST' }),

    deliver: (id: string) =>
      request<{ message: string }>(`/contractor-requests/${id}/deliver`, { method: 'POST' }),

    accept: (id: string, body: AcceptContractorPayload) =>
      request<{ message: string }>(`/contractor-requests/${id}/accept`, {
        method: 'POST',
        body: JSON.stringify(body),
    }),
    reject: (id: string) => 
      request<{ message: string }>(`/contractor-requests/${id}/reject`, { method: 'POST' }),

    cancel: (id: string) =>
      request<{ message: string }>(`/contractor-requests/${id}/cancel`, { method: 'POST' }),
  },

  // Членства підрядників (схвалення співпраці організацією)
  contractorMemberships: {
    // Адмін організації: список підрядників, що подалися/схвалені/відхилені
    list: (status?: string) =>
      request<ContractorMembership[]>(`/contractor-memberships${status ? `?status=${status}` : ''}`),

    approve: (id: string) =>
      request<{ message: string }>(`/contractor-memberships/${id}/approve`, { method: 'POST' }),

    reject: (id: string) =>
      request<{ message: string }>(`/contractor-memberships/${id}/reject`, { method: 'POST' }),

    // Self-view підрядника: з якими організаціями він співпрацює
    mine: () =>
      request<ContractorMembership[]>('/contractor-memberships/mine'),

    // Підрядник самостійно надсилає заявку на співпрацю з організацією
    apply: (tenantId: string) =>
      request<{ status: ContractorMembershipStatus; message: string }>('/contractor-memberships/apply', {
        method: 'POST',
        body: JSON.stringify({ tenant_id: tenantId }),
      }),
  },

  vehicles: {
    list: async (): Promise<Vehicle[]> => {
      return request('/vehicles');
    },
    create: async (data: { brand: string; model: string; plate_number: string; type: string; capacity_kg: number; tank_capacity: number; fuel_norm: number; driver_id?: string }): Promise<Vehicle> => {
      return request('/vehicles', {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },
    updateStatus: async (vehicleId: string, data: { status: string; reason: string }): Promise<void> => {
      return request(`/vehicles/${vehicleId}/status`, {
        method: 'PATCH',
        body: JSON.stringify(data),
      });
    },
    addFuelRecord: async (vehicleId: string, data: { liters: number; odometer_km?: number; record_type: 'REFUEL' | 'EXPENSE' }): Promise<FuelRecord> => {
      return request(`/vehicles/${vehicleId}/fuel`, {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },
    getFuelHistory: async (vehicleId: string): Promise<FuelRecord[]> => {
      return request(`/vehicles/${vehicleId}/fuel`);
    },
    performMaintenance: async (vehicleId: string, formData: FormData): Promise<MaintenanceRecord> => {
      const token = getInMemoryToken(); 
      const response = await fetch(`/api/vehicles/${vehicleId}/maintenance`, {
        method: 'POST',
        headers: {
          'Authorization': token ? `Bearer ${token}` : ''
        },
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || 'Помилка збереження акту з файлом');
      }

      return response.json();
    },
    scheduleMaintenance: async (vehicleId: string, data: { odometer_km: number; service_type: string; scheduled_for: string; description: string }): Promise<MaintenanceRecord> => {
      return request(`/vehicles/${vehicleId}/maintenance/schedule`, {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },
    
    getMaintenanceHistory: async (vehicleId: string): Promise<MaintenanceRecord[]> => {
      return request(`/vehicles/${vehicleId}/maintenance`);
    },
    assignDriver: async (vehicleId: string, driverId: string | null): Promise<void> => {
      return request(`/vehicles/${vehicleId}/driver`, {
        method: 'PATCH',
        body: JSON.stringify({ driver_id: driverId }),
      });
    },
    getDriverHistory: async (vehicleId: string): Promise<DriverHistoryRecord[]> => {
      return request(`/vehicles/${vehicleId}/drivers`);
    },
    getShipmentHistory: async (vehicleId: string): Promise<VehicleShipmentRecord[]> => {
      return request(`/vehicles/${vehicleId}/shipments`);
    },
    update: async (id: string, data: Partial<Vehicle>): Promise<void> => {
      return request(`/vehicles/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(data),
      });
    },
    delete: async (id: string): Promise<void> => {
      return request(`/vehicles/${id}`, {
        method: 'DELETE',
      });
    },
    getAvailableForRoute: (fromWarehouseID: string, toWarehouseID: string) => 
      request<Vehicle[]>(`/vehicles/available-for-route?from_warehouse_id=${fromWarehouseID}&to_warehouse_id=${toWarehouseID}`),
  },
  warehouses: {
    list: async (): Promise<Warehouse[]> => {
      return request('/warehouses');
    },
    create: async (data: CreateWarehouseData): Promise<Warehouse> => {
      return request('/warehouses', {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },
    updateLocation: async (id: string, lat: number, lng: number): Promise<void> => {
      return request(`/warehouses/${id}/location`, {
        method: 'PATCH',
        body: JSON.stringify({ latitude: lat, longitude: lng }),
      });
    },
    update: async (id: string, data: Partial<Warehouse>): Promise<void> => {
      return request(`/warehouses/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(data),
      });
    },
    delete: async (id: string): Promise<void> => {
      return request(`/warehouses/${id}`, {
        method: 'DELETE',
      });
    },
  },
  analytics: {
    getDashboard: (startDate: string, endDate: string, unitId?: string) => {
      const query = new URLSearchParams({ start: startDate, end: endDate });
      if (unitId) query.append('unit_id', unitId);
      return request<any>(`/analytics/dashboard?${query.toString()}`);
    },
    smartReplenish: (items: any[]) => 
      request<{ message: string; count: number }>('/analytics/auto-replenish', {
        method: 'POST',
        body: JSON.stringify({ items }),
      }),
    exportInventory: async (unitId?: number) => {
      const token = getInMemoryToken();
      const query = unitId ? `?unit_id=${unitId}` : '';
      
      const response = await fetch(`${API_BASE}/analytics/export/inventory${query}`, {
        method: 'GET',
        headers: { 'Authorization': token ? `Bearer ${token}` : '' },
      });

      if (!response.ok) {
        const errData = await response.json().catch(() => ({}));
        throw new Error(errData.error || 'Помилка при експорті залишків');
      }

      const blob = await response.blob();
      const contentDisposition = response.headers.get('Content-Disposition');
      let filename = 'Inventory_Report.xlsx';
      if (contentDisposition) {
        const match = contentDisposition.match(/filename="?([^"]+)"?/);
        if (match && match[1]) filename = match[1];
        else {
          const fallbackMatch = contentDisposition.match(/filename=([^;]+)/);
          if (fallbackMatch && fallbackMatch[1]) filename = fallbackMatch[1];
        }
      }
      return { blob, filename };
    },

    exportFuel: async (startDate?: string, endDate?: string, unitId?: number) => {
      const token = getInMemoryToken();
      const params = new URLSearchParams();
      if (startDate) params.append('start', startDate);
      if (endDate) params.append('end', endDate);
      if (unitId) params.append('unit_id', String(unitId));
      const query = params.toString() ? `?${params.toString()}` : '';

      const response = await fetch(`${API_BASE}/analytics/export/fuel${query}`, {
        method: 'GET',
        headers: { 'Authorization': token ? `Bearer ${token}` : '' },
      });

      if (!response.ok) {
        const errData = await response.json().catch(() => ({}));
        throw new Error(errData.error || 'Помилка при експорті пального');
      }

      const blob = await response.blob();
      const contentDisposition = response.headers.get('Content-Disposition');
      let filename = 'Fuel_Report.xlsx';
      if (contentDisposition) {
        const match = contentDisposition.match(/filename="?([^"]+)"?/);
        if (match && match[1]) filename = match[1];
        else {
          const fallbackMatch = contentDisposition.match(/filename=([^;]+)/);
          if (fallbackMatch && fallbackMatch[1]) filename = fallbackMatch[1];
        }
      }
      return { blob, filename };
    },
    getAdvancedKPIs: () => request<any>('/analytics/kpi'),
    getDemandForecast: () => request<any>('/analytics/forecast'),
    getPredictiveMaintenanceSchedule: (vehicleId?: string) => {
      const query = vehicleId ? `?vehicle_id=${vehicleId}` : '';
      return request<any>(`/analytics/maintenance${query}`);
    },
    getFuelAnomalyDetection: (vehicleId?: string) => {
      const query = vehicleId ? `?vehicle_id=${vehicleId}` : '';
      return request<any>(`/analytics/fuel-anomalies${query}`);
    },
  },

  gps: {
    recordLocation: (data: any) =>
      request<any>('/gps/locations', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    getFleetMap: () => request<any>('/gps/fleet-map'),
    getVehicleTrajectory: (vehicleId: string, startTime: string, endTime: string) => {
      const query = new URLSearchParams({ vehicle_id: vehicleId, start_time: startTime, end_time: endTime });
      return request<any>(`/gps/trajectory?${query.toString()}`);
    },
    createGeofence: (data: any) =>
      request<any>('/gps/geofences', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    getGeofences: () => request<any>('/gps/geofences'),
    getGeofenceAlerts: (hours: number = 24) =>
      request<any>(`/gps/geofence-alerts?hours=${hours}`),
    getFleetStatus: () => request<any>('/gps/fleet-status'),
    driverPing: (data: {
      latitude: number; longitude: number; altitude: number;
      speed: number; heading: number; accuracy: number;
    }) => request<any>('/gps/driver/ping', { method: 'POST', body: JSON.stringify(data) }),
    getDriverActiveShipment: () => request<any>('/gps/driver/active-shipment'),
  },

  notifications: {
    list: (limit: number = 50) => 
      request<NotificationListResponse>(`/notifications?limit=${limit}`),
    getUnreadCount: () => 
      request<{ unread_count: number }>('/notifications/unread-count'),
    markAsRead: (id: string) =>
      request<{ message: string }>(`/notifications/${id}/read`, { method: 'PATCH' }),
    markAllAsRead: () =>
      request<{ message: string }>('/notifications/mark-all-read', { method: 'POST' }),
    delete: (id: string) =>
      request<{ message: string }>(`/notifications/${id}`, { method: 'DELETE' }),
  }
}

export interface MyEquipmentItem {
  assignment_id: string;
  resource_id: string;
  resource_name: string;
  quantity: number;
  unit_type: string;
  issued_at: string;
  status: string;
}

export interface Warehouse {
  id: string;
  unit_id: number;
  name: string;
  location_type: 'STATIONARY' | 'MOBILE';
  latitude?: number;
  longitude?: number;
  capacity_level?: 'LARGE' | 'MEDIUM' | 'SMALL';
  zone_type?: 'REAR' | 'TACTICAL' | 'FORWARD';
  created_at: string;
}

export interface CreateWarehouseData {
  unit_id: number;
  name: string;
  location_type: 'STATIONARY' | 'MOBILE';
  latitude?: number;
  longitude?: number;
}

export type SubscriptionTier = 'BASIC' | 'PRO' | 'ENTERPRISE';  

// 🔥 Оновлено інтерфейс
export interface Unit {
  id: number
  parent_id?: number
  name: string
  unit_type: 'REGION' | 'BRANCH' | 'DEPARTMENT' | 'TEAM'; 
  subscription_tier: SubscriptionTier;
}

export interface ContractorRequest {
  id: string
  created_by: string
  unit_id?: number;
  unit_name?: string;
  target_warehouse_id?: string;
  warehouse_name?: string;
  title: string
  description: string
  status: string
  taken_by?: string
  taken_at?: string
  completed_at?: string
  created_at: string
  tenant_id?: string;
  tenant_name?: string;
}

export type ContractorMembershipStatus = 'PENDING' | 'APPROVED' | 'REJECTED';

export interface ContractorMembership {
  id: string;
  contractor_id: string;
  tenant_id: string;
  status: ContractorMembershipStatus;
  note?: string | null;
  requested_at: string;
  decided_at?: string | null;
  decided_by?: string | null;
  contractor_name?: string;
  contractor_email?: string;
  contractor_phone?: string | null;
  tenant_name?: string;
}

export interface AcceptContractorPayload {
  resource_id?: string;
  category_id?: string;
  category_name?: string;
  name?: string;
  quantity: number;
  unit_type: string;
  unit_price?: number;
}

export type UserStatus = 'PENDING' | 'ACTIVE' | 'BLOCKED';

export interface User {
  id: string;
  tenant_id?: string | null;     
  username?: string | null;      
  email: string;
  full_name: string;
  phone?: string | null;         
  role: UserRole;                
  status: UserStatus;            
  unit_id?: number | null;       
  created_at: string;            
  updated_at: string;
  effective_subscription_tier?: SubscriptionTier;
  subscription_expires_at?: string | null;  // ISO-8601, null = безстрокова
  unit?: {
    subscription_tier?: string;
  };
}

export interface UpdateProfileData {
  full_name?: string;
  phone?: string;
  username?: string;
  email?: string;
}

export type SystemUser = User;

export interface LoginResponse {
  token: string
  refresh_token: string
  expires_at: number
  user: User
}

export type UserRole =
  | 'SYSTEM_ADMIN'
  | 'TENANT_ADMIN'
  | 'ADMIN'
  | 'REGION_DIRECTOR'
  | 'BRANCH_MANAGER'
  | 'DEPT_MANAGER'
  | 'TEAM_LEAD'
  | 'REGION_LOGISTICIAN'
  | 'REGION_STOREKEEPER'
  | 'BRANCH_LOGISTICIAN'
  | 'BRANCH_STOREKEEPER'
  | 'DEPT_SUPERVISOR'
  | 'CONTRACTOR'
  | 'EMPLOYEE'

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
  warehouse_id?: string;
  warehouse_name?: string; // 🆕 Назва складу для відображення
  name: string
  description: string
  quantity: number
  serial_number: string
  barcode?: string
  location?: string;       // Місце зберігання (полиця, секція тощо)
  unit_type: 'PCS' | 'KIT' | 'KG' | 'L';
  condition: 'NEW' | 'USED' | 'WRITTEN_OFF';
  min_quantity: number
  assigned_to_user_id?: string;
  assigned_to_user_name?: string;
  weight_kg: number;
  unit_price: number;
  issued_quantity?: number;
}
export interface AssignResourceRequest {
  quantity: number;
  user_id: string;
}

export interface CreateResourceRequest {
  category_id: string;
  unit_id?: number;
  warehouse_id?: string;
  name: string;
  description?: string;
  quantity: number;
  unit_type: 'PCS' | 'KIT' | 'KG' | 'L';
  serial_number?: string;
  barcode?: string;
  location?: string;
  condition?: 'NEW' | 'USED' | 'WRITTEN_OFF';
  min_quantity: number;
  weight_kg: number;
  unit_price?: number;
}
export type RequestStatus = 'PENDING' | 'APPROVED' | 'LOADING' | 'DISPATCHED' | 'REJECTED' | 'COMPLETED' | 'OPEN';

export interface SupplyRequest {
  id: string
  created_by: string
  resource_id?: string | null  // Тепер nullable
  resource_name: string  // Нове поле - назва ресурсу
  resource_category_id?: string | null  // Нове поле - категорія ресурсу
  quantity: number
  status: RequestStatus 
  target_warehouse_id: string 
  approved_by?: string
  approved_at?: string
  comment?: string
  created_at: string
  updated_at?: string
}

export type FuelRecordType = 'REFUEL' | 'EXPENSE';

export type VehicleType = 'PICKUP' | 'VAN' | 'TRUCK';

export type VehicleStatus = 'ACTIVE' | 'INACTIVE' | 'IN_REPAIR' | 'ON_MISSION';

export interface Vehicle {
  id: string;
  brand: string;
  model: string;
  plate_number: string;
  status: VehicleStatus;
  tank_capacity: number;
  fuel_norm: number;
  maintenance_interval_km: number;
  last_maintenance_odometer: number;
  current_odometer: number;
  km_to_next_maintenance: number;
  maintenance_status: 'OK' | 'WARNING' | 'OVERDUE';
  driver_id?: string;
  type: VehicleType;
  capacity_kg: number;
  driver_name?: string;
  avg_km_per_day?: number;
  predicted_maint_date?: string;
  home_warehouse_id?: string;
  home_warehouse_name?: string;
  current_warehouse_id?: string;
  current_warehouse_name?: string;
  current_fuel_liters?: number;
}

export interface FuelRecord {
  id: string;
  vehicle_id: string;
  liters: number;
  odometer_km?: number;
  record_type: FuelRecordType;
  is_anomaly: boolean;
  anomaly_reason?: string;
  created_by?: string;
  created_at: string;
}

export interface MaintenanceRecord {
  id: string;
  vehicle_id: string;
  odometer_km: number;
  description: string;
  performed_by?: string;
  cost_amount: number;
  document_url?: string;
  driver_name?: string;
  status: 'COMPLETED' | 'SCHEDULED';
  service_type: string;
  scheduled_for?: string;
  created_at: string;
}

export interface DriverHistoryRecord {
  id: string;
  vehicle_id: string;
  driver_id: string | null;
  driver_name: string;
  assigned_at: string;
}

export interface VehicleShipmentRecord {
  id: string;
  vehicle_id: string;
  vehicle_plate: string;
  from_warehouse_id: string;
  from_warehouse_name: string;
  to_warehouse_id: string;
  to_warehouse_name: string;
  status: 'PENDING' | 'IN_TRANSIT' | 'DELIVERED';
  priority: string;
  direction: string;
  distance_km: number;
  actual_km: number;
  created_at: string;
  started_at?: string;
  delivered_at?: string;
}

export interface TransferResourceRequest {
  quantity: number;
  target_warehouse_id?: string;
  target_unit_id?: number;
}

export interface InventoryItem {
  id: string;
  warehouse_id: string;
  name: string;
  category: string;
  available: number;
  weight_kg: number;
}

export interface ShipmentItemPayload {
  resource_id: string;
  quantity: number;
  request_id?: string;
}

export interface CreateShipmentPayload {
  from_warehouse_id: string;
  to_warehouse_id: string;
  vehicle_id: string;
  priority: 'NORMAL' | 'URGENT';
  items: ShipmentItemPayload[];
}

export interface ShipmentRecord {
  id: string;
  from_warehouse: string;
  to_warehouse: string;
  vehicle: string;
  priority: 'NORMAL' | 'URGENT';
  status: 'DISPATCHED' | 'DELIVERED';
  created_at: string;
}

export interface ShipmentRefuel {
  id: string;
  shipment_id: string;
  vehicle_id: string;
  liters: number;
  odometer_km?: number;
  station_name?: string;
  cost_uah?: number;
  created_by?: string;
  created_at: string;
}

export interface LogShipmentRefuelPayload {
  liters: number;
  odometer_km?: number;
  station_name?: string;
  cost_uah?: number;
}

export interface AuditLog {
  id: string;
  user_id: string;
  action: string;
  entity_type: string;
  entity_id: string;
  details: string;
  created_at: string;
}

export interface AuditDiscrepancy {
  resource_id: string;
  book_quantity: number;
  actual_quantity: number;
  difference: number;
}

export interface SubmitAuditRequest {
  warehouse_id: string;
  discrepancies: AuditDiscrepancy[];
}

// ==========================================
// 🔥 СЛОВНИКИ ДЛЯ КРАСИВОГО UI
// ==========================================

export const ROLE_NAMES: Record<UserRole, string> = {
  'SYSTEM_ADMIN': 'Власник платформи',
  'TENANT_ADMIN': 'Адміністратор організації',
  'ADMIN': 'Адміністратор організації (застаріле)',
  'REGION_DIRECTOR': 'Директор регіону',
  'BRANCH_MANAGER': 'Керівник філії',
  'DEPT_MANAGER': 'Начальник відділу',
  'TEAM_LEAD': 'Керівник групи (Тімлід)',
  'REGION_LOGISTICIAN': 'Регіональний логіст',
  'REGION_STOREKEEPER': 'Завідувач рег. складом',
  'BRANCH_LOGISTICIAN': 'Логіст філії',
  'BRANCH_STOREKEEPER': 'Завідувач складом філії',
  'DEPT_SUPERVISOR': 'Супервайзер відділу',
  'CONTRACTOR': 'Підрядник (Зовнішній)',
  'EMPLOYEE': 'Співробітник'
};

export const UNIT_TYPE_NAMES: Record<string, string> = {
  'REGION': 'Регіон / Дирекція',
  'BRANCH': 'Філія',
  'DEPARTMENT': 'Відділ',
  'TEAM': 'Команда / Група'
};

export interface RequestItem {
  id: string;
  name: string;
  weight_kg: number;
}

export interface VehicleBin {
  id: string;
  name: string;
  max_weight: number;
  used_weight: number;
  items: RequestItem[];
  fuel_liters?: number;
  fuel_norm?: number;
  tank_capacity?: number;
}

export interface SmartDispatchResult {
  routes: VehicleBin[];
  unassigned: RequestItem[];
}

// ==========================================
// NOTIFICATIONS (Сповіщення)
// ==========================================

export type NotificationType = 
  | 'SHIPMENT_ASSIGNED' 
  | 'REQUEST_APPROVED' 
  | 'REQUEST_REJECTED' 
  | 'SHIPMENT_DELIVERED'
  | 'LOW_STOCK';

export interface Notification {
  id: string;
  user_id: string;
  tenant_id: string;
  type: NotificationType;
  title: string;
  message: string;
  related_id?: string;
  is_read: boolean;
  created_at: string;
  read_at?: string;
}

export interface NotificationListResponse {
  notifications: Notification[];
  unread_count: number;
}
