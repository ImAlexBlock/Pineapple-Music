import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '../api'

export const useAuthStore = defineStore('auth', () => {
  const role = ref<string | null>(null)
  const bootstrapped = ref(false)
  const checked = ref(false)

  async function checkStatus() {
    try {
      const { data } = await authApi.getStatus()
      bootstrapped.value = data.bootstrapped

      if (bootstrapped.value) {
        try {
          const { data: meData } = await authApi.me()
          role.value = meData.role
        } catch {
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

  return { role, bootstrapped, checked, checkStatus, login, logout, bootstrap }
})
