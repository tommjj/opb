package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/tommjj/go-opb/internal/cli"
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

	cli := cli.New(filestorage)
	err = cli.LoadConf()
	if err != nil {
		err := cli.Reset()
		if err != nil {
			log.Fatalln(err)
		}
	}

	args := os.Args[1:]

	if len(args) == 0 {
		cli.OpenDefault()
		return
	}

	switch args[0] {
	case "-h", "--help":
		showHelp()

	case "conf":
		showConf(cli.Config)

	case "set":
		err = cli.Set(args[1:]...)
		if err != nil {
			log.Fatal(err)
		}

	case "del":
		err = cli.Del(args[1:]...)
		if err != nil {
			log.Fatal(err)
		}

	default:
		flags, a := cutFlags(args...)

		err = cli.Open(a, flags)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// cutFlags
func cutFlags(args ...string) ([]string, []string) {
	flags := []string{}

	for idx, v := range args {
		if isFlag(v) {
			flag, ok := shortFlags[v]
			if !ok {
				flag = v
			}

			flags = append(flags, flag)
		} else {
			return flags, args[idx:]
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
	fmt.Println(`  del link 	[key] 		# delete quick link.	example:"del link ex"`)
	fmt.Println(`  del search 	[key] 		# delete search link.	example:"del search ex"`)
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
