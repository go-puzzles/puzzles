package goredis

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestPuzzleRedisClient_Lock(t *testing.T) {
	client := NewRedisClient("localhost:6379", 0)
	ctx := context.Background()
	key := "test_lock"

	// Clean up any existing lock
	client.Del(ctx, key)

	// Test successful lock acquisition
	t.Run("acquire lock success", func(t *testing.T) {
		err := client.TryLock(ctx, key, time.Second)
		assert.NoError(t, err)

		// Verify lock exists
		exists, err := client.Exists(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), exists)

		// Clean up
		err = client.Unlock(ctx, key)
		assert.NoError(t, err)
	})

	// Test lock conflict
	t.Run("lock conflict", func(t *testing.T) {
		// First lock
		err := client.TryLock(ctx, key, time.Second)
		assert.NoError(t, err)

		// Try to acquire same lock
		err = client.TryLock(ctx, key, time.Second)
		assert.ErrorIs(t, err, ErrLockAcquireFailed)

		// Clean up
		err = client.Unlock(ctx, key)
		assert.NoError(t, err)
	})

	// Test lock expiration
	t.Run("lock expiration", func(t *testing.T) {
		err := client.TryLock(ctx, key, 100*time.Millisecond)
		assert.NoError(t, err)

		// Wait for lock to expire
		time.Sleep(200 * time.Millisecond)

		// Should be able to acquire lock again
		err = client.TryLock(ctx, key, time.Second)
		assert.NoError(t, err)

		// Clean up
		err = client.Unlock(ctx, key)
		assert.NoError(t, err)
	})
}

func TestPuzzleRedisClient_Unlock(t *testing.T) {
	client := NewRedisClient("localhost:6379", 0)
	ctx := context.Background()
	key := "test_lock"

	// Clean up any existing lock
	client.Del(ctx, key)

	t.Run("unlock success", func(t *testing.T) {
		// First acquire lock
		err := client.TryLock(ctx, key, time.Second)
		assert.NoError(t, err)

		// Then unlock
		err = client.Unlock(ctx, key)
		assert.NoError(t, err)

		// Verify lock is gone
		exists, err := client.Exists(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), exists)
	})

	t.Run("unlock non-existent lock", func(t *testing.T) {
		err := client.Unlock(ctx, "non_existent_lock")
		assert.ErrorIs(t, err, ErrLockNotFound)
	})
}

func TestPuzzleRedisClient_TryLockWithTimeout(t *testing.T) {
	client := NewRedisClient("localhost:6379", 0)
	ctx := context.Background()
	key := "test_lock"

	// Clean up any existing lock
	client.Del(ctx, key)

	t.Run("acquire with timeout success", func(t *testing.T) {
		err := client.TryLockWithTimeout(ctx, key, time.Second, 500*time.Millisecond)
		assert.NoError(t, err)

		// Clean up
		err = client.Unlock(ctx, key)
		assert.NoError(t, err)
	})

	t.Run("acquire with timeout failure", func(t *testing.T) {
		// First lock
		err := client.TryLock(ctx, key, time.Second)
		assert.NoError(t, err)

		// Try to acquire with short timeout
		err = client.TryLockWithTimeout(ctx, key, time.Second, 200*time.Millisecond)
		assert.ErrorIs(t, err, ErrLockTimeout)

		// Clean up
		err = client.Unlock(ctx, key)
		assert.NoError(t, err)
	})
}

