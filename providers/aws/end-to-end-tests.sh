set -e

ZCLOUD_PROV=AWS
ZCLOUD_AWS_KEY_ID=$AWS_ACCESS_KEY_ID
ZCLOUD_AWS_SECRET_KEY=$AWS_SECRET_ACCESS_KEY
ZCLOUD_AWS_REGION="us-east-1"

echo "TESTING AWS"

chmod +x providers/end-to-end-zcloud-test.sh
providers/end-to-end-zcloud-test.sh

echo "PASSED end-to-end testing with AWS"
