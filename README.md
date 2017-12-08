## ZCloud
ZCloud is intended to prevent coupling to specific cloud providers when using cloud computing services. ZCloud provides a layer of abstraction between the code you write and the code required to use cloud services. You configure which cloud provider you want to use in exactly one place, so if you ever want to switch providers, you only need to change one line of code. In addition, using ZCloud means you don't ever have to learn more than one CLI for interacting with the cloud.

### How to Get

(hopefully, if your go env is set up right)
```bash
go get -u github.com/nathanwilk7/zcloud
```

(otherwise, from src)
```bash
git clone https://github.com/nathanwilk7/zcloud
cd zcloud
go install
```

### Example Usage
First, set the ZCloud env vars for aws:
```bash
export ZCLOUD_PROV=AWS
export ZCLOUD_AWS_KEY_ID=$AWS_ACCESS_KEY_ID
export ZCLOUD_AWS_SECRET_KEY=$AWS_SECRET_ACCESS_KEY
```
Then, copy a file to s3:
```bash
zcloud storage cp <src path> cloud://<bucket-name>/
```
Or, copy a file from s3:
```bash
zcloud storage cp cloud://<bucket-name>/<filename> <dest path>
```
Basically, usage mirrors that of `aws s3 cp` or `gsutil cp`

### Current State
- Provider Interfaces (commands/flags): in progress
- Blob Storage: in progress
- AWS Support: in progress
- gcloud Support: in progress
- Testing: in progress

### Roadmap
- Replace exec calls with something cheaper, probably go's SDK
- Compute in addition to blob storage
- ZCloud SDK's for Go, Python, Java, C#, etc.
- Support for Azure
- Support for OpenStack
- Define uniform interface for all cloud services
- Solve production and distribution of wealth