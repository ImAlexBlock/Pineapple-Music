import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  withCredentials: true,
})

// Attach CSRF token from cookie
api.interceptors.request.use((config) => {
  const csrfToken = getCookie('pm_csrf')
  if (csrfToken) {
    config.headers['X-CSRF-Token'] = csrfToken
  }
  return config
})

// Handle 401 — only redirect for explicitly protected calls, not for auth probing
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Never auto-redirect on /auth/me, /setup/*, or /play-events failures
    // These are expected to fail for unauthenticated visitors in public mode
    return Promise.reject(error)
  }
)

function getCookie(name: string): string | null {
  const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'))
  return match ? match[2] : null
}

export default api
