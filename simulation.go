package main

import (
	"fmt"

	"math"

	"math/rand"

	"time"
)

// Constants

const (
	NumPatients = 10

	NumStations = 5

	StartTime = 13 * 60 // 13:00 in minutes

	EndTime = 19 * 60 // 19:00 in minutes

	InitialTemp = 1000.0 // Initial temperature for simulated annealing

	FinalTemp = 0.01 // Final temperature

	AnnealingRate = 0.995 // Annealing rate

)

// Station times (in minutes)

var stationTimes = map[string]int{

	"General Practitioner": 30,

	"Dietitian": 30,

	"Orthopedist": 30,

	"Blood Sample": 20,

	"Electrocardiogram": 20,
}

type Schedule struct {
	Visits [NumPatients][NumStations]int // Stores the order of visits for each patient to each station

}

func main() {

	// Initialize random seed

	rand.Seed(time.Now().UnixNano())

	// Start with an initial random schedule

	initialSchedule := generateRandomSchedule()

	// Perform Simulated Annealing

	bestSchedule := simulatedAnnealing(initialSchedule)

	// Print the final best schedule and its associated cost

	printSchedule(bestSchedule)

}

func generateRandomSchedule() Schedule {

	var schedule Schedule

	for i := 0; i < NumPatients; i++ {

		// Randomize the order of stations for each patient

		for j := 0; j < NumStations; j++ {

			schedule.Visits[i][j] = j

		}

		// Shuffle stations for each patient

		rand.Shuffle(NumStations, func(k, l int) {

			schedule.Visits[i][k], schedule.Visits[i][l] = schedule.Visits[i][l], schedule.Visits[i][k]

		})

	}

	return schedule

}

// Objective function to calculate the total waiting time for a given schedule

func calculateCost(schedule Schedule) int {

	// Initialize station availability times

	stationAvailable := make([]int, NumStations)

	// Initialize patient arrival times (when they can start at each station)

	//patientArrival := make([]int, NumPatients)

	// Initialize total waiting time

	totalWaitingTime := 0

	// Process each patient's visits to each station

	for i := 0; i < NumPatients; i++ {

		arrivalTime := StartTime // All patients start at 13:00

		for j := 0; j < NumStations; j++ {

			station := schedule.Visits[i][j]

			processingTime := stationTimes[stationName(station)]

			// Update patient's arrival time considering station availability

			if arrivalTime < stationAvailable[station] {

				arrivalTime = stationAvailable[station]

			}

			// Add waiting time to the total waiting time

			waitingTime := arrivalTime - StartTime

			totalWaitingTime += waitingTime

			// Update station availability and patient's next arrival time

			stationAvailable[station] = arrivalTime + processingTime

			arrivalTime += processingTime

		}

	}

	return totalWaitingTime

}

func stationName(station int) string {

	stations := []string{

		"General Practitioner",

		"Dietitian",

		"Orthopedist",

		"Blood Sample",

		"Electrocardiogram",
	}

	return stations[station]

}

// Simulated Annealing algorithm

func simulatedAnnealing(initialSchedule Schedule) Schedule {

	currentSchedule := initialSchedule

	currentCost := calculateCost(currentSchedule)

	bestSchedule := currentSchedule

	bestCost := currentCost

	temperature := InitialTemp

	for temperature > FinalTemp {

		// Generate a new schedule by making a small change to the current schedule

		newSchedule := currentSchedule

		makeSmallChange(&newSchedule)

		// Calculate the cost of the new schedule

		newCost := calculateCost(newSchedule)

		// If the new schedule is better, accept it

		if newCost < currentCost {

			currentSchedule = newSchedule

			currentCost = newCost

			// If the new schedule is better than the best found so far, update the best schedule

			if newCost < bestCost {

				bestSchedule = newSchedule

				bestCost = newCost

			}

		} else {

			// Otherwise, accept it with a certain probability

			probability := math.Exp(float64(currentCost-newCost) / temperature)

			if rand.Float64() < probability {

				currentSchedule = newSchedule

				currentCost = newCost

			}

		}

		// Cool down the temperature

		temperature *= AnnealingRate

	}

	return bestSchedule

}

// Make a small random change to a schedule (swap two station visits for two patients)

func makeSmallChange(schedule *Schedule) {

	patient1 := rand.Intn(NumPatients)

	patient2 := rand.Intn(NumPatients)

	station1 := rand.Intn(NumStations)

	station2 := rand.Intn(NumStations)

	// Swap the visits for two patients at two random stations

	schedule.Visits[patient1][station1], schedule.Visits[patient2][station2] = schedule.Visits[patient2][station2], schedule.Visits[patient1][station1]

}

// Print the schedule in a readable format

func printSchedule(schedule Schedule) {

	for i := 0; i < NumPatients; i++ {

		fmt.Printf("Patient %d:\n", i+1)

		arrivalTime := StartTime

		for j := 0; j < NumStations; j++ {

			station := schedule.Visits[i][j]

			processingTime := stationTimes[stationName(station)]

			startTime := arrivalTime

			endTime := arrivalTime + processingTime

			fmt.Printf("  %s: %02d:%02d to %02d:%02d\n", stationName(station), startTime/60, startTime%60, endTime/60, endTime%60)

			arrivalTime = endTime

		}

	}

}
