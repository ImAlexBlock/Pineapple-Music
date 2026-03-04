<template>
  <div class="flex min-h-[calc(100vh-3rem)] items-center justify-center p-4">
    <Card class="w-full max-w-sm">
      <CardHeader class="text-center">
        <div class="mx-auto mb-3 flex h-16 w-16 items-center justify-center rounded-full bg-primary">
          <KeyRound class="h-8 w-8 text-primary-foreground" />
        </div>
        <CardTitle class="text-xl">{{ t('auth.login') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="space-y-2">
          <Label :for="'key'">{{ t('auth.adminKey') }} / {{ t('auth.guestKey') }}</Label>
          <Input
            id="key"
            v-model="key"
            type="password"
            @keyup.enter="doLogin"
          />
          <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
        </div>
        <Button class="w-full" size="lg" :disabled="loading" @click="doLogin">
          <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
          {{ t('auth.login') }}
        </Button>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { KeyRound, Loader2 } from 'lucide-vue-next'

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
