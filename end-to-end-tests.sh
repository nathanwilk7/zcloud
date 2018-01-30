mkdir -p testdata/dir
echo "natee!" > testdata/test.txt
echo "aaaaaaa" > testdata/dir/a.txt
mkdir testdata/recursive
endtest () {
	if [[ $1 == 0]]; then
		rm -rf testdata
	fi
	exit $1
}
./zcloud storage cp testdata/test.txt cloud://zcloud-testing/zcloud-test.txt
if [[ $? != "0" ]]; then
	echo "FAIL: upload failed"
	endtest 1
fi
./zcloud storage cp cloud://zcloud-testing/zcloud-test.txt testdata/download.txt
if [[ $? != "0" ]]; then
	echo "FAIL: download failed"
	endtest 1
fi
cmp testdata/test.txt testdata/download.txt
if [[ $? != "0" ]]; then
	echo "FAIL: uploaded / downloaded files don't match or there was an error comparing them"
	endtest 1
fi
./zcloud storage ls cloud://zcloud-testing/ | grep "zcloud-test.txt"
if [[ $? != "0" ]]; then
	echo "FAIL: did not list uploaded file"
	endtest 1
fi
./zcloud storage cp -r testdata/ cloud://zcloud-testing/
if [[ $? != "0" ]]; then
	echo "FAIL: did not upload recursively"
	endtest 1
fi
./zcloud storage cp -r cloud://zcloud-testing/ testdata/recursive/
if [[ $? != "0" ]]; then
	echo "FAIL: did not download recursively"
	endtest 1
fi
endtest 0
