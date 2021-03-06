// +build ignore

package main

import (
	"log"
	"os"
	"text/template"
)

var newMeasurementTemplate = template.Must(template.New("open_report").Parse(`
func TestExperimentNewMeasurement{{ .Name }}WorksAsIntended(t *testing.T) {
	sess := &Session{location: &model.LocationInfo{
		{{ .LocationInfoName }}: {{ .LocationInfoValue }},
	}}
	newExperiment := func() *Experiment {
		builder, err := sess.NewExperimentBuilder("example")
		if err != nil {
			t.Fatal(err)
		}
		return builder.NewExperiment()
	}
	t.Run("with false setting", func(t *testing.T) {
		sess.privacySettings.{{ .Setting }} = false
		exp := newExperiment()
		m := exp.newMeasurement("")
		if m.{{ .NameForTesting }} != model.Default{{ .Name }} {
			t.Fatal("not the value we expected")
		}
	})
	t.Run("with true setting", func(t *testing.T) {
		sess.privacySettings.{{ .Setting }} = true
		exp := newExperiment()
		m := exp.newMeasurement("")
		if m.{{ .NameForTesting }} != {{ .ValueForTesting }} {
			t.Fatal("not the value we expected")
		}
	})
}`))

var openReportTemplate = template.Must(template.New("new_measurement").Parse(`
func TestExperimentOpenReport{{ .Name }}WorksAsIntended(t *testing.T) {
	sess := &Session{location: &model.LocationInfo{
		{{ .LocationInfoName }}: {{ .LocationInfoValue }},
	}}
	newExperiment := func() *Experiment {
		builder, err := sess.NewExperimentBuilder("example")
		if err != nil {
			t.Fatal(err)
		}
		return builder.NewExperiment()
	}
	t.Run("with false setting", func(t *testing.T) {
		sess.privacySettings.{{ .Setting }} = false
		exp := newExperiment()
		rt := exp.newReportTemplate()
		if rt.{{ .NameForTesting }} != model.Default{{ .Name }} {
			t.Fatal("not the value we expected")
		}
	})
	t.Run("with true setting", func(t *testing.T) {
		sess.privacySettings.{{ .Setting }} = true
		exp := newExperiment()
		rt := exp.newReportTemplate()
		if rt.{{ .NameForTesting }} != {{ .ValueForTesting }} {
			t.Fatal("not the value we expected")
		}
	})
}`))

type Variable struct {
	LocationInfoName  string
	LocationInfoValue string
	Name              string
	NameForTesting    string
	NoOpenReport      bool
	Setting           string
	ValueForTesting   string
}

var Variables = []Variable{{
	LocationInfoName:  "ProbeIP",
	LocationInfoValue: `"8.8.8.8"`,
	Name:              "ProbeIP",
	NameForTesting:    "ProbeIP",
	NoOpenReport:      true,
	Setting:           "IncludeIP",
	ValueForTesting:   `"8.8.8.8"`,
}, {
	LocationInfoName:  "ASN",
	LocationInfoValue: `30722`,
	Name:              "ProbeASNString",
	NameForTesting:    "ProbeASN",
	Setting:           "IncludeASN",
	ValueForTesting:   `"AS30722"`,
}, {
	LocationInfoName:  "CountryCode",
	LocationInfoValue: `"IT"`,
	Name:              "ProbeCC",
	NameForTesting:    "ProbeCC",
	Setting:           "IncludeCountry",
	ValueForTesting:   `"IT"`,
}, {
	LocationInfoName:  "NetworkName",
	LocationInfoValue: `"Vodafone Italia"`,
	Name:              "ProbeNetworkName",
	NameForTesting:    "ProbeNetworkName",
	NoOpenReport:      true,
	Setting:           "IncludeASN",
	ValueForTesting:   `"Vodafone Italia"`,
}, {
	LocationInfoName:  "ResolverIP",
	LocationInfoValue: `"9.9.9.9"`,
	Name:              "ResolverIP",
	NameForTesting:    "ResolverIP",
	NoOpenReport:      true,
	Setting:           "IncludeIP",
	ValueForTesting:   `"9.9.9.9"`,
}, {
	LocationInfoName:  "ResolverASN",
	LocationInfoValue: `44`,
	Name:              "ResolverASNString",
	NameForTesting:    "ResolverASN",
	NoOpenReport:      true,
	Setting:           "IncludeASN",
	ValueForTesting:   `"AS44"`,
}, {
	LocationInfoName:  "ResolverNetworkName",
	LocationInfoValue: `"Google LLC"`,
	Name:              "ResolverNetworkName",
	NameForTesting:    "ResolverNetworkName",
	NoOpenReport:      true,
	Setting:           "IncludeASN",
	ValueForTesting:   `"Google LLC"`,
}}

func writestring(fp *os.File, s string) {
	if _, err := fp.Write([]byte(s)); err != nil {
		log.Fatal(err)
	}
}

func withFile(filepath string, do func(fp *os.File)) {
	fp, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	do(fp)
	if err := fp.Close(); err != nil {
		log.Fatal(err)
	}
}

func writeExperimentGeneratedTestGo() {
	withFile("experiment_generated_test.go", func(fp *os.File) {
		writestring(fp, "// Code generated by go generate; DO NOT EDIT.\n")
		writestring(fp, "\n")
		writestring(fp, "package engine\n")
		writestring(fp, "\n")
		writestring(fp, "import (\n")
		writestring(fp, "\t\"testing\"\n")
		writestring(fp, "\n")
		writestring(fp, "\t\"github.com/ooni/probe-engine/model\"\n")
		writestring(fp, ")\n")
		writestring(fp, "\n")
		writestring(fp, "//go:generate go run generate_experiment.go")
		writestring(fp, "\n")
		for _, variable := range Variables {
			if err := newMeasurementTemplate.Execute(fp, variable); err != nil {
				log.Fatal(err)
			}
			writestring(fp, "\n")
			if variable.NoOpenReport {
				continue
			}
			if err := openReportTemplate.Execute(fp, variable); err != nil {
				log.Fatal(err)
			}
			writestring(fp, "\n")
		}
	})
}

func main() {
	writeExperimentGeneratedTestGo()
}
