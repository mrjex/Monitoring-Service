package controllers

import (
	"Monitoring-service/controllers/monitoring" // Import the monitoring package
	"bufio"
	"fmt"
	"os"
	"strings"
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
		fmt.Println("q. Shut down")
		fmt.Println("--------------------------")

		fmt.Println("Enter choice:")
		//Registers choice and executes coresponding code
		fmt.Println("Enter choice:")
		if !scanner.Scan() {
			fmt.Println("Error reading input.")
			os.Exit(1)
		}
		input := scanner.Text()

		switch input {
		case "1":
			go exitListener()
			DisplayAvailability()
		case "2":
			go exitListener()
			DisplayAllReqRes()
		case "q":
			os.Exit(0)
		default:
			fmt.Println("")
			fmt.Println("Please enter valid option")

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
	fmt.Println("Notification service ...")
	fmt.Println("User service ...")
	fmt.Println("Appointment service ...")
	for {
		select {
		case flag := <-ClinicFlag:

			//Move one line up
			fmt.Print(moveUp)
			fmt.Print(moveUp)
			fmt.Print(moveUp)
			fmt.Print(moveUp)
			//Clear line
			fmt.Print(lineClear)

			if flag {
				// Makes text green
				fmt.Print(colorGreen + "Clinic service" + resetTextStyle)
			} else {
				//Makes text red
				fmt.Print(colorRed + "Clinic service" + resetTextStyle)
			}

			//Move one line down
			fmt.Print(moveDown)
			fmt.Print(moveDown)
			fmt.Print(moveDown)
			fmt.Print(moveDown)

			fmt.Print("\r")

		case flag := <-NotificationFlag:

			//Move one line up
			fmt.Print(moveUp)
			fmt.Print(moveUp)
			fmt.Print(moveUp)
			//Clear line
			fmt.Print(lineClear)

			if flag {
				// Makes text green
				fmt.Print(colorGreen + "Notification service" + resetTextStyle)
			} else {
				//Makes text red
				fmt.Print(colorRed + "Notification service" + resetTextStyle)
			}

			//Move one line down
			fmt.Print(moveDown)
			fmt.Print(moveDown)
			fmt.Print(moveDown)

			fmt.Print("\r")
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

func DisplayAllReqRes() {
	colorGreen := "\x1b[32m"
	resetTextStyle := "\x1b[0m"
	fmt.Println("Press ENTER to exit")
	fmt.Println("--------------------")
	// Display Req/Res ratio for each service
	displayReqRes("AppointmentService", colorGreen, resetTextStyle)
	displayReqRes("UserService", colorGreen, resetTextStyle)
	displayReqRes("ClinicService", colorGreen, resetTextStyle)

	<-exitChan
}

func displayReqRes(service string, colorGreen string, resetTextStyle string) {
	percentage, err := monitoring.CalculatePercentage(service)
	if err != nil {
		fmt.Printf("Error calculating percentage for %s: %v\n", service, err)
		return
	}

	fmt.Printf("%s%s: %.2f%% (res/req)%s\n", colorGreen, service, percentage, resetTextStyle)
}

func exitListener() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	exitChan <- struct{}{}
}

//timeslots, availabletimes, appointment = AppointmentService
//dentist, patient = UserService

func GetReqRes(topic string) string {
	// Done like this to make use of already existing method
	res := []string{"res"}
	if containsAny(topic, res) {
		return "res"
	}
	req := []string{"req"}
	if containsAny(topic, req) {
		return "req"
	}

	return ""
}

func GetService(topic string) string {

	appointmentTopics := []string{"timeslots", "availabletimes", "appointment"}
	if containsAny(topic, appointmentTopics) {
		return "AppointmentService"
	}

	userTopics := []string{"dentists", "patients"}
	if containsAny(topic, userTopics) {
		return "UserService"
	}

	// TODO Not sure about the topics here
	clinicTopics := []string{"clinics"}
	if containsAny(topic, clinicTopics) {
		return "ClinicService"
	}

	return "Unknown"
}

func containsAny(str string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(str, sub) {
			return true
		}
	}
	return false
}