func TestPuzzleRedisClient_ConcurrentLock(t *testing.T) {
	client := NewRedisClient("localhost:6379", 0)
	ctx := context.Background()
	key := "test_concurrent_lock"

	// Clean up any existing lock
	client.Del(ctx, key)

	t.Run("concurrent lock acquisition", func(t *testing.T) {
		numGoroutines := 10
		successCount := int32(0)
		wg := sync.WaitGroup{}
		wg.Add(numGoroutines)

		// Launch multiple goroutines to acquire lock simultaneously
		for i := 0; i < numGoroutines; i++ {
			go func(routineID int) {
				defer wg.Done()

				err := client.TryLock(ctx, key, time.Second)
				if err == nil {
					atomic.AddInt32(&successCount, 1)
					// Simulate some work
					time.Sleep(100 * time.Millisecond)
					// Release the lock
					err = client.Unlock(ctx, key)
					assert.NoError(t, err)
				}
			}(i)
		}

		wg.Wait()
		// Only one goroutine should succeed
		assert.Equal(t, int32(1), successCount)
	})

	t.Run("concurrent lock with timeout", func(t *testing.T) {
		// First acquire the lock to ensure other goroutines will timeout
		err := client.TryLock(ctx, key, 2*time.Second)
		assert.NoError(t, err)
		defer client.Unlock(ctx, key)

		numGoroutines := 5
		timeout := 500 * time.Millisecond
		successCount := int32(0)
		timeoutCount := int32(0)
		wg := sync.WaitGroup{}
		wg.Add(numGoroutines)

		start := time.Now()

		for i := 0; i < numGoroutines; i++ {
			go func(routineID int) {
				defer wg.Done()

				err := client.TryLockWithTimeout(ctx, key, time.Second, timeout)
				if err == nil {
					atomic.AddInt32(&successCount, 1)
					// Simulate some work
					time.Sleep(50 * time.Millisecond)
					// Release the lock
					err = client.Unlock(ctx, key)
					assert.NoError(t, err)
				} else if errors.Is(err, ErrLockTimeout) {
					atomic.AddInt32(&timeoutCount, 1)
				}
			}(i)
		}

		wg.Wait()
		duration := time.Since(start)

		// Verify results
		assert.Equal(t, int32(0), successCount, "No goroutine should acquire the lock")
		assert.Equal(t, int32(numGoroutines), timeoutCount, "All goroutines should timeout")
		assert.True(t, duration >= timeout, "Test should take at least the timeout duration")
	})

	t.Run("concurrent lock and unlock", func(t *testing.T) {
		numIterations := 5
		successCount := int32(0)
		wg := sync.WaitGroup{}
		wg.Add(numIterations)

		// Use channel to control concurrency
		lockChan := make(chan int, numIterations)
		// Initialize tasks
		for i := 0; i < numIterations; i++ {
			lockChan <- i
		}

		// Start multiple workers
		numWorkers := 3
		for i := 0; i < numWorkers; i++ {
			go func(workerID int) {
				for taskID := range lockChan {
					err := client.TryLock(ctx, key, time.Second)
					if err == nil {
						atomic.AddInt32(&successCount, 1)
						// Simulate some work
						time.Sleep(10 * time.Millisecond)
						// Release the lock
						err = client.Unlock(ctx, key)
						if err != nil {
							t.Logf("Unlock error in worker %d, task %d: %v", workerID, taskID, err)
						}
						wg.Done() // Decrease counter only when lock is acquired successfully
					} else {
						t.Logf("Lock error in worker %d, task %d: %v", workerID, taskID, err)
						// Put the task back to queue if lock acquisition fails
						lockChan <- taskID
					}
					// Small delay between attempts
					time.Sleep(20 * time.Millisecond)
				}
			}(i)
		}

		// Wait for all tasks to complete
		wg.Wait()
		close(lockChan)

		// Verify the total number of completed tasks
		assert.Equal(t, int32(numIterations), successCount,
			"Total successful locks should equal number of iterations")
	})
}

type TestStruct struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Tags    []string `json:"tags"`
	IsAdmin bool     `json:"is_admin"`
}

func TestPuzzleRedisClient_SetValue_GetValue(t *testing.T) {
	client := NewRedisClient("localhost:6379", 0)
	ctx := context.Background()

	tests := []struct {
		name     string
		key      string
		value    interface{}
		result   interface{}
		expected interface{}
	}{
		{
			name:     "String Value",
			key:      "test:string",
			value:    "hello world",
			result:   new(string),
			expected: "hello world",
		},
		{
			name:     "Integer Value",
			key:      "test:int",
			value:    42,
			result:   new(int),
			expected: 42,
		},
		{
			name:     "Float Value",
			key:      "test:float",
			value:    3.14,
			result:   new(float64),
			expected: 3.14,
		},
		{
			name:     "Boolean Value",
			key:      "test:bool",
			value:    true,
			result:   new(bool),
			expected: true,
		},
		{
			name: "Struct Value",
			key:  "test:struct",
			value: TestStruct{
				Name:    "John Doe",
				Age:     30,
				Tags:    []string{"admin", "user"},
				IsAdmin: true,
			},
			result: &TestStruct{},
			expected: TestStruct{
				Name:    "John Doe",
				Age:     30,
				Tags:    []string{"admin", "user"},
				IsAdmin: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test SetValue
			err := client.SetValue(ctx, tt.key, tt.value, time.Minute)
			assert.NoError(t, err)

			// Test GetValue
			err = client.GetValue(ctx, tt.key, tt.result)
			assert.NoError(t, err)

			// Compare results
			switch v := tt.result.(type) {
			case *string:
				assert.Equal(t, tt.expected, *v)
			case *int:
				assert.Equal(t, tt.expected, *v)
			case *float64:
				assert.Equal(t, tt.expected, *v)
			case *bool:
				assert.Equal(t, tt.expected, *v)
			case *TestStruct:
				assert.Equal(t, tt.expected, *v)
			}

			// Clean up
			client.Del(ctx, tt.key)
		})
	}
}

