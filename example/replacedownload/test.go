package main

import (
	"context"
	"fmt"
	gwrp "github.com/JPratama7/gwrap"
	gdocs "github.com/JPratama7/gwrap/docs"
	gdrive "github.com/JPratama7/gwrap/drive"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"

	"io"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	filepath := "./credentials.json"
	// If modifying these scopes, delete your previously saved token.json.
	//b, err := os.ReadFile(filepath)
	//if err != nil {
	//	log.Fatalf("Unable to read client secret file: %v", err)
	//}

	// Parse the client secret file and configure the OAuth2 client
	//config, err := google.ConfigFromJSON(b, drive.DriveScope, drive.DriveReadonlyScope, docs.DocumentsScope, docs.DocumentsReadonlyScope)
	//if err != nil {
	//	log.Fatalf("Unable to parse client secret file to config: %v", err)
	//}
	cfg, err := gwrp.NewGoogleConfig(filepath, drive.DriveScope, drive.DriveReadonlyScope, docs.DocumentsScope, docs.DocumentsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v\n", err)
		return
	}
	client := gwrp.GetClient(cfg, "token.json")

	srvDocs, err := docs.NewService(ctx, option.WithHTTPClient(client))
	srvDrive, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Docs client: %v", err)
		return
	}

	gdrService := gdrive.NewGoogleDrive(srvDrive)
	gdcsService := gdocs.NewGoogleDocs(srvDocs)
	// Prints the title of the requested doc:
	// https://docs.google.com/document/d/195j9eDD3ccgjQRttHhJPymLJUCOUjs-jmwTrekvdjFE/edit
	docId := "ID"

	docDup, err := gdrService.CreateDuplicate(docId, "DUP 1 MANTAP BHANK", "TESTING DUPLICATE")
	if err != nil {
		log.Fatalf("Unable to create duplicate: %v\n", err)
		return
	}
	fmt.Printf("Duplicate ID : %s\n", docDup)

	file := srvDocs.Documents.Get(docDup)
	if file == nil {
		log.Fatalf("File Not found")
		return
	}
	doc, err := file.Do()

	listReplace := make([]*docs.Request, 0, 4)
	req1 := gdocs.ReplaceTextDocs("{{NAMA}}", "CROOTT")
	req2 := gdocs.ReplaceTextDocs("{{NAMA}}", "CROOTT")
	listReplace = append(listReplace, req1, req2)
	listReplace = append(listReplace, gdocs.ReplaceTextDocs("{{TTL}}", "12-12-2022"))
	listReplace = append(listReplace, gdocs.ReplaceTextDocs("{{TTL}}", "12-12-2022"))

	err = gdcsService.FindAndReplace(docDup, listReplace...)
	if err != nil {
		log.Fatalf("Unable to find and replace: %v", err)
		return
	}

	res, err := gdrService.DownloadFile(docDup, "application/pdf")
	if err != nil {
		log.Fatalf("Unable to retrieve data from document: %v", err)
	}
	fmt.Printf("The title of the doc is: %s\n", doc.Title)

	out, err := os.Create("croott.pdf")
	if err != nil {
		log.Fatalf("Unable to create file: %v", err)
		return
	}
	defer out.Close()
	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Fatalf("Error Copying: %v\n", err)
		return
	}
}
