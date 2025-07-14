package main

import (
	"context"
	"fmt"
	"os"
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
	table := tview.NewTable().SetBorders(true)
	header := []string{"Service", "Status", "Uptime", "Port", "Health"}
	for i, h := range header {
		table.SetCell(0, i, tview.NewTableCell(h).SetSelectable(false).SetAlign(tview.AlignCenter).SetAttributes(tcell.AttrBold))
	}

	updateTable := func() {
		for i, svc := range cfg.Services {
			status, _ := docker.ContainerStatus(ctx, svc.Name)
			healthStatus := "-"
			health := health.NewMonitor(docker, cfg, nil)
			h := health.CheckService(ctx, svc)
			if h.Status != "" {
				healthStatus = string(h.Status)
			}
			// Uptime = Not implemented
			table.SetCell(i+1, 0, tview.NewTableCell(svc.Name))
			table.SetCell(i+1, 1, tview.NewTableCell(status))
			table.SetCell(i+1, 2, tview.NewTableCell("-"))
			table.SetCell(i+1, 3, tview.NewTableCell(fmt.Sprintf("%d", svc.Port)))
			table.SetCell(i+1, 4, tview.NewTableCell(healthStatus))
		}
	}

	updateTable()

	go func() {
		for {
			app.QueueUpdateDraw(updateTable)
			time.Sleep(5 * time.Second)
		}
	}()

	help := tview.NewTextView().SetText("[s] Start  [r] Restart  [l] Logs  [q] Quit").SetTextAlign(tview.AlignCenter)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(help, 1, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := table.GetSelection()
		if row == 0 {
			return event
		}
		selected := cfg.Services[row-1]
		switch event.Rune() {
		case 's':
			go docker.StartContainer(ctx, selected.Name, selected.Image, nil, nil)
		case 'r':
			go docker.RestartContainer(ctx, selected.Name)
		case 'l':
			fmt.Printf("[LOGS] Not implemented: %s\n", selected.Name)
		case 'q':
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
	}
}
