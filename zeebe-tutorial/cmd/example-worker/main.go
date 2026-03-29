// Example Zeebe job workers for tutorials:
//   - decision       — Service Task «Обработать автоматически» (Type = decision) в application.bpmn / Reunico
//   - example-task   — минимальный процесс bpmn/examples/example-task.bpmn
//
// Запуск (локальный Zeebe):
//
//	ZEEBE_ADDRESS=127.0.0.1:26500 go run ./cmd/example-worker
//
// Тест с суммой < 5000: переменные процесса { "name": "Elis", "amount": 3000 } — токен попадает на job decision.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"go.uber.org/zap"

	"zeebe-tutorial/internal/config"
	"zeebe-tutorial/internal/tutorial"
	zclient "zeebe-tutorial/internal/zeebe"
)

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	jobs := []struct {
		jobType string
		handler worker.JobHandler
	}{
		{tutorial.JobTypeDecision, decision.Handle},
		{tutorial.JobTypeExampleTask, ex.Handle},
	}

	var opened []worker.JobWorker
	for _, j := range jobs {
		w := client.NewJobWorker().
			JobType(j.jobType).
			Handler(j.handler).
			Name("example-worker-" + j.jobType).
			Open()
		opened = append(opened, w)
		logger.Info("job worker opened", zap.String("job_type", j.jobType))
	}

	logger.Info("example-worker running",
		zap.Strings("job_types", []string{tutorial.JobTypeDecision, tutorial.JobTypeExampleTask}),
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
