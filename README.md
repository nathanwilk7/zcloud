## zCloud
zCloud is intended to prevent coupling to specific cloud providers when using cloud computing services. zCloud provides a layer of abstraction between the code you write and the code required to use cloud services. You configure which cloud provider you want to use in exactly one place, so when you want to switch providers, you only need to change one line of code. In addition, zCloud saves you time because you don't have to learn more than one CLI for interacting with the cloud. But wait, there's more! zCloud runs faster than other cloud CLI's because it's natively compiled.

### How to Get
If you haven't already installed Golang, [Download and install Golang](https://golang.org/dl/). Then:
```bash
go get -u github.com/nathanwilk7/zcloud
```

(you can also install from source)
```bash
git clone https://github.com/nathanwilk7/zcloud
cd zcloud
chmod +x install.sh
./install.sh
```

### Example Usage for Amazon Web Services (AWS)
First, set the zCloud env vars for AWS (add these commands to your .bashrc to make them permanent). You'll need an AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY:
```bash
export ZCLOUD_PROV=AWS
export ZCLOUD_AWS_REGION="us-east-1"
export ZCLOUD_AWS_KEY_ID=$AWS_ACCESS_KEY_ID
export ZCLOUD_AWS_SECRET_KEY=$AWS_SECRET_ACCESS_KEY
```
Then, copy a file to S3 (make sure that the bucket you're uploading to exists):
```bash
zcloud storage cp <src path> cloud://<bucket-name>/<filename>
```
Or, copy a file from S3:
```bash
zcloud storage cp cloud://<bucket-name>/<filename> <dest path>
```
Basically, usage mirrors that of `aws s3 cp` and `gsutil cp`

### Example Usage for Google Cloud (gcloud)
Currently, you need to use gcloud's cli authentication. [Here](https://cloud.google.com/sdk/downloads#interactive) are instructions on installation. I recommend the Interactive Installer.

Then, after installing and initializing, authenticate (you may be prompted to install gcloud's beta command group and login via browser, do so):
```bash
gcloud beta auth application-default login
```
After you've authenticated gcloud, set the zCloud provider env var for gcloud (add this command to your .bashrc to make it permanent):
```bash
export ZCLOUD_PROV=GCLOUD
```
Then, copy a file to Google Cloud Storage (make sure that the bucket you're uploading to exists):
```bash
zcloud storage cp <src path> cloud://<bucket-name>/<filename>
```
Or, copy a file from Google Cloud Storage:
```bash
zcloud storage cp cloud://<bucket-name>/<filename> <dest path>
```
### Current State
- Provider Interfaces (commands/flags): in progress
- Blob Storage: in progress
- AWS Support: in progress
- gcloud Support: in progress
- Testing: in progress

### Roadmap
- Add support for OpenStack
- Add support for Azure
- ZCloud SDK's for Go, Python, Java, C#, etc.
- Compute in addition to blob storage
- Define uniform interface for all cloud services
- Solve production and distribution of wealth
