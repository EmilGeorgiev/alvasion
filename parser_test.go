package alvasion_test

import (
	"bufio"
	"log"
	"os"
	"testing"

	"github.com/EmilGeorgiev/alvasion"
	"github.com/stretchr/testify/assert"
)

func TestReadLines(t *testing.T) {
	// SetUp
	createFileWithLines("world-map.txt", []string{
		"Foo west=Baz east=Boo north=Zerty south=Hepp",
		"Baz east=Foo west=Nzas north=Lkert south=Jjer",
		"Nzas west=Jett east=Baz north=Poelk south=Xols",
		"Poelk west=Kass east=Zass north=Pass south=Nzas",
		"Kk west=Hh east=Ll north=Nn south=Pp",
	})
	lines := make(chan alvasion.Line)

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

	assert.Equal(t, alvasion.Line{Text: "Foo west=Baz east=Boo north=Zerty south=Hepp", Number: 1}, actualLine1)
	assert.Equal(t, alvasion.Line{Text: "Baz east=Foo west=Nzas north=Lkert south=Jjer", Number: 2}, actualLine2)
	assert.Equal(t, alvasion.Line{Text: "Nzas west=Jett east=Baz north=Poelk south=Xols", Number: 3}, actualLine3)
	assert.Equal(t, alvasion.Line{Text: "Poelk west=Kass east=Zass north=Pass south=Nzas", Number: 4}, actualLine4)
	assert.Equal(t, alvasion.Line{Text: "Kk west=Hh east=Ll north=Nn south=Pp", Number: 5}, actualLine5)
	assert.Empty(t, actual)
	assert.False(t, ok) // assert that the channel is closed
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
