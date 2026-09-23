import Phaser from "phaser";
import Boot from "./scenes/Boot.js";
import MainMenu from "./scenes/MainMenu.js";
import Preloader from "./scenes/Preloader.js";
import TeamSelect from "./scenes/TeamSelect.js";
import Battle from "./scenes/Battle.js"

const config = {
    type: Phaser.AUTO,
    width: 1024,
    height: 768,
    pixelArt: true,
    parent: 'game-container',
    backgroundColor: '#000000',
    scale: {
        mode: Phaser.Scale.FIT
    },
    scene: [Boot, Preloader, MainMenu, TeamSelect, Battle],
}

// Bind the game to a top-level identifier
const game = new Phaser.Game(config);