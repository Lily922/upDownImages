package operation

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/pkg/jsonmessage"
	//"github.com/docker/engine-api/client"
	//"github.com/docker/engine-api/types"
	"github.com/docker/docker/client"
	//"github.com/docker/docker/types"
	"github.com/docker/docker/api/types"
)

/*
	docker客户端
*/
type DockerClient struct {
	*client.Client
}

/*
	docker save用镜像名称save，不要用镜像id
	images是需要save的images数组，目前我们默认数组的长度为1，parentPath是镜像save的路径
*/
func (cli *DockerClient) DockerSave(images []string, ctx context.Context, parentPath string) (tarPath string, err error) {
	fmt.Println("save start")
	if len(images) == 0 {
		fmt.Println("the image list is empty")
		err := errors.New("empty images")
		return "", err
	}
	paths := strings.Split(images[0], "/")
	imagePath := ""
	if len(paths) > 1 {
		for _, p := range paths[0 : len(paths)-1] {
			imagePath = imagePath + p + "/"
		}
	}

	path := parentPath + "/" + imagePath
	err = os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Printf("%s", err)
		return "", err
	}
	ss, err := cli.ImageSave(ctx, images)
	if err != nil {
		fmt.Println("error" + err.Error())
		return "", err
	}
	tarPath = parentPath + "/" + images[0] + ".tar"
	//	w, err := os.Create(parentPath + "/" + images[0] + ".tar")
	w, err := os.Create(tarPath)
	fmt.Println(parentPath)
	if err != nil {
		fmt.Printf("xx:", err)
		return "", err
	}
	defer w.Close()
	_, err = io.Copy(w, ss)
	if err != nil {
		fmt.Println("yy:" + err.Error())
		return "", err
	}
	fmt.Println("save finished")
	return tarPath, nil
}

/*
	docker load (docker save的时候用镜像名称，不要用镜像id，这样load的镜像名称就是压缩时的镜像)
*/

func (cli *DockerClient) DockerLoad(ctx context.Context, tarPath string) (image string, err error) {
	//file, err := os.Open("/home/lily/tartest/test/tets.tar")
	//file, err := os.Open(parentPath + "/" + imageName + ".tar")
	fmt.Println("load start")
	file, err := os.Open(tarPath)
	if err != nil {
		fmt.Println("cc" + err.Error())
		return "", err
	}
	defer file.Close()

	response, err := cli.ImageLoad(ctx, file, false)
	if err != nil {
		fmt.Println("dd" + err.Error())
		return "", err
	}
	//io.Copy(os.Stdout, response.Body)
	defer response.Body.Close()
	//fmt.Printf("load finished")

	w := bytes.NewBuffer(make([]byte, 0))
	var outFd uintptr
	err = jsonmessage.DisplayJSONMessagesStream(response.Body, w, outFd, false, nil)
	if err != nil {
		fmt.Println("er" + err.Error())
		return "", err
	}
	fmt.Println(w.String())
	image, err = GetImageString(w.String())
	if err != nil {
		return "", err
	}
	fmt.Printf("load finished")

	return image, nil
}

/*
	docker pull 从镜像仓库拉取镜像
*/

func (cli *DockerClient) DockerPull(imageName string, ctx context.Context, authConfig types.AuthConfig) (err error) {
	fmt.Println("pull start")
	response, err := cli.RegistryLogin(ctx, authConfig)
	if err != nil {
		fmt.Println("sd" + err.Error())
		return err
	}
	fmt.Println(response.IdentityToken)
	fmt.Println(response.Status)

	authConfig.IdentityToken = response.IdentityToken
	encodedAuth, err := EncodeAuthToBase64(authConfig)
	options := types.ImagePullOptions{
		RegistryAuth: encodedAuth,
		//PrivilegeFunc: requestPrivilege,

	}
	//注释掉 out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	out, err := cli.ImagePull(ctx, imageName, options)
	defer out.Close()
	if err != nil {
		fmt.Println("zz:" + err.Error())
		return err
	}
	//io.Copy(os.Stdout, out)
	w := bytes.NewBuffer(make([]byte, 0))
	var outFd uintptr
	err = jsonmessage.DisplayJSONMessagesStream(out, w, outFd, false, nil)
	if err != nil {
		fmt.Println("er" + err.Error())
		return err
	}
	fmt.Println(w.String())
	fmt.Println("pull finished")
	return nil
}

