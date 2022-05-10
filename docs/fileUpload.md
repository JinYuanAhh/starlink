# 文件上传

## 格式

> **{ Args }|{ Content }**  
> **"|"** 不可省略，无论Content是否为空

### Args
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

!> 以上JSON经由**BASE64**编码后得到Args
### Content
在 **Append** 时使用，为要上传的数据

## 返回

```json
{
    "Type":"File",
    "Info": {
        "Type": "< 同你的Info.Type >",
        "Status": "是否成功 < Success / Error >",
        "Error": "仅当Status为Error时存在，为错误信息"
    }
}
```

## New - 新建文件上传请求

### Info  

```json
{
    "Type":"New",
    "Filename":"文件名",
    "Sha":"文件最终内容的Sha256",
    "Size":"文件完整大小 (KB)",
    "Args":"其他附加信息（可空，需要什么信息就添加） [建议使用json文本]"
}
```

## Append - 传输未完成上传的文件内容

### Info  

```json
{
    "Type":"New",
    "Sha":"文件最终内容的Sha256"
}
```

### Content

文件真正要追加的内容，无需编码

## Complete - 完成文件上传

### Info  

```json
{
    "Type":"New",
    "Sha":"文件最终内容的Sha256"
}
```