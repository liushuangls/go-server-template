package logdbsync

import (
	"context"
	"sync"
	"time"

	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/internal/data/ent"
)

type worker struct {
	sync.Mutex
	repo  *data.ServerLogRepo
	logs  []*ent.ServerLog
	logCh chan *ent.ServerLog
	limit int
}

func newWorker(logCh chan *ent.ServerLog, repo *data.ServerLogRepo, limit int) *worker {
	w := &worker{
		repo:  repo,
		logs:  make([]*ent.ServerLog, 0, limit),
		logCh: logCh,
		limit: limit,
	}
	go w.run()
	return w
}

func (w *worker) run() {
	for log := range w.logCh {
		w.add(log)
		if len(w.logs) >= w.limit {
			w.writeDB()
		}
	}
}

func (w *worker) add(log *ent.ServerLog) {
	w.Lock()
	defer w.Unlock()

	w.logs = append(w.logs, log)
}

func (w *worker) writeDB() {
	w.Lock()
	defer w.Unlock()

	if len(w.logs) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	_ = w.repo.Creates(ctx, w.logs...)
	w.logs = make([]*ent.ServerLog, 0, w.limit)
	cancel()
}
