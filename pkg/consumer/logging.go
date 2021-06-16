// Copyright 2021` Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package consumer

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLoggingLevel(cfg *Config) zap.AtomicLevel {
	level := cfg.LogLevel
	atom := zap.NewAtomicLevel()
	switch level {
	case "panic":
		atom.SetLevel(zapcore.PanicLevel)
		return atom
	case "fatal":
		atom.SetLevel(zapcore.FatalLevel)
		return atom
	case "error":
		atom.SetLevel(zapcore.ErrorLevel)
		return atom
	case "warn":
		atom.SetLevel(zapcore.WarnLevel)
		return atom
	case "info":
		atom.SetLevel(zapcore.InfoLevel)
		return atom
	case "debug":
		atom.SetLevel(zapcore.DebugLevel)
		return atom
	default:
		atom.SetLevel(zapcore.WarnLevel)
		return atom
	}
}

// setupLoggingZap
func setupLoggingZap(cfg *Config) *zap.Logger {
	//NewProductionConfig is a reasonable production logging configuration. Logging //is enabled at InfoLevel and above.
	//
	//It uses a JSON encoder, writes to standard error, and enables sampling. //Stacktraces are automatically included on logs of ErrorLevel and above.
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Level = getLoggingLevel(cfg)
	logger, _ := config.Build()
	return logger
}
