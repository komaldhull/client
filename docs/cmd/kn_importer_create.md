## kn importer create

Create an importer.

### Synopsis

Create an importer.

```
kn importer create NAME [flags]
```

### Options

```
      --force              Create source forcefully, replaces existing source if any.
  -h, --help               help for create
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

