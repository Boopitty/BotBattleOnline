// Bulk asset loader. The block between the marker comments below is
// auto-generated from the project's asset list (Code Editor → Assets
// tab). Don't edit between the markers — your changes will be
// overwritten on the next Run. You can still add your own
// `this.load.*` lines anywhere outside the marker block.
import Phaser from "phaser";

export default class Preloader extends Phaser.Scene
{
    constructor ()
    {
        super('Preloader');
    }

    init ()
    {
        const cx = this.scale.width / 2;
        const cy = this.scale.height / 2;
        this.add.rectangle(cx, cy, 468, 32).setStrokeStyle(1, 0xffffff);
        const bar = this.add.rectangle(cx - 230, cy, 4, 28, 0xffffff);
        this.load.on('progress', (progress) =>
        {
            bar.width = 4 + (460 * progress);
        });
    }

    preload ()
    {
        const spritesheets = [
            ["AntiTank", "Soldiers/AntiTank-Class.png", 16, 16],
            ["Assault", "Soldiers/Assault-Class.png", 16, 16],
            ["Grenadier", "Soldiers/Grenadier-Class.png", 16, 16],
            ["MachineGunner", "Soldiers/MachineGunner-Class.png", 16, 16],
            ["RadioOperator", "Soldiers/RadioOperator-Class.png", 16, 16],
            ["Sniper", "Soldiers/Sniper-Class.png", 16, 16],
            ["SquadLeader", "Soldiers/SquadLeader.png", 16, 16],
            ["Centipede", "Robots/Centipede.png", 16, 16],
            ["Hornet", "Robots/Hornet.png", 16, 16],
            ["Scarab", "Robots/Scarab.png", 16, 16],
            ["Spider", "Robots/Spider.png", 16, 16],
            ["Wasp", "Robots/Wasp.png", 16, 16],
            ["flag", "UI/objective-flag.png", 64, 64],
        ];

        for (const [fileName, filePath, width, height] of spritesheets) {
            this.load.spritesheet(
                fileName, 
                `../../assets/Robot Warfare Asset Pack 22-11-24/${filePath}`, 
                { frameWidth: width, frameHeight: height }
            );
        }
        this.load.image('Login_Button', '../../assets/Login.png');
        this.load.image('Logout_Button', '../../assets/Logout.png');
        this.load.image('dagger-64', '../../assets/dagger-64.png');
        this.load.image('dagger2-0', '../../assets/dagger2-0.png');
    }

    create ()
    {
        /* phaser:assets:setup:start */
        /* phaser:assets:setup:end */

        this.scene.start('MainMenu');
    }
}
