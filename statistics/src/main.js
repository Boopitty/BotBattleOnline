document.addEventListener("DOMContentLoaded", async () => {
    const users = await getNumUsers();
    const players = await getCurrentPlayers();
    const bots = await getNumBots();
    const games = await getNumGames();

    document.getElementById("numUsers").textContent = users;
    document.getElementById("numPlayers").textContent = players;
    document.getElementById("numBots").textContent = bots;
    document.getElementById("numGames").textContent = games;
});

document.getElementById("reset").addEventListener("click", async () => {
    console.log("Resetting database...");
    try {
        const resp = await fetch("/admin/resetUsers", {
            method: "DELETE"
        });
        const data = resp.json();
        if (!resp.ok) {
            throw Error(`Reset Failed: ${data.error}`)
        }
        console.log("Database reset successful");
    } catch (error) {
        alert(`Error: ${error}`);
        console.error("Database reset failed");
    }
});

async function getNumUsers() {
    try {
        const resp = await fetch("/admin/getNumUsers");
        const data = await resp.json();
        if (!resp.ok) {
            console.error(`Failed to get number of users: ${data.error}`)
            return "undefined";
        }
        if (data.userCount === undefined) {
            console.error("Number of users data is undefined");
            return "undefined";
        }
        return data.userCount;
    } catch (error) {
        alert(`Error: ${error}`);
        console.error("Failed to get number of users");
    }
}

async function getNumBots() {
    try {
        const resp = await fetch("/admin/getNumBots");
        const data = await resp.json();
        if (!resp.ok) {
            console.error(`Failed to get number of bots: ${data.error}`)
            return "undefined";
        }
        if (data.botCount === undefined) {
            console.error("Number of bots data is undefined");
            return "undefined";
        }
        return data.botCount;
    } catch (error) {
        alert(`Error: ${error}`);
        console.error("Failed to get number of bots");
    }
}

async function getNumGames() {
    try {
        const resp = await fetch("/admin/getNumGames");
        const data = await resp.json();
        if (!resp.ok) {
            console.error(`Failed to get number of games: ${data.error}`)
            return "undefined";
        }
        if (data.gameCount === undefined) {
            console.error("Number of games data is undefined");
            return "undefined";
        }
        return data.gameCount;
    } catch (error) {
        alert(`Error: ${error}`);
        console.error("Failed to get number of games");
    }
}

async function getCurrentPlayers() {
    try {
        const resp = await fetch("/admin/getCurrentPlayers");
        const data = await resp.json();
        if (!resp.ok) {
            console.error(`Failed to get number of current players: ${data.error}`)
            return "undefined";
        }
        if (data.playerCount === undefined) {
            console.error("Number of current players data is undefined");
            return "undefined";
        }
        return data.playerCount;
    } catch (error) {
        alert(`Error: ${error}`);
        console.error("Failed to get number of current players");
    }
}