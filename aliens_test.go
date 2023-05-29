package alvasion_test

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"

	"github.com/EmilGeorgiev/alvasion"
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

//func TestListenForSitrepWhenCityIsDestroyed(t *testing.T) {
//	sitreps := make(chan alvasion.Sitrep)
//	done := make(chan struct{})
//	a := alvasion.Alien{
//		ID:      0,
//		Sitreps: sitreps,
//		Killed:  false,
//		Trapped: false,
//	}
//	ac := &alvasion.AlienCommander{
//		Sitreps: sitreps,
//		WorldMap: map[string]alvasion.City{
//			"X1": {Name: "X1", IsDestroyed: false},
//		},
//		Soldiers:      []*alvasion.Alien{&a},
//		IterationDone: done,
//	}
//
//	go func() {
//		sitreps <- alvasion.Sitrep{
//			FromAliens: []alvasion.Alien{a},
//			CityName:   "X1",
//		}
//		ac.IterationDone <- struct{}{}
//	}()
//
//	ac.ListenForSitrep()
//
//	// ASSERTIONS
//	expectedCity := alvasion.City{Name: "X1", IsDestroyed: true}
//
//	assert.Equal(t, &expectedCity, ac.WorldMap["X1"])
//	assert.True(t, ac.Soldiers[0].Killed)
//}
//
//func TestListenForSitrepWhenCityIsNOTDestroyed(t *testing.T) {
//	sitreps := make(chan alvasion.Sitrep)
//	done := make(chan struct{})
//	a := alvasion.Alien{
//		ID:      0,
//		Sitreps: sitreps,
//		Killed:  false,
//		Trapped: false,
//	}
//	ac := &alvasion.AlienCommander{
//		Sitreps: sitreps,
//		WorldMap: map[string]alvasion.City{
//			"X1": {Name: "X1", IsDestroyed: false},
//		},
//		Soldiers:      []*alvasion.Alien{&a},
//		IterationDone: done,
//	}
//
//	go func() {
//		sitreps <- alvasion.Sitrep{
//			FromAliens: []alvasion.Alien{a},
//			CityName:   "X1",
//		}
//		ac.IterationDone <- struct{}{}
//	}()
//
//	ac.ListenForSitrep()
//
//	// ASSERTIONS
//	expectedCity := alvasion.City{Name: "X1", IsDestroyed: false}
//
//	assert.Equal(t, &expectedCity, ac.WorldMap["X1"])
//	assert.False(t, ac.Soldiers[0].Killed)
//}

type connection struct {
	DirectionName string // 0 - north, 1 - south, 2 - east, 3 -west
	City          string
}

func TestStartInvasionWith6SoldiersAnd9Cities(t *testing.T) {
	// SETUP
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

	cities := buildCitiesAndTheirConnections()
	wm := alvasion.GenerateWorldMap(parts)
	sitreps := make(chan alvasion.Sitrep, 1)
	aliens := []*alvasion.Alien{
		{ID: 0, Sitreps: sitreps}, {ID: 1, Sitreps: sitreps}, {ID: 2, Sitreps: sitreps},
		{ID: 3, Sitreps: sitreps}, {ID: 4, Sitreps: sitreps}, {ID: 5, Sitreps: sitreps},
	}
	ac := alvasion.NewAlienCommander(wm, aliens, sitreps)
	notify := make(chan string)
	done := make(chan struct{})
	listenForNotificationsFromAlienCommander(cities, notify, done)

	// ACTION
	ac.SetNotifyDestroy(notify)
	ac.StartInvasion()
	close(notify)
	<-done

	// ASSERTIONS
	actual := ac.GenerateReportForInvasion2()
	expected := generateReport(cities)
	assert.Equal(t, expected, actual)
}

func TestStartInvasionWith5AliensSoldiersAnd9Cities(t *testing.T) {
	// SETUP
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

	cities := buildCitiesAndTheirConnections()
	wm := alvasion.GenerateWorldMap(parts)
	sitreps := make(chan alvasion.Sitrep, 1)
	aliens := []*alvasion.Alien{
		{ID: 0, Sitreps: sitreps}, {ID: 1, Sitreps: sitreps}, {ID: 2, Sitreps: sitreps},
		{ID: 3, Sitreps: sitreps}, {ID: 4, Sitreps: sitreps},
	}
	ac := alvasion.NewAlienCommander(wm, aliens, sitreps)
	notify := make(chan string)
	done := make(chan struct{})
	listenForNotificationsFromAlienCommander(cities, notify, done)

	// ACTION
	ac.SetNotifyDestroy(notify)
	ac.StartInvasion()
	close(notify)
	<-done

	// ASSERTIONS
	actual := ac.GenerateReportForInvasion2()
	expected := generateReport(cities)
	assert.Equal(t, expected, actual)
}

