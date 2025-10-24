package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseStringOfInts(s string) ([]int, error) {
	nums := []int{}
	for _, num := range strings.Fields(s) {
		numInt, err := strconv.Atoi(num)
		if err != nil || numInt > 12 {
			return nil, fmt.Errorf("invalid domino value: %v", num)
		}
		nums = append(nums, numInt)
	}
	return nums, nil
}

func parseDomino(s string) (Domino, error) {
	nums, err := parseStringOfInts(s)
	if err != nil {
		return Domino{}, err
	}
	if len(nums) != 2 {
		return Domino{}, fmt.Errorf("expected 2 values for line '%v'", s)
	}
	return Domino{nums[0], nums[1]}, nil
}

func readTrainFile(filename string) (availablePlayValues []int, dominoes []Domino, err error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) < 2 {
		return nil, nil, errors.New("expected at least two lines")
	}

	availablePlayValues, err = parseStringOfInts(lines[0])
	if err != nil {
		return nil, nil, err
	}
	if len(availablePlayValues) == 0 {
		return nil, nil, errors.New("no available play values given")
	}

	dominoLines := lines[1:]
	for _, line := range dominoLines {
		if line == "" {
			continue
		}
		domino, err := parseDomino(line)
		if err != nil {
			return nil, nil, err
		}
		dominoes = append(dominoes, domino)
	}
	if len(dominoes) == 0 {
		return nil, nil, errors.New("no dominoes found")
	}

	return availablePlayValues, dominoes, nil
}

func getFileNameFromArgs() string {
	if len(os.Args) == 1 {
		printError(errors.New("must specify a filepath for the game state"))
	}
	if len(os.Args) > 2 {
		printError(errors.New("too many arguments"))
	}
	return os.Args[1]
}

func printError(err error) {
	// ANSI escape codes for red color and reset
	redColor := "\x1b[31m"
	resetColor := "\x1b[0m"
	fmt.Fprintf(os.Stderr, "%vERROR: %v%v\n", redColor, err.Error(), resetColor)
}
