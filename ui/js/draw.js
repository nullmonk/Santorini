var Point  = Isomer.Point;
var Path   = Isomer.Path;
var Shape  = Isomer.Shape;
var Vector = Isomer.Vector;
var Color  = Isomer.Color;


var red = new Color(160, 60, 50);
var blue = new Color(50, 60, 160);
var white = new Color(200, 200, 200);
var colCap = new Color(163, 69, 64);

var teamColors = [null, blue, red]

var ColorBuildHighlight = new Color(122, 204, 147);
var ColorBuilding = new Color(209, 207, 199);

Isomer.prototype._reverseTranslatePoint = function(point) {
    var x = point.x - this.originX;
    var y = point.y - this.originY;
    // Solve for original coordinates
    var originalX = (x + y) / (2 * this.transformation[0][0]);
    var originalY = (x - y) / (2 * this.transformation[1][0]);
    var originalZ = -y / this.scale;
    return new Point(originalX, originalY, originalZ);
} 

var SantoriniGame = {
    blockSize: 3, // size of the blocks
    baseHeight: .5, // base height and height of the floor
    blockHeight: 0.5, // height of each layer

    boardHeight: 0,
    boardWidth: 0,
    rotation: 0,
    zoomFactor: 5,

    // Game state stuff
    drawLayers: [], // Game layers
    currentHash: "", // The current board hash


    // Rotate the game board
    rotate: function() {
        this.rotation++
        console.log("rotating")
        this.drawBoard()
    },

    zoomIn: function() {
        this.zoomFactor--;
        this.drawBoard();
    },

    zoomOut: function() {
        this.zoomFactor--;
        this.drawBoard();
    },

    onmouse: function(e) {
        var rect = canvas.getBoundingClientRect();
        var x = e.clientX - rect.left;
        var y = e.clientY - rect.top;
        var isoCoords = {
            x: (2 * y + x) / 2,
            y: (2 * y - x) / 2,
        }
        console.log(isoCoords, x, y);
    },

    draw: function() {
        var canvas = document.getElementById("canvas")
        var ctx = canvas.getContext('2d');

        var width = this.zoomFactor*canvas.width;
        var height = this.zoomFactor*canvas.height;
        ctx.clearRect(0, 0, width, height);
    
        var iso = new Isomer(canvas);
        iso.add(Shape.Prism(Point.ORIGIN, this.boardWidth*this.blockSize, this.boardHeight*this.blockSize, this.baseHeight));
        for (let i = this.drawLayers.length-1; i >=0; i--) {
            for (t in this.drawLayers[i]) {
                iso.add(...this.drawLayers[i][t])
            }
        }
    
        // TODO: Add labels at the given points for turn notation
        let p = iso._translatePoint(Point(0, 1))
        ctx.fillText("Hello world", p.x, p.y);
    },

    drawBoard: function() {
        if (this.currentHash !== null) {
            drawLayers = Array.from(Array(this.boardWidth+this.boardHeight), () => new Array())
            for (i = 0; i < this.boardWidth*this.boardHeight; i++) {
                x = Math.floor(i/this.boardWidth);
                y = i%this.boardWidth;
    
                t = this.currentHash.charCodeAt(i+1) - 65;
                team = t>>3;
                height = t & 0x7;
                console.log(`Tile at ${t} (${x}, ${y}) = team=${team}, height=${height} ${this.currentHash[i+1]}`)
                this.drawTile(x, y, height, teamColors[team])
            }
        }
        this.draw();
    },

    drawTile: function(ox, oy, h, worker, highlight, ghostlayer) {
        ncoords = this.withRotation(ox, oy);
        x = ncoords[0];
        y = ncoords[1];
        let layer = this.drawLayers[x+y]
        if (h > 0) {
            layer.push([Shape.Prism(Point(x*this.blockSize,y*this.blockSize,this.baseHeight), this.blockSize, this.blockSize, this.blockHeight), ColorBuilding])
        }
        if (h > 1) {
            layer.push([Shape.Prism(Point(x*this.blockSize+.25,y*this.blockSize+.25,this.baseHeight+1*this.blockHeight), this.blockSize-.5, this.blockSize-.5, this.blockHeight), ColorBuilding])
        }
        if (h > 2) {
            layer.push([Shape.Prism(Point(x*this.blockSize+.5,y*this.blockSize+.5,this.baseHeight+2*this.blockHeight), this.blockSize-1, this.blockSize-1, this.blockHeight), ColorBuilding])
        }
        if (h > 3) {
            layer.push([Shape.Pyramid(Point(x*this.blockSize+.5, y*this.blockSize+.5,this.baseHeight+3*this.blockHeight), this.blockSize-1, this.blockSize-1, 1), colCap])
        }
        if (highlight && h < 4) {
            var offset = Math.max(h * 0.25 - 0.25, 0);
            layer.push([new Path([
            Point(x*this.blockSize+offset, y*this.blockSize+offset, h*this.blockHeight+this.baseHeight),
            Point((x+1)*this.blockSize-offset, y*this.blockSize+offset, h*this.blockHeight+this.baseHeight),
            Point((x+1)*this.blockSize-offset, (y+1)*this.blockSize-offset, h*this.blockHeight+this.baseHeight),
            Point(x*this.blockSize+offset, (y+1)*this.blockSize-offset, h*this.blockHeight+this.baseHeight)
            ]), highlight]);
        }
        if (worker && h < 4) {
            layer.push([Shape.Prism(Point(x*this.blockSize+1, y*this.blockSize+1, this.baseHeight+h*this.blockHeight)), worker]);
        }
    },

    // Convert an X and Y into the actual X/Y accounting for rotation
    withRotation: function(x, y) {
        if (this.rotation%4 == 0) {
            return [x, y]
        } else if (this.rotation%4 == 1){
            return [y, Grid-x-1]
        } else if (this.rotation%4 == 2){
            return [Grid-x-1, Grid-y-1]
        } else if (this.rotation%4 == 3){
            return [Grid-y-1, x]
        }
    },

    init: function(hash) {
        if (!hash) {
            hash = "CAAAAAAAIAAAQAQAAAIAAAAAAA"; // standard 2 player start
        }
        this.currentHash = hash
        this.boardWidth = hash.charCodeAt(0) - 65 + 3; // Min width is 0
        this.boardHeight = (hash.length -1)/this.boardWidth;
        this.drawLayers = Array.from(Array(this.boardHeight + this.boardWidth), () => new Array())
        this.drawBoard();
    },

}
