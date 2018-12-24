package ga

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

type node struct {
	x float64
	y float64
}

type individual struct {  // 个体
	genes []node    // 染色体 基因
	fitness float64 // 适应值
}

type individuals []individual

func (I individuals) Len() int {
	return len(I)
}

func (I individuals) Less(i, j int) bool {
	return  I[i].fitness < I[j].fitness
}

func (I individuals) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

//种群

var (
	path          individual
	cityNum       int
	MaxGen        = 300       // 遗传次数
	populationNum = 100       // 个体数
	population    individuals // 种群
	minGen        = 0		  // 最优解出现遗传代数
	minFitness    float64	  // 最优适应值
	Pc            = 0.7       // 交叉概率
	Pm            = 0.01      // 变异概率
	matingPool    individuals // 交配池
)

func Exec(xs []float64, ys []float64) {
	rand.Seed(time.Now().UnixNano())

	cityNum = len(xs)
	for i := 0; i < cityNum; i++ {
		path.genes = append(path.genes, node{xs[i], ys[i]})
	}
	ga()
}

func ga() {
	// 当前遗传次数
	nowGen := 0
	// 初始化种群
	initialize()
	for nowGen < MaxGen {
		// 个体评估
		for i:=0 ; i < populationNum; i++{
			population[i].fitness = evaluate(population[i])
		}
		// 遗传操作
		inheritance()
		nowGen++
		// 根据适应值进行排序
		sort.Sort(population)
		if population[0].fitness < minFitness {
			minFitness = population[0].fitness
			minGen = nowGen
		}
		break
	}
}

// 初始化
func initialize() {
	// 贪心生成10%初始解
	num := int(float32(populationNum)*0.1)
	for i := 0; i < num; i++{
		population = append(population, greedyIndividual())
	}

	// 随机生成初始解
	for i := 0; i < populationNum-num; i++ {
		population = append(population, getRandomIndividual())
	}
	// 根据适应值进行排序
	for i:=0 ; i < populationNum; i++{
		population[i].fitness = evaluate(population[i])
	}
	sort.Sort(population)
	minFitness = population[0].fitness
}

// 贪心初始解
func greedyIndividual() (newIndividual individual) {
	leftNode := make([]node, cityNum)
	copy(leftNode, path.genes)
	randNum := rand.Intn(len(leftNode)-1)
	newIndividual.genes = append(newIndividual.genes, leftNode[randNum])
	leftNode = append(leftNode[:randNum], leftNode[randNum+1:]...)
	nowNode := leftNode[randNum]
	for len(leftNode) > 1 {
		index := 0
		minDis := -1.0
		for i,v := range leftNode{
			dis := (v.x-nowNode.x)*(v.x-nowNode.x)+(v.y-nowNode.y)*(v.y-nowNode.y)
			if minDis == -1 || minDis > dis {
				index = i
				minDis = dis
			}
		}
		nowNode = leftNode[index]
		newIndividual.genes = append(newIndividual.genes, leftNode[index])
		leftNode = append(leftNode[:index], leftNode[index+1:]...)
	}
	newIndividual.genes = append(newIndividual.genes, leftNode[0])
	return newIndividual
}

// 随机初始解
func getRandomIndividual()(newIndividual individual){
	newIndividual.genes = make([]node, cityNum)
	randomArray := make([]int, cityNum)
	for i:=0 ; i < cityNum ; i++ {
		randomArray[i] = i
	}
	for i:= 0 ; i < cityNum ; i++ {
		randomNum1 := rand.Intn(cityNum-1)
		randomNum2 := rand.Intn(cityNum-1)
		randomArray[randomNum1], randomArray[randomNum2] = randomArray[randomNum2], randomArray[randomNum1]
	}
	for i:= 0 ; i < cityNum ; i++ {
		newIndividual.genes[i] = path.genes[randomArray[i]]
	}
	return newIndividual
}

// 适应值函数
func evaluate(paths individual) float64 {
	dis := 0.0
	for i, v := range paths.genes {
		last := (i + cityNum - 1) % cityNum
		disx := v.x - paths.genes[last].x
		disy := v.y - paths.genes[last].y
		dis += math.Sqrt(disx*disx + disy*disy)
	}
	return dis
}

// 遗传操作
func inheritance() {
	oldPopulation := make(individuals, populationNum)
	copy(oldPopulation, population)
	// 选择
	selection()
	// 交叉
	crossover()
}

// 选择操作
func selection() {
	// 前10%获得绝对交配权
	tenth := int(float32(populationNum)*0.1)
	matingPool = make(individuals, tenth)
	copy(matingPool, population)
	// 前10-60%中使用基于排名的轮盘赌选出20%
	fiftieth := tenth*5
	p := make([]float64, fiftieth)
	a := 1.1
	b := 2*(a-1)
	res := 0.0
	for i:= 0 ; i < fiftieth ; i++ {
		res += (a - b*float64(i)/(float64(fiftieth)+1))/float64(fiftieth)
		p[i] = res
	}
	for i:= 0 ; i < 2*tenth ; i++ {
		r := rand.Float64()
		for j:= 0 ; j < fiftieth ; j++ {
			if r <= p[j] {
				matingPool = append(matingPool, population[tenth+j+1])
				break
			}
		}
	}
	// 前10%直接遗传
	population = population[:tenth+1]
}

// 交叉操作
// OX
func crossover() {
	poolSize := len(matingPool)
	for len(population) < populationNum {
		r := rand.Float64()
		// 随机选择交配对象
		parent1 := matingPool[rand.Intn(poolSize-1)]
		parent2 := matingPool[rand.Intn(poolSize-1)]
		if r <= Pc {
			r1 := rand.Intn(cityNum-1)
			r2 := r1 + rand.Intn(cityNum-1-r1)
			middle :=
		}
	}
}