package main

import (
    "github.com/nathan-osman/go-sunrise"
    "os"
    "fmt"
    "bufio"
    "time"
    "strings"
    "strconv"
)

func main() {
    sunCronFilePath := "sun.cron"

    today := time.Now()
    sunrise, sunset := sunrise.SunriseSunset(
        -37.864, 144.982,
        today.Year(), today.Month(), today.Day(),
    )
    tz, _ := time.LoadLocation("Australia/Melbourne")
    sunrise.In(tz)
    fmt.Println("Location:", today.Location(), ":Time:", today)
    fmt.Println("Sunrise: " + sunrise.In(tz).Format("15:04:05"))
    fmt.Println("Sunset: " + sunset.In(tz).Format("15:04:05"))

    file, err := os.Open(sunCronFilePath)
    if err != nil {
        fmt.Println(err)
    }
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines);

    fmt.Println("Reading " + sunCronFilePath)
    var hour int
    var minute int

    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println("Line: " + line)
        r := strings.Split(line, ";")
        if len(r) != 3 {
            fmt.Println("Error: Invalid entry: " + line)
            continue
        }
        base := strings.Trim(r[0], " ")
        inptime := strings.Trim(r[1], " ")
        cmd  := strings.Trim(r[2], " ")
        if base != "sunrise" && base != "sunset" {
            fmt.Println("Error: Invalid value: " + base)
            continue
        }

        var mod = 1
        if len(inptime) > 0 && inptime[0] != '-' && inptime[0] != '+' {
            fmt.Println("Error: Invalid time: " + inptime)
            continue
        }
        if len(inptime) == 1 && inptime[0] == '-' {
            mod = -1
        }

        timeSplice := strings.Split(inptime, ":")
        if len(timeSplice) == 0 || len(timeSplice) > 2 {
            fmt.Println("Error: Invalid time: " + inptime)
            continue
        }

        minute = 0
        hour = 0
        if len(timeSplice) > 0 {
            hour, _ = strconv.Atoi(timeSplice[0])
        }
        if len(timeSplice) == 2 {
            minute, _ = strconv.Atoi(timeSplice[1])
        }

        var endTime time.Time
        if base == "sunset" {
            endTime = sunset.In(tz)
        } else {
            endTime = sunrise.In(tz)
        }
        endTime = endTime.Add(time.Hour * time.Duration(hour * mod) + time.Minute * time.Duration(minute * mod))

        fmt.Printf("result: %s %d:%d => %s: %s\n", base, hour, minute, endTime, cmd)
        fmt.Printf("%d %d * * * %s\n", endTime.Minute(), endTime.Hour(), cmd)
    }

    file.Close()
}
