package controllers

import (
	"bufio"
	"fmt"
	"os"
)

var exitChan = make(chan struct{})

func Menu(){
    scanner := bufio.NewScanner(os.Stdin)
    for{

        fmt.Println("")
        fmt.Println("--------------------------")
        fmt.Println("1. Service Availability")
        fmt.Println("--------------------------")

        fmt.Println("Enter choice:")
        //Registers choice and executes coresponding code
        scanner.Scan()
        input := scanner.Text()
        switch input{
        case "1":
            go exitListener()
            DisplayAvailability()
        default:
            return
        }
    }

}

// Displays live updates for service avalability
func DisplayAvailability(){

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
    for{
        select{
        case flag := <- ClinicFlag:

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
            } else{
                //Makes text red
                fmt.Print(colorRed + "Clinic service" + resetTextStyle)
            }

            //Move one line down
            fmt.Print(moveDown)
            fmt.Print(moveDown)
            fmt.Print(moveDown)
            fmt.Print(moveDown)

            fmt.Print("\r")

        case flag := <- NotificationFlag:

            //Move one line up
            fmt.Print(moveUp)
            fmt.Print(moveUp)
            fmt.Print(moveUp)
            //Clear line
            fmt.Print(lineClear)

            if flag {
                // Makes text green
                fmt.Print(colorGreen + "Notification service" + resetTextStyle)
            } else{
                //Makes text red
                fmt.Print(colorRed + "Notification service" + resetTextStyle)
            }

            //Move one line down
            fmt.Print(moveDown)
            fmt.Print(moveDown)
            fmt.Print(moveDown)

            fmt.Print("\r")
        case flag := <- UserFlag:

            //Move one line up
            fmt.Print(moveUp)
            fmt.Print(moveUp)
            //Clear line
            fmt.Print(lineClear)

            if flag {
                // Makes text green
                fmt.Print(colorGreen + "User service" + resetTextStyle)
            } else{
                //Makes text red
                fmt.Print(colorRed + "User service" + resetTextStyle)
            }

            //Move one line down
            fmt.Print(moveDown)
            fmt.Print(moveDown)

            fmt.Print("\r")
        case flag := <- AppointmentFlag:
            //Move one line up
            fmt.Print(moveUp)
            //Clear line
            fmt.Print(lineClear)

            if flag {
                // Makes text green
                fmt.Print(colorGreen + "Appointment service" + resetTextStyle)
            } else{
                //Makes text red
                fmt.Print(colorRed + "Appointment service" + resetTextStyle)
            }

            //Move one line down
            fmt.Print(moveDown)

            fmt.Print("\r")
        case <- exitChan:
            return
        }

    }
}

func exitListener(){
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    exitChan <- struct{}{}
}


