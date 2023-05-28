package alvasion_test

import (
	"sync"
	"testing"

	"github.com/EmilGeorgiev/alvasion"
	"github.com/stretchr/testify/assert"
)

// Tests for CleanCity ----------------------------------
func TestCleanCity(t *testing.T) {
	// SETUP
	roadEast := make(chan alvasion.Alien)
	roadWest := make(chan alvasion.Alien)
	wm := &alvasion.WorldMap{Cities: map[string]*alvasion.City{
		"Foo": {Name: "Foo", Roads: []chan alvasion.Alien{roadEast, roadWest}, IsDestroyed: false},
		"Baz": {Name: "Baz", Roads: []chan alvasion.Alien{make(chan alvasion.Alien)}, IsDestroyed: false},
	}}

	// ACTIONS
	wm.CleanCity("Foo")

	// ASSERTIONS
	_, openedEast := <-roadEast
	_, openedWest := <-roadWest
	assert.True(t, wm.Cities["Foo"].IsDestroyed)
	assert.False(t, openedEast)
	assert.False(t, openedWest)
	assert.False(t, wm.Cities["Baz"].IsDestroyed)
}

func TestCleanCityThatDoesNotExist(t *testing.T) {
	// SETUP
	roadEast := make(chan alvasion.Alien, 1)
	roadWest := make(chan alvasion.Alien, 1)
	wm := &alvasion.WorldMap{Cities: map[string]*alvasion.City{
		"Foo": {Name: "Foo", Roads: []chan alvasion.Alien{roadEast, roadWest}, IsDestroyed: false},
	}}

	// ACTIONS
	wm.CleanCity("Not exist")
	roadEast <- alvasion.Alien{} // prove that channel is opened
	roadWest <- alvasion.Alien{} // prove that channel is opened

	// ASSERTIONS
	assert.False(t, wm.Cities["Foo"].IsDestroyed)
}

// Tests for EvaluateCityDestruction --------------------
func TestEvaluateCityDestructionWhenZeroAliensVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{Name: "Foo"}
	wg := sync.WaitGroup{}

	// ACTION
	alvasion.EvaluateCityDestruction(&c, &wg)
	wg.Wait()

	// ASSERTIONS
	expected := alvasion.City{Name: "Foo"}
	assert.Equal(t, expected, c)
}

func TestEvaluateCityDestructionWhenOneAlienVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name:  "Foo",
		Roads: []chan alvasion.Alien{make(chan alvasion.Alien, 1)},
	}
	wg := sync.WaitGroup{}
	a := alvasion.Alien{ID: 55, Sitreps: make(chan alvasion.Sitrep, 1)}
	c.Roads[0] <- a

	// ACTION
	alvasion.EvaluateCityDestruction(&c, &wg)
	actual := <-a.Sitreps
	wg.Wait()

	// ASSERTION
	expected := alvasion.Sitrep{
		From:            55,
		CityName:        "Foo",
		IsCityDestroyed: false,
	}
	assert.Equal(t, expected, actual)
	assert.Equal(t, 1, len(c.Roads))
}

func TestEvaluateCityDestructionWhenTwoAliensVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name: "Baz",
		Roads: []chan alvasion.Alien{
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
		},
	}
	wg := sync.WaitGroup{}
	a55 := alvasion.Alien{ID: 55, Sitreps: make(chan alvasion.Sitrep, 1)}
	c.Roads[0] <- a55
	a100 := alvasion.Alien{ID: 100, Sitreps: make(chan alvasion.Sitrep, 1)}
	c.Roads[1] <- a100

	// ACTION
	alvasion.EvaluateCityDestruction(&c, &wg)
	actualRep55 := <-a55.Sitreps
	actualRep100 := <-a100.Sitreps
	wg.Wait()

	// ASSERTION
	expectedRep55 := alvasion.Sitrep{From: 55, CityName: "Baz", IsCityDestroyed: true}
	expectedRep100 := alvasion.Sitrep{From: 100, CityName: "Baz", IsCityDestroyed: true}
	assert.Equal(t, expectedRep55, actualRep55)
	assert.Equal(t, expectedRep100, actualRep100)
}

func TestEvaluateCityDestructionWhenFourAliensVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name: "Baz",
		Roads: []chan alvasion.Alien{
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
		},
	}
	wg := sync.WaitGroup{}
	reports := make(chan alvasion.Sitrep, 1)
	a1 := alvasion.Alien{ID: 1, Sitreps: reports}
	c.Roads[0] <- a1
	a2 := alvasion.Alien{ID: 2, Sitreps: reports}
	c.Roads[1] <- a2
	a3 := alvasion.Alien{ID: 3, Sitreps: reports}
	c.Roads[2] <- a3
	a4 := alvasion.Alien{ID: 4, Sitreps: reports}
	c.Roads[3] <- a4

	// ACTION
	alvasion.EvaluateCityDestruction(&c, &wg)
	actualRep1 := <-reports
	actualRep2 := <-reports
	actualRep3 := <-reports
	actualRep4 := <-reports
	wg.Wait()

	// ASSERTION
	expectedRep1 := alvasion.Sitrep{From: 1, CityName: "Baz", IsCityDestroyed: true}
	expectedRep2 := alvasion.Sitrep{From: 2, CityName: "Baz", IsCityDestroyed: true}
	expectedRep3 := alvasion.Sitrep{From: 3, CityName: "Baz", IsCityDestroyed: true}
	expectedRep4 := alvasion.Sitrep{From: 4, CityName: "Baz", IsCityDestroyed: true}
	assert.Equal(t, expectedRep1, actualRep1)
	assert.Equal(t, expectedRep2, actualRep2)
	assert.Equal(t, expectedRep3, actualRep3)
	assert.Equal(t, expectedRep4, actualRep4)
}

