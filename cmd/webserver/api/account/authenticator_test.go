package account

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/dbengine"
	"github.com/romeq/usva/internal/generated/db"
	"github.com/romeq/usva/internal/utils"
)

const (
	authKeepTime = time.Second / 2
)

var author UserAuthenticator

func TestMain(m *testing.M) {
	dbconfig := utils.NewTestDatabaseConfiguration()
	q, close := dbengine.Init(dbconfig)
	defer close()

	author = NewAuthenticator(q, authKeepTime)
	m.Run()
}

func TestRegister(t *testing.T) {
	tests := []struct {
		wantAuthSuccess          bool
		timeBeforeAuthentication time.Duration
	}{
		{
			wantAuthSuccess:          true,
			timeBeforeAuthentication: authKeepTime / 2,
		},
		{
			wantAuthSuccess:          false,
			timeBeforeAuthentication: authKeepTime * 2,
		},
	}

	for i, test := range tests {
		ct := context.Background()
		params := db.NewAccountParams{
			Username: fmt.Sprintf("user-%d", i),
			Password: fmt.Sprintf("password-%d", i),
		}

		session, err := author.Register(ct, params)
		if err != nil {
			t.Error(err)
		}

		u, err := author.Authenticate(ct, session)

		time.Sleep(test.timeBeforeAuthentication)

		if test.wantAuthSuccess && errors.Is(err, api.ErrAuthFailed) {
			if u.Username != params.Username {
				t.Errorf("username-expected=%s username-got=%s", params.Username, u.Username)
			}

			t.Errorf("err=%+v", err)
		}

	}
}
