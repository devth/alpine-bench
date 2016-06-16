package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	cmd          = "ab"
	postFilePath = "/tmp/ab_post_file"
	putFilePath  = "/tmp/ab_put_file"
)

func check(e error) {
	if e != nil {
		// log.Fatal(e)
		panic(e)
	}
}

func toFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 32)
	check(err)
	return f
}

func extractLatency(lines []string) map[string]map[string]float64 {

	// Matches lines like: Connect:        0    1   0.1      1       1
	statPattern, _ := regexp.Compile(`(\w+)\:\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)`)
	percPattern, _ := regexp.Compile(`\s+(\d+)\%\s+(\d+)`)

	latency := map[string]map[string]float64{
		"percentiles": map[string]float64{},
	}

	for _, v := range lines {
		// stat
		stat := statPattern.FindStringSubmatch(v)
		if stat != nil {
			log.Println("Matched stat line", stat)
			latency[strings.ToLower(stat[1])] = map[string]float64{
				"min":    toFloat(stat[2]),
				"mean":   toFloat(stat[3]),
				"stdDev": toFloat(stat[4]),
				"median": toFloat(stat[5]),
				"max":    toFloat(stat[6]),
			}
		}
		// percentile
		perc := percPattern.FindStringSubmatch(v)
		if perc != nil {
			log.Println("Matched percentline", perc)
			latency["percentiles"][perc[1]] = toFloat(perc[2])
		}
	}
	// log.Println(latency)
	return latency
}

func runAb(cmd string, args []string) map[string]map[string]float64 {
	log.Println("Running:", cmd, strings.Join(args, " "))
	// out, err := exec.Command(cmd, args...).Output()
	out, err := exec.Command(cmd, args...).CombinedOutput()
	results := fmt.Sprintf("%s", out)
	log.Println(results)
	check(err)
	lines := strings.Split(results, "\n")
	// fmt.Printf("%d lines\n", len(lines))
	latency := extractLatency(lines)
	return latency
}

func toJson(latency map[string]map[string]float64) string {
	jsonString, jsonErr := json.Marshal(latency)
	check(jsonErr)
	return fmt.Sprintf("%s", jsonString)
}

func main() {
	args := os.Args[1:]

	post := os.Getenv("POST")
	if post != "" {
		postBytes := []byte(post)
		log.Println("POST is defined, setting -p")
		err := ioutil.WriteFile(postFilePath, postBytes, 0644)
		check(err)
		args = append([]string{"-p", postFilePath}, args...)
	}

	put := os.Getenv("PUT")
	if put != "" {
		log.Println("PUT is defined, setting -u")
		err := ioutil.WriteFile(putFilePath, []byte(put), 0644)
		check(err)
		args = append([]string{"-u", putFilePath}, args...)
	}

	if post != "" && put != "" {
		panic("Cannot mix PUT and POST")
	}

	latency := runAb(cmd, args)

	// This is the only line that prints to stdout, containing latency JSON
	fmt.Println(toJson(latency))
}
