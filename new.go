
package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"
	"unsafe"

	_ "github.com/lib/pq" // PostgreSQL driver (example)
)

// 1. Advanced Channel Patterns
func channelPatterns() {
	// Fan-out pattern
	jobs := make(chan int, 10)
	results := make(chan string, 10)
	
	// Start workers
	for w := 1; w <= 3; w++ {
		go func(id int) {
			for job := range jobs {
				result := fmt.Sprintf("Worker %d processed job %d", id, job)
				time.Sleep(100 * time.Millisecond)
				results <- result
			}
		}(w)
	}
	
	// Send jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)
	
	// Collect results
	for r := 1; r <= 5; r++ {
		fmt.Println(<-results)
	}
}

// 2. Select statement with multiple channels
func selectExample() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "Message from channel 1"
	}()
	
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "Message from channel 2"
	}()
	
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		case <-time.After(3 * time.Second):
			fmt.Println("Timeout!")
		}
	}
}

// 3. Rate limiter using ticker
func rateLimiter() {
	requests := make(chan string, 5)
	limiter := time.NewTicker(200 * time.Millisecond)
	defer limiter.Stop()
	
	// Simulate requests
	go func() {
		for i := 1; i <= 10; i++ {
			requests <- fmt.Sprintf("Request %d", i)
		}
		close(requests)
	}()
	
	// Process requests with rate limiting
	for req := range requests {
		<-limiter.C
		fmt.Printf("Processing: %s at %s\n", req, time.Now().Format("15:04:05"))
	}
}

// 4. Worker pool pattern
type Job struct {
	ID   int
	Data string
}

type Result struct {
	JobID  int
	Output string
}

func workerPool() {
	numWorkers := 3
	jobs := make(chan Job, 10)
	results := make(chan Result, 10)
	
	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}
	
	// Send jobs
	for j := 1; j <= 9; j++ {
		jobs <- Job{ID: j, Data: fmt.Sprintf("task-%d", j)}
	}
	close(jobs)
	
	// Collect results
	for r := 1; r <= 9; r++ {
		result := <-results
		fmt.Printf("Job %d result: %s\n", result.JobID, result.Output)
	}
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job.ID)
		time.Sleep(500 * time.Millisecond)
		results <- Result{
			JobID:  job.ID,
			Output: fmt.Sprintf("Processed by worker %d", id),
		}
	}
}

// 5. Generic functions (Go 1.18+)
func genericMax[T comparable](a, b T) T {
	if reflect.ValueOf(a).Kind() == reflect.String {
		if strings.Compare(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b)) > 0 {
			return a
		}
		return b
	}
	// For numeric types, use reflection
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	if av.Kind() >= reflect.Int && av.Kind() <= reflect.Float64 {
		if av.Float() > bv.Float() {
			return a
		}
	}
	return b
}

// 6. Middleware pattern for HTTP
type Middleware func(http.Handler) http.Handler

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer valid-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// 7. Signal handling for graceful shutdown
func gracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		for {
			fmt.Println("Working...")
			time.Sleep(1 * time.Second)
		}
	}()
	
	fmt.Println("Press Ctrl+C to exit")
	<-c
	fmt.Println("\nShutting down gracefully...")
	time.Sleep(2 * time.Second)
	fmt.Println("Goodbye!")
}

// 8. Reflection examples
func reflectionExample() {
	type Person struct {
		Name string `json:"name" db:"full_name"`
		Age  int    `json:"age" db:"age"`
	}
	
	p := Person{Name: "Alice", Age: 30}
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	
	fmt.Printf("Type: %s\n", t.Name())
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fmt.Printf("Field: %s, Type: %s, Value: %v, JSON Tag: %s\n",
			field.Name, field.Type, value.Interface(), field.Tag.Get("json"))
	}
}

