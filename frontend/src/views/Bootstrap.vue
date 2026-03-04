<template>
  <div class="flex min-h-[calc(100vh-3rem)] items-center justify-center p-4">
    <Card class="w-full max-w-sm">
      <CardHeader class="text-center">
        <div class="text-4xl mb-2">🍍</div>
        <CardTitle class="text-xl">{{ t('auth.setupTitle') }}</CardTitle>
        <CardDescription>{{ t('auth.setupDesc') }}</CardDescription>
      </CardHeader>
      <CardContent>
        <template v-if="done">
          <Alert variant="default" class="mb-4 border-success bg-success/10">
            <CheckCircle class="h-4 w-4 text-success" />
            <AlertTitle>{{ t('auth.bootstrapDone') }}</AlertTitle>
            <AlertDescription>{{ t('auth.checkConsole') }}</AlertDescription>
          </Alert>
          <Button class="w-full" @click="goToLogin">{{ t('auth.login') }}</Button>
        </template>
        <Button v-else class="w-full" size="lg" :disabled="loading" @click="doBootstrap">
          <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
          {{ t('auth.bootstrap') }}
        </Button>
        <Alert v-if="error" variant="destructive" class="mt-4">
          <AlertCircle class="h-4 w-4" />
          <AlertDescription>{{ error }}</AlertDescription>
        </Alert>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { CheckCircle, AlertCircle, Loader2 } from 'lucide-vue-next'

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
