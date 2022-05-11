# 文件上传

## 格式

> **Args|Content**  
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

### Content
?> **可选**  

如果不为空，则会自动追加一次[Append](/fileUpload#append-传输未完成上传的文件内容)操作，若发生错误，则错误会以```APPPEND-```开头，之后错误内容同[Append错误](/fileUpload#error-错误-1)

### Error - 错误

> #### exist

为Sha **< Sha >** 的文件已存在

## Append - 传输未完成上传的文件内容

### Info  

```json
{
    "Type":"Append",
    "Sha":"文件最终内容的Sha256"
}
```

### Content

文件真正要追加的内容，无需编码

### Error - 错误

> #### not your file or no such file

没有**你上传的**Sha256为 **< Sha >** 的文件

> #### completed

此文件已经完成上传

## Complete - 完成文件上传

### Info  

```json
{
    "Type":"Complete",
    "Sha":"文件最终内容的Sha256"
}
```

### Error - 错误

> #### not your file or no such file

没有**你上传的**Sha256为 **< Sha >** 的文件

> #### completed

此文件已经完成上传

> #### sha256 not match

此文件现有内容的Sha256与[New](/fileUpload#new-新建文件上传请求)时刻传入的Sha256不匹配
