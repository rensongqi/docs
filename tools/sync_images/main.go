package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/chyroc/lark"
	"github.com/chyroc/lark/card"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	OAGitAccessToken   = "DTLUxxxxxx"
	FeiShuAppID        = "cli_a228xxxxx"
	FeiShuAppSecret    = "xyRevQ1Z069gCxxxxx"
	AdministratorEmail = "songqi.ren@rsq.com"
)

type IssueInfo struct {
	Object  ObjectAttributes `json:"object_attributes" binding:"required"`
	Project ProjectInfo      `json:"project" binding:"required"`
	User    UserInfo         `json:"user"`
}

type UserInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ProjectInfo struct {
	Id int `json:"id" binding:"required"`
}

type ObjectAttributes struct {
	IId    int    `json:"iid" binding:"required"`
	Title  string `json:"title" binding:"required"`
	State  string `json:"state" binding:"required"` // 只对opened的事件做处理
	Action string `json:"action"`                   // closed 状态下更新了issue的action 是否处理待定
}

// CommentIssue 对 oagit issue进行评论
func (i *IssueInfo) CommentIssue(comment string) error {
	url := fmt.Sprintf("https://oagit.rsq.cn/api/v4/projects/%d/issues/%d/notes", i.Project.Id, i.Object.IId)
	method := "POST"

	// 创建请求体
	data := map[string]string{
		"body": comment,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Error("Error marshalling JSON:", err)
		return nil
	}

	resp, err := handleHttpRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	logrus.Info(string(resp))
	return nil
}

// CloseIssue 关闭oagit issue
func (i *IssueInfo) CloseIssue() error {
	url := fmt.Sprintf("https://oagit.rsq.cn/api/v4/projects/%d/issues/%d", i.Project.Id, i.Object.IId)
	method := "PUT"
	// 创建请求体
	data := map[string]string{
		"state_event": "close",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Error("Error marshalling JSON:", err)
		return err
	}

	resp, err := handleHttpRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	logrus.Info(string(resp))
	return nil
}

// GetUserInfo 获取oagit用户信息
func (i *IssueInfo) GetUserInfo() ([]*UserInfo, error) {
	url := fmt.Sprintf("https://oagit.rsq.cn/api/v4/users?username=%s", i.User.Username)
	method := "GET"
	resp, err := handleHttpRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	users := make([]*UserInfo, 0)

	if err = json.Unmarshal(resp, &users); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return users, nil
}

// handleHttpRequest 标准http请求
func handleHttpRequest(method, url string, body io.Reader) ([]byte, error) {
	client := &http.Client{}

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logrus.Error("Error creating http request:", err)
		return nil, err
	}
	// 设置 Basic 认证头
	req.Header.Add("PRIVATE-TOKEN", OAGitAccessToken)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	res, err := client.Do(req)
	if err != nil {
		logrus.Error("Error making request:", err)
		return nil, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Error("Error reading response body:", err)
		return nil, err
	}
	if strings.HasPrefix(fmt.Sprintf("%d", res.StatusCode), "2") {
		logrus.Infof("request %s api successfully.", url)
	} else {
		logrus.Errorf("Failed to request api: %s status: %d, %s\n", url, res.StatusCode, respBody)
	}
	return respBody, nil
}

func feiShuCardContent(image string, err error) string {
	var (
		describe string
		content  *lark.MessageContentCard
	)

	config := lark.MessageContentCardConfig{
		UpdateMulti: true,
	}
	if err == nil {
		describe = "【酷哇运维通知】您的镜像同步成功！！！"
		content = card.Card(card.Div().SetFields(card.FieldMarkdown(image))).SetHeader(card.Header(describe).SetGreen()).SetConfig(&config)
	} else {
		describe = "【酷哇运维通知】您的镜像同步失败！！！"
		content = card.Card(card.Div().SetFields(card.FieldMarkdown(image))).SetHeader(card.Header(describe).SetRed()).SetConfig(&config)
	}

	return content.String()
}

// sendCardToUserByEmail 通过邮箱给用户发送消息
func sendCardToUserByEmail(email, content string) error {
	cli := lark.New(
		lark.WithAppCredential(FeiShuAppID, FeiShuAppSecret),
	)

	ctx := context.Background()
	_, _, err := cli.Message.Send().ToEmail(email).SendCard(ctx, content)
	if err != nil {
		logrus.Error("[sendCardToUserByEmail] failed: ", err)
		return err
	}
	return nil
}

// 获取镜像proxy地址
func getProxyAddr(imageServer string) (dockerProxy string) {
	if strings.Contains(imageServer, ".") {
		switch imageServer {
		case "registry.k8s.io":
			dockerProxy = "k8s-gcr.libcuda.so"
		case "ghcr.io":
			dockerProxy = "ghcr.libcuda.so"
		case "k8s.gcr.io":
			dockerProxy = "k8s-gcr.libcuda.so"
		case "gcr.io":
			dockerProxy = "gcr.libcuda.so"
		case "docker.io":
			dockerProxy = "docker.libcuda.so"
		case "nvcr.io":
			dockerProxy = "ngc.nju.edu.cn"
		default:
			dockerProxy = "docker.libcuda.so"
		}
	}
	return
}

