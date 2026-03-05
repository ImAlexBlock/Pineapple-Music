<template>
  <div
    class="w-full max-w-3xl rounded-2xl overflow-hidden shadow-2xl"
    :class="isLyrics
      ? 'bg-white/[0.08] backdrop-blur-xl border border-white/10'
      : 'bg-background/80 backdrop-blur-xl border border-border/50'"
  >
    <!-- Progress bar -->
    <div
      class="h-1 w-full cursor-pointer group"
      :class="isLyrics ? 'bg-white/10' : 'bg-muted'"
      @click="seekFromProgress"
    >
      <div
        class="h-full transition-all"
        :class="!isLyrics ? 'bg-primary' : ''"
        :style="isLyrics && accentColor
          ? { width: `${progress}%`, backgroundColor: accentColor }
          : isLyrics
            ? { width: `${progress}%`, backgroundColor: 'rgba(255,255,255,0.7)' }
            : { width: `${progress}%` }"
      />
    </div>

    <div class="grid grid-cols-[1fr_auto_1fr] items-center gap-2 px-3 py-2 sm:px-4 sm:py-2.5">
      <!-- Cover + track info (left) -->
      <div class="flex items-center gap-2 sm:gap-3 min-w-0">
        <div
          class="group relative h-10 w-10 sm:h-11 sm:w-11 shrink-0 rounded-lg overflow-hidden cursor-pointer"
          @click="player.toggleLyrics()"
        >
          <Avatar class="h-10 w-10 sm:h-11 sm:w-11" shape="square">
            <AvatarImage v-if="player.currentTrack?.has_cover" :src="`/api/v1/tracks/${player.currentTrack.id}/cover`" />
            <AvatarFallback class="rounded-lg" :class="isLyrics ? 'bg-white/10' : 'bg-primary/10'">
              <Music class="h-5 w-5" :class="isLyrics ? 'text-white/40' : 'text-primary'" />
            </AvatarFallback>
          </Avatar>
          <div class="absolute inset-0 flex items-center justify-center rounded-lg bg-black/50 opacity-0 transition-opacity group-hover:opacity-100">
            <Minimize2 v-if="isLyrics" class="h-4 w-4 text-white" />
            <Maximize2 v-else class="h-4 w-4 text-white" />
          </div>
        </div>
        <div class="min-w-0 hidden sm:block">
          <div
            class="text-sm font-medium truncate max-w-[160px]"
            :class="isLyrics ? 'text-white' : ''"
          >{{ player.currentTrack?.title }}</div>
          <div
            class="text-xs truncate max-w-[160px]"
            :class="isLyrics ? 'text-white/50' : 'text-muted-foreground'"
          >{{ player.currentTrack?.artist }}</div>
        </div>
      </div>

      <!-- Controls (center — always centered via grid) -->
      <div class="flex items-center gap-0.5 sm:gap-1">
        <!-- Play mode - desktop only -->
        <button
          v-if="isLyrics"
          class="h-8 w-8 hidden sm:inline-flex items-center justify-center rounded-lg transition-all active:scale-90"
          :class="isLyrics ? 'text-white/60 hover:text-white hover:bg-white/10' : ''"
          @click="player.cyclePlayMode()"
        >
          <component
            :is="playModeIconComponent"
            class="h-4 w-4"
            :style="player.playMode !== 'sequential' && accentColor ? { color: accentColor } : {}"
            :class="player.playMode !== 'sequential' && !accentColor ? 'text-primary' : ''"
          />
        </button>
        <Button v-else variant="ghost" size="icon" class="h-8 w-8 hidden sm:inline-flex" @click="player.cyclePlayMode()">
          <component
            :is="playModeIconComponent"
            class="h-4 w-4"
            :class="player.playMode !== 'sequential' ? 'text-primary' : ''"
          />
        </Button>

        <!-- Previous -->
        <button
          v-if="isLyrics"
          class="h-8 w-8 flex items-center justify-center rounded-lg text-white transition-all hover:bg-white/10 active:scale-90"
          @click="player.previous()"
        >
          <SkipBack class="h-4 w-4 sm:h-5 sm:w-5 fill-current" />
        </button>
        <Button v-else variant="ghost" size="icon" class="h-8 w-8" @click="player.previous()">
          <SkipBack class="h-4 w-4 sm:h-5 sm:w-5 fill-current" />
        </Button>

        <!-- Play/Pause -->
        <button
          v-if="isLyrics"
          class="flex h-9 w-9 sm:h-10 sm:w-10 items-center justify-center rounded-full text-black transition-all hover:scale-105 active:scale-95 shadow-lg"
          :style="accentColor ? { backgroundColor: accentColor } : {}"
          :class="!accentColor ? 'bg-white' : ''"
          @click="player.togglePlay()"
        >
          <Pause v-if="player.playing" class="h-4 w-4 sm:h-5 sm:w-5 fill-current" />
          <Play v-else class="h-4 w-4 sm:h-5 sm:w-5 fill-current ml-0.5" />
        </button>
        <Button v-else size="icon" class="h-9 w-9 sm:h-10 sm:w-10 rounded-full" @click="player.togglePlay()">
          <Pause v-if="player.playing" class="h-4 w-4 sm:h-5 sm:w-5 fill-current" />
          <Play v-else class="h-4 w-4 sm:h-5 sm:w-5 fill-current" />
        </Button>

        <!-- Next -->
        <button
          v-if="isLyrics"
          class="h-8 w-8 flex items-center justify-center rounded-lg text-white transition-all hover:bg-white/10 active:scale-90"
          @click="player.next()"
        >
          <SkipForward class="h-4 w-4 sm:h-5 sm:w-5 fill-current" />
        </button>
        <Button v-else variant="ghost" size="icon" class="h-8 w-8" @click="player.next()">
          <SkipForward class="h-4 w-4 sm:h-5 sm:w-5 fill-current" />
        </Button>
      </div>

      <!-- Time + Volume (right) -->
      <div class="flex items-center gap-2 justify-end">
        <span
          class="text-xs whitespace-nowrap"
          :class="isLyrics ? 'text-white/50' : 'text-muted-foreground'"
        >
          {{ formatTime(player.currentTime) }}
        </span>
        <div class="items-center gap-2 hidden md:flex" :class="isLyrics ? 'lyrics-slider' : ''">
          <component :is="volumeIcon" class="h-4 w-4 shrink-0" :class="isLyrics ? 'text-white/50' : 'text-muted-foreground'" />
          <Slider
            :model-value="[player.volume]"
            :min="0"
            :max="1"
            :step="0.01"
            class="w-[80px]"
            @update:model-value="(v: number[]) => player.setVolume(v[0])"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayerStore } from '../../stores/player'
