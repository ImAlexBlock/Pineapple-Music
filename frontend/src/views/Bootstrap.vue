<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="pa-6" elevation="8">
          <v-card-title class="text-center text-h4 mb-4">🍍</v-card-title>
          <v-card-subtitle class="text-center text-h6">{{ t('auth.setupTitle') }}</v-card-subtitle>
          <v-card-text class="text-center">{{ t('auth.setupDesc') }}</v-card-text>

          <v-card-text v-if="done">
            <v-alert type="success" prominent>
              <div class="text-body-1 font-weight-bold mb-2">{{ t('auth.bootstrapDone') }}</div>
              <div class="text-body-2">{{ t('auth.checkConsole') }}</div>
            </v-alert>
            <v-btn block color="primary" class="mt-4" @click="goToLogin">{{ t('auth.login') }}</v-btn>
          </v-card-text>

          <v-card-actions v-else>
            <v-btn block color="primary" size="large" :loading="loading" @click="doBootstrap">
              {{ t('auth.bootstrap') }}
            </v-btn>
          </v-card-actions>

          <v-card-text v-if="error">
            <v-alert type="error">{{ error }}</v-alert>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const { t } = useI18n()
const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const done = ref(false)
const error = ref('')

async function doBootstrap() {
  loading.value = true
  error.value = ''
  try {
    await auth.bootstrap()
    done.value = true
  } catch (e: any) {
    error.value = e?.response?.data?.message || 'Bootstrap failed'
  }
  loading.value = false
}

function goToLogin() {
  router.push('/login')
}
</script>