// 对镜像进行解析
func parseImageAndSync(images string) (err error, privateImagePath string) {
	if strings.HasPrefix(images, "nvcr.io") {
		return fmt.Errorf("nvidia镜像目前暂不支持，请等待后续优化"), ""
	}
	dockerProxy := "docker.lixd.xyz"
	imageServer := "docker.io"
	imageName := ""
	destRepo := ""
	command := []string{"--insecure-policy", "sync", "-a", "--keep-going", "--src", "docker", "--dest", "docker"}
	if !strings.Contains(images, ":") {
		images = images + ":latest"
	}
	if strings.Contains(images, "/") {
		imageSplit := strings.Split(images, "/")
		imagesContent := ""
		if len(imageSplit) >= 1 {
			firstSlashIndex := strings.Index(images, "/")
			if firstSlashIndex == -1 {
				return errors.New("镜像名有误，请检查镜像名配置是否符合规范"), ""
			}
			lastSlashIndex := strings.LastIndex(images, "/")
			if lastSlashIndex == -1 {
				return errors.New("镜像名有误，请检查镜像名配置是否符合规范"), ""
			}

			if lastSlashIndex > firstSlashIndex {
				imagesContent = images[firstSlashIndex+1 : lastSlashIndex]
			}

			if strings.Contains(imageSplit[0], ".") {
				imageServer = imageSplit[0]
				if imageServer == "k8s.gcr.io" {
					imageServer = "registry.k8s.io"
				}
				dockerProxy = getProxyAddr(imageServer)
				imageName = images[firstSlashIndex+1:]
			} else {
				imagesContent = images[:lastSlashIndex]
				imageName = images
			}
			proxyAdd := fmt.Sprintf("%s/%s", dockerProxy, imageName)
			destRepo = fmt.Sprintf("harbor.rsq.cn/%s/%s", imageServer, imagesContent)
			if strings.HasSuffix(destRepo, "/") {
				privateImagePath = fmt.Sprintf("%s%s", destRepo, images[lastSlashIndex+1:])
			} else {
				privateImagePath = fmt.Sprintf("%s/%s", destRepo, images[lastSlashIndex+1:])
			}
			command = append(command, proxyAdd, destRepo)
		}
	} else {
		proxyAdd := fmt.Sprintf("%s/%s", dockerProxy, images)
		destRepo = fmt.Sprintf("harbor.rsq.cn/docker.io")
		privateImagePath = fmt.Sprintf("%s/%s", destRepo, images)
		command = append(command, proxyAdd, destRepo)
	}

	logrus.Info("sync params: ", command)
	logrus.Info("private image path: ", privateImagePath)
	retries := 5
	var output []byte
	for i := 0; i < retries; i++ {
		cmd := exec.Command("skopeo", command...)
		output, err = cmd.CombinedOutput()
		if err == nil {
			logrus.Info("sync images success, output: ", string(output))
			break
		}
		logrus.Errorf("sync images failed, retry %d time, output: %v, err: %v", i+1, string(output), err)
		time.Sleep(time.Second * 6)
	}
	if err != nil {
		return fmt.Errorf("同步镜像 %s 失败，请检查提交镜像是否合规或联系运维人员处理, 错误信息: %v", images, string(output)), ""
	}

	return
}

// 解析oagit webhook的issue信息
func parseGitIssueInfos(c *gin.Context) {
	var issue IssueInfo
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}

	jsonData, _ := json.Marshal(&issue)
	logrus.Info(string(jsonData))

	users, err := issue.GetUserInfo()
	if err != nil {
		return
	}

	if issue.Object.State == "opened" {
		err, privateImagePath := parseImageAndSync(issue.Object.Title)
		if err != nil {
			// 给用户发送飞书失败消息
			for _, user := range users {
				if user.Name == issue.User.Name {
					_ = sendCardToUserByEmail(user.Email, feiShuCardContent(fmt.Sprintf("同步镜像 %s 失败，请检查提交镜像是否合规或联系运维人员处理", issue.Object.Title), err))
					if user.Email != AdministratorEmail {
						_ = sendCardToUserByEmail(AdministratorEmail, feiShuCardContent(fmt.Sprintf("同步镜像 %s 失败，请检查提交镜像是否合规或联系运维人员处理", issue.Object.Title), err))
					}
				}
			}

			_ = issue.CommentIssue(fmt.Sprintf("镜像同步失败，错误信息: %v", err))
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		comment := fmt.Sprintf("镜像同步完毕，内网访问地址: %s", privateImagePath)
		// 给用户发送飞书成功消息
		for _, user := range users {
			if user.Name == issue.User.Name {
				_ = sendCardToUserByEmail(user.Email, feiShuCardContent(comment, nil))
			}
		}
		if err = issue.CommentIssue(comment); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		if err = issue.CloseIssue(); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, issue)
}

func init() {
	// 设置json文本输出
	logFormat := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		PrettyPrint:     false,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := fmt.Sprintf("%v:%v", path.Base(frame.File), frame.Line)
			return frame.Function, fileName
		},
	}
	logrus.SetReportCaller(true)
	logrus.SetFormatter(logFormat)

	//同时写文件和屏幕
	stdOutWrite := io.Writer(os.Stdout)
	//fileAndStdoutWriter := io.MultiWriter(os.Stdout, logger)
	logrus.SetOutput(stdOutWrite)

	//设置最低loglevel
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("日志模块初始化完成")
}

func main() {
	r := gin.Default()
	{
		r.POST("/api/v1/issues", parseGitIssueInfos) // 解析oagit webhook的issue信息
	}
	_ = r.Run(":8888")
}
