package logic

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	DirSep     string
	NewProject string
)

func init() {
	if runtime.GOOS == "windows" {
		DirSep = "\\"
	} else {
		DirSep = "/"
	}
}

func New() {
	if NewProject == "" {
		fmt.Println("please enter project name !")
		os.Exit(0)
	}
	res, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("get current path error !")
	}

	dir := strings.Replace(res, "\n", "", -1) + DirSep + NewProject

	err = os.Mkdir(dir, 0755)

	// 写入启动脚本
	root_main(dir)

	if err != nil {
		fmt.Println(fmt.Sprintf("mkdir [%s] error : %s !", NewProject, err.Error()))
	}

	dir_map := []string{
		"app",
		"cmd",
		"server",
		"util",
	}

	// 创建文件夹
	for _, v := range dir_map {
		path := dir + DirSep + v
		fmt.Println(fmt.Sprintf("create dir: [%s]", v))
		err := os.Mkdir(path, 0755)
		//写入文件下文件
		if err != nil {
			fmt.Println(fmt.Sprintf("mkdir [%s] error : %s !", v, err.Error()))
			return
		}
	}
	app_files(dir + DirSep + "app")
	cmd_files(dir + DirSep + "cmd")
	server_files(dir + DirSep + "server")
	util_files(dir + DirSep + "util")

	fmt.Println(fmt.Sprintf("new [%s] complete !", NewProject))
}

// util文件
func util_files(path string) {
	err := os.Mkdir(path+DirSep+"config", 0755)
	if err != nil {
		fmt.Println("mkdir [config] error !")
		return
	}
	ctn := `package config

import (
	"log"
	"strconv"
	"sync"

	"github.com/Unknwon/goconfig"
)

const (
	conf_dir 	= "./app/conf/"
)

var (
	env 		string
	conf 		map[string]string
	err 		error
)

func init() {
	once := &sync.Once{}
	once.Do(loadConf)
}

func loadConf() {
	c , err := goconfig.LoadConfigFile(conf_dir + "env.conf")
	if err != nil {
		log.Fatal("load env file error")
	}
	env , err = c.GetValue("" , "runmode")

	if err != nil {
		log.Fatal("parse env file error")
	}

	c , err = goconfig.LoadConfigFile(conf_dir + "app.conf")
	if err != nil {
		log.Fatal("load config file error")
	}
	conf, err = c.GetSection(env)
	if err != nil {
		log.Fatal("parse config file error")
	}
}

// 获取运行环境
func GetEnv() string {
	return env
}

// 获取string
func GetString(key string) string {
	if val, ok := conf[key]; ok {
		return val
	}
	return ""
}

// 获取int64
func GetInt64(key string) int64 {
	if val, ok := conf[key]; ok {
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	}
	return 0
}

// 获取bool
func GetBool(key string) bool {
	var res bool
	if val, ok := conf[key]; ok {
		if val == "true" {
			res = true
		}
	}
	return res
}`
	f, err := os.Create(path + DirSep + "config" + DirSep + "config.go")
	if err != nil {
		fmt.Println("writer filer err")
		return
	}
	defer f.Close()

	f.Write([]byte(ctn))
}

// server文件
func server_files(path string) {
	ctn := `// 服务监听脚本 把pb.*（RegisterHelloWorldServer，RegisterHelloWorldHandlerFromEndpoint）方法改成对应的方法
package server

import (
	"log"
	"net"
	"net/http"
	"os"

	"` + NewProject + `/app"
	pb "` + NewProject + `/app/protos"
	"` + NewProject + `/util/config"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)


func Run() {
	GrpcEndPoint := ":" + config.GetString("grpc_port")
	HttpEndPoint := ":" + config.GetString("http_port")

	go startGrpc(GrpcEndPoint)
	startHttp(HttpEndPoint, GrpcEndPoint)
}

// gRPC
func startGrpc(GrpcEndPoint string) {
	conn, err := net.Listen("tcp", GrpcEndPoint)
	if err != nil {
		log.Printf("TCP Listen err:%v\n", err)
	}
	server := grpc.NewServer()
	// 这里需要改成对应方法
	pb.RegisterHelloWorldServer(server, app.NewHelloSrv())
	log.Println("grpc is running on ", GrpcEndPoint)
	err = server.Serve(conn)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
}

// http
func startHttp(HttpEndPoint, GrpcEndPoint string) {
	mux := runtime.NewServeMux()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	// 这里需要改成对应方法
	if err := pb.RegisterHelloWorldHandlerFromEndpoint(ctx, mux, GrpcEndPoint, opts); err != nil {
		log.Println(err)
		os.Exit(0)
	}
	log.Println("http is running on ", HttpEndPoint)
	if err := http.ListenAndServe(HttpEndPoint, mux); err != nil {
		log.Println(err)
		os.Exit(0)
	}
}`
	f, err := os.Create(path + DirSep + "server.go")
	if err != nil {
		fmt.Println("writer filer err")
		return
	}
	defer f.Close()

	f.Write([]byte(ctn))
}

