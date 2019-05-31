package util

import (
	"math"
	"math/rand"
)

func Round(value float64) int32 {
	return int32(value + 0.5)
}

func Round64(value float64) int64 {
	return int64(value + 0.5)
}

func RoundFloat64(value float64) float64 {
	return float64(int64(value + 0.5))
}

func InvSqrt(x float32) float32 {
	var xhalf float32 = 0.5 * x // get bits for floating VALUE
	i := math.Float32bits(x)    // gives initial guess y0
	i = 0x5f375a86 - (i >> 1)   // convert bits BACK to float
	x = math.Float32frombits(i) // Newton step, repeating increases accuracy
	x = x * (1.5 - xhalf*x*x)
	x = x * (1.5 - xhalf*x*x)
	x = x * (1.5 - xhalf*x*x)
	return 1 / x
}

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func Abs32(i int32) int32 {
	if i < 0 {
		return -i
	}
	return i
}

func IMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func I32Max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
func I64Max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func IMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func I32Min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
func I64Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

//取中间值
func IClamp(value, min, max int) int {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}

//取中间值
func I32Clamp(value, min, max int32) int32 {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}

//取中间值
func I64Clamp(value, min, max int64) int64 {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}

//随机范围[from, to]
func RandRange(from, to int) int {
	return rand.Intn(to-from+1) + from
}
func RandRange32(from, to int32) int32 {
	var i = rand.Intn(int(to)-int(from)+1) + int(from)
	return int32(i)
}

//随机乱序,洗牌
func Shuffle(list []interface{}) {
	var c = len(list)
	if c < 2 {
		return
	}

	for i := 0; i < c-1; i++ {
		var j = rand.Intn(c-i) + i //这里rand需要包含i自身，否则不均匀
		if i != j {
			list[i], list[j] = list[j], list[i]
		}
	}
}

//随机抽取n张
func ShuffleN(list []interface{}, randCount int) []interface{} {
	var c = len(list)
	if c < 2 {
		return list
	}

	var ct = IMin(c-1, randCount)
	for i := 0; i < ct; i++ {
		var j = rand.Intn(c-i) + i //这里rand需要包含i自身，否则不均匀
		if i != j {
			list[i], list[j] = list[j], list[i]
		}
	}

	return list[:ct]
}

func ShuffleI32(list []int32) {
	var c = len(list)
	if c < 2 {
		return
	}

	for i := 0; i < c-1; i++ {
		var j = rand.Intn(c-i) + i //这里rand需要包含i自身，否则不均匀
		if i != j {
			list[i], list[j] = list[j], list[i]
		}
	}
}

func ShuffleNI32(list []int32, randCount int) []int32 {
	var c = len(list)
	if c < 2 {
		return list
	}

	var ct = IMin(c-1, randCount)
	for i := 0; i < ct; i++ {
		var j = rand.Intn(c-i) + i //这里rand需要包含i自身，否则不均匀
		if i != j {
			list[i], list[j] = list[j], list[i]
		}
	}
	return list[:ct]
}

func ShuffleI(list []int) {
	var c = len(list)
	if c < 2 {
		return
	}

	for i := 0; i < c-1; i++ {
		var j = rand.Intn(c-i) + i //这里rand需要包含i自身，否则不均匀
		if i != j {
			list[i], list[j] = list[j], list[i]
		}
	}
}

//随机乱序,洗牌，反向，效果一样
func ShuffleR(list []interface{}) {
	var c = len(list)
	if c < 2 {
		return
	}

	for i := c - 1; i >= 1; i-- {
		var j = rand.Int() % (i + 1)
		list[i], list[j] = list[j], list[i]
	}
}

//数组求和
func SumI32(list []int32) int32 {
	var sum int32
	for _, v := range list {
		sum += v
	}
	return sum
}

//矩阵列求和
func SumMatrixColI32(mat [][]int32, col int) int32 {
	var list []int32
	for index := 0; index < len(mat); index++ {
		list = append(list, mat[index][col])
	}
	return SumI32(list)
}

//固定种子伪随机
func StaticRand(seedrare, min, max int) int {
	var seed = float64(seedrare)
	seed = seed*2045 + 1
	seed = float64(int(seed) % 1048576)
	var dis = float64(max - min)
	var ret = int(min) + int(math.Floor(seed)*dis/1048576)
	return ret
}

//随机权重
func RandomMultiWeight(weightMapping map[int32]int32, count int) []int32 {
	results := []int32{}
	for i := 0; i < count; i++ {
		result := RandomWeight(weightMapping)
		if result == 0 {
			return results
		}
		results = append(results, result)
		delete(weightMapping, result)
	}
	return results
}

//随机多个权重
func RandomWeight(weightMapping map[int32]int32) int32 {
	var totalWeight int32 = 0
	for _, weight := range weightMapping {
		totalWeight += weight
	}
	if totalWeight <= 0 {
		return 0
	}
	randomValue := rand.Int31n(totalWeight) + 1
	var currentValue int32 = 0
	for id, weight := range weightMapping {
		currentValue += weight
		if currentValue >= randomValue {
			return id
		}
	}
	return 0
}
