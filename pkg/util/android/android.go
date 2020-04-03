package android

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tinyzimmer/android-farm-operator/pkg/util/android/adb"
	"github.com/disintegration/imaging"
	"github.com/go-logr/logr"
)

var mux sync.Mutex

// DeviceSession provides an interface for interacting with an emulated device
// over ADB. It includes utility functions for searching and detecting text on
// the screen via OCR.
type DeviceSession interface {
	BootCompleted() (bool, error)
	RunCommand(bool, ...string) ([]byte, error)
	DownloadFile(string, io.Writer) error
	GetScreencap() (image.Image, error)
	GetScreencapPNG() ([]byte, error)
	GetInvertedScreencapPNG() ([]byte, error)
	GotoLauncher() error
	LaunchApp(string) error
	Tap(x, y, count int) error
	TapAtString(*TapOptions) error
	Tab() error
	InputText(string) error
	RemoveText(int) error
	Close()
}

// deviceSession implements the DeviceSession interface
type deviceSession struct {
	host   string
	sizeX  int
	sizeY  int
	logger logr.Logger
}

// NewSession returns a connected device session or any error that arises.
func NewSession(logger logr.Logger, host string, port int32) (DeviceSession, error) {
	s := &deviceSession{host: fmt.Sprintf("%s:%d", host, port), logger: logger}
	if err := s.connect(); err != nil {
		return nil, err
	}
	return s, nil
}

// connect ensures an adb server is running and connects to the remote device.
// A global lock is used to ensure a persistent server throughout the session
// usage. This was when the client is torn down it can kill the server to avoid
// stale devices in ADB.
func (d *deviceSession) connect() error {
	mux.Lock()
	if _, err := adb.NewCommand("start-server").WithTimeout(time.Duration(5) * time.Second).Execute(); err != nil {
		mux.Unlock()
		return err
	}
	out, err := adb.NewCommand("connect", d.host).WithTimeout(time.Duration(5) * time.Second).Execute()
	if err != nil {
		mux.Unlock()
		return err
	}
	if !strings.Contains(string(out), "connected") {
		mux.Unlock()
		return fmt.Errorf("Failed to connect to device %s: %s", d.host, string(out))
	}
	// leave the lock in place allowing the client an uninterrupted session.
	// TODO: A timeut on sessions should be enforced to avoid dead locking.
	return nil
}

// Close kills the adb server for this session and releases the global lock.
func (d *deviceSession) Close() {
	mux.Unlock()
	if _, err := adb.NewCommand("kill-server").WithTimeout(time.Duration(5) * time.Second).Execute(); err != nil {
		d.logger.Error(err, "Failed to cleanly stop adb server")
	}
}

// RunCommand executes a shell command inside the remote device and returns the stdout
// or any error that occurs.
func (d *deviceSession) RunCommand(root bool, cmd ...string) ([]byte, error) {
	adbcmd := adb.NewCommand(cmd...).WithDevice(d.host).WithShell().WithTimeout(time.Duration(10) * time.Second)
	if root {
		adbcmd = adbcmd.WithRoot()
	}
	return adbcmd.Execute()
}

// DownloadFile retrieves the specified file from the device and writes its contents
// to the provided buffer
func (d *deviceSession) DownloadFile(path string, writer io.Writer) error {
	_, err := adb.NewCommand(fmt.Sprintf("cat '%s'", path)).
		WithDevice(d.host).
		WithShell().
		WithTimeout(60 * time.Second).
		WithBuffer(writer).
		WithRoot().
		Execute()
	return err
}

// BootCompleted returns true if the remote device is fully booted, false if it
// isn't, or any adb error that occurs in the process.
func (d *deviceSession) BootCompleted() (bool, error) {
	out, err := d.RunCommand(false, "getprop sys.boot_completed")
	if err != nil {
		return false, err
	}
	if strings.Contains(string(out), "1") {
		return true, nil
	}
	return false, nil
}

// InputText will type the given string into the device. It is assumed that
// the text area to fill is already selected.
func (d *deviceSession) InputText(s string) error {
	d.logger.Info(fmt.Sprintf("Inputing text: %s", s))
	cmd := fmt.Sprintf("input text '%s'", strings.Replace(s, " ", "%s", -1))
	_, err := d.RunCommand(false, cmd)
	return err
}

