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

func DisplayAvailability(){


    fmt.Println("")
    fmt.Println("Press ENTER to exit")
    fmt.Println("--------------------")
    fmt.Println("User service ...")
    fmt.Println("Appointment service ...")
    for{
        select{
        case flag := <- UserFlag:

            //Move one line up
            fmt.Print("\033[A")
            fmt.Print("\033[A")
            //Clear line
            fmt.Print("\033[K")

            if flag {
                fmt.Print("\x1b[32mUser service\x1b[0m")
            } else{
                fmt.Print("\x1b[31mUser service\x1b[0m")
            }

            //Move one line down
            fmt.Print("\033[B")
            fmt.Print("\033[B")

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


