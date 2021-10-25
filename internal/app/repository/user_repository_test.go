package repository

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/portnyagin/practicum_project/internal/app/database/query"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler/mocks"
	"testing"
)

func TestUserRepository_Save(t *testing.T) {
	type args struct {
		login string
		pass  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "UserRepository. Save. Case #1",
			args:    args{"", "pass"},
			wantErr: true,
		},
		{name: "UserRepository. Save. Case #2",
			args:    args{"login11", ""},
			wantErr: true,
		},
		{name: "UserRepository. Save. Case #3",
			args:    args{"login 12", "pass"},
			wantErr: false,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	postgresHandler := mocks.NewMockDBHandler(mockCtrl)

	target := NewUserRepository(postgresHandler)
	//initDatabase(context.Background(), postgresHandler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postgresHandler.EXPECT().Execute(context.Background(), query.CreateUser, tt.args.login, tt.args.pass).Return(nil)
			if err := target.Save(context.Background(), tt.args.login, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_Check(t *testing.T) {
	type args struct {
		login string
		pass  string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "UserRepository. Check. Case #1",
			args:    args{"", "pass"},
			want:    false,
			wantErr: true,
		},
		{name: "UserRepository. Check. Case #2",
			args:    args{"login 21", ""},
			want:    false,
			wantErr: true,
		},
		{name: "UserRepository. Check. Case #3",
			args:    args{"login 22", "pass"},
			want:    true,
			wantErr: false,
		},
		{name: "UserRepository. Check. Case #4",
			args:    args{"login 23", "badPass"},
			want:    false,
			wantErr: false,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	checkRow := mocks.NewMockRow(mockCtrl)
	checkRow.EXPECT().Scan(gomock.Any()).SetArg(0, 1).Return(nil)

	emptyRow := mocks.NewMockRow(mockCtrl)
	emptyRow.EXPECT().Scan(gomock.Any()).Return(errors.New("no rows in result set"))

	postgresHandler := mocks.NewMockDBHandler(mockCtrl)

	postgresHandler.EXPECT().QueryRow(context.Background(), query.CheckUser, "login 22", "pass").Return(checkRow, nil)
	postgresHandler.EXPECT().QueryRow(context.Background(), query.CheckUser, "login 23", gomock.Any()).Return(emptyRow, nil)

	target := NewUserRepository(postgresHandler)
	//initDatabase(context.Background(), postgresHandler)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := target.Check(context.Background(), tt.args.login, tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Check() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_SaveCheckInt(t *testing.T) {
	type args struct {
		login string
		pass  string
	}
	tests := []struct {
		name      string
		argsSave  args
		argsCheck args
		want      bool
		wantErr   bool
	}{
		{name: "UserRepository. Save+Check. Case #1",
			argsSave:  args{"user31", "pass"},
			argsCheck: args{"user31", "pass"},
			want:      true,
			wantErr:   false,
		},
		{name: "UserRepository. Save+Check. Case #2",
			argsSave:  args{"user32", "pass"},
			argsCheck: args{"user32", "badPass"},
			want:      false,
			wantErr:   false,
		},
		{name: "UserRepository. Save+Check. Case #3",
			argsSave:  args{"user34", "pass"},
			argsCheck: args{"unexistUser", "badPass"},
			want:      false,
			wantErr:   false,
		},
	}
	postgresHandler, err := infrastructure.NewPostgresqlHandler(context.Background(), Datasource)

	if err != nil {
		panic(err)
	}
	target := NewUserRepository(postgresHandler)
	initDatabase(context.Background(), postgresHandler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := target.Save(context.Background(), tt.argsSave.login, tt.argsSave.pass); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := target.Check(context.Background(), tt.argsCheck.login, tt.argsCheck.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Check() got = %v, want %v", got, tt.want)
			}
		})
	}
}
