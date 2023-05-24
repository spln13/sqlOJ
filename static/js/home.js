getCookie = (cname) => {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) === 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

let createBox = (username, userType, score) => {
    const user_class = "user_" + userType;
    let mother_box = document.getElementById('ranking');
    let box = document.createElement('tr');
    box.innerHTML = '<td class="' + user_class + '">' + username + '</td><td>' + score + '</td>';
    mother_box.append(box);
}


window.onload = () => {
    // 查看登录状态，获取用户名
    // 获取所有cookie
    const username = getCookie("username");
    if (username !== "") {
        // 用户已登录，将用户名显示在页面右上角
        document.getElementById("button_username").innerHTML = '<div class="ui dropdown simple item">\n' +
            '      <div class="text">' + username + '</div>' +
            '      <i class="dropdown icon"></i>' +
            '      <div class="menu">' +
            '        <a class="item" href="/submission/">提交记录</a>' +
            '        <a class="item" href="/profile/">个人信息</a>' +
            '        <a class="item" href="/migrate/">更改信息</a>' +
            '        <a class="item" href="/logout/">登出</a>' +
            '      </div>' +
            '    </div>';
    }
    const url = '/api/ranking/get/min/';
    fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
    })
        .then(response => response.json())
        .then(data => {
            const status_code = data['status_code'];
            const status_msg = data['status_msg'];
            if (status_code === 1) {
                alert(status_msg)
                return
            }
            const list = data['list'];
            for (let i = 0; i < list.length; i++) {
                const username = list[i]['username'];
                const userType = list[i]['user_type'];
                const score = list[i]['score'];
                createBox(username, userType, score);
            }
        })
        .catch(error => console.error(error));
}

