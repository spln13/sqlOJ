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


let createBox = (idx, exercise_id, exercise_name, grade, publisher_name, publisher_type, status, contest_id) => {
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
    // 竞赛中题目的url: /contest/123/problem/123
    const problem_url = '/contest/' + contest_id + '/problem/' + exercise_id;
    box.innerHTML = '<tr><td>' + idx + '</td><td><a href="' + problem_url + '">' + exercise_name + '</a></td><td class="' + grade_class + '">'
        + grade_type + '</td><td class="' + publisher_class + '">' + publisher_name + '</td><td class="' + status_class + '">' + status_content + '</td></tr>';
}

let parseTime = (time) => {
    const originalDate = new Date(time);
    const year = originalDate.getFullYear(); // 年份
    const month = originalDate.getMonth() + 1; // 月份（注意要加1，因为月份从0开始）
    const day = originalDate.getDate(); // 日期
    const hours = originalDate.getHours(); // 小时
    const minutes = originalDate.getMinutes(); // 分钟
    const seconds = originalDate.getSeconds(); // 秒钟
    return `${year}-${month.toString().padStart(2, '0')}-${day.toString().padStart(2, '0')} ${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
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
    const submission_button_url = '/contest/' + contestID + '/my-submission';
    let submission_button_HTML = document.getElementById('submission_button');
    submission_button_HTML.setAttribute('href', submission_button_url);
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
                timeHTML.innerHTML = '开始时间: ' + parseTime(begin_at) + '; 结束时间: ' + parseTime(end_at);
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
                    createBox(i + 1, exercise_id, exercise_name, grade, publisher_name, publisher_type, status, contestID);
                }
            }
        })
        .catch(error => console.error(error));

}

