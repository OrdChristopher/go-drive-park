var map;
function initMap(bit, delta) {
	map = new google.maps.Map(document.getElementById('map'), {
		center: bit, //eg; Anaheim
		scrollwheel: false,
		zoom: 12
	});

	var directionsDisplay = new google.maps.DirectionsRenderer({
		map: map
	});

	// Set destination, origin and travel mode.
	var request = {
		destination: delta, //eg; Go Here
		origin: bit, // eg; From Here
		travelMode: 'DRIVING'
	};

	// Pass the directions request to the directions service.
	var directionsService = new google.maps.DirectionsService();
	directionsService.route(request, function(response, status) {
		if (status == 'OK') {
			// Display the route on the map.
			directionsDisplay.setDirections(response);
		}
	});
}