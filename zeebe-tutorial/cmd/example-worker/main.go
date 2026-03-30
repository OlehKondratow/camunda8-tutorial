// Example Zeebe job workers for tutorials:
//   - c8jw-golang    — Service Task „Obsłużyć automatycznie” (Type = c8jw-golang) w application.bpmn / Reunico
//   - c8jw-python    — minimalny proces bpmn/examples/example-task.bpmn
//
// Uruchomienie (lokalny Zeebe):
//
//	ZEEBE_ADDRESS=127.0.0.1:26500 go run ./cmd/example-worker
//
// Test z kwotą < 5000: zmienne { "name": "Elis", "amount": 3000 } — token trafia na job c8jw-golang.
package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"go.uber.org/zap"

	"zeebe-tutorial/internal/config"
	"zeebe-tutorial/internal/tutorial"
	zclient "zeebe-tutorial/internal/zeebe"
)

func jobTypesFromEnv() []string {
	raw := strings.TrimSpace(os.Getenv("JOB_TYPES"))
	if raw == "" {
		return []string{tutorial.JobTypeDecision, tutorial.JobTypeExampleTask}
	}
	var out []string
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{tutorial.JobTypeDecision, tutorial.JobTypeExampleTask}
	}
	return out
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("load config", zap.Error(err))
	}

	client, err := zclient.NewClient(cfg)
	if err != nil {
		logger.Fatal("create zeebe client", zap.Error(err))
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			logger.Error("close zeebe client", zap.Error(closeErr))
		}
	}()

	decision := tutorial.NewDecisionTask(logger)
	ex := tutorial.NewExampleTask(logger)

	handlers := map[string]worker.JobHandler{
		tutorial.JobTypeDecision:    decision.Handle,
		tutorial.JobTypeExampleTask: ex.Handle,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	want := jobTypesFromEnv()
	var opened []worker.JobWorker
	for _, jt := range want {
		h, ok := handlers[jt]
		if !ok {
			logger.Warn("unknown JOB_TYPES entry (ignored)", zap.String("job_type", jt))
			continue
		}
		w := client.NewJobWorker().
			JobType(jt).
			Handler(h).
			Name(cfg.WorkerName + "-" + jt).
			Open()
		opened = append(opened, w)
		logger.Info("job worker opened", zap.String("job_type", jt))
	}
	if len(opened) == 0 {
		logger.Fatal("no job workers started — check JOB_TYPES")
	}

	logger.Info("example-worker running",
		zap.Strings("job_types", want),
		zap.String("gateway", cfg.ZeebeAddress),
	)

	<-ctx.Done()
	logger.Info("shutdown, closing workers")
	for _, w := range opened {
		w.Close()
	}
	for _, w := range opened {
		w.AwaitClose()
	}
	logger.Info("example-worker stopped")
}
