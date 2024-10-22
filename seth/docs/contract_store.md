# Contract Store

Seth can be used with the contract store, but it would have very limited usage as transaction decoding and tracing cannot work with ABIs. Thus, when you initialise Seth with contract store it is necessary that it can load at least one ABI.

Another use of Contract Store is simplified contract deployment. For that we also need the contract's bytecode. The contract store can be used to store the bytecode of the contract and then deploy it using the `DeployContractFromContractStore(auth *bind.TransactOpts, name string, backend bind.ContractBackend, params ...interface{})` method. When Seth is intialisied with the contract store and no bytecode files (`*.bin`) are provided, it will log a warning, but initialise successfully nonetheless.

If bytecode file wasn't provided you need to use `DeployContract(auth *bind.TransactOpts, name string, abi abi.ABI, bytecode []byte, backend bind.ContractBackend, params ...interface{})` method, which expects you to provide contract name (best if equal to the name of the ABI file), bytecode and the ABI.
