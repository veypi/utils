package log

import (
	"github.com/rs/zerolog"
	"testing"
)

func withBenchZeroLog(b *testing.B, f func(*zerolog.Logger)) {
	//fileHook.Filename = "./log.test"
	//logger := zerolog.New(&fileHook).With().Timestamp().Logger()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f(logger)
			//f(&logger)
		}
	})
}

func BenchmarkZeroNoContext(b *testing.B) {
	withBenchZeroLog(b, func(log *zerolog.Logger) {
		log.Info().Msg("no context.")
	})
}
