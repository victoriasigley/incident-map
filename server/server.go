package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	webRootEnv = "INCIDENT_MAP_WEBROOT"

	// Postgres constants
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "postgres"
	// Postgres queries
	uploadIncidentQuery  = "SELECT * FROM incident_report.upload_report($1, $2, $3, $4, $5, $6);"
	getLastIncidentQuery = "SELECT address, apparatus, description, fire_department, version FROM incident_report.get_last_incident();"
)

type IncidentAddress struct {
	AddressId       string  `json:"address_id"`
	AddressLine1    string  `json:"address_line1"`
	City            string  `json:"city"`
	CommonPlaceName string  `json:"common_place_name"`
	CrossStreet1    string  `json:"cross_street1"`
	CrossStreet2    string  `json:"cross_street2"`
	FirstDue        string  `json:"first_due"`
	GeoHash         string  `json:"geohast"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Name            string  `json:"name"`
	Number          string  `json:"number"`
	PostalCode      string  `json:"postal_code"`
	PrefixDirection string  `json:"prefix_direction"`
	ResponseZone    string  `json:"response_zone"`
	State           string  `json:"state"`
	SuffixDirection string  `json:"suffix_direction"`
	Type            string  `json:"type"`
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

type Incident2 struct {
	Address        []byte `json:"address"`
	Apparatus      []byte `json:"apparatus"`
	Description    []byte `json:"description"`
	FireDepartment []byte `json:"fire_department"`
	Version        string `json:"version"`
}

func upload(w http.ResponseWriter, req *http.Request) {
	// Allow this function to be accessed when called from app on other port
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Parse the request form
	err := req.ParseForm()
	if err != nil {
		log.Printf("Error getting multipart form: %+v", err)
	}

	// Parse the file header
	_, fileHeader, err := req.FormFile("file")
	if err != nil {
		log.Printf("Error getting file: %+v", err)
	}

	// Open the file and defer closing it so we can read from it
	jsonFile, err := fileHeader.Open()
	defer jsonFile.Close()

	// Read the file
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Printf("Error reading file: %+v", err)
	}

	// Unmarshal file into incident struct
	var incident Incident
	json.Unmarshal(byteValue, &incident)

	jsonData, err := json.Marshal(incident)
	if err != nil {
		log.Printf("Error marshalling incident: %+v", err)
	}

	// Connect to postgres db
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Insert record to db
	address, _ := json.Marshal(incident.Address)
	apparatus, _ := json.Marshal(incident.Apparatus)
	description, _ := json.Marshal(incident.Description)
	fireDepartment, _ := json.Marshal(incident.FireDepartment)
	version, _ := json.Marshal(incident.Version)
	_, err = db.Query(uploadIncidentQuery, incident.Description.EventId, address, apparatus, description, fireDepartment, version)
	if err != nil {
		log.Printf("Error inserting incident in database: %+v", err)
	}

	fmt.Fprintf(w, string(jsonData))
}

func incident(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var incident Incident

	// Connect to postgres db
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Get latest incident data from the database
	var address []byte
	var apparatus []byte
	var description []byte
	var fireDepartment []byte
	var version []byte
	row := db.QueryRow(getLastIncidentQuery)
	row.Scan(&address, &apparatus, &description, &fireDepartment, &version)

	// Unmarshal return values into appropriate incident fields
	json.Unmarshal(address, &incident.Address)
	json.Unmarshal(apparatus, &incident.Apparatus)
	json.Unmarshal(description, &incident.Description)
	json.Unmarshal(fireDepartment, &incident.FireDepartment)
	json.Unmarshal(version, &incident.Version)

	jsonData, err := json.Marshal(incident)
	if err != nil {
		log.Printf("Error marshalling incident: %+v", err)
	}

	fmt.Fprintf(w, string(jsonData))
}

func main() {
	http.HandleFunc("/upload", upload)

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
