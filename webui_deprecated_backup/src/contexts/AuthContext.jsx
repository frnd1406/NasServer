import React, { createContext, useContext, useState, useCallback, useEffect } from 'react'

const AuthContext = createContext(null)

const LOGOUT_COUNTDOWN_SECONDS = 4

export function AuthProvider({ children }) {
  // SECURITY: Store accessToken only in memory (not localStorage)
  // This prevents XSS attacks from stealing tokens via localStorage
  const [accessToken, setAccessToken] = useState(null)
  const [showOverlay, setShowOverlay] = useState(false)
  const [countdown, setCountdown] = useState(LOGOUT_COUNTDOWN_SECONDS)

  const login = useCallback((token) => {
    setAccessToken(token)
  }, [])

  const logout = useCallback(() => {
    setAccessToken(null)
    setShowOverlay(false)
    // Clean up any legacy tokens from localStorage
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('csrfToken')
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('csrf_token')
  }, [])

  const showSessionWarning = useCallback(() => {
    setShowOverlay(true)
    setCountdown(LOGOUT_COUNTDOWN_SECONDS)
  }, [])

  useEffect(() => {
    let interval
    if (showOverlay && countdown > 0) {
      interval = setInterval(() => {
        setCountdown((c) => c - 1)
      }, 1000)
    } else if (showOverlay && countdown <= 0) {
      // Redirect happens here or in component consuming this state
      window.location.href = '/login'
    }
    return () => clearInterval(interval)
  }, [showOverlay, countdown])

  const value = {
    accessToken,
    login,
    logout,
    showSessionWarning,
    isAuthenticated: accessToken !== null,
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
      {showOverlay && (
        <div style={{
          position: 'fixed',
          inset: 0,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          background: 'rgba(0,0,0,0.35)',
          backdropFilter: 'blur(12px)',
          zIndex: 9999,
          color: 'white',
          fontFamily: 'system-ui'
        }}>
          <div style={{
            background: 'rgba(239,68,68,0.9)',
            padding: '2rem',
            borderRadius: '1rem',
            textAlign: 'center'
          }}>
            <h2>Session Expired</h2>
            <p>Redirecting to login in {countdown}s...</p>
          </div>
        </div>
      )}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
