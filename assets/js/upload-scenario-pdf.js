// 수동으로 업로드 하는 방법을 찾아야 겠다.
const dropZone = document.getElementById('drop-zone');

function clearDivContent(divId) {
    const div = document.getElementById(divId);
    div.innerHTML = '';
}

// 파일 드롭 처리 함수
async function handleFileDrop(event) {
    

    const container = document.getElementById('container');
    // 기존에 렌더링한 Text를 제거한다.
    clearDivContent('container');

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


async function fetchProjects() {
    // 인증 토큰 값을 설정합니다. 실제 토큰 값으로 교체해주세요.
    const authToken = document.getElementById("token").value;
  
    // 헤더에 인증 토큰을 추가합니다.
    const headers = new Headers({
      'Authorization': `Basic ${authToken}`,
      'Content-Type': 'application/json'
    });
  
    // REST API를 사용하여 프로젝트 목록을 가져옵니다. 헤더를 추가하여 인증합니다.
    const response = await fetch('/api2/projects', { headers });
    const projects = await response.json();
  
    // 프로젝트 목록을 사용하여 select 엘레먼트에 옵션을 추가합니다.
    const projectSelect = document.getElementById('project');
  
    projects.forEach(project => {
      // 각 프로젝트에 대해 option 엘레먼트를 생성하고 데이터를 추가합니다.
      const option = document.createElement('option');
      option.value = project.id;
      option.textContent = project.id;
  
      // 생성한 option을 select 엘레먼트에 추가합니다.
      projectSelect.appendChild(option);
    });
  }
  
  // 프로젝트 목록을 가져오고 select 엘레먼트에 옵션을 추가하는 함수를 실행합니다.
  fetchProjects();
  