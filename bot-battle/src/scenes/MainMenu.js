import Phaser from "phaser";
import { initWebSocket } from "../websocket.js"
export default class MainMenu extends Phaser.Scene
{
    constructor ()
    {
        super('MainMenu');
        this.team = [];
    }

    init (team)
    {
        this.team = team;
    }

    create ()
    {
        const cx = this.scale.width / 2;
        const cy = this.scale.height / 2;

        this.add.text(cx, cy - 40, 'Bot Battle', {
            font: '64px monospace',
            color: '#8bb2ff',
        }).setOrigin(0.5);

        // Place the start button
        this.makeButton(
            cx,
            100,
            3,
            'sci_fi_buttons',
            'start',
            'start_01',
            () => {
                if (this.team.length > 0) {
                    this.scene.start('Battle', this.team);
                } else {
                    console.warn('You cannnot fight without a team!')
                }
            }
        )

        // Handle the flag sprite
        this.createAnim("Flag", "Idle", 0, 5, -1, 6);
        this.flag = this.add.sprite(cx/2, 564, 'Flag').play('Flag_Idle').setScale(5);
        this.flag.setInteractive().on('pointerup', () => {
            this.scene.start('TeamSelect', this.team);
        }); 
    }

    /**
     * Creates an animated button.
     * @param {Number} x
     * @param {Number} y 
     * @param {Number} scale 
     * @param {String} texture 
     * @param {String} btnName 
     * @param {String} idleFrame
     * @param {func} interact 
    */
    makeButton (x, y, scale, texture, btnName, idleFrame, interact)
    {
        const button = this.add.sprite(x, y, texture, idleFrame).setScale(scale);

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
