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

const modifySubmissionButton = (exercise_id) => {
    let submissionButton = document.getElementById('submission-button');
    submissionButton.innerHTML = '<a class="item" href="/exercise/my-submission/' + exercise_id + '">我的提交</a>';
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
    const getTypeUrl = '/api/get-type/';
    // 获取所有用户身份信息
    fetch(getTypeUrl, {
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
                window.location = '/login/';
                return
            }
            const type = data['type']
            if (type < 1) {
                window.location = '/login/';
            }
        })
        .catch(error => console.error(error));
    // 获取当前页面的路径
    const path = window.location.pathname;
    // 分割路径并获取最后一个部分
    const parts = path.split('/');
    const exerciseID = parts[parts.length - 1];
    modifySubmissionButton(exerciseID)
    const url = '/api/exercise/get/one?exercise_id=' + exerciseID;
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
                const name = data['name'];  // string
                const grade = data['grade'] // int
                const description = data['description'] // string
                const publisher_name = data['publisher_name'] // string
                const publisher_type = data['publisher_type'] // int
                const submit_count = data['submit_count'] // int
                const pass_count = data['pass_count'] // int
                let grade_str;
                if (grade === 1) {
                    grade_str = 'easy';
                }
                else if (grade === 2) {
                    grade_str = 'medium';
                }
                else {
                    grade_str = 'hard';
                }

                const grade_class = 'grade_' + grade_str;
                const publisher_class = 'publisher_' + publisher_type.toString();
                document.getElementById('card_top').innerHTML = '<p><b>' + pass_count.toString() + '份提交通过</b>, 共有' + submit_count.toString() + '份提交。</p>' +
                    '<p><b>难度</b>: <b class="' + grade_class + '">' + grade_str + '</b>。</p>';
                document.getElementById('card_bottom').innerHTML = '<p><b>出题人</b>: <b class="' + publisher_class + '">' + publisher_name + "</b>。</p>"
                document.getElementById('title').innerHTML = exerciseID + '. ' + name;
                document.getElementById('content').innerHTML = marked.parse(description);
            }
        })
        .catch(error => console.error(error));
    const submitButton = document.getElementById("submit_button")
    // const sqlInput = document.getElementById("sql-input")
    const sqlEditor = CodeMirror.fromTextArea(document.getElementById("sql-input"), {
        mode: "text/x-mysql",
        lineNumbers: true
    });
    submitButton.addEventListener("click", function (e) {
        e.preventDefault();
        const sqlInputValue = sqlEditor.getValue();
        console.log("input value: ", sqlInputValue);
        const formData = new FormData;
        formData.append("exercise_id", exerciseID);
        formData.append("answer", sqlInputValue);
        fetch('/api/exercise/submit/', {
            method: 'POST',
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                const status_code = data['status_code'];
                const status_msg = data['status_msg'];
                if (status_code !== 0) {
                    alert(status_msg)
                }
                else {
                    // alert("提交成功")
                    window.location.href = '/exercise/my-submission/' + exerciseID
                }
            })
            .catch(error => console.log(error))
    })
}

