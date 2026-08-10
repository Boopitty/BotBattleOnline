// js/main.js
document.addEventListener("DOMContentLoaded", async () => {
  const username = localStorage.getItem("username");
  const token = localStorage.getItem("token");

  if (token && username) {
    console.log(`Logging in as: ${username}`);
    document.getElementById("login").style.display = "none";
    document.getElementById("cancel-login").style.display = "none";
    document.getElementById("logout").style.display = "block";

  } else {
    console.log("Username not found..")
    document.getElementById("login").style.display = "block";
    document.getElementById("cancel-login").style.display = "none";
    document.getElementById("logout").style.display = "none";
  }
});

// display login form
document.getElementById("login").addEventListener("click", () => {
    console.log("displaying login form...");
    document.getElementById("login").style.display = "none";
    document.getElementById("cancel-login").style.display = "block";
    document.getElementById("logout").style.display = "none";

    document.getElementById("login-form").style.display = "flex";
    document.getElementById("new-profile-form").style.display = "none";
});

// Cancel login process
document.getElementById("cancel-login").addEventListener("click", () => {
    console.log("cancelling login...");
    document.getElementById("login").style.display = "block";
    document.getElementById("cancel-login").style.display = "none";
    document.getElementById("logout").style.display = "none";

    document.getElementById("login-form").style.display = "none";
    document.getElementById("new-profile-form").style.display = "none";
});

// logout and clear credentials
document.getElementById("logout").addEventListener("click", () => {
    console.log("logging out...");

    localStorage.removeItem("username");
    localStorage.removeItem("token");

    document.getElementById("login").style.display = "block";
    document.getElementById("cancel-login").style.display = "none";
    document.getElementById("logout").style.display = "none";

    document.getElementById("login-form").style.display = "none";
    document.getElementById("new-profile-form").style.display = "none";
});

// switch from the "login" form to the "new profile" form
document.getElementById("switch-to-new").addEventListener("click", () => {
    document.getElementById("login-form").style.display = "none";
    document.getElementById("new-profile-form").style.display = "flex";
});

document.getElementById("switch-to-login").addEventListener("click", () => {
    document.getElementById("login-form").style.display = "flex";
    document.getElementById("new-profile-form").style.display = "none";
});

// Submit login form
document.getElementById("login-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    console.log("logging in...");

    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    try {
        const resp = await fetch("/api/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                username,
                password
            })
        });

        const data = await resp.json(); // response body
        if (!resp.ok) { // check status code
            document.getElementById("username").value = "";
            document.getElementById("password").value = "";
            throw new Error(`Login Failed: ${data.error}`);
        }

        // locally store jwt and username
        localStorage.setItem("token", data.token);
        localStorage.setItem("username", data.username);
        console.log(`logging in as ${data.username}`);
        
        // clear inputs and re-configure buttons
        document.getElementsByClassName("input-area").value = ""
        document.getElementById("login").style.display = "none";
        document.getElementById("cancel-login").style.display = "none";
        document.getElementById("logout").style.display = "block";
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
});

// submit form to create a new profile
document.getElementById("new-profile-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    console.log("creating profile...");

    const usernameBox = document.getElementById("new-username");
    const passwordBox = document.getElementById("new-password");
    const confirmPasswordBox = document.getElementById("confirm-password");

    if (passwordBox.value !== confirmPasswordBox.value) {
        alert("Error: Passwords do not match!");
        return;
    }

    try {
        const resp = await fetch("/api/createUser", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                username: usernameBox.value,
                password: passwordBox.value
            })
        });

        
        const data = await resp.json(); // response body
        if (!resp.ok) { // check status code
            document.getElementById("new-username").value = "";
            document.getElementById("new-password").value = "";
            document.getElementById("confirm-password").value = "";
            throw new Error(`Login Failed: ${data.error}`);
        }

        // locally store jwt and username
        localStorage.setItem("token", data.token);
        localStorage.setItem("username", data.username);
        console.log(`logging in as ${data.username}`);
        
        // clear inputs and re-configure buttons
        document.getElementsByClassName("input-area").value = ""
        document.getElementById("login").style.display = "none";
        document.getElementById("cancel-login").style.display = "none";
        document.getElementById("logout").style.display = "block";

    } catch (error) {
        alert(`Error: ${error.message}`);
    }
});

document.getElementById("command-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    await command();
});

async function command() {
    const command = document.getElementById("command").value;

    try{
        // make a request to the server and wait for response
        const resp = await fetch("/api/command", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ command }),
        });

        // Capture the json response data
        const data = await resp.json();
        if (!resp.ok) {
            throw new Error(`Failed to process command: ${data.error}`);
        }

        const log = data.log;
        
        if (log) {
            await commandHandler(log);
        }
        
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

function commandHandler(log) {
    const text = document.createElement("div");
    text.textContent = log;

    const logContainer = document.getElementById("command-log");
    logContainer.prepend(text);
}

async function resetUsers() {
    console.log("Resetting database...");
    try {
        const resp = await fetch("/api/resetUsers", {
            method: "DELETE"
        });
        console.log("Database reset successful");
    } catch (error) {
        alert(`Error: ${error}`);
        console.error("Database reset failed")
    }
}

