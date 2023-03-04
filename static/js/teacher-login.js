window.onload = () => {
    const btn_submit = document.querySelector('#btn_submit');
    const usernameStruct = document.querySelector('#username');
    const passwordStruct = document.querySelector('#password');
    btn_submit.addEventListener('click', function (e) {
        e.preventDefault();
        const username = usernameStruct.value;
        const password = passwordStruct.value;
        if (username === '' || password === '') {
            alert("请正确输入信息")
            return;
        }
        const url = 'http://127.0.0.1:8080/api/teacher/login/?username=' + username + '&password=' + password;
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
                const token = data['token']
                if (status_code === 1) {
                    alert(status_msg)
                }
                else {
                    const expires = new Date();
                    expires.setDate(expires.getDate() + 3); // teacher设置token过期时间为3天
                    document.cookie = `token=${token}; expires=${expires.toUTCString()}`;   // 设置token
                    alert("成功")
                }
            })
            .catch(error => console.error(error));

    })
}