package v1

import "time"

func (s *ServiceTestSuite) TestConfigValidation() {
	cfg := Config{
		RequestTimeoutDur: 500 * time.Millisecond,
		MinRetryDur:       10 * time.Millisecond,
		MaxRetryDur:       250 * time.Millisecond,
	}

	s.Run("Ok", func() {
		s.Assert().NoError(cfg.Validate())
	})

	s.Run("Fail: invalid RequestTimeoutDur", func() {
		c := cfg
		c.RequestTimeoutDur = -1 * time.Second
		s.Assert().Error(c.Validate())
	})

	s.Run("Fail: invalid MinRetryDur", func() {
		c := cfg
		c.MinRetryDur = -1 * time.Second
		s.Assert().Error(c.Validate())
	})

	s.Run("Fail: invalid MaxRetryDur", func() {
		c := cfg
		c.MaxRetryDur = -1 * time.Second
		s.Assert().Error(c.Validate())
	})

	s.Run("Fail: MaxRetryDur < MinRetryDur", func() {
		c := cfg
		c.MinRetryDur = 300 * time.Millisecond
		s.Assert().Error(c.Validate())
	})

	s.Run("Fail: MaxRetryDur > RequestTimeoutDur", func() {
		c := cfg
		c.MaxRetryDur = c.RequestTimeoutDur + 1*time.Millisecond
		s.Assert().Error(c.Validate())
	})
}
