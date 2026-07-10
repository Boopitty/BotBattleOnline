// js/websocket.js
let socket = null;

async function initWebSocket() {
    
    socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onopen = () => {
        document.getElementById("start-game-btn").style.display = "none";
        document.getElementById("stop-game-btn").style.display = "block";
        document.getElementById("attack-btn").style.display = "block";
    };
    socket.onmessage = (event) => {
        try {
            const resp = JSON.parse(event.data) 
            console.log(resp.message)
        } catch (error) {
            console.error(`Error when parsing: ${error}`)
        }
    };
    socket.onerror = (error) => console.log("WebSocket error:", error)
    socket.onclose = () => {
        document.getElementById("start-game-btn").style.display = "block";
        document.getElementById("stop-game-btn").style.display = "none";
        document.getElementById("attack-btn").style.display = "none";
        console.log("WebSocket closed");
    }
};

async function closeWebSocket() {
    if (socket) {
        socket.close();
        socket = null;
    }
}

// Send a request through the websocket. req is a struct in JSON format.
async function sendRequest(req) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(req);
    } else {
        console.error(`WebSocket ready state: ${socket.readyState}. Unable to send command: ${command}` );
    }
}