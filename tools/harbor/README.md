# 下载
```bash
wget https://github.com/goharbor/harbor/releases/download/v2.6.0/harbor-offline-installer-v2.6.0.tgz
```

# 对接S3

```bash
storage_service:
  # ca_bundle is the path to the custom root ca certificate, which will be injected into the truststore
  # of registry's and chart repository's containers.  This is usually needed when the user hosts a internal storage with self signed certificate.
  ca_bundle:
  # storage backend, default is filesystem, options include filesystem, azure, gcs, s3, swift and oss
  # for more info about this configuration please refer https://docs.docker.com/registry/configuration/
  s3:
    accesskey: harbor
    secretkey: wjI0PbalIt7PKUP0qtHxIJlZjgLx1712
    region: harbor
    regionendpoint: http://ossapi.rsq.cn:9000
    bucket: harbor
    encrypt: false
    secure: false
    v4auth: true
    chunksize: 5242880
    multipartcopychunksize: 33554432
    multipartcopymaxconcurrency: 100
    multipartcopythresholdsize: 33554432
    rootdirectory: /harbor
  # set disable to true when you want to disable registry redirect
  #redirect:
  #  disabled: true
```