package main

import( "fmt"
        "os"
        "flag"
)

type RuneScanner interface{
    ReadRune() (r rune, size int, err error)
    UnreadRune() error
}

type Lexer struct{
    file string
    line int
    r RuneScanner
    lastrune rune
}




func parseArguments() string{

	filenamePtr := flag.String("file", "", "filename to read")
	flag.Parse()
	if *filenamePtr == ""{
			fmt.Println("Error: at least argument -file is necessary.")
			os.Exit(1)
	}

	return *filenamePtr
}


func main(){

    filename := parseArguments()
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: No such file or directory.")
		return
	}
	defer file.Close()

    //scanner := bufio.NewScanner(file)

}