func TestEvaluateCityDestructionWhenTwoAliensVisitTheCityThroughTheSameChannel(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name: "Baz",
		Roads: []chan alvasion.Alien{
			make(chan alvasion.Alien, 2),
		},
	}
	wg := sync.WaitGroup{}
	reports := make(chan alvasion.Sitrep, 2)
	a1 := alvasion.Alien{ID: 1, Sitreps: reports}
	c.Roads[0] <- a1
	// Here the Sitrep channel is nil. This means that the test will panic if the Alien 2
	// push his report to the channel.
	a2 := alvasion.Alien{ID: 2, Sitreps: nil}
	c.Roads[0] <- a2

	// ACTION
	alvasion.EvaluateCityDestruction(&c, &wg)
	actualRep1 := <-reports
	wg.Wait()

	// ASSERTION
	expectedRep1 := alvasion.Sitrep{From: 1, CityName: "Baz", IsCityDestroyed: false}
	assert.Equal(t, expectedRep1, actualRep1)
}

// Tests for EvaluateRoadsDestruction -------------------
func TestEvaluateRoadsDestructionWhenZeroRoadsAreDestroyed(t *testing.T) {
	// SETUP
	east := make(chan alvasion.Alien, 1)
	west := make(chan alvasion.Alien, 1)
	north := make(chan alvasion.Alien, 1)
	south := make(chan alvasion.Alien, 1)
	c := alvasion.City{Name: "Foo", Roads: []chan alvasion.Alien{east, west, north, south}}
	wg := sync.WaitGroup{}

	// ACTION
	alvasion.EvaluateRoadsDestruction(&c, &wg)
	wg.Wait()

	// ASSERTIONS
	east <- alvasion.Alien{}  // prove that channel est is not closed/destroyed
	west <- alvasion.Alien{}  // prove that channel est is not closed/destroyed
	north <- alvasion.Alien{} // prove that channel est is not closed/destroyed
	south <- alvasion.Alien{} // prove that channel est is not closed/destroyed
	assert.Equal(t, 4, len(c.Roads))
}

func TestEvaluateRoadsDestructionWhenOneRoadsIsDestroyed(t *testing.T) {
	// SETUP
	east := make(chan alvasion.Alien)
	close(east)
	west := make(chan alvasion.Alien, 1)
	north := make(chan alvasion.Alien, 1)
	south := make(chan alvasion.Alien, 1)
	c := alvasion.City{Name: "Foo", Roads: []chan alvasion.Alien{east, west, north, south}}
	wg := sync.WaitGroup{}

	// ACTION
	alvasion.EvaluateRoadsDestruction(&c, &wg)
	wg.Wait()

	// ASSERTIONS
	_, isOpened := <-east     // prove that channel east is not opened
	west <- alvasion.Alien{}  // prove that channel est is not closed/destroyed
	north <- alvasion.Alien{} // prove that channel est is not closed/destroyed
	south <- alvasion.Alien{} // prove that channel est is not closed/destroyed
	assert.Equal(t, 3, len(c.Roads))
	assert.False(t, isOpened)
}

func TestEvaluateRoadsDestructionWhenAllRoadsAreDestroyed(t *testing.T) {
	// SETUP
	east := make(chan alvasion.Alien)
	close(east)
	west := make(chan alvasion.Alien)
	close(west)
	north := make(chan alvasion.Alien)
	close(north)
	south := make(chan alvasion.Alien)
	close(south)
	c := alvasion.City{Name: "Foo", Roads: []chan alvasion.Alien{east, west, north, south}}
	wg := sync.WaitGroup{}

	// ACTION
	alvasion.EvaluateRoadsDestruction(&c, &wg)
	wg.Wait()

	// ASSERTIONS
	_, isOpenedEast := <-east   // prove that channel east is not opened
	_, isOpenedWest := <-west   // prove that channel west is not opened
	_, isOpenedNorth := <-north // prove that channel north is not opened
	_, isOpenedSouth := <-south // prove that channel south is not opened
	assert.Equal(t, 0, len(c.Roads))
	assert.False(t, isOpenedEast)
	assert.False(t, isOpenedWest)
	assert.False(t, isOpenedNorth)
	assert.False(t, isOpenedSouth)
}

func TestEvaluateRoadsDestructionWhenHasOnlyRoad(t *testing.T) {
	// SETUP
	east := make(chan alvasion.Alien)
	close(east)
	c := alvasion.City{Name: "Foo", Roads: []chan alvasion.Alien{east}}
	wg := sync.WaitGroup{}

	// ACTION
	alvasion.EvaluateRoadsDestruction(&c, &wg)
	wg.Wait()

	// ASSERTIONS
	_, isOpenedEast := <-east // prove that channel east is not opened
	assert.Equal(t, 0, len(c.Roads))
	assert.False(t, isOpenedEast)
}
