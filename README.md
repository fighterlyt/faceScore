
# 概述
使用阿里云公开的[颜值计算API](https://market.aliyun.com/products/57124001/cmapi014897.html?spm=5176.8216963.738025.22.9wrCe2#sku=yuncode889700000)

## 处理步骤		
 1. faceScore.NewScorer(appCode) 生成一个Scorer,appCode就是购买该API所提供的验证码		
 2. faceScore.LocalScore(本地图片地址)  / Scorer.Webscore(web图片地址)
 
## API 解析
阿里云公开的API提供两种方式
*   本地图片计算  
*   网络图片计算

### 本地图片计算
要素为
*   POST
*   http://faceapi.remarkdip.com/user/faceScore
*   将base64编码的图片作为Body内容的**image_base64字段值**提交

*   Authorization 值为APPCODE空格+用户的appcode
*   Content-Type 为**application/x-www-form-urlencoded; charset=UTF-8**

### 网路图片计算

*   GET
*   http://faceapi.remarkdip.com/user/faceScore
*   将图片地址作为QUERY内容的**iimg_url字段值**提交

*   Authorization 值为APPCODE空格+用户的appcode


### 返回值

#### 正常

```json
{
  "code": "success",
  "message": "success",
  "result": 93.2
}
```
#### 失败

```json
{
  "code": "fail",
  "message": "parameter 'img_base64' is missing"
}
```
