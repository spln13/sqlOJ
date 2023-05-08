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

createBox = (exercise_id, exercise_name, publisher_id, publisher_name, publisher_type, grade, answer, description, submit_count, pass_count) => {
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
    const exercise_url = '/problem/' + exercise_id;
    const answer_url = '/teacher/exercise-answer/' + exercise_id;
    box.innerHTML = '<tr><td>' + exercise_id + '</td><td><a href="' + exercise_url + '">' + exercise_name + '</a></td>' +
        '<td class="' + publisher_class + '">' + publisher_name + '</td><td class="' + grade_class + '">' + grade_type + '</td>' +
        '<td>' + submit_count + '</td><td>' + pass_count + '</td><td><a href="' + answer_url + '">查看</a>';

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
            '        <a class="item" href="/logout/">登出</a>' +
            '      </div>' +
            '    </div>';
    }
    const url = '/api/get-type/';
    // 获取身份信息
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

    // 获取题库信息
    const get_exercises_url = '/api/exercise/teacher/all-exercises/';
    fetch(get_exercises_url, {
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
                const exercise_id = list[i]['exercise_id'];
                const exercise_name = list[i]['exercise_name'];
                const publisher_id = list[i]['publisher_id'];
                const publisher_name = list[i]['publisher_name'];
                const publisher_type = list[i]['publisher_type'];
                const grade = list[i]['grade'];
                const answer = list[i]['answer'];
                const description = list[i]['description'];
                const submit_count = list[i]['submit_count'];
                const pass_count = list[i]['pass_count'];
                createBox(exercise_id, exercise_name, publisher_id, publisher_name, publisher_type, grade, answer, description, submit_count, pass_count);
            }
        })
        .catch(error => console.error(error));

}

