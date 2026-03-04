<template>
  <v-container>
    <div class="d-flex align-center mb-6">
      <div>
        <h1 class="text-h4 font-weight-bold">{{ t('admin.scan') }}</h1>
        <p class="text-body-2 text-medium-emphasis mt-1">Scan the music directory for new files</p>
      </div>
      <v-spacer />
      <v-btn color="primary" variant="tonal" :loading="scanning" @click="startScan" size="large">
        <v-icon start>mdi-magnify-scan</v-icon>
        {{ t('admin.scan') }}
      </v-btn>
    </div>

    <v-card v-if="job" class="pa-2">
      <v-card-title class="d-flex align-center">
        <v-icon start :color="job.status === 'running' ? 'primary' : 'success'" class="mr-2">
          {{ job.status === 'running' ? 'mdi-loading mdi-spin' : 'mdi-check-circle' }}
        </v-icon>
        {{ job.status === 'running' ? t('admin.scanning') : t('admin.scanComplete') }}
      </v-card-title>
      <v-card-text>
        <v-progress-linear
          v-if="job.status === 'running'"
          :model-value="job.total ? (job.scanned / job.total) * 100 : 0"
          color="primary"
          height="8"
          rounded
          class="mb-4"
        />
        <v-row dense>
          <v-col cols="4">
            <div class="text-center">
              <div class="text-h5 font-weight-bold text-success">{{ job.added }}</div>
              <div class="text-caption text-medium-emphasis">Added</div>
            </div>
          </v-col>
          <v-col cols="4">
            <div class="text-center">
              <div class="text-h5 font-weight-bold text-info">{{ job.updated }}</div>
              <div class="text-caption text-medium-emphasis">Updated</div>
            </div>
          </v-col>
          <v-col cols="4">
            <div class="text-center">
              <div class="text-h5 font-weight-bold text-error">{{ job.errors }}</div>
              <div class="text-caption text-medium-emphasis">Errors</div>
            </div>
          </v-col>
        </v-row>
        <div class="text-caption text-medium-emphasis mt-3" v-if="job.started_at">
          Started: {{ new Date(job.started_at).toLocaleString() }}
          <span v-if="job.finished_at"> | Finished: {{ new Date(job.finished_at).toLocaleString() }}</span>
        </div>
      </v-card-text>
    </v-card>

    <div v-else class="text-center py-12">
      <v-icon size="64" color="medium-emphasis" class="mb-4">mdi-magnify-scan</v-icon>
      <div class="text-h6 text-medium-emphasis">No scans yet</div>
      <div class="text-body-2 text-medium-emphasis">Click the button above to scan your music directory</div>
    </div>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '../api'
import type { ScanJob } from '../types'

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
    startPolling()
  } catch (e) {
    console.error(e)
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
      }
    } catch {
      clearInterval(pollTimer!)
    }
  }, 2000)
}
</script>
