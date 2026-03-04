<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="pa-6">
          <div class="text-center mb-6">
            <v-avatar color="primary" size="64" class="mb-3">
              <v-icon size="32" color="on-primary">mdi-key-variant</v-icon>
            </v-avatar>
            <div class="text-h5 font-weight-bold">{{ t('auth.login') }}</div>
          </div>
          <v-text-field
            v-model="key"
            :label="t('auth.adminKey') + ' / ' + t('auth.guestKey')"
            type="password"
            prepend-inner-icon="mdi-key"
            :error-messages="error"
            @keyup.enter="doLogin"
          />
          <v-btn block color="primary" size="large" :loading="loading" class="mt-4" @click="doLogin">
            {{ t('auth.login') }}
          </v-btn>
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
const key = ref('')
const loading = ref(false)
const error = ref('')

async function doLogin() {
  loading.value = true
  error.value = ''
  try {
    await auth.login(key.value)
    router.push('/')
  } catch {
    error.value = 'Invalid key'
  }
  loading.value = false
}
</script>
