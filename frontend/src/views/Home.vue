<template>
  <v-container class="py-8">
    <div class="text-center mb-8">
      <div class="text-h3 font-weight-bold mb-2">
        <span class="text-primary">Pineapple</span> Music
      </div>
      <p class="text-body-1 text-medium-emphasis">{{ t('nav.home') }}</p>
    </div>

    <v-row justify="center">
      <v-col cols="12" sm="6" md="4" lg="3">
        <v-card to="/tracks" hover class="text-center pa-6">
          <v-icon size="48" color="primary" class="mb-3">mdi-music-note-outline</v-icon>
          <div class="text-h6">{{ t('nav.tracks') }}</div>
          <div class="text-body-2 text-medium-emphasis">Browse your library</div>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="4" lg="3">
        <v-card to="/playlists" hover class="text-center pa-6">
          <v-icon size="48" color="secondary" class="mb-3">mdi-playlist-music-outline</v-icon>
          <div class="text-h6">{{ t('nav.playlists') }}</div>
          <div class="text-body-2 text-medium-emphasis">Curated collections</div>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="4" lg="3">
        <v-card to="/search" hover class="text-center pa-6">
          <v-icon size="48" color="info" class="mb-3">mdi-magnify</v-icon>
          <div class="text-h6">{{ t('nav.search') }}</div>
          <div class="text-body-2 text-medium-emphasis">Find anything</div>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
