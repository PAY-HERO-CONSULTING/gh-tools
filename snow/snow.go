package snow

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
	SnowEnqueuer interface {
		Enqueue(ctx context.Context, job *GenericSnow) (string, error)
		EnqueueAfter(ctx context.Context, job *GenericSnow, duration time.Duration) (string, error)
	}

	snowEnqueuer struct {
		enqueuer *work.Enqueuer
	}

	GenericSnow struct {
		Name           string
		Payload        []byte
		MaxConcurrency uint
		MaxFails       uint
		Priority       uint
	}
)

func NewSnowEnqueuer(pool *redis.Pool) SnowEnqueuer {
	return &snowEnqueuer{
		enqueuer: work.NewEnqueuer(WorkerNamespace, pool),
	}
}

func (e *snowEnqueuer) Enqueue(ctx context.Context, job *GenericSnow) (string, error) {
	return e.enqueue(ctx, job, 0)
}

func (e *snowEnqueuer) EnqueueAfter(ctx context.Context, job *GenericSnow, duration time.Duration) (string, error) {
	return e.enqueue(ctx, job, duration)
}

func (e *snowEnqueuer) enqueue(ctx context.Context, job *GenericSnow, duration time.Duration) (string, error) {
	args := e.buildJobArgs(ctx, job)

	var internalJob *work.Job
	var err error

	if duration > 0 {
		var scheduledJob *work.ScheduledJob
		scheduledJob, err = e.enqueuer.EnqueueIn(job.Name, int64(duration.Seconds()), args)
		if err == nil {
			internalJob = scheduledJob.Job
		}
	} else {
		internalJob, err = e.enqueuer.Enqueue(job.Name, args)
	}

	e.logEnqueueSuccess(internalJob, job, args)
	return internalJob.ID, nil
}

func (e *snowEnqueuer) buildJobArgs(ctx context.Context, job *GenericSnow) map[string]interface{} {
	args := map[string]interface{}{
		ContextKey: map[string]interface{}{
			RequestIDKey: ctxutils.RequestId(ctx),
		},
		BodyKey: string(job.Payload),
	}
	return args
}

func (e *snowEnqueuer) enqueueScheduledJob(jobName string, args map[string]interface{}, duration time.Duration) (*work.Job, error) {
	scheduledJob, err := e.enqueuer.EnqueueIn(jobName, int64(duration.Seconds()), args)
	if err != nil {
		return nil, err
	}
	return scheduledJob.Job, nil
}

func (e *snowEnqueuer) logEnqueueError(job *GenericSnow, duration time.Duration, args map[string]interface{}, err error) {
	logger.Errorf("[x] Failed to enqueue job, snow_name: %v, duration: %v, body: %v, context: %v, err: %v",
		job.Name, duration, string(job.Payload), args[ContextKey], err)
}

func (e *snowEnqueuer) logEnqueueSuccess(internalJob *work.Job, job *GenericSnow, args map[string]interface{}) {
	logger.Infof("[x] Successfully enqueued job, job_id: %v, snow_name: %v, body: %v, context: %v",
		internalJob.ID, job.Name, string(job.Payload), args[ContextKey])
}

type WorkSnow struct {
	*GenericSnow
}

func (wj *WorkSnow) Body() (string, error) {
	return string(wj.Payload), nil
}

func (wj *WorkSnow) Name() string {
	return wj.GenericSnow.Name
}

func (wj *WorkSnow) Options() []jobs.JobOption {
	return []jobs.JobOption{
		jobs.WithMaxConcurrency(wj.GenericSnow.MaxConcurrency),
		jobs.WithMaxFails(wj.GenericSnow.MaxFails),
		jobs.WithPriority(uint(wj.GenericSnow.Priority)),
	}
}
