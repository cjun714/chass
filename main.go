package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cjun714/glog/log"
)

func main() {
	if len(os.Args) != 3 {
		log.F("not enough arguments, usage: 'ass1080 src.ass target.ass'")
	}

	srcPath := os.Args[1]
	targetPath := os.Args[2]

	byts, e := os.ReadFile(srcPath)
	if e != nil {
		log.F("read .ass failed:", e)
	}

	out, e := os.Create(targetPath)
	if e != nil {
		log.F("create .ass failed:", e)
	}
	defer out.Close()

	rd := bytes.NewReader(byts)
	sc := bufio.NewScanner(rd)

	var x, y int

	// read PlayResX and PlayResY
	for sc.Scan() {
		line := sc.Text()
		fmt.Fprintln(out, line)

		if strings.Index(line, "PlayResX:") == 0 {
			fmt.Sscanf(line, "PlayResX: %d", &x)
		}

		if strings.Index(line, "PlayResY:") == 0 {
			fmt.Sscanf(line, "PlayResY: %d", &y)
		}

		if x != 0 && y != 0 {
			break
		}
	}

	xRatio := float32(1920) / float32(x)
	yRatio := float32(1080) / float32(y)

	// replace all \pos(xxx,yyy) lines
	for sc.Scan() {
		line := sc.Text()

		if strings.Index(line, "Dialogue:") != 0 {
			fmt.Fprintln(out, line)
			continue
		}

		newline := processLine(line, xRatio, yRatio)
		fmt.Fprintln(out, newline)
	}
}

func processLine(line string, xRatio, yRatio float32) string {
	if strings.Contains(line, "\\pos(") {
		return processPos(line, xRatio, yRatio)
	}

	if strings.Contains(line, "\\move(") {
		return processMove(line, xRatio, yRatio)
	}

	return line
}

func processPos(line string, xRatio, yRatio float32) string {
	idx := strings.Index(line, "\\pos(")

	subStr := line[idx:]

	idx = strings.IndexByte(subStr, ')')
	if idx == -1 { // if no "\pos()"
		return line
	}

	str := subStr[0 : idx+1] // should be "\pos(xxx,yyy)"

	var x, y float32
	fmt.Sscanf(str, "\\pos(%f,%f)", &x, &y)

	x, y = x*xRatio, y*yRatio

	newStr := "\\pos(" + strconv.Itoa(int(x)) + "," + strconv.Itoa(int(y)) + ")"

	return strings.Replace(line, str, newStr, 1)
}

func processMove(line string, xRatio, yRatio float32) string {
	idx := strings.Index(line, "\\move(")

	subStr := line[idx:]

	idx = strings.IndexByte(subStr, ')')
	if idx == -1 { // if no "\move()"
		return line
	}

	str := subStr[0 : idx+1] // should be "\move(x0,y0,x1,y1)"

	var x0, y0, x1, y1 float32
	fmt.Sscanf(str, "\\move(%f,%f,%f,%f)", &x0, &y0, &x1, &y1)

	x0, y0 = x0*xRatio, y0*yRatio
	x1, y1 = x1*xRatio, y1*yRatio

	newStr := "\\move(" + strconv.Itoa(int(x0)) + "," + strconv.Itoa(int(y0)) + "," + strconv.Itoa(int(x1)) + "," + strconv.Itoa(int(y1)) + ")"

	return strings.Replace(line, str, newStr, 1)
}
