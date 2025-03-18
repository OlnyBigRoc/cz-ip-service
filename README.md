# cz-ip-service

#### 介绍
纯真ip社区库作为基础数据，提供接口请求服务  

#### 密钥获取
来这里 [https://www.cz88.com/geo-public](https://www.cz88.com/geo-public) 获取开发者key和fileKey和secretKey
1. 开发者key
   ![dev_key.png](img/dev_key.png)
2. fileKey
   > 点击复制我们会得到一个下载链接，这里我们只需要连接中的key部分数据即可 
   > https://www.cz88.net/api/communityIpAuthorization/communityIpDbFile?fn=czdb&key=1234567890  
   > 1234567890 就是我们需要的 fileKey

   ![file_key.png](img/file_key.png)

3. secretKey
   ![secret_key.png](img/secret_key.png)
#### 软件架构

软件架构说明


#### 安装教程

1. 下载依赖
    ``` shell
    go mod tidy
    ```
2. 构建
    ``` shell
    go build main.go
    ```
3. 启动
    ``` shell
    ./main -developerKey=developerKey -fileKey=fileKey -secretKey=secretKey
    ```
4. 访问
   ``` shell 
   curl http://127.0.0.1/json?ip=1.1.1.1
   curl http://127.0.0.1
   ```

#### 参与贡献

1. Fork 本仓库
2. 新建 Feat_xxx 分支
3. 提交代码
4. 新建 Pull Request

