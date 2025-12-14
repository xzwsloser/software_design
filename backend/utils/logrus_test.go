package utils

import "testing"

func TestLogger(t *testing.T) {
	InitLogger()

	logger.Tracef("Hello World Trace")
	logger.Debugf("Hello World Debug")
	logger.Infof("Hello World Info")
	logger.Warnf("Hello World Warn")
	// logger.Errorf("Hello World Error")
	// logger.Fatalf("Hello World Fatal")
}