<template>
  <Sheet v-model:open="player.lyricsVisible">
    <SheetContent side="right" class="w-[400px] p-0 flex flex-col">
      <SheetHeader class="p-4 pb-2 flex flex-row items-center justify-between shrink-0">
        <SheetTitle>{{ t('player.lyrics') }}</SheetTitle>
        <SheetDescription class="sr-only">Song lyrics panel</SheetDescription>
      </SheetHeader>
      <ScrollArea class="flex-1 px-4 pb-4">
        <div v-if="!lyrics" class="text-center text-muted-foreground py-8">
          {{ t('player.noLyrics') }}
        </div>
        <div v-else-if="lyrics.type === 'synced'" class="space-y-1">
          <div
            v-for="(line, i) in parsedLines"
            :key="i"
            class="cursor-pointer rounded px-2 py-1 transition-all text-sm"
            :class="isActiveLine(i) ? 'text-primary font-bold text-base' : 'text-muted-foreground'"
            @click="player.seek(line.time)"
          >
            {{ line.text }}
          </div>
        </div>
        <pre v-else class="whitespace-pre-wrap font-sans text-sm">{{ lyrics.content }}</pre>
      </ScrollArea>
    </SheetContent>
  </Sheet>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePlayerStore } from '../../stores/player'
import { trackApi } from '../../api'
import type { TrackLyric } from '../../types'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from '@/components/ui/sheet'
import { ScrollArea } from '@/components/ui/scroll-area'

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
