import api from './client'

export const authApi = {
  getStatus: () => api.get('/setup/status'),
  bootstrap: () => api.post('/setup/bootstrap'),
  login: (key: string) => api.post('/auth/login', { key }),
  logout: () => api.post('/auth/logout'),
  me: () => api.get('/auth/me'),
}

export const trackApi = {
  list: (params?: Record<string, unknown>) => api.get('/tracks', { params }),
  get: (id: number) => api.get(`/tracks/${id}`),
  streamUrl: (id: number) => `/api/v1/tracks/${id}/stream`,
  coverUrl: (id: number) => `/api/v1/tracks/${id}/cover`,
  getLyrics: (id: number) => api.get(`/tracks/${id}/lyrics`),
  getArtists: () => api.get('/artists'),
  getAlbums: (params?: Record<string, unknown>) => api.get('/albums', { params }),
}

export const playlistApi = {
  list: () => api.get('/playlists'),
  get: (id: number) => api.get(`/playlists/${id}`),
  create: (name: string) => api.post('/playlists', { name }),
  delete: (id: number) => api.delete(`/playlists/${id}`),
  addTrack: (id: number, trackId: number) => api.post(`/playlists/${id}/tracks`, { track_id: trackId }),
  removeTrack: (id: number, trackId: number) => api.delete(`/playlists/${id}/tracks/${trackId}`),
  reorder: (id: number, trackIds: number[]) => api.put(`/playlists/${id}/reorder`, { track_ids: trackIds }),
}

export const playEventApi = {
  record: (trackId: number) => api.post('/play-events', { track_id: trackId }),
}

export const adminApi = {
  dashboard: () => api.get('/admin/dashboard'),
  getSettings: () => api.get('/admin/settings'),
  updateSettings: (settings: Record<string, string>) => api.put('/admin/settings', settings),
  rotateAdminKey: () => api.post('/admin/rotate-admin-key'),
  rotateGuestKey: () => api.post('/admin/rotate-guest-key'),
  getAuditLogs: (params?: Record<string, unknown>) => api.get('/admin/audit-logs', { params }),
  upload: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
  startScan: () => api.post('/scan'),
  scanStatus: () => api.get('/scan/status'),
}
