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
    getByWarehouse: (warehouseId: string) => 
      request<InventoryItem[]>(`/inventory/warehouse/${warehouseId}`),

    downloadShipmentPDF: async (shipmentId: string) => {
      const token = getToken();
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
      const token = getToken();
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
      const token = getToken();
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

    smartDispatchPreview: (requestIds: string[]) =>
      request<SmartDispatchResult>('/requests/smart-dispatch-preview', {
        method: 'POST',
        body: JSON.stringify({ request_ids: requestIds }),
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
      const token = getToken();
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
      const token = getToken();
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
    create: (body: { resource_id: string; quantity: number; target_warehouse_id: string }) =>
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
    list: () => request<Unit[]>('/units'),
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

    create: (body: { title: string; description: string; unit_id?: number }) =>
      request<ContractorRequest>('/contractor-requests', {
        method: 'POST',
        body: JSON.stringify(body),
      }),

    take: (id: string) =>
      request<{ message: string }>(`/contractor-requests/${id}/take`, { method: 'POST' }),

    deliver: (id: string) =>
      request<{ message: string }>(`/contractor-requests/${id}/deliver`, { method: 'POST' }),

    accept: (id: string, body: { resource_id?: string; category_id: string; name: string; quantity: number; unit_type: string }) =>
      request<{ message: string }>(`/contractor-requests/${id}/accept`, {
        method: 'POST',
        body: JSON.stringify(body),
    }),
    reject: (id: string) => 
      request<{ message: string }>(`/contractor-requests/${id}/reject`, { method: 'POST' }),

    cancel: (id: string) =>
      request<{ message: string }>(`/contractor-requests/${id}/cancel`, { method: 'POST' }),
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
      const token = localStorage.getItem('token'); 
      const response = await fetch(`/api/vehicles/${vehicleId}/maintenance`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        },
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || 'Помилка збереження акту з файлом');
      }

      return response.json();
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
    getAvailableForRoute: (senderUnitID: number, receiverUnitID: number) => 
      request<Vehicle[]>(`/vehicles/available-for-route?sender_unit_id=${senderUnitID}&receiver_unit_id=${receiverUnitID}`),
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
      const token = getToken();
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

    exportFuel: async (startDate?: string, endDate?: string) => {
      const token = getToken();
      const params = new URLSearchParams();
      if (startDate) params.append('start', startDate);
      if (endDate) params.append('end', endDate);
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
  title: string
  description: string
  status: string
  taken_by?: string
  taken_at?: string
  completed_at?: string
  created_at: string
}

export interface AcceptContractorPayload {
  category_id: string;
  name: string;
  quantity: number;
  unit_type: string;
}

export type UserStatus = 'PENDING' | 'ACTIVE' | 'BLOCKED';

export interface User {
  id: string;
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
  name: string
  description: string
  quantity: number
  serial_number: string
  barcode?: string
  unit_type: 'PCS' | 'KIT' | 'KG' | 'L';
  condition: 'NEW' | 'USED' | 'WRITTEN_OFF';
  min_quantity: number
  assigned_to_user_id?: string;
  assigned_to_user_name?: string;
  weight_kg: number;
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
  barcode?: string
  condition?: 'NEW' | 'USED' | 'WRITTEN_OFF';
  min_quantity: number;
  weight_kg: number;
}
export type RequestStatus = 'PENDING' | 'APPROVED' | 'DISPATCHED' | 'REJECTED' | 'COMPLETED' | 'OPEN';

export interface SupplyRequest {
  id: string
  created_by: string
  resource_id: string
  quantity: number
  status: RequestStatus 
  target_warehouse_id: string 
  approved_by?: string
  comment: string
  created_at: string
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
  created_at: string;
}

export interface DriverHistoryRecord {
  id: string;
  vehicle_id: string;
  driver_id: string | null;
  driver_name: string;
  assigned_at: string;
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
  'ADMIN': 'Системний адміністратор',
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
}

export interface SmartDispatchResult {
  routes: VehicleBin[];
  unassigned: RequestItem[];
}