/*
	docker push 将镜像push到远程镜像仓库
*/

func (cli *DockerClient) DockerPush(imageName string, ctx context.Context, authConfig types.AuthConfig) (err error) {
	fmt.Println("push start")
	response, err := cli.RegistryLogin(ctx, authConfig)
	if err != nil {
		fmt.Println("sd" + err.Error())
		return err
	}
	fmt.Println(response.IdentityToken)
	fmt.Println(response.Status)

	authConfig.IdentityToken = response.IdentityToken
	encodedAuth, err := EncodeAuthToBase64(authConfig)
	options := types.ImagePushOptions{
		RegistryAuth: encodedAuth,
		//PrivilegeFunc: requestPrivilege,

	}

	responseBody, err := cli.ImagePush(ctx, imageName, options)
	if err != nil {
		fmt.Println("er" + err.Error())
		return err
	}
	fmt.Println("end")
	//io.Copy(os.Stdout, responseBody)
	//注释掉 w, err := os.Create("/home/lily/tartest/test.txt")
	//注释掉  defer w.Close()
	w := bytes.NewBuffer(make([]byte, 0))
	var outFd uintptr
	err = jsonmessage.DisplayJSONMessagesStream(responseBody, w, outFd, false, nil)
	if err != nil {
		fmt.Println("er" + err.Error())
		return err
	}
	fmt.Println(w.String())
	defer responseBody.Close()
	fmt.Println("push finished")
	return nil
}

/*
	docker tag，imageOld目前是imageID
*/

func (cli *DockerClient) DockerTag(imageNew, imageOld string, ctx context.Context) (err error) {
	fmt.Println("tag start")
	err = cli.ImageTag(ctx, imageOld, imageNew)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("tag finished")
	return err
}

/*
	docker delete,删除中间产生的镜像
*/

func (cli *DockerClient) DockerDelete(imageName string, ctx context.Context) (err error) {
	fmt.Println("delete start")

	options := types.ImageRemoveOptions{}
	_, err = cli.ImageRemove(ctx, imageName, options)
	if err != nil {
		return err
	}
	return nil
}

/*
   完成docker镜像由tar包上传到镜像仓库
*/

func (cli *DockerClient) UploadImage(imageName string, tarPath string, ctx context.Context, authConfig types.AuthConfig) (err error) {
	image, err := cli.DockerLoad(ctx, tarPath)

	if err != nil {
		fmt.Println("err1")
		return err
	}
	fmt.Println("iamge:" + image)
	err = cli.DockerTag(imageName, image, ctx)
	if err != nil {
		fmt.Println("err2")
		return err
	}
	err = cli.DockerPush(imageName, ctx, authConfig)
	if err != nil {
		fmt.Println("err3")
		return err
	}
	return nil
}

/*
完成docker镜像由镜像仓库下载成tar包
*/

func (cli *DockerClient) DownloadImage(imageName string, ctx context.Context, parentPath string, authConfig types.AuthConfig) (tarPath string, err error) {
	err = cli.DockerPull(imageName, ctx, authConfig)
	if err != nil {
		return "", err
	}
	images := []string{
		imageName,
	}
	tarPath, err = cli.DockerSave(images, ctx, parentPath)
	if err != nil {
		return "", err
	}
	return tarPath, nil
}

func EncodeAuthToBase64(authConfig types.AuthConfig) (string, error) {
	buf, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}

//从字符转中获取imageID
func GetImageString(in string) (out string, err error) {
	ins := strings.Split(in, ":")
	out = ""
	if len(ins) > 1 {
		out = ins[1]
		//去掉空格
		out = strings.Replace(out, " ", "", -1)
		if len(ins) == 2 {
			//去掉回车
			out = strings.Replace(out, "\n", "", -1)
			return out, nil
		} else {
			for _, i := range ins[2:len(ins)] {
				out = out + ":" + i
			}
		}
		//去掉回车
		out = strings.Replace(out, "\n", "", -1)
		return out, nil
		fmt.Println(out)
	}
	return out, errors.New("empty imageName or imageId")

}
