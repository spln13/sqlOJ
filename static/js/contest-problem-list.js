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


let createBox = (exercise_id, exercise_name, grade, pass_count, submit_count, publisher_name, publisher_type, status) => {
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
    let pass_rate = pass_count / submit_count;
    pass_rate = pass_rate.toFixed(2);
    box.innerHTML = '<tr><td>' + exercise_id + '</td><td><a href="/problem/' + exercise_id + '">' + exercise_name + '</a></td><td class=' + grade_class + '>'
        + grade_type + '</td><td class=' + publisher_class + '>' + publisher_name + '</td><td>' + pass_rate
        + '</td><td class=' + status_class + '>' + pass_count + '</td></tr>';
}

window.onload = index => {
    // 查看登录状态，获取用户名
    // 获取所有cookie
    const username = getCookie("username");
    if (username !== "") {
        // 用户已登录，将用户名显示在页面右上角
        document.getElementById("button_username").innerHTML = '<div class="ui dropdown simple item">\n' +
            '      <div class="text">spln13</div>' +
            '      <i class="dropdown icon"></i>' +
            '      <div class="menu">' +
            '        <a class="item" href="/problem/status/?user=114980">提交记录</a>' +
            '        <a class="item" href="/account/settings/profile/">个人信息</a>' +
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

    const request_contest_info_url = '/api/'


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
            if (status_code === 1) {
                alert(status_msg)
            }
            else {
                const exercise_id = data['exercise_id'] // int
                const exercise_name = data['exercise_name'] // string
                const publisher_name = data['publisher_name'] // string
                const publisher_type = data['publisher_type'] // int
                const grade = data['grade'] // int
                const status = data['status'] // int
                const submit_count = data['submit_count'] // int
                const pass_count = data['pass_count'] // int
                createBox(exercise_id, exercise_name, grade, pass_count, submit_count, publisher_name, publisher_type, status);
            }
        })
        .catch(error => console.error(error));

}

