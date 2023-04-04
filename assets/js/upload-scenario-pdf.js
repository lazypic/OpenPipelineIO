


// 수동으로 업로드 하는 방법을 찾아야 겠다.
const dropZone = document.getElementById('drop-zone');

// 파일 드롭 처리 함수
async function handleFileDrop(event) {
    event.preventDefault();

    const file = event.dataTransfer.files[0];
    if (!file) return;

    // FormData 생성 및 파일과 프로젝트 정보 추가
    const formData = new FormData();
    formData.append('file', file);
    formData.append('project', document.getElementById("project").options[document.getElementById("project").selectedIndex].value);
    formData.append('version', document.getElementById("version").value);
    formData.append('part', parseInt(document.getElementById("part").value));

    try {
        const response = await fetch('/api/pdf-to-json', {
            method: 'POST',
            headers: {
                'Authorization': "Basic "+ document.getElementById("token").value,
            },
            body: formData,
        });

        if (response.ok) {
            const jsonResponse = await response.json();
            console.log('전송 성공:', jsonResponse);
        } else {
            console.error('전송 실패:', response.statusText);
        }
    } catch (error) {
        console.error('전송 중 오류 발생:', error);
    }
}

// 드래그 앤 드롭 이벤트 처리 함수
function handleDragOver(event) {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'copy';
}

// 이벤트 리스너 등록
dropZone.addEventListener('drop', handleFileDrop);
dropZone.addEventListener('dragover', handleDragOver);

