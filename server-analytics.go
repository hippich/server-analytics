package main

import "os"
import "time"
import "github.com/jpillora/go-ogle-analytics"
import linuxproc "github.com/c9s/goprocinfo/linux"

func main() {
    uaid := os.Getenv("UAID")

    if uaid == "" {
        panic("UAID environment required. Example: UAID=UA-XXXXXX-Y ./server-analytics")
    }

    client, err := ga.NewClient(uaid)

    if err != nil {
        panic(err)
    }

    ticker := time.NewTicker(15 * time.Second)
    quit := make(chan struct{})
    stop := false

    for !stop {
       select {
        case <- ticker.C:
            diskInfo, err := linuxproc.ReadDisk("/")
            _ = err

            client.Send(ga.NewEvent("Disk", "All").Label("Total Disk Space").Value(int64(diskInfo.All)))
            client.Send(ga.NewEvent("Disk", "Used").Label("Used Disk Space").Value(int64(diskInfo.Used)))
            client.Send(ga.NewEvent("Disk", "Free").Label("Free Disk Space").Value(int64(diskInfo.Free)))

            loadAvg, err := linuxproc.ReadLoadAvg("/proc/loadavg")
            _ = err

            client.Send(ga.NewEvent("Load Average", "1m").Label("1m").Value(int64(loadAvg.Last1Min)))
            client.Send(ga.NewEvent("Load Average", "5m").Label("5m").Value(int64(loadAvg.Last5Min)))
            client.Send(ga.NewEvent("Load Average", "15m").Label("15m").Value(int64(loadAvg.Last15Min)))
            client.Send(ga.NewEvent("Load Average", "Running").Label("Process Running").Value(int64(loadAvg.ProcessRunning)))
            client.Send(ga.NewEvent("Load Average", "Total").Label("Process Total").Value(int64(loadAvg.ProcessTotal)))

            meminfo, err := linuxproc.ReadMemInfo("/proc/meminfo")
            _ = err

            client.Send(ga.NewEvent("Memory", "Total").Label("Memory Total").Value(int64(meminfo.MemTotal)))
            client.Send(ga.NewEvent("Memory", "Free").Label("Memory Free").Value(int64(meminfo.MemFree)))

        case <- quit:
            ticker.Stop()
            return
        }
    }
}
