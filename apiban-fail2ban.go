/*
 * Copyright (C) 2020-2024 Fred Posner (palner.com)
 *
 * This file is part of APIBAN.org.
 *
 * apiban-fail2ban is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version
 *
 * apiban-fail2ban is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * Example build commands:
 * GOOS=linux GOARCH=amd64 go build -o apiban-fail2ban
 * GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o apiban-fail2ban
 * GOOS=linux GOARCH=arm GOARM=7 go build -o apiban-fail2ban-pi
 */

package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/apiban/golib"
)

var configFileLocation string
var logFile string
var skipVerify string

func init() {
	flag.StringVar(&configFileLocation, "config", "", "location of configuration file")
	flag.StringVar(&logFile, "log", "/var/log/apiban-client.log", "location of log file or - for stdout")
	flag.StringVar(&skipVerify, "verify", "true", "set to false to skip verify of tls cert")

	if skipVerify == "false" {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

// ApibanConfig is the structure for the JSON config file
type ApibanConfig struct {
	APIKEY  string `json:"apikey"`
	LKID    string `json:"lkid"`
	VERSION string `json:"version"`
	SET     string `json:"set"`
	FLUSH   string `json:"flush"`
	JAIL    string `json:"jail"`

	sourceFile string
}

// Function to see if string within string
func contains(list []string, value string) bool {
	for _, val := range list {
		if val == value {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()
	defer os.Exit(0)

	// Open our Log
	if logFile != "-" && logFile != "stdout" {
		lf, err := os.OpenFile("/var/log/apiban-client.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Panic(err)
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			runtime.Goexit()
		}
		defer lf.Close()
		log.SetOutput(lf)
	}

	log.Print("** Started APIBAN FAIL2BAN")
	log.Print("** Licensed under GPLv2. See LICENSE for details.")
	now := time.Now()

	// Open our config file
	apiconfig, err := LoadConfig()
	if err != nil {
		log.Fatalln(err)
		runtime.Goexit()
	}

	// if no APIKEY, exit
	if apiconfig.APIKEY == "" {
		log.Fatalln("Invalid APIKEY. Exiting.")
		runtime.Goexit()
	}

	// if no APIKEY, exit
	if apiconfig.APIKEY == "MY API KEY" {
		log.Fatalln("Invalid APIKEY. Exiting. Please visit apiban.org and get an api key.")
		runtime.Goexit()
	}

	// allow cli of FULL to reset LKID to 100
	if len(os.Args) > 1 {
		arg1 := os.Args[1]
		if arg1 == "FULL" {
			log.Print("CLI of FULL received, resetting LKID")
			apiconfig.LKID = "100"
		}
	} else {
		log.Print("no command line arguments received")
	}

	// if no LKID, reset it to 100
	if len(apiconfig.LKID) == 0 {
		log.Print("Resetting LKID")
		apiconfig.LKID = "100"
	}

	// if no FLUSH, reset it
	if len(apiconfig.FLUSH) == 0 {
		log.Print("Resetting FLUSH")
		flushnow := now.Unix()
		apiconfig.FLUSH = strconv.FormatInt(flushnow, 10)
	}

	flushtime, _ := strconv.ParseInt(apiconfig.FLUSH, 10, 64)
	flushdiff := now.Unix() - flushtime
	if flushdiff >= 604800 {
		apiconfig.LKID = "100"
		apiconfig.FLUSH = strconv.FormatInt(now.Unix(), 10)
	}

	i := 0
	for i < 24 {
		log.Println("Checking banned list with ID", apiconfig.LKID, "settype", apiconfig.SET)
		// Get list of banned ip's from APIBAN.org (up to 24 times)
		res, err := golib.Banned(apiconfig.APIKEY, apiconfig.LKID, apiconfig.SET)
		if err != nil {
			log.Fatalln("failed to get banned list:", err)
			continue
		}

		if res.ID == apiconfig.LKID {
			log.Print("Great news... no new bans to add. Exiting...")
			if err := apiconfig.Update(); err != nil {
				log.Fatalln(err)
			}
			os.Exit(0)
		}

		if len(res.IPs) == 0 {
			log.Print("No IP addresses detected. Exiting.")
			os.Exit(0)
		}

		for _, ip := range res.IPs {
			log.Println("sending", ip, "to fail2ban", apiconfig.JAIL)
			cmd := exec.Command("fail2ban-client", "set", apiconfig.JAIL, "banip", ip)
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatalln(err)
				os.Exit(0)
			}

			log.Println("->", string(out[:]))
		}

		apiconfig.LKID = res.ID
	}

	// update config
	if err := apiconfig.Update(); err != nil {
		log.Fatalln(err)
	}

	log.Print("** Done. Exiting.")
}

// LoadConfig attempts to load the APIBAN configuration file from various locations
func LoadConfig() (*ApibanConfig, error) {
	var fileLocations []string

	// If we have a user-specified configuration file, use it preferentially
	if configFileLocation != "" {
		fileLocations = append(fileLocations, configFileLocation)
	}

	// If we can determine the user configuration directory, try there
	configDir, err := os.UserConfigDir()
	if err == nil {
		fileLocations = append(fileLocations, fmt.Sprintf("%s/apiban/config.json", configDir))
	}

	// Add standard static locations
	fileLocations = append(fileLocations,
		"/etc/apiban/config.json",
		"config.json",
		"/usr/local/bin/apiban/config.json",
	)

	for _, loc := range fileLocations {
		f, err := os.Open(loc)
		if err != nil {
			continue
		}
		defer f.Close()

		cfg := new(ApibanConfig)
		if err := json.NewDecoder(f).Decode(cfg); err != nil {
			return nil, fmt.Errorf("failed to read configuration from %s: %w", loc, err)
		}

		// Store the location of the config file so that we can update it later
		cfg.sourceFile = loc

		return cfg, nil
	}

	return nil, errors.New("failed to locate configuration file")
}

// Update rewrite the configuration file with and updated state (such as the LKID)
func (cfg *ApibanConfig) Update() error {
	f, err := os.Create(cfg.sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open configuration file for writing: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	return enc.Encode(cfg)
}
