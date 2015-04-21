package main

import (

    "github.com/alexjlockwood/gcm"
    "fmt"
    "github.com/fjl/go-couchdb"
)

func HandleChangeEvent(ID string, db *couchdb.DB) (deviceTokens []string) {

    deviceTokens = []string{}

    var list List
    if err := db.Get(ID, &list, nil); err != nil {
        fmt.Println(err)
    } else {
        var profile Profile
        db.Get(list.Owner, &profile, nil)
        deviceTokens = append(deviceTokens, profile.DeviceTokens...)


        for _, userId := range list.Members {
            var profile Profile
            db.Get(userId, &profile, nil)
            deviceTokens = append(deviceTokens, profile.DeviceTokens...)
        }
    }

    var task Task
    if err := db.Get(ID, &task, nil); err != nil {
        fmt.Println(err)
    } else {
        var list List
        db.Get(task.ListId, &list, nil)

        var profile Profile
        db.Get(list.Owner, &profile, nil)

        deviceTokens = append(deviceTokens, profile.DeviceTokens...)

        for _, userId := range list.Members {
            var profile Profile
            db.Get(userId, &profile, nil)
            deviceTokens = append(deviceTokens, profile.DeviceTokens...)
        }
    }

    return deviceTokens
}


func NotifyUsers(deviceTokens []string) (error) {
    sender := &gcm.Sender{ApiKey: "AIzaSyBEHLA1FR4OlCRQE1vPv_mfqQaIF0ICZeA"}
    message := gcm.NewMessage(nil, deviceTokens...)
    _, err := sender.Send(message, 2)
    fmt.Println("Successfully sent a notification to device token : ", deviceTokens);
    if err != nil {
        return fmt.Errorf("Error sending notifications: %s", err)
    }
    return nil
}