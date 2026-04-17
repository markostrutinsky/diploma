import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import Layout from './components/Layout'
import ProtectedRoute from './components/ProtectedRoute'
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
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
import { Toaster } from 'react-hot-toast'

function App() {
  return (
    <AuthProvider>
      <Toaster 
        position="top-right" 
        toastOptions={{
          duration: 4000,
          style: { background: '#1e293b', color: '#fff', border: '1px solid #334155' }
        }} 
      />
      <Layout>
        <Routes>
          {/* Публічні маршрути */}
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/bootstrap" element={<Bootstrap />} />
          <Route path="/setup-password" element={<SetupPassword />} />

          <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />

          <Route path="/audit" element={<ProtectedRoute><AuditLogs /></ProtectedRoute>}/>

          <Route 
            path="/my-equipment" 
            element={<ProtectedRoute forbidRoles={['CONTRACTOR']}><Profile /></ProtectedRoute>} 
          />

          <Route 
            path="/inventory" 
            element={<ProtectedRoute forbidRoles={['CONTRACTOR']}><Inventory /></ProtectedRoute>} 
          />
          <Route 
            path="/warehouses" 
            element={<ProtectedRoute forbidRoles={['CONTRACTOR']}><Warehouses /></ProtectedRoute>} 
          />
          <Route 
            path="/requests" 
            element={<ProtectedRoute forbidRoles={['CONTRACTOR']}><Requests /></ProtectedRoute>} 
          />
          <Route 
            path="/units" 
            element={<ProtectedRoute forbidRoles={['CONTRACTOR']}><Units /></ProtectedRoute>} 
          />
          <Route 
            path="/admin/users" 
            element={<ProtectedRoute forbidRoles={['CONTRACTOR']}><AdminUsers /></ProtectedRoute>} 
          />
          <Route 
            path="/vehicles" 
            element={<ProtectedRoute forbidRoles={['CONTRACTOR']}><Vehicles /></ProtectedRoute>} 
          />
          
          <Route 
            path="/analytics" 
            element={<ProtectedRoute forbidRoles={['CONTRACTOR']}><AnalyticsDashboard /></ProtectedRoute>} 
          />

          <Route 
            path="/CONTRACTOR-requests" 
            element={<ProtectedRoute><VolunteerRequests /></ProtectedRoute>} 
          />

          {/* Редирект для всього іншого */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Layout>
    </AuthProvider>
  )
}

export default App