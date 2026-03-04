export interface Track {
  id: number
  title: string
  artist: string
  album: string
  album_artist: string
  genre: string
  year: number
  track_number: number
  disc_number: number
  duration: number
  format: string
  size: number
  bitrate: number
  sample_rate: number
  has_cover: boolean
  has_lyrics: boolean
  created_at: string
  updated_at: string
}

export interface TrackLyric {
  id: number
  track_id: number
  type: 'plain' | 'synced'
  content: string
}

export interface Playlist {
  id: number
  name: string
  created_at: string
  updated_at: string
  tracks?: PlaylistTrack[]
}

export interface PlaylistTrack {
  id: number
  playlist_id: number
  track_id: number
  position: number
  track: Track
}

export interface ScanJob {
  id: number
  status: string
  total: number
  scanned: number
  added: number
  updated: number
  errors: number
  error_log?: string
  started_at: string
  finished_at?: string
}

export interface AuditLog {
  id: number
  action: string
  role: string
  ip: string
  detail: string
  created_at: string
}

export interface DashboardStats {
  tracks: number
  total_size: number
  plays: number
  playlists: number
}
