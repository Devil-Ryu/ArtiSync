# README

## 简介

ArtiSync（多平台文章同步）是一个能批量将MarkDown文章，上传至各个平台的工具，主要解决多平台上传问题，以及不同平台的防盗链机制，导致文章图片存在三方网站时，该平台不能正确显示该文章图片的问题。本工具核心如下：
1. 分析Markdown文章；
2. 提取Markdown文章的图片；
3. 根据各个平台的上传图片接口，将图片上传至平台，并获取对应链接；
4. 根据平台替换文章中的本地链接；
5. 批量上传

文章发布流程如下：
![文章发布流程](images/ArtiSync_Publish.jpg)

## 上手指南

1. 根据操作系统安装[wails](https://wails.io/zh-Hans/docs/gettingstarted/installation/)
2. 克隆此仓库

```shell
git clone https://github.com/Devil-Ryu/ArtiSync.git
```
3. 运行项目
```shell
wails dev
```
4. 构建项目
```shell
wails build
```
本项目采用`wails`开发，逻辑使用GO语言，前端页面使用VUE3，前端模板采用TDesigner

## 使用手册
软件中用请遵循[软件许可使用协议](https://github.com/Devil-Ryu/ArtiSync/blob/beta/AGREEMENT.md)，软件使用可查看[使用手册](https://github.com/Devil-Ryu/ArtiSync/blob/beta/manual.md)

## 已完成功能
- [X] 修改接口增加请求头，请求体的方法
- [X] 编辑接口后退出确认框
- [X] 编辑接口后若未保存，退出后再进去就没有改动的痕迹
- [X] 文章管理功能，新增平台，删除平台测试
- [X] 接口URL可自定义路径
- [X] 获取参数类型时，应校验对应接口的Type，而不是本接口Type
- [X] 解决打开包含有Group接口等平台非常慢的问题（排查出问题为数据库返回body字段内容太多2000多条，引起卡顿，需注意保存是否会重复保存值(解决，更新平台状态会重复保存接口参数)
- [X] 增加文章上传状态
- [X] 跳转页面时候保存上个页面状态（文章页面
- [X] 添加系统配置页面
- [X] 用store进行全局组件共享
- [X] 添加日志中心
- [X] 增加接口运行结果查看
- [X] 解决BUG：文章控制器使用的DBController，和main.go 设置的不同，需要初始化
- [X] 接口运行后增加响应校验
- [X] BUG： 日志页面，查看详情看不到（直接修改为点击查看）
- [X] 功能：根据User目录，新建配置文件和数据库文件，或者默认目录

## 待办清单
- [ ] 功能：平台增加微信、知乎、博客园
- [ ] 功能：文章更新，添加接口步骤，到详情
- [ ] 功能：文章详情页：文章名称、图片数量、各平台图片上传数量、文章上传状态、文章上传进度、各平台图片上传状态、接口数量、各平台接口运行状态
- [ ] 优化：接口URL路径参数校验
- [ ] 优化：增加接口编辑规则校验
- [ ] 优化：接口编辑界面：请求体等可折叠
- [ ] 优化：接口信息过多时打开卡顿
- [ ] 优化：处理程序中的错误(正在做)
- [ ] 优化：添加接口方法禁用
- [ ] 优化：添加文章选择上传
- [ ] 优化：一键发布弹出确认框，确认发布的文章（或者做一个综合信息页，上面有，启用的平台，选择的文章，分析的图片列表，上传进度，可参考tdesign的详情页）
- [ ] 优化：界面美化
- [ ] 优化：内存优化，把暂存结果放到本地，节省内存空间（如日志，图片进度）

