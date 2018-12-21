package bench

import (
	"testing"

	log "gopkg.in/inconshreveable/log15.v2"
)


/*
// alloc/op 单次操作分配内存次数太多，导致性能不好。排除掉！

ok  	_/home/fy/WorkDir/mygitwork/gocode/go-loggers-bench	71.981s
go test -cpu=1,2,4 -benchmem -benchtime=5s -bench "Log15.*Text"
goos: linux
goarch: amd64
BenchmarkLog15TextNegative     	 5000000	      1073 ns/op	     368 B/op	       3 allocs/op
BenchmarkLog15TextNegative-2   	10000000	       642 ns/op	     368 B/op	       3 allocs/op
BenchmarkLog15TextNegative-4   	20000000	       490 ns/op	     368 B/op	       3 allocs/op
BenchmarkLog15TextPositive     	 2000000	      3073 ns/op	     856 B/op	      14 allocs/op
BenchmarkLog15TextPositive-2   	 3000000	      2641 ns/op	     856 B/op	      14 allocs/op
BenchmarkLog15TextPositive-4   	 2000000	      3733 ns/op	     856 B/op	      14 allocs/op
PASS
*/

func BenchmarkLog15TextNegative(b *testing.B) {
	stream := &blackholeStream{}
	logger := log.New()
	logger.SetHandler(log.LvlFilterHandler(
		log.LvlError,
		log.StreamHandler(stream, log.LogfmtFormat())),
	)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("The quick brown fox jumps over the lazy dog")
		}
	})

	if stream.WriteCount() != uint64(0) {
		b.Fatalf("Log write count")
	}
}

func BenchmarkLog15TextPositive(b *testing.B) {
	stream := &blackholeStream{}
	logger := log.New()
	logger.SetHandler(log.StreamHandler(stream, log.LogfmtFormat()))
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("The quick brown fox jumps over the lazy dog")
		}
	})

	if stream.WriteCount() != uint64(b.N) {
		b.Fatalf("Log write count")
	}
}

/*
func BenchmarkLog15JSONNegative(b *testing.B) {
	stream := &blackholeStream{}
	logger := log.New()
	logger.SetHandler(log.LvlFilterHandler(
		log.LvlError,
		log.StreamHandler(stream, log.JsonFormat())),
	)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("The quick brown fox jumps over the lazy dog", "rate", 15, "low", 16, "high", 123.2)
		}
	})

	if stream.WriteCount() != uint64(0) {
		b.Fatalf("Log write count")
	}
}

func BenchmarkLog15JSONPositive(b *testing.B) {
	stream := &blackholeStream{}
	logger := log.New()
	logger.SetHandler(log.StreamHandler(stream, log.JsonFormat()))
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("The quick brown fox jumps over the lazy dog", "rate", 15, "low", 16, "high", 123.2)
		}
	})

	if stream.WriteCount() != uint64(b.N) {
		b.Fatalf("Log write count")
	}
}
*/