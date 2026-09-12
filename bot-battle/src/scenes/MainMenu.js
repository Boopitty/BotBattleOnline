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

        this.rightButton = this.add.image(cx, cy + 24, 'Login_Button').setInteractive().on('pointerup', () => {
            initWebSocket();
            this.scene.start('Game');
        });
        this.rightButton.setScale(2);
    }
}
