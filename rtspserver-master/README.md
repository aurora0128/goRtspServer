# RTSPServer

#### 介绍
600行原生GO语言代码，打造史上最简单的RTSP服务器

#### 软件架构
600行代码，没有使用任何第三方的包，使用的都是GO自带的语法


#### 安装教程

1.  基于GO语言，不依赖任何框架，放心食用
2.  推流端命令如下：ffmpeg -i rtsp://admin:admin123@192.168.1.11:554//Streaming/Channels/1 -c:v copy -rtsp_transport tcp -f rtsp rtsp://192.168.1.37:5545/1f
3.  本服务器不支持UDP协议，在推流和播放时，请注意强制使用tcp通道。至于为什么不支持UDP协议，很简单，UDP容易花屏，究其原因，是因为单包太大，约1472个字节，根据我的经验，单包超过1400个字节，就容易丢掉部分数据。
#### 使用说明

1.  本项目仅支持直播，不支持历史文件存储和点播

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
