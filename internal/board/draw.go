package board

import (
	"bytes"
	"diplomacy-api/internal/game"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

func Draw(mapBuf []byte, game *game.Game) ([]byte, error) {
	src, err := png.Decode(bytes.NewReader(mapBuf))
	if err != nil {
		return nil, err
	}

	bounds := src.Bounds()
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, src, bounds.Min, draw.Src)

	colors := map[string]color.RGBA{}
	for _, player := range game.Players {
		c, err := parseRGBColor(player.Color)
		if err != nil {
			return nil, err
		}

		colors[player.Name] = c
	}

	return nil, nil
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
