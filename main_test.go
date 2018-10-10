package main

import (
	"testing"
	"time"
)

func Test_sequental(t *testing.T) {

	assets := []struct {
		IPAddr string
		Index  int
	}{
		{"10.0.0.1", 0},
		{"10.0.0.2", 1},
		{"10.0.0.3", 2},
		{"10.0.0.1", 3},
		{"10.0.0.2", 4},
		{"10.0.0.3", 5},
		{"10.0.0.1", 6},
	}

	list := []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}

	for i := 0; i < len(assets); i++ {

		a := assets[i]
		ip := sequental(i, list)
		if ip != a.IPAddr {
			t.Error("invalid ip:", a.IPAddr)
		}
	}
}

func Test_getMetric(t *testing.T) {

	s := []map[string]time.Duration{}
	s = append(s, map[string]time.Duration{
		"DATA": time.Second * 1,
	})

	s = append(s, map[string]time.Duration{
		"DATA": time.Second * 2,
	})

	max, min, med := getMetric("DATA", s)

	if max != time.Second*1 {
		t.Error("max duration invalid")
	}

	if min != time.Second*2 {
		t.Error("min duration invalid")
	}

	if d, _ := time.ParseDuration("1.5s"); med != d {
		t.Error("med duration invalid")
	}
}

func Test_countMetric(t *testing.T) {

	s := []map[string]time.Duration{}
	s = append(s, map[string]time.Duration{
		"DATA": time.Second * 1,
	})

	s = append(s, map[string]time.Duration{
		"DATA": time.Second * 2,
	})

	s = append(s, map[string]time.Duration{
		"QUIT": time.Second * 2,
	})

	cnt := countMetric("DATA", s)

	if cnt != 2 {
		t.Error("Invalid key count:", cnt)
	}

	cnt = countMetric("QUIT", s)

	if cnt != 1 {
		t.Error("Invalid key count:", cnt)
	}
}

func Test_metricKeys(t *testing.T) {

	s := []map[string]time.Duration{}
	s = append(s, map[string]time.Duration{
		"DATA": time.Second * 1,
	})
	s = append(s, map[string]time.Duration{
		"MAIL": time.Second * 1,
	})
	s = append(s, map[string]time.Duration{
		"QUIT": time.Second * 1,
	})

	keys := metricKeys(s)

	if len(keys) != 3 {
		t.Error("invalid length")
	}
}

func Test_isContain(t *testing.T) {

	list := []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}
	if !isContain("10.0.0.1", list) {
		t.Error("Key not found")
	}

	if isContain("10.0.0.4", list) {
		t.Error("Key found! oha.")
	}

}
