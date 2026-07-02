document.getElementById('command-form').addEventListener('submit', async (event) => {
    event.preventDefault();
    await command();
});

async function command() {
    const command = document.getElementById('command').value;
    console.log(command);
    try{
        // make a request to the server and wait for response
        const res = await fetch('/api/command', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ command }),
            });

        // Capture the json response data
        const data = await res.json();
        if (!res.ok) {
            throw new Error(`Failed to create video draft: ${data.error}`);
        }

        const log = data.log;
        
        if (log) {
            await commandHandler(log);
        }
        
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
};

function commandHandler(log) {
    const text = document.createElement('div');
    text.textContent = log;

    container = document.getElementById("command-log");
    container.prepend(text);
};


