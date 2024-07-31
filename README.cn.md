# Judge

[![Go Reference](https://pkg.go.dev/badge/github.com/axliupore/judge.svg)](https://pkg.go.dev/github.com/axliupore/judge) [![Go Report Card](https://goreportcard.com/badge/github.com/axliupore/judge)](https://goreportcard.com/report/github.com/axliupore/judge)

[English](README.md)

一个基于 Golang实现的代码沙箱，旨在提供一个安全、可靠、快捷的容器化环境去运行代码
，开发者只需提供编程语言和代码内容，即可获取代码执行后的结果。本项目的底层容器使用的是 [go-judge](https://github.com/criyle/go-judge)。

## 支持语言

|     语言     |   版本   |
|:----------:|:------:|
|    C++     |   17   |
|     C      |   17   |
|   Golang   |  1.19  |
|    Java    |   17   |
|   Python   | 3.11.2 |
| JavaScript |   18   |
| TypeScript |  5.53  |

## 快速开始

### 安装和运行

克隆本项目：

```bash
git clone git@github.com:axliupore/judge
```

在根目录下执行下面的命令：

```bash
make all
make run
```

或者使用 docker，不用 `Nsq` 直接用下面的方法即可：

```bash
docker run -d --privileged --shm-size=2048m -p 6048:6048 --name=judge trialoj/judge:0.0.1
```

如果使用 `Nsq` ，需要指定 `Nsq` 运行的位置，本项目默认 `Nsq` 是放置在主机上：

```bash
docker run -d --privileged --shm-size=2048m -p 6048:6048 --name=judge --add-host="host.docker.internal:host-gateway" trialoj/judge:0.0.1
```

### http

```bash
curl -X POST -H "Content-Type: application/json" -d '{"language": "cpp", "content":"#include <iostream>\nusing namespace std;\nint main() {\ncout << 10 << endl;\n}"}' http://127.0.0.1:6048/judge
```

### Nsq

在使用 `nsq` 消息队列进行消息的传递时候，需要创建一个生产者用于传递信息给到消费者，主题为 `judge_topic` ，用 `nsq`
传递运行代码必须传递一个参数 `nsq`
给到消费者，不然无法接受到代码执行的结果，详细代码在[这里](https://github.com/axliupore/judge/blob/master/pkg/nsq/nsq_test.go)

## 架构设计

![design](./doc/design.png)

## 参数说明

### 请求参数

|     参数      |             说明              |
|:-----------:|:---------------------------:|
|  Language   | 代码语言，目前支持的可以在 `judge` 目录下查看 |
|   Content   |            文件内容             |
|    Input    |            输入内容             |
|  CpuLimit   |          CPU 时间限制           |
| MemoryLimit |            内存限制             |
| StackLimit  |            栈空间限制            |
|  ProcLimit  |            进程数限制            |
|     Nsq     | 使用消息队列时生产者接受到的响应信息的 `topic` |

### 响应参数

|  参数  |             说明              |
|:----:|:---------------------------:|
| Code |   状态码可以在 `pkg/status` 下查看   |
| Msg  |        执行失败的信息都存在在这里        |
| Data | 执行成功响应的数据，例如：执行时间、输出内容、使用内存 |

## 致谢

在本项目的开发中，感谢 [criyle](https://github.com/criyle) 在项目开发中给我提供的大量帮助，没有他的帮助，本项目也不会实现。
