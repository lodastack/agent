package trace

import (
	"fmt"
	"path/filepath"

	"github.com/jaegertracing/jaeger/cmd/agent/app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultMetricsBackend = "prometheus"
	defaultMetricsRoute   = "/metrics"
)

var defaultProcessors = []struct {
	model    app.Model
	protocol app.Protocol
	hostPort string
}{
	{model: "jaeger", protocol: "compact", hostPort: "localhost:6831"},
	{model: "jaeger", protocol: "binary", hostPort: "localhost:6832"},
}

const (
	defaultQueueSize     = 1000
	defaultMaxPacketSize = 65000
	defaultServerWorkers = 10
	defaultMinPeers      = 3

	defaultHTTPServerHostPort = "localhost:5778"
)

// Tracer global
var Tracer *tracer

// Tracer custom struct
type tracer struct {
	agent   *app.Agent
	started bool
}

// Start trace service
func Start(collector []string, logDir string) error {
	if Tracer != nil && Tracer.started {
		return nil
	}
	Tracer = &tracer{}
	if len(collector) < 1 {
		return fmt.Errorf("config collector address first: %d", len(collector))
	}

	conf := zap.NewProductionConfig()
	var level zapcore.Level
	err := (&level).UnmarshalText([]byte("info"))
	if err != nil {
		return err
	}
	conf.Level = zap.NewAtomicLevelAt(level)
	conf.OutputPaths = []string{filepath.Join(logDir, "INFO.log")}
	logger, err := conf.Build()
	if err != nil {
		return err
	}

	builder := &app.Builder{}
	builder.Metrics.Backend = defaultMetricsBackend
	builder.Metrics.HTTPRoute = defaultMetricsRoute

	for _, processor := range defaultProcessors {
		p := &app.ProcessorConfiguration{Model: processor.model, Protocol: processor.protocol}
		p.Workers = defaultServerWorkers
		p.Server.QueueSize = defaultQueueSize
		p.Server.MaxPacketSize = defaultMaxPacketSize
		p.Server.HostPort = processor.hostPort
		builder.Processors = append(builder.Processors, *p)
	}

	builder.CollectorHostPorts = collector

	builder.HTTPServer.HostPort = defaultHTTPServerHostPort
	builder.DiscoveryMinPeers = defaultMinPeers

	Tracer.agent, err = builder.CreateAgent(logger)
	if err != nil {
		return fmt.Errorf("Unable to initialize Jaeger Agent: %s", err)
	}

	// start tracer
	if err := Tracer.agent.Run(); err != nil {
		return fmt.Errorf("Failed to run the trace module: %s", err)
	}
	Tracer.started = true
	return nil
}
