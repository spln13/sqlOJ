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


let createBox = (exercise_id, exercise_name, grade, publisher_name, publisher_type, status) => {
    let mother_box = document.querySelector("#exercises");
    let box = document.createElement('tr');
    mother_box.appendChild(box);
    let grade_type;
    if (grade === 1) {
        grade_type = "easy";
    }
    else if (grade === 2) {
        grade_type = "medium";
    }
    else {
        grade_type = "hard";
    }
    const publisher_class = "publisher_" + publisher_type;
    const grade_class = "grade_" + grade_type;
    const status_class = "status_" + status;
    let status_content = '未提交';
    if (status === 1) {
        status_content = 'AC';
    }
    else if (status === 2) {
        status_content = 'WA';
    }
    else if (status === 3) {
        status_content = 'RE';
    }
    box.innerHTML = '<tr><td>' + exercise_id + '</td><td><a href="/problem/' + exercise_id + '">' + exercise_name + '</a></td><td class="' + grade_class + '">'
        + grade_type + '</td><td class="' + publisher_class + '">' + publisher_name + '</td><td class="' + status_class + '">' + status_content + '</td></tr>';
}

window.onload = index => {
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
    // 获取当前页面的路径 当前url: /contest/123
    const path = window.location.pathname;
    // 分割路径并获取最后一个部分
    const parts = path.split('/');
    const contestID = parts[parts.length - 1];
    const request_contest_info_url = '/api/contest/get/contest?contest_id=' + contestID;
    fetch(request_contest_info_url, {
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
                alert(status_msg)
            }
            else {
                const contest_id = data['contest_id'];
                const contest_name = data['contest_name'];
                const begin_at = data['begin_at'];
                const end_at = data['end_at'];
                let titleHTML = document.getElementById('title');
                titleHTML.innerHTML = contest_id.toString() + '. ' + contest_name;
                let timeHTML = document.getElementById('time');
                timeHTML.innerHTML = '开始时间: ' + begin_at + '; 结束时间: ' + end_at;
            }
        })
        .catch(error => console.error(error));


    const url = '/api/contest/get/all-exercise?contest_id=' + contestID;
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
                alert(status_msg)
                window.location = '/contest';
            }
            else {
                const list = data['list'];
                for (let i = 0; i < list.length; i++) {
                    const exercise_id = list[i]['exercise_id'] // int
                    const exercise_name = list[i]['exercise_name'] // string
                    const publisher_name = list[i]['publisher_name'] // string
                    const publisher_type = list[i]['publisher_type'] // int
                    const grade = list[i]['grade'] // int
                    const status = list[i]['status'] // int
                    const submit_count = list[i]['submit_count'] // int
                    const pass_count = list[i]['pass_count'] // int
                    createBox(i + 1, exercise_name, grade, publisher_name, publisher_type, status);
                }
            }
        })
        .catch(error => console.error(error));

}

