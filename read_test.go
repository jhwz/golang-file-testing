package ioframe_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/ioframe"
)

const filename = "test.txt"

func parseLine(in chan string, out chan int) error {
	for line := range in {
		num, err := strconv.Atoi(strings.TrimRight(line, "\n"))
		if err != nil {
			return err
		}
		out <- num
	}
	return nil
}

func TestRead(t *testing.T) {
	out := ioframe.NewFrame(filename, parseLine).ReadFile()
	fmt.Println(out)

}

func BenchmarkFrame(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ioframe.NewFrame(filename, parseLine).ReadFile()
	}
}

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ioframe.NormalFileRead(filename)
	}
}
