const getCookie = (cname) => {
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


let createBox = (id, submission_id,exercise_id, exercise_name, submit_time, status) => {
    let mother_box = document.querySelector("#submission");
    let box = document.createElement('tr');
    mother_box.appendChild(box);
    const status_class = 'status_' + status;
    let status_content = 'AC';
    if (status === 2) {
        status_content = 'WA';
    }
    else if (status === 3) {
        status_content = 'RE'
    }

    box.innerHTML = '<tr><td>' + id + '</td><td><a href="/problem/' + exercise_id + '">' + exercise_name + '</a></td><td>' + submit_time +
        '</td><td><a href="/submission/' + submission_id + '">查看</a></td>' + '<td class="' + status_class + '">' + status_content + '</td></tr>';
}

window.onload = () => {
    // 查看登 录状态，获取用户名
    // 获取所有cookie
    const username = getCookie("username");
    if (username !== "") {
        // 用户已登录，将用户名显示在页面右上角
        document.getElementById("button_username").innerHTML = '<div class="ui dropdown simple item">\n' +
            '      <div class="text">spln13</div>' +
            '      <i class="dropdown icon"></i>' +
            '      <div class="menu">' +
            '        <a class="item" href="/submission/">提交记录</a>' +
            '        <a class="item" href="/profile/">个人信息</a>' +
            '        <a class="item" href="/migrate/">更改信息</a>' +
            '        <a class="item" href="/logout/">登出</a>' +
            '      </div>' +
            '    </div>';
    }
    const url = '/api/submission/get/one-all/'
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
            if (status_code !== 0) {
                alert(status_msg);
                return
            }
            const list = data['list'];
            for (let i = 0; i < list.length; i++) {
                const submission_id = list[i]['submission_id'];
                const exercise_id = list[i]['exercise_id'];
                const exercise_name = list[i]['exercise_name'];
                const submit_time = list[i]['submit_time'];
                const status = list[i]['status'];
                createBox(i + 1, submission_id, exercise_id, exercise_name, submit_time, status);
            }
        })
        .catch(error => console.error(error));
}

