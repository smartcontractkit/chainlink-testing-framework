# Troubleshooting

## Can't run `anvil`, issue with `Rosetta`
```
2024/11/27 15:20:27 ⏳ Waiting for container id 79f8a68c07cc image: f4hrenh9it/foundry:latest. Waiting for: &{Port:8546 timeout:0x14000901278 PollInterval:100ms skipInternalCheck:false}
2024/11/27 15:20:27 container logs (all exposed ports, [8546/tcp], were not mapped in 5s: port 8546/tcp is not mapped yet
wait until ready: check target: retries: 1, port: "", last err: container exited with code 133):
rosetta error: Rosetta is only intended to run on Apple Silicon with a macOS host using Virtualization.framework with Rosetta mode enabled
```
#### Solution

Update your docker to `Docker version 27.3.1, build ce12230`
