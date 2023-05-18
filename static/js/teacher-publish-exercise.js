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

    const submit_buttonHTML = document.getElementById('submit-button');
    const nameHTML = document.getElementById('name');
    const answerHTML = document.getElementById('answer');
    const selectHTML = document.getElementById("select");
    const descriptionHTML = document.getElementById('description');
    const tableIDListHTML = document.getElementById('tables');
    submit_buttonHTML.addEventListener('click', function (e) {
        e.preventDefault();
        const grade = selectHTML.selectedIndex + 1;
        const name = nameHTML.value;
        const answer = answerHTML.value;
        const description = descriptionHTML.value;
        const tableIDStringList = tableIDListHTML.value;
        const stringArray = tableIDStringList.split(" ");
        const tableIDList = stringArray.map(str => parseInt(str, 10))
        const filteredTableIDArr = tableIDList.filter((value) => !isNaN(value));
        const dataToSent = {
            name: name,
            description: description,
            answer: answer,
            grade: grade,
            table_id_list: filteredTableIDArr,
        }
        console.log(dataToSent);
        const jsonData = JSON.stringify(dataToSent);
        fetch('/api/exercise/publish/exercise/', {
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
                    window.location.href = '/teacher/publish-exercise/';
                }
            })
            .catch(error => console.log(error))
    })
}

