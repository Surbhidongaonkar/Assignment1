package main

import (
	"testing"
	"time"
	"fmt"
)

func TESTCOUNTOFBUILDS(t *testing.T) {
	ste, er := time.Parse(time.RFC3339, "2018-11-01T08:54:40-04:00")
	ete, err := time.Parse(time.RFC3339, "2018-11-10T08:54:40-04:00")
	if er!=nil && err != nil {
		fmt.Println(err)
	}
	count:=countofBuilds(ste, ete)
	//fmt.Println("Count of Builds:", count)
	if(count!=41){
		t.Errorf("Wrong count. Count Expected is:41, Got: %d", count)
	}
}