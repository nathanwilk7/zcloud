makeTestFile () {
	mkdir -p testdata/
	echo "natee!" > testdata/test.txt
}
removeFiles () {
	rm -rf testdata
}
endTest () {
	if [[ $1 == 0 ]]; then
		removeFiles
	fi
	exit $1
}
makeTransferPrimary () {
	TEMP=$ZCLOUD_PROV
	ZCLOUD_PROV=$ZCLOUD_DEST_PROV
}
revertTransfer () {
	ZCLOUD_PROV=$TEMP
}
makeTestFile
./zcloud storage mb zcloud-testing
if [[ $? != "0" ]]; then
	echo "FAIL: make bucket zcloud-testing failed"
	endTest 1
fi
makeTransferPrimary()
./zcloud storage mb zcloud-transfer-testing
if [[ $? != "0" ]]; then
	echo "FAIL: make bucket zcloud-transfer-testing failed"
	endTest 1
fi
revertTransfer()
./zcloud storage cp testdata/test.txt cloud://zcloud-testing/zcloud-test.txt
if [[ $? != "0" ]]; then
	echo "FAIL: upload failed"
	endTest 1
fi
./zcloud storage cp cloud://zcloud-testing/zcloud-test.txt testdata/download.txt
if [[ $? != "0" ]]; then
	echo "FAIL: download failed"
	endTest 1
fi
cmp testdata/test.txt testdata/download.txt
if [[ $? != "0" ]]; then
	echo "FAIL: uploaded / downloaded files don't match or there was an error comparing them"
	endTest 1
fi
./zcloud storage ls cloud://zcloud-testing/ | grep "zcloud-test.txt"
if [[ $? != "0" ]]; then
	echo "FAIL: did not list uploaded file"
	endTest 1p
fi
./zcloud storage transfer cloud://zcloud-testing/zcloud-test.txt cloud://zcloud-transfer-testing/zcloud-transfer.txt
if [[ $? != "0" ]]; then
	echo "FAIL: did not transfer"
	endTest 1
fi
makeTransferPrimary()
./zcloud storage ls cloud://zcloud-transfer-testing/ | grep "zcloud-transfer.txt"
if [[ $? != "0" ]]; then
	echo "FAIL: did not list transferred file"
	endTest 1
fi
./zcloud storage rm cloud://zcloud-transfer-testing/zcloud-transfer.txt
if [[ $? != "0" ]]; then
	echo "FAIL: remove transfer failed"
	endTest 1
fi
revertTransfer()
./zcloud storage rm cloud://zcloud-testing/zcloud-test.txt
if [[ $? != "0" ]]; then
	echo "FAIL: remove failed"
	endTest 1
fi
./zcloud storage ls cloud://zcloud-testing/ | grep "zcloud-test.txt"
if [[ $? == "0" ]]; then
	echo "FAIL: list found zcloud-test.txt"
	endTest 1
fi
removeFiles
makeRecursiveFiles() {
	mkdir -p testdata/dir
	mkdir testdata/recursive
	echo "a" > testdata/a.txt
	echo "b" > testdata/dir/b.txt
}
makeRecursiveFiles
./zcloud storage cp -r testdata/ cloud://zcloud-testing/
if [[ $? != "0" ]]; then
	echo "FAIL: did not upload recursively"
	endTest 1
fi
./zcloud storage cp -r cloud://zcloud-testing/ testdata/recursive/
if [[ $? != "0" ]]; then
	echo "FAIL: did not download recursively"
	endTest 1
fi
cmp testdata/a.txt testdata/recursive/a.txt
if [[ $? != "0" ]]; then
	echo "FAIL: recursive upload / download failed on a.txt"
	endTest 1
fi
cmp testdata/dir/b.txt testdata/recursive/dir/b.txt
if [[ $? != "0" ]]; then
	echo "FAIL: recursive upload / download failed on b.txt"
	endTest 1
fi
./zcloud storage ls -r cloud://zcloud-testing/ | grep "a.txt"
if [[ $? != "0" ]]; then
	echo "FAIL: recursive list did not output a.txt"
	endTest 1
fi
./zcloud storage ls -r cloud://zcloud-testing/ | grep "dir/b.txt"
if [[ $? != "0" ]]; then
	echo "FAIL: recursive list did not output b.txt"
	endTest 1
fi
./zcloud storage sync cloud://zcloud-testing/ cloud://zcloud-testing/synced/
if [[ $? != "0" ]]; then
	echo "FAIL: sync failed"
	endTest 1
fi
./zcloud storage ls -r cloud://zcloud-testing/synced/ | grep "a.txt"
if [[ $? != "0" ]]; then
	echo "FAIL: recursive sync list did not output a.txt"
	endTest 1
fi
./zcloud storage ls -r cloud://zcloud-testing/synced/ | grep "dir/b.txt"
if [[ $? != "0" ]]; then
	echo "FAIL: recursive sync list did not output b.txt"
	endTest 1
fi
./zcloud storage rm cloud://zcloud-transfer-testing/synced/a.txt
if [[ $? != "0" ]]; then
	echo "FAIL: remove sync a.txt failed"
	endTest 1
fi
./zcloud storage rm cloud://zcloud-transfer-testing/synced/dir/b.txt
if [[ $? != "0" ]]; then
	echo "FAIL: remove sync dir/b.txt failed"
	endTest 1
fi
./zcloud storage transfer -r cloud://zcloud-testing/ cloud://zcloud-transfer-testing/
if [[ $? != "0" ]]; then
	echo "FAIL: recursive transfer failed"
	endTest 1
fi
makeTransferPrimary()
./zcloud storage ls -r cloud://zcloud-transfer-testing/ | grep "a.txt"
if [[ $? != "0" ]]; then
	echo "FAIL: recursive transfer list did not output a.txt"
	endTest 1
fi
./zcloud storage ls -r cloud://zcloud-transfer-testing/ | grep "dir/b.txt"
if [[ $? != "0" ]]; then
	echo "FAIL: recursive transfer list did not output b.txt"
	endTest 1
fi
./zcloud storage rm cloud://zcloud-transfer-testing/a.txt
if [[ $? != "0" ]]; then
	echo "FAIL: remove transfer a.txt failed"
	endTest 1
fi
./zcloud storage rm cloud://zcloud-transfer-testing/dir/b.txt
if [[ $? != "0" ]]; then
	echo "FAIL: remove transfer dir/b.txt failed"
	endTest 1
fi
./zcloud storage rb zcloud-transfer-testing
if [[ $? != "0" ]]; then
	echo "FAIL: remove bucket zcloud-transfer-testing failed"
	endTest 1
fi
revertTransfer()
./zcloud storage rm -r cloud://zcloud-testing/
if [[ $? != "0" ]]; then
	echo "FAIL: recursive remove"
	endTest 1
fi
./zcloud storage ls -r cloud://zcloud-testing/ | grep "dir/b.txt"
if [[ $? == "0" ]]; then
	echo "FAIL: recursive list found dir/b.txt"
	endTest 1
fi
./zcloud storage rb zcloud-testing
if [[ $? != "0" ]]; then
	echo "FAIL: remove bucket zcloud-testing failed"
	endTest 1
fi
endTest 0
