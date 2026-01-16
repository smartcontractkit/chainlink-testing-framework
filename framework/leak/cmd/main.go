package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	StepTick = 3 * time.Minute
)

func main() {
	fmt.Println("Starting CPU and Memory hog...")

	workersSchedule := os.Getenv("WORKERS")
	memorySchedule := os.Getenv("MEMORY")
	repeatStr := os.Getenv("REPEAT")

	leaks := make([][]byte, 0)
	workerCounter := 0

	go scheduleMemoryLeaks(memorySchedule, parseRepeat(repeatStr), &leaks)
	go scheduleCPUWorkers(workersSchedule, parseRepeat(repeatStr), &workerCounter)

	select {}
}

func parseRepeat(repeatStr string) int {
	if repeatStr == "" {
		return 1
	}
	repeat, _ := strconv.Atoi(repeatStr)
	if repeat < 1 {
		return 1
	}
	return repeat
}

func scheduleMemoryLeaks(schedule string, repeat int, leaks *[][]byte) {
	if schedule == "" {
		return
	}

	levels := parseSchedule(schedule)

	for r := 0; r < repeat; r++ {
		for _, target := range levels {
			timer := time.NewTimer(StepTick)

			current := len(*leaks) / 100

			if target > current {
				for i := 0; i < target-current; i++ {
					leak := make([]byte, 100*1024*1024)
					for j := range leak {
						leak[j] = byte(j % 256)
					}
					*leaks = append(*leaks, leak)
				}
			} else if target < current {
				removeCount := current - target
				if removeCount > len(*leaks)/100 {
					removeCount = len(*leaks) / 100
				}
				*leaks = (*leaks)[:len(*leaks)-removeCount*100]
			}

			log.Printf("Memory: %dx100MB", len(*leaks)/100)
			<-timer.C
		}
	}
}

func scheduleCPUWorkers(schedule string, repeat int, counter *int) {
	if schedule == "" {
		return
	}

	levels := parseSchedule(schedule)
	activeWorkers := 0

	for r := 0; r < repeat; r++ {
		for _, target := range levels {
			timer := time.NewTimer(StepTick)

			if target > activeWorkers {
				for i := 0; i < target-activeWorkers; i++ {
					*counter++
					go cpuWorker(*counter)
				}
			} else if target < activeWorkers {
				stopWorkers(activeWorkers - target)
			}

			activeWorkers = target
			log.Printf("CPU Workers: %d", activeWorkers)
			<-timer.C
		}
	}
}

var stopChan = make(chan bool, 1000)

func cpuWorker(id int) {
	for {
		select {
		case <-stopChan:
			return
		default:
			for n := 2; n < 100000; n++ {
				isPrime := true
				for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
					if n%i == 0 {
						isPrime = false
						break
					}
				}
				_ = isPrime
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func stopWorkers(count int) {
	for i := 0; i < count; i++ {
		select {
		case stopChan <- true:
		default:
		}
	}
}

func parseSchedule(schedule string) []int {
	parts := strings.Split(schedule, ",")
	var result []int
	for _, part := range parts {
		val, _ := strconv.Atoi(strings.TrimSpace(part))
		if val >= 0 {
			result = append(result, val)
		}
	}
	return result
}
