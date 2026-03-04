<template>
  <div>
    <div class="mb-6">
      <h1 class="text-2xl font-bold">{{ t('admin.dashboard') }}</h1>
      <p class="text-sm text-muted-foreground mt-1">System overview</p>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
      <Card class="p-4 text-center">
        <Music class="h-7 w-7 mx-auto mb-2 text-primary" />
        <div class="text-3xl font-bold">{{ stats.tracks }}</div>
        <div class="text-xs text-muted-foreground">{{ t('admin.totalTracks') }}</div>
      </Card>
      <Card class="p-4 text-center">
        <HardDrive class="h-7 w-7 mx-auto mb-2 text-secondary" />
        <div class="text-3xl font-bold">{{ formatSize(stats.total_size) }}</div>
        <div class="text-xs text-muted-foreground">{{ t('admin.totalSize') }}</div>
      </Card>
      <Card class="p-4 text-center">
        <Play class="h-7 w-7 mx-auto mb-2 text-info" />
        <div class="text-3xl font-bold">{{ stats.plays }}</div>
        <div class="text-xs text-muted-foreground">{{ t('admin.totalPlays') }}</div>
      </Card>
      <Card class="p-4 text-center">
        <ListMusic class="h-7 w-7 mx-auto mb-2 text-warning" />
        <div class="text-3xl font-bold">{{ stats.playlists }}</div>
        <div class="text-xs text-muted-foreground">{{ t('admin.totalPlaylists') }}</div>
      </Card>
    </div>

    <h2 class="text-lg font-bold mb-3">Quick Actions</h2>
    <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
      <Button variant="outline" size="lg" class="w-full" as-child>
        <router-link to="/admin/upload" class="flex items-center gap-2">
          <Upload class="h-4 w-4" />
          {{ t('admin.upload') }}
        </router-link>
      </Button>
      <Button variant="outline" size="lg" class="w-full" as-child>
        <router-link to="/admin/scan" class="flex items-center gap-2">
          <ScanSearch class="h-4 w-4" />
          {{ t('admin.scan') }}
        </router-link>
      </Button>
      <Button variant="outline" size="lg" class="w-full" as-child>
        <router-link to="/admin/settings" class="flex items-center gap-2">
          <Settings class="h-4 w-4" />
          {{ t('admin.settings') }}
        </router-link>
      </Button>
      <Button variant="outline" size="lg" class="w-full" as-child>
        <router-link to="/admin/audit" class="flex items-center gap-2">
          <ClipboardList class="h-4 w-4" />
          {{ t('admin.auditLog') }}
        </router-link>
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '../api'
import type { DashboardStats } from '../types'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Music, HardDrive, Play, ListMusic, Upload, ScanSearch, Settings, ClipboardList } from 'lucide-vue-next'

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
