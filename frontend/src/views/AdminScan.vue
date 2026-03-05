<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold">{{ t('admin.scan') }}</h1>
        <p class="text-sm text-muted-foreground mt-1">Scan the music directory for new files</p>
      </div>
      <Button :disabled="scanning" @click="startScan" size="lg">
        <Loader2 v-if="scanning" class="mr-2 h-4 w-4 animate-spin" />
        <ScanSearch v-else class="mr-2 h-4 w-4" />
        {{ t('admin.scan') }}
      </Button>
    </div>

    <Card v-if="job">
      <CardContent class="pt-6">
        <div class="flex items-center gap-2 mb-4">
          <Loader2 v-if="job.status === 'running'" class="h-5 w-5 animate-spin text-primary" />
          <CheckCircle v-else class="h-5 w-5 text-success" />
          <span class="font-medium">
            {{ job.status === 'running' ? t('admin.scanning') : t('admin.scanComplete') }}
          </span>
        </div>

        <Progress v-if="job.status === 'running'" :model-value="job.total ? (job.scanned / job.total) * 100 : 0" class="mb-4" />

        <div class="grid grid-cols-3 gap-4 text-center">
          <div>
            <div class="text-2xl font-bold text-success">{{ job.added }}</div>
            <div class="text-xs text-muted-foreground">Added</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-info">{{ job.updated }}</div>
            <div class="text-xs text-muted-foreground">Updated</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-destructive">{{ job.errors }}</div>
            <div class="text-xs text-muted-foreground">Errors</div>
          </div>
        </div>

        <div class="text-xs text-muted-foreground mt-4" v-if="job.started_at">
          Started: {{ new Date(job.started_at).toLocaleString() }}
          <span v-if="job.finished_at"> | Finished: {{ new Date(job.finished_at).toLocaleString() }}</span>
        </div>
      </CardContent>
    </Card>

    <div v-else class="text-center py-12">
      <ScanSearch class="h-16 w-16 mx-auto mb-4 text-muted-foreground" />
      <p class="text-lg text-muted-foreground">No scans yet</p>
      <p class="text-sm text-muted-foreground">Click the button above to scan your music directory</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '../api'
import type { ScanJob } from '../types'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { ScanSearch, Loader2, CheckCircle } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const { t } = useI18n()
const job = ref<ScanJob | null>(null)
const scanning = ref(false)
let pollTimer: ReturnType<typeof setInterval> | null = null

onMounted(loadStatus)
onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})

async function loadStatus() {
  try {
    const { data } = await adminApi.scanStatus()
    job.value = data
    if (data.status === 'running') {
      startPolling()
    }
  } catch {
    // No scan yet
  }
}

async function startScan() {
  scanning.value = true
  try {
    await adminApi.startScan()
    toast.info('Scan started')
    startPolling()
  } catch (e) {
    console.error(e)
    toast.error('Failed to start scan')
  }
  scanning.value = false
}

function startPolling() {
  if (pollTimer) clearInterval(pollTimer)
  pollTimer = setInterval(async () => {
    try {
      const { data } = await adminApi.scanStatus()
      job.value = data
      if (data.status !== 'running') {
        clearInterval(pollTimer!)
        pollTimer = null
        toast.success(`Scan complete: ${data.added} added, ${data.updated} updated`)
      }
    } catch {
      clearInterval(pollTimer!)
    }
  }, 2000)
}
</script>
