# Self-Serve EC2 VM Template

You should have been authorized in AWS.

Install [CDK](https://docs.aws.amazon.com/cdk/v2/guide/getting_started.html) and connect it to your AWS
```
npm install -g aws-cdk
cdk bootstrap -c action=bootstrap
```

Deploy VM, `connect.sh` will be generated after VM is provisioned
```
cdk deploy -c action=deploy
./connect.sh
```

Destroy VM
```
cdk destroy -c action=destroy
```
Remove `cdk-ec2-keypair.pem` manually

## Forwarding ports
Install `chisel` on your machine
```
curl https://i.jpillora.com/chisel! | bash
```

Start forwarding on the server
```
/usr/local/bin/chisel server --port 44044 &
```

Copy the fingerprint and connect from your machine
```
export CHISEL_FINGERPRINT=
./forward.sh
```
