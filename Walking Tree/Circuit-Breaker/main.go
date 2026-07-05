package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

func main() {

	cb := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name: "payment-api",

			MaxRequests: 1,

			Interval: 30 * time.Second,

			Timeout: 10 * time.Second,

			ReadyToTrip: func(
				counts gobreaker.Counts,
			) bool {

				return counts.ConsecutiveFailures >= 3
			},
		},
	)

	for i := 0; i < 5; i++ {

		result, err := cb.Execute(
			func() (interface{}, error) {

				return callPaymentAPI()
			},
		)

		fmt.Println(result, err)
	}
}

func callPaymentAPI() (interface{}, error) {

	return nil, errors.New(
		"payment api unavailable",
	)
}
