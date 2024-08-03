package newversion

type City struct {
	east  Path
	west  Path
	north Path
	south Path
}

func (c City) Live() {
	for {
		select {
		case alien := <-c.east.IncomingDirection:
		case alien := <-c.west.IncomingDirection:
		case alien := <-c.north.IncomingDirection:
		case alien := <-c.south.IncomingDirection:
		}
	}
}
