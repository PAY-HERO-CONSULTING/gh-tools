package job_enqueuer

import (
	"context"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/ctxutils"
	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
	"github.com/PAY-HERO-CONSULTING/gh-tools/pools/jobs"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

const (
	BodyKey         = "body"
	ContextKey      = "_context"
	RequestIDKey    = "request_id"
	WorkerNamespace = "gh-worker"
)

type (
	JobEnqueuer interface {
		Enqueue(ctx context.Context, job jobs.Job) (string, error)
		EnqueueAfter(ctx context.Context, job jobs.Job, duration time.Duration) (string, error)
	}

	jobEnqueuer struct {
		enqueuer *work.Enqueuer
	}
)

func NewJobEnqueuer(pool *redis.Pool) JobEnqueuer {
	return &jobEnqueuer{
		enqueuer: work.NewEnqueuer(WorkerNamespace, pool),
	}
}

func (e *jobEnqueuer) Enqueue(ctx context.Context, job jobs.Job) (string, error) {
	return e.enqueue(ctx, job, 0)
}

func (e *jobEnqueuer) EnqueueAfter(ctx context.Context, job jobs.Job, duration time.Duration) (string, error) {
	return e.enqueue(ctx, job, duration)
}

func (e *jobEnqueuer) enqueue(ctx context.Context, job jobs.Job, duration time.Duration) (string, error) {
	b, err := job.Body()
	if err != nil {
		return "", err
	}

	args := make(map[string]interface{})
	contextArgs := e.contextArgs(ctx)
	args[ContextKey] = contextArgs
	args[BodyKey] = b

	var internalJob *work.Job
	if duration > 0 {
		var scheduledJob *work.ScheduledJob
		scheduledJob, err = e.enqueuer.EnqueueIn(job.Name(), int64(duration.Seconds()), args)
		if err == nil {
			internalJob = scheduledJob.Job
		}
	} else {
		internalJob, err = e.enqueuer.Enqueue(job.Name(), args)
	}

	if err != nil {
		logger.Errorf("[x] Failed to enqueue job, job_name: %v, duration: %v, body: %v, context: %v, err: %v", job.Name(), duration, string(b), contextArgs, err)
		return "", err
	}

	logger.Infof("[x] Successfully enqueued job, job_id: %v, job_name: %v, body: %v, context: %v", internalJob.ID, job.Name(), string(b), contextArgs)

	return internalJob.ID, nil
}

func (e *jobEnqueuer) contextArgs(ctx context.Context) map[string]interface{} {
	args := make(map[string]interface{})
	args[RequestIDKey] = ctxutils.RequestId(ctx)
	return args
}
