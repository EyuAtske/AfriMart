export default defineEventHandler(async (event) => {
  const path = event.path.replace(/^\/api/, '')
  const target = `http://localhost:8080/api${path}`

  return await proxyRequest(event, target)
})
