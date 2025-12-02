import { useEffect, useState } from 'react'
import { apiRequest } from '../lib/api'

const POLL_MS = 5000

function MetricCard({ label, value, unit, color }) {
  return (
    <div style={{ flex: 1, padding: '16px', borderRadius: '12px', background: '#111827', color: '#fff', minWidth: '200px' }}>
      <div style={{ opacity: 0.8, fontSize: '14px' }}>{label}</div>
      <div style={{ fontSize: '36px', fontWeight: 700, marginTop: '8px', color }}>
        {value ?? '—'}
        {value !== null ? unit : ''}
      </div>
    </div>
  )
}

function AlertBanner({ alerts, error }) {
  const hasCritical = alerts.some((alert) => alert.severity === 'CRITICAL')
  const hasWarning = alerts.some((alert) => alert.severity === 'WARNING')

  let bg = '#064e3b'
  let text = '#d1fae5'
  let title = 'System Healthy'
  let description = 'Keine offenen Alerts.'

  if (error) {
    bg = '#7f1d1d'
    text = '#fecdd3'
    title = 'Alerts-Feed nicht verfügbar'
    description = error
  } else if (alerts.length > 0 && hasCritical) {
    bg = '#7f1d1d'
    text = '#fecdd3'
    title = 'Critical Alerts aktiv'
    description = 'Bitte sofort prüfen.'
  } else if (alerts.length > 0 && hasWarning) {
    bg = '#78350f'
    text = '#fcd34d'
    title = 'Warnings aktiv'
    description = 'Bitte beobachten.'
  }

  return (
    <div style={{ marginTop: '16px', padding: '14px 16px', borderRadius: '10px', background: bg, color: text }}>
      <div style={{ fontWeight: 700 }}>{title}</div>
      <div style={{ marginTop: '4px', opacity: 0.9 }}>{description}</div>
      {alerts.length > 0 && (
        <ul style={{ margin: '10px 0 0', paddingLeft: '18px', lineHeight: 1.5 }}>
          {alerts.map((alert) => (
            <li key={alert.id}>
              <span style={{ fontWeight: 700 }}>{alert.severity}</span>: {alert.message}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

export default function Metrics() {
  const [metric, setMetric] = useState(null)
  const [alerts, setAlerts] = useState([])
  const [error, setError] = useState('')
  const [alertsError, setAlertsError] = useState('')
  const [lastUpdated, setLastUpdated] = useState(null)

  useEffect(() => {
    let mounted = true
    let retryDelay = POLL_MS
    let timeoutId = null

    // FIX [BUG-JS-014]: Exponential backoff to prevent DoS on backend during errors
    const fetchData = async () => {
      try {
        const [metricsRes, alertsRes] = await Promise.allSettled([
          apiRequest('/api/v1/system/metrics?limit=1'),
          apiRequest('/api/v1/system/alerts'),
        ])

        if (!mounted) return

        if (metricsRes.status === 'fulfilled') {
          const item = metricsRes.value.items?.[0]
          setMetric(item || null)
          setError('')
          setLastUpdated(new Date())
          retryDelay = POLL_MS // Reset delay on success
        } else {
          setMetric(null)
          setError('API nicht erreichbar')
          setLastUpdated(null)
          // Exponential backoff on error (max 60s)
          retryDelay = Math.min(retryDelay * 2, 60000)
        }

        if (alertsRes.status === 'fulfilled') {
          setAlerts(alertsRes.value.items || [])
          setAlertsError('')
        } else {
          setAlerts([])
          setAlertsError('Alerts konnten nicht geladen werden')
          // Exponential backoff on error (max 60s)
          retryDelay = Math.min(retryDelay * 2, 60000)
        }
      } catch (err) {
        if (!mounted) return
        setError('API nicht erreichbar')
        setAlertsError('Alerts konnten nicht geladen werden')
        setMetric(null)
        setAlerts([])
        // Exponential backoff on error (max 60s)
        retryDelay = Math.min(retryDelay * 2, 60000)
      }

      // Schedule next fetch with current retry delay
      if (mounted) {
        timeoutId = setTimeout(fetchData, retryDelay)
      }
    }

    // Initial fetch
    fetchData()

    return () => {
      mounted = false
      if (timeoutId) clearTimeout(timeoutId)
    }
  }, [])

  const cpu = metric ? Number(metric.cpu_usage)?.toFixed(2) : null
  const ram = metric ? Number(metric.ram_usage)?.toFixed(2) : null
  const disk = metric ? Number(metric.disk_usage)?.toFixed(2) : null

  return (
    <div style={{ minHeight: '100vh', background: '#0b1224', color: '#e5e7eb', padding: '40px 16px' }}>
      <div style={{ maxWidth: '900px', margin: '0 auto' }}>
        <h1 style={{ margin: 0, fontSize: '28px', color: '#f3f4f6' }}>System Metrics</h1>
        <p style={{ marginTop: '8px', color: '#9ca3af' }}>Aktualisiert alle 5 Sekunden von /api/v1/system/metrics</p>

        <AlertBanner alerts={alerts} error={alertsError} />

        {error && (
          <div style={{ marginTop: '16px', padding: '12px 14px', background: '#7f1d1d', color: '#fecdd3', borderRadius: '8px' }}>
            {error}
          </div>
        )}

        <div style={{ display: 'flex', gap: '12px', marginTop: '20px', flexWrap: 'wrap' }}>
          <MetricCard label="CPU" value={cpu} unit="%" color="#22d3ee" />
          <MetricCard label="RAM" value={ram} unit="%" color="#a78bfa" />
          <MetricCard label="Disk" value={disk} unit="%" color="#fbbf24" />
        </div>

        <div style={{ marginTop: '16px', color: '#9ca3af' }}>
          {lastUpdated ? `Letztes Update: ${lastUpdated.toLocaleTimeString()}` : 'Lade...'}
        </div>
      </div>
    </div>
  )
}
