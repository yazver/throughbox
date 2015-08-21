package debug

import (
	"os"
	"io"
	"log"
	"fmt"
	"io/ioutil"
)

const (
	DFile = iota
	DConsole
	DNone
)

var DebugMode = DNone
var out io.Writer = ioutil.Discard
var file *os.File = nil
var logger *log.Logger = nil

func SetDebugMode(mode int) {
	DebugMode = mode
	initialize()
}

func init() {
	initialize()
}

func initialize() {
	if file != nil {
		file.Close()
		file = nil
	}
	switch DebugMode {
	case DFile:
		var err error
		file, err = os.Create("debugFile.log")
		out = file
		if nil != err {
			panic(err.Error())
		}
	case DConsole:
		out = os.Stdout

	default:
		out = ioutil.Discard

	}
	logger = log.New(out, "[DEBUG]", log.Lshortfile)
}


// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	logger.Println(fmt.Sprintf(format, v...))
}

// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) { logger.Println(v...) }

// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) { logger.Println(v...) }

