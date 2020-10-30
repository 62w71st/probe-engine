// Package httphostheader contains the HTTP host header network experiment.
//
// This experiment has not been specified yet. It is nonetheless available for testing
// and as a building block that other experiments could reuse.
package httphostheader

import (
	"context"
	"errors"
	"fmt"

	"github.com/ooni/probe-engine/experiment/urlgetter"
	"github.com/ooni/probe-engine/model"
)

const (
	testName    = "http_host_header"
	testVersion = "0.1.0"
)

// Config contains the experiment config.
type Config struct {
	// TestHelperURL is the address of the test helper.
	TestHelperURL string
}

// TestKeys contains httphost test keys.
type TestKeys struct {
	urlgetter.TestKeys
	THAddress string `json:"th_address"`
}

// Measurer performs the measurement.
type Measurer struct {
	config Config
}

// ExperimentName implements ExperimentMeasurer.ExperiExperimentName.
func (m *Measurer) ExperimentName() string {
	return testName
}

// ExperimentVersion implements ExperimentMeasurer.ExperimentVersion.
func (m *Measurer) ExperimentVersion() string {
	return testVersion
}

// Run implements ExperimentMeasurer.Run.
func (m *Measurer) Run(
	ctx context.Context,
	sess model.ExperimentSession,
	measurement *model.Measurement,
	callbacks model.ExperimentCallbacks,
) error {
	if measurement.Input == "" {
		return errors.New("Experiment requires measurement.Input")
	}
	if m.config.TestHelperURL == "" {
		m.config.TestHelperURL = "http://www.example.com"
	}
	urlgetter.RegisterExtensions(measurement)
	g := urlgetter.Getter{
		Begin: measurement.MeasurementStartTimeSaved,
		Config: urlgetter.Config{
			HTTPHost: string(measurement.Input),
		},
		Session: sess,
		Target:  fmt.Sprintf(m.config.TestHelperURL),
	}
	tk, _ := g.Get(ctx)
	measurement.TestKeys = tk
	return nil
}

// NewExperimentMeasurer creates a new ExperimentMeasurer.
func NewExperimentMeasurer(config Config) model.ExperimentMeasurer {
	return &Measurer{config: config}
}
