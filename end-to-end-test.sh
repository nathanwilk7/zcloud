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
makeTestFile
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
	endTest 1
fi
removeFiles
# makeRecursiveFiles() {
# 	mkdir -p testdata/dir
# 	mkdir testdata/recursive
# 	echo "a" > testdata/a.txt
# 	echo "b" > testdata/dir/b.txt
# }
# makeRecursiveFiles
# ./zcloud storage cp -r testdata/ cloud://zcloud-testing/
# if [[ $? != "0" ]]; then
# 	echo "FAIL: did not upload recursively"
# 	endTest 1
# fi
# ./zcloud storage cp -r cloud://zcloud-testing/ testdata/recursive/
# if [[ $? != "0" ]]; then
# 	echo "FAIL: did not download recursively"
# 	endTest 1
# fi
# cmp testdata/a.txt testdata/recursive/a.txt
# if [[ $? != "0" ]]; then
# 	echo "FAIL: recursive upload / download failed on a.txt"
# 	endTest 1
# fi
# cmp testdata/dir/b.txt testdata/recursive/dir/b.txt
# if [[ $? != "0" ]]; then
# 	echo "FAIL: recursive upload / download failed on b.txt"
# 	endTest 1
# fi
# ./zcloud storage ls -r cloud://zcloud-testing/ | grep "a.txt"
# if [[ $? != "0" ]]; then
# 	echo "FAIL: recursive list did not output a.txt"
# 	endTest 1
# fi
# ./zcloud storage ls -r cloud://zcloud-testing/ | grep "dir/b.txt"
# if [[ $? != "0" ]]; then
# 	echo "FAIL: recursive list did not output b.txt"
# 	endTest 1
# fi
endTest 0
