package alert

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"

	"github.com/mariusmilea/piab/app/util"
)

// Model
type Model struct {
	Id          bson.ObjectId     `json:"id" bson:"_id"`
	Name        string            `json:"name" bson:"name"`
	Created     time.Time         `json:"created" bson:"created"`
	Expression  string            `json:"expression" bson:"expression"`
	Duration    string            `json:"duration" bson:"duration"`
	Label       map[string]string `json:"label" bson:"label"`
	LastUpdated time.Time         `json:"lastUpdated" bson:"lastUpdated"`
	Summary     string            `json:"summary" bson:"summary"`
	Description string            `json:"description" bson:"description"`
	Runbook     string            `json:"runbook" bson:"runbook"`
}

// Defaults
const (
	uri                       = "/v1"
	collection                = "alerts"
	DefaultDatabase           = "piab"
	DefaultPrometheusServer   = "prometheus"
	DefaultPrometheusPort     = "9090"
	DefaultAlertmanagerServer = "alertmanager"
	DefaultAlertmanagerPort   = "9093"
)

// Interface
type Alert interface {
	Register(r *mux.Router)
}

// Attributes
type alert struct {
	database           string
	prometheusserver   string
	prometheusport     string
	alertmanagerserver string
	alertmanagerport   string
	session            *mgo.Session
}

// Options
type Options struct {
	Database           string
	PrometheusServer   string
	PrometheusPort     string
	AlertmanagerServer string
	AlertmanagerPort   string
}

// Constructor
func New(session *mgo.Session, o Options) Alert {
	self := alert{
		session: session,
	}

	if o.Database != "" {
		self.database = o.Database
	} else {
		self.database = DefaultDatabase
	}

	if o.PrometheusServer != "" {
		self.prometheusserver = o.PrometheusServer
	} else {
		self.prometheusserver = DefaultPrometheusServer
	}

	if o.PrometheusPort != "" {
		self.prometheusport = o.PrometheusPort
	} else {
		self.prometheusport = DefaultPrometheusPort
	}

	if o.AlertmanagerServer != "" {
		self.alertmanagerserver = o.AlertmanagerServer
	} else {
		self.alertmanagerserver = DefaultAlertmanagerServer
	}

	if o.AlertmanagerPort != "" {
		self.alertmanagerport = o.AlertmanagerPort
	} else {
		self.alertmanagerport = DefaultAlertmanagerPort
	}

	return &self
}

// Register
func (self *alert) Register(r *mux.Router) {
	r.HandleFunc(uri+"/alerts/{id}", self.get).Methods("GET")
	r.HandleFunc(uri+"/alerts", self.all).Methods("GET")
	r.HandleFunc(uri+"/alerts", self.create).Methods("POST")
	r.HandleFunc(uri+"/alerts/{id}", self.update).Methods("PUT")
	r.HandleFunc(uri+"/alerts/{id}", self.delete).Methods("DELETE")
	r.HandleFunc(uri+"/alerts/generate", self.generate).Methods("POST")
}

// get all alerts from DB
func (self *alert) GetAlerts() []Model {
	var result []Model
	// query all Ids
	err := self.session.DB(self.database).C(collection).Find(nil).All(&result)
	util.Check(err)

	return result
}

// get all alerts
func (self *alert) all(w http.ResponseWriter, req *http.Request) {
	result := self.GetAlerts()

	// Marshal into JSON structure
	resultj, err := json.MarshalIndent(result, "", "  ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", resultj)
}

// get a single alert ID
func (self *alert) get(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	// get id
	id := params["id"]

	// Verify ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Grab id
	oid := bson.ObjectIdHex(id)

	m := Model{}

	// query an Id
	if err := self.session.DB(self.database).C(collection).FindId(oid).One(&m); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Marshal into JSON structure
	resultj, err := json.MarshalIndent(m, "", "  ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", resultj)
}

// create alert
func (self *alert) create(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	m := Model{
		Created:     t,
		LastUpdated: t,
	}

	// Populate the alert data
	json.NewDecoder(req.Body).Decode(&m)

	// Add an Id
	m.Id = bson.NewObjectId()

	// Write the user to mongo
	self.session.DB(self.database).C(collection).Insert(m)

	// Marshal into JSON structure
	resultj, err := json.MarshalIndent(m, "", "  ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", resultj)
}

// update alert
func (self *alert) update(w http.ResponseWriter, req *http.Request) {
	m := Model{}
	params := mux.Vars(req)
	// Get id
	id := params["id"]

	// Verify ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	json.NewDecoder(req.Body).Decode(&m)
	m.LastUpdated = time.Now()

	// Update the alert
	if err := self.session.DB(self.database).C(collection).UpdateId(oid, m); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

// delete alert
func (self *alert) delete(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	// Get id
	id := params["id"]

	// Verify ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove the alert
	if err := self.session.DB(self.database).C(collection).RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

// generate alert rules
func (self *alert) generate(w http.ResponseWriter, req *http.Request) {
	alerts := self.GetAlerts()
	util.CreateConfigurationFile("/opt/config/prometheus.rules", "templates/prometheus.rules.tmpl", alerts)
	svc_api := "http://" + self.alertmanagerserver + ":" + self.alertmanagerport + "/-/reload"
	util.ReloadService(svc_api)
}
