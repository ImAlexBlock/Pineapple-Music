<template>
  <v-app>
    <!-- Top app bar -->
    <v-app-bar v-if="auth.bootstrapped" flat border density="compact" color="surface">
      <v-app-bar-nav-icon @click="drawer = !drawer" />
      <v-app-bar-title class="font-weight-bold">
        <span class="text-primary">Pineapple</span> Music
      </v-app-bar-title>
      <v-spacer />
      <v-btn icon variant="text" size="small" @click="themeStore.toggle()">
        <v-icon>{{ themeStore.isDark ? 'mdi-white-balance-sunny' : 'mdi-moon-waning-crescent' }}</v-icon>
        <v-tooltip activator="parent" location="bottom">Toggle theme</v-tooltip>
      </v-btn>
      <v-menu>
        <template #activator="{ props }">
          <v-btn icon variant="text" size="small" v-bind="props">
            <v-icon>mdi-translate</v-icon>
          </v-btn>
        </template>
        <v-list density="compact" rounded="lg">
          <v-list-item @click="locale = 'en-US'" :active="locale === 'en-US'">English</v-list-item>
          <v-list-item @click="locale = 'zh-CN'" :active="locale === 'zh-CN'">中文</v-list-item>
        </v-list>
      </v-menu>
      <v-btn v-if="auth.role" icon variant="text" size="small" @click="logout">
        <v-icon>mdi-logout</v-icon>
        <v-tooltip activator="parent" location="bottom">{{ t('nav.logout') }}</v-tooltip>
      </v-btn>
    </v-app-bar>

    <!-- Navigation drawer -->
    <v-navigation-drawer v-if="auth.bootstrapped" v-model="drawer" temporary width="280">
      <div class="pa-4 pb-2">
        <div class="text-h6 font-weight-bold"><span class="text-primary">Pineapple</span> Music</div>
        <div class="text-caption text-medium-emphasis">v0.1.0</div>
      </div>
      <v-divider />
      <v-list nav density="comfortable" class="px-2">
        <v-list-item prepend-icon="mdi-home-outline" :title="t('nav.home')" to="/" rounded="lg" />
        <v-list-item prepend-icon="mdi-music-note-outline" :title="t('nav.tracks')" to="/tracks" rounded="lg" />
        <v-list-item prepend-icon="mdi-playlist-music-outline" :title="t('nav.playlists')" to="/playlists" rounded="lg" />
        <v-list-item prepend-icon="mdi-magnify" :title="t('nav.search')" to="/search" rounded="lg" />
      </v-list>
      <template v-if="auth.role === 'admin'">
        <v-divider class="my-1" />
        <v-list nav density="comfortable" class="px-2">
          <v-list-subheader>{{ t('nav.admin') }}</v-list-subheader>
          <v-list-item prepend-icon="mdi-view-dashboard-outline" :title="t('admin.dashboard')" to="/admin" rounded="lg" />
          <v-list-item prepend-icon="mdi-upload-outline" :title="t('admin.upload')" to="/admin/upload" rounded="lg" />
          <v-list-item prepend-icon="mdi-magnify-scan" :title="t('admin.scan')" to="/admin/scan" rounded="lg" />
          <v-list-item prepend-icon="mdi-cog-outline" :title="t('admin.settings')" to="/admin/settings" rounded="lg" />
          <v-list-item prepend-icon="mdi-clipboard-text-clock-outline" :title="t('admin.auditLog')" to="/admin/audit" rounded="lg" />
        </v-list>
      </template>
    </v-navigation-drawer>

    <!-- Main content -->
    <v-main :class="{ 'pb-player': auth.bootstrapped && player.currentTrack }">
      <router-view />
    </v-main>

    <!-- Player bar -->
    <PlayerBar v-if="auth.bootstrapped" />
    <LyricsPanel />
  </v-app>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from './stores/auth'
import { useThemeStore } from './stores/theme'
import { usePlayerStore } from './stores/player'
import PlayerBar from './components/player/PlayerBar.vue'
import LyricsPanel from './components/player/LyricsPanel.vue'

const { t, locale } = useI18n()
const auth = useAuthStore()
const themeStore = useThemeStore()
const player = usePlayerStore()
const drawer = ref(false)

onMounted(() => {
  themeStore.init()
})

async function logout() {
  await auth.logout()
  window.location.href = '/login'
}
</script>

<style>
.pb-player {
  padding-bottom: 80px !important;
}
</style>