// 9. Template example
func templateExample() {
	tmpl := `
Name: {{.Name}}
Age: {{.Age}}
{{if gt .Age 18}}
Status: Adult
{{else}}
Status: Minor
{{end}}
Items:
{{range .Items}}
- {{.}}
{{end}}
`
	
	t := template.Must(template.New("person").Parse(tmpl))
	
	data := struct {
		Name  string
		Age   int
		Items []string
	}{
		Name:  "Bob",
		Age:   25,
		Items: []string{"item1", "item2", "item3"},
	}
	
	var buf bytes.Buffer
	t.Execute(&buf, data)
	fmt.Println(buf.String())
}

// 10. Custom sorting
type Person struct {
	Name string
	Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

func customSort() {
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	
	fmt.Println("Before sorting:", people)
	sort.Sort(ByAge(people))
	fmt.Println("After sorting by age:", people)
	
	// Using sort.Slice (Go 1.8+)
	sort.Slice(people, func(i, j int) bool {
		return people[i].Name < people[j].Name
	})
	fmt.Println("After sorting by name:", people)
}

// 11. Regular expressions
func regexExample() {
	text := "Email: john.doe@example.com, Phone: +1-555-123-4567"
	
	// Email regex
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := emailRegex.FindAllString(text, -1)
	fmt.Println("Found emails:", emails)
	
	// Phone regex
	phoneRegex := regexp.MustCompile(`\+?1?-?\d{3}-?\d{3}-?\d{4}`)
	phones := phoneRegex.FindAllString(text, -1)
	fmt.Println("Found phones:", phones)
	
	// Replace using regex
	result := regexp.MustCompile(`\d+`).ReplaceAllString("abc123def456", "X")
	fmt.Println("Numbers replaced:", result)
}

// 12. File system operations
func fileSystemOps() {
	// Create directory
	err := os.MkdirAll("temp/subdir", 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}
	
	// Walk directory tree
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") {
			fmt.Printf("Go file: %s (size: %d bytes)\n", path, info.Size())
		}
		return nil
	})
	
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
	
	// Clean up
	os.RemoveAll("temp")
}

// 13. URL parsing and building
func urlOperations() {
	rawURL := "https://example.com/path?param1=value1&param2=value2#fragment"
	
	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		return
	}
	
	fmt.Printf("Scheme: %s\n", u.Scheme)
	fmt.Printf("Host: %s\n", u.Host)
	fmt.Printf("Path: %s\n", u.Path)
	fmt.Printf("Query: %s\n", u.RawQuery)
	fmt.Printf("Fragment: %s\n", u.Fragment)
	
	// Parse query parameters
	params := u.Query()
	fmt.Printf("param1: %s\n", params.Get("param1"))
	
	// Build URL
	newURL := &url.URL{
		Scheme: "https",
		Host:   "api.example.com",
		Path:   "/v1/users",
	}
	
	q := newURL.Query()
	q.Add("page", "1")
	q.Add("limit", "10")
	newURL.RawQuery = q.Encode()
	
	fmt.Printf("Built URL: %s\n", newURL.String())
}

// 14. Cryptographic operations
func cryptoExample() {
	data := "Hello, World!"
	
	// MD5 hash
	hash := md5.Sum([]byte(data))
	fmt.Printf("MD5: %x\n", hash)
	
	// Random bytes
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	fmt.Printf("Random bytes: %x\n", randomBytes)
	
	// Base64 encoding
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Printf("Base64 encoded: %s\n", encoded)
	
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	fmt.Printf("Base64 decoded: %s\n", decoded)
}

// 15. Memory and performance monitoring
func memoryStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	fmt.Printf("Memory Stats:\n")
	fmt.Printf("Allocated: %d KB\n", m.Alloc/1024)
	fmt.Printf("Total Allocated: %d KB\n", m.TotalAlloc/1024)
	fmt.Printf("System Memory: %d KB\n", m.Sys/1024)
	fmt.Printf("GC Cycles: %d\n", m.NumGC)
	fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
}

