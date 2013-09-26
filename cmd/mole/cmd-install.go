// +build darwin linux

package main

import (
	"github.com/calmh/mole/table"
	"github.com/jessevdk/go-flags"
	"io"
	"io/ioutil"
	"os/exec"
	"runtime"
)

type installCommand struct{}

var installParser *flags.Parser

func init() {
	cmd := installCommand{}
	installParser = globalParser.AddCommand("install", msgInstallShort, msgInstallLong, &cmd)
}

func (c *installCommand) Usage() string {
	return "[package] [install-OPTIONS]"
}

func (c *installCommand) Execute(args []string) error {
	setup()

	cl := NewClient(serverIni.address, serverIni.fingerprint)
	if len(args) == 0 {
		pkgMap, err := cl.Packages()
		fatalErr(err)

		arch := runtime.GOOS + "-" + runtime.GOARCH
		var rows [][]string
		for _, pkg := range pkgMap[arch] {
			rows = append(rows, []string{pkg.Package, pkg.Description})
		}

		if len(rows) > 0 {
			rows = append([][]string{{"PKG", "DESCRIPTION"}}, rows...)
			infoln(table.Fmt("ll", rows))
		} else {
			infoln(msgNoPackages)
		}
	} else {
		requireRoot("install")
		name := args[0]
		fullname := args[0] + "-" + runtime.GOOS + "-" + runtime.GOARCH + ".tar.gz"
		wr, err := ioutil.TempFile("", name)
		fatalErr(err)

		rd, err := cl.Package(fullname)
		fatalErr(err)

		_, err = io.Copy(wr, rd)
		fatalErr(err)
		err = wr.Close()
		fatalErr(err)
		_ = rd.Close()

		td, err := ioutil.TempDir("", name)
		fatalErr(err)

		cmd := exec.Command("sh", "-c", "cd "+td+" && tar zxf "+wr.Name())
		_, err = cmd.CombinedOutput()
		fatalErr(err)

		becomeRoot()
		cmd = exec.Command(td+"/install.sh", td)
		_, err = cmd.CombinedOutput()
		fatalErr(err)
		dropRoot()

		okln("Installed")
	}

	return nil
}