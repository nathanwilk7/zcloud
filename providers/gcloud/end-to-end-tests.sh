set -e

ZCLOUD_PROV=GCLOUD

echo "TESTING GCLOUD"

chmod +x providers/end-to-end-zcloud-test.sh
providers/end-to-end-zcloud-test.sh

echo "PASSED end-to-end testing with GCLOUD"
