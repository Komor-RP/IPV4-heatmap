mapboxgl.accessToken = 'pk.eyJ1IjoiZmFuaCIsImEiOiJjanNrb3JyOWIxN3dhNDRscDRncGthdjE3In0.HuVODwv3RaTzjLptnEDGYQ';
var map = new mapboxgl.Map({
    container: 'map',
    style: 'mapbox://styles/mapbox/streets-v11',
    center: [-74.50, 40],
    zoom: 7
});

map.addControl(new mapboxgl.NavigationControl());



map.on("moveend", function() {
    drawHeatMap()
});




function drawHeatMap() {
    console.log("north: " + map.getBounds().getNorth());
    let bounds = map.getBounds();
    let north = bounds.getNorth();
    let south = bounds.getSouth();
    let west = bounds.getWest();
    let east = bounds.getEast()
    console.log(`north: ${north}, south: ${south}, west: ${west}, east: ${east}`);

    getIPLocations(bounds.getNorth(), bounds.getSouth(), bounds.getWest(), bounds.getEast());
}


function getIPLocations(north, south, west, east) {
    axios.get(
        '/api',
        {   
            params: {
                north: north,
                south: south,
                west: west,
                east: east
            },
            headers: {
                'Content-Type': 'application/json'
            }
        }
    ).then(function (response) {
    // handle success
        console.log(response);
    }).catch(function (error) {
    // handle error
        console.log(error);
    }).then(function () {
      // always executed
    });
}