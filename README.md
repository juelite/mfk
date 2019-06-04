## 蜜蜂控微服务脚手架说明文档

### 使用文档

#### 1. 获取脚手架工具

```$xslt
    go get github.com/juelite/mfk
```

#### 2. 编译脚手架工具并将其放入环境变量

```$xslt
    cd $GOAPTH/src/github.com/juelite/mfk
    
    go build -o $GOPATH/bin main.go
```

#### 3. 创建服务目录及结构

```$xslt
    mfk new -n demo
```

创建成功后结构如下：

```$xslt
    .
    ├── README.md               说明文档
    ├── app                     服务目录（所有服务逻辑都在此编写）
    │   ├── client              测试客户端目录
    │   │   └── client.go       测试客户端
    │   ├── conf                配置文件目录
    │   │   ├── app.conf        项目配置文件
    │   │   └── env.conf        运行环境文件
    │   ├── main.go             业务入口
    │   ├── logic               业务逻辑
    │   ├── model               业务模型
    │   └── proto               protobuf文件及生成的pb文件
    │       ├── hello.proto     服务元信息文件
    │       └── protoc.sh       将服元原信息生成pb文件脚本
    ├── cmd                     终端辅助文件目录
    │   └── server.go           
    ├── main.go                 服务启动为恩建
    ├── server                  
    │   └── server.go           服务注册脚本
    └── util                    工具包
        └── config              配置文件解析包
            └── config.go       获取配置方法

```

#### 4. 编写服务元信息文件 （app/proto）:

示例：

```$xslt
    syntax = "proto3";
    
    package proto;
    
    import "google/api/annotations.proto";
    
    // 服务名称
    service HelloWorld {
        
        // 方法名称
        rpc SayHelloWorld(HelloWorldRequest) returns (HelloWorldResponse) {
            
            // 定义http访问路由和方法
            option (google.api.http) = {
                post: "/hello_world"
                body: "*"
            };
        }
    }
    
    // 请求参数定义
    message HelloWorldRequest {
        string referer = 1;
    }
    
    // 返回参数定义
    message HelloWorldResponse {
        string message = 1;
    }
```

编写完成后 执行 ./protoc.sh （脚本里面proto文件名需要与上面文件名一致）

会生成 .pb.go 和 .pb.gw.go 两个文件 分别文grpc和grpc-gateway文件

#### 5. 编写业务入口文件（app/main.go）

```$xslt
    package app
    
    import (
    	pb "demo/app/proto"
    	"golang.org/x/net/context"
    )
    
    type helloSrv struct{}
    
    func NewHelloSrv() *helloSrv {
    	return &helloSrv{}
    }
    
    func (h helloSrv) SayHelloWorld(ctx context.Context, r *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
    	// TODO your logic
    	return &pb.HelloWorldResponse{
    		Message: "test",
    	}, nil
    }
```

#### 6. 注册服务（server/server.go）

修改 【pb.*】方法为实际的接口方法

#### 7. 运行服务

```$xslt
    // 不指定端口默认 grpc 50001 http 40001
    go run main.go server [-GrpcPort 50001] [-HttpPort 40001]  
```