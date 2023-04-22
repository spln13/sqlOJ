window.onload = index => {
    // 获取当前页面的路径
    const path = window.location.pathname;
    // 分割路径并获取最后一个部分
    const parts = path.split('/');
    const submissionID = parts[parts.length - 1];
    const url = '/api/submission/get/answer-detail?submission_id=' + submissionID;
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
                alert(status_msg)
            }
            else {
                const answer = data['answer'];
                let answerHTML = document.getElementById('answer');
                answerHTML.innerHTML = answer;
            }
        })
        .catch(error => console.error(error));

}

