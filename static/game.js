
const parent = document.getElementById("gameArea");
const config = {
  type: Phaser.AUTO,
  parent: parent,
  backgroundColor: "#8d8d8d", 
  width: parent.clientWidth,
  height: parent.clientHeight,
  scene: [LobbyScene, LoadingScene, MainScene]
};

new Phaser.Game(config);

let ws;

// input handlers

function enableInput() {
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