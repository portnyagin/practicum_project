package handler

import (
	"go.uber.org/zap"
	"os"
	"testing"
)

var log *zap.Logger

func TestMain(m *testing.M) {
	//var err error
	log, _ = zap.NewDevelopment()
	//authHandler  = NewAuthHandler(, )
	os.Exit(m.Run())
}
