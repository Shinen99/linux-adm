package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func startForking(message string, countOfProcesses, countOfIterations int64) {
	var ret uintptr
	var i, j int64
	var pids []uintptr

	// Расскомментировать, если хочется писать в файл
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	for i = 0; i < countOfProcesses; i++ {

		ret, _, _ = syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
		if ret == 0 {
			// Child process
			_, _, _ = syscall.Syscall(uintptr(syscall.SYS_PRCTL), uintptr(syscall.PR_SET_PDEATHSIG), uintptr(syscall.SIGKILL), 0)
			for j = 0; j < (i+1)*countOfIterations; j++ {
				time.Sleep(1 * time.Second)
				log.Printf("Fork №%d, Message: %s, Iteration: %d\n", i, message, j)
			}
			return
		} else {
			pids = append(pids, ret)
		}
	}

	if ret > 0 {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGCHLD)
		<-ch
		for _, pid := range pids {
			syscall.Syscall(uintptr(syscall.SYS_KILL), pid, uintptr(syscall.SIGKILL), 0)
		}
		os.Exit(0)
	}
}

func main() {
	if len(os.Args) < 4 {
		log.Println("ERROR: Not enough arguments")
		log.Printf("USAGE: %s <message> <count_of_processes> <count_of_iterations>\n", os.Args[0])
		return
	}

	countOfProcesses, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Println("Could not parse second argument")
		return
	}

	countOfIterations, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		log.Println("Could not parse third argument")
		return
	}

	startForking(os.Args[1], countOfProcesses, countOfIterations)

	log.Println("Will never execute")
}
