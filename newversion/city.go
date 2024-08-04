package newversion

import (
	"fmt"
	"io"
	"log"
)

type City struct {
	Name        string
	Log         io.StringWriter
	paths       []Path
	aliens      []Alien
	isDestroyed bool
	stop        chan struct{}
}

func NewCity(name string, paths []Path) City {
	return City{Name: name, paths: paths, stop: stop}
}

func (c *City) AddAlien(a Alien) {
	c.aliens = append(c.aliens, a)
}

func (c *City) Live() {
	for {
		for i := 0; i < len(c.paths); {
			alien, isPathOpened := c.checkPathForIncomingAlien(c.paths[i])
			if !isPathOpened {
				// Remove the element by appending the slice before and after the current index
				c.paths = append(c.paths[:i], c.paths[i+1:]...)
				continue
			}
			if alien != nil {
				c.aliens = append(c.aliens, *alien)
			}
			i++
		}
		if len(c.aliens) > 1 {
			c.destroy()
			return
		}

		if len(c.aliens) == 0 {
			continue
		}

		alien := c.aliens[0]
		c.aliens = []Alien{}
		alien.choosePath(c.paths)
	}
}

func (c *City) destroy() {
	for _, path := range c.paths {
		close(path.OutgoingDirection)
	}
	c.isDestroyed = true
	_, _ = c.Log.WriteString(fmt.Sprintf("%s has been destroyed by alien %s and alien %s!", c.Name, c.aliens[0].Name, c.aliens[1].Name))
	for _, a := range c.aliens[2:] {
		log.Println(fmt.Sprintf("Alien %s is in the same city and it is killed too.", a.Name))
	}
}

func (c *City) checkPathForIncomingAlien(path Path) (*Alien, bool) {
	select {
	case alien, ok := <-path.IncomingDirection:
		if !ok {
			return nil, false
		}
		return &alien, true
	default:
		return nil, true
	}
}
