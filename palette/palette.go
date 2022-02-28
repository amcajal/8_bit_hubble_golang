// Package palette provides color manipulation facilities for the sprites
// Algorithms and core concepts extracted from https://stackoverflow.com/questions/8507885/shift-hue-of-an-rgb-color
package palette

import (
	"math"
)

type rotationMatrix [3][3]float64
var rM rotationMatrix

// Returns RGB rotation matrix to its initial state
func resetRotationMatrix() {
	rM[0] = [3]float64{1.0, 0.0, 0.0}
	rM[1] = [3]float64{0.0, 1.0, 0.0}
	rM[2] = [3]float64{0.0, 0.0, 1.0}
}

func toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

// Sets up the RGB rotation matrix to the new hue, indicated by degrees value
func SetHueRotation(degrees int) {
	resetRotationMatrix()
	cosA := math.Cos(toRadians(float64(degrees)))
	sinA := math.Sin(toRadians(float64(degrees)))
	rM[0][0] = cosA + (1.0-cosA)/3.0
	rM[0][1] = 1./3.*(1.0-cosA) - math.Sqrt(1./3.)*sinA
	rM[0][2] = 1./3.*(1.0-cosA) + math.Sqrt(1./3.)*sinA
	rM[1][0] = 1./3.*(1.0-cosA) + math.Sqrt(1./3.)*sinA
	rM[1][1] = cosA + 1./3.*(1.0-cosA)
	rM[1][2] = 1./3.*(1.0-cosA) - math.Sqrt(1./3.)*sinA
	rM[2][0] = 1./3.*(1.0-cosA) - math.Sqrt(1./3.)*sinA
	rM[2][1] = 1./3.*(1.0-cosA) + math.Sqrt(1./3.)*sinA
	rM[2][2] = cosA + 1./3.*(1.0-cosA)
}

func clamp(value float64) uint32 {
	if value < 0 {
		return 0
	} else if value > 255 {
		return 255
	} else {
		return uint32(value + 0.5)
	}
}

// Given a RGB color, change its hue, without affecting its saturation or value
// It requires the setHueRotation function to be already executed with the new hue
func ChangeHue(r, g, b uint32) (nr, ng, nb uint32) {
	rf := float64(r)
	gf := float64(g)
	bf := float64(b)

	rx := rf*rM[0][0] + gf*rM[0][1] + bf*rM[0][2]
	gx := rf*rM[1][0] + gf*rM[1][1] + bf*rM[1][2]
	bx := rf*rM[2][0] + gf*rM[2][1] + bf*rM[2][2]
	
	return clamp(rx), clamp(gx), clamp(bx)
}

