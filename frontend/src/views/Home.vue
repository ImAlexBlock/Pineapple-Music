<template>
  <div>
    <div class="text-center mb-8">
      <h1 class="text-3xl font-bold mb-2">
        <span class="text-primary">Pineapple</span> Music
      </h1>
      <p class="text-muted-foreground">{{ t('nav.home') }}</p>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 max-w-2xl mx-auto">
      <router-link to="/tracks" class="block">
        <Card class="text-center p-6 hover:shadow-md transition-shadow cursor-pointer h-full">
          <Music class="h-12 w-12 mx-auto mb-3 text-primary" />
          <div class="font-semibold">{{ t('nav.tracks') }}</div>
          <p class="text-sm text-muted-foreground mt-1">Browse your library</p>
        </Card>
      </router-link>
      <router-link to="/playlists" class="block">
        <Card class="text-center p-6 hover:shadow-md transition-shadow cursor-pointer h-full">
          <ListMusic class="h-12 w-12 mx-auto mb-3 text-secondary" />
          <div class="font-semibold">{{ t('nav.playlists') }}</div>
          <p class="text-sm text-muted-foreground mt-1">Curated collections</p>
        </Card>
      </router-link>
      <router-link to="/search" class="block">
        <Card class="text-center p-6 hover:shadow-md transition-shadow cursor-pointer h-full">
          <SearchIcon class="h-12 w-12 mx-auto mb-3 text-info" />
          <div class="font-semibold">{{ t('nav.search') }}</div>
          <p class="text-sm text-muted-foreground mt-1">Find anything</p>
        </Card>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Card } from '@/components/ui/card'
import { Music, ListMusic, Search as SearchIcon } from 'lucide-vue-next'
const { t } = useI18n()
</script>
