<template>
  <v-container>
    <div class="d-flex align-center mb-6">
      <div>
        <h1 class="text-h4 font-weight-bold">{{ t('nav.tracks') }}</h1>
        <p class="text-body-2 text-medium-emphasis mt-1" v-if="total > 0">{{ total }} tracks</p>
      </div>
    </div>

    <v-text-field
      v-model="search"
      :placeholder="t('search.placeholder')"
      prepend-inner-icon="mdi-magnify"
      clearable
      class="mb-4"
      @update:model-value="debouncedLoad"
    />

    <v-data-table-server
      :headers="headers"
      :items="tracks"
      :items-length="total"
      :loading="loading"
      :items-per-page="50"
      hover
      @update:options="loadTracks"
    >
      <template #item.title="{ item }">
        <div class="d-flex align-center py-1">
          <v-avatar v-if="item.has_cover" size="40" rounded="lg" class="mr-3 flex-shrink-0">
            <v-img :src="coverUrl(item.id)" cover />
          </v-avatar>
          <v-avatar v-else size="40" rounded="lg" color="primary" variant="tonal" class="mr-3 flex-shrink-0">
            <v-icon size="20">mdi-music-note</v-icon>
          </v-avatar>
          <div class="text-truncate">
            <div class="text-body-2 font-weight-medium">{{ item.title }}</div>
          </div>
        </div>
      </template>
      <template #item.duration="{ item }">
        <span class="text-caption text-medium-emphasis">{{ formatDuration(item.duration) }}</span>
      </template>
      <template #item.actions="{ item }">
        <v-btn icon size="small" variant="text" color="primary" @click="playTrack(item)">
          <v-icon>mdi-play-circle-outline</v-icon>
        </v-btn>
      </template>
    </v-data-table-server>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { trackApi } from '../api'
import { usePlayerStore } from '../stores/player'
import type { Track } from '../types'

const { t } = useI18n()
const player = usePlayerStore()
const tracks = ref<Track[]>([])
const total = ref(0)
const loading = ref(false)
const search = ref('')

const headers = computed(() => [
  { title: t('track.title'), key: 'title', sortable: false },
  { title: t('track.artist'), key: 'artist', sortable: false },
  { title: t('track.album'), key: 'album', sortable: false },
  { title: t('track.duration'), key: 'duration', sortable: false, width: '80px' },
  { title: '', key: 'actions', sortable: false, width: '50px' },
])

let debounceTimer: ReturnType<typeof setTimeout>
function debouncedLoad() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => loadTracks({ page: 1, itemsPerPage: 50 }), 300)
}

async function loadTracks(options: { page: number; itemsPerPage: number }) {
  loading.value = true
  try {
    const { data } = await trackApi.list({
      offset: (options.page - 1) * options.itemsPerPage,
      limit: options.itemsPerPage,
      q: search.value || undefined,
    })
    tracks.value = data.items || []
    total.value = data.total
  } catch (e) {
    console.error(e)
  }
  loading.value = false
}

function coverUrl(id: number) {
  return trackApi.coverUrl(id)
}

function playTrack(track: Track) {
  player.play(track, tracks.value)
}

function formatDuration(sec: number): string {
  if (!sec) return '--:--'
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}
</script>
