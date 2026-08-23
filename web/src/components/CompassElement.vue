<script setup lang="ts">
import { ref } from "vue";
import IconCompass from "./icons/IconCompass.vue";

const bearing = ref(null);
const heading = ref(null);
const currentOrientation = ref(14);

function rotate() {
  if (bearing.value === null || heading.value === null) return;

  // Subtract device heading from pub bearing so the needle
  // stays pointing at the pub as the device rotates
  const direction = (bearing.value - heading.value + 360) % 360;

  // Find the shortest rotation direction
  const delta = direction - (currentOrientation.value % 360);

  let newOrientation = currentOrientation.value;
  if (delta > 180) {
    newOrientation -= 360;
  } else if (delta < -180) {
    newOrientation += 360;
  }

  newOrientation += delta;

  currentOrientation.value = newOrientation;
}
</script>

<template>
  <IconCompass v-model="currentOrientation" />
</template>
