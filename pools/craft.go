package pools

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/ctxutils"
	"github.com/PAY-HERO-CONSULTING/gh-tools/job_enqueuer"
	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
	"github.com/PAY-HERO-CONSULTING/gh-tools/pools/jobs"
	gowork "github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type WorkerPool interface {
	AddJobHandlers(handlers ...jobs.JobHandler)
	Start(ctx context.Context)
	Stop()
	AddRecurringJobs(cronTaskMap map[string]string)
}

type workerPool struct {
	ctx  context.Context
	pool *gowork.WorkerPool
}

type workerPoolContext struct{}

func log(job *gowork.Job, next gowork.NextMiddlewareFunc) error {
	logger.Infof("Starting job, job_id: %v, job_name: %v, args: %v", job.ID, job.Name, job.Args)
	return next()
}

func panicRecovery(job *gowork.Job, next gowork.NextMiddlewareFunc) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Failed to recover from panic: %v during job_id: %v, job_name: %v, args: %v", err, job.ID, job.Name, job.Args)
			debug.PrintStack()
		}
	}()

	return next()
}

func NewWorkerPool(redisPool *redis.Pool, concurrency uint, cronTaskMap map[string]string) WorkerPool {
	pool := gowork.NewWorkerPool(workerPoolContext{}, concurrency, job_enqueuer.WorkerNamespace, redisPool)
	pool.Middleware(log)
	pool.Middleware(panicRecovery) // Must go last

	for interverval, task := range cronTaskMap {
		pool.PeriodicallyEnqueue(interverval, task)
	}

	workerPool := &workerPool{
		pool: pool,
	}

	return workerPool
}

func (wp *workerPool) AddJobHandlers(jobHandlers ...jobs.JobHandler) {
	for _, jobHandler := range jobHandlers {
		// Set up default options
		jobOptions := gowork.JobOptions{
			Priority: 1,
			MaxFails: 1,
		}

		job := jobHandler.Job()
		opts := job.Options()
		for _, opt := range opts {
			opt(jobOptions)
		}

		wrappedJobHandler := wp.wrapJobHandler(jobHandler)

		wp.pool.JobWithOptions(job.Name(), jobOptions, wrappedJobHandler)
	}
}

func (wp *workerPool) AddRecurringJobs(cronTaskMap map[string]string) {
	for interval, task := range cronTaskMap {
		wp.pool.PeriodicallyEnqueue(interval, task)
	}
}

func (wp *workerPool) Start(ctx context.Context) {
	wp.ctx = ctx
	wp.pool.Start()
}

func (wp *workerPool) Stop() {
	wp.pool.Stop()
}

func (wp *workerPool) wrapJobHandler(jobHandler jobs.JobHandler) func(job *gowork.Job) error {
	return func(job *gowork.Job) error {
		startTime := time.Now()

		var requestID string
		rawContext, ok := job.Args[job_enqueuer.ContextKey]
		if ok {
			context, ok := rawContext.(map[string]interface{})
			if ok {
				rawRequestID := context[job_enqueuer.RequestIDKey]
				requestID, ok = rawRequestID.(string)
				if !ok {
					requestID = ""
				}
			}
		}

		jobContext := ctxutils.WithRequestId(wp.ctx, requestID)

		rawBody := job.Args[job_enqueuer.BodyKey]
		if rawBody == nil {
			return nil
		}

		logger.Infof("Body for request: [%+v]", rawBody)

		body, ok := rawBody.(string)
		if !ok {
			duration := time.Since(startTime)
			err := fmt.Errorf("failed to cast to string: %v", rawBody)
			logger.Errorf("[JobHandler] Job failed: %v, job_id: %v, job_name: %v, duration: %v, args: %v", err, job.ID, job.Name, duration, job.Args)

			// Since this body cannot be handled properly, return nil to prevent retries
			return nil
		}

		if len(body) < 1 {
			return nil
		}

		err := jobHandler.PerformJob(jobContext, body)
		duration := time.Since(startTime)
		if err != nil {
			logger.Infof("[JobHandler] Job completed, job_id: %v, job_name: %v, duration: %v, args: %v", job.ID, job.Name, duration, job.Args)
		}

		return err
	}
}
