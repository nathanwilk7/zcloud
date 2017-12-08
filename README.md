## ZCloud
ZCloud is intended to prevent coupling to specific cloud providers when using cloud computing services. ZCloud provides a layer of abstraction between the code you write and the code required to use cloud services. You configure which cloud provider you want to use in exactly one place, so if you ever want to switch providers, you only need to change one line of code. In addition, using ZCloud means you don't ever have to learn more than one CLI for interacting with the cloud.

### Current State
- Provider Interfaces: in progress
- Blob Storage: in progress
- AWS Support: in progress
- gcloud Support: in progress
- Testing: in progress

### Roadmap
- Replace exec calls with something cheaper, probably go's SDK
- Compute in addition to blob storage
- ZCloud SDK's for Go, Python, Java, C#, Ruby, PHP, C++, etc.
- Support Azure
- Support OpenStack
- Define uniform interface for cloud services
- Solve production and distribution of wealth