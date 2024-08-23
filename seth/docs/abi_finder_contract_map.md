# ABI Finder and Contract map

In order to be able to decode and trace transactions and calls between smart contracts we need their ABIs. Unfortunately it might happen that two or more contracts have methods with the same signatures, which might result in incorrect tracing. To make that problem less severe we have decided to add a single point of entry for contract deployment in Seth as that way we always know what contract is deployed at which address and thus avoid incorrect tracing due to potentially ambiguous method signatures.

## ABI Finder

1.  We don’t know what contract (ABI) is located at a given address. Should be the case, when the contract either wasn’t uploaded via Seth or we haven’t supplied Seth with a contract map as part of its configuration (more on that later).

    a. We sequentially iterate over all known ABIs (map: `name -> ABI_name`) checking whether it has a method with a given signature. Once we get a first match we will upsert that (`address -> ABI_name`) data into the contract map and return the ABI.

        The caveat here is that if the method we are searching for belongs is present in more than one ABI we might associate the address with an incorrect address (we will use the first match).

    b. If no match is found we will return an error.

2.  We know what ABI is located at a given address. It should be the case, when we have either uploaded the contract via Seth, provided Seth with a contract map or already traced a transaction to that address and found an ABI with matching method signature.

    a. We fetch the corresponding ABIand check if it indeed contains the method we are looking for (as mentioned earlier in some cases it might not be the case).

    b. If it does, we return the ABI.

    c. If it doesn’t we iterate over all known ABIs, in the same way as in 1a. If we find a match we update the (`address -> ABI_name`) association in the contract map and return the ABI.

        It is possible that this will happen multiple times in case we have multiple contracts with multiple identical methods, but given a sufficiently diverse set of methods that were called we should eventually arrive at a fully correct contract map.

    d. If no match is found we will return an error.

## Contract map

We support in-memory contract map and a TOML file contract map that keeps the association of (`address -> ABI_name`). The latter map is only used for non-simulated networks. Every time we deploy a contract we save (`address -> ABI_name`) entry in the in-memory map.If the network is not a simulated one we also save it in a file. That file can later be pointed to in Seth configuration and we will load the contract map from it (**currently without validating whether we have all the ABIs mentioned in the file**).

When saving contract deployment information we will either generate filename for you (if you didn’t configure Seth to use a particular file) using the pattern of `deployed_contracts_${network_name}_${timestamp}.toml` or use the filename provided in Seth TOML configuration file.

It has to be noted that the file contract map is currently updated only, when new contracts are deployed. There’s no mechanism for updating it if we found the mapping invalid (which might be the case if you manually created the entry in the file).
