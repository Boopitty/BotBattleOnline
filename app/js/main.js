// js/main.js
document.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');

  if (token) {
    console.log('token found...');
    document.getElementById('login').textContent = 'logout';
  } else {
    console.log('token not found...')
    document.getElementById('login').textContent = 'login';
  }
});

document.getElementById('login').addEventListener('click', () => {
    const btn = document.getElementById('login'); // get the button

    // Change operation depending on the button's text
    if (btn.textContent === 'login') {
        console.log('logging in...');
        btn.textContent = 'logout';
        
    } else if (btn.textContent === 'logout') {
        console.log('logging out...');
        btn.textContent = 'login';
    } else {
        console.log('Invalid');
    }
    
});

document.getElementById('command-form').addEventListener('submit', async (event) => {
    event.preventDefault();
    await command();
});

async function command() {
    const command = document.getElementById("command").value;
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
}

function commandHandler(log) {
    const text = document.createElement('div');
    text.textContent = log;

    container = document.getElementById("command-log");
    container.prepend(text);
}


