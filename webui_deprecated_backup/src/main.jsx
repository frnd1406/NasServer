import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.jsx'

// FIX [BUG-JS-010]: Disable console in production to prevent information leakage
if (import.meta.env.PROD) {
  console.log = () => {}
  console.error = () => {}
  console.warn = () => {}
  console.debug = () => {}
}

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
