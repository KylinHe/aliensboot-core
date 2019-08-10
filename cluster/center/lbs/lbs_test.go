package lbs

import (
	"github.com/KylinHe/aliensboot-core/common/util"
	"testing"
)

var hashRingLbs *HashRing
var pollingLbs *PollingLBS

func init() {
	hashRingLbs = NewHashRing(400)
	pollingLbs = NewPollingLBS()
	for i:=0; i<100; i++ {
		nodeName := "node-" + util.IntToString(i)
		hashRingLbs.AddNode(nodeName, 1)
		pollingLbs.AddNode(nodeName, 1)
	}


}

func BenchmarkHashRingGetNode(b *testing.B) {
	for i := 0; i<b.N; i++ {
		hashRingLbs.GetNode("test" + util.IntToString(i))
	}
}

func BenchmarkPollingGetNode(b *testing.B) {
	//sb := strings.Builder{}
	for i := 0; i<b.N; i++ {
		pollingLbs.GetNode("test" + util.IntToString(i))
		//sb.WriteString(node + ",")
	}
	//fmt.Print(sb.String())
}