<template>
  <v-container>
    <h1 class="text-h4 mb-4">{{ t('admin.settings') }}</h1>

    <v-card class="mb-4">
      <v-card-title>{{ t('admin.rotateGuestKey') }}</v-card-title>
      <v-card-text>
        <p class="text-body-2 mb-2">{{ t('admin.rotateGuestKeyDesc') }}</p>
        <v-alert v-if="newGuestKey" type="warning" class="mb-2">
          {{ t('auth.saveKey') }}: <code>{{ newGuestKey }}</code>
        </v-alert>
      </v-card-text>
      <v-card-actions>
        <v-btn color="warning" @click="rotateGuest">{{ t('admin.rotateGuestKey') }}</v-btn>
      </v-card-actions>
    </v-card>

    <v-card class="mb-4">
      <v-card-title>{{ t('admin.rotateAdminKey') }}</v-card-title>
      <v-card-text v-if="newAdminKey">
        <v-alert type="warning">{{ t('auth.saveKey') }}: <code>{{ newAdminKey }}</code></v-alert>
      </v-card-text>
      <v-card-actions>
        <v-btn color="error" @click="rotateAdmin">{{ t('admin.rotateAdminKey') }}</v-btn>
      </v-card-actions>
    </v-card>

    <v-card class="mb-4">
      <v-card-title>{{ t('admin.accessMode') }}</v-card-title>
      <v-card-text>
        <v-select
          v-model="accessMode"
          :items="['public', 'private']"
          :label="t('admin.accessMode')"
          @update:model-value="saveSetting('access_mode', accessMode)"
        />
      </v-card-text>
    </v-card>

    <v-card>
      <v-card-title>Subsonic Protocol</v-card-title>
      <v-card-text>
        <v-switch v-model="subsonicEnabled" label="Enable Subsonic API" @update:model-value="saveSetting('subsonic_enabled', subsonicEnabled ? 'true' : 'false')" />
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '../api'

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
