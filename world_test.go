package alvasion_test

import (
	"sync"
	"testing"

	"github.com/EmilGeorgiev/alvasion"
	"github.com/stretchr/testify/assert"
)

// Tests for DestroyCity ----------------------------------
func TestDestroyCity(t *testing.T) {
	// SETUP
	outNorth := make(chan alvasion.Alien, 1)
	outSouth := make(chan alvasion.Alien, 1)
	cityFoo := alvasion.City{
		Name:               "Foo",
		OutgoingRoads:      []chan alvasion.Alien{outNorth, outSouth, nil, nil},
		IncomingRoads:      []chan alvasion.Alien{make(chan alvasion.Alien), make(chan alvasion.Alien), nil, nil},
		OutgoingRoadsNames: []string{"north=Baz", "south=Kart"},
		IsDestroyed:        false,
		Alien:              &alvasion.Alien{ID: 77},
	}

	// ACTIONS
	actual := cityFoo.Destroy()
	_, openedNorth := <-outNorth
	_, openedSouth := <-outSouth

	// ASSERTIONS
	expected := alvasion.City{
		Name:               "Foo",
		OutgoingRoads:      make([]chan alvasion.Alien, 4),
		IncomingRoads:      make([]chan alvasion.Alien, 4),
		IsDestroyed:        true,
		Alien:              nil,
		OutgoingRoadsNames: make([]string, 4),
	}
	assert.Equal(t, expected, actual)
	assert.False(t, openedNorth)
	assert.False(t, openedSouth)
}

func TestDestroyCityDestroyedCity(t *testing.T) {
	// SETUP
	c := alvasion.City{Name: "Foo", IsDestroyed: true}

	// ACTIONS
	actual := c.Destroy()

	// ASSERTIONS
	expected := alvasion.City{
		Name:        "Foo",
		IsDestroyed: true,
	}
	assert.Equal(t, expected, actual)
}

// Tests for CheckForIncomingAliens --------------------
func TestCheckForIncomingAliensWhenZeroAliensVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{Name: "Foo"}
	wg := sync.WaitGroup{}

	// ACTION
	c.CheckForIncomingAliens(&wg)
	wg.Wait()

	// ASSERTIONS
	// this test case assert that the method will not block forever if no aliens are coming.
}

func TestCheckForIncomingAliensWhenOneAlienVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name:          "Foo",
		IncomingRoads: []chan alvasion.Alien{make(chan alvasion.Alien, 1)},
	}
	wg := sync.WaitGroup{}
	a := alvasion.Alien{ID: 55, Sitreps: make(chan alvasion.Sitrep, 1)}
	c.IncomingRoads[0] <- a

	// ACTION
	c.CheckForIncomingAliens(&wg)
	actual := <-a.Sitreps
	wg.Wait()

	// ASSERTION
	expected := alvasion.Sitrep{
		FromAliens: []alvasion.Alien{a},
		CityName:   "Foo",
	}
	assert.Equal(t, expected, actual)
}

func TestCheckForIncomingAliensWhenTwoAliensVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name: "Baz",
		IncomingRoads: []chan alvasion.Alien{
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
		},
	}
	wg := sync.WaitGroup{}
	a55 := alvasion.Alien{ID: 55, Sitreps: make(chan alvasion.Sitrep, 1)}
	c.IncomingRoads[0] <- a55
	a100 := alvasion.Alien{ID: 100, Sitreps: make(chan alvasion.Sitrep, 1)}
	c.IncomingRoads[1] <- a100

	// ACTION
	c.CheckForIncomingAliens(&wg)
	actual := <-a55.Sitreps
	wg.Wait()

	// ASSERTION
	expected := alvasion.Sitrep{FromAliens: []alvasion.Alien{a55, a100}, CityName: "Baz"}
	assert.Equal(t, expected, actual)
}

func TestCheckForIncomingAliensWhenFourAliensVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name: "Baz",
		IncomingRoads: []chan alvasion.Alien{
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
		},
	}
	wg := sync.WaitGroup{}
	reports := make(chan alvasion.Sitrep, 1)
	a1 := alvasion.Alien{ID: 1, Sitreps: reports}
	c.IncomingRoads[0] <- a1
	a2 := alvasion.Alien{ID: 2, Sitreps: reports}
	c.IncomingRoads[1] <- a2
	a3 := alvasion.Alien{ID: 3, Sitreps: reports}
	c.IncomingRoads[2] <- a3
	a4 := alvasion.Alien{ID: 4, Sitreps: reports}
	c.IncomingRoads[3] <- a4

	// ACTION
	c.CheckForIncomingAliens(&wg)
	actual := <-reports
	wg.Wait()
	close(reports)
	// if the method CheckForIncomingAliens send more than one event, even when we close the channel first values
	// in the channel will be read and finally the default value end 'false'
	_, reportsOpened := <-reports

	// ASSERTION
	expected := alvasion.Sitrep{FromAliens: []alvasion.Alien{a1, a2, a3, a4}, CityName: "Baz"}
	assert.Equal(t, expected, actual)
	assert.False(t, reportsOpened)
}

