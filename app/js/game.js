const canvas = document.getElementById("game-canvas")
const ctx = canvas.getContext("2d")

// placeholder player rectangle
ctx.fillStyle = "blue"
ctx.fillRect(100, 400, 50, 50)

// placeholder monster rectangle
ctx.fillStyle = "red"
ctx.fillRect(650, 400, 50, 50)