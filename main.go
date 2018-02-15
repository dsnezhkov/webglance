// Webglance. Screenshot a list of pages to PDF/PNG for a quick visual survey
package main

import (
	"flag"
	"fmt"
	"github.com/raff/godet"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {

	const d_urlfile = "urls.txt"
	const d_hostdbg = "localhost:9222"
	const d_pagedelay = 5
	const d_outdir = "output"
	urlfile := flag.String("urlfile", d_urlfile, "Path to Url file. ")
	hostdbg := flag.String("hostdbg", d_hostdbg, "Chrome debugger <host:port>. ")
	pagedelay := flag.Int("pagedelay", d_pagedelay, "Delay Page load before screen. ")
	outdir := flag.String("outdir", d_outdir, "Results output directory.")

	flag.Parse()

	if f := flag.Lookup("urlfile"); fmt.Sprintf("%s", f.Value) == d_urlfile {
		log.Printf("Taking defaults (%s) on url file\n", d_urlfile)
	} else {
		log.Printf("Reading URL File: %s\n", *hostdbg)
	}
	if h := flag.Lookup("hostdbg"); fmt.Sprintf("%s", h.Value) == d_hostdbg {
		log.Printf("Taking defaults (%s) on debugger location\n", d_hostdbg)
	} else {
		log.Printf("Debugger: %s\n", *hostdbg)
	}
	if o := flag.Lookup("outdir"); fmt.Sprintf("%s", o.Value) == d_outdir {
		log.Printf("Taking defaults (%s) on output location\n", d_outdir)
	} else {
		log.Printf("Output saved to: %s\n", *outdir)
	}

	// connect to Chrome instance
	var remote *godet.RemoteDebugger

	log.Printf("Attempting to connect to Chrome Debugger on %s", *hostdbg)
	remote, err := godet.Connect(*hostdbg, false)
	if err != nil {
		log.Fatal("Cannot connect to Chrome instance: ", err)
		return
	}
	log.Print("Connected.")
	// disconnect when done
	defer remote.Close()

	remote.SetVisibleSize(1280, 1696)

	// read urls from file
	// format per line http(s)://host[:port]
	content, err := ioutil.ReadFile(*urlfile)
	if err != nil {
		log.Fatal("File cannot be read")
		return
	}

	if fi, err := os.Lstat(*outdir); err != nil {
		if os.IsNotExist(err) {
			// dir does not exist
			log.Printf("Directory %s does not exist\n", *outdir)
			log.Printf("Directory %s wll be created\n", *outdir)
			os.Mkdir(*outdir, 0700)
		} else {
			if mode := fi.Mode(); mode.IsDir() {
				log.Printf("Directory %s already exists\n", *outdir)
			} else {
				log.Fatal("%s is not a directory\n", *outdir)
				return
			}
		}
	} else {
		log.Printf("Directory %s exists. Saving PDF/PNG there\n", *outdir)
	}

	urlstrings := strings.Split(string(content), "\n")

	for _, urlstr := range urlstrings {

		// empty strings
		if urlstr == "" {
			continue
		}
		// commented out strings
		if match, _ := regexp.MatchString("^#", urlstr); match == true {
			log.Printf("Commented (%s) Skipping.\n", urlstr)
			continue
		}

		u, err := url.Parse(urlstr)

		if (err) != nil {
			log.Printf("Malformed (%s) Skipping.\n", urlstr)
			continue
		}
		host := u.Host
		hostP, port, _ := net.SplitHostPort(u.Host)

		var fname string
		if (port) != "" {
			fname = strings.Join([]string{hostP, port}, "_")
		} else {
			fname = host
		}

		log.Printf("[+] Screening (%s)\n", urlstr)
		screen(remote, urlstr, fname, outdir, pagedelay)
	}

}

func screen(remote *godet.RemoteDebugger, urlstr string, fileName string, outdir *string, pagedelay *int) {
	var err error

	// create new tab
	tab, _ := remote.NewTab(urlstr)

	// navigate in existing tab
	_ = remote.ActivateTab(tab)

	// re-enable events when changing active tab
	remote.AllEvents(true) // enable all events

	time.Sleep(time.Duration(*pagedelay) * time.Second)

	// take a screenshot
	err = remote.SaveScreenshot(fmt.Sprintf("%s/%s.png", *outdir, fileName), 0644, 100, true)
	if err != nil {
		log.Printf("Unable to save PNG (%s)\n", err)
	}

	// save page as PDF (long pages / content)
	err = remote.SavePDF(fmt.Sprintf("%s/%s.pdf", *outdir, fileName), 0644, godet.PortraitMode(), godet.Scale(1.0))
	if err != nil {
		log.Printf("Unable to save PDF (%s)\n", err)
	}

	err = remote.CloseTab(tab)
	if err != nil {
		log.Printf("Unable to close tab (%s)\n", err)
	}

}
