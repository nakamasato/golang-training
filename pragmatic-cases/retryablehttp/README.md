# [go-retryablehttp](https://github.com/hashicorp/go-retryablehttp)


Example: with RetryMax 3


```go
import (
	"fmt"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

func main() {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3

	standardClient := retryClient.StandardClient() // *http.Client
	resp, err := standardClient.Get("https://nonexisting-naka.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
```

Run

```
go run main.go
2023/03/15 08:27:59 [DEBUG] GET https://nonexisting-naka.com
2023/03/15 08:27:59 [ERR] GET https://nonexisting-naka.com request failed: Get "https://nonexisting-naka.com": dial tcp: lookup nonexisting-naka.com: no such host
2023/03/15 08:27:59 [DEBUG] GET https://nonexisting-naka.com: retrying in 1s (4 left)
2023/03/15 08:28:00 [ERR] GET https://nonexisting-naka.com request failed: Get "https://nonexisting-naka.com": dial tcp: lookup nonexisting-naka.com: no such host
2023/03/15 08:28:00 [DEBUG] GET https://nonexisting-naka.com: retrying in 2s (3 left)
2023/03/15 08:28:02 [ERR] GET https://nonexisting-naka.com request failed: Get "https://nonexisting-naka.com": dial tcp: lookup nonexisting-naka.com: no such host
2023/03/15 08:28:02 [DEBUG] GET https://nonexisting-naka.com: retrying in 4s (2 left)
2023/03/15 08:28:06 [ERR] GET https://nonexisting-naka.com request failed: Get "https://nonexisting-naka.com": dial tcp: lookup nonexisting-naka.com: no such host
2023/03/15 08:28:06 [DEBUG] GET https://nonexisting-naka.com: retrying in 8s (1 left)
2023/03/15 08:28:14 [ERR] GET https://nonexisting-naka.com request failed: Get "https://nonexisting-naka.com": dial tcp: lookup nonexisting-naka.com: no such host
panic: Get "https://nonexisting-naka.com": GET https://nonexisting-naka.com giving up after 5 attempt(s): Get "https://nonexisting-naka.com": dial tcp: lookup nonexisting-naka.com: no such host
```
