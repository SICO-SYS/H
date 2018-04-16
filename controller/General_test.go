package controller

import (
	"testing"
)

func Test_getActionMap(t *testing.T) {
	value, errcode := getActionMap("aws", "ec2", "DescribeInstances")
	if errcode == 3 {
		return
	}
	if value != "DescribeInstances" || errcode != 0 {
		t.Error("transfer action failed")
	}
}

func Benchmark_getActionMap(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		getActionMap("cloud", "service", "action")
	}
}
