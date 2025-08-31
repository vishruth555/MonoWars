class MainScene extends Phaser.Scene {
    constructor() {
        super("MainScene");
        this.players = []
        this.bullets = []
    }

 /* ------------------------------- LOADING THE ASSETS HERE ---------------------------------- */
    preload() {
    this.load.image("tileW", "assets/tile_white.png"); 
    this.load.image("tileB", "assets/tile_black.png");
    this.load.image("bulletW", "assets/bullet_white.png");
    this.load.image("bulletB", "assets/bullet_black.png");

    this.load.image("playerW", "assets/player_white.png");
    this.load.image("playerB", "assets/player_black.png");
    
    // this.load.spritesheet("playerIce", "assets/player-ice.png", {
    //   frameWidth: 32,
    //   frameHeight: 32
    // });
    // this.load.spritesheet("playerLava", "assets/player-lava.png", {
    //   frameWidth: 32,
    //   frameHeight: 32
    // });
    }



   /* ------------------------------- MAIN ENTRY POINT FOR THE GAME ---------------------------------- */
   /* ------------------------------- Websocket router ---------------------------------- */

    create(){
    ws = new WebSocket(`ws://${location.host}/ws`);

    ws.onerror = (err) => {
     this.scene.start("ErrorScene");
    };

    this.createLoadingScreen()

  
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      // console.log("Received data:", data);

      //waiting for another player
      if (data.type == "waiting"){
        this.loadingScreen.setVisible(true)
      }

      //game start
      else if (data.type == "GameStart"){
        this.loadingScreen.setVisible(false)
        this.createTileMap(data.tileMap)
        setTimeout(() => {
                  this.createPlayer(data.player1Data)
        this.createPlayer(data.player2Data)
        }, 2000)

        this.showGameStartOverlay()
      }

      //ticks at regular intervals
      else if (data.type == "Tick") {
        this.handleDiffs(data.diff)
      }

      //game end
      else if(data.type == "GameEnd"){
        console.log("game has ended")
        console.log("congratulation player ", data.winner)
      }
    };
    }

   /* ------------------------------- UTILS ---------------------------------- */

    getX(x) {
      console.log("test")
      return Math.round(this.horizontalPadding + (this.tileSize * x));
    }
    getY(y) {
      return Math.round(this.verticalPadding + (this.tileSize * y));
    }


    /* ------------------------------- RENDERING ---------------------------------- */


    createTileMap(mapLayout){
      this.tileMap = mapLayout
      this.tileSprites = [];
      const rows = mapLayout.length;
      const cols = mapLayout[0].length;
      // Calculate the largest possible tile size that fits both width and height
      const tileWidth = this.sys.game.config.width / cols;
      const tileHeight = this.sys.game.config.height / rows;
      this.tileSize = Math.min(tileWidth, tileHeight);
      console.log("tileSize: ",this.tileSize)

      // Calculate total grid height and vertical padding
      const gridHeight = this.tileSize * rows;
      const gridWidth = this.tileSize * cols;
      this.verticalPadding = ((this.sys.game.config.height - gridHeight) / 2);
      this.horizontalPadding = (this.sys.game.config.width - gridWidth) / 2;

      for (let y = 0; y < rows; y++) {
        this.tileSprites[y] = []
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

          setTimeout(() => {
            // Place tile sprite with padding and square aspect ratio
          this.tileSprites[y][x] = this.add.image(
            this.getX(x),
            this.getY(y),
            tileKey
          ).setOrigin(0).setDisplaySize(this.tileSize, this.tileSize);
          }, (100*y)+(30*x))
          
        }
      }
      console.log("created tile map of size", rows, cols)
    }

    createPlayer(playerData){
    let xPos = this.getX(playerData.xPos+0.1);
    let yPos = this.getY(playerData.yPos+0.1);

    let player;
    
    if (playerData.id == 1){

      // this.anims.create({
      // key: "blinkIce",
      // frames: [
      //   { key: 'playerIce', frame: 0, duration: 600 },
      //   { key: 'playerIce', frame: 1, duration: 200 }  
      // ],
      // frameDuration: 2000,
      // repeat: -1
      // });

      player = {
      id: playerData.id,
      tilePos: { x: playerData.xPos, y: playerData.yPos },
      sprite: this.add
                .image(xPos, yPos, "playerW")
                .setOrigin(0)
                .setDisplaySize(this.tileSize*0.8, this.tileSize*0.8),
    };
    


    } else {

      // this.anims.create({
      // key: "blinkLava",
      // frames: [
      //   { key: 'playerLava', frame: 0, duration: 600 }, 
      //   { key: 'playerLava', frame: 1, duration: 200 }  
      // ],
      // frameDuration: 2000,
      // repeat: -1
      // });

      player = {
      id: playerData.id,
      tilePos: { x: playerData.xPos, y: playerData.yPos },
      sprite: this.add
                .image(xPos, yPos, "playerB")
                .setOrigin(0)
                .setDisplaySize(this.tileSize*0.8, this.tileSize*0.8),
                
    };
    
  }


    this.players.push(player)
    }

    handleDiffs(diffs) {
      if (diffs == null){
        return
      }
      console.log("Handling diffs:", diffs);
      diffs.forEach(diff => {
        switch (diff.entity) {
          case "Player1Data":
            this.movePlayer(0, diff.data.xPos, diff.data.yPos)
            break;
          case "Player2Data":
            this.movePlayer(1, diff.data.xPos, diff.data.yPos)
            break;
          case "TileMapData":
            let xPos = diff.data.xPos
            let yPos = diff.data.yPos
            if (diff.data.id == 1){
              this.tileSprites[yPos][xPos].setTexture("tileW")
            }
            else if (diff.data.id == 2){
              this.tileSprites[yPos][xPos].setTexture("tileB")
            }
            break;
          case "BulletData":
              this.renderBullet(diff.data)
            break;
        }
      });

    }

    renderBullet(data){
      console.log("handling bullet data: ", data)
      const {bulletId, id, state, xPos, yPos} = data;

      if (state == "active"){
        if(!this.bullets[bulletId]) {
          let bullet
          if (id == 1) {
            bullet = this.add
          .image(this.getX(xPos), this.getY(yPos), "bulletW")                
          .setOrigin(0.5)
          // .setDisplaySize(this.tileSize*2, this.tileSize*2);
          } else {
            bullet = this.add
          .image(this.getX(xPos), this.getY(yPos), "bulletB")                
          .setOrigin(0.5)
          // .setDisplaySize(this.tileSize*2, this.tileSize*2);
          }

          this.bullets[bulletId] = bullet;
        } else {
          this.bullets[bulletId].setPosition(this.getX(xPos), this.getY(yPos))
        }
      } else {
        if (this.bullets[bulletId]) {
          this.bullets[bulletId].destroy();
          delete this.bullets[bulletId]
        }
      }

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



    /* ------------------------------- UI Overlays ---------------------------------- */

    createLoadingScreen(){
      this.loadingScreen = this.add.text(parent.clientWidth/2, parent.clientHeight/2, "waiting for player...", {
      fontSize: "20px",
      fill: "#0f0",
      backgroundColor: "#222",
      padding: { x: 10, y: 5 }
    })
      .setOrigin(0.5)
      this.loadingScreen.setVisible(false)
    }

    showGameStartOverlay() {
    // Dark background overlay
    let overlayBg = this.add.rectangle(
        this.sys.game.config.width / 2,
        this.sys.game.config.height / 2,
        this.sys.game.config.width,
        this.sys.game.config.height,
        0x000000,
        0.5 // alpha (0 transparent, 1 solid)
    );

    // Center text on screen
    let overlayText = this.add.text(
        this.sys.game.config.width / 2,
        this.sys.game.config.height / 2,
        "READY...",
        { fontSize: "48px", color: "#ffffff" }
    ).setOrigin(0.5);
    overlayBg.setDepth(1000);
    overlayText.setDepth(1001);


    // After 1 second, change text
    this.time.delayedCall(1000, () => {
        overlayText.setText("GO!");
    });

    // After 2 seconds, remove overlay and enable input
    this.time.delayedCall(2000, () => {
        overlayText.destroy();
        overlayBg.destroy(); // remove dark background
        enableInput();
    });
}



}