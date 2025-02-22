# practice-gojek-server

## This repo implements followings

### Rate Limiter Problem Statement

Design and implement a rate-limiting service for an API. The service should:

1. Limit the number of API requests a user can make in a fixed time window (e.g., 100 requests per minute).
2. Return an appropriate response (e.g., HTTP 429 - Too Many Requests) when the limit is exceeded.
3. Support multiple users, each with their own independent rate limit tracking.

Key Considerations:

- Used in-memory data structures to track requests.
- Ensured efficiency and scalability.
- Handle edge cases like `burst traffic` (time) and `time window` overlaps.

Extensions:

- Make the time window and request limit configurable.
- this could scale in a distributed system.
