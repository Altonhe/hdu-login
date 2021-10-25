// Copyright 2021 E99p1ant. All rights reserved.
// Copyright 2021 Alton. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"io/fs"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/play175/wifiNotifier"
	log "unknwon.dev/clog/v2"

	"github.com/vidar-team/srun-login/pkg/srun"
)

type config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	SSID     string `json:"SSID"`
	Host     string `json:"host"`
}

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	host := flag.String("host", "https://login.hdu.edu.cn/", "")
	username := flag.String("username", "", "")
	password := flag.String("password", "", "")
	SSID := flag.String("ssid", "i-HDU", "")
	flag.Parse()

	var c *config
	if *username != "" && *password != "" {
		c = &config{
			Username: *username,
			Password: *password,
			SSID:     *SSID,
			Host:     *host,
		}
		f, err := json.Marshal(c)
		if err != nil {
			log.Error("Failed to create config: %v", err)
		}
		err = os.WriteFile("config.json", f, fs.ModePerm)
		if err != nil {
			log.Error("Failed to write config file: %v", err)
		}
	} else {
		f, err := os.ReadFile("config.json")
		err = json.Unmarshal(f, &c)
		if err != nil {
			log.Fatal("Empty authentication pair provided, failed to unmarshal config: %+v", err)
		}
	}

	if wifiNotifier.GetCurrentSSID() == c.SSID {
		login(c)
	}

	var prev string
	wifiNotifier.SetWifiNotifier(func(ssid string) {
		if ssid == c.SSID && ssid != prev {
			log.Info("Switch WIFI to %s, attempt to connect", c.SSID)
			time.Sleep(5 * time.Second)
			login(c)
			prev = ssid
		} else {
			log.Info("Switch WIFI to %s, ignoring...", wifiNotifier.GetCurrentSSID())
		}
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}

func login(c *config) {
	client := srun.NewClient(c.Host, c.Username, c.Password)
	challengeResp, err := client.GetChallenge()
	if err != nil {
		log.Fatal("Failed to get challenge %v", err)
	}
	challenge := challengeResp.Challenge
	log.Trace("Challenge: %q", challenge)

	portalResp, err := client.Portal(challengeResp.Challenge)
	if err != nil {
		log.Fatal("Failed to portal: %v", err)
	}
	log.Trace("%+v", portalResp)
}
