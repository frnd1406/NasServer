import express from 'express'

const app = express()
const PORT = 3001

// Health endpoint for orchestrator
app.get('/health', (req, res) => {
  res.json({
    service: 'webui',
    status: 'ok',
    timestamp: new Date().toISOString(),
    version: '1.0.0'
  })
})

app.listen(PORT, () => {
  console.log(`WebUI health server on http://localhost:${PORT}`)
})
