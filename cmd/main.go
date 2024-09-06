package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	// "time"
	"github.com/arpit/go-awesome-todo-cli"
)

const TODO_FILE = ".todo.json"

func main() {
	add := flag.Bool("add", false, "add a new task")
	complete := flag.Int("complete", 0, "mark a task as completed")
	del := flag.Int("del", 0, "delete a task")
	list := flag.Bool("list", false, "list all tasks")

	flag.Parse()

	todoList := &todo.TodoList{}

	err := todoList.LoadFromFile(TODO_FILE)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading todo list from file:", err.Error())
		os.Exit(1)
	}

	switch {
		case *add:
			input, err := getInput(os.Stdin, flag.Args()...)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error getting input for task:", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Adding a new task - %s", input)
			todoList.Add(input)

			err = todoList.SaveToFile(TODO_FILE)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error saving todo list to file:", err.Error())
				fmt.Println("Error saving todo list to file:", err)
				os.Exit(1)
			}
		case *complete > 0:
			fmt.Println("Marking a task as completed")
			err := todoList.MarkAsCompleted(*complete)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error marking task as completed:", err.Error())
				os.Exit(1)
			}
			err = todoList.SaveToFile(TODO_FILE)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error saving todo list to file:", err.Error())
				os.Exit(1)
			}
		case *del > 0:
			fmt.Println("Deleting a task")
			err := todoList.Delete(*del)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error deleting task:", err.Error())
				os.Exit(1)
			}
			err = todoList.SaveToFile(TODO_FILE)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error saving todo list to file:", err.Error())
				os.Exit(1)
			}
		case *list:
			todoList.PrintInTable()
		default:
			fmt.Fprintln(os.Stderr, "Invalid command provided")
			os.Exit(1)
	}

}


func getInput(reader io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("no task provided")
	}

	return text, nil
}