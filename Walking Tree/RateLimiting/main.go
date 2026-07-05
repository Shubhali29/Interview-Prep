/*
Design a Golang struct called RateLimitingService that implements a token bucket algorithm to control the rate of requests to a microservice. 
The struct should include methods to initialise the rate limiter with a maximum numbero of tokens and refill interval, log incoming requetss and determine if a request can be processed on the available tokens.
Ensure that the implementation is thread safe to handle concurrent requests and discuss how you would optimize the token refill process for high 
availability and performance in a cloud native environment.
*/
package main

type RateLimitingService struct {
    tokens       int
    maxTokens    int
    refillRate   int
    mu           Mutex
}

func NewRateLimiter(maxTokens int, refillRate int) *RateLimitingService {
    rl := &RateLimitingService{
        tokens: maxTokens,
        maxTokens: maxTokens,
        refillRate: refillRate,
    }

    go rl.refillTokens()

    return rl
}

func (rl *RateLimitingService) AllowRequest() bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.tokens > 0 {
        rl.tokens--
        return true
    }

    return false
}

func (rl *RateLimitingService) refillTokens() {
    ticker := NewTicker(1 second)

    for range ticker.C {
        rl.mu.Lock()

        rl.tokens += rl.refillRate

        if rl.tokens > rl.maxTokens {
            rl.tokens = rl.maxTokens
        }

        rl.mu.Unlock()
    }
}


/*
Discuss:

* Use a background goroutine with a ticker to refill tokens periodically.
* Minimize locking overhead using efficient synchronization (e.g., atomic operations).
* Store token state in a shared system like Redis for consistency across multiple instances.
* Ensure rate limits remain accurate when services scale horizontally in Kubernetes.
* Leverage Redis replication/failover to improve availability and fault tolerance.
* Keep the refill process lightweight to reduce latency and improve performance.


*/