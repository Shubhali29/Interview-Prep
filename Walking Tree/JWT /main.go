package main

/*
Design a pseudocode function in. Go that accepts a slice of strings representing user generated JSON web tokens (JWT). The function should analyse each token to detect potential vulnerabilities
related to token forgery or manipulation, such as improper signature checks, weak algorithms, missing claims. If any tokens are flagged as potentially insecure,
return a slice of these tokens along with detailed explanation for each flag. Additionally , discuss how you would implement secure token validation
mechanism within a microservice architecture running on K8s, ensuring that it mantains performance and security during high traffic senarioa

*/

func ValidateTokens(tokens []string) []string {

	flagged := []string{}

	for _, token := range tokens {

		claims, err := ParseJWT(token)

		if err != nil {
			flagged = append(flagged, token+" : invalid token")
			continue
		}

		// Verify signature
		if !VerifySignature(token) {
			flagged = append(flagged, token+" : invalid signature")
			continue
		}

		// Verify algorithm
		if claims.Alg == "none" || IsWeakAlgorithm(claims.Alg) {
			flagged = append(flagged, token+" : weak algorithm")
		}

		// Verify expiration
		if claims.Exp < CurrentTime() {
			flagged = append(flagged, token+" : expired token")
		}

		// Verify required claims
		if claims.Iss == "" || claims.Sub == "" {
			flagged = append(flagged, token+" : missing claims")
		}
	}

	return flagged
}

/*
Discuss:



Cloud-Native / Performance Notes

* Cache public keys (JWKS) to avoid repeated key lookups.
* Perform validation concurrently using goroutines for high throughput.
* Use centralized authentication services or API gateways for consistency across microservices.
* Apply rate limiting to prevent token validation abuse.

*/
