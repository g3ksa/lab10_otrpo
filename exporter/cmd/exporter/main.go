package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func main() {
	host := os.Getenv("EXPORTER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("EXPORTER_PORT")
	if port == "" {
		port = "8085"
	}

	http.HandleFunc("/", metricsHandler)

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting exporter on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting http server: %v", err)
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	cpuPercents, err := cpu.Percent(0, true)
	if err == nil && len(cpuPercents) > 0 {
		fmt.Fprintf(w, "# HELP system_cpu_usage CPU usage percentage per core (0-100)\n")
		fmt.Fprintf(w, "# TYPE system_cpu_usage gauge\n")
		for i, p := range cpuPercents {
			fmt.Fprintf(w, "metric_cpu_usage{cpu=\"%d\"} %f\n", i, p)
		}
	}

	mbModifier := uint64(1024 * 1024)
	vmStat, err := mem.VirtualMemory()
	if err == nil {
		fmt.Fprintf(w, "# HELP system_memory_megabytes System memory info in MB\n")
		fmt.Fprintf(w, "# TYPE system_memory_megabytes gauge\n")
		fmt.Fprintf(w, "metric_memory_usage{type=\"total\"} %d\n", vmStat.Total/mbModifier)
		fmt.Fprintf(w, "metric_memory_usage{type=\"used\"} %d\n", vmStat.Used/mbModifier)
	}

	gbModifier := uint64(1024 * 1024 * 1024)
	dstat, err := disk.Usage("/")
	if err == nil {
		fmt.Fprintf(w, "# HELP system_disk_gigabytes System disk usage in GB\n")
		fmt.Fprintf(w, "# TYPE system_disk_gigabytes gauge\n")
		fmt.Fprintf(w, "metric_disk_usage{mountpoint=\"/\",type=\"total\"} %d\n", dstat.Total/gbModifier)
		fmt.Fprintf(w, "metric_disk_usage{mountpoint=\"/\",type=\"used\"} %d\n", dstat.Used/gbModifier)
	}

	fmt.Fprintf(w, "# HELP system_goroutines Number of Go goroutines\n")
	fmt.Fprintf(w, "# TYPE system_goroutines gauge\n")
	fmt.Fprintf(w, "system_goroutines %d\n", runtime.NumGoroutine())
}
