package main

import (
	"errors"
	"fmt"
	"github.com/EmilGeorgiev/alvasion/newversion"
	"log"
	"os"
	"sync"

	"github.com/EmilGeorgiev/alvasion/app"
	"gopkg.in/yaml.v3"
)

type Config struct {
	WorldMap          string `yaml:"world_map"`
	ValidationWorkers int    `yaml:"validation_workers"`
	NumberOfAliens    int    `yaml:"number_of_aliens"`
}

func main() {

	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %s\n", err)
	}

	// Unmarshal YAML to the Config struct
	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Unable to unmarshal data: %s\n", err)
	}

	log.Println("Generating World Map.")
	wm, err := generateWorldMap(config)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Println("worldMap is generated.")

	log.Printf("Initialize %d number of aliens/soldiers.\n", config.NumberOfAliens)
	aliens := make([]*newversion.Alien, config.NumberOfAliens)
	for i := 0; i < config.NumberOfAliens; i++ {
		aliens[i] = &newversion.Alien{
			Name: fmt.Sprintf("%d", i),
		}
	}

	wm

	//log.Println("Initialize AlienCommander.")
	//ac := app.NewAlienCommander(wm, aliens, sitrep)
	//
	//log.Println("Start the invasion!")
	//ac.StartInvasion()

	//log.Println("Generate the report")
	//report := ac.GenerateReportForInvasion()

	//log.Println("Store the report in a file report.txt")

	//f, err := os.OpenFile("report.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	//if err != nil {
	//	log.Fatalf("os.OpenFile error: %v", err)
	//}
	//defer f.Close()

	//_, err = io.WriteString(f, report)
	//if err != nil {
	//	log.Fatalf("io.WriteString error: %v", err)
	//}
	log.Println("Finish")
}

func generateWorldMap(config Config) (map[string]newversion.City, error) {
	lines := make(chan app.Line, 1000)
	go app.ReadLines(config.WorldMap, lines)

	errCh := make(chan error)
	partsOfLine := make(chan []string, 1000)
	wg := sync.WaitGroup{}
	for i := 0; i < config.ValidationWorkers; i++ {
		wg.Add(1)
		go func() {
			app.ValidateLines(lines, partsOfLine, errCh)
			wg.Done()
		}()
	}

	var hasErr bool
	go func() {
		for e := range errCh {
			hasErr = true
			fmt.Println(e)
		}
	}()

	go func() {
		// wait until all validation workers finish their work and put
		wg.Wait()
		close(partsOfLine)
	}()
	wm := app.GenerateWorldMap(partsOfLine)

	if hasErr {
		return nil, errors.New("there are errors during parsing the file that contain cities and their outgoing roads")
	}

	return wm, nil
}
