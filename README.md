# Google-Wrapper
Simple Package that wrap google-api to make it easy to use.

currently feature :
 - Duplicate Files
 - Upload File and auto-detect mimetype if not given
 - Get Location from text
 - Get URI from file id drive
 - Delete files
 - Find and Replace Text
 - Create and Download PDF
 - Create Config from credential json file

```sh
go get -u all
go mod tidy
git tag                                 #check current version
git tag v0.0.3                          #set tag version
git push origin --tags                  #push tag version to repo
go list -m github.com/aiteung/gwrap@v0.0.3   #publish to pkg dev, replace ORG/URL with your repo URL
```