// app文件
func app_files(path string) {
	dir_map := []string{
		"client",
		"conf",
		"logics",
		"models",
		"protos",
	}

	// 创建文件夹
	for _, v := range dir_map {
		path := path + DirSep + v
		fmt.Println(fmt.Sprintf("create dir: [%s]", v))
		err := os.Mkdir(path, 0755)
		//写入文件下文件
		if err != nil {
			fmt.Println("没有权限创建目录")
		}
	}
	app_client(path + DirSep + "client")
	app_proto(path + DirSep + "protos")
	app_config(path + DirSep + "conf")
}

// app_config
func app_config(path string) {
	ctn := `# 运行模式
runmode = prod`
	f, err := os.Create(path + DirSep + "env.conf")
	if err != nil {
		fmt.Println("写入文件失败")
		return
	}
	defer f.Close()

	f.Write([]byte(ctn))

	ctn1 := `[dev]
http_port       = 40001
grpc_port       = 50001

[release]
http_port       = 40002
grpc_port       = 50002

[prod]
http_port       = 40003
grpc_port       = 50003`

	f, err = os.Create(path + DirSep + "app.conf")
	if err != nil {
		fmt.Println("写入文件失败")
		return
	}
	defer f.Close()

	f.Write([]byte(ctn1))
}

// app proto
func app_proto(path string) {
	ctn := `// proto示例文件，按需求改写
/*
syntax = "proto3";

package proto;

import "google/api/annotations.proto";

service HelloWorld {
    rpc SayHelloWorld(HelloWorldRequest) returns (HelloWorldResponse) {
        option (google.api.http) = {
            post: "/hello_world"
            body: "*"
        };
    }
}

message HelloWorldRequest {
    string referer = 1;
}

message HelloWorldResponse {
    string message = 1;
}*/`
	f, err := os.Create(path + DirSep + "demo.proto")
	if err != nil {
		fmt.Println("写入文件失败")
		return
	}
	defer f.Close()

	f.Write([]byte(ctn))

	ctn1 := `#!/usr/bin/env bash
#根据proto文件生成pb文件，请讲下面路径文件名按真实名称替换
protoc --go_out=plugins=grpc:. hello.proto

protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. demo.proto`

	f, err = os.Create(path + DirSep + "protoc.sh")
	if err != nil {
		fmt.Println("写入文件失败")
		return
	}
	defer f.Close()

	f.Write([]byte(ctn1))
}

// app client client.go
func app_client(path string) {
	ctn := `package main
//grpc测试客户端
/*import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "` + NewProject + `/app/protos"
)

func main() {
	conn, err := grpc.Dial(":50001", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// 需改成对应方法
	c := pb.NewHelloWorldClient(conn)
	ctx := context.Background()
	// 需改成对应方法
	r := &pb.HelloWorldRequest{
		Referer: "GRPC",
	}
	// 需改成对应方法
	res, err := c.SayHelloWorld(ctx, r)
	fmt.Println(res, err)
}*/`
	f, err := os.Create(path + DirSep + "client.go")
	if err != nil {
		fmt.Println("写入文件失败")
		return
	}
	defer f.Close()

	f.Write([]byte(ctn))
}

// cmd文件
func cmd_files(path string) {
	cmd_root(path)
}

// cmd/root.go
func cmd_root(path string) {
	ctn := `package cmd

import (
	"fmt"
	"os"
	"log"

	"github.com/spf13/cobra"
	"` + NewProject + `/server"
)

var rootCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recover error : ", err)
			}
		}()
		server.Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}`

	f, err := os.Create(path + DirSep + "root.go")
	if err != nil {
		fmt.Println("写入文件失败")
		return
	}
	defer f.Close()

	f.Write([]byte(ctn))
}

// 启动文件 main.go
func root_main(dir string) {
	ctn := `// 启动脚本
package main

import (
	"` + NewProject + `/cmd"
	_ "` + NewProject + `/util/config"
)

func main() {
	cmd.Execute()
}`

	f, err := os.Create(dir + DirSep + "main.go")
	if err != nil {
		fmt.Println("写入文件失败")
		return
	}
	defer f.Close()

	f.Write([]byte(ctn))

	f2, err := os.Create(dir + DirSep + "README.md")
	if err != nil {
		fmt.Println("写入文件失败")
		return
	}
	defer f2.Close()

	f2.Write([]byte(readme()))
}

