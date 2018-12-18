package sa

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
	cityNum      int
	tEnd         = 0.1    // 终止温度
	q			 = 0.99   // 降温系数
)

func Exec(xs []float64, ys []float64)  {
	cityNum = len(xs)
	for i := 0; i < cityNum; i++ {
		path = append(path, node{xs[i], ys[i]})
	}
	sa()
}

func sa()  {
	n := 0
	T := 100.0
	L := 5000
	currentDis := distance(path)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for T > tEnd {
		for i:=0 ; i < L ; i++ {
			newPath := getNewPath()
			dis := distance(newPath)
			df := dis - currentDis
			// Metropolis准则
			if df < 0 {
				// 接受新解
				currentDis = dis
				path = newPath
			} else {
				R := float64(r.Intn(100))
				R /= 100
				P := 1/(1+math.Exp(-df/T))
				// 接受恶化解
				if R > P {
					currentDis = dis
					path = newPath
				}
			}
			n++
		}
		T *= q
		fmt.Printf("迭代数: %v, %.2f%%\n", n, (currentDis-15780)/15780*100)
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