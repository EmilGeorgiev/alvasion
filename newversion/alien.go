package newversion

import (
	"errors"
	"log"
	"math/rand"
	"sync"
)

type Alien struct {
	Name         string
	maxMovements int64
	paths        chan Path
	wg           *sync.WaitGroup
}

func New(n string, maxMov int64, wg *sync.WaitGroup) Alien {
	return Alien{
		Name:         n,
		maxMovements: maxMov,
		paths:        make(chan Path),
		wg:           wg,
	}
}

func (a Alien) ChoosePath(paths []Path) {
	a.choosePath(paths)
}

func (a Alien) choosePath(paths []Path) {
	defer a.wg.Done()
	select {
	case p := <-a.paths:
		a.maxMovements++
		if a.maxMovements > a.maxMovements {
			return
		}
		p, err := chooseRandomPath(paths)
		if err != nil {
			log.Println(err)
			return
		}

		p.OutgoingDirection <- a
	}
}

type Path struct {
	OutgoingDirection chan<- Alien
	IncomingDirection <-chan Alien
	Closed            bool
}

var chooseRandomPath = NewDefaultRandomPath()

type RandomPath func(paths []Path) (Path, error)

func NewDefaultRandomPath() RandomPath {
	return defaultRandomPath
}

func defaultRandomPath(paths []Path) (Path, error) {
	var availablePaths []Path
	for _, p := range paths {
		if p.Closed {
			continue
		}
		availablePaths = append(availablePaths, p)
	}
	if len(availablePaths) == 0 {
		return Path{}, errors.New("no available paths. The alien is stuck")
	}

	i := rand.Intn(len(availablePaths))
	return availablePaths[i], nil
}
