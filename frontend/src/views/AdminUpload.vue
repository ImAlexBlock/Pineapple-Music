<template>
  <div>
    <div class="mb-6">
      <h1 class="text-2xl font-bold">{{ t('admin.upload') }}</h1>
      <p class="text-sm text-muted-foreground mt-1">Upload audio files to your library</p>
    </div>

    <Card>
      <CardContent class="pt-6 space-y-4">
        <div class="space-y-2">
          <Label for="audio-files">{{ t('admin.uploadFile') }}</Label>
          <Input id="audio-files" type="file" accept="audio/*" multiple @change="onFileChange" />
        </div>
        <Button :disabled="!files.length || uploading" @click="upload">
          <Loader2 v-if="uploading" class="mr-2 h-4 w-4 animate-spin" />
          <Upload v-else class="mr-2 h-4 w-4" />
          {{ t('admin.upload') }}
        </Button>
      </CardContent>
    </Card>

    <Card v-if="results.length" class="mt-4">
      <CardContent class="pt-6">
        <div class="space-y-2">
          <div v-for="(r, i) in results" :key="i" class="flex items-center gap-3 py-1">
            <CheckCircle v-if="r.ok" class="h-5 w-5 shrink-0 text-success" />
            <AlertCircle v-else class="h-5 w-5 shrink-0 text-destructive" />
            <div class="min-w-0">
              <div class="text-sm font-medium truncate">{{ r.name }}</div>
              <div class="text-xs" :class="r.ok ? 'text-success' : 'text-destructive'">{{ r.message }}</div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '../api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Upload, Loader2, CheckCircle, AlertCircle } from 'lucide-vue-next'

const { t } = useI18n()
const files = ref<File[]>([])
const uploading = ref(false)
const results = ref<{ name: string; ok: boolean; message: string }[]>([])

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  files.value = input.files ? Array.from(input.files) : []
}

async function upload() {
  if (!files.value.length) return
  uploading.value = true
  results.value = []

  for (const file of files.value) {
    try {
      const { data } = await adminApi.upload(file)
      results.value.push({ name: file.name, ok: true, message: data.message })
    } catch (e: unknown) {
      const msg = (e as { response?: { data?: { message?: string } } })?.response?.data?.message || 'Upload failed'
      results.value.push({ name: file.name, ok: false, message: msg })
    }
  }

  uploading.value = false
  files.value = []
}
</script>
