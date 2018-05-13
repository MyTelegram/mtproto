package main

import (
	"fmt"
	"github.com/cjongseok/mtproto"
	"github.com/cjongseok/slog"
	"os"
	"sync"
)

func usage() {
	fmt.Println("USAGE: ./playdump <AUTH> <DUMP>")
	fmt.Println("")
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		usage()
		os.Exit(1)
	}
}

type callback struct{}

func (cb *callback) OnUpdate(u mtproto.Update) {
	fmt.Println("update:", slog.Stringify(u))
}

func main() {
	if len(os.Args) != 3 {
		usage()
		return
	}

	authFileName := os.Args[1]
	dumpFilename := os.Args[2]
	slog.SetLogOutputAsFile("dump_" + dumpFilename + "_")
	handlerWg := sync.WaitGroup{}
	handlerWg.Add(1)
	out := make(chan interface{})
	go func() {
		defer handlerWg.Done()
		for {
			select {
			case u := <-out:
				if u == nil {
					return
				}
				fmt.Println("update:", slog.Stringify(u))
			}
		}
	}()
	dump, err := mtproto.NewDump(authFileName, dumpFilename, out)
	handleError(err)
	dump.Play()
	dump.Wait()
	handlerWg.Wait()
}
