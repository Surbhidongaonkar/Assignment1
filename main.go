package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

//Brecord stores single record 
type Brecord struct {
	Buildid   string
	UID       string
	reqtime   time.Time
	starttime time.Time
	endtime   time.Time
	delete    bool
	exitcode  int
	size      int
}

//UIDPair is struct to store uid as key and its build count as value
type UIDPair struct {
	Key   string
	Value int
}

// UIDPairList is a slice of pairs that implements sort.Interface to sort by values
type UIDPairList []UIDPair

// Buildrecord is array to Store each record as Brecord struct
var Buildrecord [100]Brecord
var countbuilds int //no of builds executed in given time
//var st, et time.Time

func main() {
	// Open the file
	ReadFile()
	st, er := time.Parse(time.RFC3339, os.Args[1])
	et, err := time.Parse(time.RFC3339, os.Args[2])
	if er!=nil && err != nil {
		fmt.Println(err)
	}
	// call to func to get count og build
	count:=countofBuilds(st,et)
	fmt.Println("Total builds executed in given time are: ", count)
	// call to func to get failed builds
	FailedBuilds(st,et)
	//call to func to get frequent users
	topuserlist(st,et)
}

//ReadFile reads data from CSV file and stores records into BuildRecord
func ReadFile(){
	csvfile, err := os.Open("data.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	var j int //for looping
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//copying records one by one
		Buildrecord[j].Buildid = record[0]
		Buildrecord[j].UID = record[1]
		Buildrecord[j].reqtime, err = time.Parse(time.RFC3339, record[2])
		Buildrecord[j].starttime, err = time.Parse(time.RFC3339, record[3])
		Buildrecord[j].endtime, err = time.Parse(time.RFC3339, record[4])
		Buildrecord[j].delete, err = strconv.ParseBool(record[5])
		Buildrecord[j].exitcode, err = strconv.Atoi(record[6])
		Buildrecord[j].size, err = strconv.Atoi(record[7])
		j++
	}
}

//countofBuilds counts no of builds executed in given time
func countofBuilds(st,et time.Time) int {
	countbuilds:= 0
	for i := 0; i < 100; i++ {
		if Buildrecord[i].starttime.After(st) {
			if Buildrecord[i].endtime.Before(et) {
				countbuilds++
			}
		}
	}
	return countbuilds
}

//FailedBuilds finds builds which are failed and their exit code
func FailedBuilds(st,et time.Time) {
	for i := 0; i < 100; i++ {
		if Buildrecord[i].starttime.After(st) {
			if Buildrecord[i].endtime.Before(et) {
				//finds which build fails and related exit code
				if Buildrecord[i].exitcode > 0 {
					fmt.Printf("Build %s failed with exit code %d\n", Buildrecord[i].Buildid, Buildrecord[i].exitcode)
				}
			}
		}
	}
}

func (p UIDPairList) Len() int           { return len(p) }
func (p UIDPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p UIDPairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

//topuserlist finds which 5 users have executed the max builds in given time
func topuserlist(st,et time.Time) {
	countuid := make(map[string]int) //users and no of builds for maximum builds
	for i := 0; i < 100; i++ {
		if Buildrecord[i].starttime.After(st) {
			if Buildrecord[i].endtime.Before(et) {
				if _, ok := countuid[Buildrecord[i].UID]; ok {
					//increment count of build for existing user
					countuid[Buildrecord[i].UID]++
				} else {
					//adding new user of build
					countuid[Buildrecord[i].UID] = 1
				}
			}
		}
	}
	p := make(UIDPairList, len(countuid))

	i := 0
	for k, v := range countuid {
		p[i] = UIDPair{k, v}
		i++
	}
	sort.Sort(p)
	fmt.Println("Users executed the maximum builds are:")
	i = 0
	for _, k := range p {
		fmt.Printf("%s : %d \n", k.Key, k.Value)
		if i == 4 {
			break
		}
		i++
	}
}
