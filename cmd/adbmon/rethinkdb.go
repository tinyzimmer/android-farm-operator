package main

import (
	"log"
	"strings"
	"time"

	"github.com/tinyzimmer/android-farm-operator/pkg/util/android/adb"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/rethinkdb"
)

func reconnectDevices(rdbAddr, provider string) {
	ticker := time.NewTicker(time.Duration(5 * time.Second))
	for range ticker.C {
		session, err := rethinkdb.NewSession(rdbAddr)
		if err != nil {
			log.Println("REMOTE: Failed to connect to rethinkdb, skipping:", err.Error())
			continue
		}
		devices, err := session.GetAllDevicesForProvider(provider)
		if err != nil {
			log.Println("REMOTE: Failed to query devices for provider", provider, "error:", err.Error())
			session.Close()
			continue
		}
		out, err := adb.NewCommand("devices").Execute()
		if err != nil {
			log.Println("REMOTE: Error listing ADB devices")
			session.Close()
			continue
		}
		adbDevices := strings.Split(string(out), "\n")
		for _, device := range devices {
			if len(strings.Split(device, ":")) > 1 {
				if !devicePresentAndOnline(adbDevices, device) {
					log.Println("REMOTE: Reconnecting remote device:", device)
					if out, err := adb.NewCommand("connect", device).Execute(); err != nil {
						log.Println("REMOTE: Failed to reconnect remote device:", err.Error())
					} else {
						log.Println("REMOTE:", strings.TrimSpace(string(out)))
					}
				}
			}
		}
		session.Close()
	}
}

func watchOfflineDevices(rdbAddr, provider string) {
	ticker := time.NewTicker(time.Duration(5 * time.Second))
	for range ticker.C {
		session, err := rethinkdb.NewSession(rdbAddr)
		if err != nil {
			log.Println("RDB: Could not connect to rethinkdb instance at", rdbAddr, "error:", err)
			continue
		}
		devices, err := session.GetDevicesForProviderByStatus(provider, rethinkdb.StatusOffline)
		if err != nil {
			session.Close()
			log.Println("RDB: Failed to query offline devices:", err)
			continue
		}
		session.Close()
		for _, device := range devices {
			log.Println("RDB: Reconnecting offline device:", device)
			if _, err := adb.NewCommand("reconnect").WithDevice(device).Execute(); err != nil {
				log.Println("RDB: Failed to reconnect device:", err)
			}
		}
	}
}

func watchUnauthorizedDevices(rdbAddr, provider string) {
	ticker := time.NewTicker(time.Duration(5 * time.Second))
	for range ticker.C {
		session, err := rethinkdb.NewSession(rdbAddr)
		if err != nil {
			log.Println("RDB: Could not connect to rethinkdb instance at", rdbAddr, "error:", err)
			continue
		}
		devices, err := session.GetDevicesForProviderByStatus(provider, rethinkdb.StatusUnauthorized)
		if err != nil {
			session.Close()
			log.Println("RDB: Failed to query unauthorized devices:", err)
			continue
		}
		session.Close()
		for _, device := range devices {
			log.Println("RDB: Reconnecting unauthorized device:", device)
			if _, err := adb.NewCommand("reconnect").WithDevice(device).Execute(); err != nil {
				log.Println("RDB: Failed to reconnect device:", err)
			}
		}
	}
}

func devicePresentAndOnline(adbDevices []string, device string) bool {
	for _, x := range adbDevices {
		if strings.HasPrefix(x, device) {
			return strings.HasSuffix(x, "device") || strings.HasSuffix(x, "online")
		}
	}
	return false
}
