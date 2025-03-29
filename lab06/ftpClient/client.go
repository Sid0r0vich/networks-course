package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jlaffaye/ftp"
)

func main() {
	ftpServer := "127.0.0.1"
	username := "TestUser"
	password := "Sidor2004"

	c, err := ftp.Dial(ftpServer + ":21")
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login(username, password)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	listFiles(c)

	localFilePath := "file.txt"
	err = uploadFile(c, localFilePath)
	if err != nil {
		log.Fatal(err)
	}

	remoteFilePath := "serverFile.txt"
	err = downloadFile(c, remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

func listFiles(c *ftp.ServerConn) {
	entries, err := c.List("/")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Список файлов и директорий:")
	for _, entry := range entries {
		fmt.Println(entry.Name)
	}
}

func uploadFile(c *ftp.ServerConn, localFilePath string) error {
	file, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = c.Stor(filepath.Base(localFilePath), file)
	if err != nil {
		return err
	}

	fmt.Printf("Файл %s успешно загружен на сервер\n", localFilePath)
	return nil
}

func downloadFile(c *ftp.ServerConn, remoteFilePath string) error {
	response, err := c.Retr(remoteFilePath)
	if err != nil {
		return err
	}
	defer response.Close()

	localFile, err := os.Create(remoteFilePath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, response)
	if err != nil {
		return err
	}

	fmt.Printf("Файл %s успешно загружен локально\n", remoteFilePath)
	return nil
}
