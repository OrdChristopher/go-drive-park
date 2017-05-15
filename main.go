package main

import (
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
	"html/template"
	"appengine"
	"net/http"
	"strconv"
	"time"
	"math"
	"sync"
	"fmt"
	"log"
	"os"
	
	// [START imports]
	"golang.org/x/net/context"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
	// [END imports]
)

const earthRadius = 6373
const marsRadius = 3389.5

const title = "Go Drive Park"
const limit = "-Date"
const maximum = 12
const facebook = "https://www.facebook.com/I.have.an.opulence.of.accounts"

type point struct {
	lat  float64
	long float64
}

type GoParkLocator struct {
	p1 point
	p2 point
	
	distance float64
	datetime time.Time
}

func main() {
	ctx := context.Background()
	// [START auth]
	proj := os.Getenv("go-drive-park")
	if proj == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}
	client, err := pubsub.NewClient(ctx, proj)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}
	// [END auth]
	http.HandleFunc("/", omega)
	http.HandleFunc("/a", alpha)
}

func omega(w http.ResponseWriter, r *http.Request) {
	if (r.URI.Path != "/") {
		http.NotFound(w,r)	
	}
	fmt.Fprint(w, goParkLocatorTemplate)
}

func alpha(writ http.ResponseWriter, read *http.Request) {
	var context = appengine.NewContext(read)
	
	beta1, err := strconv.ParseFloat(read.FormValue("b1"), 64)
	beta2, err := strconv.ParseFloat(read.FormValue("b2"), 64)
	var gamma = read.FormValue("c")
	var delta = read.FormValue("d")
	
	var epsilon = datastore.NewQuery(title).Ancestor(GoParkLocator(context)).Order(limit).Limit(maximum)

	var zeta = make([]GoParkLocator, 0, maximum)
	if _, err := epsilon.GetAll(context, &zeta); err != nil {
		http.Error(writ, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var location = point {beta1, beta2}
	
	var goParkResult = GoParkLocator{
		p1: location,
		p2: location,
		distance: 0,
	}
	
	for ii, iv:= range zeta {
		goParkResult.p2 = zeta[ii].p2
		goParkResult.distance = distance ( goParkResult.p1, goParkResult.p2 )
		goParkResult.datetime = zeta[ii].datetime
	}
	
	key := datastore.NewIncompleteKey(context, title, contextKey(context))
	dd, err := datastore.Put(context, key, &goParkResult)
	if err != nil {
		http.Error(writ, err.Error(), http.StatusInternalServerError)
		return
	}
	
	js, err := json.Marshal(goParkResult)
	if err != nil {
		http.Error(writ, err.Error(), http.StatusInternalServerError)
		return
	}

	writ.Header().Set("Content-Type", "application/json")
	writ.Write(js)
}
// [END func_alpha]

// distance calculation using the Spherical Law of Cosines.
func distance(p1, p2 point) float64 {
	s1, c1 := math.Sincos(rad(p1.lat))
	s2, c2 := math.Sincos(rad(p2.lat))
	clong := math.Cos(rad(p1.long - p2.long))
	return radius * math.Acos(s1*s2+c1*c2*clong)
}

// rad converts degrees to radians.
func rad(deg float64) float64 {
	return deg * math.Pi / 180
}

// contextKey returns the key used for all guestbook entries.
func contextKey(context appengine.Context ) *datastore.Key {
	return datastore.NewKey(context, title, title, 0, nil)
}

var goParkLocatorTemplate = template.Must(template.New("book").Parse('
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <title>Go, Drive &amp; Park</title>

    <!-- Bootstrap core CSS -->
    <link href="css/bootstrap.min.css" rel="stylesheet">

    <!-- Custom styles for this template -->
    <link href="css/album.css" rel="stylesheet">
    <link href="css/narrow-jumbotron.css" rel="stylesheet">
    <link rel="stylesheet" href="css/font-awesome.min.css">
  </head>

  <body>
    <div class="container">
      <div class="header clearfix">
        <nav>
          <ul class="nav nav-pills float-right">
            <li class="nav-item">
              <a class="nav-link active" href="index.html">Home <span class="sr-only">(current)</span></a>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="https://www.facebook.com/I.have.an.opulence.of.accounts">Facebook</a>
            </li>
          </ul>
        </nav>
        <h3 class="text-muted"><i class="fa fa-arrows fa-lg"></i><span> Go,</span> <i class="fa fa-car fa-lg"></i><span> Drive &amp;</span> <i class="fa fa-recycle fa-lg"></i><span> Park</span></h3>
      </div>
	 
      <div class="jumbotron">
	   <div id="choice">
		<a id="going" class="btn btn-lg btn-success" href="#go" role="button"><i class="fa fa-arrows fa-lg fa-spin fa-fw"></i> Go&bull;ing</a>
		<a id="parking" class="btn btn-lg btn-danger" href="#park" role="button">Park&bull;ing <i class="fa fa-recycle fa-lg fa-spin fa-fw"></i></a>
	   </div>
	   <div id="map">
		<script src="js/map.js"></script>
		<!-- 5zoom -->
		<iframe src="https://www.google.com/maps/embed?pb=!1m10!1m8!1m3!1d13256.033259810341!2d-117.90732429999998!3d33.837896199999996!3m2!1i1024!2i768!4f13.1!5e0!3m2!1sen!2sus!4v1494743676603" frameborder="0" style="border:0" allowfullscreen></iframe>
	   </div>
      </div>

      <div class="row marketing">
        <div class="col-lg-6">
          <h4>Go</h4>
          <p>Going is moving from one place or point to another; travel, leave & depart.</p>
		
		<h4>Drive</h4>
          <p>Drive is love towards the car for going and parking.</p>
		
          <h4>Park</h4>
          <p>Parking is to bring (a vehicle that one is driving) to a halt and leave it temporarily, typically in a parking lot or by the side of the road.</p>
        </div>

        <div class="col-lg-6">
          <h4>Go Lang</h4>
          <p>Developed with <a href="https://golang.org/">Go Lang</a></p>

          <h4>Advertisement</h4>
		<script async src="//pagead2.googlesyndication.com/pagead/js/adsbygoogle.js"></script>
		<!-- Go Drive Park -->
		<ins class="adsbygoogle"
			style="display:block"
			data-ad-client="ca-pub-6264938980367132"
			data-ad-slot="9448543605"
			data-ad-format="auto"></ins>
		<script>
		(adsbygoogle = window.adsbygoogle || []).push({});
		</script>
        </div>
      </div>

      <footer class="footer">
        <p>&copy; Mr. Rogers 2017</p>
      </footer>

    </div> <!-- /container -->

    <!-- Bootstrap core JavaScript
    ================================================== -->
    <!-- Placed at the end of the document so the pages load faster -->
    <script src="js/vendor/jquery.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/tether/1.4.0/js/tether.min.js" integrity="sha384-DztdAPBWPRXSA/3eYEEUWrWCy7G5KFbe8fFjk5JAIxUYHKkDx6Qin1DkWx51bBrb" crossorigin="anonymous"></script>
    <script src="js/main.js"></script>
    <script src="js/bootstrap.min.js"></script>
    <!-- IE10 viewport hack for Surface/desktop Windows 8 bug -->
    <script src="js/ie10-viewport-bug-workaround.js"></script>
  </body>
</html>
'))
