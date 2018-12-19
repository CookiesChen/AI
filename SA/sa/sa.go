package sa

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
	path    nodePath
	cityNum int
	tEnd    = 0.000001 // 终止温度
	q       = 0.99     // 降温系数
)

func Exec(xs []float64, ys []float64) {
	cityNum = len(xs)
	for i := 0; i < cityNum; i++ {
		path = append(path, node{xs[i], ys[i]})
	}
	nextStep := sa()
	js.Global().Set("nextStep", js.NewCallback(func(args []js.Value) {
		nextStep()
	}))
}

func sa() func() {
	n := 0
	T := 1000.0
	L := 1000
	global := js.Global()
	currentDis := distance(path)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return func() {
		if T > tEnd {
			if T > 100 {
				q = 0.9
			} else {
				q = 0.99
			}
			for i := 0; i < L; i++ {
				newPath := twoOpt()
				newPath1 := threeChange()
				dis := distance(newPath)
				dis1 := distance(newPath1)
				if dis > dis1 {
					dis = dis1
					newPath = newPath1
				}
				df := dis - currentDis
				// Metropolis准则
				if df < 0 {
					// 接受新解
					currentDis = dis
					path = newPath
				} else {
					R := r.Float64()
					P := math.Exp(-df / T)
					// 接受恶化解
					if R < P {
						currentDis = dis
						path = newPath
					}
				}
				n++
			}
			T *= q

			var nodeX strings.Builder
			var nodeY strings.Builder
			for nodeXY := 0; nodeXY < len(path); nodeXY++ {
				nodeX.WriteString(strconv.FormatFloat(path[nodeXY].x, 'g', 10, 64))
				nodeX.WriteByte(',')
				nodeY.WriteString(strconv.FormatFloat(path[nodeXY].y, 'g', 10, 64))
				nodeY.WriteByte(',')
			}
			global.Call("updateNode", js.ValueOf(n), js.ValueOf((currentDis-15780)/15780*100),
				js.ValueOf(nodeX.String()), js.ValueOf(nodeY.String()))
			// fmt.Printf("迭代数: %v, %.2f%%\n", n, (currentDis-15780)/15780*100)
		}
	}

}

func twoOpt() (newPath nodePath) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// println(cityNum - 1)
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
