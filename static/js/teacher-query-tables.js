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

let createBox = (idx, table_id, table_name, description, association_count) => {
    let mother_box = document.querySelector("#tables");
    let box = document.createElement('tr');
    mother_box.appendChild(box);
    box.innerHTML = '<tr><td>' + idx + '</td><td>' + table_name + '</td><td>' + description + '</td>' +
        '<td>' + association_count + '</td></tr>';

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
            '        <a class="item" href="/teacher/migrate/">更改信息</a>' +
            '        <a class="item" href="/teacher/add/">增加教师</a>' +
            '        <a class="item" href="/logout/">登出</a>' +
            '      </div>' +
            '    </div>';
    }
    const url = '/api/get-type/';
    // 获取所有题目信息
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
            if (status_code !== 0) {    // token出错
                window.location = '/teacher/login/';
                return
            }
            const type = data['type'];
            if (type < 2) { // 学生
                window.location = '/teacher/login/';
            }
        })
        .catch(error => console.error(error));
    fetch('/api/exercise/get/all-tables/', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
    })
        .then(response => response.json())
        .then(data => {
            const status_code = data['status_code'];
            const status_msg = data['status_msg'];
            if (status_code !== 0) {    // token出错
                alert(status_msg);
                window.location = '/teacher/tables/';
            }
            const list = data['list'];
            for (let i = 0; i < list.length; i++) {
                const table_id = list[i]['table_id'];
                const table_name = list[i]['table_name'];
                const description = list[i]['description'];
                const association_count = list[i]['association_count'];
                createBox(i + 1, table_id, table_name, description, association_count);
            }
        })
        .catch(error => console.error(error));
}

