package main

import (
	"github.com/docopt/docopt-go"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"

	"github.com/mariusmilea/piab/app/alert"
	"github.com/mariusmilea/piab/app/receiver"
	"github.com/mariusmilea/piab/app/util"
)

func main() {
	usage := `app
Usage:
 app  [--bind-addr=<addr>] [--bind-port=<port>] [--mongo-server=<server>] [--mongo-port=<port>] [--mongo-database=<database>] [--prometheus-server=<server>] [--prometheus-port=<port>] [--alertmanager-server=<server>] [--alertmanager-port=<port>]
 app -h | --help
Options:
  -h --help                     	Show this screen.
  --bind-addr=<addr>            	Bind to address. [default: 0.0.0.0]
  --bind-port=<port>            	Bind to port. [default: 12345]
  --mongo-server=<server>       	MongoDB server. [default: mongo]
  --mongo-port=<port>           	MongoDB port. [default: 27017]
  --mongo-database=<datbase>    	MongoDB database. [default: piab]
  --prometheus-server=<server>  	Prometheus server. [default: prometheus]
  --prometheus-port=<port>  		Prometheus port. [default: 9090]
  --alertmanager-server=<server>	Alertmanager server. [default: alertmanager]
  --alertmanager-port=<port>		Alertmanager port. [default: 9093]
`
	// Parse arguments
	args, err := docopt.Parse(usage, nil, true, "piab 1.0", false)
	util.Check(err)

	// Connect to MongoDB
	session, err := mgo.Dial("mongodb://" + args["--mongo-server"].(string) + ":" + args["--mongo-port"].(string))
	util.Check(err)

	router := mux.NewRouter()

	alerts := alert.New(session, alert.Options{
		Database:           args["--mongo-database"].(string),
		PrometheusServer:   args["--prometheus-server"].(string),
		PrometheusPort:     args["--prometheus-port"].(string),
		AlertmanagerServer: args["--alertmanager-server"].(string),
		AlertmanagerPort:   args["--alertmanager-port"].(string),
	})
	alerts.Register(router)

	receivers := receiver.New(session, receiver.Options{
		Database:           args["--mongo-database"].(string),
		PrometheusServer:   args["--prometheus-server"].(string),
		PrometheusPort:     args["--prometheus-port"].(string),
		AlertmanagerServer: args["--alertmanager-server"].(string),
		AlertmanagerPort:   args["--alertmanager-port"].(string),
	})
	receivers.Register(router)

	log.Fatal(http.ListenAndServe(args["--bind-addr"].(string)+":"+args["--bind-port"].(string), router))
}
