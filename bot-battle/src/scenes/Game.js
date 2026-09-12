import Phaser from "phaser";
import Boot from "./Boot.js";
import MainMenu from "./MainMenu.js";
import Preloader from "./Preloader.js";
import { initWebSocket, closeWebSocket, makeRequest, sendRequest } from "../websocket.js"

class Game extends Phaser.Scene
{
    constructor ()
    {
        super('Game');
        const team = [];
    }

    create ()
    {
        const cx = this.scale.width / 2;
        const cy = this.scale.height / 2;

        this.leftButton = this.add.image(100, 100, 'Logout_Button').setInteractive().on('pointerup', () => {       
            closeWebSocket();
            this.scene.start('MainMenu');
        })
        this.leftButton.setScale(3);

        // Handle the Assault_Class sprite
        this.createAnims("Assault", "Idle", 0, 1, -1, 2);
        this.assault = this.add.sprite(824, 500, 'Assault');
        this.assault.play('Assault_Idle');
        this.assault.setScale(15);
        this.assault.flipX = true;
        this.assault.setInteractive().on('pointerup', () => {
            sendRequest(makeRequest("leave-game"));
        });
        
        // Handle the Spider sprite
        this.createAnims("Spider", "Idle", 0, 1, -1, 2);
        this.spider = this.add.sprite(224, 500, 'Spider')
        this.spider.play('Spider_Idle');
        this.spider.setScale(15);
        this.spider.setInteractive().on('pointerup', () => {
            sendRequest(makeRequest("new-game"));
        });

        // Handle the flag sprite
        this.createAnims("flag", "Idle", 0, 5, -1, 6);
        this.flag = this.add.sprite(524, 564, 'flag');
        this.flag.play('flag_Idle');
        this.flag.setScale(5);
        this.flag.setInteractive().on('pointerup', () => {
            sendRequest(makeRequest("join-game"));
        }); 
    }

    /**
     * Creates animations for a given bot.
     * @param {string} botName 
     * @param {string} animType 
     * @param {integer} startFrame 
     * @param {integer} endFrame 
     * @param {integer} repeat 
     * @param {integer} frameRate 
     */
    createAnims(botName, animType, startFrame, endFrame, repeat, frameRate) {
        if (this.anims.exists(`${botName}_${animType}`)) {
            return;
        }
        this.anims.create({
            key: `${botName}_${animType}`,
            frames: this.anims.generateFrameNumbers(botName, {
                start: startFrame,
                end: endFrame
            }),    
            repeat: repeat,
            frameRate: frameRate
        });
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
