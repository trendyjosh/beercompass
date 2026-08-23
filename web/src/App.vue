<script setup lang="ts">
import IconCompass from "./components/icons/IconCompass.vue";
</script>

<template>
  <main>
    <div class="flex flex-col justify-center items-center h-screen w-screen top-0 p-6">
      <h1 id="location" class="mb-6 text-2xl">...</h1>
      <div class="max-w-full aspect-square flex flex-col justify-center text-center gap-2 relative">
        <div class="absolute w-full h-full flex flex-col justify-center text-center text-3xl">
          <p id="distance">...</p>
          <p id="direction">...</p>
        </div>
        <IconCompass />
      </div>
      <p class="mt-6">
        Map data ©
        <a href="https://www.openstreetmap.org/copyright" target="_blank" rel="noopener">
          OpenStreetMap
        </a>
      </p>
    </div>
  </main>
</template>
