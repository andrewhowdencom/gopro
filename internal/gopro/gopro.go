package gopro

import (
	"net"

	"github.com/pkg/errors"
)

// GoPro is the entity that represents the GoPro camera
type GoPro struct {
	// A uniqe identifier to target this instance of the GoPro camera. Not consistent across different appearances of
	// the same camera
	ID string

	Address net.Addr
}

// Option is a modifier that gets passed into the canonical constructor for this package to modify the behaviour of
// this function.
type Option func(g *GoPro) error

// New creates a new GoPro entity
func New(ID string, o ...Option) (*GoPro, error) {
	g := &GoPro{
		ID: ID,
	}

	for _, option := range o {
		if err := option(g); err != nil {
			return nil, errors.Wrap(err, "unable to modify GoPro struct")
		}
	}

	return g, nil
}

// WithAddress creates a GoPro entity that is addressable over the network
func WithAddress(address net.Addr) Option {
	return func(g *GoPro) error {
		g.Address = address

		return nil
	}
}
