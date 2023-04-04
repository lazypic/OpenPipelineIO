


// Dropzone.js 설정
Dropzone.options.pdfUpload = {
    url: '/api/pdf-to-json',
    paramName: 'file', // 서버에 전송될 파일의 매개변수 이름
    maxFilesize: 10, // MB
    acceptedFiles: 'application/pdf',
    init: function () {
        this.on('sending', function (file, xhr) {
            console.log(document.getElementById("project").options[document.getElementById("project").selectedIndex].value)
            console.log(document.getElementById("version").value)
            console.log(document.getElementById("part").value)
            const formData = new FormData();
            // 프로젝트와 버전 정보를 전송 데이터에 추가
            formData.append('project', document.getElementById("project").options[document.getElementById("project").selectedIndex].value);
            formData.append('version', document.getElementById("version").value);
            formData.append('part', parseInt(document.getElementById("part").value));

             // Authorization 헤더를 요청에 추가
             xhr.setRequestHeader('Authorization', "Basic "+ document.getElementById("token").value);

             // 파일과 함께 폼 데이터 전송
             xhr.send(formData);
        });

        // 성공적으로 업로드된 경우, 응답 JSON 출력
        this.on('success', function (file, response) {
            console.log('업로드 성공: ', response);
        });

        // 업로드 실패한 경우, 에러 메시지 출력
        this.on('error', function (file, errorMessage) {
            console.error('업로드 실패: ', errorMessage);
        });
    },
};