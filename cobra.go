package main

import (
    "fmt"
    "github.com/spf13/cobra"
)

func main() {
    var rootCmd = &cobra.Command{
        Use:   "mycli",
        Short: "A simple CLI tool",
        Long:  "A longer description of my CLI tool.",
    }

    var greetCmd = &cobra.Command{
        Use:   "greet [name]",
        Short: "Greet a user",
        Args:  cobra.MinimumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            name := args[0]
            fmt.Printf("Hello, %s!\n", name)
        },
    }

    rootCmd.AddCommand(greetCmd)
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}