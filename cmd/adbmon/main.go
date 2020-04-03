package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tinyzimmer/android-farm-operator/pkg/util/android/adb"
)

var verbose bool

func main() {
	var (
		noUSB               bool
		noRemote            bool
		noLocalOffline      bool
		noLocalUnauthorized bool

		device string
		host   string

		rethinkDBAddr string
		provider      string
	)

	if addr := os.Getenv("RETHINKDB_PORT_28015_TCP"); addr != "" {
		rethinkDBAddr = strings.TrimPrefix(addr, "tcp://")
	}
	if name := os.Getenv("STF_PROVIDER_NAME"); name != "" {
		provider = name
	}

	// Runtime arguments
	flag.StringVar(&provider, "provider", "", "The provider name to manage devices for")
	flag.BoolVar(&noUSB, "no-usb", false, "Don't run the USB device watcher")
	flag.BoolVar(&noRemote, "no-remote", false, "Don't run the stf provider remote device watcher")
	flag.BoolVar(&noLocalOffline, "no-local-offline", false, "Don't run the stf provider local offline watcher")
	flag.BoolVar(&noLocalUnauthorized, "no-local-unauthorized", false, "Don't run the stf provider local unauthorized watcher")
	flag.StringVar(&device, "connect", "", "Use to connect a device to another adb server from pods")
	flag.StringVar(&host, "host", "127.0.0.1", "Use this host to connect the provided device.")
	flag.BoolVar(&verbose, "verbose", false, "Verbose logging")

	flag.Parse()

	if device != "" {
		if err := runConnect(device, host, verbose); err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	runDaemon(daemonOpts{
		provider:            provider,
		rethinkDBAddr:       rethinkDBAddr,
		noUSB:               noUSB,
		noRemote:            noRemote,
		noLocalOffline:      noLocalOffline,
		noLocalUnauthorized: noLocalUnauthorized,
	})
}

func runConnect(device, host string, verbose bool) error {
	cmd := adb.NewCommand("connect", device).WithHost(host)
	if verbose {
		cmd = cmd.WithVerbose()
	}
	out, err := cmd.Execute()
	if err != nil {
		return err
	}
	log.Println(string(out))
	return nil
}

type daemonOpts struct {
	provider, rethinkDBAddr                              string
	noUSB, noRemote, noLocalOffline, noLocalUnauthorized bool
}

func runDaemon(opts daemonOpts) {
	// Setup stop channel and signal catcher
	stCh := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Received stop signal")
		stCh <- struct{}{}
	}()

	// Launch the USB watcher
	if !opts.noUSB {
		log.Println("Launching Local USB Device Watcher")
		go watchForRebindDevices()
	} else {
		log.Println("USB Watching is disabled")
	}

	if opts.provider == "" || opts.rethinkDBAddr == " " {
		log.Println("No provider name and/or rethinkdb configured, skipping STF functionality")
	} else {
		log.Println("Searching", opts.rethinkDBAddr, "for devices belonging to", opts.provider)

		if !opts.noLocalOffline {
			log.Println("Launching STF Local Offline Device Watcher")
			go watchOfflineDevices(opts.rethinkDBAddr, opts.provider)
		} else {
			log.Println("Offline device watching is disabled")
		}

		if !opts.noLocalUnauthorized {
			log.Println("Launching STF Local Unauthorized Device Watcher")
			go watchUnauthorizedDevices(opts.rethinkDBAddr, opts.provider)
		} else {
			log.Println("Unauthorized device watching is disabled")
		}

		if !opts.noRemote {
			log.Println("Launching STF Remote Device Watcher")
			go reconnectDevices(opts.rethinkDBAddr, opts.provider)
		} else {
			log.Println("Remote device watcher is disabled")
		}

	}

	// Run the ADB server
	runADBServer(stCh)
}
