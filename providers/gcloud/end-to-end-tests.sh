set -e

ZCLOUD_PROV=GCLOUD

echo "TESTING GCLOUD"

./end-to-end-tests.sh

echo "PASSED end-to-end testing with GCLOUD"
