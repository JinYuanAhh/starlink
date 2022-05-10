# 文件上传

## 格式

> **{ args }|{ content }**  
> **"|"** 不可省略，无论content是否为空

### args
```json
{
    "Type": "File",
    "Info": {
        ……
    }
}
```

> 关于Info：  
> 不同消息的Info不同，将会在下面分段说明

!> 以上JSON经由**BASE64**编码后得到args
### content
在 **Append** 时使用，为要上传的数据

## New - 新建文件上传请求

### Info:  
```json
{
    "Type":"New",
    "Filename":"文件名",
    "Sha":"文件最终内容的Sha256",
    "Size":"文件完整大小 (KB)",
    "Args":"其他附加信息（可空，需要什么信息就添加） [建议使用json文本]"
}
```