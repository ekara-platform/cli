# Lagoon CLI `parameters` file
___

Some CLI commands accept the flag `--param` to specify a **yaml** file used to pass parameters which will be interpreted at runtime to substitute variables into the environment descriptor.


##File format

The file can contain any valid yaml data.

##File usage

`$ cli create http://path.to.my.project/ --cert ./cert_location --host tcp://192.168.99.100:2376 --api 1.30 --client myClientName --param myParameters.yaml`


##Variables substritution

Example:

```yaml
aaa: "value1"
bbb: 
  ccc: "value2"
ddd: "value3"
```

> Passing this **yaml** content as parameter to the CLI will generate the following map:

```yaml
aaa: "value1"
bbb.ccc: "value2"
ddd: "value3"
```
> All the nested yaml levels will be concatenated to create the keys of the map. The map keys will be used to proceed to the variables substitution into the environment descriptor.


Example of descriptor:


```yaml
nodes:
  managers:
    provider:
      name: aws
      params:
        model: ${aws.size}
        region: ${aws.region}
      envvars:
        HTTP_PROXY: ${my_http_proxy}
        HTTPS_PROXY: ${my_https_proxy}
```

The substitution of ` ${aws.size}, ${aws.region}...` implies that your parameters file must contain :

```yaml
my_http_proxy: http://your_http_proxy_settings
my_https_proxy: http://your_https_proxy_settings

aws:
  region: eu-west-1
  size: t2.medium 
```




