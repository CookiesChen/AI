package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetTspData(path string)  {
	file, err := os.Open(path)
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
		cityId, coordinateX, coordinateY := lineData[0], lineData[1], lineData[2]
		fmt.Println(cityId + " " + coordinateX + " " + coordinateY)
		strconv.ParseFloat(coordinateX, 64)
	}
}
