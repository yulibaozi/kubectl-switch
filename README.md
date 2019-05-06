# kubectl-switch

kubectl-switch是一个kubernetes命令行工具插件,它做的事情非常简单,即:你可以在一台计算机上可以操作多个kubernetes集群.

### 使用快照

这是一个小示例
 
### 安装

源码编译:

```
go get -u github.com/yulibaozi/kubectl-switch

cd $GOPATH/src/github.com/yulibaozi/kubectl-switch

go build .

# move PATH
mv  ./kubectl-switch /usr/local/bin/
```

Mac :

```
brew install yulibaozi/tap/kubectl-switch
```

### 注意事项

1. `kubectl` 版本要求:1.12.0或更高,你可以使用 `kubectl version` 来查看

### 使用流程

1. 当你安装完`kubectl-switch`后, 需要向`kubectl-switch`注册需要操作的集群
```
kubectl switch register      qa     /root/yulibaozi/admin.kubeconfig
#                         集群名字   集群操作所需要config文件
```

2. 查看向`kubectl-switch`注册成功的集群列表
```
kubectl switch list
```

3. 查看当前`kubectl`操作的集群
```
kubectl switch now
```
4. 切换`kubectl`操作的集群
```
kubectl switch qa
```

### 使用惯例

1. 查看支持哪些命令
```
kubectl switch -h
```



