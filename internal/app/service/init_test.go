package service

import (
	"go.uber.org/zap"
	"os"
	"testing"
)

var log *zap.Logger

func TestMain(m *testing.M) {
	log, _ = zap.NewDevelopment()
	os.Exit(m.Run())
}
