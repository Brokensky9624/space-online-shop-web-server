## Directory Structure:
```
todolist/
├── cmd/
│   ├── add.go
│   ├── list.go
│   ├── remove.go
│   └── root.go
├── go.mod
└── main.go
```

### `cmd/add.go`:
```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
    Use:   "add [task]",
    Short: "Add a task to the to-do list",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        task := args[0]
        fmt.Printf("Added task: %s\n", task)
        // Here you can implement the logic to add the task to your to-do list
    },
}

func init() {
    rootCmd.AddCommand(addCmd)
}
```

### `cmd/list.go`:
```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List all tasks in the to-do list",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Tasks:")
        // Here you can implement the logic to list all tasks in your to-do list
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}
```

### `cmd/remove.go`:
```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
    Use:   "remove [task]",
    Short: "Remove a task from the to-do list",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        task := args[0]
        fmt.Printf("Removed task: %s\n", task)
        // Here you can implement the logic to remove the task from your to-do list
    },
}

func init() {
    rootCmd.AddCommand(removeCmd)
}
```

### `cmd/root.go`:
```go
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "todolist",
    Short: "A CLI tool for managing your to-do list",
    Long: `Todolist is a CLI tool that allows you to manage your to-do list.
You can add, list, and remove tasks.`,
    Run: func(cmd *cobra.Command, args []string) {
        // If no subcommands are provided, print usage
        cmd.Usage()
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    // This is intentionally left blank
}
```

### `main.go`:
```go
package main

import "github.com/yourusername/todolist/cmd"

func main() {
    cmd.Execute()
}
```