func readme() string {
	ctn := "## 蜜蜂控微服务脚手架说明文档\n"
	ctn += "\n"
	ctn += "### 使用文档\n"
	ctn += "\n"
	ctn += "#### 1. 获取脚手架工具\n"
	ctn += "\n"
	ctn += "```$xslt\n"
	ctn += "    go get github.com/juelite/mfk\n"
	ctn += "```\n"
	ctn += "\n"
	ctn += "#### 2. 编译脚手架工具并将其放入环境变量\n"
	ctn += "\n"
	ctn += "```$xslt\n"
	ctn += "    cd $GOAPTH/src/github.com/juelite/mfk\n"
	ctn += "    \n"
	ctn += "    go build -o $GOPATH/bin main.go\n"
	ctn += "```\n"
	ctn += "\n"
	ctn += "#### 3. 创建服务目录及机构\n"
	ctn += "\n"
	ctn += "```$xslt\n"
	ctn += "    mfk new -n demo\n"
	ctn += "```\n"
	ctn += "\n"
	ctn += "创建成功后结构如下：\n"
	ctn += "\n"
	ctn += "```$xslt\n"
	ctn += "    .\n"
	ctn += "    ├── README.md               说明文档\n"
	ctn += "    ├── app                     服务目录（所有服务逻辑都在此编写）\n"
	ctn += "    │   ├── client              测试客户端目录\n"
	ctn += "    │   │   └── client.go       测试客户端\n"
	ctn += "    │   ├── conf                配置文件目录\n"
	ctn += "    │   ├── main.go             业务入口\n"
	ctn += "    │   ├── logic               业务逻辑\n"
	ctn += "    │   ├── model               业务模型\n"
	ctn += "    │   └── proto               protobuf文件及生成的pb文件\n"
	ctn += "    │       ├── hello.proto     服务元信息文件\n"
	ctn += "    │       └── protoc.sh       将服元原信息生成pb文件脚本\n"
	ctn += "    ├── cmd                     终端辅助文件目录\n"
	ctn += "    │   ├── root.go             \n"
	ctn += "    │   └── server.go\n"
	ctn += "    ├── main.go                 服务启动为恩建\n"
	ctn += "    └── server                  server目录\n"
	ctn += "        └── server.go           server文件\n"
	ctn += "\n"
	ctn += "```\n"
	ctn += "\n"
	ctn += "#### 4. 编写服务元信息文件 （app/proto）:\n"
	ctn += "\n"
	ctn += "示例：\n"
	ctn += "\n"
	ctn += "```$xslt\n"
	ctn += "    syntax = \"proto3\";\n"
	ctn += "    \n"
	ctn += "    package proto;\n"
	ctn += "    \n"
	ctn += "    import \"google/api/annotations.proto\";\n"
	ctn += "    \n"
	ctn += "    // 服务名称\n"
	ctn += "    service HelloWorld {\n"
	ctn += "        \n"
	ctn += "        // 方法名称\n"
	ctn += "        rpc SayHelloWorld(HelloWorldRequest) returns (HelloWorldResponse) {\n"
	ctn += "            \n"
	ctn += "            // 定义http访问路由和方法\n"
	ctn += "            option (google.api.http) = {\n"
	ctn += "                post: \"/hello_world\"\n"
	ctn += "                body: \"*\"\n"
	ctn += "            };\n"
	ctn += "        }\n"
	ctn += "    }\n"
	ctn += "    \n"
	ctn += "    // 请求参数定义\n"
	ctn += "    message HelloWorldRequest {\n"
	ctn += "        string referer = 1;\n"
	ctn += "    }\n"
	ctn += "    \n"
	ctn += "    // 返回参数定义\n"
	ctn += "    message HelloWorldResponse {\n"
	ctn += "        string message = 1;\n"
	ctn += "    }\n"
	ctn += "```\n"
	ctn += "\n"
	ctn += "编写完成后 执行 ./protoc.sh （脚本里面proto文件名需要与上面文件名一致）\n"
	ctn += "\n"
	ctn += "会生成 .pb.go 和 .pb.gw.go 两个文件 分别文grpc和grpc-gateway文件\n"
	ctn += "\n"
	ctn += "#### 5. 编写业务入口文件（app/main.go）\n"
	ctn += "\n"
	ctn += "```$xslt\n"
	ctn += "    package app\n"
	ctn += "    \n"
	ctn += "    import (\n"
	ctn += "        pb \"demo/app/proto\"\n"
	ctn += "        \"golang.org/x/net/context\"\n"
	ctn += "    )\n"
	ctn += "    \n"
	ctn += "    type helloSrv struct{}\n"
	ctn += "    \n"
	ctn += "    func NewHelloSrv() *helloSrv {\n"
	ctn += "        return &helloSrv{}\n"
	ctn += "    }\n"
	ctn += "    \n"
	ctn += "    func (h helloSrv) SayHelloWorld(ctx context.Context, r *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {\n"
	ctn += "        // TODO your logic\n"
	ctn += "        return &pb.HelloWorldResponse{\n"
	ctn += "            Message: \"test\",\n"
	ctn += "        }, nil\n"
	ctn += "    }\n"
	ctn += "```\n"
	ctn += "\n"
	ctn += "#### 6. 注册服务（server/server.go）\n"
	ctn += "\n"
	ctn += "修改 【pb.*】方法为实际的接口方法\n"
	ctn += "\n"
	ctn += "#### 7. 运行服务\n"
	ctn += "\n"
	ctn += "```$xslt\n"
	ctn += "    // 不指定端口默认 grpc 50001 http 40001\n"
	ctn += "    go run main.go server [-GrpcPort 50001] [-HttpPort 40001]  \n"
	ctn += "```\n"

	return ctn
}
