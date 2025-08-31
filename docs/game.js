
const parent = document.getElementById("gameArea");
// const config = {
//   type: Phaser.AUTO,
//   parent: parent,
//   backgroundColor: "#8d8d8d", 
//   width: parent.clientWidth,
//   height: parent.clientHeight,
//   scene: [LobbyScene, LoadingScene, MainScene]
// };

const config = {
  type: Phaser.AUTO,
  parent: parent,
  backgroundColor: "#8d8d8d",
  scale: {
    mode: Phaser.Scale.FIT, // or RESIZE, depending on what you want
    autoCenter: Phaser.Scale.CENTER_BOTH,
    width: parent.clientWidth,
    height: parent.clientHeight
  },
  render: {
    pixelArt: true, // set true if you want sharp pixels
    roundPixels: true,
    antialias: false
  },
  resolution: window.devicePixelRatio,  // 👈 THIS is the main fix
  scene: [LobbyScene, ErrorScene, MainScene]
};


new Phaser.Game(config);

let ws;

// input handlers

const PLAYER_SPEED = 1

function enableInput() {
  document.getElementById('shoot').addEventListener('click',() => emitMove(1,0,0));
  hold(document.getElementById('up'),0,-PLAYER_SPEED);
  hold(document.getElementById('down'),0,PLAYER_SPEED);
  hold(document.getElementById('left'),-PLAYER_SPEED,0);
  hold(document.getElementById('right'),PLAYER_SPEED,0);
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