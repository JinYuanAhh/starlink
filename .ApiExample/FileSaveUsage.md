>数据格式：
#### Arg|Content  
### Arg:[最终需要base64编码]  
```json -- JSON
{
    "Type": "File",
    "Status": "New/Continue",
    "MD5": "<MD5>",
    "CompSegIndex": "<CompSegIndex>",
    "T": "<token>",
}
```
MD5: 完整文件MD5，文件写入完成后对比用  
CompSegIndex: 文件总共有几段  
#### New:新建文件，不需要Content，但是需要字符串结尾加入|
#### Continue:加入文件内容，content即内容
### Content:[最终需要base64编码]  
#最终：
> base64过的Arg + | + base64过的Content