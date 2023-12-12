# typora图片上传工具

typora plugin，PicUploader

在学习go语言，刚好没有找到喜欢的typora图片上传插件，便自己边查边学边练，做出来了。练手项目，写的不好
将config.json文件放在执行文件同目录即可

nextcloud使用__File sharing__插件实现直链下载，使用开放api接口上传

## nextcloud图床配置文件说明，注意复制使用时请去掉注释

```json
{
  "picBed": "nextcloud",    //nextcloud图床
  "endpoint": "[host]/remote.php/dav/files",    //nextcloud的上传地址
  "bucketName": "",     //在File sharing插件中设置的存储路径
  "accessKeyId": "",     //nextcloud的账号
  "accessKeySecret": "",     //nextcloud的密码
  "downloadUrl": "[host]/apps/sharingpath",     //nextcloud通过File sharing插件产生的下载地址
  "useSSL": false
}
```

## 阿里云OSS对象存储图床配置文件说明，注意复制使用时请去掉注释

```json
{
    "picBed": "aliyunOss",     //阿里云OSS对象存储图床
    "endpoint": "oss-cn-hangzhou.aliyuncs.com",     //阿里云OSS endpoint
    "bucketName": "",     //阿里云OSS bucket名称
    "accessKeyId": "",     //阿里云accessKeyId
    "accessKeySecret": "",     //阿里云accessKeySecret
    "downloadUrl": "",     //阿里云不填写
    "useSSL": false
}
```

## minIO OSS对象存储图床配置文件说明，注意复制使用时请去掉注释

```json
{
    "picBed": "minIO",     //minIO OSS对象存储图床
    "endpoint": "play.min.io",     //minIO OSS endpoint，带上端口
    "bucketName": "",     //minIO OSS bucket名称
    "accessKeyId": "",     //minIO OSS accessKeyId
    "accessKeySecret": "",     //minIO OSS accessKeySecret
    "downloadUrl": "",     //minIO OSS 不填写
    "useSSL": false
}
```
