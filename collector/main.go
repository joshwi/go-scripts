package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/joshwi/go-plugins/graphdb"
	"github.com/joshwi/go-utils/parser"
	"github.com/joshwi/go-utils/utils"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {

	// Init flag values
	var query string
	var name string

	// Define flag arguments for the application
	flag.StringVar(&query, `q`, ``, `Specify config. Default: <empty>`)
	flag.StringVar(&name, `c`, `pfr_team_season`, `Specify config. Default: pfr_team_season`)
	flag.Parse()

	// Pull in env variables: username, password, uri
	username := os.Getenv("NEO4J_USERNAME")
	password := os.Getenv("NEO4J_PASSWORD")
	host := os.Getenv("NEO4J_SERVICE_HOST")
	port := os.Getenv("NEO4J_SERVICE_PORT")

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

	for _, item := range parser.CONFIG_LIST {
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

	log.Println("COLLECTION - START")

	var wg sync.WaitGroup

	for _, entry := range inputs {

		wg.Add(1)

		go graphdb.RunScript(driver, entry, config, &wg)

	}

	wg.Wait()

	log.Println("COLLECTION - DONE")

	// session.Close()

}
