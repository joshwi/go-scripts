package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joshwi/go-git/gitscm"
	"github.com/joshwi/go-utils/utils"
)

var (
	DIRECTORY = utils.Env("DIRECTORY")
	GIT_URL   = utils.Env("GIT_URL")
	GIT_USER  = utils.Env("GIT_USER")
	GIT_TOKEN = utils.Env("GIT_TOKEN")
	GIT_EMAIL = utils.Env("GIT_EMAIL")
)

func main() {

	// Init flag values
	var name string

	// Define flag arguments for the application
	flag.StringVar(&name, `r`, `nfldb-backup`, `Specify repository. Default: nfldb-backup`)
	flag.Parse()

	day := time.Now().Format("2006-01-02")
	directory := fmt.Sprintf("%v/%v", DIRECTORY, name)
	url := fmt.Sprintf("%v/%v/%v.git", GIT_URL, GIT_USER, name)

	log.Println(day, directory, url)

	project := gitscm.Project{
		Name:      name,
		Directory: directory,
		Url:       url,
		Token:     GIT_TOKEN,
		User:      GIT_USER,
		Email:     GIT_EMAIL,
	}

	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		project, err = project.Clone(directory)
	} else {
		project, err = project.Open(directory)
	}

	branches, err := project.Branches(name)

	log.Println(branches)

	err = project.Branch(name, day)

	branches, err = project.Branches(name)

	log.Println(branches)

	err = utils.Write(directory, "nfl.txt", "chiefs\nchargers\nraiders\nbroncos\n", 0777)

	if err != nil {
		log.Fatal(err)
	}

	err = project.Add()

	err = project.Commit(fmt.Sprintf("DB backup: %v", time.Now().Format(time.RFC3339)))

	err = project.Push()

}
