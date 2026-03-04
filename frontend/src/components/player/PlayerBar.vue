<template>
  <v-slide-y-reverse-transition>
    <v-footer v-if="player.currentTrack" app class="pa-0" style="z-index: 100;">
      <v-card width="100%" flat rounded="0" class="player-bar" color="surface">
        <!-- Progress bar -->
        <v-progress-linear
          :model-value="progress"
          color="primary"
          height="3"
          class="cursor-pointer"
          @click="seekFromProgress"
        />
        <div class="d-flex align-center px-3 py-2 ga-2">
          <!-- Track info -->
          <div class="d-flex align-center flex-shrink-0" style="min-width: 180px; max-width: 280px;">
            <v-avatar v-if="player.currentTrack.has_cover" size="44" rounded="lg" class="mr-3 flex-shrink-0">
              <v-img :src="`/api/v1/tracks/${player.currentTrack.id}/cover`" cover />
            </v-avatar>
            <v-avatar v-else size="44" rounded="lg" color="primary" variant="tonal" class="mr-3 flex-shrink-0">
              <v-icon>mdi-music-note</v-icon>
            </v-avatar>
            <div class="overflow-hidden">
              <div class="text-body-2 font-weight-medium text-truncate">{{ player.currentTrack.title }}</div>
              <div class="text-caption text-medium-emphasis text-truncate">{{ player.currentTrack.artist }}</div>
            </div>
          </div>

          <v-spacer />

          <!-- Controls (center) -->
          <div class="d-flex align-center ga-1">
            <v-btn icon variant="text" size="small" @click="player.cyclePlayMode()" density="comfortable">
              <v-icon :color="player.playMode !== 'sequential' ? 'primary' : undefined" size="20">{{ player.playModeIcon }}</v-icon>
              <v-tooltip activator="parent" location="top">{{ player.playModeLabel }}</v-tooltip>
            </v-btn>
            <v-btn icon variant="text" size="small" @click="player.previous()" density="comfortable">
              <v-icon size="22">mdi-skip-previous</v-icon>
            </v-btn>
            <v-btn icon color="primary" size="44" @click="player.togglePlay()" elevation="0">
              <v-icon size="28">{{ player.playing ? 'mdi-pause' : 'mdi-play' }}</v-icon>
            </v-btn>
            <v-btn icon variant="text" size="small" @click="player.next()" density="comfortable">
              <v-icon size="22">mdi-skip-next</v-icon>
            </v-btn>
            <v-btn icon variant="text" size="small" @click="player.toggleLyrics()" density="comfortable">
              <v-icon :color="player.lyricsVisible ? 'primary' : undefined" size="20">mdi-text</v-icon>
              <v-tooltip activator="parent" location="top">Lyrics</v-tooltip>
            </v-btn>
          </div>

          <v-spacer />

          <!-- Time + Volume (right) -->
          <div class="d-flex align-center ga-2 flex-shrink-0" style="min-width: 180px; justify-content: flex-end;">
            <span class="text-caption text-medium-emphasis text-no-wrap">
              {{ formatTime(player.currentTime) }} / {{ formatTime(player.duration) }}
            </span>
            <v-icon size="18" class="text-medium-emphasis">
              {{ player.volume === 0 ? 'mdi-volume-off' : player.volume < 0.5 ? 'mdi-volume-medium' : 'mdi-volume-high' }}
            </v-icon>
            <v-slider
              :model-value="player.volume"
              min="0"
              max="1"
              step="0.01"
              hide-details
              density="compact"
              color="primary"
              track-color="surface-variant"
              style="max-width: 90px;"
              @update:model-value="(v: number) => player.setVolume(v)"
            />
          </div>
        </div>
      </v-card>
    </v-footer>
  </v-slide-y-reverse-transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayerStore } from '../../stores/player'

const player = usePlayerStore()

const progress = computed(() => {
  if (!player.duration) return 0
  return (player.currentTime / player.duration) * 100
})

function seekFromProgress(e: MouseEvent) {
  const target = e.currentTarget as HTMLElement
  const rect = target.getBoundingClientRect()
  const pct = (e.clientX - rect.left) / rect.width
  player.seek(pct * player.duration)
}

function formatTime(sec: number): string {
  if (!sec || isNaN(sec)) return '0:00'
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}
</script>

<style scoped>
.player-bar {
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.cursor-pointer {
  cursor: pointer;
}
</style>
