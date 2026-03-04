<template>
  <div>
    <h1 class="text-2xl font-bold mb-4">{{ t('nav.search') }}</h1>
    <div class="flex gap-2 mb-6">
      <Input
        v-model="query"
        :placeholder="t('search.placeholder')"
        class="flex-1"
        autofocus
        @keyup.enter="doSearch"
      />
      <Button @click="doSearch">
        <SearchIcon class="h-4 w-4" />
      </Button>
    </div>

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
        <Avatar class="h-11 w-11 shrink-0 rounded-lg">
          <AvatarImage v-if="track.has_cover" :src="`/api/v1/tracks/${track.id}/cover`" />
          <AvatarFallback class="rounded-lg bg-primary/10">
            <Music class="h-5 w-5 text-primary" />
          </AvatarFallback>
        </Avatar>
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
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { trackApi } from '../api'
import { usePlayerStore } from '../stores/player'
import type { Track } from '../types'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Search as SearchIcon, Music, Play } from 'lucide-vue-next'

const { t } = useI18n()
const player = usePlayerStore()
const query = ref('')
const results = ref<Track[]>([])
const searched = ref(false)

async function doSearch() {
  if (!query.value) return
  searched.value = true
  try {
    const { data } = await trackApi.list({ q: query.value, limit: 50 })
    results.value = data.items || []
  } catch (e) {
    console.error(e)
  }
}

function playTrack(track: Track) {
  player.play(track, results.value)
}
</script>
