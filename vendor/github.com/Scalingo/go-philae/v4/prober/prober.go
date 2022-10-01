package prober

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/errgo.v1"
)

// Probe define a minimal set of methods that a probe should implement
type Probe interface {
	Name() string
	Check() error
}

// Prober entrypoint of the philae api. It will retain a set of probe and run
// checks when asked to
type Prober struct {
	timeout time.Duration
	probes  map[string]Probe
}

// ProberOption is a function modifying some parameters of the Prober
type ProberOption func(p *Prober)

// WithTimeout is a ProberOption which defines a timeout the prober have to get
// executed into Recommandation: it should be higher than the timeout of the
// probes included in it, otherwise it would mask the real errors.
func WithTimeout(d time.Duration) ProberOption {
	return ProberOption(func(p *Prober) {
		p.timeout = d
	})
}

type ProberError struct {
	errs []string
}

func (e *ProberError) AddError(name string, err error) {
	msg := fmt.Sprintf("probe %s: %s", name, err.Error())
	e.errs = append(e.errs, msg)
}

func (e *ProberError) IsEmpty() bool {
	return e.errs == nil
}

func (e *ProberError) Error() string {
	if e.errs == nil {
		return "no error"
	}
	b := strings.Builder{}
	b.WriteString("prober error: ")
	b.WriteString(strings.Join(e.errs, ", "))
	return b.String()
}

// ErrProbeNotFound is emitted when a check is performed on a probe that was not
// added to the prober
var ErrProbeNotFound = errors.New("probe not found")

// Result is the data structure used to retain the data fetched from a single run of each probes
type Result struct {
	Healthy bool           `json:"healthy"`
	Probes  []*ProbeResult `json:"probes"`
	Error   error          `json:"error"`
}

// ProbeResult is the data structure used to retain the data fetched from a single probe
type ProbeResult struct {
	Name     string        `json:"name"`
	Healthy  bool          `json:"healthy"`
	Comment  string        `json:"comment"`
	Error    error         `json:"error"`
	Duration time.Duration `json:"duration"`
}

// NewProber is the default constructor of a Prober
func NewProber(opts ...ProberOption) *Prober {
	p := &Prober{
		timeout: 10 * time.Second,
		probes:  map[string]Probe{},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Prober) AddProbe(probe Probe) {
	p.probes[probe.Name()] = probe
}

// Check will run the check of each probes added and return the result in a Result struct
func (p *Prober) Check(ctx context.Context) *Result {
	probesResults := make([]*ProbeResult, len(p.probes))
	resultChan := make(chan *ProbeResult, len(p.probes))
	healthy := true
	proberErr := &ProberError{}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	for _, probe := range p.probes {
		go p.checkOneProbe(ctx, probe, resultChan)
	}

	for i := 0; i < len(p.probes); i++ {
		probeResult := <-resultChan
		if !probeResult.Healthy {
			healthy = false
			proberErr.AddError(probeResult.Name, probeResult.Error)
		}
		probesResults[i] = probeResult
	}

	var err error
	if !proberErr.IsEmpty() {
		err = proberErr
	}

	return &Result{
		Healthy: healthy,
		Probes:  probesResults,
		Error:   err,
	}
}

func (p *Prober) CheckOneProbe(ctx context.Context, probeName string) *ProbeResult {
	probe, ok := p.probes[probeName]
	if !ok {
		return &ProbeResult{
			Error: ErrProbeNotFound,
		}
	}

	resultChan := make(chan *ProbeResult, 1)
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	go p.checkOneProbe(ctx, probe, resultChan)
	probeResult := <-resultChan

	return probeResult
}

func (p *Prober) checkOneProbe(ctx context.Context, probe Probe, res chan *ProbeResult) {
	probeRes := make(chan error)
	var err error

	begin := time.Now()
	go ProberWrapper(ctx, probe, probeRes)

	select {
	case e := <-probeRes:
		err = e
	case <-ctx.Done():
		err = fmt.Errorf("prober: %s", ctx.Err())
	}

	probe_healthy := true
	duration := time.Now().Sub(begin)
	comment := fmt.Sprintf("took %v", duration)
	if err != nil {
		err = errgo.Notef(err, "probe check failed")
		comment = "error"
		probe_healthy = false
	}
	probeResult := &ProbeResult{
		Name:     probe.Name(),
		Healthy:  probe_healthy,
		Comment:  comment,
		Error:    err,
		Duration: duration,
	}

	res <- probeResult
}

//nolint:revive
func ProberWrapper(ctx context.Context, probe Probe, res chan error) {
	err := probe.Check()
	res <- err
}
