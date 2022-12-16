# Gin OAuth 2.0 Server

> 第三方应用在判断当前用户未登录的情况下跳转至，其中client_id由服务端分配
```
http://localhost:9096/authorize?client_id=000000&response_type=code&redirect_uri=http://localhost
```
>在用户正常登录的情况下浏览器会重定向至  
http://localhost/?code=R2OIB8GIPFIWKMEQWFJZIG  
第三方应用请求以下地址获取token（建议使用POST方式）

```
http://localhost:9096/token?grant_type=authorization_code&client_id=000000&client_secret=999999&scope=read&code=R2OIB8GIPFIWKMEQWFJZIG&redirect_uri=http://localhost
```

正常情况下会返回
```json
{
    "access_token":"ICDSZEVHM_IYSBNUKU9KJW",
    "expires_in":7200,
    "refresh_token":"KOP_T8GWVB2KLFHP5TE9DW",
    "token_type":"Bearer"
}

```

访问 http://localhost:9096/user 获取用户信息，在请求的header信息中应包含
```
Authorization: Bearer 0CRMNYLAMEA3AFJKTOQPVQ
```

access_token过期时可以通过refresh_token获取新的token

```
http://localhost:9096/token?client_id=000000&client_secret=999999&grant_type=refresh_token&refresh_token=OOCFVGZFWVG9S6K9VETVSQ
```