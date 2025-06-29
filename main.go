package main

/*
	This Go program implements a CRUD (Create, Read, Update, Delete) application for managing a library system.
	It supports functionalities such as adding, updating, deleting, and searching for books
 	as well as managing visitors who can rent and return books.
 	The program uses JSON files to persist data across runs, and it provides a command-line interface

	things I learned:
	1. How to use the "encoding/json" package to marshal and unmarshal data.
		Marshal means to convert Go data structures into JSON format.
		Unmarshal means to convert JSON data back into Go data structures.
	2. How to read and write files in Go using the "os" package.
		Reading files is done using os.ReadFile, and writing files is done using os.WriteFile.
	3. How to use the "bufio" package to read input from the console.
		It allows for buffered input reading, which is efficient for console applications.
	4. How to use slices in Go to manage collections of data.
		Slices are dynamic arrays that can grow and shrink in size.
*/

import (
	"bufio"         // "bufio" is used for reading input from the console
	"encoding/json" // "encoding/json" is used for encoding and decoding JSON data
	"fmt"           // "fmt" is used for formatted I/O operations
	"os"            // "os" is used for operating system functionality, like reading and writing files
	"strings"       // "strings" is used for string manipulation, such as trimming spaces and converting to lower case
)

type Book struct {
	ID     int    `json:"id"`     // ID is the unique identifier for each book
	Title  string `json:"title"`  // Title is the title of the book
	Author string `json:"author"` // Author is the author of the book
}
type Visitor struct {
	ID        int    `json:"id"`             // ID is the unique identifier for each visitor
	Name      string `json:"name"`           // Name is the name of the visitor
	RentedIDs []int  `json:"rented_book_id"` // RentedIDs is a slice of book IDs that the visitor has rented
}

var books = make(map[int]Book)       // books is a slice that holds all the books in the library
var nextID = 1                       // nextID is the next available ID for a new book
var dataFile = "books.json"          // dataFile is the name of the file where books data is stored
var visitors = make(map[int]Visitor) // visitors is a slice that holds all the visitors
var nextVisitorID = 1                // nextVisitorID is the next available ID for a new visitor
var visitorsFile = "visitors.json"   // visitorsFile is the name of the file where visitors data is stored

func waitForReturn(scanner *bufio.Scanner) {
	fmt.Print("\npress Enter to return: ")
	scanner.Scan()         // Wait for the user to press Enter
	text := scanner.Text() // Read the input
	if text != "" {        // If the input is not empty, print a message
		waitForReturn(scanner)
	}
}

func loadVisitors() {
	data, err := os.ReadFile(visitorsFile) // Read the visitors file
	if err != nil {                        // If the file does not exist, we start with an empty slice
		fmt.Println("No visitors file found.")
		return
	}
	err = json.Unmarshal(data, &visitors) // Unmarshal the JSON data into the visitors slice
	if err != nil {                       // If there is an error reading the JSON, print an error message
		fmt.Println("Error reading visitors:", err)
		return
	}
	for _, v := range visitors {
		if v.ID >= nextVisitorID {
			nextVisitorID = v.ID + 1
		}
	}
}

func saveVisitors() {
	data, err := json.MarshalIndent(visitors, "", "  ")
	if err != nil {
		fmt.Println("Error saving visitors:", err)
		return
	}
	err = os.WriteFile(visitorsFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing visitors file:", err)
	}
}

func loadBooks() {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		fmt.Println("No data file found, starting fresh.")
		return
	}

	err = json.Unmarshal(data, &books)
	if err != nil {
		fmt.Println("Error reading JSON:", err)
		return
	}
	// Find max ID to set nextID
	nextID = 1
	for id := range books {
		if id >= nextID {
			nextID = id + 1
		}
	}
}

func saveBooks() {
	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		fmt.Println("Error saving books:", err)
		return
	}
	err = os.WriteFile(dataFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}

func createBook(title, author string) {
	book := Book{ID: nextID, Title: title, Author: author}
	books[nextID] = book
	nextID++
	saveBooks()
	fmt.Println("Book created:", book)
}

func searchBooks(query string) {
	query = strings.ToLower(query)
	found := false

	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), query) {
			fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
			found = true
		}
	}

	if !found {
		fmt.Println("No books found matching your search.")
	}
}

func readBooks() {
	if len(books) == 0 {
		fmt.Println("No books found.")
		return
	}
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}
}

func updateBook(id int, newTitle, newAuthor string) {
	book, exists := books[id]
	if !exists {
		fmt.Println("Book not found")
		return
	}
	book.Title = newTitle
	book.Author = newAuthor
	books[id] = book
	saveBooks()
	fmt.Println("Book updated:", book)
}

func deleteBook(id int) {
	if _, exists := books[id]; exists {
		delete(books, id)
		saveBooks()
		fmt.Println("Book deleted:", id)
	} else {
		fmt.Println("Book not found")
	}
}

