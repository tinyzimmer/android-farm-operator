package android

import (
	"errors"
	"fmt"
	"image"
	"strings"
	"time"

	gosseract "github.com/otiai10/gosseract/v2"
	"github.com/vitali-fedulov/images"
)

// This file is super beta and not all functionality works yet

// TapLocation represents a location on the screen to tap when searching for text
type TapLocation string

const (
	// OnString taps on the found string itself
	OnString TapLocation = "OnString"
	// StartOfLine taps on the far left side of the screen where the string was found
	StartOfLine TapLocation = "StartOfLine"
	// EndOfLine taps on the far right side of the screen where the string was found
	EndOfLine TapLocation = "EndOfLine"
)

// TapOptions are options passed to a tap method
type TapOptions struct {
	// The string to search for tap
	String string
	// The number of clicks to send when the string is found
	Count int
	// Whether to attempt scrolling until the string is found
	Scroll bool
	// Whether to attempt inverting image pixels when looking for the string
	Invert bool
	// Where to send the tap event when the string is found
	TapLocation TapLocation
}

// GetCount returns the number of taps that should be sent
func (t *TapOptions) GetCount() int {
	if t.Count == 0 {
		return 1
	}
	return t.Count
}

// Tap will send a tap event to the provided coordinates on the screen or return
// any error.
func (d *deviceSession) Tap(x, y, count int) error {
	cmd := fmt.Sprintf("input tap %d %d", x, y)
	if count > 1 {
		cmd = strings.TrimSuffix(strings.Repeat(fmt.Sprintf("%s &&", cmd), count), "&&")
	}
	_, err := d.RunCommand(false, cmd)
	return err
}

// TapAtString will search the screen for a given string, and then tap it
// depending on the provided options.
func (d *deviceSession) TapAtString(opts *TapOptions) error {
	d.logger.Info(fmt.Sprintf("Searching screen for string: %s", opts.String))
	if opts.Scroll {
		return d.TapAtStringWithScroll(opts)
	}
	if err := d.ensureDimensions(); err != nil {
		return err
	}
	var err error
	var screen []byte
	if opts.Invert {
		screen, err = d.GetInvertedScreencapPNG()
	} else {
		screen, err = d.GetScreencapPNG()
	}
	if err != nil {
		return err
	}
	x, y, err := d.getStringCoordinates(opts.String, screen)
	if err == nil {
		return d.Tap(x, y, opts.GetCount())
	}
	return fmt.Errorf("Could not locate %s on the screen", opts.String)
}

// TapAtStringWithScroll is like TapAtString except it will attempt to scroll
// the screen until the string is found.
func (d *deviceSession) TapAtStringWithScroll(opts *TapOptions) error {
	var lastScreencap image.Image
	var capFunc func() (image.Image, error)
	if opts.Invert {
		capFunc = d.GetInvertedScreencap
	} else {
		capFunc = d.GetScreencap
	}
	for {
		screen, err := capFunc()
		if err != nil {
			return err
		}
		encoded, err := toPNG(screen)
		if err != nil {
			return err
		}
		x, y, err := d.getStringCoordinates(opts.String, encoded)
		if err == nil {
			return d.Tap(x, y, opts.GetCount())
		}
		if _, err := d.RunCommand(false, fmt.Sprintf(
			"input swipe %d %d %d %d", d.sizeX/2, (d.sizeY*3)/4, d.sizeX/2, d.sizeY/3,
		)); err != nil {
			return err
		}
		if lastScreencap == nil {
			lastScreencap = screen
			continue
		}
		hashA, imgSizeA := images.Hash(lastScreencap)
		hashB, imgSizeB := images.Hash(screen)
		if images.Similar(hashA, hashB, imgSizeA, imgSizeB) {
			return errors.New("Does not appear to be anywhere else to scroll")
		}
		lastScreencap = screen
		time.Sleep(time.Duration(2) * time.Second)
	}
}

// getStringCoordinates will attempt to find the coordinates of a string on the
// device screen.
func (d *deviceSession) getStringCoordinates(s string, imgBytes []byte) (x, y int, err error) {
	client := gosseract.NewClient()
	defer client.Close()
	if err := client.SetImageFromBytes(imgBytes); err != nil {
		return 0, 0, err
	}
	boxes, err := client.GetBoundingBoxes(gosseract.RIL_WORD)
	if err != nil {
		return 0, 0, err
	}
	for _, box := range boxes {
		// d.logger.Info(fmt.Sprintf("%+v", box))
		if strings.Contains(box.Word, s) {
			d.logger.Info(fmt.Sprintf("Found Box: %+v", box))
			centerX := (box.Box.Min.X + box.Box.Max.X) / 2
			centerY := (box.Box.Min.Y + box.Box.Max.Y) / 2
			return centerX, centerY, nil
		}
	}
	return 0, 0, fmt.Errorf("Could not locate %s on the screen", s)
}
