const isEmail = (str) => {
    const regex = /^[\w-]+(\.[\w-]+)*@[\w-]+(\.[\w-]+)+$/;
    return regex.test(str);
  }


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

        const code = isEmail(username) ? 2 : 1; // 1是用户名登录；2是密码登录
        const url = '/api/student/login/?username_email=' + username + '&password=' + password + '&code=' + code.toString();
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
                    expires.setDate(expires.getDate() + 7); // 设置token过期时间为7天
                    document.cookie = `token=${token}; expires=${expires.toUTCString()}`;   // 设置token
                    alert("成功")
                }
            })
            .catch(error => console.error(error));

    })
}