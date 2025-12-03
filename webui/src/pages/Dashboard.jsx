import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { apiRequest } from '../lib/api'

const HEALTH_POLL_MS = 5000
const MONITORING_POLL_MS = 5000

function StatusPill({ label, status }) {
  const colors = {
    ok: '#16a34a',
    degraded: '#dc2626',
    down: '#dc2626',
  }
  const color = colors[status] || '#4b5563'
  return (
    <span style={{
      display: 'inline-block',
      padding: '4px 10px',
      borderRadius: '12px',
      background: color,
      color: 'white',
      fontSize: '12px',
      textTransform: 'uppercase',
      letterSpacing: '0.02em',
    }}>
      {label}
    </span>
  )
}

function Dashboard() {
  const [health, setHealth] = useState(null)
  const [healthError, setHealthError] = useState('')
  const [samples, setSamples] = useState([])
  const [monitoringError, setMonitoringError] = useState('')
  const navigate = useNavigate()
  const { isAuthenticated, logout } = useAuth()

  useEffect(() => {
    // SECURITY: Check authentication via in-memory context (not localStorage)
    if (!isAuthenticated) {
      navigate('/login')
      return
    }

    let mounted = true
    const fetchHealth = async () => {
      try {
        const data = await apiRequest('/health', { method: 'GET' })
        if (mounted) {
          setHealth(data)
          setHealthError('')
        }
      } catch (err) {
        if (mounted) {
          setHealth(null)
          setHealthError(err.message || 'Health-Check fehlgeschlagen')
        }
      }
    }

    fetchHealth()
    const interval = setInterval(fetchHealth, HEALTH_POLL_MS)
    return () => {
      mounted = false
      clearInterval(interval)
    }
  }, [])

  useEffect(() => {
    let mounted = true
    const fetchMonitoring = async () => {
      try {
        const data = await apiRequest('/api/monitoring')
        if (mounted) {
          setSamples(data.items || [])
          setMonitoringError('')
        }
      } catch (err) {
        if (!mounted) return
        if (err.status === 401) {
          navigate('/login')
          return
        }
        setMonitoringError(err.message)
      }
    }

    fetchMonitoring()
    const interval = setInterval(fetchMonitoring, MONITORING_POLL_MS)
    return () => {
      mounted = false
      clearInterval(interval)
    }
  }, [navigate])

  const dependencyStatus = (deps = {}) => {
    const entries = Object.entries(deps)
    if (!entries.length) return null
    return (
      <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap', marginTop: '8px' }}>
        {entries.map(([name, status]) => (
          <StatusPill key={name} label={`${name}: ${status}`} status={status === 'ok' ? 'ok' : 'down'} />
        ))}
      </div>
    )
  }

  return (
    <div style={{ maxWidth: '900px', margin: '40px auto', padding: '0 16px' }}>
      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
          <h1 style={{ margin: 0 }}>NAS.AI Monitoring</h1>
          <span style={{ padding: '4px 8px', borderRadius: '12px', background: '#1f2937', color: '#fff', fontSize: '12px' }}>
            Version 1.0
          </span>
        </div>
        <button
          onClick={() => {
            // SECURITY: Logout via AuthContext (clears in-memory token + legacy localStorage)
            logout()
            navigate('/login')
          }}
          style={{ padding: '8px 16px', background: '#1f2937', color: '#fff', border: 'none', cursor: 'pointer' }}
        >
          Logout
        </button>
      </header>

      <section style={{ marginBottom: '20px', padding: '16px', border: '1px solid #e5e7eb', borderRadius: '8px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h2 style={{ margin: '0 0 8px 0' }}>Health</h2>
            <p style={{ margin: 0, color: '#4b5563' }}>Live-Status der API und Dependencies.</p>
          </div>
          {health && (
            <StatusPill
              label={health.status === 'ok' ? 'OK' : 'DEGRADED'}
              status={health.status === 'ok' ? 'ok' : 'degraded'}
            />
          )}
        </div>
        <div style={{ marginTop: '12px', color: '#111827' }}>
          {healthError && <div style={{ color: '#dc2626' }}>{healthError}</div>}
          {!healthError && health && (
            <div>
              <div>Service: {health.service}</div>
              <div>Version: {health.version}</div>
              {dependencyStatus(health.dependencies)}
            </div>
          )}
        </div>
      </section>

      <section style={{ padding: '16px', border: '1px solid #e5e7eb', borderRadius: '8px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h2 style={{ margin: '0 0 8px 0' }}>Monitoring Daten</h2>
            <p style={{ margin: 0, color: '#4b5563' }}>Rohe CPU/RAM Samples vom MonitoringAgent.</p>
          </div>
        </div>
        {monitoringError && <div style={{ color: '#dc2626', marginTop: '10px' }}>{monitoringError}</div>}
        {!monitoringError && samples.length === 0 && <div style={{ marginTop: '10px' }}>Noch keine Daten.</div>}
        {!monitoringError && samples.length > 0 && (
          <ul style={{ listStyle: 'none', padding: 0, marginTop: '10px' }}>
            {samples.map((s) => (
              <li key={s.id} style={{ padding: '8px 0', borderBottom: '1px solid #e5e7eb' }}>
                <div style={{ fontWeight: 600 }}>{s.source}</div>
                <div style={{ color: '#4b5563' }}>
                  CPU: {s.cpu_percent}% | RAM: {s.ram_percent}%
                </div>
                <div style={{ fontSize: '12px', color: '#6b7280' }}>{new Date(s.created_at || s.createdAt).toLocaleString()}</div>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  )
}

export default Dashboard
