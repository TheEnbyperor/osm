package main

import (
	"math"
	"fmt"
)

type Tile struct {
	Z    int
	X    int
	Y    int
	Lat  float64
	Long float64
}

func (t *Tile) URL() string {
	return fmt.Sprintf("%d/%d/%d.png", t.Z, t.X, t.Y)
}

func (t *Tile) URLBase() string {
	return fmt.Sprintf("%d/%d", t.Z, t.X)
}

func (t *Tile) Deg2num() (x int, y int) {
	x = int(math.Floor((t.Long + 180.0) / 360.0 * (math.Exp2(float64(t.Z)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(t.Lat*math.Pi/180.0)+1.0/math.Cos(t.Lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(t.Z)))))
	return
}

func (t *Tile) Num2deg() (lat float64, long float64) {
	n := math.Pi - 2.0*math.Pi*float64(t.Y)/math.Exp2(float64(t.Z))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	long = float64(t.X)/math.Exp2(float64(t.Z))*360.0 - 180.0
	return lat, long
}

func (t *Tile) NextTile() *Tile {
	return NewTileWithXY(t.X + 1, t.Y + 1, t.Z)
}

func NewTileWithXY(x int, y int, z int) (t *Tile) {
	t = new(Tile)
	t.Z = z
	t.X = x
	t.Y = y
	t.Lat, t.Long = t.Num2deg()
	return
}

func Deg2mtrs(lat float64, lon float64) (float64, float64) {
	x := lon * 20037508.34 / 180.0
	y := math.Log(math.Tan((90.0 + lat) * math.Pi / 360.0)) / (math.Pi / 180.0)
	y = y * 20037508.34 / 180.0
	return x, y
}