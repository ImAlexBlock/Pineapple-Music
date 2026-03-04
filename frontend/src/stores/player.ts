import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Track } from '../types'
import { trackApi, playEventApi } from '../api'

export type PlayMode = 'sequential' | 'repeat-all' | 'repeat-one' | 'shuffle'

export const usePlayerStore = defineStore('player', () => {
  const queue = ref<Track[]>([])
  const currentIndex = ref(-1)
  const playing = ref(false)
  const currentTime = ref(0)
  const duration = ref(0)
  const volume = ref(1)
  const playMode = ref<PlayMode>('sequential')
  const lyricsVisible = ref(false)

  const currentTrack = computed(() =>
    currentIndex.value >= 0 && currentIndex.value < queue.value.length
      ? queue.value[currentIndex.value]
      : null
  )

  const playModeIcon = computed(() => {
    switch (playMode.value) {
      case 'sequential': return 'mdi-arrow-right'
      case 'repeat-all': return 'mdi-repeat'
      case 'repeat-one': return 'mdi-repeat-once'
      case 'shuffle': return 'mdi-shuffle-variant'
    }
  })

  const playModeLabel = computed(() => {
    switch (playMode.value) {
      case 'sequential': return 'Sequential'
      case 'repeat-all': return 'Repeat All'
      case 'repeat-one': return 'Repeat One'
      case 'shuffle': return 'Shuffle'
    }
  })

  // Audio element singleton
  let audio: HTMLAudioElement | null = null

  function getAudio(): HTMLAudioElement {
    if (!audio) {
      audio = new Audio()
      audio.addEventListener('timeupdate', () => {
        currentTime.value = audio!.currentTime
      })
      audio.addEventListener('loadedmetadata', () => {
        duration.value = audio!.duration
      })
      audio.addEventListener('ended', () => {
        next()
      })
      audio.addEventListener('play', () => {
        playing.value = true
      })
      audio.addEventListener('pause', () => {
        playing.value = false
      })
    }
    return audio
  }

  function play(track: Track, tracks?: Track[]) {
    if (tracks) {
      queue.value = [...tracks]
      currentIndex.value = tracks.findIndex((t) => t.id === track.id)
    } else {
      const idx = queue.value.findIndex((t) => t.id === track.id)
      if (idx >= 0) {
        currentIndex.value = idx
      } else {
        queue.value.push(track)
        currentIndex.value = queue.value.length - 1
      }
    }

    const a = getAudio()
    a.src = trackApi.streamUrl(track.id)
    a.volume = volume.value
    a.play()

    playEventApi.record(track.id).catch(() => {})
  }

  function togglePlay() {
    const a = getAudio()
    if (playing.value) {
      a.pause()
    } else {
      a.play()
    }
  }

  function seek(time: number) {
    const a = getAudio()
    a.currentTime = time
  }

  function setVolume(v: number) {
    volume.value = v
    const a = getAudio()
    a.volume = v
  }

  function next() {
    if (queue.value.length === 0) return

    if (playMode.value === 'repeat-one') {
      const a = getAudio()
      a.currentTime = 0
      a.play()
      return
    }

    if (playMode.value === 'shuffle') {
      let idx = Math.floor(Math.random() * queue.value.length)
      if (queue.value.length > 1) {
        while (idx === currentIndex.value) {
          idx = Math.floor(Math.random() * queue.value.length)
        }
      }
      play(queue.value[idx])
      return
    }

    let nextIdx = currentIndex.value + 1
    if (nextIdx >= queue.value.length) {
      if (playMode.value === 'repeat-all') {
        nextIdx = 0
      } else {
        playing.value = false
        return
      }
    }
    play(queue.value[nextIdx])
  }

  function previous() {
    const a = getAudio()
    if (a.currentTime > 3) {
      a.currentTime = 0
      return
    }

    let prevIdx = currentIndex.value - 1
    if (prevIdx < 0) {
      prevIdx = playMode.value === 'repeat-all' ? queue.value.length - 1 : 0
    }
    play(queue.value[prevIdx])
  }

  function cyclePlayMode() {
    const modes: PlayMode[] = ['sequential', 'repeat-all', 'repeat-one', 'shuffle']
    const idx = modes.indexOf(playMode.value)
    playMode.value = modes[(idx + 1) % modes.length]
  }

  function toggleLyrics() {
    lyricsVisible.value = !lyricsVisible.value
  }

  return {
    queue, currentIndex, playing, currentTime, duration, volume,
    playMode, lyricsVisible, currentTrack,
    playModeIcon, playModeLabel,
    play, togglePlay, seek, setVolume, next, previous,
    cyclePlayMode, toggleLyrics,
  }
})
