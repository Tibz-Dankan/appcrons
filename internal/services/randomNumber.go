package services

import (
	"math/rand"
	"strconv"
	"strings"
)

type RandomNumber struct{}

func NewRandomNumber() *RandomNumber {
	return &RandomNumber{}
}

func (r *RandomNumber) generateRandomNumbers(count int) string {
	numbers := make([]string, count)
	for i := 0; i < count; i++ {
		numbers[i] = strconv.Itoa(rand.Intn(10))
	}
	return strings.Join(numbers, "")
}

func (r *RandomNumber) D2() string {
	return r.generateRandomNumbers(2)
}

func (r *RandomNumber) D3() string {
	return r.generateRandomNumbers(3)
}

func (r *RandomNumber) D4() string {
	return r.generateRandomNumbers(4)
}

func (r *RandomNumber) D5() string {
	return r.generateRandomNumbers(5)
}

func (r *RandomNumber) D6() string {
	return r.generateRandomNumbers(6)
}

func (r *RandomNumber) D7() string {
	return r.generateRandomNumbers(7)
}

func (r *RandomNumber) D8() string {
	return r.generateRandomNumbers(8)
}

func (r *RandomNumber) D9() string {
	return r.generateRandomNumbers(9)
}

func (r *RandomNumber) D10() string {
	return r.generateRandomNumbers(10)
}
