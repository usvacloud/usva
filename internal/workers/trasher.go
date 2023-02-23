package workers

import (
	"context"
	"log"
	"os"
	"path"
	"time"

	"github.com/usvacloud/usva/internal/generated/db"
)

type Trasher struct {
	TimeUntilRemove time.Duration
	Interval        time.Duration
	UploadDirectory string
	db              *db.Queries
	Running         bool
}

func NewTrasher(interval, ttl time.Duration, u string, db *db.Queries) *Trasher {
	return &Trasher{
		db:              db,
		Interval:        interval,
		TimeUntilRemove: ttl,
		UploadDirectory: u,
	}
}

func (t *Trasher) trash(ctx context.Context) {
	files, err := t.db.GetLastSeenAll(ctx)
	if err != nil {
		log.Println("Trasher:", err)
	}

	for _, file := range files {
		if time.Now().Before(file.LastSeen.Add(t.TimeUntilRemove)) {
			continue
		}

		if err := t.db.DeleteFile(ctx, file.FileUuid); err != nil {
			log.Println("Trasher:", err)
		}

		if err = os.Remove(path.Join(t.UploadDirectory, file.FileUuid)); err != nil {
			log.Println("Trasher:", err)
		}
	}
}

func (t *Trasher) worker(ctx context.Context, ch chan struct{}) {
	t.Running = true
	ticker := time.NewTicker(t.Interval)
	for t.Running {
		t.trash(ctx)
		<-ticker.C
	}
	ch <- struct{}{}
}

func (t *Trasher) Run(ctx context.Context) chan struct{} {
	ch := make(chan struct{}, 1)
	go t.worker(ctx, ch)
	return ch
}
