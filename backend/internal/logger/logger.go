package logger

import (
	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
