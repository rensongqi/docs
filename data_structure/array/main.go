package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type ValNode struct {
	raw int
	col int
	va  int
}

func main() {
	var sparAry [11][11]int
	sparAry[1][2] = 1
	sparAry[2][3] = 2
	fmt.Println("原始的数组：")
	for _, v := range sparAry {
		for _, v2 := range v {
			fmt.Printf("%d \t", v2)
		}
		fmt.Println()
	}

	var sliceAry []ValNode
	varNode1 := ValNode{
		raw: 11,
		col: 11,
		va:  0,
	}
	sliceAry = append(sliceAry, varNode1)
	for i, v := range sparAry {
		for j, v2 := range v {
			if v2 != 0 {
				varNode2 := ValNode{
					raw: i,
					col: j,
					va:  v2,
				}
				sliceAry = append(sliceAry, varNode2)
			}
		}
	}

	fmt.Println("-----------稀疏数组为----------")
	for i, v := range sliceAry {
		fmt.Printf("%d = %d %d %d\n", i, v.raw, v.col, v.va)
	}

	// 写入文件
	filePath := "chess.txt"
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	for _, v := range sliceAry {
		//fmt.Printf("%d = %d %d %d\n", i, v.raw, v.col, v.va)
		write.WriteString(fmt.Sprintf("%d %d %d\n", v.raw, v.col, v.va))
	}
	write.Flush()

	// 读取文件
	read, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("-----------从文件中的全部读取的数据----------")
	fmt.Printf("%v", string(read))

	fmt.Println("-----------从文件中的分行读取的数据----------")
	file1, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file1.Close()

	for {
		reader := bufio.NewReader(file1)
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			return
		}
		fmt.Print(str)
	}
}
