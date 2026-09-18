package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type Location struct {
	ID      int     `json:"id"`
	City    string  `json:"city"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	IP      string  `json:"ip"`
	Status  string  `json:"status"`
}

func locationsHandler(w http.ResponseWriter, r *http.Request) {
	cities := []struct {
		City, Country string
		Lat, Lng      float64
	}{
		{"New York", "USA", 40.7128, -74.0060},
		{"London", "UK", 51.5074, -0.1278},
		{"Tokyo", "Japan", 35.6762, 139.6503},
		{"Berlin", "Germany", 52.5200, 13.4050},
		{"Sydney", "Australia", -33.8688, 151.2093},
		{"Moscow", "Russia", 55.7558, 37.6173},
		{"São Paulo", "Brazil", -23.5505, -46.6333},
		{"Cairo", "Egypt", 30.0444, 31.2357},
		{"Singapore", "Singapore", 1.3521, 103.8198},
		{"Reykjavik", "Iceland", 64.1466, -21.9426},
	}

	var locs []Location
	for i, c := range cities {
		if rand.Float32() > 0.2 {
			locs = append(locs, Location{
				ID:      i + 1,
				City:    c.City,
				Country: c.Country,
				Lat:     c.Lat + (rand.Float64()-0.5)*0.4,
				Lng:     c.Lng + (rand.Float64()-0.5)*0.4,
				IP:      fmt.Sprintf("%d.%d.%d.%d", rand.Intn(180)+10, rand.Intn(255), rand.Intn(255), rand.Intn(255)),
				Status:  "ACTIVE_NODE",
			})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locs)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Amphiptere OS - Tactical Global Monitor</title>
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
    <style>
        body { margin: 0; background: #050811; color: #00ff66; font-family: monospace; }
        #header { padding: 15px 20px; background: #0b0f19; border-bottom: 1px solid #1f2937; display: flex; justify-content: space-between; align-items: center; }
        .brand { font-weight: bold; font-size: 16px; color: #ffffff; text-shadow: 0 0 8px #00ff66; }
        #map { height: calc(100vh - 60px); width: 100%; background: #050811; }
        .leaflet-layer, .leaflet-control-zoom-in, .leaflet-control-zoom-out, .leaflet-attribution-flag { 
            filter: invert(100%) hue-rotate(180deg) brightness(90%) contrast(95%); 
        }
    </style>
</head>
<body>
    <div id="header">
        <div class="brand">Amphiptere OS <span style="color: #00ff66; font-weight: normal;">// Global Tactical Monitor</span></div>
        <span id="status">STATUS: STREAMING LIVE NODES</span>
    </div>
    <div id="map"></div>
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <script>
        const map = L.map('map', { zoomControl: false }).setView([20, 10], 2);
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', { maxZoom: 18 }).addTo(map);
        
        let markers = [];
        function updateMap() {
            fetch('/locations')
                .then(res => res.json())
                .then(data => {
                    markers.forEach(m => map.removeLayer(m));
                    markers = [];
                    data.forEach(loc => {
                        let marker = L.circleMarker([loc.lat, loc.lng], {
                            color: '#00ff66',
                            fillColor: '#00ff66',
                            fillOpacity: 0.8,
                            radius: 7
                        }).addTo(map).bindPopup('<b>Target:</b> ' + loc.City + ', ' + loc.Country + '<br><b>IP:</b> ' + loc.IP + '<br><b>Status:</b> ' + loc.Status);
                        markers.push(marker);
                    });
                });
        }
        setInterval(updateMap, 3000);
        updateMap();
    </script>
</body>
</html>`
	fmt.Fprint(w, html)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/locations", locationsHandler)
	fmt.Println("[+] Amphiptere OS Tactical Map Server running at: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
