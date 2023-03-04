// Example starter JavaScript for disabling form submissions if there are invalid fields
(() => {
  'use strict'

  // Fetch all the forms we want to apply custom Bootstrap validation styles to
  const forms = document.querySelectorAll('.needs-validation')

  // Loop over them and prevent submission
  Array.from(forms).forEach(form => {
    form.addEventListener('submit', event => {
      if (!form.checkValidity()) {
        event.preventDefault()
        event.stopPropagation()
      }

      form.classList.add('was-validated')
    }, false)
  })
})()

// isEmailValid 判断邮箱是否是BUCT邮箱 (判断@前是否为10位数字;@后是否位buct.edu.cn) 大小写不敏感
const isEmailValid = (email) => {
    const pattern = /^\d{10}@buct.edu.cn$/i
    return pattern.test(email)
}
// isUsernameValid 检查用户名是否符合规范, 允许汉字
const isUsernameValid = (username) => {
    const pattern = /^[\u4e00-\u9fa5a-zA-Z0-9_]{4,30}$/
    return pattern.test(username)
}

// isPasswordValid 检查密码是否符合规范 长度为8到32位 可以包含ASCII字符
const isPasswordValid = (password) => {
    const pattern = /^[\x20-\x7E]{8,30}$/;
    return pattern.test(password);
}

window.onload = () => {
    const btn_submit = document.querySelector("#btn_submit")
    const usernameStruct = document.querySelector("#username")
    const realNameStruct = document.querySelector("#realName")
    const numberStruct = document.querySelector("#number")
    const password1Struct = document.querySelector("#password1")
    const password2Struct = document.querySelector("#password2")
    const emailStruct = document.querySelector("#email")
    const codeStruct = document.querySelector("#code")
    const btn_sendCode = document.querySelector("#send_code")
    btn_sendCode.addEventListener('click', function (e) {
        e.preventDefault()
        const email = emailStruct.value;
        if (!isEmailValid(email)) {
            alert("请输入正确邮箱")
            return
        }
        const url = '/api/student/email/send-code/?email=' + email;
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
                if (status_code === 1) {
                    alert(status_msg)
                }
                else {
                    alert('发送成功')
                }
            })
            .catch(error => console.error(error));
    })
    btn_submit.addEventListener('click', function (e) {
        e.preventDefault()
        const username = usernameStruct.value
        const number = numberStruct.value
        const realName = realNameStruct.value
        const password1 = password1Struct.value
        const password2 = password2Struct.value
        const email = emailStruct.value
        const code = codeStruct.value
        if (!isPasswordValid(password1)) {
            alert('密码不符合规范')
            return
        }
        if (password1 !== password2) {
            alert('密码不一致')
            return
        }
        if (!isEmailValid(email)) {
            alert('邮箱格式错误')
            return
        }
        if (!isUsernameValid(username)) {
            alert('用户名格式错误')
            return
        }
        if (number !== email.substring(0, 10)) {
            alert('邮箱与学号不匹配')
            return
        }
        const formData = new FormData();
        formData.append('username', username)
        formData.append('number', number)
        formData.append('real_name', realName)
        formData.append('password', password1)
        formData.append('email', email)
        formData.append('code', code)
        fetch('/api/student/register/', {
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
                    alert("注册成功")
                    window.location.href = '/login/'
                }
            })
            .catch(error => console.log(error))
    })
}