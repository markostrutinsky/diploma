import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import Layout from './components/Layout'
import ProtectedRoute from './components/ProtectedRoute'
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import SignupTenant from './pages/SignupTenant'
import PlatformAdmin from './pages/PlatformAdmin'
import Bootstrap from './pages/Bootstrap'
import SetupPassword from './pages/SetupPassword'
import AdminUsers from './pages/AdminUsers'
import Inventory from './pages/Inventory'
import Requests from './pages/Requests'
import VolunteerRequests from './pages/VolunteerRequests'
import Units from './pages/Units'
import Vehicles from './pages/Vehicles'
import Warehouses from './pages/Warehouses'
import AnalyticsDashboard from './pages/AnalyticsDashboard'
import Profile from './pages/Profile'
import AuditLogs from './pages/AuditLogs'
import Billing from './pages/Billing'
import KioskTerminal from './pages/KioskTerminal'
import KPIDashboard from './pages/KPIDashboard'
import GPSTracking from './pages/GPSTracking'
import MaintenanceSchedule from './pages/MaintenanceSchedule'
import FuelAnomalies from './pages/FuelAnomalies'
import MyShipments from './pages/MyShipments'
import { Toaster } from 'react-hot-toast'
import { ROLE_GROUPS } from './constants/roles'
import { GPSProvider } from './contexts/GPSContext'

function App() {
  return (
    <AuthProvider>
      <GPSProvider>
      <Toaster 
        position="top-right"
        containerStyle={{ zIndex: 9999 }}
        toastOptions={{
          duration: 4000,
          style: {
            background: '#1e293b',
            color: '#f1f5f9',
            border: '1px solid #334155',
            maxWidth: '420px',
            width: 'auto',
            fontSize: '14px',
            lineHeight: '1.5',
            padding: '12px 16px',
            boxShadow: '0 4px 20px rgba(0,0,0,0.4)',
          }
        }} 
      />
      <Layout>
        <Routes>
          {/* Публічні маршрути */}
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/signup" element={<SignupTenant />} />
          <Route path="/bootstrap" element={<Bootstrap />} />
          <Route path="/setup-password" element={<SetupPassword />} />

          <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />

          <Route
            path="/audit"
            element={
              <ProtectedRoute allowedRoles={[...ROLE_GROUPS.superAdmin]}>
                <AuditLogs />
              </ProtectedRoute>
            }
          />

          <Route 
            path="/my-equipment" 
            element={<ProtectedRoute forbidRoles={['CONTRACTOR']}><Profile /></ProtectedRoute>} 
          />

          <Route 
            path="/inventory" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.inventory]}><Inventory /></ProtectedRoute>} 
          />
          <Route 
            path="/warehouses" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.inventory]}><Warehouses /></ProtectedRoute>} 
          />
          <Route 
            path="/requests" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.requests]}><Requests /></ProtectedRoute>} 
          />
          <Route 
            path="/my-shipments" 
            element={<ProtectedRoute allowedRoles={['EMPLOYEE']}><MyShipments /></ProtectedRoute>} 
          />
          <Route 
            path="/units" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.units]}><Units /></ProtectedRoute>} 
          />
          <Route 
            path="/admin/users" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.users]}><AdminUsers /></ProtectedRoute>} 
          />
          <Route 
            path="/vehicles" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.transport]}><Vehicles /></ProtectedRoute>} 
          />
          
          <Route 
            path="/analytics" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.analytics]}><AnalyticsDashboard /></ProtectedRoute>} 
          />

          <Route 
            path="/kpi" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.analytics]}><KPIDashboard /></ProtectedRoute>} 
          />

          <Route 
            path="/gps" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.transport]}><GPSTracking /></ProtectedRoute>} 
          />

          <Route 
            path="/maintenance" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.transport]}><MaintenanceSchedule /></ProtectedRoute>} 
          />

          <Route 
            path="/fuel-anomalies" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.transport]}><FuelAnomalies /></ProtectedRoute>} 
          />

          <Route 
            path="/kiosk" 
            element={
              <ProtectedRoute allowedRoles={[...ROLE_GROUPS.kiosk]}>
                <KioskTerminal />
              </ProtectedRoute>
            } 
          />

          <Route 
            path="/contractor-requests" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.contractorRequestsView]}><VolunteerRequests /></ProtectedRoute>} 
          />

          <Route 
            path="/billing" 
            element={<ProtectedRoute allowedRoles={[...ROLE_GROUPS.analytics]}><Billing /></ProtectedRoute>} 
          />

          <Route
            path="/platform"
            element={
              <ProtectedRoute allowedRoles={[...ROLE_GROUPS.platform]}>
                <PlatformAdmin />
              </ProtectedRoute>
            }
          />

          {/* Редирект для всього іншого */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Layout>
      </GPSProvider>
    </AuthProvider>
  )
}

export default App