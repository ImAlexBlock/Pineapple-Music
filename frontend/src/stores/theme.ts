import { defineStore } from 'pinia'
import { ref } from 'vue'
import vuetify from '../plugins/vuetify'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(false)

  function toggle() {
    isDark.value = !isDark.value
    vuetify.theme.global.name.value = isDark.value ? 'pineappleDark' : 'pineappleLight'
    localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
  }

  function init() {
    const saved = localStorage.getItem('theme')
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    if (saved === 'dark' || (!saved && prefersDark)) {
      isDark.value = true
      vuetify.theme.global.name.value = 'pineappleDark'
    }
  }

  return { isDark, toggle, init }
})
