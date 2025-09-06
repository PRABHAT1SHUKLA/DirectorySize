package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    // Define flags
    name := flag.String("name", "Guest", "Your name")
    age := flag.Int("age", 0, "Your age")
    verbose := flag.Bool("verbose", false, "Enable verbose output")

    // Parse flags
    flag.Parse()

    // Basic validation
    if *name == "" {
        fmt.Println("Error: name is required")
        flag.Usage()
        os.Exit(1)
    }

    // Output
    if *verbose {
        fmt.Printf("Verbose mode: Hello, %s! You are %d years old.\n", *name, *age)
    } else {
        fmt.Printf("Hello, %s!\n", *name)
    }
}