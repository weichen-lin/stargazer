function getRandomHexColor() {
  return `#${Math.floor(Math.random() * 0xffffff)
    .toString(16)
    .padStart(6, '0')}`
}

export function generateRandomColors(count = 5) {
  return Array.from({ length: count }, getRandomHexColor)
}
