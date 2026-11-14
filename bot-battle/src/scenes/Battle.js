import Phaser from "phaser";
import { initWebSocket, closeWebSocket, makeRequest, sendRequest } from "../websocket.js"
import { getSavedTeam } from "../api";
export default class Battle extends Phaser.Scene
{
    constructor ()
    {
        super('Battle');
        this.team = [];
    }

    init ()
    {
        this.team = [];
        getSavedTeam().then((savedTeam) => {
            this.team = savedTeam;
        });
        initWebSocket();
    }

    create ()
    {
        window.Battle = this;
        this.makeButtons();

        // Handle the flag sprite
        this.createAnim("Flag", "Idle", 0, 5, -1, 6);
        this.flag = this.add.sprite(524, 564, 'Flag').play('Flag_Idle').setScale(5);
        this.flag.setInteractive().on('pointerup', () => {
            sendRequest(makeRequest("join-game", this.team));
        }); 

        // Set up a zone to display active bots, and dynamically change the bot sprites when the gamestate is received
        this.p1Zone1 = this.add.container(194, 250);
        this.p1Zone2 = this.add.container(194, 350);
        this.p1HugeZone = this.add.container(194, 300);

        this.p2Zone1 = this.add.container(800, 250);
        this.p2Zone2 = this.add.container(800, 350);
        this.p2HugeZone = this.add.container(800, 300);

        // Set up a zone to display idle bots, and dynamically change the bot sprites when the gamestate is received
        this.p1IdleZones = this.add.container(194, 600);
        this.p2IdleZones = this.add.container(800, 700);
    }

    clearAllZones () {
        this.p1Zone1.removeAll(true);
        this.p1Zone2.removeAll(true);
        this.p1HugeZone.removeAll(true);
        this.p2Zone1.removeAll(true);
        this.p2Zone2.removeAll(true);
        this.p2HugeZone.removeAll(true);
        this.p1IdleZones.removeAll(true);
        this.p2IdleZones.removeAll(true);
    }
 
    printGamestate (gamestate) {
        console.log(gamestate);
    }

    displayBattle (gamestate) {
        if (!gamestate) {
            console.warn("No gamestate received.");
            return;
        }

        this.clearAllZones();
        if (gamestate.finished) {
            console.log("Game finished.");
            this.flag.setVisible(true);
            return;
        } else {
            this.flag.setVisible(false);
        }

        // Create animations for all active bots
        gamestate.p1Actives.forEach(bot => {
            this.createAnim(bot.name, "Idle", 0, 1, -1, 2);
        });
        gamestate.p2Actives.forEach(bot => {
            this.createAnim(bot.name, "Idle", 0, 1, -1, 2);
        });

        // Display the idle bots for each player
        gamestate.p1Idle.forEach((botName, index) => {
            this.createAnim(botName, "Idle", 0, 1, -1, 2);
            const idleBot = this.add.sprite(0, 0, botName).setScale(5).setInteractive()
                .on('pointerover', () => {
                    idleBot.play(`${botName}_Idle`)
                })
                .on('pointerout', () => {
                    idleBot.stop();
                });

            idleBot.setPosition(index * 100, 0);
            this.p1IdleZones.add(idleBot);
        });

        gamestate.p2Idle.forEach((botName, index) => {
            this.createAnim(botName, "Idle", 0, 1, -1, 2);
            const idleBot = this.add.sprite(0, 0, botName).setScale(5).setFlip(true, false).setInteractive()
                .on('pointerover', () => {
                    idleBot.play(`${botName}_Idle`)
                })
                .on('pointerout', () => {
                    idleBot.stop();
                });
            idleBot.setPosition(index * -100, 0);
            this.p2IdleZones.add(idleBot);
        });

        // Display the active bots for each player
        if (gamestate.p1Actives.length === 1) {
            const bot = this.add.sprite(0, 0, gamestate.p1Actives[0].name).play(`${gamestate.p1Actives[0].name}_Idle`).setScale(10)
            this.p1HugeZone.add(bot);
        } else if (gamestate.p1Actives.length === 2) {
            const bot1 = this.add.sprite(0, 0, gamestate.p1Actives[0].name).play(`${gamestate.p1Actives[0].name}_Idle`).setScale(10)
            const bot2 = this.add.sprite(0, 0, gamestate.p1Actives[1].name).play(`${gamestate.p1Actives[1].name}_Idle`).setScale(10)
            this.p1Zone1.add(bot1);
            this.p1Zone2.add(bot2);
        }

        if (gamestate.p2Actives.length === 1) {
            const bot = this.add.sprite(0, 0, gamestate.p2Actives[0].name).play(`${gamestate.p2Actives[0].name}_Idle`).setScale(10).setFlip(true, false)
            this.p2HugeZone.add(bot);
        } else if (gamestate.p2Actives.length === 2) {
            const bot1 = this.add.sprite(0, 0, gamestate.p2Actives[0].name).play(`${gamestate.p2Actives[0].name}_Idle`).setScale(10).setFlip(true, false)
            const bot2 = this.add.sprite(0, 0, gamestate.p2Actives[1].name).play(`${gamestate.p2Actives[1].name}_Idle`).setScale(10).setFlip(true, false)
            this.p2Zone1.add(bot1);
            this.p2Zone2.add(bot2);
        }
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
                this.scene.start('MainMenu');
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
                this.clearAllZones();
                this.flag.setVisible(true);
            }
        );

        this.makeButton(
            cx * 1.5,
            100,
            3,
            'sci_fi_buttons',
            'start',
            () => {
                if (this.team.length > 0) {
                    sendRequest(makeRequest("new-game", this.team));
                } else {
                    console.warn('You cannnot fight without a team!')
                }
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
     * The start and end frames for the press and release animations are hardcoded.
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
