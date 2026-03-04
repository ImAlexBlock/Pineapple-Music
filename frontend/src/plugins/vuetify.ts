import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'
import { md3 } from 'vuetify/blueprints'

const pineappleLight = {
  dark: false,
  colors: {
    primary: '#E6A117',
    'primary-darken-1': '#C88A0A',
    secondary: '#4A7C2F',
    accent: '#E86A00',
    background: '#F5F5F5',
    surface: '#FFFFFF',
    'surface-variant': '#F0EBE3',
    'surface-bright': '#FFFDF7',
    error: '#C62828',
    info: '#1565C0',
    success: '#2E7D32',
    warning: '#E65100',
    'on-primary': '#FFFFFF',
    'on-secondary': '#FFFFFF',
    'on-surface': '#1C1B1F',
    'on-background': '#1C1B1F',
  },
}

const pineappleDark = {
  dark: true,
  colors: {
    primary: '#FFD54F',
    'primary-darken-1': '#FFC107',
    secondary: '#A5D6A7',
    accent: '#FFB74D',
    background: '#0F0F0F',
    surface: '#1A1A1A',
    'surface-variant': '#262320',
    'surface-bright': '#2C2C2C',
    error: '#EF5350',
    info: '#42A5F5',
    success: '#66BB6A',
    warning: '#FFA726',
    'on-primary': '#1C1B1F',
    'on-secondary': '#1C1B1F',
    'on-surface': '#E6E1E5',
    'on-background': '#E6E1E5',
  },
}

export default createVuetify({
  blueprint: md3,
  theme: {
    defaultTheme: 'pineappleLight',
    themes: {
      pineappleLight,
      pineappleDark,
    },
  },
  defaults: {
    VCard: {
      rounded: 'lg',
      elevation: 0,
      border: true,
    },
    VBtn: {
      rounded: 'lg',
    },
    VTextField: {
      variant: 'outlined',
      density: 'comfortable',
      rounded: 'lg',
    },
    VSelect: {
      variant: 'outlined',
      density: 'comfortable',
      rounded: 'lg',
    },
    VFileInput: {
      variant: 'outlined',
      density: 'comfortable',
      rounded: 'lg',
    },
    VChip: {
      rounded: 'lg',
    },
  },
})
