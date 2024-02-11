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
    configFile := "/etc/suncron.conf"
    sunCronFile := "/etc/suncron.cron"
    cronFilePath := "/etc/cron.d/suncron"

    lat, lng, timezone := loadConfig(configFile)

    today := time.Now()
    sunrise, sunset := sunrise.SunriseSunset(
        lat, lng,
        today.Year(), today.Month(), today.Day(),
    )
    tz, err := time.LoadLocation(timezone)
    if err != nil {
        fmt.Println(err)
        return
    }
    sunrise.In(tz)
    fmt.Println("Sunrise: " + sunrise.In(tz).Format("15:04:05"))
    fmt.Println("Sunset: " + sunset.In(tz).Format("15:04:05"))

    file, err := os.Open(sunCronFile)
    if err != nil {
        fmt.Println(err)
    }
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines);

    var hour int
    var minute int

    var cronRecords []string

    for scanner.Scan() {
        line := scanner.Text()
        r := strings.Split(line, ";")
        if len(r) != 4 {
            fmt.Println("Error: Invalid entry: " + line)
            continue
        }
        base    := strings.Trim(r[0], " ")
        inptime := strings.Trim(r[1], " ")
        dmw     := strings.Trim(r[2], " ")
        cmd     := strings.Trim(r[3], " ")
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

        cronRecords = append(cronRecords, fmt.Sprintf("%d %d %s %s\n", endTime.Minute(), endTime.Hour(), dmw, cmd))
    }

    writeCron(cronFilePath, cronRecords)

    file.Close()
}

func writeCron(file string, records []string) {
    // TODO: error handling
    f, err := os.Create(file)
    if err != nil {
        fmt.Println(err)
    }
    for idx := range records {
        _, err := f.WriteString(records[idx])
        if err != nil {
            fmt.Println(err)
        }
    }
    f.Close()
}

func loadConfig(configFile string) (float64, float64, string) {
    file, err := os.Open(configFile)
    if err != nil {
        panic(err)
    }
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines);
    var line string
    var lat, lng float64
    var timezone string
    for scanner.Scan() {
        line = scanner.Text()
        line = strings.TrimSpace(line)
        if !strings.Contains(line, "=") {
            continue
        }
        kv := strings.Split(line, "=")
        if (kv[0] == "timezone") {
            timezone = kv[1]
        }
        if (kv[0] == "location") {
            if !strings.Contains(kv[1], ",") {
            fmt.Println("Invalid coordinates: " + kv[1])
                continue
            }
            ll := strings.Split(kv[1], ",")
            lat, _ = strconv.ParseFloat(ll[0], 32)
            lng, _ = strconv.ParseFloat(ll[1], 32)
        }
    }
    return lat, lng, timezone
}
