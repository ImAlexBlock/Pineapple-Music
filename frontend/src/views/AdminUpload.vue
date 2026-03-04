<template>
  <v-container>
    <div class="d-flex align-center mb-6">
      <div>
        <h1 class="text-h4 font-weight-bold">{{ t('admin.upload') }}</h1>
        <p class="text-body-2 text-medium-emphasis mt-1">Upload audio files to your library</p>
      </div>
    </div>

    <v-card class="pa-2">
      <v-card-text>
        <v-file-input
          v-model="files"
          :label="t('admin.uploadFile')"
          accept="audio/*"
          prepend-icon="mdi-music"
          multiple
          show-size
        />
      </v-card-text>
      <v-card-actions class="px-4 pb-4">
        <v-btn color="primary" variant="tonal" :loading="uploading" :disabled="!files?.length" @click="upload">
          <v-icon start>mdi-upload-outline</v-icon>
          {{ t('admin.upload') }}
        </v-btn>
      </v-card-actions>
    </v-card>

    <v-card v-if="results.length" class="mt-4">
      <v-list>
        <v-list-item v-for="(r, i) in results" :key="i">
          <template #prepend>
            <v-icon :color="r.ok ? 'success' : 'error'">{{ r.ok ? 'mdi-check-circle' : 'mdi-alert-circle' }}</v-icon>
          </template>
          <v-list-item-title>{{ r.name }}</v-list-item-title>
          <v-list-item-subtitle :class="r.ok ? 'text-success' : 'text-error'">{{ r.message }}</v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminApi } from '../api'

const { t } = useI18n()
const files = ref<File[]>([])
const uploading = ref(false)
const results = ref<{ name: string; ok: boolean; message: string }[]>([])

async function upload() {
  if (!files.value?.length) return
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
