package controllers

import "strings"

//timeslots, availabletimes, appointment = AppointmentService
//dentist, patient = UserService


func GetReqRes(topic string) string {
    // Done like this to make use of already existing method
    res := []string{"res"}
    if (containsAny(topic, res)) {
        return "res"
    }
    req := []string{"req"}
    if (containsAny(topic, req)) {
        return "req"
    }

    return ""
}

func GetService(topic string) string {

    appointmentTopics := []string{"timeslots", "availabletimes", "appointment"}
    if (containsAny(topic, appointmentTopics)) {
        return "AppointmentService"
    }

    userTopics := []string{"dentists", "patients"}
    if (containsAny(topic, userTopics)) {
        return "UserService"
    }

    // TODO Not sure about the topics here
    clinicTopics := []string{"clinics"}
    if (containsAny(topic, clinicTopics)) {
        return "ClinicService"
    }

    return "Unknown"
}

func containsAny(str string, substrings []string) bool {
    for _,sub := range substrings {
        if strings.Contains(str, sub){
            return true
        }
    }
    return false
}
