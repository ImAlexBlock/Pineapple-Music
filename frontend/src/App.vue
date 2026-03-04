<template>
  <div class="min-h-screen bg-background">
    <!-- Top header bar -->
    <header v-if="auth.bootstrapped" class="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div class="flex h-12 items-center px-4">
        <Button variant="ghost" size="icon" class="mr-2 shrink-0" @click="drawerOpen = true">
          <Menu class="h-5 w-5" />
        </Button>
        <span class="font-bold text-sm">
          <span class="text-primary">Pineapple</span> Music
        </span>
        <div class="flex-1" />
        <Button variant="ghost" size="icon" class="h-8 w-8" @click="themeStore.toggle()">
          <Sun v-if="themeStore.isDark" class="h-4 w-4" />
          <Moon v-else class="h-4 w-4" />
        </Button>
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button variant="ghost" size="icon" class="h-8 w-8">
              <Languages class="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem @click="locale = 'en-US'">
              <span :class="{ 'font-bold': locale === 'en-US' }">English</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="locale = 'zh-CN'">
              <span :class="{ 'font-bold': locale === 'zh-CN' }">中文</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <Button v-if="auth.role" variant="ghost" size="icon" class="h-8 w-8" @click="logout">
          <LogOut class="h-4 w-4" />
        </Button>
      </div>
    </header>

    <!-- Navigation drawer (Sheet) -->
    <Sheet v-model:open="drawerOpen">
      <SheetContent side="left" class="w-[280px] p-0">
        <SheetHeader class="p-4 pb-2">
          <SheetTitle class="text-lg font-bold"><span class="text-primary">Pineapple</span> Music</SheetTitle>
          <SheetDescription class="text-xs text-muted-foreground">v0.1.0</SheetDescription>
        </SheetHeader>
        <Separator />
        <nav class="flex flex-col gap-1 p-2">
          <router-link v-for="(item, i) in navItems" :key="item.to" :to="item.to" custom v-slot="{ isActive, navigate }">
            <button
              class="nav-item flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm transition-all duration-200 hover:bg-accent hover:translate-x-1"
              :class="{ 'bg-accent text-accent-foreground font-medium': isActive }"
              :style="{ animationDelay: `${i * 50}ms` }"
              @click="navigate(); drawerOpen = false"
            >
              <component :is="item.icon" class="h-4 w-4 transition-transform duration-200" :class="{ 'scale-110': isActive }" />
              {{ item.label }}
            </button>
          </router-link>
        </nav>
        <template v-if="auth.role === 'admin'">
          <Separator class="my-1" />
          <div class="px-4 py-1 text-xs font-medium text-muted-foreground">{{ t('nav.admin') }}</div>
          <nav class="flex flex-col gap-1 p-2 pt-0">
            <router-link v-for="(item, i) in adminNavItems" :key="item.to" :to="item.to" custom v-slot="{ isActive, navigate }">
              <button
                class="nav-item flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm transition-all duration-200 hover:bg-accent hover:translate-x-1"
                :class="{ 'bg-accent text-accent-foreground font-medium': isActive }"
                :style="{ animationDelay: `${(i + navItems.length) * 50}ms` }"
                @click="navigate(); drawerOpen = false"
              >
                <component :is="item.icon" class="h-4 w-4 transition-transform duration-200" :class="{ 'scale-110': isActive }" />
                {{ item.label }}
              </button>
            </router-link>
          </nav>
        </template>
      </SheetContent>
    </Sheet>

    <!-- Main content -->
    <main class="mx-auto w-full max-w-5xl px-4 py-6" :class="{ 'pb-20': auth.bootstrapped && player.currentTrack }">
      <router-view />
    </main>

    <!-- Player bar -->
    <PlayerBar v-if="auth.bootstrapped" />
    <LyricsPanel />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from './stores/auth'
import { useThemeStore } from './stores/theme'
import { usePlayerStore } from './stores/player'
import PlayerBar from './components/player/PlayerBar.vue'
import LyricsPanel from './components/player/LyricsPanel.vue'
import { Button } from '@/components/ui/button'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from '@/components/ui/sheet'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Separator } from '@/components/ui/separator'
import {
  Menu, Sun, Moon, Languages, LogOut,
  Home, Music, ListMusic, Search,
  LayoutDashboard, Upload, ScanSearch, Settings, ClipboardList,
} from 'lucide-vue-next'

const { t, locale } = useI18n()
const auth = useAuthStore()
const themeStore = useThemeStore()
const player = usePlayerStore()
const drawerOpen = ref(false)

onMounted(() => {
  themeStore.init()
})

const navItems = computed(() => [
  { to: '/', icon: Home, label: t('nav.home') },
  { to: '/tracks', icon: Music, label: t('nav.tracks') },
  { to: '/playlists', icon: ListMusic, label: t('nav.playlists') },
  { to: '/search', icon: Search, label: t('nav.search') },
])

const adminNavItems = computed(() => [
  { to: '/admin', icon: LayoutDashboard, label: t('admin.dashboard') },
  { to: '/admin/upload', icon: Upload, label: t('admin.upload') },
  { to: '/admin/scan', icon: ScanSearch, label: t('admin.scan') },
  { to: '/admin/settings', icon: Settings, label: t('admin.settings') },
  { to: '/admin/audit', icon: ClipboardList, label: t('admin.auditLog') },
])

async function logout() {
  await auth.logout()
  window.location.href = '/login'
}
</script>