func TestCheckForIncomingAliensWhenTwoAliensVisitTheCityThroughTheSameChannel(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name: "Baz",
		IncomingRoads: []chan alvasion.Alien{
			make(chan alvasion.Alien, 2),
		},
	}
	wg := sync.WaitGroup{}
	reports := make(chan alvasion.Sitrep, 2)
	a1 := alvasion.Alien{ID: 1, Sitreps: reports}
	c.IncomingRoads[0] <- a1
	// Here the Sitrep channel is nil. This means that the test will panic if the Alien 2
	// push his report to the channel.
	a2 := alvasion.Alien{ID: 2, Sitreps: nil}
	c.IncomingRoads[0] <- a2

	// ACTION
	c.CheckForIncomingAliens(&wg)
	actual := <-reports
	wg.Wait()
	close(reports)
	// if the method CheckForIncomingAliens send more than one event, even when we close the channel first values
	// in the channel will be read and finally the default value end 'false'
	_, isOpened := <-reports

	// ASSERTION
	expected := alvasion.Sitrep{FromAliens: []alvasion.Alien{a1}, CityName: "Baz"}
	assert.Equal(t, expected, actual)
	assert.False(t, isOpened)
}

// Tests for CheckForDestroyedRoads -------------------
func TestCheckForDestroyedRoadsWhenZeroRoadsAreDestroyed(t *testing.T) {
	// SETUP
	northOut := make(chan alvasion.Alien, 1)
	southOut := make(chan alvasion.Alien, 1)
	eastOut := make(chan alvasion.Alien, 1)
	westOut := make(chan alvasion.Alien, 1)
	northIn := make(chan alvasion.Alien, 1)
	southIn := make(chan alvasion.Alien, 1)
	eastIn := make(chan alvasion.Alien, 1)
	westIn := make(chan alvasion.Alien, 1)
	c := alvasion.City{
		Name:               "Foo",
		OutgoingRoads:      []chan alvasion.Alien{northOut, southOut, eastOut, westOut},
		IncomingRoads:      []chan alvasion.Alien{northIn, southIn, eastIn, westIn},
		OutgoingRoadsNames: []string{"north=X1", "south=X2", "east=X3", "west=X4"},
	}

	// ACTION
	actual := c.CheckForDestroyedRoads()

	// ASSERTIONS
	assert.Equal(t, c, actual)
}

func TestEvaluateRoadsDestructionWhenOneRoadsIsDestroyed(t *testing.T) {
	// SETUP
	northOut := make(chan alvasion.Alien, 1)
	southOut := make(chan alvasion.Alien, 1)
	eastOut := make(chan alvasion.Alien, 1)
	westOut := make(chan alvasion.Alien, 1)
	northIn := make(chan alvasion.Alien, 1)
	southIn := make(chan alvasion.Alien, 1)
	eastIn := make(chan alvasion.Alien, 1)
	westIn := make(chan alvasion.Alien, 1)
	c := alvasion.City{
		Name:               "Foo",
		OutgoingRoads:      []chan alvasion.Alien{northOut, southOut, eastOut, westOut},
		IncomingRoads:      []chan alvasion.Alien{northIn, southIn, eastIn, westIn},
		OutgoingRoadsNames: []string{"north=X1", "south=X2", "east=X3", "west=X4"},
	}

	// ACTION
	close(northIn)
	c.CheckForDestroyedRoads()

	// ASSERTIONS
	// prove that channels are not closed/destroyed
	c.OutgoingRoads[1] <- alvasion.Alien{}
	c.OutgoingRoads[2] <- alvasion.Alien{}
	c.OutgoingRoads[3] <- alvasion.Alien{}
	c.IncomingRoads[1] <- alvasion.Alien{}
	c.IncomingRoads[2] <- alvasion.Alien{}
	c.IncomingRoads[3] <- alvasion.Alien{}

	assert.Nil(t, c.IncomingRoads[0])
	assert.Nil(t, c.OutgoingRoads[0])
	assert.Equal(t, []string{"", "south=X2", "east=X3", "west=X4"}, c.OutgoingRoadsNames)
}

func TestEvaluateRoadsDestructionWhenAllRoadsAreDestroyed(t *testing.T) {
	northOut := make(chan alvasion.Alien, 1)
	southOut := make(chan alvasion.Alien, 1)
	eastOut := make(chan alvasion.Alien, 1)
	westOut := make(chan alvasion.Alien, 1)
	northIn := make(chan alvasion.Alien, 1)
	southIn := make(chan alvasion.Alien, 1)
	eastIn := make(chan alvasion.Alien, 1)
	westIn := make(chan alvasion.Alien, 1)
	c := alvasion.City{
		Name:               "Foo",
		OutgoingRoads:      []chan alvasion.Alien{northOut, southOut, eastOut, westOut},
		IncomingRoads:      []chan alvasion.Alien{northIn, southIn, eastIn, westIn},
		OutgoingRoadsNames: []string{"north=X1", "south=X2", "east=X3", "west=X4"},
	}

	// ACTION
	close(northIn)
	close(southIn)
	close(eastIn)
	close(westIn)
	c.CheckForDestroyedRoads()

	// ASSERTIONS
	expectedIncoming := make([]chan alvasion.Alien, 4)
	expectedOutgoing := make([]chan alvasion.Alien, 4)
	assert.Equal(t, expectedIncoming, c.IncomingRoads)
	assert.Equal(t, expectedOutgoing, c.OutgoingRoads)
	assert.Equal(t, []string{"", "", "", ""}, c.OutgoingRoadsNames)
}
