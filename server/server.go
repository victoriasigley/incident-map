package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	webRootEnv = "INCIDENT_MAP_WEBROOT"
)

type IncidentAddress struct {
	AddressId       string `json:"address_id"`
	AddressLine1    string `json:"address_line1"`
	City            string `json:"city"`
	CommonPlaceName string `json:"common_place_name"`
	CrossStreet1    string `json:"cross_street1"`
	CrossStreet2    string `json:"cross_street2"`
	FirstDue        string `json:"first_due"`
	GeoHash         string `json:"geohast"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
	Name            string `json:"name"`
	Number          string `json:"number"`
	PostalCode      string `json:"postal_code"`
	PrefixDirection string `json:"prefix_direction"`
	ResponseZone    string `json:"response_zone"`
	State           string `json:"state"`
	SuffixDirection string `json:"suffix_direction"`
	Type            string `json:"type"`
}

type ApparatusExtendedData struct {
	EventDuration    int `json:"event_duration"`
	ResponseDuration int `json:"response_duration"`
	TravelDuration   int `json:"travel_duration"`
	TurnoutDuration  int `json:"turnout_duration"`
}

type IncidentStatusInfo struct {
	Geohash   string    `json:"geohash"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timstamp  time.Time `json:"timestamp"`
}

type IncidentUnitStatus struct {
	Acknowledge IncidentStatusInfo `json:"acknowledged"`
	Arrived     IncidentStatusInfo `json:"arrived"`
	Available   IncidentStatusInfo `json:"available"`
	Cleared     IncidentStatusInfo `json:"cleared"`
	Dispatched  IncidentStatusInfo `json:"dispatched"`
	Enroute     IncidentStatusInfo `json:"enroute"`
	Tilda       IncidentStatusInfo `json:"~"`
}

type IncidentApparatus struct {
	CarId        string                `json:"car_id"`
	ExtendedData ApparatusExtendedData `json:"extended_data"`
	Geohash      string                `json:"geohash"`
	Personnel    []string              `json:"personnel"`
	Shift        string                `json:"shift"`
	Station      string                `json:"station"`
	UnitId       string                `json:"unit_id"`
	UnitStatus   IncidentUnitStatus    `json:"unit_status"`
}

type DescriptionExtendedData struct {
	DispatchDuration int `json:"dispatch_duration"`
	EventDuration    int `json:"event_duration"`
	ResponseTime     int `json:"response_time"`
}

type IncidentDescription struct {
	Comments            string                  `json:"comments"`
	DayOfWeek           string                  `json:"day_of_week"`
	EventClosed         time.Time               `json:"event_closed"`
	EventId             string                  `json:"event_id"`
	EventOpened         time.Time               `json:"event_opened"`
	ExtendedData        DescriptionExtendedData `json:"extended_data"`
	FirstUnitArrived    time.Time               `json:"first_unit_arrived"`
	FirstUnitDispatched time.Time               `json:"first_unit_dispatched"`
	FirstUnitEnroute    time.Time               `json:"first_unit_enroute"`
	HourOfDay           int                     `json:"hour_of_day"`
	IncidentNumber      string                  `json:"incident_number"`
	LoiSearchComplete   time.Time               `json:"loi_search_complete"`
	Subtype             string                  `json:"subtype"`
	Type                string                  `json:"type"`
}

type IncidentFireDepartment struct {
	FdId        string `json:"fd_id"`
	FirecaresId string `json:"firecares_id"`
	Name        string `json:"name"`
	Shift       string `json:"shift"`
	State       string `json:"state"`
	Timezone    string `json:"timezone"`
}

type Incident struct {
	Address        IncidentAddress        `json:"address"`
	Apparatus      []IncidentApparatus    `json:"apparatus"`
	Description    IncidentDescription    `json:"description"`
	FireDepartment IncidentFireDepartment `json:"fire_department"`
	Version        string                 `json:"version"`
}

func incident(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	incident := new(Incident)
	jsonData, err := json.Marshal(incident)
	if err != nil {
		log.Printf("Error marshlling incident: %+v", err)
	}
	fmt.Fprintf(w, string(jsonData))
}

func main() {

	http.HandleFunc("/incident", incident)

	indexLoc := os.Getenv(webRootEnv)
	fs := http.FileServer(http.Dir(indexLoc))
	http.Handle("/", fs)

	log.Printf("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error running server: %+v", err)
	}
}