func TestStartInvasionWith1AliensSoldiersAnd9Cities(t *testing.T) {
	// SETUP
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

	cities := buildCitiesAndTheirConnections()
	wm := alvasion.GenerateWorldMap(parts)
	sitreps := make(chan alvasion.Sitrep, 1)
	aliens := []*alvasion.Alien{{ID: 0, Sitreps: sitreps}}
	ac := alvasion.NewAlienCommander(wm, aliens, sitreps)
	notify := make(chan string)
	done := make(chan struct{})
	listenForNotificationsFromAlienCommander(cities, notify, done)

	// ACTION
	ac.SetNotifyDestroy(notify)
	ac.StartInvasion()
	close(notify)
	<-done

	// ASSERTIONS
	actual := ac.GenerateReportForInvasion2()
	expected := generateReport(cities)
	assert.Equal(t, expected, actual)
}

func TestStartInvasionWith11AliensSoldiersAnd9Cities(t *testing.T) {
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

	cities := buildCitiesAndTheirConnections()
	wm := alvasion.GenerateWorldMap(parts)
	sitreps := make(chan alvasion.Sitrep, 1)
	aliens := []*alvasion.Alien{
		{ID: 0, Sitreps: sitreps}, {ID: 1, Sitreps: sitreps}, {ID: 2, Sitreps: sitreps}, {ID: 3, Sitreps: sitreps},
		{ID: 4, Sitreps: sitreps}, {ID: 5, Sitreps: sitreps}, {ID: 6, Sitreps: sitreps}, {ID: 7, Sitreps: sitreps},
		{ID: 8, Sitreps: sitreps}, {ID: 9, Sitreps: sitreps}, {ID: 10, Sitreps: sitreps},
	}
	ac := alvasion.NewAlienCommander(wm, aliens, sitreps)
	notify := make(chan string)
	done := make(chan struct{})
	listenForNotificationsFromAlienCommander(cities, notify, done)

	// ACTION
	ac.SetNotifyDestroy(notify)
	ac.StartInvasion()
	close(notify)
	<-done

	// ASSERTIONS
	actual := ac.GenerateReportForInvasion2()
	expected := generateReport(cities)
	assert.Equal(t, expected, actual)
}

func buildCitiesAndTheirConnections() map[string][]*connection {
	return map[string][]*connection{
		"X1": {{DirectionName: "south=X4", City: "X4"}, {DirectionName: "east=X2", City: "X2"}},
		"X2": {{DirectionName: "south=X5", City: "X5"}, {DirectionName: "east=X3", City: "X3"}, {DirectionName: "west=X1", City: "X1"}},
		"X3": {{DirectionName: "south=X6", City: "X6"}, {DirectionName: "west=X2", City: "X2"}},
		"X4": {{DirectionName: "north=X1", City: "X1"}, {DirectionName: "south=X7", City: "X7"}, {DirectionName: "east=X5", City: "X5"}},
		"X5": {{DirectionName: "north=X2", City: "X2"}, {DirectionName: "south=X8", City: "X8"}, {DirectionName: "east=X6", City: "X6"}, {DirectionName: "west=X4", City: "X4"}},
		"X6": {{DirectionName: "north=X3", City: "X3"}, {DirectionName: "south=X9", City: "X9"}, {DirectionName: "west=X5", City: "X5"}},
		"X7": {{DirectionName: "north=X4", City: "X4"}, {DirectionName: "east=X8", City: "X8"}},
		"X8": {{DirectionName: "north=X5", City: "X5"}, {DirectionName: "east=X9", City: "X9"}, {DirectionName: "west=X7", City: "X7"}},
		"X9": {{DirectionName: "north=X6", City: "X6"}, {DirectionName: "west=X8", City: "X8"}},
	}
}

func listenForNotificationsFromAlienCommander(cities map[string][]*connection, notify chan string, done chan struct{}) {
	go func() {
		for n := range notify {
			for _, conn := range cities[n] {
				if conn == nil {
					continue
				}
				conn2 := cities[conn.City]
				for i, c := range conn2 {
					if c == nil {
						continue
					}
					if c.City == n {
						conn2[i] = nil
					}
				}
			}
			delete(cities, n)
		}
		done <- struct{}{}
	}()
}

func generateReport(cities map[string][]*connection) map[string][]string {
	var keys []string
	for k, _ := range cities {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := map[string][]string{}
	for _, name := range keys {
		connections := cities[name]
		for _, road := range connections {
			if road == nil {
				continue
			}
			result[name] = append(result[name], road.DirectionName)
		}
	}

	return result
}
