package controllers

import (
	"Monitoring-service/controllers/monitoring" // Import the monitoring package
	"bufio"
	"fmt"
	"os"
)

var exitChan = make(chan struct{})

func Menu() {
	scanner := bufio.NewScanner(os.Stdin)
	for {

		fmt.Println("")
		fmt.Println("--------------------------")
		fmt.Println("1. Service Availability")
		fmt.Println("--------------------------")
		fmt.Println("")
		fmt.Println("--------------------------")
		fmt.Println("2. Req/res ratio")
		fmt.Println("--------------------------")

		fmt.Println("Enter choice:")
		//Registers choice and executes coresponding code
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "1":
			go exitListener()
			DisplayAvailability()
		case "2":
			go exitListener()
			DisplayReqResMenu()
		default:
			return
		}
	}

}

// Displays live updates for service avalability
func DisplayAvailability() {

	moveUp := "\033[A"
	moveDown := "\033[B"
	lineClear := "\033[K"
	colorGreen := "\x1b[32m"
	colorRed := "\x1b[31m"
	resetTextStyle := "\x1b[0m"

	fmt.Println("")
	fmt.Println("Press ENTER to exit")
	fmt.Println("--------------------")
	fmt.Println("Clinic service ...")
	fmt.Println("User service ...")
	fmt.Println("Appointment service ...")
	for {
		select {
		case flag := <-UserFlag:

			//Move one line up
			fmt.Print(moveUp)
			fmt.Print(moveUp)
			//Clear line
			fmt.Print(lineClear)

			if flag {
				// Makes text green
				fmt.Print(colorGreen + "User service" + resetTextStyle)
			} else {
				//Makes text red
				fmt.Print(colorRed + "User service" + resetTextStyle)
			}

			//Move one line down
			fmt.Print(moveDown)
			fmt.Print(moveDown)

			fmt.Print("\r")
		case flag := <-AppointmentFlag:
			//Move one line up
			fmt.Print(moveUp)
			//Clear line
			fmt.Print(lineClear)

			if flag {
				// Makes text green
				fmt.Print(colorGreen + "Appointment service" + resetTextStyle)
			} else {
				//Makes text red
				fmt.Print(colorRed + "Appointment service" + resetTextStyle)
			}

			//Move one line down
			fmt.Print(moveDown)

			fmt.Print("\r")
		case <-exitChan:
			return
		}

	}
}

func DisplayReqResMenu() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Select a service for Req/Res ratio:")
	fmt.Println("1. Appointment Service")
	fmt.Println("2. User Service")
	fmt.Println("3. Clinic Service")
	fmt.Println("4. Back to main menu")

	scanner.Scan()
	serviceChoice := scanner.Text()

	switch serviceChoice {
	case "1":
		DisplayReqRes("AppointmentService")
	case "2":
		DisplayReqRes("UserService")
	case "3":
		DisplayReqRes("ClinicService")
	case "4":
		return
	default:
		fmt.Println("Invalid choice. Please enter a valid option.")
		DisplayReqResMenu()
	}
}

func DisplayReqRes(service string) {
	colorGreen := "\x1b[32m"
	resetTextStyle := "\x1b[0m"

	percentage, err := monitoring.CalculatePercentage(service)
	if err != nil {
		fmt.Println("Error calculating percentage:", err)
		return
	}

	fmt.Println(fmt.Sprintf(colorGreen+"Request to response ratio for %s: %.2f%%"+resetTextStyle, service, percentage))

	<-exitChan
}

func exitListener() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	exitChan <- struct{}{}
}
