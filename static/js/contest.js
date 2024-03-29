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
let parseTime = (time) => {
    const originalDate = new Date(time);
    const year = originalDate.getFullYear(); // 年份
    const month = originalDate.getMonth() + 1; // 月份（注意要加1，因为月份从0开始）
    const day = originalDate.getDate(); // 日期
    const hours = originalDate.getHours(); // 小时
    const minutes = originalDate.getMinutes(); // 分钟
    const seconds = originalDate.getSeconds(); // 秒钟
    return `${year}-${month.toString().padStart(2, '0')}-${day.toString().padStart(2, '0')} ${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
}

const createBox = (idx, contest_id, contest_name, publisher_name, publisher_type, begin_at, end_at) => {
    let mother_box = document.querySelector('#contests');
    let box = document.createElement('tr');
    const publisher_class = "publisher_" + publisher_type;
    box.innerHTML = '<tr><td>' + idx + '</td><td><a href="/contest/' + contest_id.toString() + '">' + contest_name + '</a></td>' +
        '<td class="' + publisher_class + '">' + publisher_name + '</td><td>' + begin_at +
        '</td><td>' + end_at + '</td></tr>';
    mother_box.append(box);
}
window.onload = () => {
    // 查看登 录状态，获取用户名
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
    // 获取所有竞赛信息
    const url = '/api/contest/get/all/'
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
            if (status_code !== 0) {
                alert(status_msg);
                return
            }
            const list = data['list'];
            for (let i = 0; i < list.length; i++) {
                const contest_id = list[i]['contest_id'];
                const contest_name = list[i]['contest_name'];
                const publisher_name = list[i]['publisher_name']
                const publisher_type = list[i]['publisher_type'] // 根据不同的发布者类型渲染不同颜色
                const begin_at = list[i]['begin_at']
                const end_at = list[i]['end_at']
                createBox(i + 1, contest_id, contest_name, publisher_name, publisher_type, parseTime(begin_at), parseTime(end_at));
            }
        })
        .catch(error => console.error(error));
}
