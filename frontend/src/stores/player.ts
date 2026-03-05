import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
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

  // Next track for preloading
  const nextTrack = computed(() => {
    if (queue.value.length === 0 || playMode.value === 'shuffle') return null
    if (playMode.value === 'repeat-one') return currentTrack.value
    const nextIdx = currentIndex.value + 1
    if (nextIdx < queue.value.length) return queue.value[nextIdx]
    if (playMode.value === 'repeat-all') return queue.value[0]
    return null
  })

  const playModeIcon = computed(() => {
    switch (playMode.value) {
      case 'sequential': return 'ArrowRight'
      case 'repeat-all': return 'Repeat'
      case 'repeat-one': return 'Repeat1'
      case 'shuffle': return 'Shuffle'
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
  // Preload audio element for next track
  let preloadAudio: HTMLAudioElement | null = null
  let preloadedTrackId: number | null = null

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

  function preloadNext() {
    const next = nextTrack.value
    if (!next || next.id === preloadedTrackId) return
    if (!preloadAudio) {
      preloadAudio = new Audio()
      preloadAudio.preload = 'auto'
    }
    preloadAudio.src = trackApi.streamUrl(next.id)
    preloadedTrackId = next.id
  }

  // Watch for track changes and preload next
  watch(currentIndex, () => {
    preloadNext()
  })

  function updateMediaSession(track: Track) {
    if (!('mediaSession' in navigator)) return
    navigator.mediaSession.metadata = new MediaMetadata({
      title: track.title,
      artist: track.artist,
      album: track.album,
      artwork: track.has_cover
        ? [{ src: trackApi.coverUrl(track.id), sizes: '512x512', type: 'image/jpeg' }]
        : [],
    })
    navigator.mediaSession.setActionHandler('play', () => togglePlay())
    navigator.mediaSession.setActionHandler('pause', () => togglePlay())
    navigator.mediaSession.setActionHandler('previoustrack', () => previous())
    navigator.mediaSession.setActionHandler('nexttrack', () => next())
    navigator.mediaSession.setActionHandler('seekto', (details) => {
      if (details.seekTime != null) seek(details.seekTime)
    })
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

    // If the next track was preloaded and matches, swap the audio elements
    if (preloadAudio && preloadedTrackId === track.id) {
      // Swap preloaded audio in as the main player
      const oldAudio = audio!
      // Copy event listeners by replacing
      preloadAudio.volume = volume.value
      preloadAudio.addEventListener('timeupdate', () => {
        currentTime.value = preloadAudio!.currentTime
      })
      preloadAudio.addEventListener('loadedmetadata', () => {
        duration.value = preloadAudio!.duration
      })
      preloadAudio.addEventListener('ended', () => {
        next()
      })
      preloadAudio.addEventListener('play', () => {
        playing.value = true
      })
      preloadAudio.addEventListener('pause', () => {
        playing.value = false
      })
      audio = preloadAudio
      preloadAudio = null
      preloadedTrackId = null
      audio.play()
      // Clean up old audio
      oldAudio.pause()
      oldAudio.removeAttribute('src')
      oldAudio.load()
    } else {
      a.src = trackApi.streamUrl(track.id)
      a.volume = volume.value
      a.play()
    }

    updateMediaSession(track)
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
    playMode, lyricsVisible, currentTrack, nextTrack,
    playModeIcon, playModeLabel,
    play, togglePlay, seek, setVolume, next, previous,
    cyclePlayMode, toggleLyrics,
  }
})
