<template>
  <Transition name="player">
    <div v-if="player.currentTrack" class="fixed bottom-0 left-0 right-0 z-50 border-t bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/80">
      <!-- Progress bar (clickable) -->
      <div
        class="h-1 w-full cursor-pointer bg-muted group"
        @click="seekFromProgress"
      >
        <div
          class="h-full bg-primary transition-all"
          :style="{ width: `${progress}%` }"
        />
      </div>

      <div class="flex items-center gap-2 px-3 py-2">
        <!-- Track info (left) -->
        <div class="flex items-center gap-3 min-w-0 w-[180px] max-w-[280px] shrink-0">
          <Avatar class="h-11 w-11 shrink-0 rounded-lg">
            <AvatarImage v-if="player.currentTrack.has_cover" :src="`/api/v1/tracks/${player.currentTrack.id}/cover`" />
            <AvatarFallback class="rounded-lg bg-primary/10">
              <Music class="h-5 w-5 text-primary" />
            </AvatarFallback>
          </Avatar>
          <div class="min-w-0">
            <div class="text-sm font-medium truncate">{{ player.currentTrack.title }}</div>
            <div class="text-xs text-muted-foreground truncate">{{ player.currentTrack.artist }}</div>
          </div>
        </div>

        <div class="flex-1" />

        <!-- Controls (center) -->
        <div class="flex items-center gap-1">
          <Button variant="ghost" size="icon" class="h-8 w-8" @click="player.cyclePlayMode()">
            <component
              :is="playModeIconComponent"
              class="h-4 w-4"
              :class="player.playMode !== 'sequential' ? 'text-primary' : ''"
            />
          </Button>
          <Button variant="ghost" size="icon" class="h-8 w-8" @click="player.previous()">
            <SkipBack class="h-5 w-5" />
          </Button>
          <Button size="icon" class="h-10 w-10 rounded-full" @click="player.togglePlay()">
            <Pause v-if="player.playing" class="h-5 w-5" />
            <Play v-else class="h-5 w-5" />
          </Button>
          <Button variant="ghost" size="icon" class="h-8 w-8" @click="player.next()">
            <SkipForward class="h-5 w-5" />
          </Button>
          <Button variant="ghost" size="icon" class="h-8 w-8" @click="player.toggleLyrics()">
            <Type class="h-4 w-4" :class="player.lyricsVisible ? 'text-primary' : ''" />
          </Button>
        </div>

        <div class="flex-1" />

        <!-- Time + Volume (right) -->
        <div class="flex items-center gap-2 shrink-0 w-[180px] justify-end">
          <span class="text-xs text-muted-foreground whitespace-nowrap">
            {{ formatTime(player.currentTime) }} / {{ formatTime(player.duration) }}
          </span>
          <component :is="volumeIcon" class="h-4 w-4 text-muted-foreground shrink-0" />
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
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayerStore } from '../../stores/player'
import { Button } from '@/components/ui/button'
import { Slider } from '@/components/ui/slider'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  Music, Play, Pause, SkipBack, SkipForward, Type,
  ArrowRight, Repeat, Repeat1, Shuffle,
  Volume2, Volume1, VolumeX,
} from 'lucide-vue-next'

const player = usePlayerStore()

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
