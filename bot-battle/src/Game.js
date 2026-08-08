import Phaser from "phaser";
import Boot from "./scenes/Boot.js";
import MainMenu from "./scenes/MainMenu.js";
import Preloader from "./scenes/Preloader.js";
import {initWebSocket, closeWebSocket, sendRequest} from "./websocket.js"

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

        this.add.text(cx, cy, 'Game starts here', {
            font: '24px monospace',
            color: '#ffe1da',
        }).setOrigin(0.5);

        const rightButton = this.add.image(900, 100, 'Login_Button').setInteractive().on('pointerdown', () => {
            initWebSocket();
        });

        const leftButton = this.add.image(100, 100, 'Login_Button').setInteractive().on('pointerdown', () => {
            closeWebSocket();
        })

        this.anims.create({
            key: 'Assault_Idle',
            frames: this.anims.generateFrameNumbers('Assault_Class', {
            start: 0,
            end: 1
        }),    
            repeat: -1,
            frameRate: 2
        });

        const assault = this.add.sprite(824, 500, 'Assault_Class')
        assault.play('Assault_Idle'); // Play the idle animation for the Assault_Class sprite
        assault.setScale(15); // Scale the sprite up by a factor of 15
        assault.flipX = true; // Flip the sprite horizontally

        this.anims.create({
            key: 'Spider_Idle',
            frames: this.anims.generateFrameNumbers('Spider', {
            start: 0,
            end: 1
        }),    
            repeat: -1,
            frameRate: 2
        });

        const spider = this.add.sprite(224, 500, 'Spider')
        spider.play('Spider_Idle'); // Play the idle animation for the Spider sprite
        spider.setScale(15); // Scale the sprite up by a factor of 15
    }
}

const config = {
    type: Phaser.AUTO,
    width: 1024,
    height: 768,
    pixelArt: true,
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
