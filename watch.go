package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	version = "1.0.0"
	sleep   = 2 * time.Second
)

func main() {
	var command string

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-h", "--help":
			fmt.Println("Usage: watch [command] [-interval N]")
			fmt.Println("Options:")
			fmt.Println("  -interval N    Set refresh rate in seconds (default: 2)")
			fmt.Println("  -h, --help     Show this help message")
			os.Exit(0)
		case "-interval":
			if i+1 < len(os.Args) {
				if sec, err := strconv.Atoi(os.Args[i+1]); err == nil && sec > 0 {
					sleep = time.Duration(sec) * time.Second
					i++
				}
			}
		}
	}

	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stdout, "Welcome to watch %v.\nType command to watch.\n> ", version)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			command = scanner.Text()
			break
		}
	} else {
		command = strings.Join(os.Args[1:], " ")
	}

	shell := defaultShell
	if s, ok := os.LookupEnv("WATCH_COMMAND"); ok {
		shell = s
	}
	sh := strings.Split(shell, " ")
	sh = append(sh, command)

	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = screen.Init()
	if err != nil {
		panic(err)
	}

	app := tview.NewApplication()
	app.SetScreen(screen)
	viewer := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetTextColor(tcell.ColorDefault)
	viewer.SetBackgroundColor(tcell.ColorDefault)

	var statusBar *tview.Flex
	var elapsed *tview.TextView

	elapsed = tview.NewTextView().
		SetTextColor(tcell.ColorLightCyan).
		SetTextAlign(tview.AlignRight).
		SetText("0s")
	elapsed.SetBackgroundColor(tcell.ColorBlack)

	title := tview.NewTextView().
		SetTextColor(tcell.ColorLightCyan).
		SetText(command)
	title.SetBackgroundColor(tcell.ColorBlack)

	statusBar = tview.NewFlex().
		AddItem(title, 0, 1, false).
		AddItem(elapsed, 7, 1, false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(viewer, 0, 1, true)
	flex.AddItem(statusBar, 1, 1, false)
	app.SetRoot(flex, true)

	go func() {
		for {
			cmd := exec.Command(sh[0], sh[1:]...)

			var buf bytes.Buffer
			err := cmdOutput(cmd, &buf)
			if err != nil {
				panic(err)
			}

			app.QueueUpdateDraw(func() {
				screen.Clear()
				viewer.SetText(tview.TranslateANSI(buf.String()))
			})

			startTime := time.Now()
			smoothCheckerDuration := time.Duration(369)

			for {
				remainingTime := sleep - time.Since(startTime)

				if remainingTime <= 0 {
					if elapsed != nil {
						app.QueueUpdateDraw(func() {
							elapsed.SetText(fmt.Sprintf("%.2f s", 0.0))
						})
					}
					break
				}

				if elapsed != nil {
					app.QueueUpdateDraw(func() {
						elapsed.SetText(fmt.Sprintf("%.2f s", remainingTime.Seconds()))
					})
				}

				time.Sleep(smoothCheckerDuration * time.Millisecond)
			}
		}
	}()

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
