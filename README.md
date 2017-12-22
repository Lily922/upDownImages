文件
	test.go文件为main函数文件,主要涉及创建Dockerclient
	operation.go文件为方法文件包括上传镜像与下载镜像方法以及子方法


上传镜像
UploadImage(imageName string, tarPath string, ctx context.Context, authConfig types.AuthConfig)
参数：
	imageName	为上传到镜像仓库的镜像名称：（例：”10.10.10.101/library/ubuntu:14-04"）
	tarPath  为上传到镜像仓库的镜像的tar包绝对路径：（例：“/home/tom/imagetar/ubuntu-14-04.tar”）
	ctx		上下文参数
	authConfig	仓库权限认证
返回值：
	err：是否出错
	
下载镜像
DownloadImage(imageName string, ctx context.Context, parentPath string, authConfig types.AuthConfig)
参数：
	imageName	从镜像仓库下载下来的镜像名称：（例：“10.10.10.101/library/ubuntu:14.04"）
	parentPath	服务器存放tar文件的目录，使用绝对路径（例："/home/tom/tartest/test"）
返回值：
	tarPath		服务器上tar文件的位置，使用绝对路径（例："/home/lily/tartest/test/ubuntu-14-04.tar"）
