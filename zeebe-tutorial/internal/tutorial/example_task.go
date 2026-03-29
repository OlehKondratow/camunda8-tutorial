package tutorial

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"go.uber.org/zap"
	zjobs "zeebe-tutorial/internal/zeebe"
)

// JobTypeExampleTask is the Zeebe job type for the minimal tutorial worker (BPMN service task must match).
const JobTypeExampleTask = "example-task"

type exampleInput struct {
	Name string `json:"name"`
}

type exampleOutput struct {
	Message string `json:"message"`
	OK      bool   `json:"ok"`
}

// ExampleTask is a minimal job worker for demos (Reunico-style service task).
type ExampleTask struct {
	log *zap.Logger
}

func NewExampleTask(log *zap.Logger) *ExampleTask {
	return &ExampleTask{log: log}
}

// Handle completes the job with a greeting; input variable "name" is optional (default "world").
func (h *ExampleTask) Handle(client worker.JobClient, job entities.Job) {
	ctx := context.Background()
	h.log.Info("example-task activated",
		zap.Int64("job_key", job.GetKey()),
		zap.Int64("process_instance_key", job.GetProcessInstanceKey()),
	)

	name := "world"
	raw := strings.TrimSpace(job.GetVariables())
	if raw != "" && raw != "{}" {
		var in exampleInput
		if err := json.Unmarshal([]byte(raw), &in); err != nil {
			h.fail(ctx, client, job, fmt.Errorf("variables: %w", err))
			return
		}
		if strings.TrimSpace(in.Name) != "" {
			name = strings.TrimSpace(in.Name)
		}
	}

	out, err := json.Marshal(exampleOutput{
		Message: fmt.Sprintf("Hello, %s!", name),
		OK:      true,
	})
	if err != nil {
		h.fail(ctx, client, job, err)
		return
	}

	if err := zjobs.CompleteJSON(ctx, client, job, string(out)); err != nil {
		h.fail(ctx, client, job, err)
		return
	}

	h.log.Info("example-task completed", zap.Int64("job_key", job.GetKey()))
}

func (h *ExampleTask) fail(ctx context.Context, client worker.JobClient, job entities.Job, jobErr error) {
	h.log.Error("example-task failed", zap.Error(jobErr), zap.Int64("job_key", job.GetKey()))
	if err := zjobs.Fail(ctx, client, job, jobErr); err != nil {
		h.log.Error("fail job command failed", zap.Error(err), zap.Int64("job_key", job.GetKey()))
	}
}
