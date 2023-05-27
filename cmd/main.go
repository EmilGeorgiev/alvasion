package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/EmilGeorgiev/alvasion"
	"gopkg.in/yaml.v3"
)

type Config struct {
	WorldMap          string `yaml:"world_map"`
	ValidationWorkers int    `yaml:"validation_workers"`
}

func main() {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %s\n", err)
	}

	// Unmarshal YAML to the Config struct
	var config Config
	if err = yaml.Unmarshal([]byte(data), &config); err != nil {
		log.Fatalf("Unable to unmarshal data: %s\n", err)
	}

	lines := make(chan alvasion.Line, 1000)
	go alvasion.ReadLines(config.WorldMap, lines)

	errCh := make(chan error)
	partsOfLine := make(chan []string, 1000)
	wg := sync.WaitGroup{}
	for i := 0; i < config.ValidationWorkers; i++ {
		wg.Add(1)
		go func() {
			alvasion.ValidateLines(lines, partsOfLine, errCh)
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
	wm := alvasion.GenerateWorldMap(partsOfLine)

	if hasErr {
		log.Fatal("there are errors during parsing the file that contain cities and their outgoing roads")
	}

	_ = wm
	fmt.Println("WorldMap is generated:")

	var number int
	fmt.Print("Enter a number of aliens: ")
	_, err = fmt.Scanf("%d", &number)
	if err != nil {
		log.Fatal("You entered an invalid number")
	}

}
