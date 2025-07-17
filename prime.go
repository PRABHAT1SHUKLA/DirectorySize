package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 1. Simple CLI Tool - Word Counter
func wordCounter() {
	fmt.Print("Enter text: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	
	words := strings.Fields(text)
	wordCount := make(map[string]int)
	
	for _, word := range words {
		word = strings.ToLower(strings.Trim(word, ".,!?"))
		wordCount[word]++
	}
	
	fmt.Printf("\nWord count:\n")
	for word, count := range wordCount {
		fmt.Printf("%s: %d\n", word, count)
	}
	fmt.Printf("Total words: %d\n", len(words))
}

// 2. REST API Server
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{package main

	import (
		"bufio"
		"context"
		"encoding/csv"
		"encoding/json"
		"fmt"
		"io"
		"log"
		"net/http"
		"os"
		"sort"
		"strconv"
		"strings"
		"sync"
		"time"
	)
	
	// 1. Simple CLI Tool - Word Counter
	func wordCounter() {
		fmt.Print("Enter text: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		
		words := strings.Fields(text)
		wordCount := make(map[string]int)
		
		for _, word := range words {
			word = strings.ToLower(strings.Trim(word, ".,!?"))
			wordCount[word]++
		}
		
		fmt.Printf("\nWord count:\n")
		for word, count := range wordCount {
			fmt.Printf("%s: %d\n", word, count)
		}
		fmt.Printf("Total words: %d\n", len(words))
	}
	
	// 2. REST API Server
	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	
	var users = []User{