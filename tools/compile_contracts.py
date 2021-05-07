import json
import subprocess
import os

rootdir = "./artifacts/contracts/ethereum/src"
targetdir = "./contracts/ethereum"

for subdir, dirs, files in os.walk(rootdir):
    for f in files:
        if ".dbg." not in f:
            print(f)
            compile_contract = open(subdir + "/" + f, "r")
            data = json.load(compile_contract)
            contract_name = data["contractName"]

            abi_name = targetdir + "/abi/" + contract_name + ".abi"
            abi_file = open(abi_name, "w")
            abi_file.write(json.dumps(data["abi"], indent=2))

            bin_name = targetdir + "/bin/" + contract_name + ".bin"
            bin_file = open(bin_name, "w")
            bin_file.write(str(data["bytecode"]))
            abi_file.close()
            bin_file.close()

            subprocess.run("abigen --bin=" + bin_name + " --abi=" + abi_name + " --pkg=" + contract_name + " --out=" + 
            targetdir + "/" + contract_name + ".go", shell=True, check=True)