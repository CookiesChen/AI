package hill_climbing

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type node struct {
	x float64
	y float64
}

type nodePath []node

var(
	path         nodePath
	MaxIteration = 1500
	cityNum      int
)
func Exec(xs []float64, ys []float64) {
	cityNum = len(xs)
	for i := 0; i < cityNum; i++ {
		path = append(path, node{xs[i], ys[i]})
	}
	hillClimbing()
}

// 爬山法
// 2-opt算法生成领域
func hillClimbing() {
	itCount := 0
	currentDis := distance(path)
	n := 0
	for itCount < MaxIteration {
		n++
		newPath := getNewPath()
		e := distance(newPath)
		if e < currentDis {
			itCount = 0
			currentDis = e
			path = newPath
		} else {
			itCount++
		}
		fmt.Print(n)
		fmt.Print(" ")
		fmt.Println(currentDis)
		fmt.Println(path)
	}
}


func getNewPath() (newPath nodePath) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	start := r.Intn(cityNum - 1)
	end := start + r.Intn(cityNum - 1 - start)
	// 第一段不变
	newPath = make(nodePath, cityNum)
	copy(newPath, path[:start])
	// 第二段反转
	for i:= 0; i <= end - start; i++ {
		newPath[start+i]= path[end-i]
	}
	// 第三段不变
	for i:= end + 1; i < cityNum; i++ {
		newPath[i] = path[i]
	}
	return newPath
}

func distance(paths nodePath) float64{
	dis := 0.0
	for i, v := range paths {
		last := (i+cityNum-1)%cityNum
		disx := math.Abs(v.x - paths[last].x)
		disy := math.Abs(v.y - paths[last].y)
		dis += math.Sqrt(disx*disx + disy*disy)
	}
	return dis
}



