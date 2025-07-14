package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"dock_pilot/internal/config"
	"dock_pilot/pkg/health"
	"dock_pilot/pkg/services"
)

type ServiceRow struct {
	Name   string
	Status string
	Uptime string
	Port   string
	Health string
}

func RunTUI(ctx context.Context, cfg *config.Config, docker *services.DockerService) {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 0)

	table.SetBorder(true).SetTitle(" 🚀 DockPilot - Serviços ").SetTitleAlign(tview.AlignLeft)

	header := []string{"Service", "Status", "Uptime", "Port", "Health", "Último Log"}
	for i, h := range header {
		table.SetCell(0, i, tview.NewTableCell(h).
			SetSelectable(false).
			SetAlign(tview.AlignCenter).
			SetAttributes(tcell.AttrBold).
			SetTextColor(tcell.ColorYellow))
	}

	updateTable := func() {
		for i, svc := range cfg.Services {
			status, _ := docker.ContainerStatus(ctx, svc.Name)
			healthStatus := "-"
			healthColor := tcell.ColorWhite
			health := health.NewMonitor(docker, cfg, nil)
			h := health.CheckService(ctx, svc)
			if h.Status != "" {
				healthStatus = string(h.Status)
				switch h.Status {
				case "healthy":
					healthColor = tcell.ColorGreen
				case "degraded":
					healthColor = tcell.ColorOrange
				case "unreachable":
					healthColor = tcell.ColorRed
				}
			}
			statusColor := tcell.ColorWhite
			switch status {
			case "running":
				statusColor = tcell.ColorGreen
			case "exited":
				statusColor = tcell.ColorRed
			case "created":
				statusColor = tcell.ColorBlue
			}
			// Lê a última linha do log
			logPath := fmt.Sprintf("./logs/%s.log", svc.Name)
			lastLog := "-"
			if data, err := os.ReadFile(logPath); err == nil {
				lines := strings.Split(string(data), "\n")
				if len(lines) > 1 {
					lastLog = lines[len(lines)-2] // penúltima linha, pois a última é vazia
				}
			}
			table.SetCell(i+1, 0, tview.NewTableCell(svc.Name))
			table.SetCell(i+1, 1, tview.NewTableCell(status).SetTextColor(statusColor))
			table.SetCell(i+1, 2, tview.NewTableCell("-").SetTextColor(tcell.ColorGray))
			table.SetCell(i+1, 3, tview.NewTableCell(fmt.Sprintf("%d", svc.Port)).SetTextColor(tcell.ColorAqua))
			table.SetCell(i+1, 4, tview.NewTableCell(healthStatus).SetTextColor(healthColor))
			table.SetCell(i+1, 5, tview.NewTableCell(lastLog).SetTextColor(tcell.ColorGray))
		}
	}

	updateTable()

	go func() {
		for {
			app.QueueUpdateDraw(updateTable)
			time.Sleep(5 * time.Second)
		}
	}()

	help := tview.NewTextView().
		SetText("[s] Start  [r] Restart  [l] Logs  [q] Quit").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetTextColor(tcell.ColorLightCyan)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(help, 1, 1, false)

	app.SetFocus(table)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := table.GetSelection()
		if row == 0 {
			return event
		}
		selected := cfg.Services[row-1]
		switch event.Rune() {
		case 's':
			go func() {
				if err := docker.StartContainer(ctx, selected.Name, selected.Image, selected.Ports, nil); err != nil {
					fmt.Printf("[ERROR] failed to start container: %v\n", err)
				}
			}()
		case 'r':
			go func() {
				if err := docker.RestartContainer(ctx, selected.Name); err != nil {
					fmt.Printf("[ERROR] failed to restart container: %v\n", err)
				}
			}()
		case 'l':
			go func() {
				logPath := fmt.Sprintf("./logs/%s.log", selected.Name)
				data, err := os.ReadFile(logPath)
				lines := []string{}
				if err == nil {
					lines = strings.Split(string(data), "\n")
				}

				tableLog := tview.NewTable().
					SetSelectable(true, false).
					SetFixed(1, 0)

				tableLog.SetBorder(true).SetTitle(fmt.Sprintf("Logs - %s (q para sair)", selected.Name)).SetTitleAlign(tview.AlignLeft)

				// Header
				tableLog.SetCell(0, 0, tview.NewTableCell("Timestamp").SetAttributes(tcell.AttrBold).SetTextColor(tcell.ColorYellow))
				tableLog.SetCell(0, 1, tview.NewTableCell("Mensagem").SetAttributes(tcell.AttrBold).SetTextColor(tcell.ColorYellow))
				// Preencher linhas
				for i, line := range lines {
					if line == "" {
						continue
					}
					// Separar timestamp e mensagem, se possível
					ts, msg := line, ""
					if idx := strings.Index(line, "] "); idx != -1 {
						ts = line[:idx+1]
						msg = line[idx+2:]
					}
					tableLog.SetCell(i+1, 0, tview.NewTableCell(ts).SetTextColor(tcell.ColorGray))
					tableLog.SetCell(i+1, 1, tview.NewTableCell(msg))
				}
				tableLog.SetFixed(1, 0)
				tableLog.SetSelectable(true, false)
				tableLog.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					if event.Rune() == 'q' || event.Key() == tcell.KeyEsc {
						app.SetRoot(flex, true)
						app.SetFocus(table)
						return nil
					}
					return event
				})
				app.QueueUpdateDraw(func() {
					app.SetRoot(tableLog, true)
					app.SetFocus(tableLog)
				})
			}()
		case 'q':
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
	}
}
