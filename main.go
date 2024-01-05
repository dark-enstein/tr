package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/andrew-d/go-termutil"
	"github.com/dark-enstein/tr/pkg/r"
	"github.com/dark-enstein/tr/pkg/w"
	"github.com/spf13/pflag"
)

const (
	CONSOLE = iota
	STDIN
	FILE
)

var (
	DefaultRawString = "iuwidninjnfinrefnjkrvref"
	DefaultFrom      = "injn"
	DefaultTo        = "high"
)

func main() {
	test := true
	var ctx = context.Background()
	f := r.Flags{}
	switch test {
	case true:
		initFlags(&f)
		if len(f.DelString) > 0 {
			if f.Action >= 0 {
				f.Action = r.Action_DELETE
			} else {
				log.Printf("Flag action already set: %d\n", f.Action)
				return
			}
		}
		if len(f.SqueezeString) > 0 {
			fmt.Println("enter squeeze:", f.SqueezeString)
			if f.Action >= 0 {
				f.Action = r.Action_SQUEEZE
				f.SqueezeBytes = []byte(f.SqueezeString)
			} else {
				log.Printf("Flag action already set: %d\n", f.Action)
				return
			}
		}
		_main(&f, ctx)
	case false:
		initFlags(&f)
		if len(f.DelString) > 0 {
			if f.Action >= 0 {
				f.Action = r.Action_DELETE
			} else {
				log.Printf("Flag action already set: %d\n", f.Action)
				return
			}
		}
		if len(f.SqueezeString) > 0 {
			if f.Action >= 0 {
				f.Action = r.Action_SQUEEZE
				f.SqueezeBytes = []byte(f.SqueezeString)
			} else {
				log.Printf("Flag action already set: %d\n", f.Action)
				return
			}
		}
		f.DelString = "project"
		fmt.Println(_mainDebug(&f, ctx))
	}
}

// initFlags initializes the flags defined at start time
func initFlags(f *r.Flags) {
	pflag.StringVarP(&f.DelString, "delete", "d",
		"90808", "delete all occurrence of a string in input text")
	pflag.StringVarP(&f.SqueezeString, "squeeze", "s",
		"", "EXPERIMENTAL: reduce all repeated char occurence of any of char"+
			" in value"+
			" string in input text (doesn't work as intended)")
	pflag.Parse()
}

func _main(f *r.Flags, ctx context.Context) {
	rep := r.R{}
	var arg []string
	by, class := whichClass()
	if len(f.DelString) > 1 {
		rep.FlagEnabled = true
		rep.Flag = f
	}
	if len(f.SqueezeString) > 1 {
		rep.FlagEnabled = true
		rep.Flag = f
	}
	switch class {
	case CONSOLE:
		//fmt.Println("enter console")
		arg = os.Args[1:]
		b := OpenConsole()
		rep.RawBytes = b
		rep.RawString = string(b)
		rep.From = []byte(arg[0])
		rep.To = []byte(arg[1])
		rep.Churn(ctx)
		w.Write(rep.DestString)
		_main(f, ctx)
	case STDIN:
		//fmt.Println("enter stdin")
		arg = os.Args[1:]
		if len(arg) != 2 {
			log.Printf("expecting two arguments. got: %v\n", arg)
			os.Exit(1)
		}
		rep.RawBytes = by
		rep.RawString = string(by)
		rep.From = []byte(arg[0])
		rep.To = []byte(arg[1])
		rep.Churn(ctx)
		w.Write(rep.DestString)
	case FILE:
		//fmt.Println("enter file")
		fileName, err := filepath.Abs(os.Args[1])
		if err != nil {
			log.Printf("err with expanding file path: %s\n", err.Error())
			os.Exit(1)
		}
		_, err = os.Stat(fileName)
		if os.IsNotExist(err) {
			log.Printf("err file doesn't exist: %s\n",
				err.Error())
			os.Exit(1)
		}
		fileB, err := os.ReadFile(fileName)
		if err != nil {
			log.Printf("err with reading file: %s\n", err.Error())
			os.Exit(1)
		}
		rep.RawBytes = fileB
		rep.RawString = string(fileB)
		arg = os.Args[2:]
		if len(arg) != 2 {
			log.Printf("expecting two arguments. got: %v\n", arg)
			os.Exit(1)
		}
		rep.From = []byte(arg[0])
		rep.To = []byte(arg[1])
		rep.Churn(ctx)
		w.Write(rep.DestString)
	}
	return
}

func whichClass() ([]byte, int) {
	byt, err := []byte{}, errors.New("")
	switch termutil.Isatty(os.Stdin.Fd()) {
	case true:
		if !(len(os.Args) == 3) {
			log.Println("no stdin")
			return []byte(""), FILE
		} else {
			return []byte(""), CONSOLE
		}
	case false:
		//log.Println("something in stdin")
		byt, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Printf("error encountered reading Stdin: %s\n", err.Error())
		}
		byt = bytes.TrimSpace(byt)
	}
	return byt, STDIN
}

func _mainDebug(f *r.Flags, ctx context.Context) string {
	rep := r.R{}
	rep.RawString, rep.From, rep.To = DefaultRawString, []byte("a-Z"),
		[]byte("a-z")
	rep.Churn(ctx)
	var arg []string
	by, class := whichClass()
	if len(f.DelString) > 1 {
		rep.FlagEnabled = true
		rep.Flag = f
	}
	if len(f.SqueezeString) > 1 {
		rep.FlagEnabled = true
		rep.Flag = f
	}
	switch class {
	case CONSOLE:
		//fmt.Println("enter console")
		arg = os.Args[1:]
		b := OpenConsole()
		rep.RawBytes = b
		rep.RawString = string(b)
		rep.From = []byte(arg[0])
		rep.To = []byte(arg[1])
		rep.Churn(ctx)
		w.Write(rep.DestString)
		_main(f, ctx)
	case STDIN:
		//fmt.Println("enter stdin")
		arg = os.Args[1:]
		if len(arg) != 2 {
			log.Printf("expecting two arguments. got: %v\n", arg)
			os.Exit(1)
		}
		rep.RawBytes = by
		rep.RawString = string(by)
		rep.From = []byte(arg[0])
		rep.To = []byte(arg[1])
		rep.Churn(ctx)
		w.Write(rep.DestString)
	case FILE:
		//fmt.Println("enter file")
		fileName, err := filepath.Abs(os.Args[1])
		if err != nil {
			log.Printf("err with expanding file path: %s\n", err.Error())
			os.Exit(1)
		}
		_, err = os.Stat(fileName)
		if os.IsNotExist(err) {
			log.Printf("err file doesn't exist: %s\n",
				err.Error())
			os.Exit(1)
		}
		fileB, err := os.ReadFile(fileName)
		if err != nil {
			log.Printf("err with reading file: %s\n", err.Error())
			os.Exit(1)
		}
		rep.RawBytes = fileB
		rep.RawString = string(fileB)
		arg = os.Args[2:]
		if len(arg) != 2 {
			log.Printf("expecting two arguments. got: %v\n", arg)
			os.Exit(1)
		}
		rep.From = []byte(arg[0])
		rep.To = []byte(arg[1])
		rep.Churn(ctx)
		w.Write(rep.DestString)
	}
	return rep.DestString
}

func OpenConsole() []byte {
	var byt = []byte{}
	bufScan := bufio.NewScanner(os.Stdin)
cursor:
	for bufScan.Scan() {
		byt = bufScan.Bytes()
		if len(byt) == 0 {
			break cursor
		}
		return byt
	}
	return byt
}
