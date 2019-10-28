package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	inv "github.com/redhat-cop/dash/pkg/inventory"
)

var version = "undefined"
var invPath string

func init() {
	flag.StringVar(&invPath, "i", "./", "Path to Inventory, relative or absolute")
	flag.Parse()
}

func main() {

	var i inv.Inventory
	var ns string

	yamlFile, err := ioutil.ReadFile(invPath + "dash.yaml")
	if err != nil {
		fmt.Printf("Error: Couldn't load dash inventory: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	i.Load(yamlFile, invPath)
	err = i.Process(&ns)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		flag.Usage()
		os.Exit(1)
	}

}
