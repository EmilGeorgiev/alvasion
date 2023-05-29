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
	go alvasion.ReadLines("world-map.txt", lines)

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
	cases := []struct {
		Name               string
		NorthIsNil         bool
		SouthIsNil         bool
		EastIsNil          bool
		WestIsNil          bool
		OutgoingRoadsNames []string
	}{
		{Name: "X1", NorthIsNil: true, SouthIsNil: false, EastIsNil: false, WestIsNil: true, OutgoingRoadsNames: []string{"", "south=X4", "east=X2", ""}},
		{Name: "X2", NorthIsNil: true, SouthIsNil: false, EastIsNil: false, WestIsNil: false, OutgoingRoadsNames: []string{"", "south=X5", "east=X3", "west=X1"}},
		{Name: "X3", NorthIsNil: true, SouthIsNil: false, EastIsNil: true, WestIsNil: false, OutgoingRoadsNames: []string{"", "south=X6", "", "west=X2"}},
		{Name: "X4", NorthIsNil: false, SouthIsNil: false, EastIsNil: false, WestIsNil: true, OutgoingRoadsNames: []string{"north=X1", "south=X7", "east=X5", ""}},
		{Name: "X5", NorthIsNil: false, SouthIsNil: false, EastIsNil: false, WestIsNil: false, OutgoingRoadsNames: []string{"north=X2", "south=X8", "east=X6", "west=X4"}},
		{Name: "X6", NorthIsNil: false, SouthIsNil: false, EastIsNil: true, WestIsNil: false, OutgoingRoadsNames: []string{"north=X3", "south=X9", "", "west=X5"}},
		{Name: "X7", NorthIsNil: false, SouthIsNil: true, EastIsNil: false, WestIsNil: true, OutgoingRoadsNames: []string{"north=X4", "", "east=X8", ""}},
		{Name: "X8", NorthIsNil: false, SouthIsNil: true, EastIsNil: false, WestIsNil: false, OutgoingRoadsNames: []string{"north=X5", "", "east=X9", "west=X7"}},
		{Name: "X9", NorthIsNil: false, SouthIsNil: true, EastIsNil: true, WestIsNil: false, OutgoingRoadsNames: []string{"north=X6", "", "", "west=X8"}},
	}
	// 0 (north), 1 (south), 2 (east), 3 (west)
	parts := make(chan []string)
	//var wm alvasion.WorldMap

	// ACTION
	go func() {
		parts <- []string{"X1", "east=X2", "south=X4"}
		parts <- []string{"X2", "east=X3", "west=X1", "south=X5"}
		parts <- []string{"X3", "west=X2", "south=X6"}
		parts <- []string{"X4", "east=X5", "north=X1", "south=X7"}
		parts <- []string{"X5", "west=X4", "east=X6", "north=X2", "south=X8"}
		parts <- []string{"X6", "west=X5", "north=X3", "south=X9"}
		parts <- []string{"X7", "east=X8", "north=X4"}
		parts <- []string{"X8", "west=X7", "east=X9", "north=X5"}
		parts <- []string{"X9", "west=X8", "north=X6"}
		close(parts)
	}()
	wm := alvasion.GenerateWorldMap(parts)

	// ASSERTIONS
	for _, c := range cases {
		t.Run("Assert "+c.Name, func(t *testing.T) {
			assert.Equal(t, c.Name, wm[c.Name].Name)
			assert.Equal(t, c.OutgoingRoadsNames, wm[c.Name].OutgoingRoadsNames)

			assert.Equal(t, c.NorthIsNil, wm[c.Name].OutgoingRoads[0] == nil)
			assert.Equal(t, c.NorthIsNil, wm[c.Name].IncomingRoads[0] == nil)

			assert.Equal(t, c.SouthIsNil, wm[c.Name].OutgoingRoads[1] == nil)
			assert.Equal(t, c.SouthIsNil, wm[c.Name].IncomingRoads[1] == nil)

			assert.Equal(t, c.EastIsNil, wm[c.Name].OutgoingRoads[2] == nil)
			assert.Equal(t, c.EastIsNil, wm[c.Name].IncomingRoads[2] == nil)

			assert.Equal(t, c.WestIsNil, wm[c.Name].OutgoingRoads[3] == nil)
			assert.Equal(t, c.WestIsNil, wm[c.Name].IncomingRoads[3] == nil)
		})
	}
	assert.Equal(t, 9, len(wm))
}

func TestGenerateWorldMapConnectProperlyCitiesOnDirectionNorthSouth(t *testing.T) {
	// SETUP

	// 0 (north), 1 (south), 2 (east), 3 (west)
	parts := make(chan []string)

	// ACTION
	go func() {
		parts <- []string{"X1", "south=X4"}
		parts <- []string{"X4", "north=X1"}
		close(parts)
	}()
	wm := alvasion.GenerateWorldMap(parts)
	expectedAlien1 := alvasion.Alien{ID: 1}
	wm["X1"].OutgoingRoads[1] <- expectedAlien1
	actualAlien1 := <-wm["X4"].IncomingRoads[0]

	expectedAlien2 := alvasion.Alien{ID: 2}
	wm["X4"].OutgoingRoads[0] <- expectedAlien2
	actualAlien2 := <-wm["X1"].IncomingRoads[1]

	// ASSERTIONS
	assert.Equal(t, expectedAlien1, actualAlien1)
	assert.Equal(t, expectedAlien2, actualAlien2)
}

func TestGenerateWorldMapConnectProperlyCitiesOnDirectionEastWest(t *testing.T) {
	// SETUP

	// 0 (north), 1 (south), 2 (east), 3 (west)
	parts := make(chan []string)

	// ACTION
	go func() {
		parts <- []string{"X1", "east=X2"}
		parts <- []string{"X2", "west=X1"}
		close(parts)
	}()
	wm := alvasion.GenerateWorldMap(parts)
	expectedAlien1 := alvasion.Alien{ID: 1}
	wm["X1"].OutgoingRoads[2] <- expectedAlien1
	actualAlien1 := <-wm["X2"].IncomingRoads[3]

	expectedAlien2 := alvasion.Alien{ID: 2}
	wm["X2"].OutgoingRoads[3] <- expectedAlien2
	actualAlien2 := <-wm["X1"].IncomingRoads[2]

	// ASSERTIONS
	assert.Equal(t, expectedAlien1, actualAlien1)
	assert.Equal(t, expectedAlien2, actualAlien2)
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
