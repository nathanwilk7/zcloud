set -e -x

export ZCLOUD_PROV=AWS
export ZCLOUD_AWS_KEY_ID=$AWS_ACCESS_KEY_ID
export ZCLOUD_AWS_SECRET_KEY=$AWS_SECRET_ACCESS_KEY
export ZCLOUD_AWS_REGION="us-east-1"
export ZCLOUD_GCLOUD_PROJECT_ID="zcloud-testing"
export ZCLOUD_DEST_PROV=GCLOUD

echo "TESTING AWS"

. end-to-end-test.sh

echo "PASSED end-to-end testing with AWS"

ZCLOUD_PROV=GCLOUD
ZCLOUD_DEST_PROV=AWS

echo "TESTING GCLOUD"

. end-to-end-test.sh

echo "PASSED end-to-end testing with GCLOUD"
