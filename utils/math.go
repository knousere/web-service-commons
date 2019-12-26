package utils

import (
	"math"
)

// EarthRadius is the radius of the earth at the equator in meters
const EarthRadius = 6378137.0

// RoundFloat64 rounds a 64 bit floating point number
// to the nearest integer.
func RoundFloat64(number float64) int {
	// Truncate the float and use this to get the decimal component
	truncNum := math.Trunc(number)
	decimalComp := number - truncNum

	var result float64
	if decimalComp >= 0.5 {
		// Use Ceil
		result = math.Ceil(number)
	} else {
		// Use Floor
		result = math.Floor(number)
	}

	return int(result)
}

// Round is a generalized Rounding function conspicuously missing from Go math package
func Round(value float64, prec int) float64 {
	multiplier := math.Pow10(prec)
	interim := math.Floor(value*multiplier + 0.5)
	return interim / multiplier
}

// SquareFloat64 squares a float^$.
func SquareFloat64(number float64) float64 {
	return number * number
}

// Haversine calculates the great circle arc between two lat/long coordinates.
// It returns the arc converted to meters.
func Haversine(latA float64, longA float64, latB float64, longB float64) int {

	latARad := ToRadians(latA)
	longARad := ToRadians(longA)

	latBRad := ToRadians(latB)
	longBRad := ToRadians(longB)

	deltaLat := latBRad - latARad
	deltaLong := longBRad - longARad

	aHaver :=
		SquareFloat64(math.Sin(deltaLat/2)) +
			math.Cos(latARad)*
				math.Cos(latBRad)*
				SquareFloat64(math.Sin(deltaLong/2))

	cHaver := 2 * math.Atan2(math.Sqrt(aHaver), math.Sqrt(1-aHaver))

	distance := RoundFloat64(EarthRadius * cHaver)

	return distance
}

// ToDegrees converts radians to degrees.
func ToDegrees(dblRadians float64) float64 {
	return dblRadians * (180 / math.Pi)
}

// ToRadians converts degrees to radians.
func ToRadians(dblDegrees float64) float64 {
	return dblDegrees * (math.Pi / 180)
}

// DistanceToRadians converts arc distance in meters to radians.
func DistanceToRadians(dblDistance float64) float64 {
	return dblDistance / 378137.0
}
