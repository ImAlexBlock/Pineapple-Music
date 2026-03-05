import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '../api'

export const useAuthStore = defineStore('auth', () => {
  const role = ref<string | null>(null)
  const bootstrapped = ref(false)
  const checked = ref(false)
  const accessMode = ref<string>('public')
  let checkPromise: Promise<void> | null = null

  async function checkStatus() {
    // Deduplicate: if already checking or checked, return the same promise
    if (checked.value) return
    if (checkPromise) return checkPromise

    checkPromise = _doCheck()
    return checkPromise
  }

  async function _doCheck() {
    try {
      const { data } = await authApi.getStatus()
      bootstrapped.value = data.bootstrapped
      accessMode.value = data.access_mode || 'public'

      if (bootstrapped.value) {
        try {
          const { data: meData } = await authApi.me()
          role.value = meData.role
        } catch {
          // 401 is expected for unauthenticated visitors in public mode
          role.value = null
        }
      }
    } catch {
      // Server unreachable
    }
    checked.value = true
  }

  async function login(key: string) {
    const { data } = await authApi.login(key)
    role.value = data.role
  }

  async function logout() {
    await authApi.logout()
    role.value = null
  }

  async function bootstrap() {
    await authApi.bootstrap()
    bootstrapped.value = true
  }

  return { role, bootstrapped, checked, accessMode, checkStatus, login, logout, bootstrap }
})
