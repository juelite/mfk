package logic

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"mfk/common"
)

var (
	DirSep 			string
	NewProject 		string
)

func init() {
	if runtime.GOOS == "windows" {
		DirSep = "\\"
	} else {
		DirSep = "/"
	}
}

func Run() {
	if NewProject == "" {
		fmt.Println("请输入项目名称")
		os.Exit(0)
	}
	err, res, _ := common.ExecCommand("pwd")
	if err != nil {
		fmt.Println("获取当前路径出错")
	}



	dir := strings.Replace(res, "\n", "", -1) + DirSep + NewProject

	err = os.Mkdir(dir, 0755)

	// 写入启动脚本
	root_main(dir)

	if err != nil {
		fmt.Println("目录创建失败：", err.Error())
	}

	dir_map := []string{
		"app",
		"cmd",
		"server",
	}

	// 创建文件夹
	for _, v := range dir_map {
		path := dir + DirSep + v
		fmt.Println("写入目录：", v)
		err := os.Mkdir(path, 0755)
		//写入文件下文件
		if err != nil {
			fmt.Println("没有权限创建目录")
		}
	}
	app_files(dir + DirSep + "app")
	cmd_files(dir + DirSep + "cmd")
	server_files(dir + DirSep + "server")

	fmt.Println("初始化完成！")
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

	"`+NewProject+`/app"
	pb "`+NewProject+`/app/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	GrpcPort 		string
	HttpPort 		string
)

func Run() {
	GrpcEndPoint := ":" + GrpcPort
	HttpEndPoint := ":" + HttpPort

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
		fmt.Println("写入文件失败")
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
		fmt.Println("写入目录：", v)
		err := os.Mkdir(path, 0755)
		//写入文件下文件
		if err != nil {
			fmt.Println("没有权限创建目录")
		}
	}
	app_client(path + DirSep + "client")
	app_proto(path + DirSep + "protos")
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
	pb "`+NewProject+`/app/proto"
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
	cmd_server(path)
}

// cmd/server.go
func cmd_server(path string) {
	ctn := `package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"demo/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
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

func init() {
	serverCmd.Flags().StringVarP(&server.GrpcPort, "GrpcPort", "", "50001", "grpc port")
	serverCmd.Flags().StringVarP(&server.HttpPort, "HttpPort", "", "40001", "http port")
	rootCmd.AddCommand(serverCmd)
}`

	f, err := os.Create(path + DirSep + "server.go")
	if err != nil {
		fmt.Println("写入文件失败")
		return
	}
	defer f.Close()

	f.Write([]byte(ctn))
}

// cmd/root.go
func cmd_root(path string) {
	ctn := `package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the gRPC server",
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
}