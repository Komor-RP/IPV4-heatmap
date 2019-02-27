let myMap;
let heat;
init();

function init() {
    myMap = L.map('mapid').setView([35.99, -78.89], 12);
    L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token={accessToken}', {
        maxZoom: 17,
        id: 'mapbox.streets',
        accessToken:    'pk.eyJ1IjoiZmFuaCIsImEiOiJjanNrb3JyOWIxN3dhNDRscDRncGthdjE3In0.HuVODwv3RaTzjLptnEDGYQ'
    }).addTo(myMap);

    myMap.on('zoom', getIPLocations);
    myMap.on('moveend', getIPLocations);

    getIPLocations();
}

function getIPLocations() {
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

    if (data) {
        data.forEach(function(point) {
            locations.push(Object.values(point));
        });
    }
    if (heat) {
        heat.remove();
    }
    heat = L.heatLayer(locations, {radius: 15}).addTo(myMap);
}

function getMapBounds() {
    let bounds = myMap.getBounds();

    return {
        north: bounds.getNorth(),
        south: bounds.getSouth(),
        west: bounds.getWest(),
        east: bounds.getEast()
    }
}

