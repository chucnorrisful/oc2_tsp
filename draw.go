package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func draw() {
	cmd := exec.Command("python3", "graph/draw.py")
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	fmt.Println(string(out))
}

func exportD(d [][]int) {
	os.Remove("graph/matrix.csv")
	file, err := os.OpenFile("graph/matrix.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range d {
		for k, v := range data {
			if k == len(data)-1 {
				_, _ = datawriter.WriteString(strconv.Itoa(v))
			} else {
				_, _ = datawriter.WriteString(strconv.Itoa(v) + ",")
			}

		}
		_, _ = datawriter.WriteString("\n")
	}

	datawriter.Flush()
	file.Close()
}

func exportPath(p []int) {
	os.Remove("graph/path.csv")
	file, err := os.OpenFile("graph/path.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for k, data := range p {
		if k == len(p)-1 {
			_, _ = datawriter.WriteString(strconv.Itoa(data))
		} else {
			_, _ = datawriter.WriteString(strconv.Itoa(data) + ",")
		}
	}

	datawriter.Flush()
	file.Close()
}
