package pool

// Coordinates represents the location for a driver using latitude and longitude parameters
type Coordinates struct {
	Latitude  int `json:"latitude"`
	Longitude int `json:"longitude"`
	DriverID  int `json:"driver_id"`
}

// GetLatitude returns the latitude parameter of the coordinate
func (coord *Coordinates) GetLatitude() int {
	return coord.Latitude
}

// GetLongitude returns the longitude parameter of the coordinate
func (coord *Coordinates) GetLongitude() int {
	return coord.Longitude
}

// GetDriverID returns the driver ID of the coordinate
func (coord *Coordinates) GetDriverID() int {
	return coord.DriverID
}