// RemoveText will remove the provided count of characters from the currently
// selected text area.
func (d *deviceSession) RemoveText(count int) error {
	cmd := "input keyevent 67"
	if count > 1 {
		cmd = strings.TrimSuffix(strings.Repeat(fmt.Sprintf("%s &&", cmd), count), "&&")
	}
	_, err := d.RunCommand(false, cmd)
	return err
}

// Tab sends a <Tab> event to the device. Usually to switch to the next field
// in some form or input.
func (d *deviceSession) Tab() error {
	cmd := "input keyevent 61"
	_, err := d.RunCommand(false, cmd)
	return err
}

// GotoLauncher will tap at 20 pixels above the bottom of the screen in the center
// which is where almost all android devices have their home button.
// TODO: This could be done just as easily by broadcasting an Intent over adb.
func (d *deviceSession) GotoLauncher() error {
	if err := d.ensureDimensions(); err != nil {
		return err
	}
	return d.Tap(d.sizeX/2, d.sizeY-20, 1)
}

// LaunchApp will launch the provided app on the device, bringing it to focus.
func (d *deviceSession) LaunchApp(app string) error {
	d.logger.Info(fmt.Sprintf("Launching app: %s", app))
	_, err := d.RunCommand(false, fmt.Sprintf("monkey -p %s 1", app))
	return err
}

// GetScreencap will return a raw screenshot of the device's screen.
func (d *deviceSession) GetScreencap() (image.Image, error) {
	if err := d.ensureDimensions(); err != nil {
		return nil, err
	}
	img := image.NewNRGBA(image.Rect(0, 0, d.sizeX, d.sizeY))
	raw, err := d.RunCommand(false, "screencap")
	if err != nil {
		return nil, err
	}
	img.Pix = raw
	return img, nil
}

// GetScreencapPNG returns a PNG encoded screen capture to be used in OCR
// functions.
func (d *deviceSession) GetScreencapPNG() ([]byte, error) {
	screencap, err := d.GetScreencap()
	if err != nil {
		return nil, err
	}
	return toPNG(screencap)
}

// GetInvertedScreencap return a raw capture with the pixels inverted.
// Useful for finding white or other light colored text.
func (d *deviceSession) GetInvertedScreencap() (image.Image, error) {
	screencap, err := d.GetScreencap()
	if err != nil {
		return nil, err
	}
	return imaging.Invert(imaging.Grayscale(screencap)), nil
}

// GetInvertedScreencapPNG return a PNG screen capture with the pixels inverted.
// Useful for finding white or other light colored text.
func (d *deviceSession) GetInvertedScreencapPNG() ([]byte, error) {
	screencap, err := d.GetInvertedScreencap()
	if err != nil {
		return nil, err
	}
	return toPNG(screencap)
}

// setDimensions retrieves and sets the dimensions of the screen to the current
// session.
func (d *deviceSession) setDimensions() error {
	out, err := d.RunCommand(false, "wm size")
	if err != nil {
		return err
	}
	fields := strings.Fields(string(out))
	spl := strings.Split(fields[len(fields)-1], "x")
	if len(spl) != 2 {
		return errors.New("Could not parse screen dimensions")
	}
	d.sizeX, err = strconv.Atoi(spl[0])
	if err != nil {
		return err
	}
	d.sizeY, err = strconv.Atoi(spl[1])
	if err != nil {
		return err
	}
	return nil
}

// ensureDimensions will check if we do not yet know the dimensions of the screen,
// and set them if required.
func (d *deviceSession) ensureDimensions() error {
	if d.dimensionsUnknown() {
		if err := d.setDimensions(); err != nil {
			return err
		}
	}
	return nil
}

// dimensionsUnknown returns true if the screen dimensions have not been captured
// yet.
func (d *deviceSession) dimensionsUnknown() bool {
	return d.sizeX == 0 || d.sizeY == 0
}

// toPNG converts the given raw image to PNG bytes
func toPNG(img image.Image) ([]byte, error) {
	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
