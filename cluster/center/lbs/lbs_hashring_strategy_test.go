package lbs

import (
	"github.com/KylinHe/aliensboot-core/common/util"
	"testing"
)

var lbs *HashRing

func init() {
	lbs = NewHashRing(400)
	lbs.AddNode("node-1", 1)

}

func BenchmarkGetNode(b *testing.B) {
	for i := 0; i<b.N; i++ {
		lbs.GetNode("test" + util.IntToString(i))
	}
}