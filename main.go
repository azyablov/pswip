package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

var wg sync.WaitGroup = sync.WaitGroup{}
var appLogger *log.Logger

type AppParams struct {
	NumOfPackets int `json:"numofpackets"`
}

type AppConfig struct {
	Targets []string  `json:"targets"`
	Params  AppParams `json:"params"`
}

func main() {
	// Parsing arguments

	var logFile string
	var configFile string
	flag.StringVar(&configFile, "c", "pswip.conf", "application config file.")
	flag.StringVar(&logFile, "l", "pswip.log", "log file.")
	flag.Parse()

	// Logging
	logfh, err := os.Create(logFile)
	if err != nil {
		fmt.Println(err, "\nCan't create log file, switching to STDOUT.")
	}
	appLogger = log.New(logfh, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	appLogger.Println("App started...")

	// Reading config file
	ac := new(AppConfig)
	err = openJSONConfig(ac, configFile)
	if err != nil {
		appLogger.Fatalln(err)
	}
	if len(ac.Targets) == 0 {
		appLogger.Println("Nothing to do, since sero targets were identified.")
		os.Exit(0)
	}

	targets := ac.Targets
	// Setup channels for stats
	type statsCh chan *ping.Statistics
	slRCh := make([]statsCh, len(targets))
	for i := range slRCh {
		slRCh[i] = make(chan *ping.Statistics, 1)
	}

	// channel for errors
	type errCh chan error
	slErrCh := make([]errCh, len(targets))
	for i := range slErrCh {
		slErrCh[i] = make(chan error, 1)
	}

	// WG setup
	for i, t := range targets {
		wg.Add(1)
		go pingTarget(t, true, ac.Params.NumOfPackets, slRCh[i], slErrCh[i])
	}
	wg.Wait()

	// Retriving results
	for i := 0; i < len(targets); i++ {
		err = <-slErrCh[i]
		if err != nil {
			appLogger.Printf("%v\n", err)
			continue
		}
		err = printStats(<-slRCh[i], targets[i])
		if err != nil {
			appLogger.Println(err)
		}
	}
}

func pingTarget(t string, p bool, numOfPackets int, ch chan<- *ping.Statistics, errCh chan<- error) {
	defer wg.Done()
	defer appLogger.Printf("pingTarget: %v handling is done", t)
	pinger, err := ping.NewPinger(t)
	if err != nil {
		errCh <- fmt.Errorf("pingTarget: %v", err)
		ch <- nil
		return
	}

	pinger.SetPrivileged(p)
	if !pinger.Privileged() {
		errCh <- fmt.Errorf("pingTarget: can't enable privilege mode")
		ch <- nil
		return
	}
	pinger.Timeout = time.Duration(4000000000)
	pinger.Count = numOfPackets

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		errCh <- fmt.Errorf("pingTarget: %v", err)
		ch <- nil
		return
	}

	ch <- pinger.Statistics()
	errCh <- nil
}

func printStats(stats *ping.Statistics, t string) error {
	var color string
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"
	switch {
	case stats.PacketLoss == 0.0:
		color = colorGreen
	case stats.PacketLoss == 100.0:
		color = colorRed
	case stats.PacketLoss < 100.0:
		color = colorYellow
	default:
		return fmt.Errorf("printStats: unexpected stats value for PacketLoss")
	}
	fmt.Printf("%vTarget: %v\nIP@: %v\n", string(color), t, stats.IPAddr)
	fmt.Printf("%v %d packets transmitted, %d packets received, %v%% packet loss\n",
		string(color), stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
	fmt.Printf("%v round-trip min/avg/max/stddev = %v/%v/%v/%v\n%v",
		string(color), stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt, string(colorReset))
	fmt.Println(strings.Repeat("=", 80))
	return nil
}

func openJSONConfig(ac *AppConfig, fileName string) error {

	// Opening file
	fh, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("openJSONConfig: can't open specifies file %v: %v", fileName, err)
	}
	defer fh.Close()

	// Reading contents
	fileStream, err := ioutil.ReadAll(fh)
	if err != nil {
		return fmt.Errorf("openJSONConfig: can't read contents of %v: %v", fh.Name(), err)
	}

	err = json.Unmarshal(fileStream, ac)
	if err != nil {
		return fmt.Errorf("openJSONConfig: can't unmarshall file %v: %v", fh.Name(), err)
	}

	return nil
}
