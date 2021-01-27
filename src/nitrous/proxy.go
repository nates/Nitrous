package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"

	"h12.io/socks"
)

var (
	proxies []string
	cycle   int
)

func loadProxies() error {
	file, err := os.Open("proxies.txt")

	if err != nil {
		return errors.New("could not open proxies.txt")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err == nil {
			proxies = append(proxies, scanner.Text())
		}
	}

	if len(proxies) == 0 {
		return errors.New("no proxies")
	}

	return nil
}

func getTransport() (*http.Transport, error) {
	switch proxyType {
	case "http":
		proxy, err := getProxy()
		if err != nil {
			return nil, err
		}
		parsed, err := url.Parse(fmt.Sprintf("http://%v", proxy))
		if err != nil {
			return nil, err
		}
		return &http.Transport{
			Proxy: http.ProxyURL(parsed),
		}, nil
	case "socks4":
		proxy, err := getProxy()
		if err != nil {
			return nil, err
		}
		dial := socks.Dial(fmt.Sprintf("socks4://%v", proxy))
		return &http.Transport{
			Dial: dial,
		}, nil
	case "socks5":
		proxy, err := getProxy()
		if err != nil {
			return nil, err
		}
		dial := socks.Dial(fmt.Sprintf("socks5://%v", proxy))
		return &http.Transport{
			Dial: dial,
		}, nil
	}
	return nil, errors.New("unkown error")
}

func getProxy() (string, error) {
	if proxyMethod == "cycle" {
		if cycle >= len(proxies) {
			cycle = 0
		}
		proxy := proxies[cycle]
		cycle++
		return proxy, nil
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(proxies))))
	if err != nil {
		return "", err
	}
	proxy := proxies[int(nBig.Int64())]
	return proxy, nil
}
