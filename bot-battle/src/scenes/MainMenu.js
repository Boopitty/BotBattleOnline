import Phaser from "phaser";
import { initWebSocket } from "../websocket.js"
export default class MainMenu extends Phaser.Scene
{
    constructor ()
    {
        super('MainMenu');
    }

    create ()
    {
        const cx = this.scale.width / 2;
        const cy = this.scale.height / 2;

        this.add.text(cx, cy - 40, 'Bot Battle', {
            font: '64px monospace',
            color: '#8bb2ff',
        }).setOrigin(0.5);

        // Place login button
        this.rightButton = this.add.image(100, 100, 'Login_Button').setInteractive().on('pointerup', () => {
            initWebSocket();
            this.scene.start('Game');
        });
        this.rightButton.setScale(2);

        // Place the start button
        const startButton = this.makeButton(
            cx/2 + 100,
            100,
            3,
            'sci_fi_buttons',
            'start',
            'start_01', 'start_02', 'start_03'
        )

        // Handle the flag sprite
        this.createAnim("Flag", "Idle", 0, 5, -1, 6);
        this.flag = this.add.sprite(cx/2, 564, 'Flag').play('Flag_Idle').setScale(5);
        this.flag.setInteractive().on('pointerup', () => {
            this.scene.start('TeamSelect');
        }); 
    }

    makeButton (x, y, scale, texture, btnName, idleFrame, movingFrame, downFrame)
    {
        const button = this.add.sprite(x, y, texture, idleFrame).setScale(scale);

        button.setInteractive({ useHandCursor: true });

        this.createButtonAnim(btnName, 'press', texture, 1, 3);
        this.createButtonAnim(btnName, 'release', texture, 3, 1);

        button.on('pointerdown', () => {
            button.play(`${btnName}_press`)
        })
        button.on('pointerup', () => {
            button.play(`${btnName}_release`)
        })
    }

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
            frameRate: 30
        })
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
}
