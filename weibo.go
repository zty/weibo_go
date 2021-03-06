package main

import (
"bytes"
"encoding/json"
"fmt"
"io"
"io/ioutil"
"mime/multipart"
"net/http"
"net/url"
// "strconv"
"os"
"strings"
)

type AccessToken struct {
Access_token string
Expires_in int
Uid string
}

type UnreadCount struct {
Status int
Follower int
Cmt int
Dm int
Mention_status int
Mention_cmt int
Group int
Notice int
Invite int
Badge int
Photo int
}

func (a AccessToken) String() string {
return a.Access_token
}

var (
access_token string = ""
uid string = ""
config_path string = os.Getenv("HOME") + "/.clifaweibo_oauth_2"
)

const (
client_id string = "2369288261"
client_secret string = "8aa208e6b2df21e5a3e65326d3568d9b"
grant_type string = "authorization_code"
redirect_uri string = "http://clifaweibo.sinaapp.com/ilovecliofubuntu.php"
access_token_url string = "https://api.weibo.com/oauth2/access_token"
statuses_update_url string = "https://api.weibo.com/2/statuses/update.json"
statuses_upload_url string = "https://upload.api.weibo.com/2/statuses/upload.json"
unread_count_url string = "https://rm.api.weibo.com/2/remind/unread_count.json"
)

func get_access_token_from_file() bool {
config_file, err := os.OpenFile(config_path, os.O_RDONLY, 0644)
if err != nil {
return false
}
defer config_file.Close()
a := new(AccessToken)
err = json.NewDecoder(config_file).Decode(a)
if err != nil {
fmt.Println(err.Error())
return false
}
access_token = a.Access_token
uid = a.Uid
return true
}

func get_access_token_from_http() bool {
fmt.Printf("请复制下面的地址在浏览器中打开完成授权\nhttp://t.cn/zlUnqCC\n请输入(右键粘贴)授权码:")
var grant_code string
fmt.Scanf("%s", &grant_code)
grant_url_values := make(url.Values)
grant_url_values.Set("client_id", client_id)
grant_url_values.Set("client_secret", client_secret)
grant_url_values.Set("grant_type", grant_type)
grant_url_values.Set("code", grant_code)
grant_url_values.Set("redirect_uri", redirect_uri)
r, _ := http.Post(access_token_url, "application/x-www-form-urlencoded", strings.NewReader(grant_url_values.Encode()))
body, _ := ioutil.ReadAll(r.Body)
fmt.Println(r,r.Body,body)
fmt.Println(string(body))
defer r.Body.Close()
a := new(AccessToken)
err := json.Unmarshal(body,&a)
if err != nil {
fmt.Println(err.Error())
return false
}
access_token = a.Access_token
uid = a.Uid
config_file, err := os.OpenFile(config_path, os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
return false
}
defer config_file.Close()
config_file.Write(body)
return true
}

func send_text_weibo(text string) bool {
weibo_values := make(url.Values)
weibo_values.Set("access_token", access_token)
weibo_values.Set("status", text)
r, err := http.Post(statuses_update_url, "application/x-www-form-urlencoded", strings.NewReader(weibo_values.Encode()))
defer r.Body.Close()
if err != nil {
return false
}
return true
}
func send_pic_weibo(text, file_path string) bool {
var b bytes.Buffer
body_writer := multipart.NewWriter(&b)
body_writer.WriteField("access_token", access_token)
body_writer.WriteField("status", text)

pic_writer, _ := body_writer.CreateFormFile("pic", file_path)
fh, _ := os.Open(file_path)
io.Copy(pic_writer, fh)
form_type := body_writer.FormDataContentType()
body_writer.Close()
resp, err := http.Post(statuses_upload_url, form_type, &b)
defer resp.Body.Close()
if err != nil {
return false
}
return true
}
func get_unread_count() (u *UnreadCount, err error) {
get_url := unread_count_url + "?access_token=" + access_token + "&uid=" + uid
r, err := http.Get(get_url)
defer r.Body.Close()
if err != nil {
fmt.Println(err.Error())
return u, err
}
u = new(UnreadCount)
err = json.NewDecoder(r.Body).Decode(u)
if err != nil {
fmt.Println(err.Error())
return u, err
}
return u, nil
}

func send_weibo(argc int) {
switch os.Args[1] {
case "-m", "m":
u, err := get_unread_count()
if err != nil {
fmt.Println(err.Error())
return
}
fmt.Println(u.Status)
case "-t", "t":
text := ""
for i := 2; i < argc; i++ {
text += os.Args[i]
if i < argc-1 {
text += " "
}
}
if true == send_text_weibo(text) {
fmt.Println("发送成功")
}
case "-p", "p":
if true == send_pic_weibo("分享图片", os.Args[argc-1]) {
fmt.Println("发送成功")
}
case "-tp", "-pt", "pt", "tp":
text := ""
for i := 2; i < argc-1; i++ {
text += os.Args[i]
if i < argc-1 {
text += " "
}
}
if true == send_pic_weibo(text, os.Args[argc-1]) {
fmt.Println("发送成功")
}
default:
return
}
}

func show_unread_count(show_from_num int) {
u, _ := get_unread_count()
um := map[int]string{
// u.Status : "新微博未读",
u.Follower: "新粉丝",
u.Cmt: "新评论",
u.Dm: "新私信",
u.Mention_status: "新提及我的微博",
u.Mention_cmt: "新提及我的评论",
u.Group: "微群消息未读",
u.Notice: "新通知未读",
u.Invite: "新邀请未读",
u.Badge: "新勋章",
u.Photo: "相册消息未读",
}
if show_from_num == 0 {
um[u.Status] = "新微博未读"
}
for k, v := range um {
if k != 0 {
fmt.Printf("%s:%d\n", v, k)
}
}
}

func main() {
if false == get_access_token_from_file() {
if false == get_access_token_from_http() {
return
}
}
argc := len(os.Args)
switch argc {
case 1:
return
case 2:
if os.Args[1] == "m" || os.Args[1] == "-m" {
show_unread_count(0)
}
if os.Args[1] == "a" || os.Args[1] == "-a" {
show_unread_count(1)
}
return
default:
send_weibo(argc)
}
}
