<template>
  <div>
    <div class="mb-4">
      <h1 class="text-2xl font-bold">{{ t('nav.tracks') }}</h1>
      <p v-if="total > 0" class="text-sm text-muted-foreground mt-1">{{ total }} tracks</p>
    </div>

    <Input
      v-model="search"
      :placeholder="t('search.placeholder')"
      class="mb-4"
      @input="debouncedLoad"
    />

    <div class="rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[40%]">{{ t('track.title') }}</TableHead>
            <TableHead>{{ t('track.artist') }}</TableHead>
            <TableHead>{{ t('track.album') }}</TableHead>
            <TableHead class="w-20">{{ t('track.duration') }}</TableHead>
            <TableHead class="w-12"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-if="loading">
            <TableCell :colspan="5" class="text-center py-8">
              <Loader2 class="h-6 w-6 animate-spin mx-auto text-muted-foreground" />
            </TableCell>
          </TableRow>
          <TableRow v-else-if="tracks.length === 0">
            <TableCell :colspan="5" class="text-center py-8 text-muted-foreground">
              No tracks found
            </TableCell>
          </TableRow>
          <TableRow
            v-for="track in tracks"
            :key="track.id"
            class="cursor-pointer"
            @click="playTrack(track)"
          >
            <TableCell>
              <div class="flex items-center gap-3">
                <div class="h-10 w-10 shrink-0 rounded-lg overflow-hidden bg-primary/10 flex items-center justify-center">
                  <img v-if="track.has_cover" :src="coverUrl(track.id)" loading="lazy" class="h-full w-full object-cover" />
                  <Music v-else class="h-4 w-4 text-primary" />
                </div>
                <span class="truncate text-sm font-medium">{{ track.title }}</span>
              </div>
            </TableCell>
            <TableCell class="text-sm">{{ track.artist }}</TableCell>
            <TableCell class="text-sm">{{ track.album }}</TableCell>
            <TableCell class="text-xs text-muted-foreground">{{ formatDuration(track.duration) }}</TableCell>
            <TableCell>
              <Button variant="ghost" size="icon" class="h-8 w-8" @click.stop="playTrack(track)">
                <Play class="h-4 w-4 text-primary" />
              </Button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 mt-4">
      <Button variant="outline" size="sm" :disabled="page <= 1" @click="goToPage(page - 1)">
        <ArrowLeft class="h-4 w-4" />
      </Button>
      <span class="text-sm text-muted-foreground">{{ page }} / {{ totalPages }}</span>
      <Button variant="outline" size="sm" :disabled="page >= totalPages" @click="goToPage(page + 1)">
        <ArrowRight class="h-4 w-4" />
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { trackApi } from '../api'
import { usePlayerStore } from '../stores/player'
import type { Track } from '../types'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Music, Play, ArrowLeft, ArrowRight, Loader2 } from 'lucide-vue-next'

const { t } = useI18n()
const player = usePlayerStore()
const tracks = ref<Track[]>([])
const total = ref(0)
const loading = ref(false)
const search = ref('')
const page = ref(1)
const perPage = 50

const totalPages = computed(() => Math.ceil(total.value / perPage))

onMounted(() => loadTracks())

let debounceTimer: ReturnType<typeof setTimeout>
function debouncedLoad() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    page.value = 1
    loadTracks()
  }, 300)
}

function goToPage(p: number) {
  page.value = p
  loadTracks()
}

async function loadTracks() {
  loading.value = true
  try {
    const { data } = await trackApi.list({
      offset: (page.value - 1) * perPage,
      limit: perPage,
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
