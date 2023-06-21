package main

import (
	"context"
	"fmt"
	gwrp "github.com/JPratama7/gwrap"
	gdocs "github.com/JPratama7/gwrap/docs"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"

	"log"
)

func replaceTextWithTable(docID string, requests []*docs.Request) error {
	srv, err := docs.NewService(context.Background())
	if err != nil {
		return err
	}

	_, err = srv.Documents.BatchUpdate(docID, &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}

	return nil
}

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
	if err != nil {
		log.Fatalf("Unable to retrieve Docs client: %v", err)
		return
	}

	gdrService := gdocs.NewGoogleDocs(srvDocs)
	//gdcsService := gdocs.NewGoogleDocs(srvDocs)
	// Prints the title of the requested doc:
	// https://docs.google.com/document/d/195j9eDD3ccgjQRttHhJPymLJUCOUjs-jmwTrekvdjFE/edit

	loc, err := gdrService.FindTextLocation("ID", "{{CROT}}")

	fmt.Printf("crooot : %s\n", err)
	fmt.Printf("loc : %+v\n", loc)
}
