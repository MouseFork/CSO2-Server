## Counter-Strike Online 2 Server 

### 一、介绍

CSOL2 服务器 v0.2.0

*用于韩版 Counter-Strike Online 2

目前客户端请使用L-Leite的启动器。

这是我的第一个Go语言项目，用来练习，不知道会不会弄下去，参考了L-Leite的数据结构。

如果大家有什么建议或问题，欢迎提出。

![Image](./photos/main.png)

![Image](./photos/intro.png)

![Image](./photos/channel.png)

![Image](./photos/ingame.jpg)

![Image](./photos/result.jpg)

### 二、项目计划

    1.先实现基本的游戏游玩功能和联机功能 √
    2.重构代码 ...(进行中)
    3.添加局域网频道和局域网搜索联机功能
    4.制作启动器，同时与服务器连接获取反馈，例如登录失败、当前用户已登录等
    5.添加仓库管理器，管理自己的武器

### 三、基本已完成的功能

    登录、频道、房间、仓库、UDP、角色战绩(游戏结果界面)、数据库、个人信息

### 四、正在编写的功能

    玩家积分、聊天

### 五、已知问题

    1.可能存在多协程共享变量不安全的问题
    2.由于房间用户与主管理器的用户重复，可能造成性能浪费
    3.房间ID和房间NUM在多频道下可能冲突（虽然目前是单频道）
    4.房主离开后，其余玩家会卡住直到炸出房间

### 六、已修复的问题

    1.主机开始游戏后，其他玩家不能加入，显示超时。需要和主机一起开始游戏才能加入。
    2.房间列表显示的房间信息及状态不准确，待刷新
    3.玩家仓库数据不准确
    4.房间列表的房主名字显示中文乱码
    5.每局结束显示的角色战绩可能存在一些问题
    6.房间无法加密码

### 七、使用方法

1.需要有韩版的CSOL2，同时使用第三方启动器，[点击这里下载](https://pan.baidu.com/s/13wEMinbj6E2Z9lds20NU3A) 提取码：picf

2.进入本项目的release页面下载最新版本的程序（ https://github.com/KouKouChan/CSO2-Server/releases ）

3.建立bat文件，和游戏的bin目录同级，里面写入：

```shell
START ./bin/launcher.exe -masterip IP地址 -enablecustom -username 用户名 -password 密码
```

4.IP地址指的是你的服务端IP，如果是本地那么就填192.168开头的你的IP地址，如果你要连接局域网别人的服务端那么就填别人的IP地址，如果你安装了汉化包，也可以再加上以下语句：

```shell
-lang schinese
```

5.先运行本项目的exe文件启动服务器，然后打开bat文件启动客户端即可

注意：从网盘里面下载得到的start-cso2.bat文件需要修改里面的IP地址和用户名！

### 八、自定义文件方法

1.下载CSOL2解包工具，[点击这里下载](https://pan.baidu.com/s/14q1SoIdHwp1casMWG2OS-w) 提取码：41bs

2.解压后，打开工具，点击左上角File选项，点击Open folder，选中csol2的data文件夹即可

3.解压你需要的文件，并且将解压后的文件按你的想法进行修改

4.将文件放入csol2目录的custom文件夹下，打开游戏

### 九、Docker下使用方法

1.首先你需要拥有Docker,请下载并安装Docker,同时配置好Docker,比如Docker源

2.输入以下命令拉取最新版的服务端:

```shell
docker pull koukouchan/cso2-server:latest
```

3.运行服务端

```shell
docker run -p 30001:30001 -p 30002:30002 koukouchan/cso2-server:latest
```

4.接下来打开客户端，连接服务器

### 十、编译环境

*Go 1.14.2*

当你要架设局域网或外网时，请打开防火墙的端口。30001-TCP类型端口、30002-UDP类型端口

貌似建立互联网服务器需要双方玩家都能内网穿透，实测局域网能够连接，互联网无法房间内加入主机，可能需要架设虚拟局域网。

### 十一、编译方法

```shell
1. 进入目录
2. 执行命令 go build
3. 运行生成的可执行文件即可
```

### 十二、Docker下编译方法

1.首先你需要拥有Docker,请下载并安装Docker,同时配置好Docker,比如Docker源

2.在终端下进入项目目录，输入以下命令:

```shell
docker build -t cso2-server .
```

3.在第2步后，如果运行正常，会显示所有步骤都运行完毕。接下来是运行服务端，为了能够让游戏和Docker容器里面的服务端相连，你需要打开相应的端口映射，使用以下命令运行：

```shell
docker run -p 30001:30001 -p 30002:30002 cso2-server
```

4.接下来打开客户端，连接服务器

*声明：Counter-Strike Online 2 归NEXON所有*
