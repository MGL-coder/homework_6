package homework_6

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

// taskGenerator generates random n functions with error probability p
// t determines the maximum time needed for one task in milliseconds
func taskGenerator(n int, p float32, t int) []func() error {
	funcs := make([]func() error, n)
	for i := range funcs {
		funcs[i] = func() error {
			sleepDur := rand.Int() % t + 1
			time.Sleep(time.Millisecond * time.Duration(sleepDur))
			if rand.Float32() < p {
				return errors.New("error occurred")
			}
			return nil
		}
	}
	return funcs
}

func equalError(e1, e2 error) bool {
	if e1 != nil && e2 != nil {
		return e1.Error() == e2.Error()
	}
	return e1 == e2
}

func TestExecute(t *testing.T) {
	testTable := []struct{
		tasks []func() error
		N int
		E int
		expectedError error
	}{
		{
			tasks:         taskGenerator(10, 0.0, 100),
			N:             3,
			E:             1,
			expectedError: nil,
		},
		{
			tasks:         taskGenerator(10, 0.1, 100),
			N:             3,
			E:             5,
			expectedError: nil,
		},
		{
			tasks:         taskGenerator(10, 0.9, 100),
			N:             3,
			E:             1,
			expectedError: errors.New("exceeded the permissible number of errors"),
		},
		{
			tasks:         taskGenerator(10, 0.9, 100),
			N:             0,
			E:             1,
			expectedError: errors.New("incorrect arguments: N, E, len(tasks) must be positive"),
		},
	}

	for i, testCase := range testTable {
		t.Logf("Executing test case %v", i)
		result := Execute(testCase.tasks, testCase.N, testCase.E)

		if !equalError(result, testCase.expectedError) {
			t.Fatalf("Incorrect result: expected \"%s\", got \"%s\"", result, testCase.expectedError)
		}
	}
}


func benchmarkExecute(n int, b *testing.B) {
	// running Execute with n concurrent tasks
	// each task takes maximum about 10 milliseconds of time
	tasks := taskGenerator(n, 0.0, 10)
	for i := 0; i < b.N; i++ {
		Execute(tasks, n, 1)
	}
}

func BenchmarkExecute1(b *testing.B) { benchmarkExecute(1, b)}
func BenchmarkExecute100(b *testing.B) { benchmarkExecute(100, b)}
func BenchmarkExecute1000(b *testing.B) { benchmarkExecute(1000, b)}
func BenchmarkExecute10000(b *testing.B) { benchmarkExecute(10000, b)}
func BenchmarkExecute1000000(b *testing.B) { benchmarkExecute(1000000, b)}
