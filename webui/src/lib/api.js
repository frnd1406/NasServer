const envBaseUrl = (import.meta.env.VITE_API_BASE_URL || '').trim()

// SECURITY: AuthContext reference for in-memory token access
// Set by setAuthContext() when app initializes
let authContextRef = null

export function setAuthContext(authContext) {
  authContextRef = authContext
}

function deriveDefaultBaseUrl() {
  // When no env is provided, fall back to the current host but the API port (8080 for dev).
  if (typeof window === 'undefined') {
    return 'http://localhost:8080'
  }
  const { protocol, hostname } = window.location
  const defaultPort = protocol === 'https:' ? '8443' : '8080'
  return `${protocol}//${hostname}:${defaultPort}`
}

function normalizeBaseUrl(url) {
  const base = url || deriveDefaultBaseUrl()
  return base.replace(/\/+$/, '')
}

const API_BASE_URL = normalizeBaseUrl(envBaseUrl)
const LOGOUT_COUNTDOWN_SECONDS = 4

// FIX [BUG-JS-019]: Validate and warn about API URL configuration
if (!envBaseUrl && import.meta.env.PROD) {
  console.warn('⚠️  VITE_API_BASE_URL not configured, using derived URL:', API_BASE_URL)
}

function buildUrl(path = '') {
  if (!path.startsWith('/')) {
    return `${API_BASE_URL}/${path}`
  }
  return `${API_BASE_URL}${path}`
}

export function getApiBaseUrl() {
  return API_BASE_URL
}

// FIX [BUG-JS-003]: Removed global state and DOM manipulation
// Session warning is now handled by AuthContext

async function refreshAccessToken() {
  const refreshToken = localStorage.getItem('refreshToken')
  if (!refreshToken) return null

  try {
    const res = await fetch(buildUrl('/auth/refresh'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
      // CRITICAL FIX: Send cookies with refresh request (BUG-JS-002)
      credentials: 'include',
    })

    if (!res.ok) {
      return null
    }

    const data = await res.json().catch(() => null)
    const newAccessToken = data?.access_token || data?.token

    if (!newAccessToken) {
      return null
    }

    // SECURITY: Store accessToken in-memory via AuthContext (not localStorage)
    if (authContextRef?.login) {
      authContextRef.login(newAccessToken)
    }

    // Keep refresh_token in localStorage temporarily (will move to HttpOnly cookie)
    if (data?.refresh_token) {
      localStorage.setItem('refreshToken', data.refresh_token)
    }
    if (data?.csrf_token) {
      localStorage.setItem('csrfToken', data.csrf_token)
    }

    return newAccessToken
  } catch (err) {
    return null
  }
}

function redirectToLogin() {
  if (typeof window === 'undefined') {
    return
  }

  // Use AuthContext to trigger session warning UI
  if (authContextRef?.showSessionWarning) {
    authContextRef.showSessionWarning()
  } else {
    // Fallback if context not available
    window.location.href = '/login'
  }
}

function buildHeaders(accessToken, headersOverride = {}) {
  const csrfToken = localStorage.getItem('csrfToken') || localStorage.getItem('csrf_token') || ''
  const headers = {
    'Content-Type': 'application/json',
    ...headersOverride,
  }

  if (accessToken) {
    headers.Authorization = `Bearer ${accessToken}`
  }
  if (csrfToken && !headers['X-CSRF-Token']) {
    headers['X-CSRF-Token'] = csrfToken
  }

  return headers
}

function extractErrorMessage(res, data) {
  return data?.error?.message || data?.error || res.statusText || 'Request failed'
}

async function performRequest(path, options, tokenOverride) {
  // SECURITY: Get accessToken from in-memory AuthContext (not localStorage)
  const accessToken = tokenOverride || authContextRef?.accessToken || null
  // FIX [BUG-JS-012]: Add timeout to prevent hanging requests
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), 10000) // 10s timeout

  let res
  try {
    res = await fetch(buildUrl(path), {
      ...options,
      headers: buildHeaders(accessToken, options.headers),
      // CRITICAL FIX: Send cookies with every request (BUG-JS-002)
      credentials: 'include',
      signal: controller.signal,
    })
  } catch (err) {
    if (err.name === 'AbortError') {
      throw new Error(`Request timed out after 10s`)
    }
    throw new Error(`Cannot reach API at ${API_BASE_URL} (${err.message})`)
  } finally {
    clearTimeout(timeoutId)
  }

  const isJson = res.headers.get('content-type')?.includes('application/json')
  let data = null

  if (isJson) {
    try {
      data = await res.json()
    } catch (err) {
      data = null
    }
  }

  return { res, data }
}

export async function apiRequest(path, options = {}) {
  const firstAttempt = await performRequest(path, options)

  if (firstAttempt.res.ok) {
    return firstAttempt.data
  }

  if (firstAttempt.res.status === 401) {
    const newAccessToken = await refreshAccessToken()

    if (newAccessToken) {
      const retry = await performRequest(path, options, newAccessToken)

      if (retry.res.ok) {
        return retry.data
      }

      if (retry.res.status === 401) {
        redirectToLogin()
      }

      const retryMessage = extractErrorMessage(retry.res, retry.data)
      const retryError = new Error(retryMessage)
      retryError.status = retry.res.status
      throw retryError
    }

    redirectToLogin()
    const refreshError = new Error('Session expired. Please log in again.')
    refreshError.status = 401
    throw refreshError
  }

  const message = extractErrorMessage(firstAttempt.res, firstAttempt.data)
  const error = new Error(message)
  error.status = firstAttempt.res.status
  throw error
}
