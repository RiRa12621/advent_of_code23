package main

import (
	"bufio"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	filePtr := flag.String("input-file", "", "input file to read")
	debugPtr := flag.Bool("debug", false, "decide if debug should be enabled")

	// Parse the flags
	flag.Parse()
	// check if we want debug logs
	if *debugPtr == true {
		log.SetLevel(log.DebugLevel)
	}
	// check if a file is specified
	if *filePtr == "" {
		log.Fatalf("You need to specify an input file")
	}

	// read the given file
	file, err := os.Open(*filePtr)
	if err != nil {
		log.Fatalf("Couldn't open file: %v", err)
	}
	defer file.Close()
	lines, err := readLines(file)
	if err != nil {
		log.Fatalf("Couldn't read lines in file: %v", err)
	}
	var numbers []int

	// go over each line and get the number from it, then add the number to the numbers
	for linenumber, line := range lines {
		thisSum, err := getNum(line)
		if err != nil {
			log.Fatalf("Couldn't get number for line %v: %v", linenumber, err)
		}
		log.Debugf("Line %v, has number %v", linenumber, thisSum)
		numbers = append(numbers, thisSum)
	}

	// add the numbers up
	log.Debugf("Total amount of numbers: %v", len(numbers))
	result := sumArray(numbers)
	log.Infof("The result is %v", result)
}

// readLines read a line from a file and adds each lines as string element into an array
func readLines(file *os.File) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// getNum gets concatenated number from a string and returns it as int
func getNum(line string) (int, error) {
	// regex match all numbers

	log.Debugf("Line before replacement: %v", line)
	replacements := map[string]string{
		"one":   "o1e",
		"two":   "t2o",
		"seven": "s7n",
		"eight": "e8t",
	}

	for word, replacement := range replacements {
		line = regexp.MustCompile(word).ReplaceAllString(line, replacement)
	}
	log.Debugf("Line after replacement: %v", line)

	expression := regexp.MustCompile("[0-9]|one|two|three|four|five|six|seven|eight|nine")

	// find the all numbers in the string
	allNumbers := expression.FindAllString(line, -1)
	log.Debugf("All numbers in line: %v", allNumbers)

	// just get the first and last number
	numbers := []string{allNumbers[0], allNumbers[len(allNumbers)-1]}
	log.Debugf("Got number %v", numbers)

	// Transform words to digits
	// First digit
	switch numbers[0] {
	case "one":
		numbers[0] = "1"
	case "two":
		numbers[0] = "2"
	case "three":
		numbers[0] = "3"
	case "four":
		numbers[0] = "4"
	case "five":
		numbers[0] = "5"
	case "six":
		numbers[0] = "6"
	case "seven":
		numbers[0] = "7"
	case "eight":
		numbers[0] = "8"
	case "nine":
		numbers[0] = "9"
	default:
		numbers[0] = numbers[0]
	}

	// Second digit
	switch numbers[1] {
	case "one":
		numbers[1] = "1"
	case "two":
		numbers[1] = "2"
	case "three":
		numbers[1] = "3"
	case "four":
		numbers[1] = "4"
	case "five":
		numbers[1] = "5"
	case "six":
		numbers[1] = "6"
	case "seven":
		numbers[1] = "7"
	case "eight":
		numbers[1] = "8"
	case "nine":
		numbers[1] = "9"
	default:
		numbers[1] = numbers[1]
	}

	// concatenate the strings
	var str strings.Builder
	for _, number := range numbers {
		str.WriteString(number)
	}
	// return int converted from string
	result, err := strconv.Atoi(str.String())
	if err != nil {
		return result, err
	}
	return result, nil
}

// sumArray sums all elements in an array of integers and returns the result as new int
func sumArray(sums []int) int {
	result := 0

	for _, number := range sums {
		result += number
	}
	return result
}
