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


let createBox = (idx, number, realName, exerciseScore, contestScore, score) => {
    let mother_box = document.querySelector("#students");
    let box = document.createElement('tr');
    box.innerHTML = '<td>' + idx + '</td><td>' + number + '</td><td>' + realName + '</td><td>' + exerciseScore + '</td>' +
        '<td>' + contestScore + '</td><td>' + score + '</td>';
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
            const type = data['type']
            if (type < 2) { // 学生
                window.location = '/teacher/login/';
            }
        })
        .catch(error => console.error(error));
    const computingButton = document.getElementById('computing')
    const downloadButton = document.getElementById('download')
    downloadButton.addEventListener('click', function (e) {
        e.preventDefault();
        window.location.href = '/download/score';
    })
    computingButton.addEventListener('click', function (e) {
        e.preventDefault();
        const requestURL = '/api/student/get/score/';
        fetch(requestURL, {
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
                    return
                }
                const list = data['list'];
                for (let i = 0; i < list.length; i++) {
                    const number = list[i]['number'];
                    const realName = list[i]['real_name'];
                    const exerciseScore = list[i]['exercise_score'];
                    const contestScore = list[i]['contest_score'];
                    const score = list[i]['score'];
                    createBox(i + 1, number, realName, exerciseScore, contestScore, score);
                }
            })
            .catch(error => console.error(error));
    })
}

