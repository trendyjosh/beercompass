export {};

// Fix issue with Safari non-standard webkit addition to fix
// TS error throw.
declare global {
  interface DeviceOrientationEvent {
    readonly webkitCompassHeading?: number;
    readonly webkitCompassAccuracy?: number;
  }

  interface DeviceOrientationEventStatic {
    requestPermission?: () => Promise<"granted" | "denied">;
  }

  interface Window {
    DeviceOrientationEvent: typeof DeviceOrientationEvent & DeviceOrientationEventStatic;
  }
}
