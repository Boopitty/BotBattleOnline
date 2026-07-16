import Boot from "./scenes/Boot.js"
import MainMenu from "./scenes/MainMenu.js"
import Preloader from "./scenes/Preloader.js"

class Game extends Phaser.Scene
{
    constructor ()
    {
        super('Game');
    }

    create ()
    {
        const cx = this.scale.width / 2;
        const cy = this.scale.height / 2;

        this.add.text(cx, cy, 'Your game starts here', {
            font: '24px monospace',
            color: '#ffe1da',
        }).setOrigin(0.5);

        this.add.image(900, 100, 'Login_Button').setInteractive().on('pointerdown', () => {
            console.log('logging in...')
        })
        this.add.image(100, 600, 'dagger-64').setScale(5, 5)
    }
}

const config = {
    type: Phaser.AUTO,
    width: 1024,
    height: 768,
    parent: 'game-container',
    backgroundColor: '#000000',
    scale: {
        mode: Phaser.Scale.FIT
    },
    scene: [Boot, Preloader, MainMenu, Game],
}

// Bind the game to a top-level identifier
const game = new Phaser.Game(config);

// Event listender for the websocket. Salvaged from the deleted Canvas.js, to-be adjusted for Game.js.
/*
document.getElementById("game-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const command = document.getElementById("game-input");

    try {
        
        const res = await sendRequest(
            JSON.stringify(
                {
                    user: user,
                    command: command.value
                })
            );

    } catch (error) {
        console.error("Error sending command:", error);
    };
    
    command.value = "";
});
*/
