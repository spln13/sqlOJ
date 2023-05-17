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

// isPasswordValid 检查密码是否符合规范 长度为8到32位 可以包含ASCII字符
const isPasswordValid = (password) => {
    const pattern = /^[\x20-\x7E]{8,30}$/;
    return pattern.test(password);
}

// isUsernameValid 检查用户名是否符合规范, 允许汉字
const isUsernameValid = (username) => {
    const pattern = /^[\u4e00-\u9fa5a-zA-Z0-9_]{4,30}$/
    return pattern.test(username)
}

window.onload = () => {
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

    const submit_buttonHTML = document.getElementById('submit-button');
    const usernameHTML = document.getElementById('username');
    const realNameHTML = document.getElementById('real-name');
    const passwordHTML = document.getElementById("password");
    const passwordAgainHTML = document.getElementById('password-again');
    submit_buttonHTML.addEventListener('click', function (e) {
        e.preventDefault();
        const username = usernameHTML.value;
        const realName = realNameHTML.value;
        const password = passwordHTML.value;
        const passwordAgain = passwordAgainHTML.value;

        if (!isUsernameValid(username)) {
            alert('用户名不合法');
            return;
        }
        if (password !== passwordAgain) {
            alert('密码不相等');
            return;
        }
        if (!isPasswordValid(password)) {
            alert('密码不合法');
            return;
        }


        const url = '/api/teacher/add?username=' + username + '&password=' + password + '&real_name=' + realName;
        fetch(url, {
            method: 'POST',
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
                }
                else {
                    alert("添加成功");
                    window.location.href = '/teacher/upload-table/';
                }
            })
            .catch(error => console.log(error))
    })
}