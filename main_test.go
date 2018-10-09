package main

import (
	"fmt"
	"testing"
)

func Test_iplist(t *testing.T) {

	list, err := ipv4list()

	if err != nil {
		t.Error(err)
	}

	t.Log(list)
}

func Test_sequental(t *testing.T) {

	list, err := ipv4list()

	if err != nil {
		t.Error(err)
	}

	total := 20

	for i := 0; i < total; i++ {

		ob := sequental(i, list)
		fmt.Println(i, ob)
	}
}