func showVisitors(scanner *bufio.Scanner) {
	for _, v := range visitors {
		renting := "none"
		if len(v.RentedIDs) > 0 {
			ids := []string{}
			for _, id := range v.RentedIDs {
				ids = append(ids, fmt.Sprintf("%d", id))
			}
			renting = "Book ID(s) " + strings.Join(ids, ", ")
		}
		fmt.Printf("ID: %d, Name: %s, Renting: %s\n", v.ID, v.Name, renting)
	}
	waitForReturn(scanner)
}

func addVisitor(scanner *bufio.Scanner) {
	fmt.Print("Enter visitor name: ")
	scanner.Scan()
	name := scanner.Text()

	visitor := Visitor{ID: nextVisitorID, Name: name}
	visitors[nextVisitorID] = visitor
	nextVisitorID++
	saveVisitors()
	fmt.Println("Visitor added.")
}

func rentBook(scanner *bufio.Scanner) {
	fmt.Print("Visitor ID: ")
	var vid int
	fmt.Scanln(&vid)

	visitor, exists := visitors[vid]
	if !exists {
		fmt.Println("Visitor not found.")
		return
	}

	fmt.Print("Book ID to rent: ")
	var bid int
	fmt.Scanln(&bid)

	var bookExists bool
	for _, book := range books {
		if book.ID == bid {
			bookExists = true
			break
		}
	}
	if !bookExists {
		fmt.Println("Book not found.")
		return
	}

	for _, rid := range visitor.RentedIDs {
		if rid == bid {
			fmt.Println("Visitor already rented this book.")
			return
		}
	}
	visitor.RentedIDs = append(visitor.RentedIDs, bid)

	// Important: Save updated visitor back to map
	visitors[vid] = visitor
	saveVisitors()
	fmt.Println("Book rented.")
}

func returnBook(scanner *bufio.Scanner) {
	fmt.Print("Visitor ID: ")
	var vid int
	fmt.Scanln(&vid)

	visitor, found := visitors[vid]
	if !found {
		fmt.Println("Visitor not found.")
		waitForReturn(scanner)
		return
	}

	fmt.Print("Book ID to return: ")
	var bid int
	fmt.Scanln(&bid)

	index := -1
	for i, id := range visitor.RentedIDs {
		if id == bid {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("This book is not currently rented by the visitor.")
	} else {
		// Remove the book ID from the RentedIDs slice
		visitor.RentedIDs = append(visitor.RentedIDs[:index], visitor.RentedIDs[index+1:]...)
		// Save the updated visitor struct back into the map
		visitors[vid] = visitor
		saveVisitors()
		fmt.Println("Book returned.")
	}

	waitForReturn(scanner)
}

func handleCreate(scanner *bufio.Scanner) {
	fmt.Print("Enter title: ")
	scanner.Scan()
	title := scanner.Text()

	fmt.Print("Enter author: ")
	scanner.Scan()
	author := scanner.Text()

	createBook(title, author)
}

func handleUpdate(scanner *bufio.Scanner) {
	fmt.Print("Enter ID to update: ")
	var id int
	fmt.Scanln(&id)

	fmt.Print("Enter new title: ")
	scanner.Scan()
	newTitle := scanner.Text()

	fmt.Print("Enter new author: ")
	scanner.Scan()
	newAuthor := scanner.Text()
	updateBook(id, newTitle, newAuthor)
}

const Green = "\033[32m"
const Reset = "\033[0m"

func main() {
	loadBooks()
	loadVisitors()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println(Green + "\nAvailable commands: \n\nVisitors Commands\n[VISITORS] [ADDVISITOR] [RENT] \n[RETURN]\n\nBooks Commands\n[CREATE] [READ] [SEARCH] \n[UPDATE] [DELETE] [EXIT]\n" + Reset)
		fmt.Print("Enter command: ")

		if !scanner.Scan() {
			break
		}
		cmd := strings.ToUpper(strings.TrimSpace(scanner.Text()))
		switch cmd {
		case "VISITORS":
			showVisitors(scanner)

		case "ADDVISITOR":
			addVisitor(scanner)

		case "RENT":
			rentBook(scanner)

		case "RETURN":
			returnBook(scanner)

		case "CREATE":
			handleCreate(scanner)

		case "READ":
			readBooks()
			waitForReturn(scanner)

		case "SEARCH":
			fmt.Print("Enter title keyword to search: ")
			scanner.Scan()
			query := scanner.Text()
			searchBooks(query)
			waitForReturn(scanner)

		case "UPDATE":
			handleUpdate(scanner)

		case "DELETE":
			fmt.Print("Enter ID to delete: ")
			var id int
			fmt.Scanln(&id)
			deleteBook(id)

		case "EXIT":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Unknown command.")
		}
	}
}
