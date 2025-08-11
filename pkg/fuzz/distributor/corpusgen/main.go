package main

import (
	"flag"
	"fmt"
	"github.com/grafana/pyroscope/pkg/pprof"
	"os"
	"path/filepath"
	"strings"
)

var testdata = flag.String("testdata", "", "path to testdata directory")
var corpus = flag.String("corpus", "", "path to corpus dir")

func main() {
	flag.Parse()
	if *corpus == "" {
		panic("corpus path is required")
	}
	if *testdata == "" {
		panic("testdata path is required")
	}
	profiles := loadProfilesFromTestdata()
	_ = profiles
}

func loadProfilesFromTestdata() []*pprof.Profile {
	var profiles []*pprof.Profile
	files, err := os.ReadDir(*testdata)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() || strings.HasSuffix(file.Name(), ".txt") {
			continue
		}
		var data []byte
		var p *pprof.Profile
		pp := filepath.Join(*testdata, file.Name())
		data, err = os.ReadFile(pp)
		if err != nil {
			panic(err)
		}

		p, err = pprof.RawFromBytes(data)

		if err != nil {
			panic(fmt.Errorf("could not parse %s: %w", pp, err))
		}
		fmt.Printf("loaded pprof from %s\n", pp)
		profiles = append(profiles, p)
	}
	return profiles
}
