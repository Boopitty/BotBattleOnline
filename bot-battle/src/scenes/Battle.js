import Phaser from "phaser";
import { initWebSocket, closeWebSocket, makeRequest, sendRequest } from "../websocket.js"
export default class Battle extends Phaser.Scene
{
    constructor ()
    {
        super('Battle');
        this.team = [];
    }

    init (team = [])
    {
        this.team = team;
        initWebSocket();
    }

    create ()
    {
        this.makeButtons();

        // Handle the flag sprite
        this.createAnim("Flag", "Idle", 0, 5, -1, 6);
        this.flag = this.add.sprite(524, 564, 'Flag').play('Flag_Idle').setScale(5);
        this.flag.setInteractive().on('pointerup', () => {
            sendRequest(makeRequest("join-game"));
        }); 
    }

    // Creates all buttons for this scene
    makeButtons ()
    {
        const cx = this.scale.width / 2;

        this.makeButton(
            cx / 2,
            100,
            3,
            'sci_fi_buttons',
            'quit',
            () => {
                closeWebSocket();
                this.scene.start('MainMenu', this.team);
            }
        );

        this.makeButton(
            cx,
            100,
            3,
            'sci_fi_buttons',
            'leave',
            () => {
                sendRequest(makeRequest("leave-game"));
            }
        );

        this.makeButton(
            cx * 1.5,
            100,
            3,
            'sci_fi_buttons',
            'start',
            () => {
                sendRequest(makeRequest("new-game"));
            }
        );
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
    createAnim (botName, animType, startFrame, endFrame, repeat, frameRate)
    {
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

    /**
     * Creates an animated button.
     * @param {Number} x
     * @param {Number} y 
     * @param {Number} scale 
     * @param {String} texture 
     * @param {String} btnName 
     * @param {func} interact 
    */
    makeButton (x, y, scale, texture, btnName, interact)
    {
        const button = this.add.sprite(x, y, texture, `${btnName}_01`).setScale(scale);

        button.setInteractive({ useHandCursor: true });

        this.createButtonAnim(btnName, 'press', texture, 1, 3);
        this.createButtonAnim(btnName, 'release', texture, 3, 1);

        button.on('pointerdown', () => {
            button.play(`${btnName}_press`);
            button.play(`${btnName}_release`);
        })
        button.on('pointerup', interact);
        return button
    }

    /**
     * Constructs animation for given button
     * @param {String} btnName
     * @param {String} animType
     * @param {String} texture
     * @param {Number} start
     * @param {Number} end
    */ 
    createButtonAnim (btnName, animType, texture, start, end)
    {
        const animName = `${btnName}_${animType}`
        if (this.anims.exists(animName)) {
            return;
        }

        this.anims.create({
            key: animName,
            frames: this.anims.generateFrameNames(texture, {
                prefix: `${btnName}_`,
                start: start,
                end: end,
                zeroPad: 2
            }),
            repeat: 0,
            frameRate: 10
        });
    }
}
