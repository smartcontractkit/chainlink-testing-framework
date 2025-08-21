# WASP - How to configure requests per minute or hour

By default, WASP schedule describes requests **per second**, but what if you want to define the load in terms of requests **per minute** or **per hour**? `RateLimitUnitDuration` to the rescue!

If you want to execute 10 requests per minute, you'd use this generator config:
```go
gen, err := wasp.NewGenerator(&wasp.Config{
    LoadType: wasp.RPS,
    Schedule:   wasp.Plain(10, 2*time.Hour),      // plain line profile - 10 requests per minute for 2h
    Gun:        NewExampleHTTPGun(srv.URL()),
    Labels:     labels,
    LokiConfig: wasp.NewEnvLokiConfig(),
    RateLimitUnitDuration: time.Minute,           // <---- this is the key setting
})
if err != nil {
    panic(err)
}
```

In other words you could say that `RateLimitUnitDuration` represents the *denominator* of rate limit duration and RPS value in the schedule the *numeral*.