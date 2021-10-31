package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

type Config struct {
	configFile string
}

func main() {
	// https://qiita.com/kumatch/items/258d7984c0270f6dd73a
	// https://github.com/prometheus/prometheus/blob/c2d1c858577cdac009162d095ec80402b049949c/cmd/prometheus/main.go#L186-L338
	cfg := Config{}
	a := kingpin.New("test", "help")
	a.Version("v1.0.0")
	a.HelpFlag.Short('h')
	a.Flag("config.file", "Prometheus configuration file path.").
		Default("prometheus.yml").StringVar(&cfg.configFile)
	a.Parse(os.Args[1:])
	fmt.Printf("configFile: %s", cfg.configFile)
}
