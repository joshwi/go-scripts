package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/joshwi/go-plugins/graphdb"
	"github.com/joshwi/go-utils/logger"
	"github.com/joshwi/go-utils/parser"
	"github.com/joshwi/go-utils/utils"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

var (
	// Pull in env variables: username, password, uri
	username = os.Getenv("NEO4J_USERNAME")
	password = os.Getenv("NEO4J_PASSWORD")
	host     = os.Getenv("NEO4J_SERVICE_HOST")
	port     = os.Getenv("NEO4J_SERVICE_PORT")

	// Init flag values
	query    string
	name     string
	filename string
	logfile  string
)

func init() {

	// Define flag arguments for the application
	flag.StringVar(&query, `q`, ``, `Run query to DB for input parameters. Default: <empty>`)
	flag.StringVar(&name, `c`, `pfr_team_season`, `Specify config. Default: pfr_team_season`)
	flag.StringVar(&filename, `f`, ``, `Location of parsing config file. Default: <empty>`)
	flag.StringVar(&logfile, `l`, `./collection.log`, `Location of script logfile. Default: ./collection.log`)
	flag.Parse()

	// Initialize logfile at user given path. Default: ./collection.log
	logger.InitLog(logfile)

	logger.Logger.Info().Str("config", name).Str("query", query).Str("status", "start").Msg("COLLECTION")
}

func main() {

	// Open file with parsing configurations
	fileBytes, err := utils.Read("pfr.json")
	if err != nil {
		log.Println(err)
	}

	var CONFIG []utils.Config
	json.Unmarshal(fileBytes, &CONFIG)

	// Create application session with Neo4j
	uri := "bolt://" + host + ":" + port
	driver := graphdb.Connect(uri, username, password)
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	session, err := driver.NewSession(sessionConfig)
	if err != nil {
		log.Println(err)
	}

	// Find parsing config requested by user
	config := utils.Config{Name: "", Urls: []string{}, Params: []string{}, Parser: []utils.Parser{}}

	for _, item := range CONFIG {
		if name == item.Name {
			config = item
		}
	}

	// Compile parser config into regexp
	config.Parser = parser.Compile(config.Parser)

	// Grab input parameters from  Neo4j
	inputs := [][]utils.Tag{{utils.Tag{Name: "name", Value: config.Name}}}

	if len(query) > 0 {
		inputs = graphdb.RunCypher(session, query)
	}

	var wg sync.WaitGroup

	for _, entry := range inputs {

		wg.Add(1)

		go graphdb.RunScript(driver, entry, config, &wg)

	}

	wg.Wait()

	logger.Logger.Info().Str("config", name).Str("query", query).Str("status", "finish").Msg("COLLECTION")

	session.Close()

}
