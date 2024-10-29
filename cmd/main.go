package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/tommjj/go-opb/internal/browser"
	"github.com/tommjj/go-opb/internal/config"
	"github.com/tommjj/go-opb/internal/filestorage"
)

var shortFlags = map[string]string{
	"-n": "--new-window",
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	filestorage := filestorage.New(path.Join(homeDir, ".opb-conf.json"))

	conf := config.New(filestorage)
	err = conf.Load()
	if err != nil {
		err := browser.Reset(conf)
		if err != nil {
			log.Fatalln(err)
		}
	}

	agr := os.Args[1:]

	if len(agr) == 0 {
		browser.OpenDefault(conf.Config)
		return
	}

	// help mode
	if agr[0] == "-h" || agr[0] == "--help" {
		showHelp()
		return
	}

	// show conf mode
	if agr[0] == "conf" {
		showConf(conf.Config)
		return
	}

	// set mode
	if agr[0] == "set" {
		err = conf.Set(agr[1:]...)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	flags, a := cutFlags(agr...)

	if len(a) == 0 {
		err := browser.OpenDefault(conf.Config, flags...)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = browser.Open(conf.Config, agr, flags)
	if err != nil {
		log.Fatal(err)
	}
}

// cutFlags
func cutFlags(agr ...string) ([]string, []string) {
	flags := []string{}

	for idx, v := range agr {
		if isFlag(v) {
			flag, ok := shortFlags[v]
			if !ok {
				flag = v
			}

			flags = append(flags, flag)
		} else {
			return flags, agr[idx:]
		}
	}
	return flags, []string{}
}

func isFlag(str string) bool {
	if len(str) == 0 {
		return false
	}
	return string(str[0]) == "-"
}

func showHelp() {
	fmt.Println("usage: ")
	fmt.Println(`  If you don't set a key it will use the default value. `)
	fmt.Println(`  [url]		# open with url`)
	fmt.Println(`  [key]		# open quick link`)
	fmt.Println(`  [key] [query]	# open search link`)
	fmt.Println()

	fmt.Println("conf:")
	fmt.Println("  If you set key as default it will be default value.")
	fmt.Println(`  set browser	[browser path]	# set browser path.	example:"set browser browser.exe"`)
	fmt.Println(`  set link 	[key] [url]	# set quick link.	example:"set link ex https://example.com"`)
	fmt.Println(`  set search 	[key] [url]	# set search link.	example:"set search ex https://example.com?q=$"`)
	fmt.Println(`  conf				# show conf`)

}

func showConf(conf config.Config) {
	fmt.Println("config:")
	fmt.Println("browser path:")
	fmt.Printf("	%v\n", conf.Browser)
	fmt.Println()

	fmt.Println("quick link:")
	for k, v := range conf.Links {
		fmt.Printf("	%v: %v\n", k, v)
	}
	fmt.Println()

	fmt.Println("search link:")
	for k, v := range conf.Search {
		fmt.Printf("	%v: %v\n", k, v)
	}
}
