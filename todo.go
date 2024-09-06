package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type TodoItem struct {
	Task        string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type TodoList []TodoItem

func (t *TodoList) Add(task string) {
	item := TodoItem{Task: task, 
		CreatedAt:   time.Now(),
		Completed:   false,
		CompletedAt: time.Time{},
	}
	*t = append(*t, item)

}

func (t *TodoList) MarkAsCompleted(index int) error{
	if index < 0 || index >= len(*t) {
		return fmt.Errorf("invalid index: %d", index)
	}
	(*t)[index-1].Completed = true
	(*t)[index-1].CompletedAt = time.Now()
	return nil
}

func (t *TodoList) Delete(index int) error {
	if index < 0 || index >= len(*t) {
		return fmt.Errorf("invalid index: %d", index)
	}
	*t = append((*t)[:index-1], (*t)[index:]...)
	return nil
}

func (t *TodoList) LoadFromFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("error reading file: %w", err)
	}
	return json.Unmarshal(content, t)
}

func (t *TodoList) SaveToFile(filename string) error {
	content, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("error marshalling todo list: %w", err)
	}
	return os.WriteFile(filename, content, 0644)
}

func (t *TodoList) Print() {
	headerFmt := color.New(color.FgGreen, color.Bold).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Task", "Completed", "Created At", "Completed At")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for i, item := range *t {
		completedStr := "No"
		completedAt := "-"
		if item.Completed {
			completedStr = "Yes"
			completedAt = item.CompletedAt.Format("2006-01-02 15:04:05")
		}
		tbl.AddRow(i+1, item.Task, completedStr, item.CreatedAt.Format("2006-01-02 15:04:05"), completedAt)
	}

	tbl.Print()
}

func (t *TodoList) PrintInTable() {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: color.HiCyanString("ID")},
			{Align: simpletable.AlignCenter, Text: color.HiMagentaString("Task")},
			{Align: simpletable.AlignCenter, Text: color.HiYellowString("Completed")},
			{Align: simpletable.AlignCenter, Text: color.HiGreenString("Created At")},
			{Align: simpletable.AlignCenter, Text: color.HiBlueString("Completed At")},
		},
	}

	completedCount := 0
	for i, item := range *t {
		completedStr := color.RedString("✗")
		completedAt := "-"
		if item.Completed {
			completedStr = color.GreenString("✓")
			completedAt = item.CompletedAt.Format("2006-01-02 15:04:05")
			completedCount++
		}
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", i+1)},
			{Text: item.Task},
			{Align: simpletable.AlignCenter, Text: completedStr},
			{Align: simpletable.AlignRight, Text: item.CreatedAt.Format("2006-01-02 15:04:05")},
			{Align: simpletable.AlignRight, Text: completedAt},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 5, Text: color.HiWhiteString(
				fmt.Sprintf("Total: %d | Completed: %d | Pending: %d",
					len(*t), completedCount, len(*t)-completedCount),
			)},
		},
	}

	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

func (t *TodoItem) MarkAsCompleted() {
	t.Completed = true
	t.CompletedAt = time.Now()
}