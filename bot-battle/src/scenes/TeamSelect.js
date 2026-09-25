import Phaser from "phaser"

export default class TeamSelect extends Phaser.Scene
{
    constructor ()
    {
        super('TeamSelect');
        this.team = [];
        this.teamSprites = [];
    }

    init (team = [])
    {
        this.team = team;
        this.teamSprites = [];
    }

    create ()
    {
        this.leftButton = this.add.image(100, 100, 'Logout_Button').setInteractive().on('pointerup', () => {       
            this.scene.start('MainMenu', this.team);
        })
        this.leftButton.setScale(3);

        this.displayTeam();
        this.displayBots();
    }

    // Clear the displayed team sprites, then re-display them based on the current team contents.
    displayTeam ()
    {
        for (const i in this.teamSprites) {
            this.teamSprites[i].removeFromDisplayList();
        }
        for (const i in this.team) {
            this.teamSprites.push(this.add.sprite(150 * (Number(i)+1), 300, this.team[i]).setScale(10))
        }
    }

    // Display all available bots.
    displayBots ()
    {
        this.createAnim("Assault", "Idle", 0, 1, -1, 2);
        this.assault = this.add.sprite(75, 568, 'Assault').play('Assault_Idle').setScale(5);
        this.assault.setInteractive().on('pointerup', () => {
            if (this.team.includes("Assault")) {
                this.removeFromTeam("Assault");
            } else {
                this.addToTeam("Assault");
            }
        });

        this.createAnim("Grenadier", "Idle", 0, 1, -1, 2);
        this.grenadier = this.add.sprite(150, 568, 'Grenadier').play('Grenadier_Idle').setScale(5);
        this.grenadier.setInteractive().on('pointerup', () => {
            if (this.team.includes("Grenadier")) {
                this.removeFromTeam("Grenadier");
            } else {
                this.addToTeam("Grenadier");
            }
        });

        this.createAnim("Sniper", "Idle", 0, 1, -1, 2);
        this.sniper = this.add.sprite(225, 568, 'Sniper').play('Sniper_Idle').setScale(5);
        this.sniper.setInteractive().on('pointerup', () => {
            if (this.team.includes("Sniper")) {
                this.removeFromTeam("Sniper");
            } else {
                this.addToTeam("Sniper");
            }
        });

        this.createAnim("Spider", "Idle", 0, 1, -1, 2);
        this.spider = this.add.sprite(75, 668, 'Spider').play('Spider_Idle').setScale(5);
        this.spider.setInteractive().on('pointerup', () => {
            if (this.team.includes("Spider")) {
                this.removeFromTeam("Spider");
            } else {
                this.addToTeam("Spider");
            }
        });
    }

    addToTeam (botName)
    {
        if (this.team.length < 6) {
            if (!this.team.includes(botName)) {
                this.team.push(botName);
                console.log(`Added ${botName} to team. Current team: ${this.team}`);
            } else {
                console.log(`${botName} is already in the team.`);
            }
        } else {
            console.log("Team is full. Cannot add more bots.");
        }
        this.displayTeam();
    }

    removeFromTeam (botName)
    {
        const index = this.team.indexOf(botName);
        if (index > -1) {
            this.team.splice(index, 1);
            console.log(`Removed ${botName} from team. Current team: ${this.team}`);
        } else {
            console.log(`${botName} is not in the team.`);
        }
        this.displayTeam();
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