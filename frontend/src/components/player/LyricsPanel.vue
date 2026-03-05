<template>
  <Teleport to="body">
    <Transition name="lyrics-fullscreen">
      <div
        v-if="player.lyricsVisible && player.currentTrack"
        class="fixed inset-0 z-[100] flex flex-col overflow-hidden"
        :style="accentVars"
      >
        <!-- Opaque background with blurred cover tint -->
        <div class="absolute inset-0 bg-black">
          <img
            v-if="player.currentTrack.has_cover"
            ref="coverImg"
            :src="`/api/v1/tracks/${player.currentTrack.id}/cover`"
            crossorigin="anonymous"
            class="absolute inset-0 h-full w-full object-cover scale-110 blur-[80px] saturate-150 opacity-40"
            @load="onCoverLoad"
          />
        </div>

        <!-- Top bar -->
        <div class="relative z-10 flex items-center justify-between px-6 pt-6 pb-2 shrink-0">
          <button
            class="rounded-full bg-white/10 p-2 backdrop-blur-sm transition-all hover:bg-white/20 active:scale-95"
            @click="player.toggleLyrics()"
          >
            <ChevronDown class="h-5 w-5 text-white" />
          </button>
          <div class="text-center min-w-0 flex-1 mx-4">
            <div class="text-sm font-semibold text-white truncate">{{ player.currentTrack.title }}</div>
            <div class="text-xs text-white/60 truncate">{{ player.currentTrack.artist }}</div>
          </div>
          <div class="w-9" />
        </div>

        <!-- Lyrics body -->
        <div ref="lyricsContainer" class="relative z-10 flex-1 overflow-y-auto px-6 py-8 lyrics-scroll">
          <!-- No lyrics -->
          <div v-if="!lyrics" class="flex h-full items-center justify-center">
            <div class="text-center">
              <Music class="mx-auto mb-4 h-16 w-16 text-white/30" />
              <p class="text-lg text-white/50">{{ t('player.noLyrics') }}</p>
            </div>
          </div>

          <!-- Synced lyrics -->
          <div v-else-if="lyrics.type === 'synced'" class="mx-auto max-w-2xl space-y-4 py-[30vh]">
            <div
              v-for="(group, gi) in lineGroups"
              :key="gi"
              :ref="el => setGroupRef(gi, el as HTMLElement)"
              class="lyrics-line cursor-pointer rounded-xl px-4 py-2 transition-all duration-500 ease-out"
              :class="getGroupClass(gi)"
              @click="player.seek(group.time)"
            >
              <div v-for="(text, ti) in group.texts" :key="ti" :class="ti > 0 ? 'mt-1 text-[0.85em] opacity-80' : ''">
                {{ text }}
              </div>
            </div>
          </div>

          <!-- Plain lyrics -->
          <div v-else class="mx-auto max-w-2xl py-[30vh]">
            <pre class="whitespace-pre-wrap font-sans text-xl leading-relaxed text-white/80">{{ lyrics.content }}</pre>
          </div>
        </div>

        <!-- Bottom floating controls — reuse PlayerBar -->
        <div class="relative z-10 shrink-0 flex justify-center px-2 pb-4 sm:px-4 sm:pb-6 pt-2">
          <PlayerBar variant="lyrics" :accent-color="accentColor" />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePlayerStore } from '../../stores/player'
import { trackApi } from '../../api'
import PlayerBar from './PlayerBar.vue'
import type { TrackLyric } from '../../types'
import { ChevronDown, Music } from 'lucide-vue-next'

const { t } = useI18n()
const player = usePlayerStore()
const lyrics = ref<TrackLyric | null>(null)
const lyricsContainer = ref<HTMLElement | null>(null)
const groupRefs = ref<Map<number, HTMLElement>>(new Map())
const coverImg = ref<HTMLImageElement | null>(null)
const accentColor = ref<string | null>(null)

interface LyricLine {
  time: number
  text: string
}

interface LyricGroup {
  time: number
  texts: string[]
}

function setGroupRef(index: number, el: HTMLElement | null) {
  if (el) {
    groupRefs.value.set(index, el)
  }
}

