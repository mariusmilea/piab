package receiver

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
	Created     time.Time         `json:"created" bson:"created"`
	Email       string            `json:"email" bson: "email"`
	Label       map[string]string `json:"label" bson:"label"`
	LastUpdated time.Time         `json:"lastUpdated" bson:"lastUpdated"`
}

// Defaults
const (
	uri                       = "/v1"
	collection                = "receivers"
	DefaultDatabase           = "piab"
	DefaultPrometheusServer   = "prometheus"
	DefaultPrometheusPort     = "9090"
	DefaultAlertmanagerServer = "alertmanager"
	DefaultAlertmanagerPort   = "9093"
)

// Interface
type Receiver interface {
	Register(r *mux.Router)
}

// Attributes
type receiver struct {
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
func New(session *mgo.Session, o Options) Receiver {
	self := receiver{
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
func (self *receiver) Register(r *mux.Router) {
	r.HandleFunc(uri+"/receivers/{id}", self.get).Methods("GET")
	r.HandleFunc(uri+"/receivers", self.all).Methods("GET")
	r.HandleFunc(uri+"/receivers", self.create).Methods("POST")
	r.HandleFunc(uri+"/receivers/{id}", self.update).Methods("PUT")
	r.HandleFunc(uri+"/receivers/{id}", self.delete).Methods("DELETE")
	r.HandleFunc(uri+"/receivers/generate", self.generate).Methods("POST")
}

// get all receivers from DB
func (self *receiver) GetReceivers() []Model {
	var result []Model
	// query all Ids
	err := self.session.DB(self.database).C(collection).Find(nil).All(&result)
	util.Check(err)

	return result
}

// get all receivers
func (self *receiver) all(w http.ResponseWriter, req *http.Request) {
	result := self.GetReceivers()

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

// get a single receiver ID
func (self *receiver) get(w http.ResponseWriter, req *http.Request) {
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

// create receiver
func (self *receiver) create(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	m := Model{
		Created:     t,
		LastUpdated: t,
	}

	// Populate the receiver data
	json.NewDecoder(req.Body).Decode(&m)

	// Add an Id
	m.Id = bson.NewObjectId()

	// Write to mongo
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

// update receiver
func (self *receiver) update(w http.ResponseWriter, req *http.Request) {
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

	// Update the receiver
	if err := self.session.DB(self.database).C(collection).UpdateId(oid, m); err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

// delete receiver
func (self *receiver) delete(w http.ResponseWriter, req *http.Request) {

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

	// Remove the receiver
	if err := self.session.DB(self.database).C(collection).RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

// generate receiver rules
func (self *receiver) generate(w http.ResponseWriter, req *http.Request) {
	receivers := self.GetReceivers()
	util.CreateConfigurationFile("/opt/config/alertmanager.yml", "templates/alertmanager.yml.tmpl", receivers)
	svc_api := "http://" + self.prometheusserver + ":" + self.prometheusport + "/-/reload"
	util.ReloadService(svc_api)
}
