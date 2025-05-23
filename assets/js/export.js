// changeReportExcelURI 함수는 프로젝트를 선택하면 reportexcel url을 설정한다.
function changeReportExcelURI() {
    let project = CurrentProject();
    document.getElementById("reportexcelURI").href = "/reportexcel?project=" + project;
}

// changeReportJSONURI 함수는 프로젝트를 선택하면 reportjson url을 설정한다.
function changeReportJSONURI() {
    let project = CurrentProject();
    document.getElementById("reportJSONURI").href = "/reportjson?project=" + project;
}

// include, exclude modal에서 드레그앤 드롭을 할 때 사용하는 코드
Sortable.create(include, {
    group: {
        name: 'include',
        put: ['exclude1','exclude2','exclude3'],
    },
    animation: 100
});

Sortable.create(exclude1, {
    group: {
        name: 'exclude1',
        put: ['include','exclude2','exclude3'],
    },
    animation: 100
});

Sortable.create(exclude2, {
    group: {
        name: 'exclude2',
        put: ['include','exclude1','exclude3'],
    },
    animation: 100
});

Sortable.create(exclude3, {
    group: {
        name: 'exclude3',
        put: ['include','exclude1','exclude2'],
    },
    animation: 100
});

// changeExportFormatType 함수는 ExportExcel의 포멧 옵션이 변경될 때 실행되는 함수이다.
function changeExportFormatType() {
    let token = document.getElementById("token").value;
    let format = document.getElementById("format").value;
    if (format === "all") {
        let sel = document.getElementById('task');
        sel.innerHTML = "";
        let allOpt = document.createElement('option');
        allOpt.appendChild(document.createTextNode("All"));
        allOpt.value = "all";
        sel.appendChild(allOpt);
    }
    if (format === "shot") {
        $.ajax({
            url: `/api/shottasksetting`,
            type: "get",
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let tasks = data.tasksettings
                // task를 갱신한다.
                let sel = document.getElementById('task');
                sel.innerHTML = "";
                // all 추가
                let allOpt = document.createElement('option');
                allOpt.appendChild(document.createTextNode("all"));
                allOpt.value = "all";
                sel.appendChild(allOpt);
                // shot Task 추가
                for (let i = 0; i < tasks.length; i++) {
                    let opt = document.createElement('option');
                    opt.appendChild( document.createTextNode(tasks[i].name));
                    opt.value = tasks[i].name;
                    sel.appendChild(opt);
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
    if (format === "asset") {
        $.ajax({
            url: `/api/assettasksetting`,
            type: "get",
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let tasks = data.tasksettings
                // task를 갱신한다.
                let sel = document.getElementById('task');
                sel.innerHTML = "";
                // all 추가
                let allOpt = document.createElement('option');
                allOpt.appendChild(document.createTextNode("all"));
                allOpt.value = "all";
                sel.appendChild(allOpt);
                // shot Task 추가
                for (let i = 0; i < tasks.length; i++) {
                    let opt = document.createElement('option');
                    opt.appendChild( document.createTextNode(tasks[i].name));
                    opt.value = tasks[i].name;
                    sel.appendChild(opt);
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}


// exportExcelCurrentPage는 현재 페이지를 엑셀로 뽑는다.
function exportExcelCurrentPage() {
    let project = document.getElementById("searchbox-project").value
    let task = document.getElementById("searchbox-task").value
    let searchword = document.getElementById("searchbox-searchword").value
    let sortkey = document.getElementById("searchbox-sortkey").value
    let sortorder = document.getElementById("searchbox-sortorder").value
    let truestatusList = []
    
    let checkStatus = document.querySelectorAll('*[id^="searchbox-checkbox-"]');
    for (i=0;i<checkStatus.length;i++) {
        if (checkStatus[i].checked) {
            truestatusList.push(checkStatus[i].getAttribute("status"))
        }
    }
    truestatus = truestatusList.join(",")
    
    // 요청
    let url = `/download-excel-file?project=${project}&task=${task}&searchword=${searchword}&sortkey=${sortkey}&sortorder=${sortorder}&truestatus=${truestatus}`
    location.href = url
}

// exportJsonCurrentPage는 현재 페이지를 .json 으로 뽑는다.
function exportJsonCurrentPage() {
    let project = document.getElementById("searchbox-project").value
    let task = document.getElementById("searchbox-task").value
    let searchword = document.getElementById("searchbox-searchword").value
    let sortkey = document.getElementById("searchbox-sortkey").value
    let truestatusList = []
    let truestatus = "" // ver1 검색바 때문에 이 값이 필요하다.
    
    let checkStatus = document.querySelectorAll('*[id^="searchbox-checkbox-"]');
    for (i=0;i<checkStatus.length;i++) {
        if (checkStatus[i].checked) {
            truestatusList.push(checkStatus[i].getAttribute("status"))
        }
    }
    truestatus = truestatusList.join(",")
    
    // 요청
    let url = `/download-json-file?project=${project}&task=${task}&searchword=${searchword}&sortkey=${sortkey}&truestatus=${truestatus}`
    location.href = url
}

// exportCsvCurrentPage는 현재 페이지를 .csv 으로 뽑는다.
function exportCsvCurrentPage() {
    let project = document.getElementById("searchbox-project").value
    let task = document.getElementById("searchbox-task").value
    let searchword = document.getElementById("searchbox-searchword").value
    let sortkey = document.getElementById("searchbox-sortkey").value
    let sortorder = document.getElementById("searchbox-sortorder").value
    let truestatusList = []
    let truestatus = ""
    let checkStatus = document.querySelectorAll('*[id^="searchbox-checkbox-"]');
    for (i=0;i<checkStatus.length;i++) {
        if (checkStatus[i].checked) {
            truestatusList.push(checkStatus[i].getAttribute("status"))
        }
    }
    truestatus = truestatusList.join(",")
    
    // 사용자가 CSV를 뽑기 위해서 include 영역에 드레그한 아이템 가지고 오기.
    let ul = document.getElementById("include")
    let li = ul.getElementsByTagName("li")
    let titles = []
    for (let i = 0; i <= li.length -1; i++) {
        titles.push(li[i].getAttribute("id"))
    }
    let titleString = titles.join(",")
    // 요청
    let url = `/download-csv-file?project=${project}&task=${task}&searchword=${searchword}&sortkey=${sortkey}&sortorder=${sortorder}&truestatus=${truestatus}&titles=${titleString}`
    location.href = url
}

