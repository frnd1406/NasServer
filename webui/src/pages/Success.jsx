import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

function Success() {
  const [user, setUser] = useState(null)
  const navigate = useNavigate()
  const { isAuthenticated, logout } = useAuth()

  useEffect(() => {
    // SECURITY: Check authentication via in-memory context (not localStorage)
    if (!isAuthenticated) {
      navigate('/login')
      return
    }

    // Fetch user info (optional - you can skip this if API doesn't provide user endpoint)
    // For now, just show success message
    setUser({ email: 'user@example.com' })
  }, [navigate, isAuthenticated])

  const handleLogout = () => {
    // SECURITY: Logout via AuthContext (clears in-memory token + legacy localStorage)
    logout()
    navigate('/login')
  }

  if (!user) {
    return <div>Loading...</div>
  }

  return (
    <div style={{ maxWidth: '600px', margin: '100px auto', padding: '40px', border: '1px solid #ccc', textAlign: 'center' }}>
      <h1 style={{ color: '#28a745', marginBottom: '20px' }}>Login Successful!</h1>
      <p style={{ fontSize: '18px', marginBottom: '30px' }}>
        Welcome to NAS.AI
      </p>
      <div style={{ padding: '20px', background: '#f8f9fa', borderRadius: '4px', marginBottom: '20px' }}>
        <p>You are now logged in.</p>
        <p style={{ color: '#666', marginTop: '10px' }}>Gehe zum Dashboard für Health & Monitoring.</p>
      </div>
      <button
        onClick={() => navigate('/dashboard')}
        style={{ padding: '10px 30px', background: '#0f766e', color: 'white', border: 'none', cursor: 'pointer', fontSize: '16px', marginRight: '10px' }}
      >
        Zum Dashboard
      </button>
      <button
        onClick={handleLogout}
        style={{ padding: '10px 30px', background: '#dc3545', color: 'white', border: 'none', cursor: 'pointer', fontSize: '16px' }}
      >
        Logout
      </button>
    </div>
  )
}

export default Success
