package main

import (
	"flag"
	"github.com/GoSploit/obscurate"
	"io/ioutil"
)

var Ops = flag.Int("ops", 20, "Operations to perform")

var Obs = flag.Bool("obs", true, "Create obscuration code")
var ObsName = flag.String("obsname", "obs", "Name of the obscuration routine")

var Deobs = flag.Bool("deobs", true, "Create deobscuration code")
var DeobsName = flag.String("deobsname", "deobs", "Name of the deobscuration routine")

var Package = flag.String("package", "main", "Package to build the files as a part of")
var File = flag.String("file", "obs.go", "File to save code in")

func main() {
	flag.Parse()
	key, _ := obscurate.GenerateKey(*Ops)
	filecode := "package " + *Package
	if *Obs {
		filecode += "\n\n" + key.ObscurateFunc(*ObsName)
	}
	if *Deobs {
		filecode += "\n\n" + key.DeobscurateFunc(*DeobsName)
	}

	ioutil.WriteFile(*File, []byte(filecode), 0666)
}
