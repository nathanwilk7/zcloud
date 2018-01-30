set -e

ZCLOUD_PROV=TEST go test ./...

chmod +x build.sh
./build.sh

chmod +x providers/aws/end-to-end-tests.sh
./providers/aws/end-to-end-tests.sh

chmod +x providers/gcloud/end-to-end-tests.sh
./providers/gcloud/end-to-end-tests.sh
