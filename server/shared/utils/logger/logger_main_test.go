package logger

import (
	"testing"
	"time"

	"space.online.shop.web.server/shared/utils/logger/internal"

	"github.com/stretchr/testify/suite"
)

type MyLoggerTestSuite struct {
	suite.Suite
	logger *internal.MyLogger
}

func (suite *MyLoggerTestSuite) TestNewMyLogger() {
	suite.logger = internal.NewMyLogger("server_log.properties", internal.WithPropertyUpdateInterval(3*time.Second))
	suite.NotNil(suite.logger, "logger should not be nil")

	// print something and check in stdout
	suite.logger.Debug("debug")
	suite.logger.Info("info")
	suite.logger.Warn("warn")
	suite.logger.Error("error")

	suite.logger.Close()
}

func TestMyLogger(t *testing.T) {
	suite.Run(t, new(MyLoggerTestSuite))
}
