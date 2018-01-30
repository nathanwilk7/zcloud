./zcloud storage cp testdata/test.txt cloud://zcloud-testing/zcloud-test.txt
if [[ $? != "0" ]]; then
	echo "FAIL: upload failed"
	exit 1
fi
./zcloud storage cp cloud://zcloud-testing/zcloud-test.txt testdata/download.txt
if [[ $? != "0" ]]; then
	echo "FAIL: download failed"
	exit 1
fi
cmp testdata/test.txt testdata/download.txt
if [[ $? != "0" ]]; then
	echo "FAIL: uploaded / downloaded files don't match or there was an error comparing them"
	exit 1
fi
./zcloud storage ls cloud://zcloud-testing/ | grep "zcloud-test.txt"
if [[ $? != "0" ]]; then
	echo "FAIL: did not list uploaded file"
	exit 1
fi
./zcloud storage cp -r testdata/ cloud://zcloud-testing/
if [[ $? != "0" ]]; then
	echo "FAIL: did not upload recursively"
	exit 1
fi
