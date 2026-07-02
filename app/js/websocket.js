async function initWebSocket() {
    document.getElementById("start-game-button").style.display = "none";
    const socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onopen = () => {
        socket.send("hello");
    };
    socket.onmessage = (event) => console.log(event.data);
    socket.onerror = (error) => console.log("WebSocket error:", error)
    socket.onclose = () => console.log("WebSocket closed")
};