package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joshwi/go-utils/utils"
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

	audits = map[string][]utils.Tag{
		"nfl": {
			{Name: "games", Value: "MATCH (n:games) RETURN n.label as label ORDER BY label"},
			{Name: "colors", Value: "MATCH (n:colors) RETURN n.label as label ORDER BY label"},
		},
	}
)

func init() {
	// Define flag arguments for the application
	flag.StringVar(&repo, `r`, `deltadb-backup`, `Specify config. Default: nfldb-backup`)
	flag.StringVar(&collection, `a`, `nfl`, `Specify collection audit. Default: nfl`)
	flag.Parse()
}

func main() {

	// Create application session with Neo4j
	// uri := "bolt://" + HOST + ":" + PORT
	// driver := graphdb.Connect(uri, USERNAME, PASSWORD)
	// sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	// session, err := driver.NewSession(sessionConfig)
	// if err != nil {
	// 	log.Println(err)
	// }

	// output := [][]string{}
	// changelog := map[string][]string{}
	// audit := audits[collection]
	directory := fmt.Sprintf("%v/%v/%v", DIRECTORY, repo, collection)

	err := utils.Unzip(directory+".zip", directory)

	if err != nil {
		log.Fatal(err)
	}

	// for _, item := range audit {
	// 	// log.Println(n, item.Name, item.Value)

	// 	old_text := utils.ReadTxt(directory + "/" + item.Name + ".txt")

	// 	diff_file := fmt.Sprintf("%v.txt", item.Name)

	// 	cypher_response := graphdb.RunCypher(session, item.Value)

	// 	new_text := utils.Strip(cypher_response)

	// 	old := strings.Split(old_text, "\n")
	// 	new := strings.Split(new_text, "\n")

	// 	diff := utils.Difference(old, new)

	// 	err = utils.Write(directory, diff_file, new_text, 0777)

	// 	if err != nil {
	// 		log.Println(err)
	// 	}

	// 	changelog[item.Name] = diff

	// }

	// output := utils.Rotate(changelog)
	// err = utils.WriteCsv(directory, "audit.csv", output, 0777)
	// if err != nil {
	// 	log.Println(err)
	// }

	// utils.Zip(directory, directory+".zip")

}
