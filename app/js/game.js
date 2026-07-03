// js/main.js
document.getElementById("start-game-btn").addEventListener("click", () => {
    renderBlocks()
    initWebSocket()
})

document.getElementById("stop-game-btn").addEventListener("click", () => {
    clearCanvas()
    closeWebSocket()
})

const canvas = document.getElementById("game-canvas");
const ctx = canvas.getContext("2d");

function renderBlocks() {
    clearCanvas();

    // background
    ctx.fillStyle = "black";
    ctx.fillRect(0, 0, 800, 600);

    // placeholder player rectangle
    ctx.fillStyle = "blue";
    ctx.fillRect(0, 0, 50, 50);

    // placeholder monster rectangle
    ctx.fillStyle = "red";
    ctx.fillRect(750, 350, 50, 50);
};

function clearCanvas() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
}