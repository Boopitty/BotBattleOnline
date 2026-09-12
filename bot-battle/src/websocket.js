// js/websocket.js
export { initWebSocket, closeWebSocket, makeRequest, sendRequest };
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
        socket = null;
        console.log("WebSocket Closed");
    }
};

async function closeWebSocket() {
    if (socket) {
        socket.close();
    } else {
        console.warn("Websocket not found.");
    }
}

/** 
 * @Param {string} command - The command to send
 * @Param {number} input - The input for the command
 * @returns {string} - The request string in JSON format
 */
function makeRequest (command, input = 0) {
    return JSON.stringify({
        User: {
            username: localStorage.getItem("username"),
        },
        command: command,
        input: input
    });
}

/**
 * @param {string} req - The request string in JSON format
 * @returns {Promise} - The response from the server
 */
async function sendRequest(req) {
    if (socket) {
        if (socket.readyState === WebSocket.OPEN) {
            socket.send(req);
        } else {
            console.warn(`WebSocket ready state: ${socket.readyState}. Unable to send command: ${req}` );
        }
    } else {
        console.warn(`WebSocket not found. Unable to send request: ${req}`);
    }
}
