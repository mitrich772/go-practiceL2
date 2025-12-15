package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var (
	currentCmd *exec.Cmd
	cmdMu      sync.Mutex
)

func setupSignalHandler() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		for range sigCh {
			cmdMu.Lock()
			if currentCmd != nil && currentCmd.Process != nil {
				_ = currentCmd.Process.Signal(syscall.SIGINT)
			}
			cmdMu.Unlock()
		}
	}()
}

func runInteractive() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		commands := parseLineToCommands(line)
		runCommands(commands)
	}
}

func runScript(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "не удалось открыть файл: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		commands := parseLineToCommands(line)
		runCommands(commands)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка чтения: %v\n", err)
	}
}

func parseLineToCommands(line string) [][]string {
	var result [][]string
	commands := strings.Split(line, "|")
	for _, v := range commands {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		command := strings.Fields(v)
		result = append(result, command)
	}
	return result
}

func runCommand(name string, args []string, in io.Reader, out io.Writer) int {
	switch name {

	case "cd":
		path := ""
		if len(args) > 0 {
			path = args[0]
		} else {
			path = os.Getenv("HOME")
			if path == "" {
				path = os.Getenv("USERPROFILE")
			}
		}

		if err := os.Chdir(path); err != nil {
			fmt.Fprintf(out, "cd: %v\n", err)
			return 1
		}
		return 0

	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(out, "pwd: %v\n", err)
			return 1
		}
		fmt.Fprintln(out, dir)
		return 0

	case "echo":
		fmt.Fprintln(out, strings.Join(args, " "))
		return 0

	case "kill":
		if len(args) == 0 {
			fmt.Fprintln(out, "kill: не указан PID")
			return 1
		}
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(out, "kill: неверный PID %s\n", args[0])
			return 1
		}

		if runtime.GOOS == "windows" {
			cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F")
			cmd.Stdout = out
			cmd.Stderr = out
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(out, "kill: %v\n", err)
				return 1
			}
			return 0
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			fmt.Fprintf(out, "kill: %v\n", err)
			return 1
		}
		if err := proc.Signal(syscall.SIGTERM); err != nil {
			fmt.Fprintf(out, "kill: %v\n", err)
			return 1
		}
		return 0

	case "upper":
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			fmt.Fprintln(out, strings.ToUpper(scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(out, "upper: %v\n", err)
			return 1
		}
		return 0

	case "ps":
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("tasklist")
		} else {
			cmd = exec.Command("ps", "aux")
		}
		cmd.Stdout = out
		cmd.Stderr = out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(out, "ps: %v\n", err)
			return 1
		}
		return 0

	case "exit":
		os.Exit(0)

	default:
		cmd := exec.Command(name, args...)
		cmd.Stdin = in
		cmd.Stdout = out
		cmd.Stderr = os.Stderr

		cmdMu.Lock()
		currentCmd = cmd
		cmdMu.Unlock()

		err := cmd.Run()

		cmdMu.Lock()
		currentCmd = nil
		cmdMu.Unlock()

		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if exitErr.ExitCode() == -1 {
					return 130
				}
			}
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
	}
	return 0
}

func runCommands(commands [][]string) int {
	if len(commands) == 0 {
		return 0
	}
	if len(commands) == 1 {
		return runCommand(commands[0][0], commands[0][1:], os.Stdin, os.Stdout)
	}

	readers := make([]*io.PipeReader, len(commands)-1)
	writers := make([]*io.PipeWriter, len(commands)-1)

	for i := 0; i < len(commands)-1; i++ {
		readers[i], writers[i] = io.Pipe()
	}

	var wg sync.WaitGroup
	exitCodes := make(chan int, len(commands))

	for i := range commands {
		wg.Add(1)

		go func(i int, cmd []string) {
			defer wg.Done()

			var in io.Reader = os.Stdin
			var out io.Writer = os.Stdout

			if i > 0 {
				in = readers[i-1]
			}
			if i < len(commands)-1 {
				out = writers[i]
			}

			code := runCommand(cmd[0], cmd[1:], in, out)

			if i < len(commands)-1 {
				_ = writers[i].Close()
			}

			exitCodes <- code
		}(i, commands[i])
	}

	wg.Wait()
	close(exitCodes)

	var lastCode int
	for code := range exitCodes {
		lastCode = code
	}

	return lastCode
}

func main() {
	setupSignalHandler()

	if len(os.Args) > 1 {
		runScript(os.Args[1])
	} else {
		runInteractive()
	}
}
