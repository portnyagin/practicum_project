package repository

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler/mocks"
	"github.com/stretchr/testify/assert"
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
			args:    args{"login 12", "pass"},
			wantErr: false,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	postgresHandler := mocks.NewMockDBHandler(mockCtrl)

	target, _ := NewUserRepository(postgresHandler, Log)

	checkRow := mocks.NewMockRow(mockCtrl)
	checkRow.EXPECT().Scan(gomock.Any()).SetArg(0, 10).Return(nil).AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postgresHandler.EXPECT().Execute(context.Background(), CreateUser, gomock.Any(), tt.args.login, tt.args.pass).Return(nil)
			postgresHandler.EXPECT().Execute(context.Background(), CreateAccount, gomock.Any()).Return(nil)
			postgresHandler.EXPECT().QueryRow(context.Background(), GetNextUserID).Return(checkRow, nil)
			userID, err := target.Save(context.Background(), tt.args.login, tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if userID == 0 {
				t.Errorf("Got userID=0")
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

	mockPostgresHandler := mocks.NewMockDBHandler(mockCtrl)

	mockPostgresHandler.EXPECT().QueryRow(context.Background(), CheckUser, "login 22", "pass").Return(checkRow, nil)
	mockPostgresHandler.EXPECT().QueryRow(context.Background(), CheckUser, "login 23", gomock.Any()).Return(emptyRow, nil)

	target, _ := NewUserRepository(mockPostgresHandler, Log)
	//initDatabase(context.Background(), postgresHandler)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := target.Check(context.Background(), tt.args.login, tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserRepository Check() got = %v, want %v", got, tt.want)
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

	target, _ := NewUserRepository(postgresHandler, Log)
	initDatabase(context.Background(), postgresHandler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			userID, err := target.Save(context.Background(), tt.argsSave.login, tt.argsSave.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if userID == 0 {
				t.Errorf("Got userID=0")
			}

			got, err := target.Check(context.Background(), tt.argsCheck.login, tt.argsCheck.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserRepository Check() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_SaveGetInt(t *testing.T) {
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
		{name: "UserRepository. Save+Get. Case #1",
			argsSave:  args{"user31", "pass"},
			argsCheck: args{"user31", "pass"},
			want:      true,
			wantErr:   false,
		},
		{name: "UserRepository. Save+Get. Case #2",
			argsSave:  args{"user32", "pass"},
			argsCheck: args{"user32", "badPass"},
			want:      false,
			wantErr:   false,
		},
		{name: "UserRepository. Save+Get. Case #3",
			argsSave:  args{"user34", "pass"},
			argsCheck: args{"unexistUser", "badPass"},
			want:      false,
			wantErr:   true,
		},
	}

	target, _ := NewUserRepository(postgresHandler, Log)
	initDatabase(context.Background(), postgresHandler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := target.Save(context.Background(), tt.argsSave.login, tt.argsSave.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if userID == 0 {
				t.Errorf("Got userID=0")
			}

			user, err := target.GetUserByLogin(context.Background(), tt.argsCheck.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.argsSave.login, user.Login, "UserRepository. Get() compare result. login want = %s, got %s", tt.argsSave.login, user.Login)
				assert.Equal(t, tt.argsSave.pass, user.Pass, "UserRepository. Get() compare result. pass want = %s, got %s", tt.argsSave.pass, user.Pass)
			}
		})
	}
}
