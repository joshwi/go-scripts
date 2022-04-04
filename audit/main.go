package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joshwi/go-plugins/graphdb"
	"github.com/joshwi/go-utils/logger"
	"github.com/joshwi/go-utils/utils"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var (
	DIRECTORY = os.Getenv("DIRECTORY")
	USERNAME  = os.Getenv("NEO4J_USERNAME")
	PASSWORD  = os.Getenv("NEO4J_PASSWORD")
	HOST      = os.Getenv("NEO4J_SERVICE_HOST")
	PORT      = os.Getenv("NEO4J_SERVICE_PORT")

	// Init flag values
	repo       string
	collection string
	logfile    string

	audits = map[string][]utils.Tag{
		"nfl": {
			{Name: "conferences", Value: "MATCH (n:nfl_conferences) RETURN n.label as label ORDER BY label"},
			{Name: "divisions", Value: "MATCH (n:nfl_divisions) RETURN n.label as label ORDER BY label"},
			{Name: "seasons", Value: "MATCH (n:nfl_seasons) RETURN n.label as label ORDER BY label"},
			{Name: "teams", Value: "MATCH (n:nfl_teams) RETURN n.label as label ORDER BY label"},
			{Name: "games", Value: "MATCH (n:nfl_games) RETURN n.label as label ORDER BY label"},
			{Name: "colors", Value: "MATCH (n:nfl_colors) RETURN n.label as label ORDER BY label"},
		},
	}
)

func init() {
	// Define flag arguments for the application
	flag.StringVar(&repo, `r`, `deltadb-backup`, `Specify config. Default: nfldb-backup`)
	flag.StringVar(&collection, `a`, `nfl`, `Specify collection audit. Default: nfl`)
	flag.StringVar(&logfile, `l`, `./audit.log`, `Location of script logfile. Default: ./audit.log`)
	flag.Parse()

	// Initialize logfile at user given path. Default: ./collection.log
	logger.InitLog(logfile)

	logger.Logger.Info().Str("repo", repo).Str("collection", collection).Str("status", "start").Msg("AUDIT")
}

func main() {

	// Create application session with Neo4j
	uri := "bolt://" + HOST + ":" + PORT
	driver := graphdb.Connect(uri, USERNAME, PASSWORD)
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	session := driver.NewSession(sessionConfig)

	changelog := map[string][]string{}
	audit := audits[collection]
	directory := fmt.Sprintf("%v/%v/%v", DIRECTORY, repo, collection)

	// Unzip folder
	// err = utils.Unzip(directory+".zip", directory)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	for _, item := range audit {

		old_text, _ := utils.Read(directory + "/nodes/" + item.Name + ".txt")

		cypher_response := graphdb.RunCypher(session, item.Value)

		new_text := utils.Strip(cypher_response)

		old := strings.Split(string(old_text), "\n")
		new := strings.Split(new_text, "\n")

		diff := utils.Difference(new, old)

		err := utils.Write(directory+"/nodes/"+item.Name+".txt", []byte(new_text), 0755)
		if err != nil {
			log.Println(err)
		}

		changelog[item.Name] = diff

	}

	// Write changelog to audit.json file
	output, _ := json.Marshal(changelog)
	utils.Write(directory+"/audit.json", output, 0755)

	// Zip folder
	// err = utils.Zip(directory, directory+".zip")
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
