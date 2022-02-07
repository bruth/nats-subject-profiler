# NATS Subject Profiler

## Example

Connect it to the demo NATS server.

```sh
nats-subject-profiler -server demo.nats.io:4222
```

You may see some subjects come through. However, if not, you can publish to the server in another shell using the [NATS CLI](https://github.com/nats-io/natscli)

```sh
nats --server demo.nats.io pub foo 'hello'
```

*Note, that subject deduplication is on by default, so if you publish multiple times you will only see the profile print it out once.*
