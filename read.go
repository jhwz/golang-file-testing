package ioframe

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Packge works in three steps. Read - parse - store

// Parser is the type of function that needs to be provided for parsing step
type Parser func(in chan string, out chan int) error

// Frame is the parent object for ochestrating the io
type Frame struct {
	// Number of goroutines to use
	NumRoutines int

	// File to read
	Filename string

	err error

	inChan  chan string
	outChan chan int

	Parser Parser

	store []int
}

// NewFrame returns a frame from the elements we need to run
func NewFrame(filename string, f Parser) *Frame {
	return &Frame{
		Filename:    filename,
		NumRoutines: 8,
		inChan:      make(chan string, 50),
		outChan:     make(chan int, 50),
		Parser:      f,
	}
}

// ReadFile is a high-level function that ochestrates the sub-routines for reading the file
func (f *Frame) ReadFile() []int {
	var wg, wgp sync.WaitGroup

	wgp.Add(1)
	go func() {
		defer wgp.Done()
		f.populateStore()
	}()

	wg.Add(f.NumRoutines)
	for i := 0; i < f.NumRoutines; i++ {
		go func() {
			defer wg.Done()
			if err := f.Parser(f.inChan, f.outChan); err != nil {
				f.err = err
			}
		}()
	}

	f.readfile()
	close(f.inChan)
	wg.Wait()
	close(f.outChan)
	wgp.Wait()
	return f.store
}

func (f *Frame) readfile() {
	file, err := os.Open(f.Filename)
	if err != nil {
		f.err = err
		return
	}
	defer file.Close()

	bufr := bufio.NewReader(file)
	var line string
	for {
		line, err = bufr.ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			f.err = err
			return
		}

		f.inChan <- line
	}
}

func (f *Frame) populateStore() {
	for el := range f.outChan {
		f.store = append(f.store, el)
	}
}

// Error will be nil if successful, otherwise the reason for failure
func (f *Frame) Error() error {
	return f.err
}

// NormalFileRead basic approach
func NormalFileRead(f string) ([]int, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var out []int
	bufr := bufio.NewReader(file)
	var line string
	for {
		line, err = bufr.ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		conv, err := strconv.Atoi(strings.TrimRight(line, "\n"))
		if err != nil {
			return nil, err
		}
		out = append(out, conv)
	}
	return out, nil
}
