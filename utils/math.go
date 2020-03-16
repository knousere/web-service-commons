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
func Round(value float64, precision int) float64 {
	multiplier := math.Pow10(precision)
	interim := math.Floor(value*multiplier + 0.5)
	return interim / multiplier
}

// SquareFloat64 squares a float64.
func SquareFloat64(number float64) float64 {
	return number * number
}

// DistanceMeters calculates distance between two lat/long coordinates
// in degrees using the haversine formula.
// Return the arc converted to meters as a positive integer.
// Return -1 if either (lat long) is not initialized.
func DistanceMeters(latA float64, longA float64, latB float64, longB float64) int {
	if (latA == 0 && longA == 0) || (latB == 0 && longB == 0) {
		return -1
	}
	radLatA := ToRadians(latA)
	radLongA := ToRadians(longA)

	radLatB := ToRadians(latB)
	radLongB := ToRadians(longB)
	haversine := Haversine(radLatA, radLongA, radLatB, radLongB)

	return RoundFloat64(EarthRadius * haversine)
}

// Haversine returns the great circle arc (haversine) in radians
// between a pair of coordinates in radians.
func Haversine(radLatA float64, radLongA float64, radLatB float64, radLongB float64) int {
	deltaLat := radLatB - radLatA
	deltaLong := radLongB - radLongA
	if deltaLat == 0 && deltaLong == 0 {
		return 0
	}

	aHaver :=
		SquareFloat64(math.Sin(deltaLat/2)) +
			math.Cos(radLatA)*
				math.Cos(radLatB)*
				SquareFloat64(math.Sin(deltaLong/2))

	haversine := 2 * math.Atan2(math.Sqrt(aHaver), math.Sqrt(1-aHaver))

	return haversine
}

// ToDegrees converts radians to degrees.
func ToDegrees(dblRadians float64) float64 {
	return dblRadians * 180 / math.Pi
}

// ToRadians converts degrees to radians.
func ToRadians(dblDegrees float64) float64 {
	return dblDegrees * math.Pi / 180
}

// DistanceToRadians converts arc distance in meters to radians.
func DistanceToRadians(dblDistance float64) float64 {
	return dblDistance * math.Pi * 2 / EarthRadius
}
