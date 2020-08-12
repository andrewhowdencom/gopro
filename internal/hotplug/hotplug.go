// Package hotplug provides a wrapper that notifies the program whenever a GoPro camera has been
// connected or disconnected via USB
package hotplug

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"net"

	"github.com/andrewhowdencom/gopro/internal/gopro"
	"github.com/jochenvg/go-udev"
	"github.com/pkg/errors"
)

// GoProPciVendorID is the vendor that is attached to the hardware GoPro makes
const GoProPciVendorID = "2672"

// Supported Hardware

// GoProPciDeviceIDHero8BlackAsModem -- https://gopro.com/en/us/shop/cameras/hero8-black/CHDHX-801-master.html (webcam firmware)
const GoProPciDeviceIDHero8BlackAsModem = "0050"

// GoProPciDeviceIDHero8BlackAsFilesystem -- https://gopro.com/en/us/shop/cameras/hero8-black/CHDHX-801-master.html (normal firmware)
const GoProPciDeviceIDHero8BlackAsFilesystem = "0049"

// GoProIP is the IP that the device apparently continues to use always
//
// Todo: Verify if this is true *across* gopro
const GoProIP = "172.26.169.51"

const (
	// Connected means listen only for connect events
	Connected = iota

	// Disconnected means listen only for disconnect events
	Disconnected
)

// Hotplug provides the entity that listens for the changes
type Hotplug struct {
	// A list of "events" to filter for, such as
	// - connected
	// - disconnected
	Filters []int
}

// Event represents the GoPro having been connected
type Event struct {
	// Whether the camera was connected or disconnected
	Type int

	Entity *gopro.GoPro
}

// Option is a "functional argument".
// See:
// - https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type Option func(*Hotplug) error

// New creates a new "hotplug"
func New(options ...Option) (*Hotplug, error) {
	return &Hotplug{
		Filters: []int{Connected, Disconnected},
	}, nil
}

// Listen listens for connection or disconnection events, returning an "GoPro" entity.
func (p Hotplug) Listen(ctx context.Context) (<-chan *Event, error) {
	c := make(chan *Event)
	u := udev.Udev{}
	m := u.NewMonitorFromNetlink("udev")

	ch, err := m.DeviceChan(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable create listener for ")
	}

	// Dispatch a goroutine to run the wait
	go func() {

		// Check in a loop
		for {

			// Await the results of either the cancellation or the plug/unplug events
			select {
			case <-ctx.Done():
				close(c)
				return
			case d := <-ch:
				// Filter out the events to the ones we really care about
				if p.isGoPro(d) && p.isEvent(d) && p.isNetworkSubsystem(d) {
					g, e := getGoPro(d)

					if e != nil {
						log.Printf("unable to generate gopro: %s", e.Error())
						continue
					}

					c <- &Event{
						Type: getEvent(d),

						Entity: g,
					}
				}
			}
		}
	}()

	return c, nil
}

func getGoPro(dev *udev.Device) (*gopro.GoPro, error) {
	var opts []gopro.Option

	// Construct an identifier
	h := sha1.New()
	h.Write([]byte(dev.Devpath()))

	// If it was a connection event, assign it an address
	if getEvent(dev) == Connected {
		opts = append(opts, gopro.WithAddress(&net.IPAddr{
			IP: net.ParseIP(GoProIP),
		}))
	}

	// Return the constructed object
	return gopro.New(fmt.Sprintf("%x", h.Sum(nil)), opts...)
}

func getEvent(dev *udev.Device) int {
	var result int

	// Check Action
	switch dev.Action() {
	case "add":
		result = Connected
	case "remove":
		result = Disconnected

	// Unspuported
	default:
		result = -1
	}

	return result
}

// Takes a hardware event, and returns if it is one we're interested in.
func (p Hotplug) isEvent(dev *udev.Device) bool {
	// See if action is among filters
	for _, filter := range p.Filters {
		if getEvent(dev) == filter {
			return true
		}
	}

	return false
}

// Takes a hardware event, and returns if it is a "supported" gopro
func (p Hotplug) isGoPro(dev *udev.Device) bool {
	// Verify Vendor
	if dev.PropertyValue("ID_VENDOR_ID") != GoProPciVendorID {
		return false
	}

	// Verify Model
	model := dev.PropertyValue("ID_MODEL_ID")
	for _, id := range []string{GoProPciDeviceIDHero8BlackAsModem} {
		if id == model {
			return true
		}
	}

	return false
}

func (p Hotplug) isNetworkSubsystem(dev *udev.Device) bool {
	if dev.Subsystem() == "net" {
		return true
	}

	return false
}
