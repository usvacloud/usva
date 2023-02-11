package workers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/romeq/usva/internal/dbengine"
	"github.com/romeq/usva/internal/generated/db"
	"github.com/romeq/usva/internal/utils"
)

func TestTrasher_Run(t *testing.T) {
	upd := t.TempDir()
	dbq := dbengine.Init(utils.NewTestDatabaseConfiguration())

	type fields struct {
		TimeUntilRemove time.Duration
		Interval        time.Duration
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		runtime   time.Duration
		wantSaved bool
	}{
		{
			name: "test-1",
			fields: fields{
				TimeUntilRemove: time.Second * 2,
				Interval:        time.Millisecond * 2,
			},
			args: args{
				ctx: context.Background(),
			},
			runtime:   time.Second,
			wantSaved: true,
		},
		{
			name: "test-2",
			fields: fields{
				TimeUntilRemove: time.Millisecond * 5,
				Interval:        time.Millisecond * 2,
			},
			args: args{
				ctx: context.Background(),
			},
			runtime:   time.Second,
			wantSaved: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Trasher{
				TimeUntilRemove: tt.fields.TimeUntilRemove,
				Interval:        tt.fields.Interval,
				UploadDirectory: upd,
				db:              dbq,
				Running:         true,
			}

			fuuid := uuid.NewString()
			if err := dbq.NewFile(context.Background(), db.NewFileParams{
				FileUuid:     fuuid,
				AccessToken:  uuid.NewString(),
				EncryptionIv: []byte{},
			}); err != nil {
				t.Errorf("Trasher.Run(): %+v", err)
			}

			ctx, cancel := context.WithTimeout(tt.args.ctx, tt.runtime)
			defer cancel()
			ch := tr.Run(ctx)

			select {
			case <-ch:
				break
			case <-ctx.Done():
				break
			}

			file, err := dbq.FileInformation(context.Background(), fuuid)
			if !tt.wantSaved && err == pgx.ErrNoRows {
				return
			}

			if err != nil {
				t.Errorf("Trasher.Run(): %+v", err)
			}
			if tt.wantSaved && file.FileUuid != fuuid {
				t.Errorf("Trasher.Run() = %v, want %v", file.FileUuid, fuuid)
			}
		})
	}
}
