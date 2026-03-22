package worker

import (
	"context"
	"log"
	"sync"

	"github.com/gewall/short-url/internal/domain"

	"github.com/gewall/short-url/pkg"
)

type ClickWorker struct {
	Jobs     chan domain.Clicks
	wg       sync.WaitGroup
	proccess ProcessFunc
}

type ProcessFunc func(context.Context, domain.Clicks, any) error

func NewClickWorker(workers, size int, fn ProcessFunc, repo any) *ClickWorker {
	cw := &ClickWorker{
		Jobs:     make(chan domain.Clicks, size),
		proccess: fn,
	}

	for range workers {
		cw.wg.Add(1)
		go cw.run(context.Background(), repo)
	}

	return cw
}

func (cw *ClickWorker) Submit(clicks domain.Clicks) error {
	select {
	case cw.Jobs <- clicks:
		return nil
	default:
		return pkg.ErrWorkerQueueFull
	}
}

func (cw *ClickWorker) Shutdown() {
	close(cw.Jobs)
	cw.wg.Wait()
}

func (cw *ClickWorker) run(ctx context.Context, repo any) {
	defer cw.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case clicks, ok := <-cw.Jobs:
			if !ok {
				return
			}
			if err := cw.proccess(ctx, clicks, repo); err != nil {
				log.Printf("click worker error: %v", err)

			}
			log.Printf("click worker processed: %v", clicks)
		}
	}
}
