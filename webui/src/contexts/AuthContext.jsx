import React, { createContext, useContext, useState, useCallback } from 'react'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  // SECURITY: Store accessToken only in memory (not localStorage)
  // This prevents XSS attacks from stealing tokens via localStorage
  const [accessToken, setAccessToken] = useState(null)

  const login = useCallback((token) => {
    setAccessToken(token)
  }, [])

  const logout = useCallback(() => {
    setAccessToken(null)
    // Clean up any legacy tokens from localStorage
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('csrfToken')
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('csrf_token')
  }, [])

  const value = {
    accessToken,
    login,
    logout,
    isAuthenticated: accessToken !== null,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
