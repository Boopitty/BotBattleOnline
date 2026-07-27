// js/websocket.js
let socket = null;

async function initWebSocket() {
    if (socket) {
        console.log("Websocket already exists.");
        return
    }

    socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onopen = () => {
        console.log("WebSocket Opened");
    };
    socket.onmessage = (event) => {
        try {
            const resp = JSON.parse(event.data);
            console.log(resp.message);
        } catch (error) {
            console.error(`Error when parsing: ${error}`);
        }
    };
    socket.onerror = (error) => console.log("WebSocket error:", error);
    socket.onclose = () => {
        console.log("WebSocket Closed");
    }
};

async function closeWebSocket() {
    if (socket) {
        socket.close();
        socket = null;
    } else {
        console.log("Websocket not found.");
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
