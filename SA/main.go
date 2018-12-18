package main

import (
	"bufio"
	"github.com/CookiesChen/AI/SA/hill_climbing"
	"os"
	"strconv"
	"strings"
)

func main() {
	filePath := "tsp/d198.tsp"
	xs, ys := getData(filePath)
	hill_climbing.Exec(xs, ys)
}

// 读取文件
func getData(filePath string)(xs []float64, ys []float64){
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	for {
		line, _, _ := rd.ReadLine()
		if strings.Contains(string(line), "EOF") {
			break
		}
		data := string(line)
		lineData := strings.Split(data, " ")
		coordinateX, _ := strconv.ParseFloat(lineData[1], 64)
		coordinateY, _ := strconv.ParseFloat(lineData[2], 64)
		xs = append(xs, coordinateX)
		ys = append(ys, coordinateY)
	}
	return xs, ys
}
