package drive

import (
	"google.golang.org/api/drive/v3"
	"net/http"
)

func CreateDuplicate(srvDrive *drive.Service, fileId, filename, desc string) (fileDupId string, err error) {
	file, err := srvDrive.Files.Copy(fileId, &drive.File{Name: filename, Description: desc}).Do()
	if err != nil {
		return
	}
	switch file.Shared {
	case true:
		fileDupId = file.DriveId
		return
	case false:
		fileDupId = file.Id
		return
	}
	return
}

func DeleteFiles(srvDrive *drive.Service, fileId ...string) (fileStatus map[string]bool, err error) {
	fileStat := make(map[string]bool, len(fileId))
	for _, v := range fileId {
		er := srvDrive.Files.Delete(v).Do()
		if er != nil {
			fileStat[v] = false
			continue
		}
		fileStat[v] = true
	}
	fileStatus = fileStat

	return
}

func DownloadFile(srvDrive *drive.Service, fileId, mimeType string) (res *http.Response, err error) {
	res, err = srvDrive.Files.Export(fileId, mimeType).Download()

	return
}
