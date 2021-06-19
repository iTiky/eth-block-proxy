package service

import "github.com/rs/zerolog"

// NewTestBaseSvc creates a BaseSvc for tests with NOOP logger.
func NewTestBaseSvc() BaseSvc {
	return BaseSvc{
		logger: zerolog.Nop(),
	}
}
