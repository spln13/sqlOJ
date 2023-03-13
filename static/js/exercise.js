window.onload = () => {
    // 查看登录状态，获取用户名
    // 获取所有cookie
    let cookies = document.cookie;
// 将cookie字符串分割成一个数组
    let cookieArray = cookies.split("; ");
// 定义变量来存储token和username
    let token = "";
    let tokenExpires = "";
    let username = "";
    let usernameExpires = "";
// 遍历数组，查找名为"token"和"username"的cookie
    for (let i = 0; i < cookieArray.length; i++) {
        let cookie = cookieArray[i];
        let name = cookie.split("=")[0].trim();
        let value = cookie.split("=")[1];
        if (name === "token") {
            // 找到名为"token"的cookie，获取其值和过期时间
            token = value;
            tokenExpires = cookie.split("expires=")[1];
        } else if (name === "username") {
            // 找到名为"username"的cookie，获取其值和过期时间
            username = value;
            usernameExpires = cookie.split("expires=")[1];
        }
    }

// 判断token是否过期
    if (token && tokenExpires) {
        let now = new Date().getTime();
        let expires = new Date(tokenExpires).getTime();
        if (expires < now) {
            // token已过期，清除cookie
            document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
            token = "";
            tokenExpires = "";
        }
    }

// 判断username是否过期
    if (username && usernameExpires) {
        let now = new Date().getTime();
        let expires = new Date(usernameExpires).getTime();
        if (expires < now) {
            // username已过期，清除cookie
            document.cookie = "username=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
            username = "";
            usernameExpires = "";
        }
    }

// 输出结果
    console.log("Token: " + token);
    console.log("Token expires: " + tokenExpires);
    console.log("Username: " + username);
    console.log("Username expires: " + usernameExpires);

}