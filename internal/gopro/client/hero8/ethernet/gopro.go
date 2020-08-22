// Package ethernet is the build of GoPro Hero 8 which, instead of the normal usbfs driver, instead presents itself as
// a "CDC Ethernet" device.
package ethernet

// Ethernet is the client.GoPro implementation for the GoPro Hero 8 ethernet build
type Ethernet struct {
}

// StreamStart makes the gopro start the stream
func (g Ethernet) StreamStart() {
}

// StreamSpot makes the gopro stop the stream
func (g Ethernet) StreamStop() {

}
