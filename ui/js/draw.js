var Point  = Isomer.Point;
var Path   = Isomer.Path;
var Shape  = Isomer.Shape;
var Vector = Isomer.Vector;
var Color  = Isomer.Color;


var red = new Color(160, 60, 50);
var blue = new Color(50, 60, 160);
var white = new Color(200, 200, 200);
var colCap = new Color(163, 69, 64);

var ColorBuildHighlight = new Color(122, 204, 147);
var ColorBuilding = new Color(209, 207, 199);

var BS = 3; // size of the blocks
var BH = .5; // base height and height of the floor
var BlockHeight = .5 // height of each layer
var Grid = 5;


// Layers to draw on
var drawLayers = Array.from(Array(Grid+Grid), () => new Array())
var rotation = 0;
var zoomfactor = 5;

function Rotate() {
    rotation++
    DrawBoard()
}

function ZoomIn() {
    zoomfactor--;
}

function ZoomOut() {
    zoomfactor++;
}
// Convert an X and Y into the actual X/Y accounting for rotation
function withRotation(x, y) {
    if (rotation%4 == 0) {
        return [x, y]
    } else if (rotation%4 == 1){
        return [y, Grid-x-1]
    } else if (rotation%4 == 2){
        return [Grid-x-1, Grid-y-1]
    } else if (rotation%4 == 3){
        return [Grid-y-1, x]
    }
}

function drawtile(ox, oy, h, worker, highlight, ghostLayer) {
    ncoords = withRotation(ox, oy);
    x = ncoords[0];
    y = ncoords[1];
    let layer = drawLayers[x+y]
    if (h > 0) {
        layer.push([Shape.Prism(Point(x*BS,y*BS,BH), BS, BS, BlockHeight), ColorBuilding])
    }
    if (h > 1) {
        layer.push([Shape.Prism(Point(x*BS+.25,y*BS+.25,BH+1*BlockHeight), BS-.5, BS-.5, BlockHeight), ColorBuilding])
    }
    if (h > 2) {
        layer.push([Shape.Prism(Point(x*BS+.5,y*BS+.5,BH+2*BlockHeight), BS-1, BS-1, BlockHeight), ColorBuilding])
    }
    if (h > 3) {
        layer.push([Shape.Pyramid(Point(x*BS+.5, y*BS+.5,BH+3*BlockHeight), BS-1, BS-1, 1), colCap])
    }
    if (highlight && h < 4) {
        var offset = Math.max(h * 0.25 - 0.25, 0);
         layer.push([new Path([
          Point(x*BS+offset, y*BS+offset, h*BlockHeight+BH),
          Point((x+1)*BS-offset, y*BS+offset, h*BlockHeight+BH),
          Point((x+1)*BS-offset, (y+1)*BS-offset, h*BlockHeight+BH),
          Point(x*BS+offset, (y+1)*BS-offset, h*BlockHeight+BH)
        ]), highlight]);
    }
    if (worker && h < 4) {
        layer.push([Shape.Prism(Point(x*BS+1, y*BS+1, BH+h*BlockHeight)), worker]);
    }
}

function DrawBoard(data) {
    drawLayers = Array.from(Array(Grid+Grid), () => new Array())
    drawtile(4,4,4,red)
    drawtile(3,3,4,red)
    drawtile(3,2,2,null, ColorBuildHighlight)
    drawtile(2,2,2,null, ColorBuildHighlight)
    drawtile(2,1,4,red)
    drawtile(0,0,4,red)
    Redraw()
}

function Redraw() {
    var ctx = canvas.getContext('2d');
    var width = zoomfactor*canvas.width;
    var height = zoomfactor*canvas.height;
    ctx.clearRect(0, 0, width, height);

    var iso = new Isomer(document.getElementById("canvas"));
    iso.add(Shape.Prism(Point.ORIGIN, Grid*BS, Grid*BS, BH));
    for (let i = drawLayers.length-1; i >=0; i--) {
        console.log(i, drawLayers[i].length)
        for (t in drawLayers[i]) {
            iso.add(...drawLayers[i][t])
        }
    }

    // TODO: Add labels at the given points for turn notation
    let p = iso._translatePoint(Point(0, 1))
    ctx.fillText("Hello world", p.x, p.y);
}


