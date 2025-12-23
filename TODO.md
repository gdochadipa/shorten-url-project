# System Design of a URL Shortener Overview

---

## 1. Goal of the System
Convert a long URL like:
`https://example.com/articles/life-hacks-2025`

into a short, shareable link like:
`short.ly/abc123`

---

## 2. Core Backend Workflow

### a. User Submits a Long URL
- A long URL is sent to the server via an API.

### b. Generate a Unique Short Code
Two common methods:

- **Counter + Base62**:
  Start from 1, 2, 3... and convert to Base62 (e.g., `a, b, Z, 10`, etc.)

- **Random Base62 String**:
  Generate a 6–8 character string like `abc123`, check database for uniqueness.

| Method             | Pros                               | Cons                                    |
|--------------------|------------------------------------|-----------------------------------------|
| Counter + Base62   | Simple, no duplicates              | Needs centralized counter               |
| Random String      | Easy to scale, looks clean         | Low chance of collision, still needs check |
| Hashing Long URL   | Same input = same code             | Collisions if truncated, not customizable |

### c. Save to Database
Store mapping in a table:

| ShortCode | LongURL                                      | Clicks |
|-----------|----------------------------------------------|--------|
| abc123    | https://example.com/articles/life-hacks-2025 | 0      |

---

## 3. Redirection Flow
When someone opens `short.ly/abc123`:

1. Server extracts the `abc123` code.
2. Looks it up in the database.
3. Increments the click counter (optional).
4. Redirects the user to the long URL instantly.

---

## 4. Scaling the System (How to Handle Millions of Users Smoothly)

- **Caching (e.g., Redis):**
  Store popular short codes in memory for faster lookup. Reduces DB load and speeds up redirection.

- **Database Sharding:**
  Split large databases into smaller chunks (shards) based on short code ranges or user IDs to spread the load.

- **Read Replicas:**
  Use secondary databases for read-heavy operations like redirection, while writes go to the primary DB.

- **Load Balancers:**
  Distribute incoming traffic across multiple servers to avoid bottlenecks.

- **Stateless Backend Servers:**
  Keep servers lightweight and stateless so they can scale horizontally (just add more servers when needed).

- **Asynchronous Processing:**
  For analytics or logging (click tracking), use background jobs to avoid slowing down redirection.

- **CDN Integration (Optional):**
  If static redirection rules become common, use Content Delivery Networks to handle them globally with ultra-low latency.

- **Rate Limiting:**
  Prevent abuse by limiting how many short URLs a user can generate in a certain time.

- **Monitoring & Alerts:**
  Track system health (latency, error rates, traffic spikes) and set up alerts for anomalies.
