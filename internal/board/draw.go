package board

import (
	"bytes"
	"diplomacy-api/internal/game"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
)

const (
	armyRadius         = 8
	supplyCenterRadius = 4
)

var black = color.RGBA{A: 255}

func Draw(mapBuf []byte, g *game.Game) ([]byte, error) {
	src, err := png.Decode(bytes.NewReader(mapBuf))
	if err != nil {
		return nil, err
	}

	bounds := src.Bounds()
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, src, bounds.Min, draw.Src)

	colors := map[string]color.RGBA{}
	for _, player := range g.Players {
		c, err := parseRGBColor(player.Color)
		if err != nil {
			return nil, err
		}

		colors[player.Name] = c
	}

	for _, p := range g.Board {
		if p.SupplyCenter != nil && p.SupplyCenter.ControlledBy != nil {
			c, ok := colors[*p.SupplyCenter.ControlledBy]
			if !ok {
				return nil, fmt.Errorf("missing color for player %q", *p.SupplyCenter.ControlledBy)
			}

			drawCircle(canvas, p.SupplyCenter.Coordinates, supplyCenterRadius, c)
		}

		if p.Unit != nil {
			c, ok := colors[p.Unit.ControlledBy]
			if !ok {
				return nil, fmt.Errorf("missing color for player %q", p.Unit.ControlledBy)
			}

			switch p.Unit.Type {
			case game.UnitArmy:
				drawCircle(canvas, p.Coordinates, armyRadius, c)
			case game.UnitFleet:
				drawTriangle(canvas, p.Coordinates, c)
			default:
				return nil, fmt.Errorf("unknown unit type %q", p.Unit.Type)
			}
		}
	}

	var out bytes.Buffer
	if err := png.Encode(&out, canvas); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func parseRGBColor(s string) (color.RGBA, error) {
	var r, g, b int

	n, err := fmt.Sscanf(s, "rgb(%d,%d,%d)", &r, &g, &b)
	if err != nil || n != 3 {
		return color.RGBA{}, fmt.Errorf("invalid rgb color %q", s)
	}

	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return color.RGBA{}, fmt.Errorf("invalid rgb color %q", s)
	}

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}

func drawCircle(img *image.RGBA, coords game.Coordinates, radius int, fill color.RGBA) {
	cx := int(math.Round(coords.X))
	cy := int(math.Round(coords.Y))
	r2 := radius * radius
	inner := (radius - 1) * (radius - 1)

	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			dx := x - cx
			dy := y - cy
			d2 := dx*dx + dy*dy

			if d2 > r2 {
				continue
			}
			if d2 >= inner {
				img.Set(x, y, black)
			} else {
				img.Set(x, y, fill)
			}
		}
	}
}

func drawTriangle(img *image.RGBA, coords game.Coordinates, fill color.RGBA) {
	top := image.Point{X: int(math.Round(coords.X)), Y: int(math.Round(coords.Y)) - 5}
	left := image.Point{X: int(math.Round(coords.X)) - 12, Y: int(math.Round(coords.Y)) + 5}
	right := image.Point{X: int(math.Round(coords.X)) + 12, Y: int(math.Round(coords.Y)) + 5}

	minX := min(left.X, top.X, right.X)
	maxX := max(left.X, top.X, right.X)
	minY := min(left.Y, top.Y, right.Y)
	maxY := max(left.Y, top.Y, right.Y)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if pointInTriangle(image.Point{X: x, Y: y}, left, top, right) {
				img.Set(x, y, fill)
			}
		}
	}

	drawLine(img, left, top, black)
	drawLine(img, top, right, black)
	drawLine(img, right, left, black)
}

func pointInTriangle(p, a, b, c image.Point) bool {
	d1 := sign(p, a, b)
	d2 := sign(p, b, c)
	d3 := sign(p, c, a)

	hasNeg := d1 < 0 || d2 < 0 || d3 < 0
	hasPos := d1 > 0 || d2 > 0 || d3 > 0

	return !(hasNeg && hasPos)
}

func sign(p1, p2, p3 image.Point) int {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}

func drawLine(img *image.RGBA, from, to image.Point, c color.RGBA) {
	x0 := from.X
	y0 := from.Y
	x1 := to.X
	y1 := to.Y

	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)

	sx := -1
	if x0 < x1 {
		sx = 1
	}

	sy := -1
	if y0 < y1 {
		sy = 1
	}

	err := dx + dy
	for {
		img.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}
