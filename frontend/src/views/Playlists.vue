<template>
  <v-container>
    <div class="d-flex align-center mb-6">
      <div>
        <h1 class="text-h4 font-weight-bold">{{ t('nav.playlists') }}</h1>
      </div>
      <v-spacer />
      <v-btn v-if="auth.role === 'admin'" color="primary" @click="showCreate = true">
        <v-icon start>mdi-plus</v-icon>
        {{ t('playlist.create') }}
      </v-btn>
    </div>

    <div v-if="playlists.length === 0" class="text-center py-12">
      <v-icon size="64" color="medium-emphasis" class="mb-4">mdi-playlist-music-outline</v-icon>
      <div class="text-h6 text-medium-emphasis">{{ t('playlist.empty') }}</div>
    </div>

    <v-row>
      <v-col v-for="pl in playlists" :key="pl.id" cols="12" sm="6" md="4">
        <v-card :to="`/playlists/${pl.id}`" hover class="pa-4">
          <div class="d-flex align-center">
            <v-avatar color="primary" variant="tonal" rounded="lg" size="48" class="mr-3">
              <v-icon>mdi-playlist-music</v-icon>
            </v-avatar>
            <div>
              <div class="text-body-1 font-weight-medium">{{ pl.name }}</div>
              <div class="text-caption text-medium-emphasis">{{ pl.created_at?.slice(0, 10) }}</div>
            </div>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <v-dialog v-model="showCreate" max-width="420">
      <v-card class="pa-2">
        <v-card-title class="text-h6">{{ t('playlist.create') }}</v-card-title>
        <v-card-text>
          <v-text-field v-model="newName" :label="t('playlist.name')" autofocus @keyup.enter="createPlaylist" />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showCreate = false">Cancel</v-btn>
          <v-btn color="primary" variant="tonal" @click="createPlaylist">Create</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { playlistApi } from '../api'
import { useAuthStore } from '../stores/auth'
import type { Playlist } from '../types'

const { t } = useI18n()
const auth = useAuthStore()
const playlists = ref<Playlist[]>([])
const showCreate = ref(false)
const newName = ref('')

onMounted(loadPlaylists)

async function loadPlaylists() {
  try {
    const { data } = await playlistApi.list()
    playlists.value = data
  } catch (e) {
    console.error(e)
  }
}

async function createPlaylist() {
  if (!newName.value) return
  try {
    await playlistApi.create(newName.value)
    newName.value = ''
    showCreate.value = false
    await loadPlaylists()
  } catch (e) {
    console.error(e)
  }
}
</script>
