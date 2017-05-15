$(function () {
	var bit;
	var park;

	var interval;
	function countRefresh ( ) {
		interval = setInterval ( function ( ) { locateParking ( 1 ); }, 27000 );
	}

	function initInteract ( ) {
		bit = getLocation ( );
		$ ( "#going" ).click ( function ( ) {
			countRefresh ( );
			locateParking ( 0 );
			clearInterval ( interval );
		} );
		$ ( "#parking" ).click ( function ( ) {
			locateParking ( 1 );
			clearInterval ( interval );
		} );
	}

	function locateParking ( choice ) {
		park = choice;
		$.getJSON( "a/b1=" + bit.lat + "&b2=" + bit.long + "&c=" + choice + "&d=" + Math.floor ( Date.now ( ) / 1000 ), function( zeta ){
			for ( var i = 0; i < zeta.length; i++ ) {
				var delta = zeta [ i ];
				if ( i == 0 ) {
					bit = getLocation ( );
					if ( bit.error !== -1 ) {
						initMap ( bit, delta );
					}
				}
				marking(map, delta.lat, delta.long );
			}
		} );
	}

	function getLocation ( ) {
		if ( navigator.geolocation ) {
			return navigator.geolocation.getCurrentPosition ( setPoint );
		}
		return { error: -1 };
	}

	function setPoint ( position ) {
		return {lat:position.coords.latitude, long:position.coords.longitude};
	}

	function marking ( map, latitude, longitude ) {
		var marker = new google.maps.Marker( {
			map: map,
			position: {lat: latitude, long: longitude},
			title: 'Observed Vacant Parking'
		} );
	}

	$ ( document ).ready ( function ( ) {
		initInteract ( );
	} );
});