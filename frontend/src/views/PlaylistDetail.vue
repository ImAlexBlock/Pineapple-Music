<template>
  <div>
    <div class="flex items-center gap-2 mb-6">
      <Button variant="ghost" size="icon" @click="$router.back()">
        <ArrowLeft class="h-5 w-5" />
      </Button>
      <h1 class="text-2xl font-bold flex-1">{{ playlist?.name }}</h1>
      <Button v-if="auth.role === 'admin'" variant="ghost" size="icon" class="text-destructive" @click="deletePlaylist">
        <Trash2 class="h-5 w-5" />
      </Button>
    </div>

    <div v-if="!playlist?.tracks?.length" class="text-center py-12">
      <ListMusic class="h-16 w-16 mx-auto mb-4 text-muted-foreground" />
      <p class="text-lg text-muted-foreground">{{ t('playlist.empty') }}</p>
    </div>

    <div v-else class="space-y-1">
      <div
        v-for="(pt, index) in playlist.tracks"
        :key="pt.id"
        class="flex w-full items-center gap-3 rounded-lg p-2 text-left transition-colors hover:bg-accent group"
        :class="dragOverIndex === index ? 'border-t-2 border-primary' : ''"
        draggable="true"
        @dragstart="onDragStart(index, $event)"
        @dragover.prevent="onDragOver(index)"
        @dragend="onDragEnd"
        @click="playTrack(pt.track)"
      >
        <GripVertical class="h-4 w-4 shrink-0 text-muted-foreground opacity-0 group-hover:opacity-100 cursor-grab transition-opacity" />
        <div class="h-11 w-11 shrink-0 rounded-lg overflow-hidden bg-primary/10 flex items-center justify-center">
          <img v-if="pt.track.has_cover" :src="`/api/v1/tracks/${pt.track.id}/cover`" loading="lazy" class="h-full w-full object-cover" />
          <Music v-else class="h-5 w-5 text-primary" />
        </div>
        <div class="min-w-0 flex-1">
          <div class="truncate text-sm font-medium">{{ pt.track.title }}</div>
          <div class="truncate text-xs text-muted-foreground">{{ pt.track.artist }}</div>
        </div>
        <Play class="h-5 w-5 shrink-0 text-primary" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { playlistApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { usePlayerStore } from '../stores/player'
import type { Playlist, Track } from '../types'
import { Button } from '@/components/ui/button'
import { ArrowLeft, Trash2, ListMusic, Music, Play, GripVertical } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const player = usePlayerStore()
const playlist = ref<Playlist | null>(null)
const dragIndex = ref(-1)
const dragOverIndex = ref(-1)

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    const { data } = await playlistApi.get(id)
    playlist.value = data
  } catch (e) {
    console.error(e)
  }
})

function onDragStart(index: number, e: DragEvent) {
  dragIndex.value = index
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
  }
}

function onDragOver(index: number) {
  dragOverIndex.value = index
}

async function onDragEnd() {
  const from = dragIndex.value
  const to = dragOverIndex.value
  dragIndex.value = -1
  dragOverIndex.value = -1

  if (from < 0 || to < 0 || from === to || !playlist.value?.tracks) return

  const tracks = [...playlist.value.tracks]
  const [moved] = tracks.splice(from, 1)
  tracks.splice(to, 0, moved)
  playlist.value.tracks = tracks

  try {
    await playlistApi.reorder(playlist.value.id, tracks.map(pt => pt.track_id))
  } catch (e) {
    console.error(e)
    toast.error(t('error.generic'))
  }
}

function playTrack(track: Track) {
  const tracks = playlist.value?.tracks?.map((pt) => pt.track) || []
  player.play(track, tracks)
}

async function deletePlaylist() {
  if (!playlist.value) return
  try {
    await playlistApi.delete(playlist.value.id)
    toast.success(t('playlist.deleted'))
    router.push('/playlists')
  } catch (e) {
    console.error(e)
    toast.error(t('error.generic'))
  }
}
</script>
