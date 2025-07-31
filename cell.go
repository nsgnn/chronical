package main

import "log"

const (
	given uint = iota
	empty
	filled
	invalid
)

type Cell struct {
	x     int
	y     int
	state uint
	value rune
}

func (c *Cell) EnterValue(v rune) {
	switch c.state {
	case given:
		log.Printf("event=\"enter_value_failed\" reason=\"cell is given\" x=%d y=%d", c.x, c.y)
	default:
		c.value = v
		c.state = filled
		log.Printf("event=\"enter_value_success\" x=%d y=%d value=%c", c.x, c.y, c.value)
	}
}

func (c *Cell) Clear() {
	switch c.state {
	case given:
		log.Printf("event=\"clear_failed\" reason=\"cell is given\" x=%d y=%d", c.x, c.y)
	case empty:
		return
	default:
		log.Printf("event=\"clear_success\" x=%d y=%d value=%c", c.x, c.y, c.value)
		c.value = ' '
		c.state = empty
	}
}

func (c *Cell) RunValidation(result bool) {
	switch c.state {
	case filled:
		if !result {
			c.state = invalid
			log.Printf("event=\"validation_failed\" x=%d y=%d value=%c", c.x, c.y, c.value)
		}
	case invalid:
		if result {
			c.state = filled
			log.Printf("event=\"validation_passed\" x=%d y=%d value=%c", c.x, c.y, c.value)
		}
	}
}

func (c *Cell) View() string {
	return string(c.value)
}

func NewCell(x, y int, g *rune) *Cell {
	cell := &Cell{
		x:     x,
		y:     y,
		state: empty,
		value: ' ',
	}
	if g != nil && *g != '.' {
		cell.value = *g
		cell.state = given
	}
	log.Printf("event=\"new_cell\" x=%d y=%d value=%c state=%d", cell.x, cell.y, cell.value, cell.state)
	return cell
}
