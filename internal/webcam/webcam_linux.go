package webcam

import (
	"fmt"
	"net/url"

	"github.com/andrewhowdencom/gopro/internal/gopro"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

// Mental Notes
//
// Need to accept a list of v4l2 devices that will be available by default, as well as document instructions for how to do this.
//
// Debian packages:
// v4l2loopback-dkms v4l2loopback-utils

// Webcam is the entity that represents the camera
type Webcam struct {
	gopro *gopro.GoPro

	device string
}

// Start "switches on the webcam
func (w Webcam) Start() error {

	// Start the actual camera

	// Start the camera
	_, e := retryablehttp.Get(fmt.Sprintf("%s", &url.URL{
		Scheme: "http",
		Host:   w.gopro.Address.String(),
		Path:   "gp/gpWebcam/START",
	}))

	if e != nil {
		return errors.Wrap(e, "unable to start webcam via client")
	}

	return nil
}

// Stop "switches off" the webcam
func (w Webcam) Stop() error {
	_, e := retryablehttp.Get(fmt.Sprintf("%s", &url.URL{
		Scheme: "http",
		Host:   w.gopro.Address.String(),
		Path:   "gp/gpWebcam/STOP",
	}))

	if e != nil {
		return errors.Wrap(e, "unable to start webcam via client")
	}

	return nil
}

// Initialization functions

// WithDevice sets the device for the v4l2 output
func WithDevice(path string) Option {
	return func(w *Webcam) error {
		w.device = path

		return nil
	}
}
