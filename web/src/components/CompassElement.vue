<script setup lang="ts">
import { onMounted, ref, type Ref } from "vue";
import IconCompass from "./icons/IconCompass.vue";
import { distanceBetween } from "@/utils/distance.ts";
import axios from "axios";

const locationMessage = defineModel("locationMessage");
const distanceMessage = defineModel("distanceMessage");
const directionMessage = defineModel("directionMessage");

// Direction of compass
const bearing: Ref<number | null> = ref(null);
// Direction of nearest pub
const heading: Ref<number | null> = ref(null);
// Store time bearing was last updated
let bearingUpdated: number | null = null;
// Delay before updating bearing
const bearingDelay: number = 700; // 0.7s
// Number of readings taken
let orientationReadings: number = 0;
// Minimum number of readings before first rotating
const stabiliseAfter: number = 2;
const currentOrientation: Ref<number> = ref(14);
// Movement threshold
const threshold: number = 0.05; // 50m
// Last stored latitude
let lastLat: number | null = null;
// Last stored longitude
let lastLng: number | null = null;
// Fetching state to avoid duplicate requests
let isFetching: boolean = false;

/**
 * Retrieve nearest pub from endpoint and update interface.
 */
async function getNearestPub(lat: number, lng: number) {
  if (isFetching) return;
  isFetching = true;

  axios
    .post("/nearest", {
      lat: lat,
      lng: lng,
    })
    .then((response) => {
      const pub = response.data;

      // Print output
      locationMessage.value = pub.name;
      distanceMessage.value = pub.distanceKm + "km";
      directionMessage.value = `${pub.bearing}° ${pub.direction}`;

      // Rotate compass to bearing
      setBearing(pub.bearing);
    })
    .catch((error: Error) => {
      locationMessage.value = "Server error";
      distanceMessage.value = "No pubs found";
      directionMessage.value = null;
    })
    .finally(() => {
      isFetching = false;
    });
}

/**
 * Get the device's geolocation.
 */
async function getLocation() {
  // Check navigator enabled
  if (!navigator.geolocation) {
    locationMessage.value = "Geolocation not supported";
    distanceMessage.value = null;
    directionMessage.value = null;
    return;
  }

  const geoOptions = {
    enableHighAccuracy: false,
    maximumAge: 300000,
  };

  navigator.geolocation.watchPosition(
    async (position: GeolocationPosition) => {
      return positionFound(position);
    },
    (error: GeolocationPositionError) => {
      locationMessage.value = "Location not found";
      distanceMessage.value = "Please enable location services";
      directionMessage.value = null;
      return;
    },
    geoOptions,
  );
}

/**
 * Trigger nearest pub retrieval based on device's location.
 */
async function positionFound(position: GeolocationPosition) {
  const lat: number = position.coords.latitude;
  const lng: number = position.coords.longitude;

  // Only refetch if the device has moved more than threshold
  if (
    lastLat !== null &&
    lastLng !== null &&
    distanceBetween(lastLat, lastLng, lat, lng) < threshold
  ) {
    return;
  }

  // Cache last location
  lastLat = lat;
  lastLng = lng;

  // Get nearest pub
  await getNearestPub(lat, lng);
}

/**
 * Initialise the listener for orientation change.
 */
function start() {
  if (typeof DeviceOrientationEvent === "undefined") {
    console.warn("Device orientation not supported");
    return;
  }

  if (typeof DeviceOrientationEvent.requestPermission === "function") {
    // iOS 13+ requires explicit permission
    requestIOSPermission();
  } else {
    // Android — try absolute first, fall back to relative
    attachListener();
  }
}

/**
 * Get permissions from the user to enable orientation.
 */
function requestIOSPermission() {
  // iOS requires permission to be requested from a user gesture
  // so we create a button if needed
  const backdrop = document.createElement("div");
  backdrop.classList.add("compass-permission-backdrop");
  const button = document.createElement("button");
  button.innerText = "Enable Compass";
  button.classList.add("compass-permission-btn");

  button.addEventListener("click", async () => {
    try {
      const permission = await DeviceOrientationEvent.requestPermission();
      if (permission === "granted") {
        button.remove();
        backdrop.remove();
        attachListener();
      } else {
        button.innerText = "Compass permission denied";
      }
    } catch (error) {
      console.error("iOS orientation permission error:", error);
    }
  });

  document.body.appendChild(backdrop);
  document.body.appendChild(button);
}

/**
 * Create the listener for orientation change.
 */
function attachListener() {
  let usingAbsolute: boolean = false;

  const absoluteHandler = (event: DeviceOrientationEvent) => {
    if (event.absolute === true && event.alpha !== null) {
      usingAbsolute = true;
      handleOrientation(event);
    }
  };

  const relativeHandler = (event: DeviceOrientationEvent) => {
    // Only use relative if absolute never fired
    if (!usingAbsolute) {
      handleOrientation(event);
    }
  };

  // Try absolute first
  window.addEventListener("deviceorientationabsolute", absoluteHandler, true);

  // Attach relative as fallback — relativeHandler checks if absolute is working
  window.addEventListener("deviceorientation", relativeHandler, true);

  // After 1 second, if absolute never fired, log that we are using fallback
  setTimeout(() => {
    if (!usingAbsolute) {
      console.warn(
        "deviceorientationabsolute not supported, using relative fallback — compass may not point true north",
      );
    }
  }, 1000);
}

/**
 * Get the orientation and update compass.
 */
function handleOrientation(event: DeviceOrientationEvent) {
  // Throttle bearing update
  const now: number = Date.now();
  if (bearingUpdated && now - bearingUpdated < bearingDelay) return;
  bearingUpdated = now;

  let heading = null;

  if (event.webkitCompassHeading !== undefined && event.webkitCompassHeading !== null) {
    // iOS — webkitCompassHeading is already 0-360 relative to true north
    heading = event.webkitCompassHeading;
  } else if (event.absolute && event.alpha !== null) {
    // Android absolute — convert alpha to compass heading
    heading = 360 - event.alpha;
  } else if (event.alpha !== null) {
    // Relative fallback — not true north, just device rotation
    heading = 360 - event.alpha;
  }

  if (heading === null) return;

  orientationReadings++;
  heading = heading;

  // Only rotate once stabalised
  if (orientationReadings < stabiliseAfter) return;

  rotate();
}

/**
 * Set the bearing value.
 */
function setBearing(newBearing: number) {
  bearing.value = newBearing;
  rotate();
}

/**
 * Apply rotation to the compass element.
 */
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

onMounted(start);
</script>

<template>
  <IconCompass v-model="currentOrientation" />
</template>
