## kn importer update

Update an importer.

### Synopsis

Update an importer.

```
kn importer update NAME [flags]
```

### Options

```
  -h, --help               help for update
      --image string       If you want to create a containersource, image to run.
  -n, --namespace string   List the requested object(s) in given namespace.
      --sink string        Name and type of the sink, ex. 'broker:default'. Defaults to namespace broker when not specified 
      --type string        Type of the source. Currently supported option: 'container'
```

### Options inherited from parent commands

```
      --config string       config file (default is $HOME/.kn/config.yaml)
      --kubeconfig string   kubectl config file (default is $HOME/.kube/config)
```

### SEE ALSO

* [kn importer](kn_importer.md)	 - Importer command group