// --- Color extraction from album cover ---
function extractDominantColor(img: HTMLImageElement): string | null {
  try {
    const canvas = document.createElement('canvas')
    const size = 64
    canvas.width = size
    canvas.height = size
    const ctx = canvas.getContext('2d', { willReadFrequently: true })
    if (!ctx) return null

    ctx.drawImage(img, 0, 0, size, size)
    const data = ctx.getImageData(0, 0, size, size).data

    const buckets = new Map<string, { r: number; g: number; b: number; count: number; saturation: number }>()

    for (let i = 0; i < data.length; i += 16) {
      const r = data[i]
      const g = data[i + 1]
      const b = data[i + 2]

      const brightness = (r + g + b) / 3
      if (brightness < 30 || brightness > 230) continue

      const qr = Math.round(r / 32) * 32
      const qg = Math.round(g / 32) * 32
      const qb = Math.round(b / 32) * 32
      const key = `${qr},${qg},${qb}`

      const max = Math.max(r, g, b)
      const min = Math.min(r, g, b)
      const sat = max === 0 ? 0 : (max - min) / max

      const existing = buckets.get(key)
      if (existing) {
        existing.r += r
        existing.g += g
        existing.b += b
        existing.count++
        existing.saturation = Math.max(existing.saturation, sat)
      } else {
        buckets.set(key, { r, g, b, count: 1, saturation: sat })
      }
    }

    if (buckets.size === 0) return null

    let best: { r: number; g: number; b: number; count: number; saturation: number } | null = null
    let bestScore = 0

    for (const bucket of buckets.values()) {
      const score = bucket.count * (1 + bucket.saturation * 3)
      if (score > bestScore) {
        bestScore = score
        best = bucket
      }
    }

    if (!best) return null

    const r = Math.round(best.r / best.count)
    const g = Math.round(best.g / best.count)
    const b = Math.round(best.b / best.count)

    const max = Math.max(r, g, b)
    const min = Math.min(r, g, b)
    if (max - min < 20) return null

    const lum = 0.299 * r + 0.587 * g + 0.114 * b
    const factor = lum < 120 ? 120 / lum : 1
    const fr = Math.min(255, Math.round(r * factor))
    const fg = Math.min(255, Math.round(g * factor))
    const fb = Math.min(255, Math.round(b * factor))

    return `rgb(${fr}, ${fg}, ${fb})`
  } catch {
    return null
  }
}

function onCoverLoad() {
  if (coverImg.value) {
    accentColor.value = extractDominantColor(coverImg.value)
  }
}

const accentVars = computed(() => {
  if (!accentColor.value) return {}
  return { '--lyrics-accent': accentColor.value } as Record<string, string>
})

const parsedLines = computed<LyricLine[]>(() => {
  if (!lyrics.value || lyrics.value.type !== 'synced') return []
  const lines: LyricLine[] = []
  const timeReg = /\[(\d{1,2}):(\d{2})(?:\.(\d{1,3}))?\]/g
  for (const rawLine of lyrics.value.content.split('\n')) {
    let match: RegExpExecArray | null
    let lastIndex = 0
    const times: number[] = []
    while ((match = timeReg.exec(rawLine)) !== null) {
      const min = parseInt(match[1])
      const sec = parseInt(match[2])
      const ms = match[3] ? parseInt(match[3].padEnd(3, '0')) : 0
      times.push(min * 60 + sec + ms / 1000)
      lastIndex = timeReg.lastIndex
    }
    timeReg.lastIndex = 0
    const text = rawLine.slice(lastIndex).trim()
    if (text) {
      for (const t of times) {
        lines.push({ time: t, text })
      }
    }
  }
  return lines.sort((a, b) => a.time - b.time)
})

const lineGroups = computed<LyricGroup[]>(() => {
  const groups: LyricGroup[] = []
  for (const line of parsedLines.value) {
    const last = groups[groups.length - 1]
    if (last && Math.abs(last.time - line.time) < 0.05) {
      last.texts.push(line.text)
    } else {
      groups.push({ time: line.time, texts: [line.text] })
    }
  }
  return groups
})

const activeGroupIndex = computed(() => {
  const current = player.currentTime
  for (let i = lineGroups.value.length - 1; i >= 0; i--) {
    if (current >= lineGroups.value[i].time) return i
  }
  return -1
})

function getGroupClass(index: number): string {
  const active = activeGroupIndex.value
  if (index === active) {
    return 'font-bold text-2xl scale-[1.02] origin-left lyrics-active-line'
  }
  const dist = Math.abs(index - active)
  if (dist === 1) return 'text-white/40 text-xl'
  if (dist === 2) return 'text-white/25 text-lg'
  return 'text-white/15 text-lg'
}

watch(activeGroupIndex, (idx) => {
  if (idx < 0) return
  const el = groupRefs.value.get(idx)
  const container = lyricsContainer.value
  if (!el || !container) return
  const containerRect = container.getBoundingClientRect()
  const elRect = el.getBoundingClientRect()
  const targetY = elRect.top - containerRect.top - containerRect.height / 3
  container.scrollBy({ top: targetY, behavior: 'smooth' })
})

watch(() => player.currentTrack?.id, async (id) => {
  lyrics.value = null
  groupRefs.value.clear()
  accentColor.value = null
  if (!id) return
  try {
    const { data } = await trackApi.getLyrics(id)
    lyrics.value = data
  } catch {
    lyrics.value = null
  }
})
</script>

<style scoped>
.lyrics-scroll {
  scrollbar-width: none;
}
.lyrics-scroll::-webkit-scrollbar {
  display: none;
}

.lyrics-line {
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.lyrics-active-line {
  color: var(--lyrics-accent, white);
}

/* Fullscreen transition */
.lyrics-fullscreen-enter-active {
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}
.lyrics-fullscreen-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.6, 1);
}
.lyrics-fullscreen-enter-from {
  opacity: 0;
  transform: translateY(100%) scale(0.95);
}
.lyrics-fullscreen-leave-to {
  opacity: 0;
  transform: translateY(60%) scale(0.98);
}
</style>
