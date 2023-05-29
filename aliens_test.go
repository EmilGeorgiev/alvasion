package alvasion_test

import (
	"fmt"
	"testing"

	"github.com/EmilGeorgiev/alvasion"
	"github.com/stretchr/testify/assert"
)

//func TestGiveOrders(t *testing.T) {
//	north := make(chan alvasion.Alien, 1)
//	south := make(chan alvasion.Alien, 1)
//	east := make(chan alvasion.Alien, 1)
//	west := make(chan alvasion.Alien, 1)
//
//	c := alvasion.City{
//		Name:          "Foo",
//		OutgoingRoads: []chan alvasion.Alien{north, south, east, west},
//		IsDestroyed:   false,
//		Alien:         alvasion.Alien{ID: 55, Sitreps: make(chan alvasion.Sitrep)},
//	}
//
//	al := &alvasion.AlienCommander{}
//	roadIndex := al.GiveOrdersToTheAlienIn(c)
//
//	actual := <-c.OutgoingRoads[roadIndex]
//
//	// ASSERTIONS
//	// approve that the method put only one alien in one randomly picked channel
//	close(north)
//	close(south)
//	close(east)
//	close(west)
//	_, northIsOpen := <-north
//	_, southIsOpen := <-south
//	_, eastIsOpen := <-east
//	_, westIsOpen := <-west
//
//	assert.Equal(t, 55, actual.ID)
//	assert.False(t, northIsOpen)
//	assert.False(t, southIsOpen)
//	assert.False(t, eastIsOpen)
//	assert.False(t, westIsOpen)
//}
//
//func TestGiveOrdersIfCityIsNotOccupiedByAlien(t *testing.T) {
//	north := make(chan alvasion.Alien, 1)
//	south := make(chan alvasion.Alien, 1)
//	east := make(chan alvasion.Alien, 1)
//	west := make(chan alvasion.Alien, 1)
//
//	c := alvasion.City{
//		Name:          "Foo",
//		OutgoingRoads: []chan alvasion.Alien{north, south, east, west},
//		IsDestroyed:   false,
//	}
//
//	al := &alvasion.AlienCommander{}
//	randIndex := al.GiveOrdersToTheAlienIn(c)
//
//	// ASSERTIONS
//	// approve that the method doesn't put an alien in one randomly picked channel
//	close(north)
//	close(south)
//	close(east)
//	close(west)
//	_, northIsOpen := <-north
//	_, southIsOpen := <-south
//	_, eastIsOpen := <-east
//	_, westIsOpen := <-west
//
//	assert.False(t, northIsOpen)
//	assert.False(t, southIsOpen)
//	assert.False(t, eastIsOpen)
//	assert.False(t, westIsOpen)
//	assert.Equal(t, -1, randIndex)
//}

func TestListenForSitrepWhenCityIsDestroyed(t *testing.T) {
	sitreps := make(chan alvasion.Sitrep)
	done := make(chan struct{})
	ac := &alvasion.AlienCommander{
		Sitreps: sitreps,
		WorldMap: alvasion.WorldMap{Cities: map[string]*alvasion.City{
			"X1": {
				Name:        "X1",
				IsDestroyed: false,
			},
		}},
		Soldiers: []*alvasion.Alien{{
			ID:      0,
			Sitreps: sitreps,
			Killed:  false,
			Trapped: false,
		}},
		IterationDone: done,
	}

	go func() {
		sitreps <- alvasion.Sitrep{
			From:            0,
			CityName:        "X1",
			IsCityDestroyed: true,
		}
		ac.IterationDone <- struct{}{}
	}()

	ac.ListenForSitrep()

	// ASSERTIONS
	expectedCity := alvasion.City{Name: "X1", IsDestroyed: true}

	assert.Equal(t, &expectedCity, ac.WorldMap.Cities["X1"])
	assert.True(t, ac.Soldiers[0].Killed)
}

func TestListenForSitrepWhenCityIsNOTDestroyed(t *testing.T) {
	sitreps := make(chan alvasion.Sitrep)
	done := make(chan struct{})
	ac := &alvasion.AlienCommander{
		Sitreps: sitreps,
		WorldMap: alvasion.WorldMap{Cities: map[string]*alvasion.City{
			"X1": {
				Name:        "X1",
				IsDestroyed: false,
			},
		}},
		Soldiers: []*alvasion.Alien{{
			ID:      0,
			Sitreps: sitreps,
			Killed:  false,
			Trapped: false,
		}},
		IterationDone: done,
	}

	go func() {
		sitreps <- alvasion.Sitrep{
			From:            0,
			CityName:        "X1",
			IsCityDestroyed: false,
		}
		ac.IterationDone <- struct{}{}
	}()

	ac.ListenForSitrep()

	// ASSERTIONS
	expectedCity := alvasion.City{Name: "X1", IsDestroyed: false}

	assert.Equal(t, &expectedCity, ac.WorldMap.Cities["X1"])
	assert.False(t, ac.Soldiers[0].Killed)
}

func TestStartInvasion(t *testing.T) {
	parts := make(chan []string)
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

	sitreps := make(chan alvasion.Sitrep, 1)
	aliens := []*alvasion.Alien{
		{ID: 0, Sitreps: sitreps},
		{ID: 1, Sitreps: sitreps},
		{ID: 2, Sitreps: sitreps},
		{ID: 3, Sitreps: sitreps},
		{ID: 4, Sitreps: sitreps},
		{ID: 5, Sitreps: sitreps},
	}

	//for i := 0; i < 6; i++ {
	//	aliens[0] = alvasion.Alien{
	//		ID:      0,
	//		Sitreps: sitreps,
	//		Killed:  false,
	//		Trapped: false,
	//	}
	//}

	ac := alvasion.NewAlienCommander(wm, aliens, sitreps)

	ac.StartInvasion()

	report := ac.GenerateReportForInvasion()
	fmt.Println("---------------------------")
	fmt.Println(report)
}
