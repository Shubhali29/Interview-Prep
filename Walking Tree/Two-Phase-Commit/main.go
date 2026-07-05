/*
Design a pseudocode or Golang Func that implements a service for handling distributed transactions across multiple microservice in a Kubernetes environment.
The function should ensure that all operations within a transaction either complete successfully or rollback in case of failure.
Use goroutines to manage concurrent operations and implement a two-phase commit protocol for ensuring consistency. Handle potential failures and timeouts gracefully,
providing a report that indicates the outcome of each transaction (committed, rolled back, failed).
Include comments to explain your design choices, particularly focusing on fault tolerance, consistency gurantees and concurrency management.

*/

package main

import (
	"context"
	"sync"
	"time"
)

type Service struct {
	Name string
}

func (s Service) Prepare(ctx context.Context) error {
	// Call service's prepare endpoint
	// Return nil if ready
	return nil
}

func (s Service) Commit(ctx context.Context) error {
	// Call service's commit endpoint
	// Persist changes
	return nil
}

func (s Service) Rollback(ctx context.Context) error {
	// Undo prepared changes
	return nil
}

type Result struct {
	Service string
	Status  string
	Err     error
}

func DistributedTransaction(ctx context.Context, services []Service) string {

	prepareResults := make(chan Result, len(services))

	// Phase 1: PREPARE
	for _, svc := range services {
		go func(s Service) {
			ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			err := s.Prepare(ctxTimeout)

			if err != nil {
				prepareResults <- Result{
					Service: s.Name,
					Status:  "failed",
					Err:     err,
				}
				return
			}

			prepareResults <- Result{
				Service: s.Name,
				Status:  "prepared",
			}
		}(svc)
	}

	prepared := []Service{}

	for i := 0; i < len(services); i++ {
		result := <-prepareResults

		if result.Err != nil {

			// Rollback all prepared services
			rollback(prepared)

			return "rolled back"
		}

		prepared = append(prepared,
			findService(services, result.Service))
	}

	// Phase 2: COMMIT
	commitResults := make(chan Result, len(prepared))

	for _, svc := range prepared {
		go func(s Service) {

			ctxTimeout, cancel :=
				context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			err := s.Commit(ctxTimeout)

			commitResults <- Result{
				Service: s.Name,
				Status:  "committed",
				Err:     err,
			}
		}(svc)
	}

	for i := 0; i < len(prepared); i++ {

		result := <-commitResults

		if result.Err != nil {

			rollback(prepared)

			return "failed"
		}
	}

	return "committed"
}

func rollback(services []Service) {
	var wg sync.WaitGroup

	for _, svc := range services {
		wg.Add(1)

		go func(s Service) {
			defer wg.Done()
			s.Rollback()
		}(svc)
	}

	wg.Wait()
}
