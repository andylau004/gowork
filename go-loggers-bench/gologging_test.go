package bench

import (
	"testing"
	log "github.com/op/go-logging"
)
/*
ok  	_/home/fy/WorkDir/mygitwork/gocode/go-loggers-bench	54.640s
go test -cpu=1,2,4 -benchmem -benchtime=5s -bench "Gologging.*Text"
goos: linux
goarch: amd64
BenchmarkGologgingTextNegative     	30000000	       205 ns/op	     144 B/op	       2 allocs/op
BenchmarkGologgingTextNegative-2   	50000000	       131 ns/op	     144 B/op	       2 allocs/op
BenchmarkGologgingTextNegative-4   	100000000	        86.2 ns/op	     144 B/op	       2 allocs/op
BenchmarkGologgingTextPositive     	 3000000	      1973 ns/op	     920 B/op	      15 allocs/op
BenchmarkGologgingTextPositive-2   	 5000000	      1740 ns/op	     920 B/op	      15 allocs/op
BenchmarkGologgingTextPositive-4   	 5000000	      1183 ns/op	     920 B/op	      15 allocs/op
PASS

*/
func BenchmarkGologgingTextNegative(b *testing.B) {
	stream := &blackholeStream{}
	logger := log.MustGetLogger("")
	subBackend := log.NewLogBackend(stream, "", 0)
	formatter := log.MustStringFormatter("%{time:2006-01-02T15:04:05Z07:00} %{level} %{message}")
	backend := log.NewBackendFormatter(subBackend, formatter)
	leveled := log.AddModuleLevel(backend)
	leveled.SetLevel(log.ERROR, "")
	logger.SetBackend(leveled)
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

func BenchmarkGologgingTextPositive(b *testing.B) {
	stream := &blackholeStream{}
	logger := log.MustGetLogger("")
	subBackend := log.NewLogBackend(stream, "", 0)
	formatter := log.MustStringFormatter("%{time:2006-01-02T15:04:05Z07:00} %{level} %{message}")
	backend := log.NewBackendFormatter(subBackend, formatter)
	leveled := log.AddModuleLevel(backend)
	logger.SetBackend(leveled)
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
