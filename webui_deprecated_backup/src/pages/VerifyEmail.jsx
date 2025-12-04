import { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { apiRequest } from '../lib/api'

function VerifyEmail() {
  const location = useLocation()
  const navigate = useNavigate()
  const [status, setStatus] = useState('pending') // pending | success | error
  const [message, setMessage] = useState('')

  useEffect(() => {
    const params = new URLSearchParams(location.search)
    const token = params.get('token')
    if (!token) {
      setStatus('error')
      setMessage('Kein Token übergeben.')
      return
    }

    const verify = async () => {
      try {
        await apiRequest('/auth/verify-email', {
          method: 'POST',
          body: JSON.stringify({ token }),
        })
        setStatus('success')
        setMessage('E-Mail erfolgreich verifiziert. Du kannst dich jetzt einloggen.')
      } catch (err) {
        setStatus('error')
        setMessage(err.message || 'Verifizierung fehlgeschlagen')
      }
    }

    verify()
  }, [location.search])

  return (
    <div style={{ maxWidth: '420px', margin: '120px auto', padding: '24px', border: '1px solid #e5e7eb', borderRadius: '8px' }}>
      <h1 style={{ marginTop: 0 }}>E-Mail Verifizieren</h1>
      {status === 'pending' && <p>Verifiziere Token...</p>}
      {status === 'success' && <p style={{ color: '#16a34a' }}>{message}</p>}
      {status === 'error' && <p style={{ color: '#dc2626' }}>{message}</p>}
      <button
        style={{ marginTop: '16px', padding: '10px 16px', background: '#1f2937', color: 'white', border: 'none', cursor: 'pointer' }}
        onClick={() => navigate('/login')}
      >
        Zur Anmeldung
      </button>
    </div>
  )
}

export default VerifyEmail
