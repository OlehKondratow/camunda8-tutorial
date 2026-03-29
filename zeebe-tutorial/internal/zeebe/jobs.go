package zeebe

import (
	"context"
	"fmt"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
)

// CompleteJSON completes a job with JSON variables (merge into process scope).
func CompleteJSON(ctx context.Context, client worker.JobClient, job entities.Job, variablesJSON string) error {
	cmd, err := client.NewCompleteJobCommand().
		JobKey(job.GetKey()).
		VariablesFromString(variablesJSON)
	if err != nil {
		return fmt.Errorf("complete job variables: %w", err)
	}
	if _, err := cmd.Send(ctx); err != nil {
		return fmt.Errorf("complete job: %w", err)
	}
	return nil
}

// Fail reports job failure with decremented retries.
func Fail(ctx context.Context, client worker.JobClient, job entities.Job, jobErr error) error {
	retries := job.GetRetries() - 1
	if retries < 0 {
		retries = 0
	}
	_, err := client.NewFailJobCommand().
		JobKey(job.GetKey()).
		Retries(retries).
		ErrorMessage(jobErr.Error()).
		Send(ctx)
	if err != nil {
		return fmt.Errorf("fail job: %w", err)
	}
	return nil
}
