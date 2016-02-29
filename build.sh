VERSION=$(git describe --tags)
gox -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}" -os="windows darwin linux" -ldflags="-s -w -X main.version=${VERSION}"
cd dist
mv  darwin_386_logcat2csv        logcat2csv
zip darwin_386_logcat2csv        logcat2csv     -qm
mv  darwin_amd64_logcat2csv      logcat2csv
zip darwin_amd64_logcat2csv      logcat2csv     -qm
mv  linux_386_logcat2csv         logcat2csv
zip linux_386_logcat2csv         logcat2csv     -qm
mv  linux_amd64_logcat2csv       logcat2csv
zip linux_amd64_logcat2csv       logcat2csv     -qm
mv  linux_arm_logcat2csv         logcat2csv
zip linux_arm_logcat2csv         logcat2csv     -qm
mv  windows_386_logcat2csv.exe   logcat2csv.exe
zip windows_386_logcat2csv       logcat2csv.exe -qm
mv  windows_amd64_logcat2csv.exe logcat2csv.exe
zip windows_amd64_logcat2csv     logcat2csv.exe -qm
cd -
