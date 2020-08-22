// Package webcam takes a GoPro object and sets up an associated device that can be used by the operationg system
// as a webcam
package webcam

import (
	"github.com/andrewhowdencom/gopro/internal/gopro"
	"github.com/pkg/errors"
)

// Option is a modifier used to construct the webcam
type Option func(*Webcam) error

// New is the constructor for the webcam,
func New(gopro *gopro.GoPro, options ...Option) (*Webcam, error) {
	w := &Webcam{
		gopro: gopro,
	}

	for _, o := range options {
		if e := o(w); e != nil {
			return nil, errors.Wrap(e, "unable to apply option")
		}
	}

	return w, nil
}
