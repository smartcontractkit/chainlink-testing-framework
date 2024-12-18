# Concurrency

It's a small library that simplifies dividing N tasks between X workers with but two options:
* payload
* early exit on first error

It was created for parallelising chain interaction with [Seth](../libs/seth.md), where strict ordering
and association of a given private key with a specific contract is required.

> [!NOTE]
> This library is an overkill if all you need to do is to deploy 10 contract, where each deployment
> consists of a single transaction.
> But... if your deployment flow is a multi-stepped one and it's crucial that all operations are executed
> using the same private key (e.g. due to privilleged access) it might be a good pick. Especially, if
> you don't want to extensively test a native `WaitGroup`/`ErrGroup`-based solution.

## No payload
If the task to be executed requires no payload (or it's the same for each task) using the tool is much simpler.

First you need to create an instance of the executor:
```go
l := logging.GetTestLogger(nil)

executor := concurrency.NewConcurrentExecutor[ContractIntstance, contractResult, concurrency.NoTaskType](l)
```

Where generic parameters represent (from left to right):
* type of execution result
* type of channel that holds the results
* type of task payload

In our case, we want the execution to return `ContractInstance`s, that will be stored by this type:
```go
type contractResult struct {
	instance ContractIntstance
}
```

And which won't use any payload, as indicated by a no-op `concurrency.NoTaskType`.

Then, we need to define a function that will be executed for each task. For example:
```go
var deployContractFn = func(channel chan contractResult, errorCh chan error, executorNum int) {
    keyNum := executorNum + 1 // key 0 is the root key

    instance, err := client.deployContractFromKey(keyNum)
    if err != nil {
        errorCh <- err
        return
    }

    channel <- contractResult{instance: instance}
}
```

It needs to have the following signature:
```go
type SimpleTaskProcessorFn[ResultChannelType any] func(resultCh chan ResultChannelType, errorCh chan error, executorNum int)
```
and send results of successful execution to `resultCh` and errors to `errorCh`.

Once the processing function is defined all that's left is the execution:
```go
results, err := executor.ExecuteSimple(client.getConcurrency(), numberOfContracts, deployContractFn)
```

Parameters for `ExecuteSimple` (without payload) are as follows(from left to right):
* concurrency count (number of parallel executors)
* total number of executions
* function to execute

`results` contain a slice with results of each execution with `ContractInstance` type.
`err` will be non-nil if any of the executions failed. To get all errors you should call `executor.GetErrors()`.

## With payload
If your tasks need payload, then two things change.

First, you need to pass task type, when creating the executor instance:
```go
executor := concurrency.NewConcurrentExecutor[ContractIntstance, contractResult, contractConfiguration](l)
```

Here, it's set to dummy:
```go
type contractConfiguration struct{}
```

Second, the signature of processing function:
```go
type TaskProcessorFn[ResultChannelType, TaskType any] func(resultCh chan ResultChannelType, errorCh chan error, executorNum int, payload TaskType)
```
Which now includes a forth parameter representing the payload. And that function's implementation (making use of the payload).

> [!NOTE]
> You can find the full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/concurrency/example_test.go).