func TestPuzzleRedisClient_GetValue_TypeConversionErrors(t *testing.T) {
	client := NewRedisClient("localhost:6379", 0)
	ctx := context.Background()

	tests := []struct {
		name        string
		value       string
		result      interface{}
		expectedErr string
	}{
		{
			name:        "Invalid Int Conversion",
			value:       "not a number",
			result:      new(int),
			expectedErr: "strconv.Atoi: parsing \"not a number\": invalid syntax",
		},
		{
			name:        "Invalid Float Conversion",
			value:       "not a float",
			result:      new(float64),
			expectedErr: "strconv.ParseFloat: parsing \"not a float\": invalid syntax",
		},
		{
			name:        "Invalid Bool Conversion",
			value:       "not a bool",
			result:      new(bool),
			expectedErr: "strconv.ParseBool: parsing \"not a bool\": invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "test:conversion"

			// Set the test value
			err := client.Client.Set(ctx, key, tt.value, time.Minute).Err()
			assert.NoError(t, err)

			// Try to get with wrong type
			err = client.GetValue(ctx, key, tt.result)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)

			// Clean up
			client.Del(ctx, key)
		})
	}
}

func TestPuzzleRedisClient_GetValue_NotFound(t *testing.T) {
	client := NewRedisClient("localhost:6379", 0)
	ctx := context.Background()

	var result string
	err := client.GetValue(ctx, "non:existent:key", &result)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

func TestPuzzleRedisClient_SetValue_InvalidJSON(t *testing.T) {
	client := NewRedisClient("localhost:6379", 0)
	ctx := context.Background()

	// Create a struct with a channel which cannot be JSON marshaled
	invalidStruct := struct {
		Ch chan int
	}{
		Ch: make(chan int),
	}

	err := client.SetValue(ctx, "test:invalid", invalidStruct, time.Minute)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json marshal failed")
}

func TestPuzzleRedisClient_GetValue_InvalidJSON(t *testing.T) {
	client := NewRedisClient("localhost:6379", 0)
	ctx := context.Background()

	// Set invalid JSON string
	err := client.Client.Set(ctx, "test:invalid:json", "{invalid json}", time.Minute).Err()
	assert.NoError(t, err)

	var result TestStruct
	err = client.GetValue(ctx, "test:invalid:json", &result)
	assert.Error(t, err)

	// Clean up
	client.Del(ctx, "test:invalid:json")
}

func TestPuzzleRedisClient_SetGetValue_Bytes(t *testing.T) {
	ctx := context.Background()
	client := setupTestClient(t)

	tests := []struct {
		name    string
		key     string
		value   []byte
		wantErr bool
	}{
		{
			name:    "normal binary data",
			key:     "test_binary",
			value:   []byte{0x00, 0x01, 0x02, 0x03},
			wantErr: false,
		},
		{
			name:    "empty binary data",
			key:     "test_empty_binary",
			value:   []byte{},
			wantErr: false,
		},
		{
			name:    "text as binary",
			key:     "test_text_binary",
			value:   []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name:    "binary with null bytes",
			key:     "test_null_binary",
			value:   []byte{0x00, 0xFF, 0x00, 0xFF},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test SetValue
			err := client.SetValue(ctx, tt.key, tt.value, time.Minute)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Test GetValue
			var got []byte
			err = client.GetValue(ctx, tt.key, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !bytes.Equal(got, tt.value) {
				t.Errorf("GetValue() got = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestPuzzleRedisClient_SetGetValue_Mixed(t *testing.T) {
	ctx := context.Background()
	client := setupTestClient(t)

	t.Run("set string get bytes", func(t *testing.T) {
		key := "test_str_bytes"
		value := "Hello, World!"

		// Set as string
		err := client.SetValue(ctx, key, value, time.Minute)
		if err != nil {
			t.Fatalf("SetValue() error = %v", err)
		}

		// Get as bytes
		var got []byte
		err = client.GetValue(ctx, key, &got)
		if err != nil {
			t.Fatalf("GetValue() error = %v", err)
		}

		if !bytes.Equal(got, []byte(value)) {
			t.Errorf("GetValue() got = %v, want %v", got, []byte(value))
		}
	})

	t.Run("set bytes get string", func(t *testing.T) {
		key := "test_bytes_str"
		value := []byte("Hello, World!")

		// Set as bytes
		err := client.SetValue(ctx, key, value, time.Minute)
		if err != nil {
			t.Fatalf("SetValue() error = %v", err)
		}

		// Get as string
		var got string
		err = client.GetValue(ctx, key, &got)
		if err != nil {
			t.Fatalf("GetValue() error = %v", err)
		}

		if got != string(value) {
			t.Errorf("GetValue() got = %v, want %v", got, string(value))
		}
	})
}

func setupTestClient(t *testing.T) *PuzzleRedisClient {
	// 使用测试环境的Redis配置
	client := NewRedisClient("localhost:6379", 0)

	// 确保Redis连接正常
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}
