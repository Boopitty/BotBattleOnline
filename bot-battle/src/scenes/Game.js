import Phaser from "phaser";
import { initWebSocket, closeWebSocket, makeRequest, sendRequest } from "../websocket.js"
export default class Game extends Phaser.Scene
{
    constructor ()
    {
        super('Game');
        this.team = [];
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
        this.createAnim("Assault", "Idle", 0, 1, -1, 2);
        this.assault = this.add.sprite(824, 500, 'Assault').play('Assault_Idle').setScale(15);
        this.assault.flipX = true;
        this.assault.setInteractive().on('pointerup', () => {
            sendRequest(makeRequest("leave-game"));
        });
        
        // Handle the Spider sprite
        this.createAnim("Spider", "Idle", 0, 1, -1, 2);
        this.spider = this.add.sprite(224, 500, 'Spider').play('Spider_Idle').setScale(15);
        this.spider.setInteractive().on('pointerup', () => {
            sendRequest(makeRequest("new-game"));
        });

        // Handle the flag sprite
        this.createAnim("Flag", "Idle", 0, 5, -1, 6);
        this.flag = this.add.sprite(524, 564, 'Flag').play('Flag_Idle').setScale(5);
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
    createAnim(botName, animType, startFrame, endFrame, repeat, frameRate) {
        const animName = `${botName}_${animType}`
        if (this.anims.exists(animName)) {
            return;
        }
        this.anims.create({
            key: animName,
            frames: this.anims.generateFrameNumbers(botName, {
                start: startFrame,
                end: endFrame
            }),    
            repeat: repeat,
            frameRate: frameRate
        });
        return;
    };
}
