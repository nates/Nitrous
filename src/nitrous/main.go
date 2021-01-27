package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	proxyType   string = "http"
	proxyMethod string = "cycle"
	threads     int    = 10
	checks      int    = 0
	timeout     int    = 5
	debug       bool   = false
	start       time.Time
)

func main() {
	flag.StringVar(&proxyType, "type", proxyType, "Type of proxies to use. [http | socks4 | socks5]")
	flag.StringVar(&proxyMethod, "method", proxyMethod, "Proxy rotation method. [cycle | random]")
	flag.IntVar(&threads, "threads", threads, "Amount of threads to use.")
	flag.IntVar(&timeout, "timeout", timeout, "Timeout in seconds for each request.")
	flag.BoolVar(&debug, "debug", debug, "Debug to console.")
	flag.Parse()

	if proxyType != "http" && proxyType != "socks4" && proxyType != "socks5" {
		log.Fatalf("Error: %v", "Proxy type must be http, socks4, or socks5.")
	}

	if proxyMethod != "cycle" && proxyMethod != "random" {
		log.Fatalf("Error: %v", "Proxy method must be cycle or random.")
	}

	if threads <= 0 {
		log.Fatalf("Error: %v", "Threads must be greater than 0.")
	}

	err := loadProxies()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("[Main] Starting with %v threads and %v proxies.", threads, len(proxies))
	start = time.Now()

	for i := 0; i < threads; i++ {
		go worker()
	}

	select {}
}

func worker() {
	for {
		code, err := genCode()
		if err != nil {
			if debug {
				log.Printf("Error: %v", err)
			}
			continue
		}
		valid, err := check(code)
		if err != nil {
			if debug {
				log.Printf("Error: %v", err)
			}
			continue
		}
		checks++
		log.Printf("%v | CPS: %.1f | Valid: %v", code, float64(checks)/(float64(time.Since(start).Seconds())+1), valid)
		if valid {
			if _, err := os.Stat("codes.txt"); os.IsNotExist(err) {
				_, err := os.Create("codes.txt")
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
			}

			f, err := os.OpenFile("codes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			if _, err := f.Write([]byte(fmt.Sprintf("%v\n", code))); err != nil {
				log.Fatalf("Error: %v", err)
			}
			if err := f.Close(); err != nil {
				log.Fatalf("Error: %v", err)
			}
		}
	}
}
