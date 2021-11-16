// Copyright 2021 E99p1ant. All rights reserved.
// Copyright 2021 Alton. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"
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

LOGIN:
	for {
		redirectURL, err := getMicrosoftRedirectResponse()
		if err != nil {
			log.Error("getMicrosoftRedirectResponse failed: %v", err)
		}

		u, err := url.Parse(*host)
		if err != nil {
			log.Error("Failed to parse host: %v", err)
			return
		}

		switch redirectURL.Host {
		case u.Host:
			err := login(c)
			if err != nil {
				log.Error("Failed to login: %v", err)
			}
		case "go.microsoft.com":
			log.Info("Succeed Login!")
			break LOGIN
		case "2.2.2.2":
			log.Error("unexpected redirect host, expected: %v, actual %v", u.Host, redirectURL.Host)
			err := oldLogin(c)
			if err != nil {
				log.Error("Failed to authenticate: %v", err)
			}
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}

func login(c *config) error {
	client := srun.NewClient(c.Host, c.Username, c.Password)
	challengeResp, err := client.GetChallenge()
	if err != nil {
		return errors.Wrap(err, "get challenge response")
	}
	challenge := challengeResp.Challenge
	log.Trace("Challenge: %q", challenge)

	portalResp, err := client.Portal(challengeResp.Challenge)
	if err != nil {
		return errors.Wrap(err, "portal challenge")
	}

	log.Trace("%+v", portalResp)
	return nil
}

func getMicrosoftRedirectResponse() (*url.URL, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get("http://www.msftconnecttest.com/redirect")
	if err != nil {
		return nil, err
	}

	location, err := resp.Location()
	if err != nil {
		return nil, err
	}
	log.Trace("location: %v", location.String())
	return location, nil
}

func oldLogin(c *config) error {
	client := http.Client{}
	reader := strings.NewReader(fmt.Sprintf("opr=pwdLogin&userName=%s&pwd=%s&rememberPwd=0", c.Username, c.Password))
	req, err := http.NewRequest(http.MethodPost, "http://2.2.2.2/ac_portal/login.php", reader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
