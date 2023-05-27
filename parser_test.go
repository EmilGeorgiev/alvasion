package alvasion_test

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/EmilGeorgiev/alvasion"
	"github.com/stretchr/testify/assert"
)

func TestReadLines(t *testing.T) {
	// SETUP
	createFileWithLines("world-map.txt", []string{
		"Foo west=Baz east=Boo north=Zerty south=Hepp",
		"Baz east=Foo west=Nzas north=Lkert south=Jjer",
		"Nzas west=Jett east=Baz north=Poelk south=Xols",
		"Poelk west=Kass east=Zass north=Pass south=Nzas",
		"Kk west=Hh east=Ll north=Nn south=Pp",
	})
	lines := make(chan alvasion.Line)

	// ACTION
	go func() {
		if err := alvasion.ReadLines("world-map.txt", lines); err != nil {
			log.Fatal(err)
		}
	}()

	actualLine1 := <-lines
	actualLine2 := <-lines
	actualLine3 := <-lines
	actualLine4 := <-lines
	actualLine5 := <-lines
	actual, ok := <-lines // channel MUST be closed.

	// ASSERTION
	assert.Equal(t, alvasion.Line{Text: "Foo west=Baz east=Boo north=Zerty south=Hepp", Number: 1}, actualLine1)
	assert.Equal(t, alvasion.Line{Text: "Baz east=Foo west=Nzas north=Lkert south=Jjer", Number: 2}, actualLine2)
	assert.Equal(t, alvasion.Line{Text: "Nzas west=Jett east=Baz north=Poelk south=Xols", Number: 3}, actualLine3)
	assert.Equal(t, alvasion.Line{Text: "Poelk west=Kass east=Zass north=Pass south=Nzas", Number: 4}, actualLine4)
	assert.Equal(t, alvasion.Line{Text: "Kk west=Hh east=Ll north=Nn south=Pp", Number: 5}, actualLine5)
	assert.Empty(t, actual)
	assert.False(t, ok) // assert that the channel is closed
}

// Test cases for Validate Lines

func TestValidateCorrectLines(t *testing.T) {
	// SETUP
	lines := make(chan alvasion.Line)
	parts := make(chan []string, 3)
	errs := make(chan error)
	var actualErrs []error

	go func() {
		err := <-errs
		actualErrs = append(actualErrs, err)
	}()

	// ACTION
	go alvasion.ValidateLines(lines, parts, errs)

	lines <- alvasion.Line{Text: "Foo west=Baz east=Boo north=Zerty south=Hepp", Number: 1}
	lines <- alvasion.Line{Text: "Baz east=Foo west=Nzas north=Lkert", Number: 2}
	lines <- alvasion.Line{Text: "Nzas west=Jett", Number: 3}

	actualPartsForLine1 := <-parts
	actualPartsForLine2 := <-parts
	actualPartsForLine3 := <-parts

	// ASSERTION
	assert.Equal(t, actualPartsForLine1, []string{"Foo", "west=Baz", "east=Boo", "north=Zerty", "south=Hepp"})
	assert.Equal(t, actualPartsForLine2, []string{"Baz", "east=Foo", "west=Nzas", "north=Lkert"})
	assert.Equal(t, actualPartsForLine3, []string{"Nzas", "west=Jett"})
	assert.Nil(t, actualErrs)
}

func TestValidateLinesWithNoSpaces(t *testing.T) {
	// SETUP
	lines := make(chan alvasion.Line)
	parts := make(chan []string, 3)
	errs := make(chan error)
	done := make(chan bool)
	var actualErrs []error

	go func() {
		for {
			err, ok := <-errs
			if !ok {
				done <- true
				return
			}
			actualErrs = append(actualErrs, err)
		}
	}()

	// ACTION
	go alvasion.ValidateLines(lines, parts, errs)

	lines <- alvasion.Line{Text: "Foowest=Bazeast=Boonorth=Zertysouth=Hepp", Number: 1}
	lines <- alvasion.Line{Text: "Bazeast=Foowest=Nzasnorth=Lkert", Number: 2}
	lines <- alvasion.Line{Text: "Nzas west=Jett", Number: 3}

	actualPartsForLine3 := <-parts

	close(errs)
	<-done // be sure that the goroutine that read from the channel errs will read all errors.
	// ASSERTION
	expectedErrs := []error{
		fmt.Errorf("line number: %d has wrong format. A line should contains a city name and at least "+
			"one road that leading out of the city. Expect something like 'Foo west=Bar north=Baz' got: %s\n", 1, "Foowest=Bazeast=Boonorth=Zertysouth=Hepp"),

		fmt.Errorf("line number: %d has wrong format. A line should contains a city name and at least "+
			"one road that leading out of the city. Expect something like 'Foo west=Bar north=Baz' got: %s\n", 2, "Bazeast=Foowest=Nzasnorth=Lkert"),
	}
	assert.Equal(t, []string{"Nzas", "west=Jett"}, actualPartsForLine3)
	assert.Equal(t, expectedErrs, actualErrs)
}

