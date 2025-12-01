import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import { setAuthContext } from './lib/api'
import Login from './pages/Login'
import Register from './pages/Register'
import Success from './pages/Success'
import Dashboard from './pages/Dashboard'
import VerifyEmail from './pages/VerifyEmail'
import Metrics from './pages/Metrics'
import { useEffect } from 'react'

function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Metrics />} />
      <Route path="/metrics" element={<Metrics />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route path="/success" element={<Success />} />
      <Route path="/dashboard" element={<Dashboard />} />
      <Route path="/verify-email" element={<VerifyEmail />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

// Bridge component to connect AuthContext to api.js
function AuthContextBridge({ children }) {
  const auth = useAuth()

  useEffect(() => {
    // Connect AuthContext to api.js so it can access in-memory tokens
    setAuthContext(auth)
  }, [auth])

  return children
}

function App() {
  return (
    <AuthProvider>
      <AuthContextBridge>
        <Router>
          <AppRoutes />
        </Router>
      </AuthContextBridge>
    </AuthProvider>
  )
}

export default App
