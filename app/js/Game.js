import Boot from "./scenes/Boot.js"
import MainMenu from "./scenes/MainMenu.js"
import Preloader from "./scenes/Preloader.js"

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

        this.add.text(cx, cy, 'Your game starts here', {
            font: '24px monospace',
            color: '#ffe1da',
        }).setOrigin(0.5);

        this.add.image(900, 100, 'Login_Button').setInteractive().on('pointerdown', () => {
            console.log('logging in...')
        })
    }
}

const config = {
    type: Phaser.AUTO,
    width: 1024,
    height: 768,
    parent: 'game-container',
    backgroundColor: '#000000',
    scale: {
        mode: Phaser.Scale.FIT
    },
    scene: [Boot, Preloader, MainMenu, Game],
};

// Bind the game to a top-level identifier
const game = new Phaser.Game(config);
