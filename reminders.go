package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Name struct {
	Date        int64
	Description string
}

var (
	newFlag   = flag.Bool("n", false, "create new reminder")
	clearFlag = flag.Bool("c", false, "clear reminders")
	listFlag  = flag.Bool("l", false, "list all reminders")
	helpFlag  = flag.Bool("h", false, "help")
)

var (
	path = os.Getenv("GOPATH")
)

func checkFile() {
	if _, err := os.Stat(path + "/reminders.json"); os.IsNotExist(err) {
		os.Create(path + "/reminders.json")
	}
}

func newReminder() {
	checkFile()

	var text string

	file, err := os.OpenFile(path+"/reminders.json", os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Print("Description: ")

	scanner := bufio.NewScanner(os.Stdin)
	if ok := scanner.Scan(); ok {
		text = scanner.Text()
	}

	reminder := &Name{time.Now().Unix(), text}
	newReminder, err := json.Marshal(&reminder)
	if err != nil {
		log.Fatal(err)
	}

	file.Write(newReminder)
}

func removeReminders() {
	err := os.Remove(path + "/reminders.json")
	if err != nil {
		log.Fatal(err)
	}
}

func listReminders() error {
	file, err := os.Open(path + "/reminders.json")
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	for {
		var content Name
		switch decoder.Decode(&content) {
		case nil:
			t := time.Unix(content.Date, 0)
			fmt.Printf("%s: %s\n", t.Format("2006-01-02 15:04:05"), content.Description)
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	if *newFlag {
		newReminder()
	}

	if *clearFlag {
		removeReminders()
	}

	if *helpFlag {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	}

	if *listFlag {
		err := listReminders()
		if err != nil {
			log.Fatal(err)
		}
	}

}
