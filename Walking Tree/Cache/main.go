package main

import (
	"fmt"
	"sync"
)

/*
Design a pseudocode function in Golang that implements a service managing distributed cache with error handling in a microservice architecture
The function should accept a list of operations, each containing a 'cacheKey', 'operationType' (set, get delete) and an optional value.
If an operation fails due to a non. existent key or an unsupported operation type then it should panic. Implement a recovery mechanism that logs the error details while allowing valid cache operations to proceed,
Additionally, ensure that the service can handle concurrent requests data corruption and return a summary indicating how many operations were successfully executed, how many failed and the resons of any failures.

*/

type Operation struct {
	CacheKey string
	Type     string // set, get, delete
	Value    string
}

func ProcessOperations(ops []Operation) Summary {

	cache := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	success := 0
	failed := 0
	failures := []string{}

	for _, op := range ops {

		wg.Add(1)

		go func(op Operation) {
			defer wg.Done()

			defer func() {
				if r := recover(); r != nil {
					failed++
					failures = append(failures,
						fmt.Sprintf("panic: %v", r))
				}
			}()

			mu.Lock()
			defer mu.Unlock()

			switch op.Type {

			case "set":
				cache[op.CacheKey] = op.Value
				success++

			case "get":
				_, ok := cache[op.CacheKey]
				if !ok {
					failed++
					failures = append(failures,
						"key not found: "+op.CacheKey)
					return
				}
				success++

			case "delete":
				delete(cache, op.CacheKey)
				success++

			default:
				failed++
				failures = append(failures,
					"invalid operation: "+op.Type)
			}

		}(op)
	}

	wg.Wait()

	return Summary{
		Success: success,
		Failed:  failed,
		Errors:  failures,
	}
}
