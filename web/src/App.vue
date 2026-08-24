<script setup lang="ts">
import { ref } from "vue";
import CompassElement from "./components/CompassElement.vue";

const locationMessage = ref("...");
const distanceMessage = ref(null);
const directionMessage = ref(null);
</script>

<template>
  <main>
    <div class="flex flex-col justify-center items-center h-screen w-screen top-0 p-6">
      <h1 class="mb-6 text-2xl">{{ locationMessage }}</h1>
      <div class="max-w-full aspect-square flex flex-col justify-center text-center gap-2 relative">
        <div class="absolute w-full h-full flex flex-col justify-center text-center text-3xl">
          <p v-show="distanceMessage">{{ distanceMessage }}</p>
          <p v-show="directionMessage">{{ directionMessage }}</p>
        </div>
        <CompassElement :locationMessage :distanceMessage :directionMessage />
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
