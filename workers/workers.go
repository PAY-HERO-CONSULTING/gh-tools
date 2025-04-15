package workers

import (
	"context"
	"sync"

	"github.com/PAY-HERO-CONSULTING/gh-tools/db"
	"github.com/PAY-HERO-CONSULTING/gh-tools/queue_manager"
)

type Worker interface {
	Shutdown()
	Start(ctx context.Context, dB db.DB, wg *sync.WaitGroup, queueManager queue_manager.QueueManager)
}
