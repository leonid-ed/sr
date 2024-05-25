package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"sync"
	"time"

	"github.com/cheynewallace/tabby"
	"github.com/logrusorgru/aurora/v3"
	"github.com/xeonx/timeago"
)

var (
	done = make(chan struct{})

	flagNoColors bool
	flagReverse  bool
	flagJson     bool
	flagDigital  bool
	flagFuzzy    bool
	flagUTC      bool
	flagAll      bool
)

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

type fileMsg struct {
	parentID     int
	filename     string
	isDir        bool
	lastModified time.Time
}

// walkDir recursively walks the file tree rooted at dir
// and sends the info of each found file on fileMsg.
func walkDir(currentDir string, parentDir string, maxLevel int, parentID int, wg *sync.WaitGroup, fileMsgs chan<- fileMsg) {
	defer wg.Done()

	if cancelled() {
		return
	}

	numParents := 0
	for _, entry := range handleDir(currentDir) {
		newParentID := parentID
		entryName := entry.Name()

		if strings.HasPrefix(entryName, ".") && !flagAll {
			continue
		}

		msg := fileMsg{
			parentID:     newParentID,
			filename:     filepath.Join(parentDir, entryName),
			isDir:        entry.IsDir(),
			lastModified: entry.ModTime()}

		if parentID == 0 {
			msg.filename = entryName
			fileMsgs <- msg

			numParents++
			if !entry.IsDir() {
				continue
			}
			newParentID = numParents
		}

		if entry.IsDir() {
			newCurrentDir := filepath.Join(currentDir, entryName)
			newParentDir := filepath.Join(parentDir, entryName)
			if maxLevel != 0 {
				wg.Add(1)
				walkDir(newCurrentDir, newParentDir, maxLevel-1, newParentID, wg, fileMsgs)
			}
		} else {
			fileMsgs <- msg
		}
	}
}

// sema is a counting semaphore for limiting concurrency in handleDir.
var sema = make(chan struct{}, 4)

// handleDir returns the entries of directory dir.
func handleDir(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}: // acquire token
	case <-done:
		return nil // cancelled
	}

	defer func() { <-sema }() // release token

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("cannot read the dir '%s': %v\n", dir, err)
	}
	fileInfos := make([]fs.FileInfo, 0, len(dirEntries))
	for _, entry := range dirEntries {
		info, err := entry.Info()
		if err != nil {
			log.Fatalf("cannot read the FileInfo entry '%s': %v\n", entry.Name(), err)
		}
		fileInfos = append(fileInfos, info)
	}

	return fileInfos
}

type Record struct {
	Name               string    `json:"name"`
	IsDir              bool      `json:"isDir"`
	ChildName          string    `json:"child"`
	LastModified       time.Time `json:"-"`
	LastModifiedString string    `json:"lastModified"`
	NumChildren        uint      `json:"numChildren"`
}

func main() {
	// Parse flags.
	var flagMaxLevel int
	flag.BoolVar(&flagNoColors, "n", false, "turn colors off")
	flag.BoolVar(&flagReverse, "r", false, "reverse the order of items")
	flag.BoolVar(&flagJson, "j", false, "show results in json format")
	flag.BoolVar(&flagDigital, "d", false, "show dates in digital format")
	flag.BoolVar(&flagFuzzy, "f", false, "show dates in fuzzy format (e.g. '12 minutes ago')")
	flag.BoolVar(&flagUTC, "u", false, "show time in UTC")
	flag.BoolVar(&flagAll, "a", false, "show all files and directories (including hidden ones)")
	flag.IntVar(&flagMaxLevel, "L", -1, "the max depth of the directory tree; -1 if no depth limit")
	flag.Parse()

	if flagMaxLevel < -1 {
		flagMaxLevel = -1
	}

	roots := flag.Args()
	var rootDir string
	if len(roots) == 0 {
		rootDir = "."
	} else {
		rootDir = roots[0]
	}

	// Cancel traversal when input is detected.
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(done)
		fmt.Println("Cancelled!")
	}()

	// Traverse the file tree in parallel.
	fileMsgs := make(chan fileMsg)
	var wg sync.WaitGroup
	wg.Add(1)
	go walkDir(rootDir, ".", flagMaxLevel, 0, &wg, fileMsgs)

	go func() {
		wg.Wait()
		close(fileMsgs)
	}()

	records := make(map[int]Record)
	currentParentID := 0
loop:
	for {
		select {
		case <-done:
			// Drain fileMsgs to allow existing goroutines to finish.
			for range fileMsgs {
				// Do nothing.
			}
		case msg, ok := <-fileMsgs:
			if !ok {
				break loop // fileMsgs was closed.
			}

			if msg.parentID == 0 {
				currentParentID++
				record := Record{
					Name:         msg.filename,
					IsDir:        msg.isDir,
					ChildName:    msg.filename,
					LastModified: msg.lastModified}

				records[currentParentID] = record
			} else {
				el, ok := records[msg.parentID]
				if !ok {
					log.Fatalf("unknown parentID (%d) (filename: %v, records: %v)\n", msg.parentID, msg.filename, records)
				}
				el.NumChildren++
				if el.NumChildren == 1 || el.LastModified.Before(msg.lastModified) {
					el.ChildName = msg.filename
					el.LastModified = msg.lastModified
				}
				records[msg.parentID] = el
			}
		}
	}
	printResults(records)
}

func printResults(records map[int]Record) {
	if len(records) == 0 {
		return
	}

	dateFormat := time.RFC822Z
	if flagDigital {
		dateFormat = "2006-01-02T15:04:05-0700"
	}

	// Setup shorter forms of timeago values.
	timeagoFormat := timeago.English
	timeagoFormat.DefaultLayout = dateFormat
	timeagoFormat.Periods[0].One = "a second"
	timeagoFormat.Periods[1].One = "a minute"
	timeagoFormat.Periods[2].One = "an hour"
	timeagoFormat.Max = 24 * 365 * 2 * time.Hour

	sliceRecords := make([]Record, 0, len(records))
	for _, v := range records {
		modifiedAt := v.LastModified
		if flagUTC {
			modifiedAt = modifiedAt.UTC()
		}

		if flagFuzzy {
			v.LastModifiedString = timeagoFormat.Format(modifiedAt)
		} else {
			v.LastModifiedString = modifiedAt.Format(dateFormat)
		}

		sliceRecords = append(sliceRecords, v)
	}

	sort.Slice(sliceRecords, func(i, j int) bool {
		if flagReverse {
			i, j = j, i
		}
		return sliceRecords[i].LastModified.After(sliceRecords[j].LastModified)
	})

	au := aurora.NewAurora(!flagNoColors)

	if flagJson {
		jsonString, err := json.Marshal(sliceRecords)
		if err != nil {
			log.Fatalf("cannot marshall: %v\n", err)
		}
		fmt.Printf("%s\n", jsonString)
	} else {
		t := tabby.New()
		var outputName, pathSeparator string
		var auName aurora.Value
		for i := range sliceRecords {
			outputName = fmt.Sprintf("%s%%s", sliceRecords[i].Name)
			if sliceRecords[i].IsDir {
				pathSeparator = string(os.PathSeparator)
				auName = au.Blue(outputName)
			} else {
				pathSeparator = ""
				auName = au.White(outputName)
			}
			outputName = au.Sprintf(auName, au.White(pathSeparator))

			t.AddLine(
				outputName,
				sliceRecords[i].LastModifiedString,
				sliceRecords[i].NumChildren,
				sliceRecords[i].ChildName)
		}
		t.Print()
	}
}
