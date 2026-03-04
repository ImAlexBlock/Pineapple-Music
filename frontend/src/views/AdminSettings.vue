<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">{{ t('admin.settings') }}</h1>

    <div class="space-y-4">
      <!-- Rotate Guest Key -->
      <Card>
        <CardHeader>
          <CardTitle class="text-base">{{ t('admin.rotateGuestKey') }}</CardTitle>
          <CardDescription>{{ t('admin.rotateGuestKeyDesc') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <Alert v-if="newGuestKey" class="mb-4 border-warning bg-warning/10">
            <AlertCircle class="h-4 w-4 text-warning" />
            <AlertDescription>{{ t('auth.saveKey') }}: <code class="font-mono bg-muted px-1 rounded">{{ newGuestKey }}</code></AlertDescription>
          </Alert>
          <Button variant="outline" @click="rotateGuest">{{ t('admin.rotateGuestKey') }}</Button>
        </CardContent>
      </Card>

      <!-- Rotate Admin Key -->
      <Card>
        <CardHeader>
          <CardTitle class="text-base">{{ t('admin.rotateAdminKey') }}</CardTitle>
        </CardHeader>
        <CardContent>
          <Alert v-if="newAdminKey" class="mb-4 border-warning bg-warning/10">
            <AlertCircle class="h-4 w-4 text-warning" />
            <AlertDescription>{{ t('auth.saveKey') }}: <code class="font-mono bg-muted px-1 rounded">{{ newAdminKey }}</code></AlertDescription>
          </Alert>
          <Button variant="destructive" @click="rotateAdmin">{{ t('admin.rotateAdminKey') }}</Button>
        </CardContent>
      </Card>

      <!-- Access Mode -->
      <Card>
        <CardHeader>
          <CardTitle class="text-base">{{ t('admin.accessMode') }}</CardTitle>
        </CardHeader>
        <CardContent>
          <Select v-model="accessMode" @update:model-value="(v: string) => saveSetting('access_mode', v)">
            <SelectTrigger class="w-48">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="public">public</SelectItem>
              <SelectItem value="private">private</SelectItem>
            </SelectContent>
          </Select>
        </CardContent>
      </Card>

      <!-- Subsonic Protocol -->
      <Card>
        <CardHeader>
          <CardTitle class="text-base">Subsonic Protocol</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="flex items-center gap-3">
            <Switch
              :checked="subsonicEnabled"
              @update:checked="(v: boolean) => { subsonicEnabled = v; saveSetting('subsonic_enabled', v ? 'true' : 'false') }"
            />
            <Label>Enable Subsonic API</Label>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '../api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { AlertCircle } from 'lucide-vue-next'

const { t } = useI18n()
const newGuestKey = ref('')
const newAdminKey = ref('')
const subsonicEnabled = ref(true)
const accessMode = ref('public')

onMounted(async () => {
  try {
    const { data: settings } = await adminApi.getSettings()
    subsonicEnabled.value = settings.subsonic_enabled !== 'false'
    accessMode.value = settings.access_mode || 'public'
  } catch (e) {
    console.error(e)
  }
})

async function rotateGuest() {
  const { data } = await adminApi.rotateGuestKey()
  newGuestKey.value = data.guest_key
}

async function rotateAdmin() {
  const { data } = await adminApi.rotateAdminKey()
  newAdminKey.value = data.admin_key
}

async function saveSetting(key: string, value: string) {
  await adminApi.updateSettings({ [key]: value })
}
</script>
