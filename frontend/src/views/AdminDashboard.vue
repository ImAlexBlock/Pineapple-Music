<template>
  <v-container>
    <div class="d-flex align-center mb-6">
      <div>
        <h1 class="text-h4 font-weight-bold">{{ t('admin.dashboard') }}</h1>
        <p class="text-body-2 text-medium-emphasis mt-1">System overview</p>
      </div>
    </div>

    <v-row class="mb-6">
      <v-col cols="6" md="3">
        <v-card class="pa-4 text-center">
          <v-icon size="28" color="primary" class="mb-2">mdi-music-note-outline</v-icon>
          <div class="text-h4 font-weight-bold">{{ stats.tracks }}</div>
          <div class="text-caption text-medium-emphasis">{{ t('admin.totalTracks') }}</div>
        </v-card>
      </v-col>
      <v-col cols="6" md="3">
        <v-card class="pa-4 text-center">
          <v-icon size="28" color="secondary" class="mb-2">mdi-harddisk</v-icon>
          <div class="text-h4 font-weight-bold">{{ formatSize(stats.total_size) }}</div>
          <div class="text-caption text-medium-emphasis">{{ t('admin.totalSize') }}</div>
        </v-card>
      </v-col>
      <v-col cols="6" md="3">
        <v-card class="pa-4 text-center">
          <v-icon size="28" color="info" class="mb-2">mdi-play-circle-outline</v-icon>
          <div class="text-h4 font-weight-bold">{{ stats.plays }}</div>
          <div class="text-caption text-medium-emphasis">{{ t('admin.totalPlays') }}</div>
        </v-card>
      </v-col>
      <v-col cols="6" md="3">
        <v-card class="pa-4 text-center">
          <v-icon size="28" color="warning" class="mb-2">mdi-playlist-music-outline</v-icon>
          <div class="text-h4 font-weight-bold">{{ stats.playlists }}</div>
          <div class="text-caption text-medium-emphasis">{{ t('admin.totalPlaylists') }}</div>
        </v-card>
      </v-col>
    </v-row>

    <h2 class="text-h6 font-weight-bold mb-3">Quick Actions</h2>
    <v-row>
      <v-col cols="12" sm="6" md="3">
        <v-btn block variant="tonal" color="primary" to="/admin/upload" size="large">
          <v-icon start>mdi-upload-outline</v-icon>
          {{ t('admin.upload') }}
        </v-btn>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-btn block variant="tonal" color="primary" to="/admin/scan" size="large">
          <v-icon start>mdi-magnify-scan</v-icon>
          {{ t('admin.scan') }}
        </v-btn>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-btn block variant="tonal" color="primary" to="/admin/settings" size="large">
          <v-icon start>mdi-cog-outline</v-icon>
          {{ t('admin.settings') }}
        </v-btn>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-btn block variant="tonal" color="primary" to="/admin/audit" size="large">
          <v-icon start>mdi-clipboard-text-clock-outline</v-icon>
          {{ t('admin.auditLog') }}
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '../api'
import type { DashboardStats } from '../types'

const { t } = useI18n()
const stats = ref<DashboardStats>({ tracks: 0, total_size: 0, plays: 0, playlists: 0 })

onMounted(async () => {
  try {
    const { data } = await adminApi.dashboard()
    stats.value = data
  } catch (e) {
    console.error(e)
  }
})

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}
</script>
