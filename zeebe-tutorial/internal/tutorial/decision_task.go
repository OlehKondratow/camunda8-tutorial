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

// JobTypeDecision — pole Type w Modeler dla gałęzi „Obsłużyć automatycznie” (Reunico / credit-application tutorial).
const JobTypeDecision = "c8jw-golang"

type decisionInput struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type decisionOutput struct {
	Processed bool   `json:"processed"`
	Message   string `json:"message"`
	Mode      string `json:"mode"`
}

// DecisionTask handles the automatic branch service task (amount < 5000 in BPMN gateway).
type DecisionTask struct {
	log *zap.Logger
}

func NewDecisionTask(log *zap.Logger) *DecisionTask {
	return &DecisionTask{log: log}
}

// Handle completes with stub approval — replace with real rules / scoring.
func (h *DecisionTask) Handle(client worker.JobClient, job entities.Job) {
	ctx := context.Background()
	h.log.Info("decision job activated",
		zap.Int64("job_key", job.GetKey()),
		zap.Int64("process_instance_key", job.GetProcessInstanceKey()),
	)

	raw := strings.TrimSpace(job.GetVariables())
	if raw == "" || raw == "{}" {
		h.fail(ctx, client, job, fmt.Errorf("variables: expected name and amount"))
		return
	}

	var in decisionInput
	if err := json.Unmarshal([]byte(raw), &in); err != nil {
		h.fail(ctx, client, job, fmt.Errorf("variables: %w", err))
		return
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		name = "applicant"
	}

	msg := fmt.Sprintf("automatic processing for %s, amount=%v", name, in.Amount)
	out := decisionOutput{
		Processed: true,
		Message:   msg,
		Mode:      "automatic",
	}
	payload, err := json.Marshal(out)
	if err != nil {
		h.fail(ctx, client, job, err)
		return
	}

	if err := zjobs.CompleteJSON(ctx, client, job, string(payload)); err != nil {
		h.fail(ctx, client, job, err)
		return
	}

	h.log.Info("decision job completed", zap.Int64("job_key", job.GetKey()), zap.Float64("amount", in.Amount))
}

func (h *DecisionTask) fail(ctx context.Context, client worker.JobClient, job entities.Job, jobErr error) {
	h.log.Error("decision job failed", zap.Error(jobErr), zap.Int64("job_key", job.GetKey()))
	if err := zjobs.Fail(ctx, client, job, jobErr); err != nil {
		h.log.Error("fail job command failed", zap.Error(err), zap.Int64("job_key", job.GetKey()))
	}
}
