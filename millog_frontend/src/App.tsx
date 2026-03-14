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

function App() {
  return (
    <AuthProvider>
      <Layout>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/bootstrap" element={<Bootstrap />} />
          <Route path="/setup-password" element={<SetupPassword />} />
          <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />
          <Route path="/admin/users" element={<ProtectedRoute><AdminUsers /></ProtectedRoute>} />
          <Route path="/inventory" element={<ProtectedRoute><Inventory /></ProtectedRoute>} />
          <Route path="/requests" element={<ProtectedRoute><Requests /></ProtectedRoute>} />
          <Route path="/volunteer-requests" element={<ProtectedRoute><VolunteerRequests /></ProtectedRoute>} />
          <Route path="/units" element={<ProtectedRoute><Units /></ProtectedRoute>} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Layout>
    </AuthProvider>
  )
}

export default App
