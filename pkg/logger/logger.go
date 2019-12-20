package logger

import zap "go.uber.org/zap"

func New() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

	return sugar
}
