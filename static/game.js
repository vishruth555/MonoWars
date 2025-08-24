let ws;

class GameScene extends Phaser.Scene {
  constructor() {
    super("GameScene");
  }
  
  preload() {
    this.load.image("tileW", "assets/player_white.png"); 
    this.load.image("tileB", "assets/player_black.png");
    this.load.image("player", "Player.png");
  }
  create() {

    const mapLayout = [
      [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
      [0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0],
      [0, 2, 2, 2, 2, 2, 2, 1, 1, 1, 2, 0],
      [0, 2, 2, 2, 2, 2, 2, 1, 1, 1, 2, 0],
      [0, 2, 2, 2, 2, 2, 2, 1, 1, 1, 2, 0],
      [0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0],
      [0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0],
      [0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0],
      [0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0],
      [0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0],
      [0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0],
      [0, 1, 2, 2, 2, 1, 1, 1, 1, 1, 1, 0],
      [0, 1, 2, 2, 2, 1, 1, 1, 1, 1, 1, 0],
      [0, 1, 2, 2, 2, 1, 1, 1, 1, 1, 1, 0],
      [0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0],
      [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    ];

    const rows = mapLayout.length;
    const cols = mapLayout[0].length;

    // Calculate the largest possible tile size that fits both width and height
    const tileWidth = this.sys.game.config.width / cols;
    const tileHeight = this.sys.game.config.height / rows;
    const tileSize = Math.min(tileWidth, tileHeight);

    // Calculate total grid height and vertical padding
    const gridHeight = tileSize * rows;
    const gridWidth = tileSize * cols;
    const verticalPadding = ((this.sys.game.config.height - gridHeight) / 2);
    const horizontalPadding = (this.sys.game.config.width - gridWidth) / 2;

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
          x * tileSize + horizontalPadding,
          y * tileSize + verticalPadding,
          tileKey
        ).setOrigin(0).setDisplaySize(tileSize, tileSize);
      }
    }

    let xPos = horizontalPadding + (tileSize * 3);
    let yPos = verticalPadding + (tileSize * 12);
    this.playerTilePos = { x: 3, y: 12 };
    console.log("Player tile position:", this.playerTilePos);

    this.playerSprite = this.add.image(xPos, yPos, "player").setOrigin(0).setDisplaySize(tileSize, tileSize);
    console.log(this.playerSprite);

      // Button event handlers
  document.getElementById("up").onclick = () => this.movePlayer(0, -0.5, tileSize, horizontalPadding, verticalPadding);
  document.getElementById("down").onclick = () => this.movePlayer(0, 0.5, tileSize, horizontalPadding, verticalPadding);
  document.getElementById("left").onclick = () => this.movePlayer(-0.5, 0, tileSize, horizontalPadding, verticalPadding);
  document.getElementById("right").onclick = () => this.movePlayer(0.5, 0, tileSize, horizontalPadding, verticalPadding);


  }
  movePlayer(dx, dy, tileSize, horizontalPadding, verticalPadding) {
    console.log(this.playerTilePos)
  this.playerTilePos.x += dx;
  this.playerTilePos.y += dy;
  this.playerSprite.x = horizontalPadding + (tileSize * this.playerTilePos.x);
  this.playerSprite.y = verticalPadding + (tileSize * this.playerTilePos.y);
  console.log(this.playerSprite);
}

  update() {

  }
}

// class ShooterScene extends Phaser.Scene {
//   constructor() {
//     super("ShooterScene");
//   }

//   preload() {
//     this.load.image("player", "Player.png"); // add your assets
//     this.load.image("tile", "tile.png");
//   }

//   create() {
//     // Render repeating tile background
//     this.add.tileSprite(0, 0, this.sys.game.config.width, this.sys.game.config.height, "tile").setOrigin(0, 0);

//     // Connect to Go backend WebSocket
//     ws = new WebSocket("ws://localhost:8080/ws");

//     ws.onmessage = (event) => {
//       const data = JSON.parse(event.data);
//       console.log("Received data:", data);
//       // this.add.sprite(100, 100, "tile");

//       // Handle diffs: update players
//       if (data.players) {
//         this.players.forEach(p => p.destroy());
//         this.players = data.players.map(playerData =>
//           this.add.sprite(playerData.x, playerData.y, "player")
//         );
//       }
//     };

//     this.players = [];
//   }

//   update() {
//     // handle input locally, send to server if needed
//     // if (ws && ws.readyState === WebSocket.OPEN) {
//     //   if (this.input.keyboard.isDown(Phaser.Input.Keyboard.KeyCodes.W)) {
//     //     ws.send(JSON.stringify({ action: "move", dir: "up" }));
//     //   }
//     // }
//   }
// }

class TempScene extends Phaser.Scene {
  constructor() {
    super("TempScene");
  }

  preload() {
  }

  create() {


    // Connect to Go backend WebSocket
    ws = new WebSocket(`ws://${location.host}/ws`);
    console.log("WebSocket initialized:", ws);
    sendInput()


    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("Received data:", data);
      if (data.type == "ping"){
        ws.send(JSON.stringify({type: "pong"}))
      }

      // // Handle diffs: update players
    };

    this.players = [];
  }

  update() {
    // handle input locally, send to server if needed
    // if (ws && ws.readyState === WebSocket.OPEN) {
    //   if (this.input.keyboard.isDown(Phaser.Input.Keyboard.KeyCodes.W)) {
    //     ws.send(JSON.stringify({ action: "move", dir: "up" }));
    //   }
    // }
  }
}












const parent = document.getElementById("gameArea");
const config = {
  type: Phaser.AUTO,
  parent: parent,
  backgroundColor: "#8d8d8d", 
  width: parent.clientWidth,
  height: parent.clientHeight,
  scene: TempScene
};

new Phaser.Game(config);



// input handlers


function sendInput() {
  document.getElementById('shoot').addEventListener('click',() => emitMove(1,0,0));
  hold(document.getElementById('up'),0,-0.1);
  hold(document.getElementById('down'),0,0.1);
  hold(document.getElementById('left'),-0.1,0);
  hold(document.getElementById('right'),0.1,0);
}
function emitMove(type, dx, dy){
  console.log(JSON.stringify({type:type,dx:dx,dy:dy}))
  ws.send(JSON.stringify({type:type,dx:dx,dy:dy}));
}
function hold(btn,dx,dy){
  btn.addEventListener('touchstart',e=>{e.preventDefault();emitMove(0,dx,dy);});
  btn.addEventListener('touchend',e=>{e.preventDefault();emitMove(0,0,0);});
  btn.addEventListener('mousedown',()=>{emitMove(0,dx,dy);});
  btn.addEventListener('mouseup',()=>{emitMove(0,0,0);});
}