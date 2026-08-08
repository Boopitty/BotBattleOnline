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
        const spritesheets = new Map([
            ["AnitTank_Class", "Soldiers/AnitTank-Class.png"],
            ["Assault_Class", "Soldiers/Assault-Class.png"],
            ["Grenadier_Class", "Soldiers/Grenadier-Class.png"],
            ["MachineGunner_Class", "Soldiers/MachineGunner-Class.png"],
            ["RadioOperator_Class", "Soldiers/RadioOperator-Class.png"],
            ["Sniper_Class", "Soldiers/Sniper-Class.png"],
            ["SquadLeader_Class", "Soldiers/SquadLeader-Class.png"],
            ["Centipede", "Robots/Centipede.png"],
            ["Hornet", "Robots/Hornet.png"],
            ["Scarab", "Robots/Scarab.png"],
            ["Spider", "Robots/Spider.png"],
            ["Wasp", "Robots/Wasp.png"],
        ]);

        for (const [fileName, filePath] of spritesheets) {
            this.load.spritesheet(
                fileName, 
                `../../assets/Robot Warfare Asset Pack 22-11-24/${filePath}`, 
                { frameWidth: 16, frameHeight: 16 }
            );
        }
        this.load.image('Login_Button', '../../assets/Login_Button.png');
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
