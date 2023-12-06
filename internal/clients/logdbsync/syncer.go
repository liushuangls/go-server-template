package logdbsync

import (
	"encoding/json"

	"github.com/spf13/cast"
	"go.uber.org/atomic"

	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/internal/data/ent"
	"github.com/liushuangls/go-server-template/internal/data/ent/serverlog"
	"github.com/liushuangls/go-server-template/pkg/xstrings"
)

type Client struct {
	Repo *data.ServerLogRepo

	syncing  atomic.Bool
	logCh    chan *ent.ServerLog
	logLimit int
	workers  []*worker
}

func NewClient(repo *data.ServerLogRepo) *Client {
	c := &Client{
		Repo:     repo,
		logLimit: 1000,
		logCh:    make(chan *ent.ServerLog, 1000),
	}
	for i := 0; i < 5; i++ {
		c.workers = append(c.workers, newWorker(c.logCh, repo, 10))
	}
	return c
}

func (c *Client) Write(p []byte) (n int, err error) {
	// channel 写满后直接返回
	if len(c.logCh) >= c.logLimit {
		return len(p), nil
	}

	log := &ent.ServerLog{From: serverlog.FromLog}
	if err = json.Unmarshal(p, &log.Extra); err != nil {
		log.Level = "unknown"
		log.ErrMsg = xstrings.BytesToString(p)
		log.Extra = map[string]any{
			"unmarshal_err": err,
		}
	} else {
		log.Level = serverlog.Level(cast.ToString(log.Extra["level"]))
		log.ErrMsg = cast.ToString(log.Extra["msg"])
		log.UserID = cast.ToInt(log.Extra["user_id"])
		delete(log.Extra, "level")
		delete(log.Extra, "msg")
		delete(log.Extra, "user_id")
	}

	c.logCh <- log
	return len(p), nil
}

func (c *Client) Sync() error {
	if c.syncing.Load() {
		return nil
	}

	c.syncing.Store(true)
	for _, w := range c.workers {
		w.writeDB()
	}
	c.syncing.Store(false)
	return nil
}
