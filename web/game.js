////// Here, we initialize the pixi application
var pixiRoot = new PIXI.Application({
    width: 1200,
    height: 800,
    backgroundColor: 0x6BACDE,
});

// add the renderer view element to the DOM
document.body.appendChild(pixiRoot.view);

////// Here, we create our traviso instance and add on top of pixi

// engine-instance configuration object
var instanceConfig =
{
    mapDataPath: "map.json", // the path to the json file that defines map data, required
    assetsToLoad: ["assets"], // array of paths to the assets that are desired to be loaded by traviso, no need to use if assets are already loaded to PIXI cache, default null

    engineInstanceReadyCallback: onEngineInstanceReady, // callback function that will be called once everything is loaded and engine instance is ready, default null
    tileSelectCallback: onTileSelect, // callback function that will be called when a tile is selected (call params will be the row and column indexes of the tile selected), default null
    objectSelectCallback: onObjectSelect, // callback function that will be called when a tile with an interactive map-object on it is selected (call param will be the object selected), default null
    objectReachedDestinationCallback: onObjectReachedDestination, // callback function that will be called when any moving object reaches its destination (call param will be the moving object itself), default null
    otherObjectsOnTheNextTileCallback: onOtherObjectsOnTheNextTile // callback function that will be called when any moving object is in move and there are other objects on the next tile (call params will be the moving object and an array of objects on the next tile), default null
};

var engine = TRAVISO.getEngineInstance(instanceConfig, { logEnabled: true });


// this method will be called when the engine is ready
function onEngineInstanceReady() {
    pixiRoot.stage.addChild(engine);
}

// this method will be called every time a tile is selected in the engine
function onTileSelect(rowIndex, columnIndex) {

}

// this method will be called when a tile with an interactive map-object on it is selected
function onObjectSelect(obj) {
    engine.moveCurrentControllableToLocation(obj.mapPos);
}

// this method will be called when any moving object reaches its destination
function onObjectReachedDestination(obj) {
    var objectsOnDestination = engine.getObjectsAtLocation(obj.mapPos);
    for (var i = 0; i < objectsOnDestination.length; i++) {
        // check if there is a flag on the destination tile
        if (objectsOnDestination[i].type === '10') {
            obj.changeVisual("flip", false, true, flipAnimFinished);
            break;
        }
    }
}

// this method will be called when the custom flip anim finished
function flipAnimFinished(obj) {
    // change the visual of the object so that it will face its last direction
    obj.changeVisualToDirection(obj.currentDirection, false);
}

// this method will be called when any moving object is in move and there are other objects on the next tile
function onOtherObjectsOnTheNextTile(movingObject, objectsOnNewTile) {
    var boxAnim;
    for (var i = 0; i < objectsOnNewTile.length; i++) {
        // check if there are boxes on the next tile
        if (objectsOnNewTile[i].type === '12') {
            boxAnim = createAndStartBoxAnim();
            engine.addCustomObjectToLocation(boxAnim, objectsOnNewTile[i].mapPos);
            engine.removeObjectFromLocation(objectsOnNewTile[i]);
        }
    }
}