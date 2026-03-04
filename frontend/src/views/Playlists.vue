<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">{{ t('nav.playlists') }}</h1>
      <Button v-if="auth.role === 'admin'" @click="showCreate = true">
        <Plus class="mr-2 h-4 w-4" />
        {{ t('playlist.create') }}
      </Button>
    </div>

    <div v-if="playlists.length === 0" class="text-center py-12">
      <ListMusic class="h-16 w-16 mx-auto mb-4 text-muted-foreground" />
      <p class="text-lg text-muted-foreground">{{ t('playlist.empty') }}</p>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
      <router-link v-for="pl in playlists" :key="pl.id" :to="`/playlists/${pl.id}`" class="block">
        <Card class="p-4 hover:shadow-md transition-shadow cursor-pointer">
          <div class="flex items-center gap-3">
            <div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-primary/10">
              <ListMusic class="h-6 w-6 text-primary" />
            </div>
            <div class="min-w-0">
              <div class="font-medium truncate">{{ pl.name }}</div>
              <div class="text-xs text-muted-foreground">{{ pl.created_at?.slice(0, 10) }}</div>
            </div>
          </div>
        </Card>
      </router-link>
    </div>

    <!-- Create dialog -->
    <Dialog v-model:open="showCreate">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{{ t('playlist.create') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="playlist-name">{{ t('playlist.name') }}</Label>
            <Input id="playlist-name" v-model="newName" @keyup.enter="createPlaylist" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showCreate = false">Cancel</Button>
          <Button @click="createPlaylist">Create</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { playlistApi } from '../api'
import { useAuthStore } from '../stores/auth'
import type { Playlist } from '../types'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Plus, ListMusic } from 'lucide-vue-next'

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
