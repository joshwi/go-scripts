package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joshwi/go-utils/utils"
)

var (
	DIRECTORY = os.Getenv("DIRECTORY")
	USERNAME  = os.Getenv("NEO4J_USERNAME")
	PASSWORD  = os.Getenv("NEO4J_PASSWORD")
	HOST      = os.Getenv("NEO4J_SERVICE_HOST")
	PORT      = os.Getenv("NEO4J_SERVICE_PORT")
)

var audits = map[string][]utils.Tag{
	"nfl": {
		{Name: "games", Value: "MATCH (n:games) RETURN n.label as label ORDER BY label"},
		{Name: "colors", Value: "MATCH (n:colors) RETURN n.label as label ORDER BY label"},
	},
}

func Unzip(src, dest string) error {

	response := fmt.Sprintf(`[ Function: Unzip ] [ Source: %v ] [ Destination: %v ] [ Status: Success ]`, src, dest)

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0777)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		log.Println(f.Name)
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		log.Println("ILIANA")

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			response = fmt.Sprintf(`[ Function: Unzip ] [ Source: %v ] [ Destination: %v ] [ Status: Failed ] [ Error: %v ]`, src, dest, err)
			log.Println(response)
			return err
		}
	}

	return nil
}

func main() {

	// Init flag values
	var repo string
	var collection string

	// Define flag arguments for the application
	flag.StringVar(&repo, `r`, `deltadb-backup`, `Specify config. Default: nfldb-backup`)
	flag.StringVar(&collection, `a`, `nfl`, `Specify collection audit. Default: nfl`)
	flag.Parse()

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

	err := Unzip(directory+".zip", directory)

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
