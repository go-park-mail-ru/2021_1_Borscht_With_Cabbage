package calcDistance

import (
	restModel "github.com/borscht/backend/internal/restaurant"
	"math"
	"strconv"
)

type TwoAddresses struct {
	Latitude1  float64
	Longitude1 float64
	Latitude2  float64
	Longitude2 float64
}

func GetDeliveryTime(latitudeUser, longitudeUser, latitudeRest, longitudeRest string) int {
	latitudeU, latitudeErrU := strconv.ParseFloat(latitudeUser, 64)
	longitudeU, longitudeErrU := strconv.ParseFloat(longitudeUser, 64)
	latitudeR, latitudeErrR := strconv.ParseFloat(latitudeRest, 64)
	longitudeR, longitudeErrR := strconv.ParseFloat(longitudeRest, 64)
	if longitudeErrU == nil && latitudeErrU == nil && latitudeErrR == nil && longitudeErrR == nil {
		distanse := GetDistanceFromLatLonInKm(TwoAddresses{
			Latitude1: latitudeU, Longitude1: longitudeU, Latitude2: latitudeR, Longitude2: longitudeR})
		return int(restModel.MinutesInHour*distanse/restModel.CourierSpeed + restModel.CookingTime)
	}
	return 0
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func GetDistanceFromLatLonInKm(coordinates TwoAddresses) float64 {
	R := 6371.0 // Radius of the Earth in km
	dLat := deg2rad(coordinates.Latitude2 - coordinates.Latitude1)
	dLon := deg2rad(coordinates.Longitude2 - coordinates.Longitude1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(coordinates.Latitude1))*math.Cos(deg2rad(coordinates.Latitude2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c // Distance in km
	return d
}
