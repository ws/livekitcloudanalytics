# livekitcloudanalytics

[LiveKit](https://livekit.io/) has a [managed cloud offering](https://livekit.io/cloud). That managed cloud offering provides a [basic analytics API](https://docs.livekit.io/cloud/analytics-api/). This project acts as a wrapper around those endpoints (as well as documenting some quirks in the docs). I would bet the LiveKit team rolls out analytics support [in the official SDK](https://github.com/livekit/server-sdk-go) soon, at which point this library will be archived.

Full disclosure, this is my first real Go project- if I'm violating any conventions or missing anything obvious, please don't hesitate to file an issue and let me know! I'm here to learn. I would not use this in production without a quick audit by a Go pro.

## Usage

Basic

```
client := lkanalytics.NewClientWithToken(token)

projectId := "fillthisin" // looks like "p_{gibberish}"

// List sessions
response, err := client.ListSessions(projectId)

// Get a single session
sessionId := "fillthisin" // looks like "RM_{gibberish}"
response, err := client.ListSessionDetails(projectId, sessionId)
```

With a rate limiter
```
import "golang.org/x/time/rate"

limiter := rate.NewLimiter(rate.Every(time.Minute/50), 1)
client := lkanalytics.NewClient().WithToken(token).WithRateLimiter(limiter)
```