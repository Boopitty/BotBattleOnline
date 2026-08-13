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

        this.rightButton = this.add.image(900, 100, 'Login_Button').setInteractive().on('pointerdown', () => {
            initWebSocket();
        });
        this.rightButton.setScale(2);

        this.leftButton = this.add.image(100, 100, 'Logout_Button').setInteractive().on('pointerdown', () => {
            closeWebSocket();
        })
        this.leftButton.setScale(3);

        this.anims.create({
            key: 'Assault_Idle',
            frames: this.anims.generateFrameNumbers('Assault', {
            start: 0,
            end: 1
        }),    
            repeat: -1,
            frameRate: 2
        });

        // Handle the Assault_Class sprite
        this.assault = this.add.sprite(824, 500, 'Assault')
        this.assault.play('Assault_Idle');
        this.assault.setScale(15);
        this.assault.flipX = true;
        this.assault.setInteractive().on('pointerdown', () => {
            this.makeRequest("attack");
        });

        this.anims.create({
            key: 'Spider_Idle',
            frames: this.anims.generateFrameNumbers('Spider', {
            start: 0,
            end: 1
        }),    
            repeat: -1,
            frameRate: 2
        });

        // Handle the Spider sprite
        this.spider = this.add.sprite(224, 500, 'Spider')
        this.spider.play('Spider_Idle');
        this.spider.setScale(15);
        this.spider.setInteractive().on('pointerdown', () => {
            this.makeRequest("attack");
        });
    }

    makeRequest (type) 
    {
        sendCommand(JSON.stringify({
            command: type
        }));
    }
}

// Send a command through the websocket. req is a struct in JSON format.
async function sendCommand(req) {
    try {
        const resp = await sendRequest(req);

    } catch (error) {
        console.error("Error sending command: ", error);
    };
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
