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


let createBox = (idx, exercise_id, exercise_name, status, submit_time, contest_id, submission_id, on_chain) => {
    let mother_box = document.querySelector("#submission");
    let box = document.createElement('tr');
    mother_box.appendChild(box);
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
    let on_chain_content = '未上链';
    if (on_chain === 1) {
        on_chain_content = '已上链';
    }
    // 竞赛中题目的url: /contest/123/problem/123
    const problem_url = '/contest/' + contest_id + '/problem/' + exercise_id;
    const contest_submission_url = '/contest/submission-detail/' + submission_id;
    box.innerHTML = '<tr><td>' + idx + '</td><td><a href="' + problem_url + '">' + exercise_name + '</a></td><td>' + submit_time + '</td>' +
    '<td><a href="' + contest_submission_url + '">查看' + '</a></td><td>' + on_chain_content + '</td>'
     +   '<td class="' + status_class + '">' + status_content + '</td>';
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
    // 获取当前页面的路径 当前url: /contest/123/my-submission
    const path = window.location.pathname;
    // 分割路径并获取最后一个部分
    const parts = path.split('/');
    const contestID = parts[parts.length - 2];
    const problem_list_url = '/contest/' + contestID;
    let problem_list_button = document.getElementById('problem-list');
    problem_list_button.setAttribute('href', problem_list_url);

    const request_contest_submission_url = '/api/submission/contest/get-my?contest_id=' + contestID;

    fetch(request_contest_submission_url, {
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
                    const submission_id = list[i]['submission_id']
                    const exercise_id = list[i]['exercise_id'] // int
                    const exercise_name = list[i]['exercise_name'] // string
                    const answer = list[i]['answer'] // string
                    const status = list[i]['status'] // int
                    const submit_time = list[i]['submit_time'] // int
                    const on_chain = list[i]['on_chain']
                    createBox(i + 1, exercise_id, exercise_name, status, parseTime(submit_time), contestID, submission_id, on_chain);
                }
            }
        })
        .catch(error => console.error(error));

}

