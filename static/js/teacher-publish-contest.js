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
    const submitButton = document.getElementById('submit-button');
    const nameHTML = document.getElementById('name');
    const beginAtHTML = document.getElementById('begin_at');
    const endAtHTML = document.getElementById('end_at');
    const problemIDList = document.getElementById('problems');
    const classIDList = document.getElementById('classes');
    submitButton.addEventListener('click', function (e) {
        e.preventDefault();
        const name = nameHTML.value;
        const beginAtValue = beginAtHTML.value;
        const endAtValue = endAtHTML.value;
        const beginAtDate = new Date(beginAtValue);
        const endAtDate = new Date(endAtValue);
        const beginAt = beginAtDate.toISOString();
        const endAt = endAtDate.toISOString();
        const problemIDStringList = problemIDList.value;
        const classIDStringList = classIDList.value;
        const problemIDStringArray = problemIDStringList.split(" ");
        const classIDStringArray = classIDStringList.split(" ");
        const problemIntIDList = problemIDStringArray.map(str => parseInt(str, 10));
        const classIntIDList = classIDStringArray.map(str => parseInt(str, 10));
        const filteredProblemIDArr = problemIntIDList.filter((value) => !isNaN(value));
        const filteredClassIDArr = classIntIDList.filter((value) => !isNaN(value));
        const dataToSend = {
            contest_name: name,
            begin_at: beginAt,
            end_at: endAt,
            exercise_id_list: filteredProblemIDArr,
            class_id_list: filteredClassIDArr
        }
        console.log(dataToSend);
        const jsonData = JSON.stringify(dataToSend);
        fetch('/api/contest/create/', {
            method: 'POST',
            body: jsonData
        })
            .then(response => response.json())
            .then(data => {
                const status_code = data['status_code'];
                const status_msg = data['status_msg'];
                if (status_code !== 0) {
                    alert(status_msg);
                }
                else {
                    alert("发布成功");
                    window.location.href = '/teacher/publish-contest/';
                }
            })
            .catch(error => console.log(error))
    })
}

