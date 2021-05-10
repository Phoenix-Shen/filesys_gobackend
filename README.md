# filesys_gobackend
写GO后端还不如睡大觉

# OSS的地址
可以在conf/app.conf 中的<br>
accesskeyid = 填入你的ID<br>
accesskeysecret = 填入你的密码<br>
endpoint = oss-cn-hangzhou.aliyuncs.com<br>
这些东西里头修改
# 换代理
git config --global http.proxy http://127.0.0.1:8889 <br>
git config --global --unset http.proxy

# 使用方法
在控制台中进入当前目录，输入go  run  main.go <br>
或者使用bee工具，在控制台中输入bee run即可

# 关于生成“昂首阔步地走”的API接口文档
运行项目之后 在浏览器中输入 localhost:8080/swagger/# 就可以访问辣！


