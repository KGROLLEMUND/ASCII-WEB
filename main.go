package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

//ReadFile returns an array of string which is the same as the file (line = line)
func ReadFile(StylizedFile string) []string {
	var source []string
	file, _ := os.Open(StylizedFile)  // opens the .txt
	scanner := bufio.NewScanner(file) // scanner scans the file
	scanner.Split(bufio.ScanLines)    // sets-up scanner preference to read the file line-by-line
	for scanner.Scan() {              // loop that performs a line-by-line scan on each new iteration
		if scanner.Text() != "" {
			source = append(source, scanner.Text()) // adds the value of scanner (that contains the characters from StylizedFile) to source
		}
	}
	file.Close() // closes the file
	return source
}

//PrintASCII prints the stylized characters
func PrintASCII(arguments []rune, txtfile []string, w http.ResponseWriter) {
	for ligne := 0; ligne < 8; ligne++ { // Each character is composed of 8 lines
		for index, char := range arguments {
			fmt.Fprintf(w, txtfile[ligne+(int(char)-32)*8]) // Recovers the stylized characters from the files and displays them on the website
			if index == len(arguments)-1 && ligne != 7 {    // Jumps a newline when it is required
				fmt.Fprintln(w)
				break
			}
		}
	}
}

type Page struct { //Creates the Page structure which allows the personalization of the website "/ascii-art"
	ColorTxt string //value used to add color selection
	FontSize string //value used to add a size to the Ascii Generator
	ColorBG  string //value used to add color to the background of the website "/ascii-art"
}

//BuildTemplate edit the configuration of the ascii-art page
func BuildTemplate(w http.ResponseWriter, r *http.Request, ColorTxt string, FontSize string, ColorBG string) {
	p := Page{ColorTxt, FontSize, ColorBG}                               // Associates the elements of ascii-art.html with main.go
	parsedTemplate, _ := template.ParseFiles("templates/ascii-art.html") // States the files that requires a modification

	// Executes the modification
	err := parsedTemplate.Execute(w, p)
	// The variable p will be depicted as "." inside the layout
	// Exemple : {{.}} == p

	if err != nil { // If the program contains an error, it will display one
		log.Fatalf("Template execution: %s", err)
	}
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" { // Checks if we are in  /ascii-art
		http.Error(w, "404 not found.", http.StatusNotFound) // Sends a 404 error (page not found)
		return
	}

	if err := r.ParseForm(); err != nil { //If an error occurs during the POST request
		fmt.Fprintf(w, "Can't get input data")
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	// Saves the value from the form tag of the root page ("index.html")
	FontSize := r.FormValue("FontSize")
	ColorBG := r.FormValue("ColorBG")
	ColorTxt := r.FormValue("ColorList")
	Fontlist := r.FormValue("Fontlist")
	Text := r.FormValue("Text")

	if Text == "" { // Tells to the user that he didn't follow the instructions and quit
		fmt.Fprintf(w, "	Please enter the text in the right position then hit \"go run main.go\" button")
		return
	}

	arguments := []rune(Text)
	for index := range arguments {
		if arguments[index] < 32 || arguments[index] > 126 { // Conditions to check if there are non-printable characters
			if arguments[index] != '\n' {
				if arguments[index] != '\r' {
					fmt.Fprintf(w, "	You wrote a non-printable character... Please, try again")
					return
				}
			}
		}
	}

	BuildTemplate(w, r, ColorTxt, FontSize, ColorBG) // Calls the configuration of the webpage

	txtfile := ReadFile(Fontlist) // Recovers the right txt file

	var start int
	w.Write([]byte("<pre>")) // Uses "pre"(html tag) to handle breakline
	for index := range arguments {
		if arguments[index] == '\r' {
			PrintASCII(arguments[start:index], txtfile, w) // Calls the PrintASCII function
			fmt.Fprintln(w)
			start = index + 2
		} else if index == len(arguments)-1 {
			PrintASCII(arguments[start:], txtfile, w)
		}
	}
	w.Write([]byte("</pre>"))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./templates"))) // handle root (main) webpage
	http.HandleFunc("/ascii-art", asciiHandler)                // handle ascii-art page

	fmt.Printf("Starting server at port 8080\n")
	fmt.Println("Go on http://127.0.0.1:8080") // Prints the link of the website on the command prompt
	fmt.Printf("\nTo shutdown the server and exit the code hit \"crtl+C\"\n")
	if err := http.ListenAndServe(":8080", nil); err != nil { // Launches the server on port 8080 if port 8080 is not already busy, else quit
		log.Fatal(err)
	}
}