func TestValidateLineContainsOnlyCityName(t *testing.T) {
	// SETUP
	lines := make(chan alvasion.Line)
	parts := make(chan []string, 3)
	errs := make(chan error)
	done := make(chan bool)
	var actualErrs []error

	go func() {
		for {
			err, ok := <-errs
			if !ok {
				done <- true
				return
			}
			actualErrs = append(actualErrs, err)
		}
	}()

	// ACTION
	go alvasion.ValidateLines(lines, parts, errs)

	lines <- alvasion.Line{Text: "Foo", Number: 1}
	lines <- alvasion.Line{Text: "Nzas west=Jett", Number: 2}

	actualPartsForLine2 := <-parts
	close(errs)
	<-done // be sure that the goroutine that reds from the channel errs will finish his work before assertion

	// ASSERTION
	expectedErrs := []error{
		fmt.Errorf("line number: %d has wrong format. A line should contains a city name and at least "+
			"one road that leading out of the city. Expect something like 'Foo west=Bar north=Baz' got: %s\n", 1, "Foo"),
	}
	assert.Equal(t, []string{"Nzas", "west=Jett"}, actualPartsForLine2)
	assert.Equal(t, expectedErrs, actualErrs)
}

func TestValidateLineContainsMoreThan5Parts(t *testing.T) {
	// SETUP
	lines := make(chan alvasion.Line)
	parts := make(chan []string, 3)
	errs := make(chan error)
	done := make(chan bool)
	var actualErrs []error

	go func() {
		for {
			err, ok := <-errs
			if !ok {
				done <- true
				return
			}
			actualErrs = append(actualErrs, err)
		}
	}()

	// ACTION
	go alvasion.ValidateLines(lines, parts, errs)

	// this line contains more than 4 roads leading out of the city
	lines <- alvasion.Line{Text: "Foo west=Baz east=Boo north=Zerty south=Hepp west=Kop", Number: 1}
	lines <- alvasion.Line{Text: "Nzas west=Jett", Number: 2}

	actualPartsForLine2 := <-parts
	close(errs)
	<-done // be sure that the goroutine that reds from the channel errs will finish his work before assertion

	// ASSERTION
	expectedErrs := []error{
		fmt.Errorf("line number: %d has wrong format. A line should contains a city name and maximum "+
			"4 road that leading out of the city. Expect something like 'Foo west=Bar north=Baz' got: %s\n", 1, "Foo west=Baz east=Boo north=Zerty south=Hepp west=Kop"),
	}
	assert.Equal(t, []string{"Nzas", "west=Jett"}, actualPartsForLine2)
	assert.Equal(t, expectedErrs, actualErrs)
}

func TestValidateLinesWithWrongRoadsFormat(t *testing.T) {
	// SETUP
	lines := make(chan alvasion.Line)
	parts := make(chan []string, 3)
	errs := make(chan error)
	done := make(chan bool)
	var actualErrs []error

	go func() {
		for {
			err, ok := <-errs
			if !ok {
				done <- true
				return
			}
			actualErrs = append(actualErrs, err)
		}
	}()

	// ACTION
	go alvasion.ValidateLines(lines, parts, errs)

	// this doesn't have sigh '=' in the road eastBoo
	lines <- alvasion.Line{Text: "Foo west=Baz eastBoo north=Zerty south=Hepp", Number: 1}
	// this doesn't have sigh '=' in the road northLkert
	lines <- alvasion.Line{Text: "Baz east=Foo west=Nzas northLkert", Number: 2}
	// this line is correct
	lines <- alvasion.Line{Text: "Nzas west=Jett", Number: 3}

	actualPartsForLine3 := <-parts

	close(errs)
	<-done // be sure that the goroutines that read from errs channel finished his work

	// ASSERTION
	expectedErrs := []error{
		fmt.Errorf("on line %d the road number %d has wrong format. Expected something like 'west=Baz' got %s", 1, 2, "eastBoo"),
		fmt.Errorf("on line %d the road number %d has wrong format. Expected something like 'west=Baz' got %s", 2, 3, "northLkert"),
	}

	assert.Equal(t, actualPartsForLine3, []string{"Nzas", "west=Jett"})
	assert.Equal(t, expectedErrs, actualErrs)
}

