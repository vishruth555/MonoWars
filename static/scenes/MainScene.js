class MainScene extends Phaser.Scene {
    constructor() {
        super("MainScene");
        this.players = []
    }

    preload() {
    this.load.image("tileW", "assets/player_white.png"); 
    this.load.image("tileB", "assets/player_black.png");
    this.load.image("player", "Player.png");
    }

    create(){
    ws = new WebSocket(`ws://${location.host}/ws`);
    console.log("WebSocket initialized:", ws);
    
    this.createLoadingScreen()

  
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      // console.log("Received data:", data);
      if (data.type == "waiting"){
        this.loadingScreen.setVisible(true)
      }
      else if (data.type == "GameStart"){
        this.loadingScreen.setVisible(false)
        this.createTileMap(data.tileMap)
        this.createPlayer(data.player1Data)
        this.createPlayer(data.player2Data)
        this.showGameStartOverlay()
      }
      else if (data.type == "Tick") {
        this.handleDiffs(data.diff)
      }

    };


    }

    createTileMap(mapLayout){
      this.tileMap = mapLayout
      const rows = mapLayout.length;
      const cols = mapLayout[0].length;
      // Calculate the largest possible tile size that fits both width and height
      const tileWidth = this.sys.game.config.width / cols;
      const tileHeight = this.sys.game.config.height / rows;
      this.tileSize = Math.min(tileWidth, tileHeight);

      // Calculate total grid height and vertical padding
      const gridHeight = this.tileSize * rows;
      const gridWidth = this.tileSize * cols;
      this.verticalPadding = ((this.sys.game.config.height - gridHeight) / 2);
      this.horizontalPadding = (this.sys.game.config.width - gridWidth) / 2;

      for (let y = 0; y < rows; y++) {
        for (let x = 0; x < cols; x++) {
          const cell = mapLayout[y][x];

          if (cell === 0) {
            continue;
          }

          let tileKey;
          if (cell === 1) {
            tileKey = "tileW";
          } else if (cell === 2) {
            tileKey = "tileB";
          }

          // Place tile sprite with padding and square aspect ratio
          this.add.image(
            x * this.tileSize + this.horizontalPadding,
            y * this.tileSize + this.verticalPadding,
            tileKey
          ).setOrigin(0).setDisplaySize(this.tileSize, this.tileSize);
        }
      }
      console.log("created tile map of size", rows, cols)
    }

    getX(x) {
      return this.horizontalPadding + (this.tileSize * x);
    }
    getY(y) {
      return this.verticalPadding + (this.tileSize * y);
    }

    createPlayer(playerData){
    let xPos = this.getX(playerData.xPos);
    let yPos = this.getY(playerData.yPos);

    const player = {
      id: playerData.id,
      tilePos: { x: playerData.xPos, y: playerData.yPos },
      sprite: this.add
                .image(xPos, yPos, "player")
                .setOrigin(0)
                .setDisplaySize(this.tileSize, this.tileSize),

    }
    this.players.push(player)
    }

    createLoadingScreen(){
      this.loadingScreen = this.add.container(400, 300);
      this.loadingScreen.setVisible(false)
      const bg = this.add.rectangle(0, 0, 800, 600, 0x000000, 0.7);
      const text = this.add.text(0, 0, "Waiting for opponent...", {
        fontSize: "28px",
        fill: "#fff"
      }).setOrigin(0.5);

      this.loadingScreen.add([bg, text]);
    }

    showGameStartOverlay() {
    // Center text on screen
    let overlayText = this.add.text(
        this.sys.game.config.width / 2,
        this.sys.game.config.height / 2,
        "READY...",
        { fontSize: "48px", color: "#ffffff" }
    ).setOrigin(0.5);

    // After 1 second, change text
    this.time.delayedCall(1000, () => {
        overlayText.setText("GO!");
    });

    // After 2 seconds, remove overlay and enable input
    this.time.delayedCall(2000, () => {
        overlayText.destroy();
        enableInput(); 
    });
}
    handleDiffs(diffs) {
      if (diffs == null){
        return
      }
      console.log("Handling diffs:", diffs);
      diffs.forEach(diff => {
        switch (diff.entity) {
          case "player1Data":
            this.movePlayer(0, diff.data.xPos, diff.data.yPos)
            break;
          case "player2Data":
            this.movePlayer(1, diff.data.xPos, diff.data.yPos)
            break;
        }
      });

    }

    movePlayer(index, newX, newY) {
      console.log("moving player ", index, " from ",this.players[index].tilePos, " to ", newX, newY)
      let xPos = this.getX(newX)
      let yPos = this.getY(newY)
      this.players[index].tilePos = {x: newX, y: newY}
      // this.players[index].sprite.setPosition(xPos, yPos)
      this.tweens.add({
        targets: this.players[index].sprite,
        x: xPos,
        y: yPos,
        duration: 80, // duration of the tween in milliseconds
        ease: 'Linear' // easing function
      });

    }



    

}