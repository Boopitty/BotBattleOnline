let socket = null;

async function initWebSocket() {
    
    socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onopen = () => {
        socket.send("hello");
        document.getElementById("start-game-btn").style.display = "none";
        document.getElementById("stop-game-btn").style.display = "block";
    };
    socket.onmessage = (event) => console.log(event.data);
    socket.onerror = (error) => console.log("WebSocket error:", error)
    socket.onclose = () => {
        console.log("WebSocket closed");
        document.getElementById("start-game-btn").style.display = "block";
        document.getElementById("stop-game-btn").style.display = "none";
    }
};

async function closeWebSocket() {
    if (socket) {
        socket.close();
        socket = null;
    }
}

async function sendCommand(command) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({command}));
    } else {
        console.error("WebSocket is not open. Unable to send command:", command);
    }
}