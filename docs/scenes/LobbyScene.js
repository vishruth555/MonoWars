
class LobbyScene extends Phaser.Scene {
    constructor () {
        super("LobbyScene")
    }

    preload() {}

    create() {

    const playButton = this.add.text(parent.clientWidth/2, parent.clientHeight/2, "▶ PLAY", {
      fontSize: "40px",
      fill: "#0f0",
      backgroundColor: "#222",
      padding: { x: 10, y: 5 }
    })
      .setOrigin(0.5)
      .setInteractive({ useHandCursor: true }) 
      .on("pointerdown", () => {
        this.scene.start("MainScene"); // Switch scene
      })
      .on("pointerover", () => playButton.setStyle({ fill: "#ff0" }))
      .on("pointerout", () => playButton.setStyle({ fill: "#0f0" }));
    }

}

class LoadingScene extends Phaser.Scene {
  constructor() {
    super("LoadingScene");
  }

  create() {
    this.loadingText = this.add.text(parent.clientWidth/2, parent.clientHeight/2, "test", {
      fontSize: "40px",
      fill: "#0f0",
      backgroundColor: "#222",
      padding: { x: 10, y: 5 }
    })
      .setOrigin(0.5)
    }
}