<template>
  <v-container>
    <div class="d-flex align-center mb-6">
      <v-btn icon variant="text" @click="$router.back()">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <h1 class="text-h4 font-weight-bold ml-2">{{ playlist?.name }}</h1>
      <v-spacer />
      <v-btn v-if="auth.role === 'admin'" icon variant="text" color="error" @click="deletePlaylist">
        <v-icon>mdi-delete-outline</v-icon>
        <v-tooltip activator="parent" location="bottom">{{ t('playlist.delete') }}</v-tooltip>
      </v-btn>
    </div>

    <div v-if="!playlist?.tracks?.length" class="text-center py-12">
      <v-icon size="64" color="medium-emphasis" class="mb-4">mdi-playlist-music-outline</v-icon>
      <div class="text-h6 text-medium-emphasis">{{ t('playlist.empty') }}</div>
    </div>

    <v-list v-else lines="two">
      <v-list-item
        v-for="pt in playlist.tracks"
        :key="pt.id"
        rounded="lg"
        class="mb-1"
        @click="playTrack(pt.track)"
      >
        <template #prepend>
          <v-avatar v-if="pt.track.has_cover" size="44" rounded="lg" class="mr-3">
            <v-img :src="`/api/v1/tracks/${pt.track.id}/cover`" cover />
          </v-avatar>
          <v-avatar v-else size="44" rounded="lg" color="primary" variant="tonal" class="mr-3">
            <v-icon>mdi-music-note</v-icon>
          </v-avatar>
        </template>
        <v-list-item-title class="font-weight-medium">{{ pt.track.title }}</v-list-item-title>
        <v-list-item-subtitle>{{ pt.track.artist }}</v-list-item-subtitle>
        <template #append>
          <v-btn icon size="small" variant="text" color="primary" @click.stop="playTrack(pt.track)">
            <v-icon>mdi-play-circle-outline</v-icon>
          </v-btn>
        </template>
      </v-list-item>
    </v-list>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { playlistApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { usePlayerStore } from '../stores/player'
import type { Playlist, Track } from '../types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const player = usePlayerStore()
const playlist = ref<Playlist | null>(null)

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    const { data } = await playlistApi.get(id)
    playlist.value = data
  } catch (e) {
    console.error(e)
  }
})

function playTrack(track: Track) {
  const tracks = playlist.value?.tracks?.map((pt) => pt.track) || []
  player.play(track, tracks)
}

async function deletePlaylist() {
  if (!playlist.value) return
  try {
    await playlistApi.delete(playlist.value.id)
    router.push('/playlists')
  } catch (e) {
    console.error(e)
  }
}
</script>
