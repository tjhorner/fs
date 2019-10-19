package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/tjhorner/fs/shorty"

	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
)

var globalConfig *config

func oops(err interface{}) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	globalConfig = loadConfig(defaultConfigPath())
	gcsClient, err := storage.NewClient(context.Background(), option.WithCredentialsFile(defaultServiceAccountPath()))
	if err != nil {
		oops(err)
	}

	bucket := gcsClient.Bucket(globalConfig.BucketName)

	var name string
	var shorten bool
	var verbose bool
	var del string
	flag.BoolVar(&shorten, "shorten", false, "Shorten the resulting URL with shorty after uploading")
	flag.StringVar(&name, "name", "", "Name of the file. Useful if piping into stdin")
	flag.BoolVar(&verbose, "verbose", false, "Make it loud")
	flag.StringVar(&del, "delete", "", "When passed, this function will delete the specified object")
	flag.Parse()

	if del != "" {
		err = deleteObject(bucket, del)
		if err != nil {
			oops(err)
		}

		fmt.Printf("%s is gone\n", del)
		return
	}

	if name == "" {
		name = flag.Arg(0)
	}

	if name == "" {
		oops("Please pass a file name when piping to stdin (if this isn't what you were expecting, see --help).")
	}

	usingStdin := flag.Arg(0) == ""

	var obj *storage.ObjectHandle
	if usingStdin {
		if verbose {
			fmt.Printf("Uploading file %s from stdin...\n", name)
		}
		obj, err = upload(bucket, os.Stdin, name)
	} else {
		fil, err := os.Open(name)
		if err != nil {
			oops(err)
		}

		fn := path.Base(name)

		if verbose {
			fmt.Printf("Uploading %s...\n", fn)
		}
		obj, err = upload(bucket, fil, fn)
	}

	if err != nil {
		oops(err)
	}

	fqURL, err := url.Parse(fmt.Sprintf("%s/%s", globalConfig.Host, obj.ObjectName()))
	if err != nil {
		oops(err)
	}

	if globalConfig.Shorten || shorten {
		if verbose {
			fmt.Println("Shortening...")
		}
		short, err := shorty.Shorten(fqURL.String(), globalConfig.ShortyConfig)
		if err != nil {
			fmt.Printf("Failed to shorten (%v). Long URL: %s\n", err, fqURL.String())
			return
		}

		fmt.Println(short)
		return
	}

	fmt.Println(fqURL)
}

func upload(bucket *storage.BucketHandle, file *os.File, name string) (*storage.ObjectHandle, error) {
	obj := newObject(bucket, name)

	w := obj.NewWriter(context.Background())

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func newObject(bucket *storage.BucketHandle, name string) *storage.ObjectHandle {
	now := time.Now()
	path := fmt.Sprintf("%d/%d/%d/%d_%s", now.Year(), now.Month(), now.Day(), now.Unix(), name)

	obj := bucket.Object(path)
	return obj
}

func deleteObject(bucket *storage.BucketHandle, name string) error {
	obj := bucket.Object(name)
	return obj.Delete(context.Background())
}
