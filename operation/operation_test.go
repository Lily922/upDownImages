package operation

import (
	"context"
	"testing"

	"net/http"
	"os"
	"path/filepath"

	//"github.com/docker/engine-api/client"
	//"github.com/docker/engine-api/types"
	"github.com/docker/docker/client"
	//"github.com/docker/docker/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/tlsconfig"
	//"github.com/docker/go-connections/tlsconfig"
)

const DefaultDockerHost = "tcp://10.10.101.194:2376"

const DefaultVersion string = "1.23"

const parentPath string = "/home/lily/tartest/test"

var ()

func TestDockerPull(t *testing.T) {
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		t.Fatalf(err.Error())
	}
	var cli DockerClient = DockerClient{client}
	/*
		docker load
	*/
	imageName := "10.10.101.175/lily/mysql:harbor"
	authConfig := types.AuthConfig{}
	authConfig.Username = "admin"
	authConfig.Password = "123"
	authConfig.ServerAddress = "10.10.101.175"
	authConfig.IdentityToken = ""
	err = cli.DockerPull(imageName, ctx, authConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDockerSave(t *testing.T) {
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		t.Fatalf(err.Error())
	}
	var cli DockerClient = DockerClient{client}

	var images = []string{
		"10.10.101.175/lily/mysql:harbor",
	}
	/*
		docker save
	*/
	err = cli.DockerSave(images, ctx, parentPath)
	if err != nil {
		t.Fatalf(err.Error())
	}

}

func TestDockerLoad(t *testing.T) {
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		t.Fatalf(err.Error())
	}
	var cli DockerClient = DockerClient{client}
	/*
		docker load
	*/
	//imageName := "10.10.101.175/lily/mysql:harbor"
	tarPath := "/home/lily/tartest/test/10.10.101.175/lily/mysql:harbor.tar"
	err = cli.DockerLoad(ctx, tarPath)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDockerTag(t *testing.T) {
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		t.Fatalf(err.Error())
	}
	var cli DockerClient = DockerClient{client}
	/*
		docker tag
	*/

	imageNew := "10.10.101.175/lily/mysql:harbor"
	imageOld := "10.10.101.175/lily/mysql:test"
	err = cli.DockerTag(imageNew, imageOld, ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
func TestDockerPush(t *testing.T) {
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		t.Fatalf(err.Error())
	}
	var cli DockerClient = DockerClient{client}
	imageName := "10.10.101.175/lily/mysql:harbor"
	authConfig := types.AuthConfig{}
	authConfig.Username = "admin"
	authConfig.Password = "123"
	authConfig.ServerAddress = "10.10.101.175"
	authConfig.IdentityToken = ""
	err = cli.DockerPush(imageName, ctx, authConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func newClient() (*client.Client, error) {
	var client1 *http.Client
	dockerCertPath := ""
	options := tlsconfig.Options{
		InsecureSkipVerify: os.Getenv("DOCKER_TLS_VERIFY") == "",
		CAFile:             filepath.Join(dockerCertPath, "/home/lily/cert/docker2376/ca.pem"),
		CertFile:           filepath.Join(dockerCertPath, "/home/lily/cert/docker2376/client-cert.pem"),
		KeyFile:            filepath.Join(dockerCertPath, "/home/lily/cert/docker2376/client-key.pem"),
	}
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
