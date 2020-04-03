package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tinyzimmer/android-farm-operator/pkg/util/android/adb"
)

const devicePath = "/sys/devices"
const unbindPath = "/sys/bus/usb/drivers/usb/unbind"
const bindPath = "/sys/bus/usb/drivers/usb/bind"

func usbLog(msg ...interface{}) {
	if verbose {
		out := []interface{}{}
		out = append(out, "USB:")
		out = append(out, msg...)
		log.Println(out...)
	}
}

func watchForRebindDevices() {
	ticker := time.NewTicker(time.Duration(5) * time.Second)
	for range ticker.C {
		usbLog("Checking for disconnected USB devices")
		if err := rebindDevices(); err != nil {
			usbLog("Error running usb rebind:", err)
		}
	}
}

func rebindDevices() error {
	if err := filepath.Walk(devicePath, func(path string, info os.FileInfo, err error) error {
		if info.Name() == "bInterfaceSubClass" {
			details, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			if !strings.Contains(string(details), "42") {
				return nil
			}
			usbLog("Checking matched device:", path)
			class, err := ioutil.ReadFile(filepath.Join(filepath.Dir(path), "bInterfaceClass"))
			if err != nil {
				return err
			}
			classStr := strings.TrimSpace(string(class))
			usbLog("Device Class:", classStr)
			proto, err := ioutil.ReadFile(filepath.Join(filepath.Dir(path), "bInterfaceProtocol"))
			if err != nil {
				return err
			}
			protoStr := strings.TrimSpace(string(proto))
			usbLog("Device Protocol:", protoStr)
			if !strings.Contains(classStr, "ff") || !strings.Contains(protoStr, "01") {
				usbLog("Device class and protocol do not match, skipping")
				return nil
			}
			serial, err := ioutil.ReadFile(filepath.Join(filepath.Dir(path), "..", "serial"))
			if err != nil {
				return err
			}
			serialStr := strings.TrimSpace(string(serial))
			usbLog("Class and protocol match, checking ADB for device with serial:", serialStr)
			out, err := adb.NewCommand("devices").Execute()
			if err != nil {
				return err
			}
			if !strings.Contains(string(out), serialStr) {
				usbLog("Device", serialStr, "is not connected to ADB, rebinding...")
				pathSplit := strings.Split(path, "/")
				deviceID := pathSplit[len(pathSplit)-2]
				usbLog("Unbinding and rebinding USB device with ID:", deviceID)
				if err := ioutil.WriteFile(unbindPath, []byte(deviceID), 0666); err != nil {
					return err
				}
				if err := ioutil.WriteFile(bindPath, []byte(deviceID), 0666); err != nil {
					return err
				}
			} else {
				usbLog("Device", serialStr, "is connected to ADB")
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
