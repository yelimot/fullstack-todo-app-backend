package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/yelimot/fullstack-todo-app-backend/pkg/api"
	"github.com/yelimot/fullstack-todo-app-backend/pkg/app"
	"github.com/yelimot/fullstack-todo-app-backend/pkg/repository"
	"github.com/yelimot/fullstack-todo-app-backend/pkg/version"
	"gopkg.in/yaml.v2"
)

var (
	dbFileFlag     = flag.String("db", "db.json", "Db file.")
	configFileFlag = flag.String("config", "config.yml", "Path to the configuration file.")
	logFileFlag    = flag.String("log", "todo.log", "Path to the log file.")
	debugFlag      = flag.Bool("debug", false, "Show debug information.")
	versionFlag    = flag.Bool("version", false, "Show version information.")
)

func init() {
	// Parse command-line flags
	flag.Parse()

	// Log settings
	if *debugFlag {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.TraceLevel)
	} else {
		logrus.SetReportCaller(false)
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logFile, err := os.OpenFile(*logFileFlag, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logrus.WithError(err).Fatal("Could not open log file")
	}

	logrus.SetOutput(logFile)
}

func main() {

	// Show version information
	if *versionFlag {
		fmt.Fprintln(os.Stdout, version.Print("todo"))
		os.Exit(0)
	}

	// Load configuration file
	data, err := ioutil.ReadFile(*configFileFlag)
	if err != nil {
		logrus.WithError(err).Fatal("Could not load configuration")
	}
	var cfg api.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Could not load configuration")
	}

	dbFile, err := os.OpenFile(*dbFileFlag, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		logrus.WithError(err).Fatal("Could not open db file")
	}

	repo, err := repository.New(dbFile)
	if err != nil {
		logrus.WithError(err).Fatal("Could not create repository")
	}

	// Create a new todo app
	appInstance := app.New(repo)

	// Create a new api
	apiInstance, err := api.New(&cfg, appInstance)
	if err != nil {
		panic(err)
	}

	if err := apiInstance.Start(); err != nil {
		logrus.WithError(err).Fatal("Could not start api")
	}

}
