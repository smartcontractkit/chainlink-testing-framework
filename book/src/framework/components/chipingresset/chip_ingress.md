# Chip Ingress Set

Chip Ingress Set is a composite component that collects Beholder events. It is a thin `testcontainers-go` wrapper over a Docker Compose file copied from the [Atlas](https://github.com/smartcontractkit/atlas/blob/master/chip-ingress/docker-compose.yml) repo.

It consists of 3 components:
- Chip Ingress
- Red Panda
- Red Panda Console

## Configuration

To add it to your stack use following TOML:
```toml
[chip_ingress]
  # using a local docker-compose file
  compose_file='file://../../components/chip_ingress_set/docker-compose.yml'
  # using a remote file
  # compose_file='https://my.awesome.resource.io/docker-compose.yml'
  extra_docker_networks = ["my-existing-network"]
```

Where compose file indicates the location of the `docker-compose.yml` file (remote URLs are supported) and `extra_docker_networks` an optional slice of existing Docker networks, to which whole stack should be connected to.

## Exposed ports

These 3 components expose a variety of ports, but the most important ones from the point of view of user interaction are:
- schema registry port: `18081`
- Kafka port: `19092`
- Red Panda console port: `8080`

## Useful helper methods

Packge contains also a bunch of helper functions tha can:
- create and delete Kafka topics
- fetch `.proto` files from remote repositories and register them with Red Panda


### Topic management
```go
import chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"

topicsErr := chipingressset.DeleteAllTopics(cmd.Context(), redPandaKafkaURLFlag)
if topicsErr != nil {
    panic(topicsErr)
}

createTopicsErr := chipingressset.CreateTopics(ctx, out.RedPanda.KafkaExternalURL, []string{"cre"})
if createTopicsErr != nil {
    panic(createTopicsErr)
}
```

### Protobuf schema registration
```go
out, outErr := chipingressset.New(in.ChipIngress)
if outErr != nil {
    panic(outErr)
}

ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

protoErr := chipingressset.DefaultRegisterAndFetchProtos(
    ctx,
    nil, // GH client will be created dynamically, if needed
    []chipingressset.RepoConfiguration{
    {
        URI:   "https://github.com/smartcontractkit/chainlink-protostractkit",
        Ref:     "626c42d55bdcb36dffe0077fff58abba40acc3e5",
        Folders: []string{"workflows"},
    },
}, out.RedPanda.SchemaRegistryExternalURL)
if protoErr != nil {
    panic(protoErr)
}
```

Since `ProtoSchemaSet` has TOML tags you can also read it from a TOML file with this content:
```toml
[[proto_schema_set]]
# reading from remote registry (only github.com supported)
uri = 'https://github.com/smartcontractkit/chainlink-protos'
ref = '626c42d55bdcb36dffe0077fff58abba40acc3e5'
folders = ['workflows']
subject_prefix = 'cre-'

[[proto_schema_set]]
# reading from local folder
uri = 'file://../../chainlink-protos'
# ref is not supported, when reading from local folders
folders = ['workflows']
subject_prefix = 'cre-'
```

And then use this Go code to register them:
```go
var protoSchemaSets []chipingressset.ProtoSchemaSet
for _, schemaSet := range configFiles {
    file, fileErr := os.ReadFile(schemaSet)
    if fileErr != nil {
        return errors.Wrapf(fileErr, "failed to read proto schema set config file: %s", schemaSet)
    }

    type protoSchemaSets struct {
        Sets []chipingressset.ProtoSchemaSet `toml:"proto_schema_set"`
    }

    var sets protoSchemaSets
    if err := toml.Unmarshal(file, &sets); err != nil {
        return errors.Wrapf(err, "failed to unmarshal proto config file: %s", protoConfig)
    }

    protoSchemaSets = append(reposConfigs, sets.Sets...)
}
```

Registration logic is very simple and should handle cases of protos that import other protos as long they are all available in the `ProtoSchemaSet`s provided to the registration function. That function uses an algorithm called "topological sorting by trail", which will try to register all protos in a loop until it cannot register any more protos or it has registered all of them. That allows us to skip dependency parsing completely.

Kafka doesn't have any automatic discoverability mechanism for subject - schema relationship (it has to be provided out-of-band). Currenly, we create the subject in the following way: <subject_prefix>.<package>.<1st-message-name>. Subject prefix is optional and if it's not present, then subject is equal to: <package>.<1st-message-name>. Only the first message in the `.proto` file is ever registered.

## Protobuf caching

Once fetched from `https://github.com` protobuf files will be saved in `.local/share/beholder/protobufs/<OWNER>/<REPOSTIORY>/<SHA>` folder and subsequently used. If saving to cache or reading from it fails, we will load files from the original source.