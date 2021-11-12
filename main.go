package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joshwi/go-plugins/graphdb"
	"github.com/joshwi/go-utils/utils"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

var (
	DIRECTORY = utils.Env("DIRECTORY")
	USERNAME  = utils.Env("NEO4J_USERNAME")
	PASSWORD  = utils.Env("NEO4J_PASSWORD")
	HOST      = utils.Env("NEO4J_SERVICE_HOST")
	PORT      = utils.Env("NEO4J_SERVICE_PORT")
)

//Read contents of a file
func ReadCsv(filename string) [][]string {

	/*
		Input:
			(filename) string - Path of file to read
		Output:
			map[string]interface{} - JSON structured output
	*/

	output := [][]string{}

	csv_file, _ := os.Open(filename)
	r := csv.NewReader(csv_file)
	record, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	output = record

	return output

}

//Write contents of a file
func WriteCsv(filepath string, filename string, data [][]string, mode int) error {

	/*
		Input:
			(filename) string - Path of file to read
		Output:
			map[string]interface{} - JSON structured output
	*/

	response := fmt.Sprintf(`[ Function: Write ] [ Directory: %v ] [ File: %v ] [ Status: Success ]`, filepath, filename)

	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		os.MkdirAll(filepath, os.FileMode(mode))
	}

	path := fmt.Sprintf("%v/%v", filepath, filename)

	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	writer := csv.NewWriter(f)

	err = writer.WriteAll(data)

	if err != nil {
		response = fmt.Sprintf(`[ Function: Write ] [ Directory: %v ] [ File: %v ] [ Status: Failed ] [ Error: %v ]`, filepath, filename, err)
		log.Println(response)
		return err
	}

	log.Println(response)

	return nil

}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func strip(input [][]utils.Tag) string {
	output := ``

	for _, item := range input {
		for _, elem := range item {
			output = output + elem.Value + "\n"
		}
	}

	return output
}

var audits = map[string][]utils.Tag{
	"nfl": {
		// {Name: "games", Value: "MATCH (n:games) RETURN n.label as label ORDER BY label"},
		{Name: "colors", Value: "MATCH (n:colors) RETURN n.label as label ORDER BY label"},
	},
}

func main() {

	// Init flag values
	var repo string
	var collection string

	// Define flag arguments for the application
	flag.StringVar(&repo, `r`, `nfldb-backup`, `Specify config. Default: nfldb-backup`)
	flag.StringVar(&collection, `a`, `nfl`, `Specify collection audit. Default: nfl`)
	flag.Parse()

	// Create application session with Neo4j
	uri := "bolt://" + HOST + ":" + PORT
	driver := graphdb.Connect(uri, USERNAME, PASSWORD)
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	session, err := driver.NewSession(sessionConfig)
	if err != nil {
		log.Println(err)
	}

	output := [][]string{}
	directory := fmt.Sprintf("%v/%v/%v", DIRECTORY, repo, collection)
	audit := audits[collection]

	for _, item := range audit {
		// log.Println(n, item.Name, item.Value)

		old_text := utils.ReadTxt(directory + "/" + item.Name + ".txt")

		cypher_response := graphdb.RunCypher(session, item.Value)

		new_text := strip(cypher_response)

		old := strings.Split(old_text, "\n")
		new := strings.Split(new_text, "\n")

		diff := difference(old, new)
		// log.Println(diff)

		if len(output) == 0 {
			output = append(output, []string{item.Name})
		} else {
			output[0] = append(output[0], item.Name)
		}

		for m, elem := range diff {
			if len(output)-1 <= m {
				// log.Println(elem)
				output = append(output, []string{elem})
			}
		}

	}

	log.Println(output)

	// test := ReadCsv(directory + "/colors.csv")
	// log.Println(test)

	// err := WriteCsv(directory, "test.csv", [][]string{{"test", "test"}}, 0777)
	// if err != nil {
	// 	log.Println(err)
	// }

	// for n := range diff {
	// 	if len(output)+1 > n {
	// 		log.Println(8)
	// 		if output[n] == nil {
	// 			log.Println(n, output[n])
	// 		}
	// 	}
	// }

	// err = utils.Write(directory, "games.txt", new_text, 0777)

	// if err != nil {
	// 	log.Println(err)
	// }

}
