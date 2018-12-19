package hill_climbing

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

type node struct {
	x float64
	y float64
}

type nodePath []node

var (
	path         nodePath
	MaxIteration = 10000
	cityNum      int
)

func Exec(xs []float64, ys []float64) {
	cityNum = len(xs)
	for i := 0; i < cityNum; i++ {
		path = append(path, node{xs[i], ys[i]})
	}
	nextStep := hillClimbing()
	js.Global().Set("nextStep", js.NewCallback(func(args []js.Value) {
		nextStep()
	}))
}

// 爬山法
// 2-opt算法生成领域
func hillClimbing() func() {
	itCount := 0
	currentDis := distance(path)
	n := 0
	return func() {
		for i := 0; i < 1000; i++ {
			n++
			newPath := twoOpt()
			newPath1 := threeChange()
			e1 := distance(newPath)
			e2 := distance(newPath1)
			if e1 > e2 {
				newPath = newPath1
				e1 = e2
			}
			if e1 < currentDis {
				itCount = 0
				currentDis = e1
				path = newPath
			} else {
				itCount++
			}
		}
		var nodeX strings.Builder
		var nodeY strings.Builder
		for nodeXY := 0; nodeXY < len(path); nodeXY++ {
			nodeX.WriteString(strconv.FormatFloat(path[nodeXY].x, 'g', 10, 64))
			nodeX.WriteByte(',')
			nodeY.WriteString(strconv.FormatFloat(path[nodeXY].y, 'g', 10, 64))
			nodeY.WriteByte(',')
		}
		js.Global().Call("updateNode", js.ValueOf(n), js.ValueOf((currentDis-15780)/15780*100),
			js.ValueOf(nodeX.String()), js.ValueOf(nodeY.String()))
		//fmt.Println(path)
	}
}

func twoOpt() (newPath nodePath) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	start := r.Intn(cityNum - 1)
	end := start + r.Intn(cityNum-1-start)
	// 第一段不变
	newPath = make(nodePath, cityNum)
	copy(newPath, path[:start])
	// 第二段反转
	for i := 0; i <= end-start; i++ {
		newPath[start+i] = path[end-i]
	}
	// 第三段不变
	for i := end + 1; i < cityNum; i++ {
		newPath[i] = path[i]
	}
	return newPath
}

func threeChange() (newPath nodePath) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	start := r.Intn(cityNum - 1)
	middle := start + r.Intn(cityNum-1-start)
	end := middle + r.Intn(cityNum-1-middle)
	// 第一段不变
	newPath = make(nodePath, cityNum)
	for i := 0; i < start; i++ {
		newPath[i] = path[i]
	}
	// 第二段第三段交换
	count := 0
	for i := 0; i <= end-middle; i++ {
		newPath[start+i] = path[middle+i]
		count = start + i
	}
	for i := 0; i < middle-start; i++ {
		newPath[count+1+i] = path[start+i]
	}
	// 第四段不变
	for i := end + 1; i < cityNum; i++ {
		newPath[i] = path[i]
	}
	return newPath
}

func distance(paths nodePath) float64 {
	dis := 0.0
	for i, v := range paths {
		last := (i + cityNum - 1) % cityNum
		disx := v.x - paths[last].x
		disy := v.y - paths[last].y
		dis += math.Sqrt(disx*disx + disy*disy)
	}
	return dis
}
