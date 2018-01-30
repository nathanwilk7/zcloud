ZCLOUD_PROV=TEST

echo "TESTING: Test Provider"

./zcloud storage cp a.txt cloud://my-bucket/b.txt
if [[ $? == "0" ]]; then
	echo "FAIL: a.txt does not exist but cp did not fail"
	exit 1
fi
./zcloud storage cp asdf://my-bucket/b.txt c.txt
if [[ $? == "0" ]]; then
	echo "FAIL: source url is not valid, but test provider did not fail"
	exit 1
fi
./zcloud storage cp README.md cloud://my-bucket/readme.md
if [[ $? != "0" ]]; then
	echo "FAIL: did not upload file with test provider"
	exit 1
fi
./zcloud storage ls cloud://my-bucket/
if [[ $? != "0" ]]; then
	echo "FAIL: did not list files with test provider"
	exit 1
fi
echo "PASSED end-to-end testing with test provider"
