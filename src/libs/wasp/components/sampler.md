# WASP - Sampler

By default, WASP saves all successful, failed, and timed-out responses in each [Generator](./generator.md).  
The **Sampler** component allows you to programmatically set the sampling ratio for successful responses (a value between 0 and 100).

For example, if you set the sampling ratio to 10, only 10% of successful responses will be saved:

```go
samplerCfg := &SamplerConfig{
    SuccessfulCallResultRecordRatio: 10,
}

g := &Generator{
    sampler: NewSampler(samplerCfg),
    // other fields
}
