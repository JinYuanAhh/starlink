>FirstConnectMessage

```json -- JSON
{
    "Type": "ConnectMsg",
}
```
>Signup

```json -- JSON
{
    "Type": "Signup",
    "Info": {
        "Account": "<Account>",
        "Pwd": "<Password>"
    }
}
```
>Signin

```json -- JSON
{
    "Type": "Signin",
    "Info": {
        "Account": "<Account>",
        "Pwd": "<Password>",
        "P": "Web/Mobile/PC/...."
    }
}
```
>Loggout

```json -- JSON
{
    "Type": "Logout",
    "T": "<token>"
}
```