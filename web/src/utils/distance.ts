/**
 * Get the distance between 2 points in Km.
 * Calculation is the Haversine distance formula.
 */
export function distanceBetween(oldLat: number, oldLng: number, newLat: number, newLng: number) {
  const earthRadius = 6371;

  const oldLatRad = deg2rad(oldLat);
  const oldLngRad = deg2rad(oldLng);
  const newLatRad = deg2rad(newLat);
  const newLngRad = deg2rad(newLng);

  const latDelta = newLatRad - oldLatRad;
  const lngDelta = newLngRad - oldLngRad;

  const haversine =
    Math.pow(Math.sin(latDelta / 2), 2) +
    Math.cos(oldLatRad) * Math.cos(newLatRad) * Math.pow(Math.sin(lngDelta / 2), 2);

  const angle = 2 * Math.atan2(Math.sqrt(haversine), Math.sqrt(1 - haversine));

  return angle * earthRadius;
}

/**
 * Converts the number in degrees to the radian equivalent.
 */
function deg2rad(value: number) {
  return (value * Math.PI) / 180;
}
