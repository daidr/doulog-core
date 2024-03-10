package filter

import (
	"testing"
)

func BenchmarkText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Text.Replace("日本AV演员兼电视、电影演员。苍井空AV女优是龟儿子xx出龟儿子道, 日本AV女优牛逼们最精彩的表演是AV演员色情表演牛逼牛逼牛逼")
	}
}
