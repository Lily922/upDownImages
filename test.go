package main

import (
	"context"
	"fmt"

	"net/http"
	"os"
	"path/filepath"

	"upDownImage/operation"

	"github.com/docker/docker/client"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/tlsconfig"
)

/*
	docker deamon开放2375端口，在/lib/systemd/system/docker.service中进行配置。
*/

var DefaultDockerHost = "tcp://10.10.101.194:2376"

var DefaultVersion string = "1.23"

var parentPath string = "/home/lily/tartest/test"

var authConfig types.AuthConfig = types.AuthConfig{
	Username:      "admin",
	Password:      "123",
	ServerAddress: "10.10.101.175",
	IdentityToken: "",
}
var dockerCertPath string = ""
var options tlsconfig.Options = tlsconfig.Options{
	//CAFile:             filepath.Join(dockerCertPath, "ca.pem"),
	//CertFile:           filepath.Join(dockerCertPath, "cert.pem"),
	//KeyFile:            filepath.Join(dockerCertPath, "key.pem"),
	InsecureSkipVerify: os.Getenv("DOCKER_TLS_VERIFY") == "",
	CAFile:             filepath.Join(dockerCertPath, "/home/lily/cert/docker2376/ca.pem"),
	CertFile:           filepath.Join(dockerCertPath, "/home/lily/cert/docker2376/client-cert.pem"),
	KeyFile:            filepath.Join(dockerCertPath, "/home/lily/cert/docker2376/client-key.pem"),
	//InsecureSkipVerify: false,
}

var imageName string = "10.10.101.175/lily/scheduler:v1.5.3"
var uploadtarPath string = "/home/lily/tartest/test/gcr.io/google_containers/kubedns-amd64:1.9.tar"

func main() {
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		fmt.Printf("error happens when create newClient")
	}
	var cli operation.DockerClient = operation.DockerClient{client}
	/*
		var images = []string{
			//"gcr.io/google_containers/kubedns-amd64:1.9",
			"90e98111e9f1",
		}
	*/
	//imageName := "10.10.101.175/lily/ubedns-amd64:1.9"
	/*
		authConfig := types.AuthConfig{}
		authConfig.Username = "admin"
		authConfig.Password = "123"
		authConfig.ServerAddress = "10.10.101.175"
		authConfig.IdentityToken = ""
	*/
	/*
	   docker tag
	*/
	/*
		imageNew := "nginx:lily"
		imageOld := "nginx:1.9"
		err = cli.DockerTag(imageNew, imageOld, ctx)
		if err != nil {
			fmt.Println("main" + err.Error())
		}
	*/
	/*
		docker save
	*/
	/*
		err = cli.DockerSave(images, ctx, parentPath)
		if err != nil {
			fmt.Println("main" + err.Error())
		}
	*/
	/*
		docker load
	*/
	/*
		tarPath := "/home/lily/tartest/test/gcr.io/google_containers/kubedns-amd64:1.9.tar"
		image, err := cli.DockerLoad(ctx, tarPath)
		if err != nil {
			fmt.Println("main" + err.Error())
		}
		fmt.Println(image)
	*/
	/*
		docker pull
	*/
	/*
		imageName := "10.10.101.175/lily/mysql:harbor"
		err = cli.DockerPull(imageName, ctx, authConfig)
		if err != nil {
			fmt.Println("main" + err.Error())
		}
	*/
	/*
		docker push
	*/
	/*
		err = cli.DockerPush(imageName, ctx, authConfig)
		if err != nil {
			fmt.Println("main" + err.Error())
		}
	*/
	/*
		docker delete
	*/
	/*
		err = cli.DockerDelete(imageName, ctx)
		if err != nil {
			fmt.Println("main" + err.Error())
		}
	*/

	/*
		image upload
	*/
	/*
		imageName := "10.10.101.175/lily/kubedns-amd64:test"
		uploadtarPath := "/home/lily/tartest/test/gcr.io/google_containers/kubedns-amd64:1.9.tar"
	*/

	err = cli.UploadImage(imageName, uploadtarPath, ctx, authConfig)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("finished upload")
	}

	/*
		image download
	*/

	//imageName := "10.10.101.175/lily/scheduler:v1.5.3"

	tarPath, err := cli.DownloadImage(imageName, ctx, parentPath, authConfig)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("finihed download" + tarPath)
	}

}

func newClient() (*client.Client, error) {
	var client1 *http.Client
	/*
		dockerCertPath := ""
		options := tlsconfig.Options{
			InsecureSkipVerify: os.Getenv("DOCKER_TLS_VERIFY") == "",
			CAFile:             filepath.Join(dockerCertPath, "/home/lily/cert/docker2376/ca.pem"),
			CertFile:           filepath.Join(dockerCertPath, "/home/lily/cert/docker2376/client-cert.pem"),
			KeyFile:            filepath.Join(dockerCertPath, "/home/lily/cert/docker2376/client-key.pem"),
		}
	*/
	tlsc, err := tlsconfig.Client(options)
	if err != nil {
		return nil, err
	}

	client1 = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsc,
		},
	}

	host := DefaultDockerHost
	version := DefaultVersion
	return client.NewClient(host, version, client1, nil)
}
