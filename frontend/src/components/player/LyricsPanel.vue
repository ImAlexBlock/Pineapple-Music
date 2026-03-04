<template>
  <v-navigation-drawer
    v-model="player.lyricsVisible"
    location="right"
    temporary
    width="400"
  >
    <v-card flat>
      <v-card-title class="d-flex align-center">
        <span>{{ t('player.lyrics') }}</span>
        <v-spacer />
        <v-btn icon size="small" @click="player.toggleLyrics()">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>
      <v-card-text>
        <div v-if="!lyrics" class="text-center text-medium-emphasis">
          {{ t('player.noLyrics') }}
        </div>
        <div v-else-if="lyrics.type === 'synced'" class="lyrics-synced">
          <div
            v-for="(line, i) in parsedLines"
            :key="i"
            :class="{ 'lyrics-active': isActiveLine(i) }"
            class="lyrics-line"
            @click="player.seek(line.time)"
          >
            {{ line.text }}
          </div>
        </div>
        <pre v-else class="lyrics-plain">{{ lyrics.content }}</pre>
      </v-card-text>
    </v-card>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePlayerStore } from '../../stores/player'
import { trackApi } from '../../api'
import type { TrackLyric } from '../../types'

const { t } = useI18n()
const player = usePlayerStore()
const lyrics = ref<TrackLyric | null>(null)

interface LyricLine {
  time: number
  text: string
}

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

function isActiveLine(index: number): boolean {
  const current = player.currentTime
  const line = parsedLines.value[index]
  const next = parsedLines.value[index + 1]
  if (!line) return false
  return current >= line.time && (!next || current < next.time)
}

watch(() => player.currentTrack?.id, async (id) => {
  lyrics.value = null
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
.lyrics-line {
  padding: 4px 0;
  cursor: pointer;
  transition: all 0.3s;
  opacity: 0.5;
}
.lyrics-active {
  opacity: 1;
  font-weight: bold;
  color: rgb(var(--v-theme-primary));
  font-size: 1.1em;
}
.lyrics-plain {
  white-space: pre-wrap;
  font-family: inherit;
}
</style>
