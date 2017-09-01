//Create the renderer
var renderer = PIXI.autoDetectRenderer(800, 600);
renderer.backgroundColor = 0xffffff;

//Add the canvas to the HTML document
document.body.appendChild(renderer.view);

//Create a container object called the `stage`
var stage = new PIXI.Container();



var graphics = new PIXI.Graphics();

graphics.beginFill(0xFFFF00);

// set the line style to have a width of 5 and set the color to red
graphics.lineStyle(5, 0xFF0000);

// draw a rectangle
graphics.drawRect(50, 50, 10, 10);

stage.addChild(graphics);



renderer.render(stage);

var ws = new WebSocket("ws://localhost:3000/ws");
ws.onopen = function() {
  send("PING", {"name": "jaska"}, {"x": 1});
};

function send(name, opts, payload) {
  str = name + "\n"

  arr = []
  for (var key in opts) {
    arr.push(key + ":" + opts[key]);
  }
  opts = arr.join(';');

  str += opts + "\n"

  str += JSON.stringify(payload) + "\n"

  ws.send(str);
}



function onKeyUp(key) {
  //console.log('sending');
  //ws.send(`hello\n`);
}

document.addEventListener('keyup', onKeyUp);
