package main

import (
	"bufio"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Game struct {
	Id    int  `json:"id"`
	Red   int  `json:"red"`
	Blue  int  `json:"blue"`
	Green int  `json:"green"`
	Valid bool `json:"valid"`
}

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

	// Create an array of games
	var games []Game

	for _, line := range lines {
		game, err := parseLine(line)
		if err != nil {
			log.Debugf("Couldn't parse line: %v", err)
		}

		// validate if it's a possible game
		games = append(games, game)
	}

	log.Infof("The solution is: %v", sumGameIds(games))
}

func sumGameIds(games []Game) int {
	var result int
	for _, game := range games {
		if game.Valid {
			log.Debugf("valid game at number  %v", game.Id)
			result += game.Id
		}
	}
	return result
}

// parseLine reads a line as string and parses into a Game
func parseLine(line string) (Game, error) {
	var game Game
	var err error

	// assume every game is valid at the beginning
	game.Valid = true

	parts := strings.Split(line, ":")
	log.Debugf("Parts after first split: %v", parts)
	if len(parts) != 2 {
		return game, fmt.Errorf("invalid line format")
	}

	// Match a number
	expression := regexp.MustCompile(`\d+`)
	gameNumber := expression.FindAllString(parts[0], -1)
	game.Id, err = strconv.Atoi(gameNumber[0])
	if err != nil {
		log.Debugf("Got game number %v", gameNumber)
		log.Debugf("Using %v as gameID", gameNumber[0])
		log.Debugf("Failed to parse gameID")
		return game, err
	}

	// split on each semicolon for sets of cubes
	// example structure looks like this now: ["13 blue", "6 green, 4 blue, 11 red", "15 red"]
	cubeSets := strings.Split(parts[1], ";")

	// split every cubeSet into it's color/value pairs
	for _, cubeSet := range cubeSets {
		pairs := strings.Split(cubeSet, ",")
		// pairs is now like so [" 6 green", " 4 blue", " 11 red"]

		for _, pair := range pairs {
			colorValue := strings.Split(pair, " ")
			count, err := strconv.Atoi(colorValue[1])
			if err != nil {
				log.Debugf("Couldn't convert color value")
				log.Debugf("Received %v", colorValue[1])
				return game, err
			}

			// Add count depending on value we got.
			switch colorValue[2] {
			case "red":
				if count > 12 {
					log.Debugf("Count %v for %v, resulted in invalid game", count, colorValue[2])
					game.Valid = false
				}

				game.Red += count
			case "blue":
				if count > 14 {
					log.Debugf("Count %v for %v, resulted in invalid game", count, colorValue[2])
					game.Valid = false
				}

				game.Blue += count
			case "green":
				if count > 13 {
					log.Debugf("Count %v for %v, resulted in invalid game", count, colorValue[2])
					game.Valid = false
				}

				game.Green += count
			}
			log.Debugf("Game validation is %v", game.Valid)

		}
	}
	return game, nil
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
