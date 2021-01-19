package location

// Location is the current position for any driver
type Location interface {
	GetLatitude() int
	GetLongitude() int
	GetDriverID() int
}
