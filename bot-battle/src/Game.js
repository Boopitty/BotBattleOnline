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
        this.assault = this.add.sprite(824, 500, 'Assault');
        this.assault.play('Assault_Idle');
        this.assault.setScale(15);
        this.assault.flipX = true;
        this.assault.setInteractive().on('pointerup', () => {
            sendRequest(makeRequest("leave-game"));
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
            const resp = sendRequest(makeRequest("new-game"));
            console.log(resp.message);
        });

        this.anims.create({
            key: 'flag_Idle',
            frames: this.anims.generateFrameNumbers('flag', {
                start: 0,
                end: 1
            }),    
            repeat: -1,
            frameRate: 4
        });

        // Handle the flag sprite
        this.flag = this.add.sprite(524, 564, 'flag');
        this.flag.play('flag_Idle');
        this.flag.setScale(5);
        this.flag.setInteractive().on('pointerdown', () => {
            const resp = sendRequest(makeRequest("join-game"));
            console.log(resp.message);
        }); 
    }
}

/** 
 * @Param {string} command - The command to send
 * @Param {number} input - The input for the command
 * @returns {string} - The request string in JSON format
 */
function makeRequest (command, input = 0) {
    const req = JSON.stringify({
        User: {
            username: localStorage.getItem("username"),
        },
        command: command,
        input: input
    });
    return req;
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
