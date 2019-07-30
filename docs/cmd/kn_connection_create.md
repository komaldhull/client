## kn connection create

Create a connection.

### Synopsis

Create a connection.

```
kn connection create NAME --image IMAGE [flags]
```

### Options

```
      --async              Create connection and don't wait for it to become ready.
      --broker string      Broker to subscribe to
      --force              Create trigger forcefully, replaces existing trigger if any.
  -h, --help               help for create
  -n, --namespace string   List the requested object(s) in given namespace.
      --sequence string    Sequence that is the subscriber
      --service string     Service that is the subscriber
      --source string      Event source filter
      --type string        Event type filter
      --wait-timeout int   Seconds to wait before giving up on waiting for connection to be ready. (default 60)
```

### Options inherited from parent commands

```
      --config string       config file (default is $HOME/.kn/config.yaml)
      --kubeconfig string   kubectl config file (default is $HOME/.kube/config)
```

### SEE ALSO

* [kn connection](kn_connection.md)	 - Connection command group