// 16. Database operations (example with PostgreSQL)
func databaseExample() {
	// This is a conceptual example - you'd need actual database
	/*
	db, err := sql.Open("postgres", "user=username dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	// Insert
	_, err = db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", "John", "john@example.com")
	if err != nil {
		log.Fatal(err)
	}
	
	// Query
	rows, err := db.Query("SELECT id, name, email FROM users WHERE age > $1", 18)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var id int
		var name, email string
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("User: %d, %s, %s\n", id, name, email)
	}
	*/
	fmt.Println("Database example (commented out - requires actual database)")
}

// 17. Type assertions and switches
func typeAssertions() {
	var i interface{} = 42
	
	// Type assertion
	if num, ok := i.(int); ok {
		fmt.Printf("i is an int: %d\n", num)
	}
	
	// Type switch
	switch v := i.(type) {
	case int:
		fmt.Printf("int: %d\n", v)
	case string:
		fmt.Printf("string: %s\n", v)
	case bool:
		fmt.Printf("bool: %t\n", v)
	default:
		fmt.Printf("unknown type: %T\n", v)
	}
}

// 18. Unsafe operations (use with caution)
func unsafeExample() {
	s := "hello"
	
	// Get string length using unsafe
	header := (*reflect.StringHeader)(unsafe.Pointer(&s))
	fmt.Printf("String length (unsafe): %d\n", header.Len)
	
	// Convert string to byte slice without copying
	bytes := *(*[]byte)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)}))
	
	fmt.Printf("Bytes: %v\n", bytes)
}

// 19. Benchmarking helper
func benchmarkExample() {
	// Measure function execution time
	start := time.Now()
	
	// Simulate work
	sum := 0
	for i := 0; i < 1000000; i++ {
		sum += i
	}
	
	duration := time.Since(start)
	fmt.Printf("Operation took: %v, Result: %d\n", duration, sum)
}

// 20. Error wrapping (Go 1.13+)
type CustomError struct {
	Op  string
	Err error
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("operation %s failed: %v", e.Op, e.Err)
}

func (e *CustomError) Unwrap() error {
	return e.Err
}

func errorWrapping() {
	err := fmt.Errorf("original error")
	wrappedErr := &CustomError{Op: "database", Err: err}
	
	fmt.Printf("Wrapped error: %v\n", wrappedErr)
	fmt.Printf("Unwrapped error: %v\n", wrappedErr.Unwrap())
}

func main() {
	fmt.Println("=== Advanced Go Code Snippets ===")
	
	// Uncomment to run specific examples:
	// channelPatterns()
	// selectExample()
	// rateLimiter()
	// workerPool()
	// gracefulShutdown()
	
	fmt.Println("\n--- Generic Max Example ---")
	fmt.Println("Max of 5 and 3:", genericMax(5, 3))
	fmt.Println("Max of 'hello' and 'world':", genericMax("hello", "world"))
	
	fmt.Println("\n--- Reflection Example ---")
	reflectionExample()
	
	fmt.Println("\n--- Template Example ---")
	templateExample()
	
	fmt.Println("\n--- Custom Sort Example ---")
	customSort()
	
	fmt.Println("\n--- Regex Example ---")
	regexExample()
	
	fmt.Println("\n--- File System Example ---")
	fileSystemOps()
	
	fmt.Println("\n--- URL Operations Example ---")
	urlOperations()
	
	fmt.Println("\n--- Crypto Example ---")
	cryptoExample()
	
	fmt.Println("\n--- Memory Stats Example ---")
	memoryStats()
	
	fmt.Println("\n--- Database Example ---")
	databaseExample()
	
	fmt.Println("\n--- Type Assertions Example ---")
	typeAssertions()
	
	fmt.Println("\n--- Unsafe Example ---")
	unsafeExample()
	
	fmt.Println("\n--- Benchmark Example ---")
	benchmarkExample()
	
	fmt.Println("\n--- Error Wrapping Example ---")
	errorWrapping()
}
