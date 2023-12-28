package main

import (
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {
	s := NewMemory()
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("Testing %d", i+1)
		index, err := s.Put([]byte(key))
		if err != nil {
			fmt.Println(err)
		}
		data, err := s.Pull(index)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(data))
	}
}
