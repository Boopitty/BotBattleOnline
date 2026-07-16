// Bulk asset loader. The block between the marker comments below is
// auto-generated from the project's asset list (Code Editor → Assets
// tab). Don't edit between the markers — your changes will be
// overwritten on the next Run. You can still add your own
// `this.load.*` lines anywhere outside the marker block.

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
        this.load.image('Login_Button', '../../assets/Login_Button.png');
        this.load.image('dagger-64', '../../assets/dagger-64.png');
        /* phaser:assets:start */
        /* phaser:assets:end */
    }

    create ()
    {
        /* phaser:assets:setup:start */
        /* phaser:assets:setup:end */

        this.scene.start('MainMenu');
    }
}
