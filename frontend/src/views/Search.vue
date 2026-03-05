<template>
  <div>
    <h1 class="text-2xl font-bold mb-4">{{ t('nav.search') }}</h1>

    <div v-if="results.length === 0 && searched" class="text-center py-12">
      <SearchIcon class="h-16 w-16 mx-auto mb-4 text-muted-foreground" />
      <p class="text-lg text-muted-foreground">{{ t('search.noResults') }}</p>
    </div>

    <div v-if="results.length > 0" class="space-y-1">
      <button
        v-for="track in results"
        :key="track.id"
        class="flex w-full items-center gap-3 rounded-lg p-2 text-left transition-colors hover:bg-accent"
        @click="playTrack(track)"
      >
        <div class="h-11 w-11 shrink-0 rounded-lg overflow-hidden bg-primary/10 flex items-center justify-center">
          <img v-if="track.has_cover" :src="`/api/v1/tracks/${track.id}/cover`" loading="lazy" class="h-full w-full object-cover" />
          <Music v-else class="h-5 w-5 text-primary" />
        </div>
        <div class="min-w-0 flex-1">
          <div class="truncate text-sm font-medium">{{ track.title }}</div>
          <div class="truncate text-xs text-muted-foreground">{{ track.artist }} - {{ track.album }}</div>
        </div>
        <Play class="h-5 w-5 shrink-0 text-primary" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { trackApi } from '../api'
import { usePlayerStore } from '../stores/player'
import type { Track } from '../types'
import { Search as SearchIcon, Music, Play } from 'lucide-vue-next'

const { t } = useI18n()
const route = useRoute()
const player = usePlayerStore()
const results = ref<Track[]>([])
const searched = ref(false)

async function doSearch(q: string) {
  if (!q) return
  searched.value = true
  try {
    const { data } = await trackApi.list({ q, limit: 50 })
    results.value = data.items || []
  } catch (e) {
    console.error(e)
  }
}

// React to query param changes (from header search bar)
watch(() => route.query.q, (q) => {
  if (q && typeof q === 'string') {
    doSearch(q)
  }
}, { immediate: true })

function playTrack(track: Track) {
  player.play(track, results.value)
}
</script>
