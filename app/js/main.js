// js/main.js
document.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');

  if (token) {
    console.log('token found...');
    document.getElementById('login').style.display = 'none';
    document.getElementById('cancel-login').style.display = 'none';
    document.getElementById('logout').style.display = 'block';
  } else {
    console.log('token not found...')
    document.getElementById('login').style.display = 'block';
    document.getElementById('cancel-login').style.display = 'none';
    document.getElementById('logout').style.display = 'none';
  }
});

// display login form
document.getElementById('login').addEventListener('click', () => {
    console.log('displaying login form...');
    document.getElementById('login').style.display = 'none';
    document.getElementById('cancel-login').style.display = 'block';
    document.getElementById('logout').style.display = 'none';

    document.getElementById('login-form').style.display = 'flex';
    document.getElementById('new-profile-form').style.display = 'none';
});

// Cancel login process
document.getElementById('cancel-login').addEventListener('click', () => {
    console.log('cancelling login...');
    document.getElementById('login').style.display = 'block';
    document.getElementById('cancel-login').style.display = 'none';
    document.getElementById('logout').style.display = 'none';

    document.getElementById('login-form').style.display = 'none';
    document.getElementById('new-profile-form').style.display = 'none';
});

// logout and clear credentials
document.getElementById('logout').addEventListener('click', () => {
    console.log('logging out...');
    document.getElementById('login').style.display = 'block';
    document.getElementById('cancel-login').style.display = 'none';
    document.getElementById('logout').style.display = 'none';

    document.getElementById('login-form').style.display = 'none';
    document.getElementById('new-profile-form').style.display = 'none';
});

// switch from the 'login' form to the 'new profile' form
document.getElementById('switch-to-new').addEventListener('click', () => {
    document.getElementById('login-form').style.display = 'none';
    document.getElementById('new-profile-form').style.display = 'flex';
})

document.getElementById('switch-to-login').addEventListener('click', () => {
    document.getElementById('login-form').style.display = 'flex';
    document.getElementById('new-profile-form').style.display = 'none';
})

document.getElementById('submit-login').addEventListener('click', () => {
    console.log('logging in...');
})

document.getElementById('submit-new-profile').addEventListener('click', () => {
    console.log('creating profile...');
})

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


