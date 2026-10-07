document.addEventListener("DOMContentLoaded", async () => {
    const users = await getNumUsers();
    document.getElementById("numUsers").textContent = users;
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