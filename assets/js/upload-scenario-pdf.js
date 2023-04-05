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
    formData.append('ignore', document.getElementById("ignore").value);

    try {
        const response = await fetch('/api/pdf-to-json', {
            method: 'POST',
            headers: {
                'Authorization': "Basic "+ document.getElementById("token").value,
            },
            body: formData,
        });

        if (response.ok) {
            const data = await response.json();
            const container = document.getElementById('container');
            data.forEach(item => {
                // 각 JSON 항목에 대해 div를 생성하고 데이터를 추가합니다.
                const div = document.createElement('div');
                div.textContent = `Page: ${item.pagenum}, Line: ${item.linenum}, Text: ${item.text}`;
            
                // 생성한 div를 컨테이너에 추가합니다.
                container.appendChild(div);
            });
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

