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

let createBox = (idx, status, user_type, username, submit_time, exercise_id, exercise_name, on_chain, contest_id, answer) => {
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
    let on_chain_content = '未上链';
    if (on_chain === 1) {
        on_chain_content = '已上链';
    }
    const user_class = "user_" + user_type;
    const problemURL = '/contest/' + contest_id + '/problem/' + exercise_id;
    let idxDOM = document.createElement('td');
    idxDOM.innerHTML = idx;
    box.append(idxDOM);
    let problemDOM = document.createElement('td');
    problemDOM.innerHTML = '<a href="' + problemURL + '">' + exercise_name + '</a>';
    box.append(problemDOM);
    let userDOM = document.createElement('td');
    userDOM.innerHTML = username;
    userDOM.setAttribute('class', user_class);
    box.append(userDOM);
    let submitTimeDOM = document.createElement('td');
    submitTimeDOM.innerHTML = submit_time;
    box.append(submitTimeDOM);
    let answerDOM = document.createElement('td');
    answerDOM.innerHTML = answer;
    box.append(answerDOM);
    let onChainDOM = document.createElement('td');
    onChainDOM.innerHTML = on_chain_content;
    box.append(onChainDOM);
    let statusDOM = document.createElement('td');
    statusDOM.innerHTML = status_content;
    statusDOM.setAttribute('class', status_class);
    box.append(statusDOM);
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
            '        <a class="item" href="/migrate/">更改信息</a>' +
            '        <a class="item" href="/logout/">登出</a>' +
            '      </div>' +
            '    </div>';
    }
    const url = '/api/get-type/';
    // 获取所有用户身份信息
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

    const path = window.location.pathname;
    // 分割路径并获取最后一个部分
    const parts = path.split('/');
    const contestID = parts[parts.length - 1];
    const getContestSubmissionURL = '/api/submission/contest/get-all?contest_id=' + contestID;
    fetch(getContestSubmissionURL, {
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
                const answer = list[i]['answer'];
                const status = list[i]['status'];
                const user_type = list[i]['user_type'];
                const username = list[i]['username'];
                const submit_time = list[i]['submit_time'];
                const exercise_id = list[i]['exercise_id'];
                const exercise_name = list[i]['exercise_name'];
                const on_chain = list[i]['on_chain'];
                createBox(i + 1, status, user_type, username, submit_time, exercise_id, exercise_name, on_chain, contestID, answer);
            }
        })
        .catch(error => console.error(error));
}
