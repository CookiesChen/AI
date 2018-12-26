---
title: AI | GA遗传算法解决TSP问题
date: 2018-12-25 08:52:00
tags: 人工智能
---

## AI | GA遗传算法解决TSP问题

用遗传算法求解TSP问题（问题规模等和模拟退火求解TSP实验同），要求：

1.设计较好的交叉操作，并且引入多种局部搜索操作（可替换通常遗传算法的变异操作）

2.和之前的模拟退火算法（采用相同的局部搜索操作）进行比较

3.得出设计高效遗传算法的一些经验，并比较单点搜索和多点搜索的优缺点。

[源代码](https://github.com/CookiesChen/AI/blob/master/SA/ga)

### 数据

在TSPLIB（http://comopt.ifi.uni-heidelberg.de/software/TSPLIB95/）中选一个大于100个城市数的TSP问题，本次选取的数据集`d198.tsp`。

### 实验环境

语言：`Golang`

可视化：`js`+`html`

### 算法流程

遗传算法设计五大要素：

* 个体编码

* 初始种群的产生
* 适应度函数设计
* 遗传操作设计
* 控制参数设定

**流程图：**

![process](https://cookieschen.cn/img/AI/ga/process.png)

### 个体编码

本次实验解决tsp问题，编码直接使用路径即可。

### 初始种群的产生

遗传算法是对群体进行的进化操作，需要给其淮备一些表示起始搜索点的初始群体数据。本次使用`贪心算法`和随机法生成初始解。贪心算法生成的解占种群数的`5%`。

贪心算法的策略是，随机选取一个结点p，从剩余的结点中选取离p最近的结点，然后把该节点当做p，从剩余的结点中选取离p最近的结点。

### 适应度函数设计

本次实验解决tsp问题，适应度直接使用路径长度即可。

### 遗传操作设计

遗传操作主要涉及三个操作：

* 选择操作
* 交叉操作
* 变异操作

#### 选择操作

选择操作从当前群体中选出个体以生成交配池，所选出的个体具有良好的特征，以便生成优良的后代。

选择操作大致分为三种：

* 基于适应值比例
* 基于排名
* 基于局部竞争机制

本次使用第二种，基于排名的`线性排名选择`。线性排名选择需要先将种群成员按照适应值进行排序，然后根据一个线性函数分配选择概率pi，线性函数`pi=(a-(b*i)/(N+1))/N`，其中i=1,2,3....N，N为种群成员数量，a = 1.1，b=2(a-1)。然后计算累积概率，使用轮盘赌随机挑选个体。

`Golang`代码如下：

```go
// 选择操作
func selection() {
	// 前10%获得绝对交配权
	tenth := int(float32(populationNum)*0.1)
	matingPool = make(individuals, tenth)
	copy(matingPool, population)
	// 使用基于排名的轮盘赌选出20%
	fiftieth := tenth*7
	p := make([]float64, fiftieth)
	a := 1.1
	b := 2*(a-1)
	res := 0.0
	for i:= 0 ; i < fiftieth ; i++ {
		res += (a - b*float64(i)/(float64(fiftieth)+1))/float64(fiftieth)
		p[i] = res
	}
	for i:= 0 ; i < tenth*3; i++ {
		r := rand.Float64()
		for j:= 0 ; j < fiftieth ; j++ {
			if r <= p[j] {
				matingPool = append(matingPool, population[tenth+j+1])
				break
			}
		}
	}
	// 前10%直接遗传
	population = population[:tenth]
}
```

#### 交叉操作

交叉操作是将两个个体的遗传物质交换产生新的个体，可以把两个个体的优良格式传递到下一代的某个个体中，该个体具有优于前驱的性能。如果交叉后得到的个体性能不佳，可以在后面的复制操作中将其淘汰。

本次实验中使用的是`Order Crossover (OX)`，过程分为以下三部：

* 随机选择一对染色体（父代）中几个基因的起止位置（两染色体被选位置相同）：

![1545668654120](https://cookieschen.cn/img/AI/ga/1545668654120.png)

* 生成一个子代，并保证子代中被选中的基因的位置与父代相同：

![1545668662187](https://cookieschen.cn/img/AI/ga/1545668662187.png)

* 先找出第一步选中的基因在另一个父代中的位置，再将其余基因按顺序放入上一步生成的子代中：

![1545668673840](https://cookieschen.cn/img/AI/ga/1545668673840.png)

该操作只生成一个子代。

Golang代码如下：

```go
// 交叉操作和变异操作
// OX
func crossover() {
	poolSize := len(matingPool)
	for len(population) < 2*populationNum {
		r := rand.Float64()
		// 随机选择交配对象
		parent1 := matingPool[rand.Intn(poolSize-1)]
		parent2 := matingPool[rand.Intn(poolSize-1)]
		r1 := rand.Intn(cityNum-1)
		r2 := r1 + rand.Intn(cityNum-1-r1)
		var newIndividual individual
		if r <= Pc {
			middle := parent1.genes[r1:r2+1]
			newIndividual.genes = make([]node, 0)
			count := 0
			once := 0
			for i:=0; i < cityNum; i++ {
				flag := 0
				for _,v := range middle{
					if v.x == parent2.genes[i].x && v.y == parent2.genes[i].y {
						flag = 1
						break
					}
				}
				if count == r1 && once == 0 {
					once = 1
					newIndividual.genes = append(newIndividual.genes, middle...)
				}
				if flag == 0 {
					newIndividual.genes = append(newIndividual.genes,  parent2.genes[i])
					count++
				}
			}
		} else {
			newIndividual = parent1
		}
		r = rand.Float64()
		if r < Pm {
			newIndividual = twoOpt(newIndividual)
		}
		population = append(population, newIndividual)
	}
}
```

#### 变异操作

变异是在个体中遗传物质被改变，他可以使运算过程中丢弃的个体的某些重要特性得以恢复。本次变异的操作使用的是局部搜索和SA算法中使用的`2-opt`。

- 2-opt：在当前路径path中任意取两点，将两点间的路径逆序。

### 控制参数设定

种群大小：150

迭代次数：10000

交叉率：98%

变异率：0.75%

### 实验测试

| 实验截图                                                     |
| ------------------------------------------------------------ |
| ![1545797787011](https://cookieschen.cn/img/AI/ga/1545797787011.png) |
| ![1545797901433](https://cookieschen.cn/img/AI/ga/1545797901433.png) |

### 实验分析

这是使用贪心算法生成初始解的结果。因为使用了贪心算法，种群内一开始就有了比较好的解，误差大概在16%左右，再通过遗传算法，解的误差可以稳定的进入到10%以内，如果进化的代数足够多，生成的解的质量也越优。

`与SA相比：`

* 计算时间比SA的长。
* 误差与SA的。

### 设计经验

遗传算法理解起来并不难，该算法的潜力也非常大，但是设计起来比较难，本次实验总结了一些设计经验：

* 交叉率：交叉率可以设置的高一点，因为本次实现的交叉产生新的优秀的个体概率比较低，因此可以通过多进行几次交叉操作。
* 进行筛选：生成新生个体的数量可以是种群的几倍，这个可以类比于自然界，动物交配产生的后代肯定不止有1个，并且产生的后代也不是一定能存活下来的，本次实验中生成了种群数量两倍的新生个体，再筛选一半的新生个体，这可以理解为新生个体夭折。
* 贪心算法生成初始解：使用贪心算法生成部分的初始解，就相当于引入了优良基因，再这基础进行遗传，可以更高效的产生优质解。