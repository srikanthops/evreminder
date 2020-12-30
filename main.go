package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	ev "github.com/srikanthops/evreminder/event"
	"gopkg.in/yaml.v2"
)

type evServer struct {
	Enabled bool   `yaml:"enabled"` // enabled http server
	Port    string `yaml:"port"`    // http listen port
	Verbose bool   `yaml:"verbose"`
}

// Configuration Exported
type Configuration struct {
	CS evServer `yaml:"cloudrunserver"`
}

var c = Configuration{}

func loadConfig(filename string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return err
	}
	return nil
}

// Get overwritten by `-ldflags "-X main.version=${VERSION}"`
var (
	Version = "dev"
)

var (
	flags        = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	help         = flags.Bool("help", false, "Print usage instructions and exit.")
	printVersion = flags.Bool("version", false, "Print version")
	v            = flags.Bool("v", false, "Enable verbose output.")
	d            = flags.Bool("d", false, "Enable debug output.")
	s            = flags.Bool("s", false, "Enable server mode")
	configfile   = flags.String("c", "config.yaml", "ConfigFile")
	eventsfile   = flags.String("e", "events.json", "Events File in  json")
)

var conf = Configuration{}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `Usage: %s [flags] [address]`, os.Args[0])
	flags.PrintDefaults()
}

type serverConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    string `yaml:"port"`
	Verbose bool   `yaml:"verbose"`
}

func main() {
	flags.Usage = usage
	if err := flags.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Flags parse error %s", err)
		usage()
		os.Exit(1)
	}

	if *help {
		usage()
		os.Exit(0)
	}

	if *printVersion {
		_, _ = fmt.Fprintf(os.Stderr, "%s %s\n", filepath.Base(os.Args[0]), Version)
		os.Exit(0)
	}

	if err := loadConfig(*configfile); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing configfile:%s err:%s\n",
			*configfile, err)
		os.Exit(1)
	}

	if conf.CS.Port == "" {
		conf.CS.Port = "8080"
	}

	// Load Events

	if *eventsfile != "" {
		if err := ev.PopulateEvents(*eventsfile); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to populate %s err%s", *eventsfile, err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stdout, "No Events configured ")
	}

	// StartTime HTTP server.
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK\n")
	})
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", Version)
	})

	http.HandleFunc("/", todayEvents)
	http.HandleFunc("/today", todayEvents)
	http.HandleFunc("/yesterday", yEvents)
	http.HandleFunc("/tomorrow", tomEvents)
	http.HandleFunc("/month", mEvents)

	fmt.Fprintf(os.Stdout, "Started listening at port %s\n", conf.CS.Port)
	if err := http.ListenAndServe(":"+conf.CS.Port, nil); err != nil {
		log.Fatal(err)
	}
}
