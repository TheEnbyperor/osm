package main

import (
	"math"
)

type Tile struct {
	Z    int
	X    int
	Y    int
	Lat  float64
	Long float64
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
	lat, lng := t.Num2deg()
	return &Tile{
		X: t.X + 1,
		Y: t.Y + 1,
		Z: t.Z,
		Lat: lat,
		Long: lng,
	}
}

func NewTileWithLatLong(lat float64, long float64, z int) (t *Tile) {
	t = new(Tile)
	t.Lat = lat
	t.Long = long
	t.Z = z
	t.X, t.Y = t.Deg2num()
	return
}

func NewTileWithXY(x int, y int, z int) (t *Tile) {
	t = new(Tile)
	t.Z = z
	t.X = x
	t.Y = y
	t.Lat, t.Long = t.Num2deg()
	return
}
