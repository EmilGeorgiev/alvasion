package app_test

import (
	"sort"
	"testing"

	"github.com/EmilGeorgiev/alvasion/app"
	"github.com/stretchr/testify/assert"
)

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
	wm := app.GenerateWorldMap(parts)
	sitreps := make(chan app.Sitrep, 1)
	aliens := []*app.Alien{
		{ID: 0, Sitreps: sitreps}, {ID: 1, Sitreps: sitreps}, {ID: 2, Sitreps: sitreps},
		{ID: 3, Sitreps: sitreps}, {ID: 4, Sitreps: sitreps}, {ID: 5, Sitreps: sitreps},
	}
	ac := app.NewAlienCommander(wm, aliens, sitreps)
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
	wm := app.GenerateWorldMap(parts)
	sitreps := make(chan app.Sitrep, 1)
	aliens := []*app.Alien{
		{ID: 0, Sitreps: sitreps}, {ID: 1, Sitreps: sitreps}, {ID: 2, Sitreps: sitreps},
		{ID: 3, Sitreps: sitreps}, {ID: 4, Sitreps: sitreps},
	}
	ac := app.NewAlienCommander(wm, aliens, sitreps)
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
	wm := app.GenerateWorldMap(parts)
	sitreps := make(chan app.Sitrep, 1)
	aliens := []*app.Alien{{ID: 0, Sitreps: sitreps}}
	ac := app.NewAlienCommander(wm, aliens, sitreps)
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
	wm := app.GenerateWorldMap(parts)
	sitreps := make(chan app.Sitrep, 1)
	aliens := []*app.Alien{
		{ID: 0, Sitreps: sitreps}, {ID: 1, Sitreps: sitreps}, {ID: 2, Sitreps: sitreps}, {ID: 3, Sitreps: sitreps},
		{ID: 4, Sitreps: sitreps}, {ID: 5, Sitreps: sitreps}, {ID: 6, Sitreps: sitreps}, {ID: 7, Sitreps: sitreps},
		{ID: 8, Sitreps: sitreps}, {ID: 9, Sitreps: sitreps}, {ID: 10, Sitreps: sitreps},
	}
	ac := app.NewAlienCommander(wm, aliens, sitreps)
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
