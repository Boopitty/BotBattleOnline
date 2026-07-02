// js/main.js
document.getElementById("start-game-btn").addEventListener("click", () => {
    renderBlocks()
    initWebSocket()
})

function renderBlocks() {
    const canvas = document.getElementById("game-canvas");
    const ctx = canvas.getContext("2d");

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