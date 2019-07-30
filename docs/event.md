## kn Eventing Support

This is a basic prototype implementation of the proposal [here](https://docs.google.com/document/d/1InzYoOG33z8tCkLtd7SqeBWhAkLte148n9tDYcRRHFw/edit)

Includes: 
* broker list/describe commands
* importer CRUD commands (currently supports containersources)
* connection CRUD commands (currently only exposes triggers)

Example workflows: 

Direct importer to service connection: 

```bash
# Create a new importer with sink "mysvc"

kn importer create src --type container  --image docker.io/kdhull/source --sink service:mysvc

# View created importer

kn importer list 

```

Connect service to broker:


```bash
# Create a new importer with sink broker 'default'"

kn importer create src --type container  --image docker.io/kdhull/source --sink broker:default

# View created importer

kn importer list 

# Create a connection between broker default and mysvc 

kn connection create conn --service mysvc	--broker default

# View created connection

kn connection list

```


