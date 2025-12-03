import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import { setAuthContext } from './lib/api'
import Login from './pages/Login'
import Register from './pages/Register'
import Success from './pages/Success'
import Dashboard from './pages/Dashboard'
import VerifyEmail from './pages/VerifyEmail'
import Metrics from './pages/Metrics'
import NotFound from './pages/NotFound'
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
      {/* FIX [BUG-JS-020]: Add proper 404 page instead of redirect */}
      <Route path="*" element={<NotFound />} />
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

// FIX [BUG-JS-008]: Simple ErrorBoundary component
class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error) {
    return { hasError: true }
  }

  componentDidCatch(error, errorInfo) {
    console.error("Uncaught error:", error, errorInfo)
  }

  render() {
    if (this.state.hasError) {
      return (
        <div style={{ padding: '2rem', textAlign: 'center', color: '#dc2626' }}>
          <h1>Something went wrong.</h1>
          <p>Please refresh the page.</p>
        </div>
      )
    }

    return this.props.children
  }
}

function App() {
  return (
    <AuthProvider>
      <AuthContextBridge>
        <ErrorBoundary>
          <Router>
            <AppRoutes />
          </Router>
        </ErrorBoundary>
      </AuthContextBridge>
    </AuthProvider>
  )
}

export default App
