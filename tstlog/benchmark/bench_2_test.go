package benchmarks

import (
	// "io/ioutil"
	// "log"

	"testing"

	"github.com/Sirupsen/logrus"
)

func Benchmark_file2_1(b *testing.B) {
	b.Run("bench_1_tst", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logrus.Infoln("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
			}
		})
	})
}
