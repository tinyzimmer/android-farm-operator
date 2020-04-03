package main

import (
	"context"
	"log"

	"github.com/tinyzimmer/android-farm-operator/pkg/util/android/adb"
)

func adbLog(msg ...interface{}) {
	out := []interface{}{}
	out = append(out, "ADB:")
	out = append(out, msg...)
	log.Println(out...)
}

func runADBServer(stCh <-chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	adbLog("Starting ADB server watcher")
	go func() {
		for {
			if _, err := adb.NewCommand("kill-server").Execute(); err != nil {
				adbLog("Failed to run kill-server")
			}
			adbLog("Launching ADB server on 0.0.0.0:5037")
			if out, err := adb.NewCommand("-a", "-P", "5037", "server", "nodaemon").
				WithContext(ctx).
				Execute(); err != nil {
				adbLog("ADB server exited with error:", err)
			} else {
				adbLog("ADB server exited:", string(out))
			}
		}
	}()
	<-stCh
	cancel()
}
