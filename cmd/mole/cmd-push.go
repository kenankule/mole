package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"path/filepath"
)

type cmdPush struct{}

var pushParser *flags.Parser

func init() {
	cmd := cmdPush{}
	pushParser = globalParser.AddCommand("push", msgPushShort, msgPushLong, &cmd)
}

func (c *cmdPush) Usage() string {
	return "<tunnelfile>"
}

func (c *cmdPush) Execute(args []string) error {
	setup()

	if len(args) != 1 {
		digParser.WriteHelp(os.Stdout)
		infoln()
		fatalln("push: missing required option <tunnelfile>")
	}

	cert := certificate()
	cl := NewClient(serverIni.address, cert)
	_ = cl

	filename := filepath.Base(args[0])
	if ext := filepath.Ext(filename); ext != ".ini" {
		fatalf(msgFileNotInit, filename)
	}
	tunnelname := filename[:len(filename)-4]

	file, err := os.Open(args[0])
	fatalErr(err)

	cl.Put(tunnelname, file)
	okf(msgOkPushed, tunnelname)
	return nil
}