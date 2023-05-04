let createRecords = (statusList) => {
    let mother_box = document.querySelector("#records");
    for (let i = 0; i < statusList.length; i++) {
        let line = document.createElement('tr');
        mother_box.appendChild(line);
        const number = statusList[i][0];
        let numberCell = document.createElement('td');
        numberCell.innerHTML = number;
        line.append(numberCell);
        for (let j = 1; j < statusList[i].length; j++) {
             let cell = document.createElement('td');
             const status_class = 'status_' + statusList[i][j];
             let status = '未提交';
             if (statusList[i][j] === '1') {
                 status = 'AC';
             }
             else if (statusList[i][j] === '2') {
                 status = 'WA';
             }
             else if (statusList[i][j] === '3') {
                 status = 'RE';
             }
             cell.innerHTML = status;
             cell.setAttribute('class', status_class);
             line.append(cell);
        }
    }
}

let createFirstRow = (problemIDList, colCount, contestID) => {
    let table_head = document.querySelector('#first_row');
    const oneWideClass = 'center aligned one wide';
    for (let i = 0; i < colCount; i++) {
        const problemURL = '/contest/' + contestID + '/problem/' + problemIDList[i];
        let box = document.createElement('th');
        const idx = (i + 1).toString();
        box.innerHTML = '<th class="' + oneWideClass + '"><a href="' + problemURL + '">问题' + idx + '</a></th>';
        table_head.append(box);
    }
}

window.onload = () => {
    const getTypeUrl = '/api/get-type/';
    // 获取所有用户身份信息
    fetch(getTypeUrl, {
        method: 'GET',
        headers:  {
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

    // 获取竞赛id
    const path = window.location.pathname;
    // 分割路径并获取最后一个部分
    const parts = path.split('/');
    const contestID = parts[parts.length - 1];
    const getStatusUrl = '/api/contest/status?contest_id=' + contestID;
    fetch(getStatusUrl, {
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
                window.location = '/teacher/contests/';
                return
            }
            const problemIDList = data['problem_list'];
            const statusList = data['status_list']; // [][]string
            const colCount = problemIDList.length;
            createFirstRow(problemIDList, colCount, contestID);
            createRecords(statusList);
        })
        .catch(error => console.error(error));
}