import { Button } from '@/components/ui/button'
import { Slider } from '@/components/ui/slider'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  Music, Play, Pause, SkipBack, SkipForward, Maximize2, Minimize2,
  ArrowRight, Repeat, Repeat1, Shuffle,
  Volume2, Volume1, VolumeX,
} from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  variant?: 'default' | 'lyrics'
  accentColor?: string | null
}>(), {
  variant: 'default',
  accentColor: null,
})

const player = usePlayerStore()

const isLyrics = computed(() => props.variant === 'lyrics')

const progress = computed(() => {
  if (!player.duration) return 0
  return (player.currentTime / player.duration) * 100
})

const playModeIconComponent = computed(() => {
  const map: Record<string, typeof ArrowRight> = {
    ArrowRight, Repeat, Repeat1, Shuffle,
  }
  return map[player.playModeIcon] || ArrowRight
})

const volumeIcon = computed(() => {
  if (player.volume === 0) return VolumeX
  if (player.volume < 0.5) return Volume1
  return Volume2
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
/* Override slider colors for dark lyrics context */
.lyrics-slider :deep(span[data-orientation="horizontal"]) {
  background: rgba(255, 255, 255, 0.15) !important;
}
.lyrics-slider :deep(span[role="slider"]) {
  border-color: var(--lyrics-accent, rgba(255, 255, 255, 0.8)) !important;
  background: white !important;
  height: 14px !important;
  width: 14px !important;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.4);
}
.lyrics-slider :deep(span[data-orientation="horizontal"] > span) {
  background: var(--lyrics-accent, rgba(255, 255, 255, 0.7)) !important;
}
</style>
