<template>
  <v-container>
    <div class="mb-6">
      <h1 class="text-h4 font-weight-bold mb-4">{{ t('nav.search') }}</h1>
      <v-text-field
        v-model="query"
        :placeholder="t('search.placeholder')"
        prepend-inner-icon="mdi-magnify"
        clearable
        autofocus
        @keyup.enter="doSearch"
      />
    </div>

    <div v-if="results.length === 0 && searched" class="text-center py-12">
      <v-icon size="64" color="medium-emphasis" class="mb-4">mdi-magnify-close</v-icon>
      <div class="text-h6 text-medium-emphasis">{{ t('search.noResults') }}</div>
    </div>

    <v-list v-if="results.length > 0" lines="two">
      <v-list-item
        v-for="track in results"
        :key="track.id"
        rounded="lg"
        class="mb-1"
        @click="playTrack(track)"
      >
        <template #prepend>
          <v-avatar v-if="track.has_cover" size="44" rounded="lg" class="mr-3">
            <v-img :src="`/api/v1/tracks/${track.id}/cover`" cover />
          </v-avatar>
          <v-avatar v-else size="44" rounded="lg" color="primary" variant="tonal" class="mr-3">
            <v-icon>mdi-music-note</v-icon>
          </v-avatar>
        </template>
        <v-list-item-title class="font-weight-medium">{{ track.title }}</v-list-item-title>
        <v-list-item-subtitle>{{ track.artist }} - {{ track.album }}</v-list-item-subtitle>
        <template #append>
          <v-btn icon size="small" variant="text" color="primary">
            <v-icon>mdi-play-circle-outline</v-icon>
          </v-btn>
        </template>
      </v-list-item>
    </v-list>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { trackApi } from '../api'
import { usePlayerStore } from '../stores/player'
import type { Track } from '../types'

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
