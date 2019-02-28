let myMap;
let heat;
init();

function init() {
    // set down map
    myMap = L.map('mapid').setView([35.99, -78.89], 12);
    L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token={accessToken}', {
        maxZoom: 17,
        id: 'mapbox.streets',
        accessToken:    'pk.eyJ1IjoiZmFuaCIsImEiOiJjanNrb3JyOWIxN3dhNDRscDRncGthdjE3In0.HuVODwv3RaTzjLptnEDGYQ'
    }).addTo(myMap);

    // add movement listeners
    myMap.on('zoom', drawHeatMap);
    myMap.on('moveend', drawHeatMap);

    // initiate first load
    drawHeatMap();
}

// sends a get request for coordinate information
// and plots the coordinates onto the map
function drawHeatMap() {
    let coordinates = getMapBounds();

    axios.get(
        '/api',
        {   
            params: {
                north: coordinates.north,
                south: coordinates.south,
                west: coordinates.west,
                east: coordinates.east
            },
            headers: {
                'Content-Type': 'application/json'
            }
        }
    ).then(function (response) {
        plotHeat(response.data);

    }).catch(function (error) {
        console.log(error);
    });
}

function plotHeat(data) {
    let locations = [];

    // transforms data to expected Leaflet-heat form
    // [latitude, longitude, intensity]
    if (data) {
        data.forEach(function(point) {
            locations.push(Object.values(point));
        });
    }

    // removes old heat layer so they don't stack
    if (heat) {
        heat.remove();
    }

    // plots heat layer using constructed data
    heat = L.heatLayer(locations, {
        radius: 15,
        minOpacity: .2,

        gradient: {
            0.2: 'blue',
            0.4: 'cyan',
            0.6: 'lime',
            0.8: 'yellow', 
            1: 'red'}
    }).addTo(myMap);
}

// gets the map boundaries and returns as an anonymous object
function getMapBounds() {
    let bounds = myMap.getBounds();

    return {
        north: bounds.getNorth(),
        south: bounds.getSouth(),
        west: bounds.getWest(),
        east: bounds.getEast()
    }
}

