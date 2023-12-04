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
	idx := strings.Index(line, "\\pos(")

	if idx == -1 { // no "\pos()"
		return line
	}

	subStr := line[idx:]

	idx = strings.IndexByte(subStr, ')')
	if idx == -1 { // if no "\pos()"
		return line
	}

	posStr := subStr[0 : idx+1] // should be "\pos(xxx,yyy)"

	var x, y float32
	fmt.Sscanf(posStr, "\\pos(%f,%f)", &x, &y)

	x = x * xRatio
	y = y * yRatio

	newPosStr := "\\pos(" + strconv.Itoa(int(x)) + "," + strconv.Itoa(int(y)) + ")"

	return strings.Replace(line, posStr, newPosStr, 1)
}
