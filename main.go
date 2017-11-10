package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/burntsushi/toml"
	"github.com/masci/flickr"
	"log"
	"os"
	"os/user"
	"strings"
)

type Config struct {
	Apikey      string
	Apisecret   string
	Oauthtoken  string
	Oauthsecret string
	Basepath    string
}

func main() {
	// retrieve Flickr credentials from env vars

	privacy := flag.Bool("private", false, "Privacy settings")
	flag.Parse()

	config := ReadConfig()
	apik := config.Apikey
	apisec := config.Apisecret
	token := config.Oauthtoken
	tokenSecret := config.Oauthsecret

	// create an API client with credentials
	client := flickr.NewFlickrClient(apik, apisec)
	client.OAuthToken = token
	client.OAuthTokenSecret = tokenSecret

	// get a list of filepaths
	// upload first filepath (add some sort of progress meter here)
	// add to a set
	// for each remaining filepath (upload and add to set w/progress meter)

	filepaths := getFilePaths(config.Basepath)
	number_files := len(filepaths)
	fmt.Println("Number of images to upload: ", number_files)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter photoset: ")
	photoset_name, _ := reader.ReadString('\n')
	photoset_name = strings.TrimSpace(photoset_name)

	// Run a shift on the filepaths slice
	first_image, filepaths := filepaths[0], filepaths[1:]
	fmt.Println("Uploading initial image and creating photoset..")
	photoset_id, err := uploadImageAndCreateSet(config.Basepath, first_image, client, photoset_name, *privacy)
	if err != nil {
		fmt.Println("Error uploading first photo and creating photoset", err)
	} else {
		fmt.Println("Done!")
	}

	x := 2
	fmt.Println("Number of images to upload: ", number_files)

	for _, f := range filepaths {
		fmt.Printf("Uploading %d of %d images.\n", x, number_files)
		_, err := uploadImageToSet(config.Basepath, f, client, photoset_id, *privacy)
		if err != nil {
			fmt.Println("Error uploading image", f.Name())
		}
		x++
	}
	fmt.Println("Uploading complete. All done!")
}

func ReadConfig() Config {
	usr, _ := user.Current()
	dir := usr.HomeDir
	configfile := dir + "/.flickrfolder.toml"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
