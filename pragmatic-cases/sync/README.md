# sync

`Wait()` returns when `Signal` is called.

```
go run main.go
Increment 1000 times
Print10: Waiting
Print10: Finished waiting
Print10: Waiting
Print10: Finished waiting
Print10: Waiting
Print10: Finished waiting
Print10: Waiting
Print10: Finished waiting
Print10: Waiting
Print10: Finished waiting
Print10: Waiting
Print10: Finished waiting
Print10: num reached 1000
Counter.num: 1000, NaiveCounter.num: 980, expected: 1000
```
