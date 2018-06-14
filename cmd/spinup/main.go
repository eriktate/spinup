package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/eriktate/spinup/gen"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Flags are possible command flags that can be passed in to spinup.
type Flags struct {
	Impl  string
	Pkg   string
	PB    string
	Print bool
}

var flags Flags

func init() {
	flag.StringVar(&flags.Impl, "impl", "", "Name of the struct implementing service")
	flag.StringVar(&flags.Pkg, "pkg", "", "Name of the generated package")
	flag.StringVar(&flags.PB, "pb", "", "Path of .pb.go file")
	flag.BoolVar(&flags.Print, "print", false, "Print output file to stdout")
}

func main() {
	logrus.SetOutput(os.Stdout)
	log := logrus.New()

	flag.Parse()

	if err := genService(log, flags); err != nil {
		log.Error(err)
	}
}

func genService(log *logrus.Logger, f Flags) error {
	if f.Impl == "" {
		return errors.New("You must include an impl!")
	}

	if f.Pkg == "" {
		return errors.New("You must include a pkg!")
	}

	if f.PB == "" {
		// TODO: Maybe default to searching the current directory.
		return errors.New("You must include a .pb.go path!")
	}

	def := gen.GenerateServiceDef(f.PB)

	def.Implementor = f.Impl
	def.PackageName = f.Pkg

	serviceTemplate, err := ioutil.ReadFile("templates/service.gotemplate")
	if err != nil {
		return errors.Wrap(err, "Failed to load service template")
	}

	tmpl := template.New("service")
	tmpl, err = tmpl.Parse(string(serviceTemplate))
	if err != nil {
		return errors.Wrap(err, "Failed to parse service template")
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, &def); err != nil {
		return errors.Wrap(err, "Failed to execute template")
	}

	if f.Print {
		fmt.Fprintln(os.Stdout, buf.String())
	}

	return nil
}
