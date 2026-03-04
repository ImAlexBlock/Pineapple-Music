<template>
  <v-container>
    <div class="d-flex align-center mb-6">
      <div>
        <h1 class="text-h4 font-weight-bold">{{ t('admin.auditLog') }}</h1>
      </div>
    </div>

    <v-card>
      <v-data-table-server
        :headers="headers"
        :items="logs"
        :items-length="total"
        :loading="loading"
        :items-per-page="50"
        hover
        @update:options="loadLogs"
      >
        <template #item.action="{ item }">
          <v-chip size="small" variant="tonal" :color="actionColor(item.action)">{{ item.action }}</v-chip>
        </template>
        <template #item.created_at="{ item }">
          <span class="text-caption">{{ new Date(item.created_at).toLocaleString() }}</span>
        </template>
      </v-data-table-server>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '../api'
import type { AuditLog } from '../types'

const { t } = useI18n()
const logs = ref<AuditLog[]>([])
const total = ref(0)
const loading = ref(false)

const headers = [
  { title: 'Action', key: 'action', sortable: false },
  { title: 'Role', key: 'role', sortable: false },
  { title: 'IP', key: 'ip', sortable: false },
  { title: 'Detail', key: 'detail', sortable: false },
  { title: 'Time', key: 'created_at', sortable: false },
]

function actionColor(action: string): string {
  if (action.includes('login')) return 'info'
  if (action.includes('rotate') || action.includes('bootstrap')) return 'warning'
  if (action.includes('scan') || action.includes('upload')) return 'success'
  if (action.includes('settings')) return 'primary'
  return 'default'
}

async function loadLogs(options: { page: number; itemsPerPage: number }) {
  loading.value = true
  try {
    const { data } = await adminApi.getAuditLogs({
      offset: (options.page - 1) * options.itemsPerPage,
      limit: options.itemsPerPage,
    })
    logs.value = data.items || []
    total.value = data.total
  } catch (e) {
    console.error(e)
  }
  loading.value = false
}
</script>
