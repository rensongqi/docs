- [1 git分支操作](#1-git分支操作)
- [2 git submodule命令](#2-git-submodule命令)
- [3 git 创建之前版本的新分支](#3-git-创建之前版本的新分支)
- [4 git checkout 报错](#4-git-checkout-报错)
- [5 git cherry-pick](#5-git-cherry-pick)
- [6 git rebase](#6-git-rebase)
- [7 git 免秘钥pull push](#7-git-免秘钥pull-push)
- [8 git pull代码至指定目录下](#8-git-pull代码至指定目录下)
- [9 git 设置 proxy](#9-git-设置-proxy)
- [10 git 忽略 .idea文件](#10-git-忽略-idea文件)

# 1 git分支操作
```bash
# 1、拉取一个新项目
git clone http://xx.xx.xx/rsq/rsq.git

# 2、切换branch
git checkout rsq_test1

#3、创建一个新branch并push到git仓库
git checkout -b rsq_test2
git status
git add .
git push --set-upstream origin rsq_test2
```
# 2 git submodule命令
```bash
# 1、切换分支，pull最新的代码
git checkout rsq_test1
git branch
git pull

# 2、更改rsq/docs的submodule id
cd rsq/docs
git pull
git checkout f5b91959......
cd ..
git status 
git add docs
git commit -m "Refresh rsq/docs commit id to docs latest commit id f5b91959......"
git push
```

# 3 git 创建之前版本的新分支
```bash
git checkout rsq_test1
git pull

# 找到 rsq_test1 之前提交的 commit id
git checkout 7eb99a12....

# 在此分支之上创建新分支
git checkout -b rsq_test3
git status
git push --set-upstream origin rsq_test3
```

# 4 git checkout 报错
这种是本地有新的修改，解决办法就是提交修改或者删除文件
若是清理可按照如下操作
```bash
1、先强制切换到一个branch
git checkout -f master

2、将工作区、暂存区和HEAD保持一致
git reset --hard HEAD

3、获取最新代码
git pull

4、清除文件预览
git chean n

5、强制清除文件
git chean -f
```

# 5 git cherry-pick
用到cherry-pick的场景有很多，比如代码中写死了某个ip，由于不可抗拒因素，ip发生改变，所有branch的某一部分代码都要更改，这个时候就可以先找一个branch更改完毕，然后把这个commit cherry-pick到所有其它branch即可。
```
  a - b - c - d     Master
       \
         e - f - g  Release
```
如上所示，现在两个分支，要把Release分支中的f commit到Master分支，可按如下操作
```bash
# 切换到 master 分支
$ git checkout master

# Cherry pick 操作
$ git cherry-pick f
```
执行完的代码库变成如下
```
  a - b - c - d - f   Master
       \
         e - f - g    Release
```
git cherry-pick 后边可以跟多个commit id，形如
```bash
# 有两个branch A B
git cherry-pick A B

# A~B之间的所有提交全部cherry-pick 至当前分支，以下不包含A这个commit
# 需要注意的是A 的commit要早于B的commit，否则会失败但不会报错
git cherry-pick A..B 

# 包含A这个commit 写法如下
git cherry-pick A^..B 
```
>[git cherry-pick 教程](http://www.ruanyifeng.com/blog/2020/04/git-cherry-pick.html)

# 6 git rebase
**查看提交历史**
```bash
# git log --oneline --graph
* a8d5157 (HEAD -> release, origin/release) 
* 0194960 test1
* 2e8b3b3 test2
* e939b86 test3
```
当前分支在`a8d5157` ，此时若想回退到上一个commit：`0194960`，可用git rebase
```bash
# git rebase --onto 0194960
Successfully rebased and updated refs/heads/release.

# git log --oneline --graph
* 0194960 (HEAD -> release) delete pip install
* 2e8b3b3 add realtoedit copy scripts
* e939b86 change pip repo to douban.com
# git add .
# git commit -m "reset"
# git push --force
```

# 7 git 免秘钥pull push
```bash
git config --global user.email "xxx"
git config --global user.name “xxx”
git config --global credential.helper store
```

# 8 git pull代码至指定目录下
```bash
mkdir custom
cd custom
git init
git remote add -f origin http://git.rsq.local/rsq/test.git
git pull origin master
```

# 9 git 设置 proxy
```bash
# 设置http和https:
git config --global http.proxy http://127.0.0.1:7890
git config --global https.proxy https://127.0.0.1:7890

# 取消设置
git config --global --unset http.proxy
git config --global --unset https.proxy
```

# 10 git 忽略 .idea文件

> 若.idea没有被git跟踪，则直接在.gitignore文件中添加.idea
> 
> 若.idea已经被git跟踪，之后再加入.gitignore后是没有作用的，需要清除.idea的git缓存

```bash
git rm -r --cached .idea
```
然后在.gitignore中添加.idea