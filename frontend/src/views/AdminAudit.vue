<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">{{ t('admin.auditLog') }}</h1>

    <Card>
      <div class="rounded-lg border-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Action</TableHead>
              <TableHead>Role</TableHead>
              <TableHead>IP</TableHead>
              <TableHead>Detail</TableHead>
              <TableHead>Time</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-if="loading">
              <TableCell :colspan="5" class="text-center py-8">
                <Loader2 class="h-6 w-6 animate-spin mx-auto text-muted-foreground" />
              </TableCell>
            </TableRow>
            <TableRow v-else-if="logs.length === 0">
              <TableCell :colspan="5" class="text-center py-8 text-muted-foreground">
                No audit logs
              </TableCell>
            </TableRow>
            <TableRow v-for="log in logs" :key="log.id">
              <TableCell>
                <Badge :variant="actionVariant(log.action)">{{ log.action }}</Badge>
              </TableCell>
              <TableCell class="text-sm">{{ log.role }}</TableCell>
              <TableCell class="text-sm">{{ log.ip }}</TableCell>
              <TableCell class="text-sm max-w-[200px] truncate">{{ log.detail }}</TableCell>
              <TableCell class="text-xs text-muted-foreground whitespace-nowrap">{{ new Date(log.created_at).toLocaleString() }}</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>
    </Card>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 mt-4">
      <Button variant="outline" size="sm" :disabled="page <= 1" @click="goToPage(page - 1)">
        <ArrowLeft class="h-4 w-4" />
      </Button>
      <span class="text-sm text-muted-foreground">{{ page }} / {{ totalPages }}</span>
      <Button variant="outline" size="sm" :disabled="page >= totalPages" @click="goToPage(page + 1)">
        <ArrowRight class="h-4 w-4" />
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '../api'
import type { AuditLog } from '../types'
import { Card } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ArrowLeft, ArrowRight, Loader2 } from 'lucide-vue-next'

const { t } = useI18n()
const logs = ref<AuditLog[]>([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const perPage = 50

const totalPages = computed(() => Math.ceil(total.value / perPage))

onMounted(() => loadLogs())

function goToPage(p: number) {
  page.value = p
  loadLogs()
}

function actionVariant(action: string): 'default' | 'secondary' | 'destructive' | 'outline' {
  if (action.includes('login')) return 'default'
  if (action.includes('rotate') || action.includes('bootstrap')) return 'destructive'
  if (action.includes('scan') || action.includes('upload')) return 'secondary'
  return 'outline'
}

async function loadLogs() {
  loading.value = true
  try {
    const { data } = await adminApi.getAuditLogs({
      offset: (page.value - 1) * perPage,
      limit: perPage,
    })
    logs.value = data.items || []
    total.value = data.total
  } catch (e) {
    console.error(e)
  }
  loading.value = false
}
</script>