func TestValidateLinesWithWrongRoadDirection(t *testing.T) {
	// SETUP
	lines := make(chan alvasion.Line)
	parts := make(chan []string, 3)
	errs := make(chan error)
	done := make(chan bool)
	var actualErrs []error

	go func() {
		for {
			err, ok := <-errs
			if !ok {
				done <- true
				return
			}
			actualErrs = append(actualErrs, err)
		}
	}()

	// ACTION
	go alvasion.ValidateLines(lines, parts, errs)

	// in this line the second road has wrong direction 'eastt'. It should be 'east'
	lines <- alvasion.Line{Text: "Foo west=Baz eastt=Boo north=Zerty south=Hepp", Number: 1}
	// in this line the last road has wrong direction nor. It should be north
	lines <- alvasion.Line{Text: "Baz east=Foo west=Nzas nor=Lkert", Number: 2}
	// in this line the first road has wrong direction nor. It should be west
	lines <- alvasion.Line{Text: "Too westt=Baz east=Boo north=Zerty south=Hepp", Number: 3}
	//in this line the last road has wrong direction nor. It should be south
	lines <- alvasion.Line{Text: "Dooo west=Baz east=Boo north=Zerty outh=Hepp", Number: 4}
	// this line is correct
	lines <- alvasion.Line{Text: "Nzas west=Jett", Number: 5}

	actualPartsForLine5 := <-parts

	close(errs)
	<-done // be sure that the goroutines that read from errs channel finished his work

	// ASSERTION
	expectedErrs := []error{
		fmt.Errorf("on the line %d the road number %d has wrong direction. Expected 'west/north/east/south' got %s", 1, 2, "eastt"),
		fmt.Errorf("on the line %d the road number %d has wrong direction. Expected 'west/north/east/south' got %s", 2, 3, "nor"),
		fmt.Errorf("on the line %d the road number %d has wrong direction. Expected 'west/north/east/south' got %s", 3, 1, "westt"),
		fmt.Errorf("on the line %d the road number %d has wrong direction. Expected 'west/north/east/south' got %s", 4, 4, "outh"),
	}

	assert.Equal(t, actualPartsForLine5, []string{"Nzas", "west=Jett"})
	assert.Equal(t, expectedErrs, actualErrs)
}

// Test cases for Generate Word Map
func TestGenerateWorldMap(t *testing.T) {
	// SETUP
	parts := make(chan []string)
	done := make(chan struct{})
	var wm alvasion.WorldMap

	// ACTION
	go func() {
		parts <- []string{"Foo", "west=Baz", "east=Boo", "north=Zerty", "south=Hepp"}
		parts <- []string{"Baz", "east=Foo", "west=Nzas", "north=Lkert", "south=Jjer"}
		parts <- []string{"Nzas", "west=Jett", "east=Baz", "north=Poelk", "south=Xols"}
		parts <- []string{"Poelk", "west=Kass", "east=Zass", "north=Pass", "south=Nzas"}
		parts <- []string{"Kk", "west=Hh", "east=Ll", "north=Nn", "south=Pp"}
		close(done)
	}()

	wm = alvasion.GenerateWorldMap(parts, done)

	// ASSERTION
	_, okFoo := wm.Cities["Foo"]
	_, okBaz := wm.Cities["Baz"]
	_, okNzas := wm.Cities["Nzas"]
	_, okPoelk := wm.Cities["Poelk"]
	_, okKk := wm.Cities["Kk"]

	assert.Equal(t, 5, len(wm.Cities))
	assert.True(t, okFoo)
	assert.True(t, okBaz)
	assert.True(t, okNzas)
	assert.True(t, okPoelk)
	assert.True(t, okKk)
}

func createFileWithLines(fileName string, lines []string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}
	defer file.Close()

	// Use a buffered writer to write lines to the file
	writer := bufio.NewWriter(file)

	for _, line := range lines {
		if _, err = writer.WriteString(line + "\n"); err != nil {
			log.Fatalf("failed writing to file: %s", err)
		}
	}

	// Make sure all data is written to the underlying writer
	err = writer.Flush()
	if err != nil {
		log.Fatalf("failed flushing writer: %s", err)
	}
}
