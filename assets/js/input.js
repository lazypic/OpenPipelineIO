
var OSName="Linux";
if (navigator.appVersion.indexOf("Win") != -1) OSName="Win";
if (navigator.appVersion.indexOf("Mac") != -1) OSName="Mac";
if (navigator.appVersion.indexOf("X11") != -1) OSName="Linux";
if (navigator.appVersion.indexOf("Linux") != -1) OSName="Linux";

// changeProject 함수는 검색바의 프로젝트가 바뀔 때 발생하는 이벤트를 처리한다.
function changeProject() {
	// 제목설정
	let appname = document.getElementById("appname").innerText +": "
	let e = document.getElementById("searchbox-project");
	let project = e.options[e.selectedIndex].value;
	if (project != "") {
		document.title = appname + project;
	}
}

// selectmodeV2 함수는 검색바에서 "선택모드" 버튼을 눌렀을 때 실행되는 함수이다.
function selectmodeV2() {
	var checknum = 0
	// 체크된 Status 수를 구한다.
	for(var i=0; i<document.getElementsByClassName('StatusCheckBox').length;i++) {
		if (document.getElementsByClassName('StatusCheckBox')[i].checked == true) {
			checknum += 1
		}
	}
	// 다 켜져있을 때는 다 끈다.
	if (checknum == document.getElementsByClassName('StatusCheckBox').length) {
		for(var i=0; i<document.getElementsByClassName('StatusCheckBox').length;i++) {
			document.getElementsByClassName('StatusCheckBox')[i].checked=false
		}
	} else if (checknum == 0) {
		// 다 꺼져 있을 때는 몇개만 켠다.
		for(var i=0; i<document.getElementsByClassName('StatusCheckBox').length;i++) {
			document.getElementsByClassName('StatusCheckBox')[i].checked=false
		}
		for(var i=0; i<document.getElementsByClassName('DefaultStatusCheckBox').length;i++) {
			document.getElementsByClassName('DefaultStatusCheckBox')[i].checked=true
		}
	} else {
		// 일부만 켜있다면 다 켠다.
		for(var i=0; i<document.getElementsByClassName('StatusCheckBox').length;i++) {
			document.getElementsByClassName('StatusCheckBox')[i].checked=true
		}
	}
}


//샷 체크박스 F5 눌렀을때 없애는 기능.
function uncheck() {
	var checkboxes = document.getElementsByName('select_slug');
	for (var i=0; i<checkboxes.length; i++) {
		console.log(checkboxes[i].type)
		if (checkboxes[i].type == 'checkbox') {
			checkboxes[i].checked = false;
		}
	}
}


function playmov(address)
{
	myWindow=window.open(address,"PlayWindows","width=1280, height=720, toolbar=no, menubar=no, location=no");
	myWindow.focus();
}

function keypressed(){
	if(event.keyCode==122) self.close();
	else return false;
}
document.omkeydown=keypressed;


function onlyNumber(event) {
	event = event || window.event;
	var keyID = (event.which) ? event.which : event.keyCode;
	if ( (keyID >=48 && keyID <= 57) || (keyID >= 96 && keyID <= 105) || keyID == 8 || keyID == 46 || keyID == 37 || keyID == 39)
		return;
	else
		return false;
}

function removeChar(event) {
	event = event || window.event;
	var keyID = (event.which) ? event.which : event.keyCode;
	if (keyID == 8 || keyID == 46 || keyID == 37 || keyID == 39)
		return;
	else
		event.target.value = event.target.value.replace(/[^0-9]/g,"");
}


function removeWhiteSpace(event) {
	event.value = event.value.replace(/ /g, '');
}

// *문자를 x문자로 바꾼다.
// X를 x문자로 바꾼다.
// 공백을 제거한다.
// 렌즈디스토션값을 입력시 2048*1280 -> 2048x1280 형태로 바꾸기 위함이다.
// 숫자와 x를 제외한 영문입력시 삭제됩니다.
function widthxHeight(event) {
	event = event || window.event;
	event.target.value = event.target.value.replace("*","x");
	event.target.value = event.target.value.replace("X","x");
	event.target.value = event.target.value.replace(/[^\d\x]/gi,"");
}

//drop된 file의 "file://" 제거
function rmFileProtocol(event) {	
	event = event || window.event;
	event.preventDefault();
	
	var data= event.dataTransfer.getData('text/plain'); //드래그한 데이터 자료를 얻는다.
	event.target.value = "";
	event.target.value = data.replace("file://","");
}

//버튼을 클릭 하면 editItem 언디스토션사이즈 form에 placesize(scansize) 값을 입력한다.
//Checkbox를 체크하면 undistort value에 platesize가 입력된다.
//Checkbox의 체크를 해제하면 undistort value가 ""이 된다.
function inputNone(checkbox) {
	if (checkbox.checked) {
		document.getElementById("undistort").value = document.getElementById("platesize").value;
	} else {
		document.getElementById("undistort").value = "";
	}
}

// ScreenX 버튼이 클릭될때 체크 여부에 따라 이벤트가 발생한다.
// ScreenX가 체크되면 ScreenxOverlay가 활성화 된다.
// ScreenX가체크가 해제되면 ScreenxOverlay가 비활성화되고 1.0의 값이 들어간다.
function checkScreenx(event) {
	event = event || window.event;
	if (event.target.checked) {
		document.getElementById("screenxoverlay").readOnly = false;
	} else {
		document.getElementById("screenxoverlay").readOnly = true;
		document.getElementById("screenxoverlay").value = 1.0;
	}
}


function wfs(host, task, type, assettype, project, name, seq, cut, token) {
	let WFSPATH = "";
	$.ajax({
		url: "/api/tasksetting",
		type: "POST",
		data: {
			os: "", // os를 설정하지 않으면 WFSPath를 불러온다.
			task: task,
			type: type,
			assettype: assettype,
			project: project,
			name: name,
			seq: seq,
			cut: cut,
		},
		headers: {
			"Authorization": "Basic "+ token
		},
		async:false,
		dataType: "json",
		success: function(data) {
			WFSPATH = data.path
		},
		error: function(request,status,error){
			alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
		}
	});
	window.open(host + WFSPATH, 'WebFileSystem', 'left=20, top=20, width=500, height=500, toolbar=0, resizable=1, scrollbars=1').focus();
}
	

/* RESTAPI Template
fetch('/api/addassettag', {
    method: 'POST',
    headers: {
        "Authorization": "Basic "+ document.getElementById("token").value,
    },
    body: new URLSearchParams({
        id: id,
        assettag: assettag,
    })
})
.then((response) => {
    if (!response.ok) {
        throw Error(response.statusText + " - " + response.url);
    }
    return response.json()
})
.then((data) => {
    // process
})
.catch((err) => {
    alert(err)
});
*/

var SELECT_COLOR = "rgb(255, 196, 35)" // 선택된 색상
var NON_SELECT_COLOR = "rgb(167, 165, 157)" // 기본 색상

// modal이 뜨면 오토포커스가 되어야 한다.
if (document.getElementById("modal-addcomment") !== null) {
    $('#modal-addcomment').on('shown.bs.modal', function () {
        $('#modal-addcomment-text').trigger('focus')
    })
}
if (document.getElementById("modal-editcomment") !== null) {
    $('#modal-editcomment').on('shown.bs.modal', function () {
        $('#modal-editcomment-text').trigger('focus')
    })
}
if (document.getElementById("modal-setnote") !== null) {
    $('#modal-setnote').on('shown.bs.modal', function () {
        $('#modal-setnote-text').trigger('focus')
    })
}
if (document.getElementById("modal-addsource") !== null) {
    $('#modal-addsource').on('shown.bs.modal', function () {
        $('#modal-addsource-subject').trigger('focus')
    })
}
if (document.getElementById("modal-rmsource") !== null) {
    $('#modal-rmsource').on('shown.bs.modal', function () {
        $('#modal-rmsource-subject').trigger('focus')
    })
}
if (document.getElementById("modal-setrnum") !== null) {
    $('#modal-setrnum').on('shown.bs.modal', function () {
        $('#modal-setrnum-text').trigger('focus')
    })
}
if (document.getElementById("modal-addtag") !== null) {
    $('#modal-addtag').on('shown.bs.modal', function () {
        $('#modal-addtag-text').trigger('focus')
    })
}
if (document.getElementById("modal-rmtag") !== null) {
    $('#modal-rmtag').on('shown.bs.modal', function () {
        $('#modal-rmtag-text').trigger('focus')
    })
}
if (document.getElementById("modal-deadline2d") !== null) {
    $('#modal-deadline2d').on('shown.bs.modal', function () {
        $('#modal-deadline2d-date').trigger('focus')
    })
}
if (document.getElementById("modal-deadline3d") !== null) {
    $('#modal-deadline3d').on('shown.bs.modal', function () {
        $('#modal-deadline3d-date').trigger('focus')
    })
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

function setCookie(name,value,days) {
    var expires = "";
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days*24*60*60*1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + (value || "")  + expires + "; path=/";
}

function eraseCookie(name) {   
    document.cookie = name +'=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}

// Hotkey: http://gcctech.org/csc/javascript/javascript_keycodes.htm
document.onkeyup = function(e) {
    if (e.ctrlKey && e.shiftKey && e.which == 65) {
        selectCheckboxAll()
    } else if (e.ctrlKey && e.altKey && e.shiftKey && e.which == 68) {
        selectCheckboxNone()
    } else if (e.ctrlKey && e.altKey && e.shiftKey && e.which == 73) {
        selectCheckboxInvert()
    } else if (e.ctrlKey && e.altKey && e.shiftKey && e.which == 84) {
        scroll(0,0)
    } else if (e.ctrlKey && e.altKey && e.shiftKey && e.which == 77) {
        selectmodeV2()
    } else if (e.ctrlKey && e.altKey && e.shiftKey && e.which == 69) {
        OpenEditfolder()
    } else if (e.which == 119) { // F8
        if ($('#modal-addcomment').hasClass('show')) {
            document.getElementById("modal-addcomment-addbutton").click();
        }
        if ($('#modal-editcomment').hasClass('show')) {
            document.getElementById("modal-editcomment-editbutton").click();
        }
        if ($('#modal-setnote').hasClass('show')) {
            document.getElementById("modal-setnote-editbutton").click();
        }
        if ($('#modal-editnote').hasClass('show')) {
            document.getElementById("modal-editnote-editbutton").click();
        }
    }
};

function OpenEditfolder() {
    let uri = document.getElementById("edit").href;
    window.location = uri;
}

function padNumber(number) {
    number = number.toString();
    while(number.length < 4) {
        number = "0" + number;
    }
    return number;
}

function id2name(id) {
    l = id.split("_");
    l.pop()
    return l.join("_")
}

// addTagAtSearchbox 함수는 태그가 searchbox-searchbox 에 없다면 검색어로 추가한다.
function addTagAtSearchbox(tag) {
    let searchboxWord = document.getElementById("searchbox-searchword")
    if (!searchboxWord.value.includes("tag:"+tag)) {
        searchboxWord.value = searchboxWord.value + " tag:"+tag
    }
}

// addAssetTagAtSearchbox 함수는 태그가 searchbox-searchbox 에 없다면 검색어로 추가한다.
function addAssetTagAtSearchbox(assettag) {
    let searchboxWord = document.getElementById("searchbox-searchword")
    if (!searchboxWord.value.includes("assettag:"+assettag)) {
        searchboxWord.value = searchboxWord.value + " assettag:"+assettag
    }
}

// setModal 함수는 modalID와 value를 받아서 modal에 셋팅한다.
function setModal(modalID, value) {
    document.getElementById(modalID).value = value;
}


// setEditTaskModal 함수는 item id, task 정보를 가지고 와서 Edit Task Modal에 값을 채운다.
function setEditTaskModal(id, task) {
    document.getElementById("modal-edittask-id").value = id;
    document.getElementById("modal-edittask-title").innerHTML = "Edit Task" + multiInputTitle(id);
    $.ajax({
        url: "/api/task",
        type: "POST",
        data: {
            id: id,
            task: task,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-edittask-start').value=data.task.start;
            document.getElementById('modal-edittask-end').value=data.task.end;
            document.getElementById('modal-edittask-task').value=data.task.title;            
            document.getElementById('modal-edittask-path').value=data.task.mov;
            document.getElementById('modal-edittask-usernote').value=data.task.usernote;
            document.getElementById('modal-edittask-user').value=data.task.user;
            document.getElementById('modal-edittask-usercomment').value=data.task.usercomment;
            document.getElementById('modal-edittask-id').value=data.id;
            document.getElementById("modal-edittask-statusv2").value=data.task.statusv2;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setTimeModal 함수는 id 정보를 가지고 와서 Edit Time Modal 값을 채운다.
function setTimeModal(id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-edittime-id").value = id;
    document.getElementById("modal-edittime-title").innerHTML = "Edit Time" + multiInputTitle(id);
    $.ajax({
        url: "/api/timeinfo",
        type: "POST",
        data: {
            "id": id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('scanin').value = data.scanin;
            document.getElementById('scanout').value = data.scanout;
            document.getElementById('scanframe').value = data.scanframe;
            document.getElementById('scantimecodein').value = data.scantimecodein;
            document.getElementById('scantimecodeout').value = data.scantimecodeout;
            document.getElementById('platein').value = data.platein;
            document.getElementById('plateout').value = data.plateout;
            document.getElementById('handlein').value = data.handlein;
            document.getElementById('handleout').value = data.handleout;
            document.getElementById('justin').value = data.justin;
            document.getElementById('justout').value = data.justout;
            document.getElementById('justtimecodein').value = data.justtimecodein;
            document.getElementById('justtimecodeout').value = data.justtimecodeout;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setShottypeModal 함수는 item id 정보를 이용해 Edit Shottype Modal에 값을 채운다.
function setShottypeModal(id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-shottype-title").innerHTML = "Shot Type" + multiInputTitle(id);
    $.ajax({
        url: "/api/shottype",
        type: "POST",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-shottype-id').value=id;
            document.getElementById("modal-shottype-type").value=data.shottype;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setUsetypeModal 함수는 project, name을 받아서 Usetype Modal을 설정한다.
function setUsetypeModal(id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-usetype-id").value = id
    document.getElementById("modal-usetype-title").innerHTML = "Use Type: " + id;
    $.ajax({
        url: `/api/usetypes?id=${id}`,
        type: "GET",
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            let types = data.types;
            // selectbox를 채운다.
            let sel = document.getElementById('modal-usetype-type');
            sel.innerHTML = "";
            for (let i = 0; i < types.length; i++) {
                let opt = document.createElement('option');
                opt.appendChild( document.createTextNode(types[i]) );
                opt.value = types[i]; 
                sel.appendChild(opt); 
            }
            // 이미 선택된 옵션을 selectbox에서 선택한다.
            document.getElementById("modal-usetype-type").value=data.usetype;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setModalCheckbox(modalID, value) {
    if (value === "true") {
        document.getElementById(modalID).checked = true;
    } else {
        document.getElementById(modalID).checked = false;
    }
}

// *문자를 x문자로 바꾼다.
// X를 x문자로 바꾼다.
// 공백을 제거한다.
// 렌즈디스토션값을 입력시 2048*1280 -> 2048x1280 형태로 바꾸기 위함이다.
// 숫자와 x를 제외한 영문입력시 삭제됩니다.
function widthxHeight(event) {
	event = event || window.event;
	event.target.value = event.target.value.replace("*","x");
	event.target.value = event.target.value.replace("X","x");
	event.target.value = event.target.value.replace(/[^\d\x]/gi,"");
}

function sleep( millisecondsToWait ) {
    var now = new Date().getTime();
    while ( new Date().getTime() < now + millisecondsToWait ) {
        /* do nothing; this will exit once it reaches the time limit */
        /* if you want you could do something and exit */
    }
}

function multiInputTitle(id) {
    let checknum = 0;
    let cboxes = document.getElementsByName('selectID');
    for (let i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === true) {
            checknum += 1
        }
    }
    if (checknum === 0) {
        return ": " + id2name(id)
    } else if (checknum === 1) {
        return ": " + id2name(id)        
    } else {
        let name = id2name(id);
        let num = checknum - 1
        return `: ${name}외 ${num}건`
    }
}



function addTask() {
    let token = document.getElementById("token").value;
    
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);

            fetch('/api/addtask', {
                method: 'POST',
                headers: {
                    "Authorization": "Basic "+ document.getElementById("token").value,
                },
                body: new URLSearchParams({
                    id: id,
                    task: document.getElementById('modal-addtask-taskname').value,
                })
            })
            .then((response) => {
                return response.json()
            })
            .then((data) => {
                let newItem = `<div class="row" id="${data.id}-task-${data.task}">
                <div id="${data.id}-task-${data.task}-status">
                    <span class="finger mt-1 badge badge-${data.status} statusbox">${data.task}</span>
                </div>
                <div id="${data.id}-task-${data.task}-predate"></div>
                <div id="${data.id}-task-${data.task}-date"></div>
                <div id="${data.id}-task-${data.task}-user"></div>
                <div id="${data.id}-task-${data.task}-playbutton"></div>
                <div class="ml-1">
                    <span class="add" data-toggle="modal" data-target="#modal-edittask" onclick="
                    setEditTaskModal('${data.id}', '${data.task}');
                    ">≡</span>
                </div>
                </div>`;
                document.getElementById(`${data.id}-tasks`).innerHTML = newItem + document.getElementById(`${data.id}-tasks`).innerHTML;
            })
            .catch((error) => {
                alert(error)
            });
        }
    } else {
        fetch('/api/addtask', {
            method: 'POST',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
            body: new URLSearchParams({
                id: document.getElementById('modal-addtask-id').value,
                task: document.getElementById('modal-addtask-taskname').value,
            })
        })
        .then((response) => {
            return response.json()
        })
        .then((data) => {
            let newItem = `<div class="row" id="${data.id}-task-${data.task}">
            <div id="${data.id}-task-${data.task}-status">
                <span class="finger mt-1 badge badge-${data.status} statusbox">${data.task}</span>
            </div>
            <div id="${data.id}-task-${data.task}-predate"></div>
            <div id="${data.id}-task-${data.task}-date"></div>
            <div id="${data.id}-task-${data.task}-user"></div>
            <div id="${data.id}-task-${data.task}-playbutton"></div>
            <div class="ml-1">
                <span class="add" data-toggle="modal" data-target="#modal-edittask" onclick="
                setEditTaskModal('${data.id}', '${data.task}');
                ">≡</span>
            </div>
            </div>`;
            document.getElementById(`${data.id}-tasks`).innerHTML = newItem + document.getElementById(`${data.id}-tasks`).innerHTML;
        })
        .catch((error) => {
            alert(error)
        });
    }
}

function rmTask(project, id, task) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/rmtask",
                type: "POST",
                data: {
                    id: id,
                    task: task,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`${data.id}-task-${data.task}`).remove();
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/rmtask",
            type: "POST",
            data: {
                id: id,
                task: task,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`${data.id}-task-${data.task}`).remove();
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setFrame(mode, id, frame) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/" + mode,
        type: "POST",
        data: {
            id: id,
            frame: frame,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (mode === "setscanin") {
                document.getElementById("scanin-"+data.id).innerHTML = `<span class="text-badge ml-1" title="scanin">${data.frame}</span>`;
            }
            if (mode === "setscanout") {
                document.getElementById("scanout-"+data.id).innerHTML = `<span class="text-badge ml-1" title="scanout">${data.frame}</span>`;
            }
            if (mode === "setscanframe") {
                document.getElementById("scanframe-"+data.id).innerHTML = `<span class="text-badge ml-1" title="scanframe">(${data.frame})</span>`;
            }
            if (mode === "sethandlein") {
                document.getElementById("handlein-"+data.id).innerHTML = data.frame;
            }
            if (mode === "sethandleout") {
                document.getElementById("handleout-"+data.id).innerHTML = data.frame;
            }
            if (mode === "setplatein") {
                document.getElementById("platein-"+data.id).innerHTML = `<span class="text-white black-opbg" title="platein">${data.frame}</span>`;
            }
            if (mode === "setplateout") {
                document.getElementById("plateout-"+data.id).innerHTML = `<span class="text-white black-opbg" title="plateout">${data.frame}</span>`;
            }
            if (mode === "setjustin") {
                document.getElementById("justin-"+data.id).innerHTML = `<span class="text-warning black-opbg" title="justin">${data.frame}</span>`;
            }
            if (mode === "setjustout") {
                document.getElementById("justout-"+data.id).innerHTML = `<span class="text-warning black-opbg" title="justout">${data.frame}</span>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setScanTimecodeIn(id, timecode) {
    $.ajax({
        url: "/api/setscantimecodein",
        type: "POST",
        data: {
            id: id,
            timecode: timecode,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("scantimecodein-"+data.id).innerHTML = `<span class="text-badge ml-1">${data.timecode}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraPubTask() {
    $.ajax({
        url: "/api/setcamerapubtask",
        type: "POST",
        data: {
            project: document.getElementById('modal-cameraoption-project').value,
            id: document.getElementById('modal-cameraoption-id').value,
            task: document.getElementById('modal-cameraoption-pubtask').value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("campubtask-"+data.id).innerHTML = `<span class="text-badge ml-1">${data.task}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraLensmm() {
    $.ajax({
        url: "/api/setcameralensmm",
        type: "POST",
        data: {
            project: document.getElementById('modal-cameraoption-project').value,
            id: document.getElementById('modal-cameraoption-id').value,
            lensmm: document.getElementById('modal-cameraoption-lensmm').value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("camlensmm-"+data.id).innerHTML = `<span class="text-badge ml-1">${data.lensmm}mm</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraPubPath() {
    $.ajax({
        url: "/api/setcamerapubpath",
        type: "POST",
        data: {
            project: document.getElementById('modal-cameraoption-project').value,
            id: document.getElementById('modal-cameraoption-id').value,
            path: document.getElementById('modal-cameraoption-pubpath').value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("campubpath-"+data.id).innerHTML = `<a href="${data.protocol}://${data.path}" class="text-badge ml-1">${data.path}</a>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setOutputModal(id) {
    document.getElementById("modal-output-id").value = id;
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-output-finver').value = data.finver;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setCameraOptionModal(project, id) {
    document.getElementById("modal-cameraoption-project").value = project;
    document.getElementById("modal-cameraoption-id").value = id;
    document.getElementById("modal-cameraoption-title").innerHTML = "Camera Option" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-cameraoption-pubtask').value = data.productioncam.pubtask;
            document.getElementById('modal-cameraoption-lensmm').value = data.productioncam.lensmm;
            document.getElementById('modal-cameraoption-pubpath').value = data.productioncam.pubpath;
            if (data.productioncam.projection) {
                document.getElementById("modal-cameraoption-projection").checked = true;
            } else {
                document.getElementById("modal-cameraoption-projection").checked = false;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraProjection() {
    $.ajax({
        url: "/api/setcameraprojection",
        type: "POST",
        data: {
            project: document.getElementById('modal-cameraoption-project').value,
            id: document.getElementById('modal-cameraoption-id').value,
            projection: document.getElementById("modal-cameraoption-projection").checked,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            if (data.projection) {
                document.getElementById("camprojection-"+data.id).innerHTML = `<span class="text-badge ml-1">ProjectionCam</span>`;
            } else {
                document.getElementById("camprojection-"+data.id).innerHTML = "";
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setScanTimecodeOut(id, timecode) {
    $.ajax({
        url: "/api/setscantimecodeout",
        type: "POST",
        data: {
            id: id,
            timecode: timecode,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("scantimecodeout-"+data.id).innerHTML = `<span class="text-badge ml-1">${data.timecode}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setJustTimecodeIn(id, timecode) {
    $.ajax({
        url: "/api/setjusttimecodein",
        type: "POST",
        data: {
            id: id,
            timecode: timecode,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("justtimecodein-"+data.id).innerHTML = `<span class="text-warning black-opbg">${data.timecode}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setJustTimecodeOut(id, timecode) {
    $.ajax({
        url: "/api/setjusttimecodeout",
        type: "POST",
        data: {
            id: id,
            timecode: timecode,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("justtimecodeout-"+data.id).innerHTML = `<span class="text-warning black-opbg">${data.timecode}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setNoteModal(id) {
    document.getElementById("modal-setnote-id").value = id;
    document.getElementById("modal-setnote-title").innerHTML = "Set Note" + multiInputTitle(id);
    document.getElementById("modal-setnote-text").value = "";
}

function editNoteModal(id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-editnote-id").value = id;
    document.getElementById("modal-editnote-title").innerHTML = "Set Note" + multiInputTitle(id);
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("modal-editnote-text").value = data.note.text;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setNote(id, text) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let currentID = cboxes[i].getAttribute("id");
            let overwrite = document.getElementById("modal-setnote-overwrite").checked;
            sleep(200);
            $.ajax({
                url: "/api/setnote",
                type: "POST",
                data: {
                    id: currentID,
                    text: text,
                    userid: userid,
                    overwrite: overwrite,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    if (overwrite) {
                        // note-{{.Name}} 내부 내용을 교체한다.
                        document.getElementById("note-"+data.id).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                    } else {
                        // note-{{.Name}} 내부 내용에 추가한다.
                        document.getElementById("note-"+data.id).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>') + "<br>" + document.getElementById("note-"+data.id).innerHTML;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        let overwrite = document.getElementById("modal-setnote-overwrite").checked;
        $.ajax({
            url: "/api/setnote",
            type: "POST",
            data: {
                id: id,
                text: text,
                userid: userid,
                overwrite: overwrite,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                if (overwrite) {
                    // note-{{.Name}} 내부 내용을 교체한다.
                    document.getElementById("note-"+data.id).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                } else {
                    // note-{{.Name}} 내부 내용에 추가한다.
                    document.getElementById("note-"+data.id).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>') + "<br>" + document.getElementById("note-"+data.id).innerHTML;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function editNote(id, text) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/setnote",
        type: "POST",
        data: {
            id: id,
            text: text,
            userid: userid,
            overwrite: true,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("note-"+data.id).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// initPublishModal 함수는 modal-addpublish 를 초기화 한다.
function initPublishModal() {
    document.getElementById("modal-addpublish-project").value = ""
    document.getElementById("modal-addpublish-name").value = ""
    document.getElementById("modal-addpublish-task").value = ""
    document.getElementById("modal-addpublish-secondarykey").value = ""
    document.getElementById("modal-addpublish-path").value = ""
    document.getElementById("modal-addpublish-status").value= "usethis"
    document.getElementById("modal-addpublish-tasktouse").value= ""
    document.getElementById("modal-addpublish-subject").value= ""
    document.getElementById("modal-addpublish-mainversion").value= 1
    document.getElementById("modal-addpublish-subversion").value= 0
    document.getElementById("modal-addpublish-filetype").value= ""
    document.getElementById("modal-addpublish-kindofusd").value= ""
    document.getElementById("modal-addpublish-outputdatapath").value= ""
}

function setAddPublishModal(project, name, task) {
    initPublishModal() // AddPublishModal을 한번 초기화 한다.
    document.getElementById("modal-addpublish-project").value = project
    document.getElementById("modal-addpublish-name").value = name
    document.getElementById("modal-addpublish-task").value = task
    fetch("/api/publishkeys", {
        methode: "GET",
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
    })
    .then((response) => {
        return response.json()
    })
    .then((datas) => {
        if (datas.length == 0) {
            alert("PublishKey registration is required.");
            document.getElementById('modal-addpublish-addbutton').disabled = true;
            return
        }
        let keys = document.getElementById('modal-addpublish-key');
        keys.innerHTML = "";
        for (let i = 0; i < datas.length; i++){
            let opt = document.createElement('option');
            opt.value = datas[i].id;
            opt.innerHTML = datas[i].id;
            keys.appendChild(opt);
        }
    })
}

function setEditPublishModal(project, id, task, tasktouse, key, createtime, path) {
    document.getElementById("modal-editpublish-project").value = project
    document.getElementById("modal-editpublish-id").value = id
    document.getElementById("modal-editpublish-task").value = task
    document.getElementById("modal-editpublish-tasktouse").value = tasktouse
    document.getElementById("modal-editpublish-path").value = path
    // publishkey를 셋팅한다.
    fetch('/api/publishkeys', {
        method: 'GET',
        headers: {
            "Authorization": "Basic " + document.getElementById("token").value,
        },
    })
    .then((response) => {
        return response.json()
    })
    .then((datas) => {
        if (datas.length == 0) {
            alert("PublishKey registration is required.");
            document.getElementById('modal-editpublish-editbutton').disabled = true;
            return
        }
        let keys = document.getElementById('modal-editpublish-key');
        keys.innerHTML = "";
        for (let i = 0; i < datas.length; i++){
            let opt = document.createElement('option');
            opt.value = datas[i].id;
            opt.innerHTML = datas[i].id;
            keys.appendChild(opt);
        }
    })
    .catch((error) => {
        alert(error)
    });

    // 아이템 정보를 가지고 와서 modal-editpublish를 채운다.
    fetch('/api/getpublish', {
        method: 'post',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: new URLSearchParams({
            "project": project,
            "id": id,
            "task": task,
            "key": key,
            "path": path,
            "createtime": createtime,
        })
    })
    .then((response) => {
        return response.json()
    })
    .then((data) => {
        document.getElementById('modal-editpublish-key').value = key;
        document.getElementById('modal-editpublish-secondarykey').value = data.secondarykey;
        document.getElementById('modal-editpublish-path').value = data.path;
        document.getElementById('modal-editpublish-status').value = data.status;
        document.getElementById('modal-editpublish-tasktouse').value = data.tasktouse;
        document.getElementById('modal-editpublish-subject').value = data.subject;
        document.getElementById('modal-editpublish-mainversion').value = data.mainversion;
        document.getElementById('modal-editpublish-subversion').value = data.subversion;
        document.getElementById('modal-editpublish-filetype').value = data.filetype;
        document.getElementById('modal-editpublish-filetype').value = data.filetype;
        document.getElementById('modal-editpublish-kindofusd').value = data.kindofusd;
        document.getElementById('modal-editpublish-createtime').value = data.createtime;
        document.getElementById('modal-editpublish-isoutput').checked = data.isoutput;
        document.getElementById('modal-editpublish-outputdatapath').value = data.outputdatapath;
    })
    .catch((error) => {
        alert(error)
    });
}

function editPublish() {
    let project = document.getElementById('modal-editpublish-project').value
    let id = document.getElementById('modal-editpublish-id').value
    let task = document.getElementById('modal-editpublish-task').value
    let key = document.getElementById('modal-editpublish-key').value
    let createtime = document.getElementById('modal-editpublish-createtime').value
    let secondarykey = document.getElementById('modal-editpublish-secondarykey').value
    let path = document.getElementById('modal-editpublish-path').value
    let status = document.getElementById('modal-editpublish-status').value
    let tasktouse = document.getElementById('modal-editpublish-tasktouse').value
    let subject = document.getElementById('modal-editpublish-subject').value
    let mainversion = document.getElementById('modal-editpublish-mainversion').value
    let subversion = document.getElementById('modal-editpublish-subversion').value
    let filetype = document.getElementById('modal-editpublish-filetype').value
    let kindofusd = document.getElementById('modal-editpublish-kindofusd').value
    let isoutput = false
    if (document.getElementById('modal-editpublish-isoutput').checked) {
        isoutput = true
    }
    let outputdatapath = document.getElementById('modal-editpublish-outputdatapath').value
    // 기존 데이터를 삭제한다.
    fetch('/api/rmpublish', {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: new URLSearchParams({
            project: project,
            id: id,
            task: task,
            key: key,
            path: path,
            createtime: createtime,
        })
    })
    .then((response) => {
        return response.json()
    })
    .then((data) => {
        return
    })
    .catch((error) => {
        alert(error)
    });

    // 새로운 데이터를 추가한다.
    fetch('/api/addpublish', {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: new URLSearchParams({
            project: project,
            name: id2name(id),
            task: task,
            key: key,
            createtime: createtime,
            secondarykey: secondarykey,
            path: path,
            status: status,
            tasktouse: tasktouse,
            subject: subject,
            mainversion: mainversion,
            subversion: subversion,
            filetype: filetype,
            kindofusd: kindofusd,
            isoutput: isoutput,
            outputdatapath: outputdatapath,
        })
    })
    .then((response) => {
        return response.json()
    })
    .then((data) => {
        location.reload()
        return
    })
    .catch((error) => {
        alert(error)
    });
}

function setAddCommentModal(id) {
    document.getElementById("modal-addcomment-id").value = id;
    document.getElementById("modal-addcomment-title").innerHTML = "Add Comment" + multiInputTitle(id);
    document.getElementById("modal-addcomment-mediatitle").value = "";
    document.getElementById("modal-addcomment-media").value = "";
}

function addComment() {
    let token = document.getElementById("token").value;
    let id = document.getElementById('modal-addcomment-id').value;
    let text = document.getElementById('modal-addcomment-text').value
    let media = document.getElementById('modal-addcomment-media').value
    let mediatitle = document.getElementById('modal-addcomment-mediatitle').value
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (let i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let currentID = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/addcomment",
                type: "POST",
                data: {                    
                    id: currentID,
                    text: text,
                    mediatitle: mediatitle,
                    media: media,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    // comments-{{.Name}} 내부 내용에 추가한다.
                    let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                    let newComment = `<div id="comment-${data.id}-${data.date}">
                    <span class="text-badge">${data.date} / <a href="/user?id=${data.userid}" class="text-darkmode">${data.authorname}</a></span>
                    <span class="edit" data-toggle="modal" data-target="#modal-editcomment" onclick="setEditCommentModal('${data.id}', '${data.date}', '${data.text}', '${data.mediatitle}', '${data.media}')">≡</span>
                    <span class="remove" data-toggle="modal" data-target="#modal-rmcomment" onclick="setRmCommentModal('${data.id}', '${data.date}', '${data.text}')">×</span>
                    <br><small class="text-warning">${body}</small>`
                    if (data.media != "") {
                        if (data.media.includes("http")) {
                            newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                        } else {
                            newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="${data.protocol}://${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                        }
                    }
                    newComment += `<hr class="my-1 p-0 m-0 divider"></hr></div>`
                    document.getElementById("comments-"+data.id).innerHTML = newComment + document.getElementById("comments-"+data.id).innerHTML;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/addcomment",
            type: "POST",
            data: {
                id: id,
                text: text,
                mediatitle: mediatitle,
                media: media,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                // comments-{{$id}} 내부 내용에 추가한다.
                let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                let newComment = `<div id="comment-${data.id}-${data.date}">
                <span class="text-badge">${data.date} / <a href="/user?id=${data.userid}" class="text-darkmode">${data.authorname}</a></span>
                <span class="edit" data-toggle="modal" data-target="#modal-editcomment" onclick="setEditCommentModal('${data.id}', '${data.date}', '${data.text}', '${data.mediatitle}', '${data.media}')">≡</span>
                <span class="remove" data-toggle="modal" data-target="#modal-rmcomment" onclick="setRmCommentModal('${data.id}', '${data.date}', '${data.text}')">×</span>
                <br><div class="text-warning small">${body}</div>`
                if (data.media != "") {
                    if (data.media.includes("http")) {
                        newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                    } else {
                        newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="${data.protocol}://${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                    }
                }
                newComment += `<hr class="my-1 p-0 m-0 divider"></hr></div>`
                document.getElementById("comments-"+data.id).innerHTML = newComment + document.getElementById("comments-"+data.id).innerHTML;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function editComment() {
    let token = document.getElementById("token").value;
    let id = document.getElementById('modal-editcomment-id').value;
    let time = document.getElementById('modal-editcomment-time').value
    let text = document.getElementById('modal-editcomment-text').value
    let mediatitle = document.getElementById('modal-editcomment-mediatitle').value
    let media = document.getElementById('modal-editcomment-media').value
    $.ajax({
        url: "/api/editcomment",
        type: "POST",
        data: {
            id: id,
            time: time,
            text: text,
            mediatitle: mediatitle,
            media: media,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            // comments-${data.id}}-${data.time} 내부 내용을 업데이트 한다.
            let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
            let newComment = `<span class="text-badge">${data.time} / <a href="/user?id=${data.userid}" class="text-darkmode">${data.authorname}</a></span>
            <span class="edit" data-toggle="modal" data-target="#modal-editcomment" onclick="setEditCommentModal('${data.id}', '${data.time}', '${data.text}', '${data.mediatitle}', '${data.media}')">≡</span>
            <span class="remove" data-toggle="modal" data-target="#modal-rmcomment" onclick="setRmCommentModal('${data.id}', '${data.time}', '${data.text}')">×</span>
            <br><div class="text-warning small">${body}</div>`
            if (data.media != "") {
                if (data.media.includes("http")) {
                    newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                } else {
                    newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="${data.protocol}://${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                }
            }
            document.getElementById(`comment-${data.id}-${data.time}`).innerHTML = newComment
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setRmCommentModal(id, time, text) {
    document.getElementById("modal-rmcomment-id").value = id;
    document.getElementById("modal-rmcomment-time").value = time;
    document.getElementById("modal-rmcomment-text").value = text;
    document.getElementById("modal-rmcomment-title").innerHTML = "Rm Comment" + multiInputTitle(id);
}

function setEditReviewModal(id) {
    document.getElementById("modal-editreview-id").value = id;
    // review id의 데이터를 가지고 와서 모달을 설정한다.
    fetch('/api/review', {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: new URLSearchParams({
            id: id,
        })
    })
    .then((response) => {
        return response.json()
    })
    .then((data) => {
        document.getElementById("modal-editreview-project").value = data.project;
        document.getElementById("modal-editreview-task").value = data.task;
        document.getElementById("modal-editreview-name").value = data.name;
        document.getElementById("modal-editreview-itemstatus").value = data.itemstatus;
        document.getElementById("modal-editreview-createtime").value = data.createtime;
        document.getElementById("modal-editreview-path").value = data.path;
        document.getElementById("modal-editreview-mainversion").value = data.mainversion;
        document.getElementById("modal-editreview-subversion").value = data.subversion;
        document.getElementById("modal-editreview-fps").value = data.fps;
        document.getElementById("modal-editreview-description").value = data.description;
        document.getElementById("modal-editreview-camerainfo").value = data.camerainfo;
        document.getElementById("modal-editreview-outputdatapath").value = data.outputdatapath;
    })
    .catch((error) => {
        alert(error)
    });
}

function setRmReviewCommentModal(id, time) {
    document.getElementById("modal-rmreviewcomment-id").value = id;
    document.getElementById("modal-rmreviewcomment-time").value = time;
    // review id의 데이터를 가지고 와서 모달을 설정한다.
    $.ajax({
        url: "/api/review",
        type: "POST",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("modal-rmreviewcomment-project").value = data.project;
            document.getElementById("modal-rmreviewcomment-name").value = data.name;
            for (let i = 0; i < data.comments.length; i++) {
                if (data.comments[i].date == time) {
                    document.getElementById("modal-rmreviewcomment-text").value = data.comments[i].text;
                    break
                }
            }
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

function setEditReviewCommentModal(id, time) {
    document.getElementById("modal-editreviewcomment-id").value = id;
    document.getElementById("modal-editreviewcomment-time").value = time;
    // review id의 데이터를 가지고 와서 모달을 설정한다.
    $.ajax({
        url: "/api/review",
        type: "POST",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            for (let i = 0; i < data.comments.length; i++) {
                if (data.comments[i].date == time) {
                    document.getElementById("modal-editreviewcomment-text").value = data.comments[i].text;
                    document.getElementById("modal-editreviewcomment-media").value = data.comments[i].media;
                    document.getElementById("modal-editreviewcomment-frame").value = data.comments[i].frame;
                    break
                }
            }
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

function setRmReviewModal(id) {
    // review id의 데이터를 가지고 와서 모달을 설정한다.
    $.ajax({
        url: "/api/review",
        type: "POST",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("modal-rmreview-id").innerHTML = "ID: " + id;
            document.getElementById("modal-rmreview-id").value = id;
            document.getElementById("modal-rmreview-project").innerHTML = "Project: " + data.project;
            document.getElementById("modal-rmreview-name").innerHTML = "Name: " + data.name;
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

function setRmPublishKeyModal(project, id, task, key) {
    document.getElementById("modal-rmpublishkey-project").value = project;
    document.getElementById("modal-rmpublishkey-id").value = id;
    document.getElementById("modal-rmpublishkey-task").value = task;
    document.getElementById("modal-rmpublishkey-key").value = key;
    document.getElementById("modal-rmpublishkey-title").innerHTML = "Rm Publish Key" + multiInputTitle(id);
}

function setRmPublishModal(project, id, task, key, createtime, path) {
    document.getElementById("modal-rmpublish-project").value = project
    document.getElementById("modal-rmpublish-id").value = id
    document.getElementById("modal-rmpublish-task").value = task
    document.getElementById("modal-rmpublish-key").value = key
    document.getElementById("modal-rmpublish-createtime").value = createtime
    document.getElementById("modal-rmpublish-path").value = path
}

function setPublishModal(project, id, task, key, path, createtime, status) {
    document.getElementById("modal-setpublish-project").value = project;
    document.getElementById("modal-setpublish-id").value = id;
    document.getElementById("modal-setpublish-task").value = task;
    document.getElementById("modal-setpublish-key").value = key;
    document.getElementById("modal-setpublish-path").value = path;
    document.getElementById("modal-setpublish-createtime").value = createtime;
    document.getElementById("modal-setpublish-status").value = status;
    document.getElementById("modal-setpublish-status").innerHTML = status;
}

function setEditCommentModal(id, time, text, mediatitle, media) {
    document.getElementById("modal-editcomment-id").value = id;
    document.getElementById("modal-editcomment-time").value = time;
    document.getElementById("modal-editcomment-text").value = text;
    document.getElementById("modal-editcomment-mediatitle").value = mediatitle;
    document.getElementById("modal-editcomment-media").value = media;
    document.getElementById("modal-editcomment-title").innerHTML = "Edit Comment" + multiInputTitle(id);
}

function setDetailCommentsModal(project, id) {
    document.getElementById("modal-detailcomments-title").innerHTML = "Detail Comments" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic " + token
        },
        dataType: "json",
        success: function(data) {
            // 기존 디테일을 지운다.
            document.getElementById('modal-detailcomments-body').innerHTML = "";
            // 코멘트를 추가한다.
            let comments = data.comments
            comments.reverse();
            for (var i = 0; i < comments.length; ++i) {
                let br = document.createElement("br")
                // elements way
                let cmt = document.createElement("div")
                cmt.setAttribute("id", "comment-"+data.id+"-"+comments[i].date)
                let info = document.createElement("span")
                info.setAttribute("class","text-badge")
                info.innerHTML = comments[i].date + " / "
                let userinfo = document.createElement("a")
                userinfo.setAttribute("href", "/user?id="+data.comments[i].author)
                userinfo.setAttribute("class","text-darkmode")
                userinfo.innerHTML = comments[i].author
                info.append(userinfo)
                cmt.append(info)
                cmt.append(br)
                let text = document.createElement("div")
                text.setAttribute("class","text-darkmode small")
                text.innerHTML = "<br />" + comments[i].text.replace(/\n/g, "<br />")
                cmt.append(text)
                cmt.append(br)
                if (comments[i].media !== "") {
                    let link = document.createElement("a")
                    let protocol = document.getElementById("protocol").value + "://"
                    if (comments[i].media.startsWith("http")) {
                        protocol = ""
                    }
                    link.setAttribute("href", protocol + comments[i].media)
                    link.innerHTML = `<img src="/assets/img/link.svg" class="finger">`
                    cmt.append(link)
                    let span = document.createElement("span")
                    span.setAttribute("class","text-darkmode small pl-2")
                    span.innerHTML = comments[i].mediatitle
                    cmt.append(span)
                }
                let line = document.createElement("hr")
                line.setAttribute("class","my-1 p-0 m-0 divider")
                cmt.append(line)
                let parents= document.getElementById('modal-detailcomments-body')
                parents.append(cmt)
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setRmUserModal(id) {
    document.getElementById("modal-rmuser-id").value = id;
}

function rmComment(id, date) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/rmcomment",
        type: "POST",
        data: {
            id: id,
            date: date
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`comment-${data.id}-${data.date}`).remove();
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function rmReview() {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/rmreview",
        type: "POST",
        data: {
            id: document.getElementById("modal-rmreview-id").value,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`review-${data.id}`).remove();
            initCanvas();
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function editReviewComment() {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/editreviewcomment",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreviewcomment-id").value,
            time: document.getElementById("modal-editreviewcomment-time").value,
            text: document.getElementById("modal-editreviewcomment-text").value,
            media: document.getElementById("modal-editreviewcomment-media").value,
            frame: document.getElementById("modal-editreviewcomment-frame").value,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`reviewcomment-${data.id}-${data.time}-text`).innerText = data.text
            if (data.frame == 0) {
                document.getElementById(`reviewcomment-${data.id}-${data.time}-frame`).remove();
            } else {
                document.getElementById(`reviewcomment-${data.id}-${data.time}-frame`).innerText = data.frame + "/" + (data.frame + data.productionstartframe - 1)
            }
            
            if (data.media.startsWith("http") || data.media.startsWith("rvlink")) {
                document.getElementById(`reviewcomment-${data.id}-${data.time}-media`).setAttribute("href", data.media)
            } else {
                document.getElementById(`reviewcomment-${data.id}-${data.time}-media`).setAttribute("href", data.protocol + "://" + data.media)
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function rmReviewComment() {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/rmreviewcomment",
        type: "POST",
        data: {
            id: document.getElementById("modal-rmreviewcomment-id").value,
            time: document.getElementById("modal-rmreviewcomment-time").value,
            project: document.getElementById("modal-rmreviewcomment-project").value,
            name: document.getElementById("modal-rmreviewcomment-name").value,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`reviewcomment-${data.id}-${data.time}`).remove();
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function rmPublishKey() {
    $.ajax({
        url: "/api/rmpublishkey",
        type: "POST",
        data: {
            project: document.getElementById('modal-rmpublishkey-project').value,
            id: document.getElementById('modal-rmpublishkey-id').value,
            task: document.getElementById('modal-rmpublishkey-task').value,
            key: document.getElementById('modal-rmpublishkey-key').value
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            location.reload()
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function rmPublish() {
    $.ajax({
        url: "/api/rmpublish",
        type: "POST",
        data: {
            project: document.getElementById('modal-rmpublish-project').value,
            id: document.getElementById('modal-rmpublish-id').value,
            task: document.getElementById('modal-rmpublish-task').value,
            key: document.getElementById('modal-rmpublish-key').value,
            createtime: document.getElementById('modal-rmpublish-createtime').value,
            path: document.getElementById('modal-rmpublish-path').value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function() {
            location.reload()
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setAddSourceModal(id) {
    document.getElementById("modal-addsource-id").value = id;
    document.getElementById("modal-addsource-subject").value = "";
    document.getElementById("modal-addsource-path").value = "";
    document.getElementById("modal-addsource-title").innerHTML = "Add Source" + multiInputTitle(id);
}

function addSource(id, title, path) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200);
            let currentID = cboxes[i].getAttribute("id")
            $.ajax({
                url: "/api/addsource",
                type: "POST",
                data: {
                    id: currentID,
                    title: title,
                    path: path,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    // 기존 Sources 추가된다.
                    let source = "";
                    if (path.startsWith("http")) {
                        source = `<div id="source-${data.id}-${data.title}"><a href="${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
                    } else {
                        source = `<div id="source-${data.id}-${data.title}"><a href="${data.protocol}://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}" onclick="copyClipboard('${data.path}')">${data.title}</a></div>`;
                    }
                    document.getElementById("sources-"+data.id).innerHTML = document.getElementById("sources-"+data.id).innerHTML + source;
                    // 요소갯수에 따라 버튼을 설정한다.
                    if (document.getElementById(`sources-${data.id}`).childElementCount > 0) {
                        document.getElementById("source-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="(setAddSourceModal('${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmsource" onclick="setRmSourceModal('${data.id}')">－</span>
                        `
                    } else {
                        document.getElementById("source-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.id}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
            
        }
    } else {
        $.ajax({
            url: "/api/addsource",
            type: "POST",
            data: {
                id: id,
                title: title,
                path: path,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                // 기존 Sources 추가된다.
                let source = "";
                if (path.startsWith("http")) {
                    source = `<div id="source-${data.id}-${data.title}"><a href="${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
                } else {
                    source = `<div id="source-${data.id}-${data.title}"><a href="${data.protocol}://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}" onclick="copyClipboard('${data.path}')">${data.title}</a></div>`;
                }
                document.getElementById("sources-"+data.id).innerHTML = document.getElementById("sources-"+data.id).innerHTML + source;
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`sources-${data.id}`).childElementCount > 0) {
                    document.getElementById("source-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmsource" onclick="setRmSourceModal('${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("source-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.id}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setRmSourceModal(id) {
    document.getElementById("modal-rmsource-id").value = id;
    document.getElementById("modal-rmsource-subject").value = "";
    document.getElementById("modal-rmsource-title").innerHTML = "Rm Source" + multiInputTitle(id);
}

function rmSource(id, title) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        var cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200);
            currentID = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/rmsource",
                type: "POST",
                data: {
                    id: currentID,
                    title: title,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`source-${data.id}-${data.title}`).remove();
                    if (document.getElementById(`sources-${data.id}`).childElementCount > 0) {
                        document.getElementById("source-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmsource" onclick="setRmSourceModal('${data.id}')>－</span>
                        `
                    } else {
                        document.getElementById("source-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.id}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
            
        }
    } else {
        $.ajax({
            url: "/api/rmsource",
            type: "POST",
            data: {
                id: id,
                title: title,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`source-${data.id}-${data.title}`).remove();
                if (document.getElementById(`sources-${data.id}`).childElementCount > 0) {
                    document.getElementById("source-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmsource" onclick="setRmSourceModal('${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("source-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.id}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setAddReferenceModal(id) {
    document.getElementById("modal-addreference-id").value = id;
    document.getElementById("modal-addreference-subject").value = "";
    document.getElementById("modal-addreference-path").value = "";
    document.getElementById("modal-addreference-title").innerHTML = "Add Reference" + multiInputTitle(id);
}

function addReference(id, title, path) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200);
            let currentID = cboxes[i].getAttribute("id")
            $.ajax({
                url: "/api/addreference",
                type: "POST",
                data: {
                    id: currentID,
                    title: title,
                    path: path,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    // 기존 References 추가된다.
                    let ref = "";
                    if (path.startsWith("http")) {
                        ref = `<div id="reference-${data.id}-${data.title}"><a href="${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
                    } else {
                        ref = `<div id="reference-${data.id}-${data.title}"><a href="${data.protocol}://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}" onclick="copyClipboard('${data.path}')">${data.title}</a></div>`;
                    }
                    document.getElementById("references-"+data.id).innerHTML = document.getElementById("references-"+data.id).innerHTML + ref;
                    // 요소갯수에 따라 버튼을 설정한다.
                    if (document.getElementById(`references-${data.id}`).childElementCount > 0) {
                        document.getElementById("reference-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmreference" onclick="setRmReferenceModal('${data.id}')">－</span>
                        `
                    } else {
                        document.getElementById("reference-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.id}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
            
        }
    } else {
        $.ajax({
            url: "/api/addreference",
            type: "POST",
            data: {
                id: id,
                title: title,
                path: path,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                // 기존 References 추가된다.
                let ref = "";
                if (path.startsWith("http")) {
                    ref = `<div id="reference-${data.id}-${data.title}"><a href="${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
                } else {
                    ref = `<div id="reference-${data.id}-${data.title}"><a href="${data.protocol}://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}" onclick="copyClipboard('${data.path}')">${data.title}</a></div>`;
                }
                document.getElementById("references-"+data.id).innerHTML = document.getElementById("references-"+data.id).innerHTML + ref;
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`references-${data.id}`).childElementCount > 0) {
                    document.getElementById("reference-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference"  onclick="setAddReferenceModal('${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmreference"  onclick="setRmReferenceModal('${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("reference-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.id}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setRmReferenceModal(id) {
    document.getElementById("modal-rmreference-id").value = id;
    document.getElementById("modal-rmreference-subject").value = "";
    document.getElementById("modal-rmreference-title").innerHTML = "Rm Source" + multiInputTitle(id);
}

function rmReference(id, title) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        var cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200);
            currentID = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/rmreference",
                type: "POST",
                data: {
                    id: currentID,
                    title: title,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`reference-${data.id}-${data.title}`).remove();
                    if (document.getElementById(`references-${data.id}`).childElementCount > 0) {
                        document.getElementById("reference-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmreference" onclick="setRmReferenceModal('${data.id}')">－</span>
                        `
                    } else {
                        document.getElementById("reference-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.id}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
            
        }
    } else {
        $.ajax({
            url: "/api/rmreference",
            type: "POST",
            data: {
                id: id,
                title: title,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`reference-${data.id}-${data.title}`).remove();
                if (document.getElementById(`references-${data.id}`).childElementCount > 0) {
                    document.getElementById("reference-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmreference" onclick="setRmReferenceModal('${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("reference-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.id}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setSeq(seq) {
    $.ajax({
        url: "/api/setseq",
        type: "POST",
        data: {
            id: document.getElementById('modal-iteminfo-id').value,
            seq: seq,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            return data;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setScene(scene) {
    $.ajax({
        url: "/api/setscene",
        type: "POST",
        data: {
            id: document.getElementById('modal-iteminfo-id').value,
            scene: scene,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            return data;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCut(cut) {
    $.ajax({
        url: "/api/setcut",
        type: "POST",
        data: {
            id: document.getElementById('modal-iteminfo-id').value,
            cut: cut,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            return data;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setSeason(season) {
    $.ajax({
        url: "/api/setseason",
        type: "POST",
        data: {
            project: document.getElementById('modal-iteminfo-project').value,
            id: document.getElementById('modal-iteminfo-id').value,
            season: season,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            return data;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setEpisode(episode) {
    $.ajax({
        url: "/api/setepisode",
        type: "POST",
        data: {
            id: document.getElementById('modal-iteminfo-id').value,
            episode: episode,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            return data;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setPlatePath(path) {
    $.ajax({
        url: "/api/setplatepath",
        type: "POST",
        data: {
            id: document.getElementById('modal-iteminfo-id').value,
            path: path,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            return data;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setThummov(path) {
    $.ajax({
        url: "/api/setthummov",
        type: "POST",
        data: {
            id: document.getElementById('modal-iteminfo-id').value,
            path: path,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("button-thumbplay-"+data.id).innerHTML = `<a href="${data.protocol}://${data.path}" class="play">PLAY</a>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setBeforemov(path) {
    $.ajax({
        url: "/api/setbeforemov",
        type: "POST",
        data: {
            id: document.getElementById('modal-iteminfo-id').value,
            path: path,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            //console.info(data);
            return
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setAftermov(path) {
    $.ajax({
        url: "/api/setaftermov",
        type: "POST",
        data: {
            id: document.getElementById('modal-iteminfo-id').value,
            path: path,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            //console.info(data);
            return
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setEditmov(path) {
    $.ajax({
        url: "/api/seteditmov",
        type: "POST",
        data: {
            id: document.getElementById('modal-iteminfo-id').value,
            path: path,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            //console.info(data);
            return
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setRetimeplate(id, path) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/setretimeplate",
        type: "POST",
        data: {
            id: id,
            path: path,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.path === "") {
                document.getElementById("button-retime-"+data.id).innerHTML = "";
            } else {
                document.getElementById("button-retime-"+data.id).innerHTML = `<a href="${data.protocol}://${data.path}" class="badge badge-danger">R</a>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setScanname(id, scanname) {
    $.ajax({
        url: "/api/setscanname",
        type: "POST",
        data: {
            id: id,
            scanname: scanname,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.id}-scanname`).innerHTML = data.scanname;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setIteminfoModal(project, id) {
    document.getElementById("modal-iteminfo-project").value = project;
    document.getElementById("modal-iteminfo-id").value = id;
    document.getElementById("modal-iteminfo-title").innerHTML = "Iteminfo" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic " + token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-iteminfo-name').value = data.name;
            document.getElementById('modal-iteminfo-type').value = data.type;
            document.getElementById('modal-iteminfo-seq').value = data.seq;
            document.getElementById('modal-iteminfo-scene').value = data.scene;
            document.getElementById('modal-iteminfo-cut').value = data.cut;
            document.getElementById('modal-iteminfo-episode').value = data.episode;
            document.getElementById('modal-iteminfo-platepath').value = data.platepath;
            document.getElementById('modal-iteminfo-thummov').value = data.thummov;
            document.getElementById('modal-iteminfo-beforemov').value = data.beforemov;
            document.getElementById('modal-iteminfo-aftermov').value = data.aftermov;
            document.getElementById('modal-iteminfo-editmov').value = data.editmov;
            document.getElementById('modal-iteminfo-retimeplate').value = data.retimeplate;
            document.getElementById('modal-iteminfo-scanname').value = data.scanname;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setTaskMov(id, task, mov) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api2/settaskmov",
        type: "POST",
        data: {
            id: id,
            task: task,
            mov: mov,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.mov === "") {
                document.getElementById(`${data.id}-task-${data.task}-playbutton`).innerHTML = "";
            } else {
                document.getElementById(`${data.id}-task-${data.task}-playbutton`).innerHTML = `<a class="mt-1 ml-1 badge badge-light" href="${data.protocol}://${data.mov}">▶</a>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}



function setTaskUser() {
    let id = document.getElementById('modal-edittask-id').value
    let task = document.getElementById('modal-edittask-task').value
    let user = document.getElementById('modal-edittask-user').value
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            fetch('/api/settaskuser', {
                method: 'POST',
                headers: {
                    "Authorization": "Basic "+ document.getElementById("token").value,
                },
                body: new URLSearchParams({
                    id: id,
                    task: task,
                    user: user,
                })
            })
            .then((response) => {
                if (!response.ok) {
                    throw Error(response.statusText + " - " + response.url);
                }
                return response.json()
            })
            .then((data) => {
                if (data.usernameandteam === "") {
                    document.getElementById(`${data.id}-task-${data.task}-user`).innerHTML = "";
                } else {
                    document.getElementById(`${data.id}-task-${data.task}-user`).innerHTML = `<span class="mt-1 ml-1 badge badge-light">${data.usernameandteam}</span>`;
                }
            })
            .catch((err) => {
                alert(err)
            });
        }
    } else {
        fetch('/api/settaskuser', {
            method: 'POST',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
            body: new URLSearchParams({
                id: id,
                task: task,
                user: user,
            })
        })
        .then((response) => {
            if (!response.ok) {
                throw Error(response.statusText + " - " + response.url);
            }
            return response.json()
        })
        .then((data) => {
            if (data.usernameandteam === "") {
                document.getElementById(`${data.id}-task-${data.task}-user`).innerHTML = "";
            } else {
                document.getElementById(`${data.id}-task-${data.task}-user`).innerHTML = `<span class="mt-1 ml-1 badge badge-light">${data.usernameandteam}</span>`;
            }
        })
        .catch((err) => {
            alert(err)
        });
    }
}

function setTaskUserComment() {
    let token = document.getElementById("token").value;
    let id = document.getElementById('modal-edittask-id').value
    let task = document.getElementById('modal-edittask-task').value
    let usercomment = document.getElementById('modal-edittask-usercomment').value
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/settaskusercomment",
                type: "POST",
                data: {
                    id: id,
                    task: task,
                    usercomment: usercomment,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    if (data.username === "") {
                        document.getElementById(`${data.id}-task-${data.task}-usercomment`).innerHTML = "";
                    } else {
                        document.getElementById(`${data.id}-task-${data.task}-usercomment`).innerHTML = `<span class="mt-1 ml-1 badge badge-darkmode">${data.usercomment}</span>`;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskusercomment",
            type: "POST",
            data: {
                id: id,
                task: task,
                usercomment: usercomment,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                if (data.username === "") {
                    document.getElementById(`${data.id}-task-${data.task}-usercomment`).innerHTML = "";
                } else {
                    document.getElementById(`${data.id}-task-${data.task}-usercomment`).innerHTML = `<span class="mt-1 ml-1 badge badge-darkmode">${data.usercomment}</span>`;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}


function setTaskStatusV2(id, task, status) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (let i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api2/settaskstatus",
                type: "POST",
                data: {
                    id: id,
                    task: task,
                    status: status,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`${data.id}-task-${data.task}-status`).innerHTML = `<a class="mt-1 badge badge-${data.status} statusbox" title="${data.status}">${data.task}</a>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api2/settaskstatus",
            type: "POST",
            data: {
                id: id,
                task: task,
                status: status,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`${data.id}-task-${data.task}-status`).innerHTML = `<a class="mt-1 badge badge-${data.status} statusbox" title="${data.status}">${data.task}</a>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskDate(id, task, date) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200)
            let id = cboxes[i].getAttribute("id")
            $.ajax({
                url: "/api/settaskdate",
                type: "POST",
                data: {
                    id: id,
                    task: task,
                    date: date,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`${data.id}-task-${data.task}-date`).innerHTML = `<span class="mt-1 ml-1 badge badge-darkmode">${data.shortdate}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskdate",
            type: "POST",
            data: {
                id: id,
                task: task,
                date: date,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`${data.id}-task-${data.task}-date`).innerHTML = `<span class="mt-1 ml-1 badge badge-darkmode">${data.shortdate}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskStart(id, task, date) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200)
            let id = cboxes[i].getAttribute("id")
            $.ajax({
                url: "/api/settaskstart",
                type: "POST",
                data: {
                    id: id,
                    task: task,
                    date: date,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    //console.log(data.date);
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskstart",
            type: "POST",
            data: {
                id: id,
                task: task,
                date: date,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.log(data.date);
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}


function setTaskUserNote(id, task, usernote) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200)
            let id = cboxes[i].getAttribute("id")
            $.ajax({
                url: "/api/settaskusernote",
                type: "POST",
                data: {
                    id: id,
                    task: task,
                    usernote: usernote,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    console.info(data)
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskusernote",
            type: "POST",
            data: {
                id: id,
                task: task,
                usernote: usernote,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}



function setTaskEnd(id, task, date) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (let i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200)
            $.ajax({
                url: "/api/settaskend",
                type: "POST",
                data: {
                    id: id,
                    task: task,
                    date: date,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    if (data.shortdate === "") {
                        document.getElementById(`${data.id}-task-${data.task}-end`).innerHTML = "";
                    } else {
                        document.getElementById(`${data.id}-task-${data.task}-end`).innerHTML = `<span class="mt-1 ml-1 badge badge-outline-darkmode">${data.shortdate}</span>`;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskend",
            type: "POST",
            data: {
                id: id,
                task: task,
                date: date,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                if (data.shortdate === "") {
                    document.getElementById(`${data.id}-task-${data.task}-end`).innerHTML = "";
                } else {
                    document.getElementById(`${data.id}-task-${data.task}-end`).innerHTML = `<span class="mt-1 ml-1 badge badge-outline-darkmode">${data.shortdate}</span>`;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function rfc3339toNormaltime(t) {
    if (t.includes("T")) {
        return t.split("T")[0]
    }
    return t
}

// setDeadline2dModal 함수는 project, id 정보를 이용해서 Deadline2d Modal 값을 채운다.
function setDeadline2dModal(id) {
    document.getElementById("modal-deadline2d-id").value = id;
    document.getElementById("modal-deadline2d-title").innerHTML = "Set Deadline 2D" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic " + token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-deadline2d-date').value = rfc3339toNormaltime(data.ddline2d);
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setDeadline2D(id, date) {
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            fetch('/api/setdeadline2d', {
                method: 'POST',
                headers: {
                    "Authorization": "Basic "+ document.getElementById("token").value,
                },
                body: new URLSearchParams({
                    id: id,
                    date: date,
                })
            })
            .then((response) => {
                if (!response.ok) {
                    throw Error(response.statusText + " - " + response.url);
                }
                return response.json()
            })
            .then((data) => {
                document.getElementById("deadline2d-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline2d" onclick="setDeadline2dModal('${data.id}')">2D:${data.shortdate}</span>`;
            })
            .catch((err) => {
                alert(err)
            });
        }
    } else {
        fetch('/api/setdeadline2d', {
            method: 'POST',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
            body: new URLSearchParams({
                id: id,
                date: date,
            })
        })
        .then((response) => {
            if (!response.ok) {
                throw Error(response.statusText + " - " + response.url);
            }
            return response.json()
        })
        .then((data) => {
            document.getElementById("deadline2d-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline2d" onclick="setDeadline2dModal('${data.id}')">2D:${data.shortdate}</span>`;
        })
        .catch((err) => {
            alert(err)
        });
    }
}

// setDeadline3dModal 함수는 project, id 정보를 이용해서 Deadline3d Modal 값을 채운다.
function setDeadline3dModal(id) {
    document.getElementById("modal-deadline3d-id").value = id;
    document.getElementById("modal-deadline3d-title").innerHTML = "Set Deadline 3D" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic " + token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-deadline3d-date').value = rfc3339toNormaltime(data.ddline3d);
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setDeadline3D(id, date) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setdeadline3d",
                type: "POST",
                data: {
                    id: id,
                    date: date,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("deadline3d-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline3d" onclick="setDeadline3dModal('${data.id}')">3D:${data.shortdate}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setdeadline3d",
            type: "POST",
            data: {
                id: id,
                date: date,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("deadline3d-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline3d" onclick="setDeadline3dModal('${data.id}')">3D:${data.shortdate}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setShottype(id) {
    let token = document.getElementById("token").value;
    let e = document.getElementById("modal-shottype-type");
    let shottype = e.options[e.selectedIndex].value;
    
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setshottype",
                type: "POST",
                data: {
                    id: id,
                    shottype: shottype,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("shottype-"+data.name).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#modal-shottype" onclick="setShottypeModal('${data.id}')">${data.type}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setshottype",
            type: "POST",
            data: {
                id: id,
                shottype: shottype,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("shottype-"+data.id).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#modal-shottype" onclick="setShottypeModal('${data.id}')">${data.type}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setUsetype(id) {
    let token = document.getElementById("token").value;
    let e = document.getElementById("modal-usetype-type");
    let type = e.options[e.selectedIndex].value;
    
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setusetype",
                type: "POST",
                data: {
                    id: id,
                    type: type,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`${data.project}-${data.id}-usetype`).innerHTML = `<span class="badge badge-warning ml-1" data-toggle="modal" data-target="#modal-usetype" onclick="setUsetypeModal('${data.id}')">${data.type}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setusetype",
            type: "POST",
            data: {
                id: id,
                type: type,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`${data.id}-usetype`).innerHTML = `<span class="badge badge-warning ml-1" data-toggle="modal" data-target="#modal-usetype" onclick="setUsetypeModal('${data.id}')">${data.type}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setAssettypeModal(id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-assettype-title").innerHTML = "Assettype Type" + multiInputTitle(id);
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-assettype-id').value=id;
            document.getElementById("modal-assettype-type").value=data.assettype;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setAssettype(id) {
    let token = document.getElementById("token").value;
    let types = document.getElementById("modal-assettype-type");
    let assettype = types.options[types.selectedIndex].value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setassettype",
                type: "POST",
                data: {
                    id: id,
                    assettype: assettype,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    // assettype button update
                    document.getElementById("assettype-"+data.id).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#modal-assettype" onclick="setAssettypeModal('${data.id}')">${data.type}</span>`;
                    // remove old assettype tag
                    document.getElementById(`assettag-${data.id}-${data.oldtype}`).remove();
                    // add new assettype tag
                    let url = `/inputmode?project=${data.project}&searchword=assettags:${data.type}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&task=`;
                    source = `<div id="tag-${data.id}-${data.type}"><a href="${url}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.type}</a></div>`;
                    document.getElementById("assettags-"+data.id).innerHTML = document.getElementById("assettags-"+data.id).innerHTML + source;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setassettype",
            type: "POST",
            data: {
                id: id,
                assettype: assettype,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                // assettype button update
                document.getElementById("assettype-"+data.id).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#modal-assettype" onclick="setAssettypeModal('${data.id}')">${data.type}</span>`;
                // remove old assettype tag
                document.getElementById(`assettag-${data.id}-${data.oldtype}`).remove();
                // add new assettype tag
                let url = `/inputmode?project=${data.project}&searchword=assettags:${data.type}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&task=`;
                source = `<div id="tag-${data.name}-${data.type}"><a href="${url}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.type}</a></div>`;
                document.getElementById("assettags-"+data.id).innerHTML = document.getElementById("assettags-"+data.id).innerHTML + source;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

// setRnumModal 함수는 project, id 정보를 이용해서 Edit Rnum Modal 값을 채운다.
function setRnumModal(id) {
    document.getElementById("modal-setrnum-id").value = id;
    document.getElementById("modal-setrnum-title").innerHTML = "Set Rnum number" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-setrnum-text').value = data.rnum;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setRnum() {
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            fetch('/api2/setrnum', {
                method: 'POST',
                headers: {
                    "Authorization": "Basic "+ document.getElementById("token").value,
                },
                body: new URLSearchParams({
                    id: id,
                    rnum: document.getElementById('modal-setrnum-text').value,
                })
            })
            .then((response) => {
                return response.json()
            })
            .then((data) => {
                if (data.rnum !== "") {
                    document.getElementById("rnum-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-setrnum" onclick="setModal('modal-setrnum-text', '${data.rnum}' );setModal('modal-setrnum-id', '${data.id}')"{{end}}>${data.rnum}</span>`;
                } else {
                    document.getElementById("rnum-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-setrnum" onclick="setModal('modal-setrnum-text', '${data.rnum}' );setModal('modal-setrnum-id', '${data.id}')"{{end}}>no rnum</span>`;
                }
            })
            .catch((error) => {
                alert(error)
            });
        }
    } else {
        fetch('/api2/setrnum', {
            method: 'POST',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
            body: new URLSearchParams({
                id: document.getElementById('modal-setrnum-id').value,
                rnum: document.getElementById('modal-setrnum-text').value,
            })
        })
        .then((response) => {
            return response.json()
        })
        .then((data) => {
            if (data.rnum !== "") {
                document.getElementById("rnum-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-setrnum" onclick="setModal('modal-setrnum-text', '${data.rnum}' );setModal('modal-setrnum-id', '${data.id}')"{{end}}>${data.rnum}</span>`;
            } else {
                document.getElementById("rnum-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-setrnum" onclick="setModal('modal-setrnum-text', '${data.rnum}' );setModal('modal-setrnum-id', '${data.id}')"{{end}}>no rnum</span>`;
            }
        })
        .catch((error) => {
            alert(error)
        });
    }
}

function setAddTagModal(id) {
    document.getElementById("modal-addtag-id").value = id;
    document.getElementById("modal-addtag-text").value = "";
    document.getElementById("modal-addtag-title").innerHTML = "Add Tag" + multiInputTitle(id);
}

function setAddAssetTagModal(id) {
    document.getElementById("modal-addassettag-id").value = id;
    document.getElementById("modal-addassettag-text").value = "";
    document.getElementById("modal-addassettag-title").innerHTML = "Add Asset Tag" + multiInputTitle(id);
}

function setRenameTagModal(project) {
    document.getElementById("modal-renametag-project").value = project;
    document.getElementById("modal-renametag-beforetag").value = "";
    document.getElementById("modal-renametag-aftertag").value = "";
}

function addTag(id, tag) {
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            fetch('/api/addtag', {
                method: 'POST',
                headers: {
                    "Authorization": "Basic "+ document.getElementById("token").value,
                },
                body: new URLSearchParams({
                    id: id,
                    tag: tag,
                })
            })
            .then((response) => {
                if (!response.ok) {
                    throw Error(response.statusText + " - " + response.url);
                }
                return response.json()
            })
            .then((data) => {
                // 기존 Tags에 추가된다.
                let url = `/inputmode?project=${data.project}&searchword=tag:${data.tag}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&task=`
                source = `<div id="tag-${data.id}-${data.tag}"><a href="${url}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.tag}</a></div>`;
                document.getElementById("tags-"+data.id).innerHTML = document.getElementById("tags-"+data.id).innerHTML + source;
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`tags-${data.id}`).childElementCount > 0) {
                    document.getElementById("tag-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("tag-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.id}')">＋</span>
                    `
                }
            })
            .catch((err) => {
                alert(err)
            });
        }
    } else {
        fetch('/api/addtag', {
            method: 'POST',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
            body: new URLSearchParams({
                id: id,
                tag: tag,
            })
        })
        .then((response) => {
            if (!response.ok) {
                throw Error(response.statusText + " - " + response.url);
            }
            return response.json()
        })
        .then((data) => {
            // 기존 Tags에 추가된다.
            let url = `/inputmode?project=${data.project}&searchword=tag:${data.tag}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&task=`
            let source = `<div id="tag-${data.id}-${data.tag}"><a href="${url}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.tag}</a></div>`;
            document.getElementById("tags-"+data.id).innerHTML = document.getElementById("tags-"+data.id).innerHTML + source;
            // 요소갯수에 따라 버튼을 설정한다.
            if (document.getElementById(`tags-${data.id}`).childElementCount > 0) {
                document.getElementById("tag-button-"+data.id).innerHTML = `
                <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.id}')">＋</span>
                <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.id}')">－</span>
                `
            } else {
                document.getElementById("tag-button-"+data.id).innerHTML = `
                <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.id}')">＋</span>
                `
            }
        })
        .catch((err) => {
            alert(err)
        });
    }
}



function renameTag(project, before, after) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/renametag",
        type: "POST",
        
        data: {
            project: project,
            before: before,
            after: after,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            alert(`The "${data.before}" tag in the "${data.project}" project has been changed to "${data.after}" tag.\nPlease refresh the page.`);
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setRmTagModal(id) {
    document.getElementById("modal-rmtag-id").value = id;
    document.getElementById("modal-rmtag-tag").value = "";
    document.getElementById("modal-rmtag-title").innerHTML = "Rm Tag" + multiInputTitle(id);
    document.getElementById("modal-rmtag-iscontain").value = false;
}

function setRmAssetTagModal(id) {
    document.getElementById("modal-rmassettag-id").value = id;
    document.getElementById("modal-rmassettag-tag").value = "";
    document.getElementById("modal-rmassettag-title").innerHTML = "Rm Asset Tag" + multiInputTitle(id);
    document.getElementById("modal-rmassettag-iscontain").value = false;
}

function rmTag() {
    let id = document.getElementById('modal-rmtag-id').value
    let tag = document.getElementById('modal-rmtag-tag').value
    let isContain = document.getElementById('modal-rmtag-iscontain').checked
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            fetch('/api/rmtag', {
                method: 'POST',
                headers: {
                    "Authorization": "Basic "+ document.getElementById("token").value,
                },
                body: new URLSearchParams({
                    id: id,
                    tag: tag,
                    iscontain: isContain,
                })
            })
            .then((response) => {
                if (!response.ok) {
                    throw Error(response.statusText + " - " + response.url);
                }
                return response.json()
            })
            .then((data) => {
                if (isContain) {
                    document.querySelectorAll(`[id^="tag-${data.id}-"][id*="${data.tag}"]`).forEach(el => el.remove());
                } else {
                    document.getElementById(`tag-${data.id}-${data.tag}`).remove();
                }
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`tags-${data.id}`).childElementCount > 0) {
                    document.getElementById("tag-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.id}')">－</span>
                    `;
                } else {
                    document.getElementById("tag-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.id}')">＋</span>
                    `;
                }
            })
            .catch((err) => {
                alert(err)
            });
        }
    } else {
        fetch('/api/rmtag', {
            method: 'POST',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
            body: new URLSearchParams({
                id: id,
                tag: tag,
                iscontain: isContain,
            })
        })
        .then((response) => {
            if (!response.ok) {
                throw Error(response.statusText + " - " + response.url);
            }
            return response.json()
        })
        .then((data) => {
            if (isContain) {
                document.querySelectorAll(`[id^="tag-${data.id}-"][id*="${data.tag}"]`).forEach(el => el.remove());
            } else {
                document.getElementById(`tag-${data.id}-${data.tag}`).remove();
            }
            // 요소갯수에 따라 버튼을 설정한다.
            if (document.getElementById(`tags-${data.id}`).childElementCount > 0) {
                document.getElementById("tag-button-"+data.id).innerHTML = `
                <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.id}')">＋</span>
                <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.id}')">－</span>
                `;
            } else {
                document.getElementById("tag-button-"+data.id).innerHTML = `
                <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.id}')">＋</span>
                `;
            }
        })
        .catch((err) => {
            alert(err)
        });
    }
}

function isMultiInput() {
    var cboxes = document.getElementsByName('selectID');
    for (let i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === true) {
            return true;
        }
    }
    return false;
}


function selectCheckbox(id) {
    let beforeNum = document.querySelectorAll('input[name=selectID]:checked').length;
    if (document.getElementById(id).checked) {
        if ((beforeNum - 1) === 0) {
            document.getElementById("topbtn").innerHTML = "Top"
        } else {
            document.getElementById("topbtn").innerHTML = "Top<br>" + (beforeNum - 1)
        }
    } else {
        if ((beforeNum + 1) === 0) {
            document.getElementById("topbtn").innerHTML = "Top"
        } else {
            document.getElementById("topbtn").innerHTML = "Top<br>" + (beforeNum + 1)
        }
        
    }
    // 선택된 아이템의 숫자를 갱신한다.
}

function selectCheckboxAll() {
    let cboxes = document.getElementsByName('selectID');
    for (let i = 0; i < cboxes.length; ++i) {
        cboxes[i].checked = true;
    }
    document.getElementById("topbtn").innerHTML = "Top<br>" + (document.querySelectorAll('input[name=selectID]:checked').length)
}

function selectCheckboxNone() {
    let cboxes = document.getElementsByName('selectID');
    for (let i = 0; i < cboxes.length; ++i) {
        cboxes[i].checked = false;
    }
    document.getElementById("topbtn").innerHTML = "Top"
}

function selectCheckboxInvert() {
    let cboxes = document.getElementsByName('selectID');
    let invertNum = 0
    for (let i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            cboxes[i].checked = true;
            invertNum += 1
        } else {
            cboxes[i].checked = false;
        }
    }
    document.getElementById("topbtn").innerHTML = "Top<br>" + (invertNum)
}


function setObjectIDModal(project, id) {
    document.getElementById("modal-objectid-project").value = project;
    document.getElementById("modal-objectid-id").value = id;
    document.getElementById("modal-objectid-title").innerHTML = "Object ID" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-objectid-in').value = data.objectidin;
            document.getElementById('modal-objectid-out').value = data.objectidout;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setFinver(id, version) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/setfinver",
        type: "POST",
        data: {
            id: id,
            version: version,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("finver-"+data.id).innerHTML = `<span>v${data.version}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setObjectID(project, id, innum, outnum) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/setobjectid",
        type: "POST",
        data: {
            project: project,
            name: id2name(id),
            in: innum,
            out: outnum,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("objectidnum-"+data.name).innerHTML = `<span class="text-badge ml-1">${data.in}-${data.out}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCrowdAsset(project, id) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/setcrowdasset",
        type: "POST",
        data: {
            project: project,
            name: id2name(id),
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.crowdasset) {
                document.getElementById("crowdasset-"+data.id).innerHTML = `<span class="badge badge-warning finger" onclick="setCrowdAsset('${data.project}', '${data.id}')">Crowd</span>`;
            } else {
                document.getElementById("crowdasset-"+data.id).innerHTML = `<span class="badge badge-light fade finger" onclick="setCrowdAsset('${data.project}', '${data.id}')">Crowd</span>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setPlatesizeModal(id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-platesize-title").innerHTML = "Platesize" + multiInputTitle(id);
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-platesize-id').value = id;
            document.getElementById("modal-platesize-size").value = data.platesize;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setPlatesize(id, size) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setplatesize",
                type: "POST",
                data: {
                    id: id,
                    size: size,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("platesize-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-platesize" onclick="setPlatesizeModal('${data.id}')">S:${data.size}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setplatesize",
            type: "POST",
            data: {
                id: id,
                size: size,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("platesize-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-platesize" onclick="setPlatesizeModal('${data.id}')">S:${data.size}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setUndistortionsizeModal(id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-undistortionsize-title").innerHTML = "Undistortionsize" + multiInputTitle(id);
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-undistortionsize-id').value = id;
            document.getElementById("modal-undistortionsize-size").value = data.undistortionsize;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setUndistortionsize(id, size) {
    let token = document.getElementById("token").value;

    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setundistortionsize",
                type: "POST",
                data: {
                    id: id,
                    size: size,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("undistortionsize-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-undistortionsize" onclick="setUndistortionsizeModal('${data.id}', '${data.size}')">U:${data.size}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setundistortionsize",
            type: "POST",
            
            data: {
                id: id,
                size: size,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("undistortionsize-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-undistortionsize" onclick="setUndistortionsizeModal('${data.id}', '${data.size}')">U:${data.size}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setRendersizeModal(id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-rendersize-title").innerHTML = "Rendersize" + multiInputTitle(id);
    $.ajax({
        url: `/api2/item?id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-rendersize-id').value = id;
            document.getElementById("modal-rendersize-size").value = data.rendersize;
            document.getElementById("modal-rendersize-overscanratio").value = data.overscanratio;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setRendersize() {
    let id = document.getElementById('modal-rendersize-id').value;
    let size = document.getElementById('modal-rendersize-size').value;
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api2/setrendersize",
                type: "POST",
                
                data: {
                    id: id,
                    size: size,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("rendersize-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-rendersize" onclick="setRendersizeModal('${data.id}')">R:${data.size}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api2/setrendersize",
            type: "POST",
            
            data: {
                id: id,
                size: size,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("rendersize-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-rendersize" onclick="setRendersizeModal('${data.id}')">R:${data.size}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setOverscanRatio() {
    let id = document.getElementById('modal-rendersize-id').value;
    let ratio = document.getElementById('modal-rendersize-overscanratio').value;
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setoverscanratio",
                type: "POST",
                data: {
                    id: id,
                    ratio: ratio,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("overscanratio-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-rendersize" onclick="setRendersizeModal('${data.id}')">OSR:${data.overscanratio}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setoverscanratio",
            type: "POST",
            
            data: {
                id: id,
                ratio: ratio,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("overscanratio-"+data.id).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-rendersize" onclick="setRendersizeModal('${data.id}')">OSR:${data.overscanratio}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function CurrentProject() {
    let e = document.getElementById("searchbox-project");
    return e.options[e.selectedIndex].value;
}

function rmItem() {
    let token = document.getElementById("token").value;
    let cboxes = document.getElementsByName('selectID');
    let selectNum = 0;
    for (let i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === true) {
            selectNum += 1
        }
    }
    if (selectNum > 0) {
        for (let i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/rmitemid",
                type: "POST",
                data: {

                    id: id,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`item-${data.id}`).remove();
                    document.getElementById("topbtn").innerHTML = "Top" // 삭제가되면 Top버튼에서 선택된 아이템 갯수를 리셋한다.
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        alert("삭제할 아이템을 선택해주세요.");
    }
}

function autocomplete(inp) {
    let arr
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/autocompliteusers",
        type: "get",
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            arr = data.users;
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
    /*the autocomplete function takes two arguments,
    the text field element and an array of possible autocompleted values:*/
    var currentFocus;
    /*execute a function when someone writes in the text field:*/
    inp.addEventListener("input", function(e) {
        var a, b, i, val = this.value; // a:유저를 표기하는 박스, b:유저를 표기하는 박스 내부의 한 요소, i: for문에 필요한 인수, val: input창에 입력한 값
        /*이미 열려있는 리스트를 닫는다.*/
        closeAllLists();
        if (!val) { return false;}
        currentFocus = -1;
        // DIV 하나를 생성한다.
        a = document.createElement("DIV");
        a.setAttribute("id", this.id + "autocomplete-list");
        a.setAttribute("class", "autocomplete-items");
        /*위에서 생성한 검색창을 부모에 붙힌다.*/
        this.parentNode.appendChild(a);
        /*각각의 아이템을 순환한다.*/
        for (i = 0; i < arr.length; i++) {
          /*검색어가 아이템에 포함되어 있다면, div를 생성한다.*/
          if (arr[i].searchword.includes(val)) {
            /*create a DIV element for each matching element:*/
            b = document.createElement("DIV");
            let userInfo = arr[i].id + "(" + arr[i].name + "," + arr[i].team + ")";
            /*make the matching letters bold:*/
            let otherList = userInfo.split(val)
            b.innerHTML = otherList[0] + "<span class='text-warning'>" + val + "</span>" + otherList[1];
            /*insert a input field that will hold the current array item's value:*/
            b.innerHTML += `<input type='hidden' value='${userInfo}'>`;
            /*execute a function when someone clicks on the item value (DIV element):*/
            b.addEventListener("click", function(e) {
                /*insert the value for the autocomplete text field:*/
                inp.value = this.getElementsByTagName("input")[0].value;
                /*close the list of autocompleted values,
                (or any other open lists of autocompleted values:*/
                closeAllLists();
            });
            a.appendChild(b);
          }
        }
    });
    /*execute a function presses a key on the keyboard:*/
    inp.addEventListener("keydown", function(e) {
        var x = document.getElementById(this.id + "autocomplete-list");
        if (x) x = x.getElementsByTagName("div");
        if (e.keyCode == 40) {
          /*If the arrow DOWN key is pressed,
          increase the currentFocus variable:*/
          currentFocus++;
          /*and and make the current item more visible:*/
          addActive(x);
        } else if (e.keyCode == 38) { //up
          /*If the arrow UP key is pressed,
          decrease the currentFocus variable:*/
          currentFocus--;
          /*and and make the current item more visible:*/
          addActive(x);
        } else if (e.keyCode == 13) {
          /*If the ENTER key is pressed, prevent the form from being submitted,*/
          e.preventDefault();
          if (currentFocus > -1) {
            /*and simulate a click on the "active" item:*/
            if (x) x[currentFocus].click();
          }
        }
    });
    function addActive(x) {
      /*a function to classify an item as "active":*/
      if (!x) return false;
      /*start by removing the "active" class on all items:*/
      removeActive(x);
      if (currentFocus >= x.length) currentFocus = 0;
      if (currentFocus < 0) currentFocus = (x.length - 1);
      /*add class "autocomplete-active":*/
      x[currentFocus].classList.add("autocomplete-active");
    }
    function removeActive(x) {
      /*a function to remove the "active" class from all autocomplete items:*/
      for (var i = 0; i < x.length; i++) {
        x[i].classList.remove("autocomplete-active");
      }
    }
    function closeAllLists(elmnt) {
      /*close all autocomplete lists in the document,
      except the one passed as an argument:*/
      var x = document.getElementsByClassName("autocomplete-items");
      for (var i = 0; i < x.length; i++) {
        if (elmnt != x[i] && elmnt != inp) {
        x[i].parentNode.removeChild(x[i]);
      }
    }
  }
  
  /*execute a function when someone clicks in the document:*/
  document.addEventListener("click", function (e) {
      closeAllLists(e.target);
  });
}

let input = document.getElementsByClassName("searchuser");
for (var i = 0; i < input.length; i++) {
    autocomplete(input[i])
}

function setAddTaskModal(id, type) {
    // 모달을 설정한다.
    document.getElementById("modal-addtask-id").value = id;
    document.getElementById("modal-addtask-title").innerHTML = "Add Task" + multiInputTitle(id);
    if (type === "org" || type === "main" || type === "mp" || type === "plt" || type === "plate" || type === "left") {
        // Task 셋팅
        fetch('/api/shottasksetting', {
            method: 'GET',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
        })
        .then((response) => {
            if (!response.ok) {
                throw Error(response.statusText + " - " + response.url);
            }
            return response.json()
        })
        .then((data) => {
            let tasks = data["tasksettings"];
            let addtasks = document.getElementById('modal-addtask-taskname');
            addtasks.innerHTML = "";
            for (let i = 0; i < tasks.length; i++){
                let opt = document.createElement('option');
                opt.value = tasks[i].name;
                opt.innerHTML = tasks[i].name;
                addtasks.appendChild(opt);
            }
        })
        .catch((err) => {
            alert(err)
        });
    }
    if (type === "asset") {
        // Task 셋팅
        fetch('/api/assettasksetting', {
            method: 'GET',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
        })
        .then((response) => {
            if (!response.ok) {
                throw Error(response.statusText + " - " + response.url);
            }
            return response.json()
        })
        .then((data) => {
            let tasks = data["tasksettings"]
            let addtasks = document.getElementById('modal-addtask-taskname');
            addtasks.innerHTML = "";
            for (let i = 0; i < tasks.length; i++){
                let opt = document.createElement('option');
                opt.value = tasks[i].name;
                opt.innerHTML = tasks[i].name;
                addtasks.appendChild(opt);
            }
        })
        .catch((err) => {
            alert(err)
        });
    }
}

function setRmTaskModal(id, type) {
    document.getElementById("modal-rmtask-id").value = id;
    document.getElementById("modal-rmtask-title").innerHTML = "Rm Task" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    if (type === "org" || type === "main" || type === "mp" || type === "plt" || type === "plate" || type === "left") {
        $.ajax({
            url: "/api/shottasksetting",
            type: "get",
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let tasks = data["tasksettings"];
                let rmtasks = document.getElementById('modal-rmtask-taskname');
                rmtasks.innerHTML = "";
                for (let i = 0; i < tasks.length; i++){
                    let opt = document.createElement('option');
                    opt.value = tasks[i].name;
                    opt.innerHTML = tasks[i].name;
                    rmtasks.appendChild(opt);
                }
            },
            error: function(request,status,error){
                alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
    if (type === "asset") {
        $.ajax({
            url: "/api/assettasksetting",
            type: "get",
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let tasks = data["tasksettings"]
                let rmtasks = document.getElementById('modal-rmtask-taskname');
                rmtasks.innerHTML = "";
                for (let i = 0; i < tasks.length; i++){
                    let opt = document.createElement('option');
                    opt.value = tasks[i].name;
                    opt.innerHTML = tasks[i].name;
                    rmtasks.appendChild(opt);
                }
            },
            error: function(request,status,error){
                alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function changeStatusURI(status) {
    let tags = document.getElementsByClassName("statusuri");
    for ( var i = 0; i < tags.length; i++) {
        let c = document.getElementById("searchbox-checkbox-" + status);
        if (tags[i].href.includes(status + "=true")) {
            tags[i].href = tags[i].href.replace(status + "=true", status + "=" + c.checked)
        } else {
            tags[i].href = tags[i].href.replace(status + "=false", status + "=" + c.checked)
        }
    }
}

function mailInfo(project, id) {
    $.ajax({
        url: "/api/mailinfo",
        type: "POST",
        data: {
            "project": project,
            "id": id,
            "lang": "ko",
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            let mailString = "mailto:"
            // 메일 보낼 사람을 추가한다.
            if (data.mails) {
                mailString += data.mails.join(",")
            }
            mailString += `?subject=[${data.header}] ${data.title}&`;
            // 메일을 참조할 사람을 추가한다.
            if (data.cc) {
                mailString += `cc=${data.cc.join(",")}`
            }
            // 브라우저가 크롬이라면 _blank로 열리게 함. 크롬을 사용한다면 웹메일 클라이언트를 사용할 확률이 높기 때문에 현재 작업중인 창에서 메일 작성창이 열리면 안된다.
            // https://stackoverflow.com/questions/4565112/javascript-how-to-find-out-if-the-user-browser-is-chrome/13348618
            
            if (/Chrome/.test(navigator.userAgent) && /Google Inc/.test(navigator.vendor)) {
                document.getElementById("web-mail-link").href = unescape(mailString)
                document.getElementById("web-mail-link").target = "_blank"
            } else {
                document.getElementById("web-mail-link").href = unescape(mailString)
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

/** enable listview */
function listview() {
    document.getElementById("listview").style.display="block";
    document.getElementById("page").style.display="block";
    document.getElementById("calendar").innerHTML = "";
    document.location.reload(); // 새로고침을 해야한다. 달력 간트챠트의 데이터를 새로 그리기 위해
}


function foldingmenu() {
    if(searchbox.style.display=="none") {
        searchbox.style.display='block'; // 펼치기
        document.getElementById("foldoption").innerText = "Collapse Searchbox ▲" // 글씨 변경
        setCookie("searchboxVisable", "true")// 쿠키저장
    } else {
        searchbox.style.display='none'; // 접기
        document.getElementById("foldoption").innerText = "Expand Searchbox ▼" // 글씨 변경
        setCookie("searchboxVisable", "false")// 쿠키저장
    }
    let clientSearchboxHeight = document.getElementById('floatingmenu').clientHeight;
    document.getElementById("blinkspace").style.height = clientSearchboxHeight + "px";
}

// TopClick 함수는 스크롤시 보여지는 Top 버튼을 누를 때 발생하는 이벤트이다.
function TopClick() {
    document.body.scrollTop = 0;
    document.documentElement.scrollTop = 0;
    let searchbox = document.getElementById("searchbox")
    if (searchbox !== null) {
        searchbox.style.display = "block";
        document.getElementById("blinkspace").style.height = "550px";
        document.getElementById("foldoption").innerText = "Collapse Searchbox ▲"
    }
}

function setPublish() {
    fetch('/api/setpublishstatus', {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: new URLSearchParams({
            project: document.getElementById('modal-setpublish-project').value,
            id: document.getElementById('modal-setpublish-id').value,
            task: document.getElementById('modal-setpublish-task').value,
            key: document.getElementById('modal-setpublish-key').value,
            path: document.getElementById('modal-setpublish-path').value,
            createtime: document.getElementById('modal-setpublish-createtime').value,
            status: document.getElementById("modal-setpublish-status").value,
        })
    })
    .then((response) => {
        return response.json()
    })
    .then((data) => {
        location.reload()
        return
    })
    .catch((error) => {
        alert(error)
    });
}

function addPublish() {
    fetch('/api/addpublish', {
        method: 'POST',
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify({
            id: document.getElementById('itemid').value,
            task: document.getElementById('modal-addpublish-task').value,
            key: document.getElementById('modal-addpublish-key').value,
            secondarykey: document.getElementById('modal-addpublish-secondarykey').value,
            path: document.getElementById('modal-addpublish-path').value,
            status: document.getElementById('modal-addpublish-status').value,
            tasktouse: document.getElementById('modal-addpublish-tasktouse').value,
            subject: document.getElementById('modal-addpublish-subject').value,
            mainversion: document.getElementById('modal-addpublish-mainversion').value,
            subversion: document.getElementById('modal-addpublish-subversion').value,
            filetype: document.getElementById('modal-addpublish-filetype').value,
            kindofusd: document.getElementById('modal-addpublish-kindofusd').value,
            createtime: "",
            isoutput: document.getElementById('modal-addpublish-isoutput').checked,
            outputdatapath: document.getElementById('modal-addpublish-outputdatapath').value,
        })
    })
    .then((response) => {
        if (!response.ok) {
            return response.text().then((errorText) => {
                throw new Error(errorText || `HTTP error! status: ${response.status}`);
            });
        }
        return response.json();
    })
    .then((data) => {
        location.reload()
        return
    })
    .catch((error) => {
        alert(`Error: ${error.message}`);
    });
}

function addReviewStatusMode() {
    let token = document.getElementById("token").value
    let reviewFps = document.getElementById("modal-addreview-statusmode-fps")
    $.ajax({
        url: "/api/addreview",
        type: "POST",
        data: {
            project: document.getElementById("modal-addreview-statusmode-project").value,
            name: document.getElementById("modal-addreview-statusmode-name").value,
            task: document.getElementById("modal-addreview-statusmode-task").value,
            type: document.getElementById("modal-addreview-statusmode-type").value,
            ext: document.getElementById("modal-addreview-statusmode-ext").value,
            author: document.getElementById("modal-addreview-statusmode-author").value,
            path: document.getElementById("modal-addreview-statusmode-path").value,
            description: document.getElementById("modal-addreview-statusmode-description").value,
            camerainfo: document.getElementById("modal-addreview-statusmode-camerainfo").value,
            fps: reviewFps.options[reviewFps.selectedIndex].value,
            mainversion: document.getElementById("modal-addreview-statusmode-mainversion").value,
            subversion: document.getElementById("modal-addreview-statusmode-subversion").value,
            outputdatapath: document.getElementById("modal-addreview-statusmode-outputdatapath").value,
            removeafterprocess: document.getElementById("modal-addreview-statusmode-removeafterprocess").checked,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function() {
            alert("리뷰가 정상적으로 등록되었습니다.");
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function clickReviewStatusModeCommentButton() {
    addReviewStatusModeComment()
}

function setReviewProcessStatus(id, status) {
    $.ajax({
        url: "/api/setreviewprocessstatus",
        type: "POST",
        data: {
            status: status,
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            let item = document.getElementById("reviewstatus-"+data.id)
            // 상태 내부 글씨를 바꾼다.
            item.innerHTML = data.status
            // 상태의 색상을 바꾼다.
            if (data.processstatus === "wait") {
                item.setAttribute("class","ml-1 badge badge-danger")
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewProject 함수는 리뷰데이터의 Project를 변경한다.
function setReviewProject() {
    $.ajax({
        url: "/api/setreviewproject",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            project: document.getElementById("modal-editreview-project").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.id}-project`).innerText = data.project
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewTask 함수는 리뷰데이터의 Task을 변경한다.
function setReviewTask() {
    $.ajax({
        url: "/api/setreviewtask",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            task: document.getElementById("modal-editreview-task").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.id}-task`).innerText = data.task
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewPath 함수는 리뷰데이터의 Path를 변경한다.
function setReviewPath() {
    $.ajax({
        url: "/api/setreviewpath",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            path: document.getElementById("modal-editreview-path").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            setReviewProcessStatus(data.id, "wait") // 만약 Path가 수정되면 Status가 wait으로 바뀌고 동영상이 다시 연산이 되어야 한다.
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewCreatetime 함수는 리뷰데이터의 Createtime을 변경한다.
function setReviewCreatetime() {
    $.ajax({
        url: "/api/setreviewcreatetime",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            createtime: document.getElementById("modal-editreview-createtime").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            return
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// rfc3339 함수는 date 값을 받아서 RFC3339 포멧으로 반환한다.
function rfc3339(d) {
    function pad(n) {
        return n < 10 ? "0" + n : n;
    }
    function timezoneOffset(offset) {
        var sign;
        if (offset === 0) {
            return "Z";
        }
        sign = (offset > 0) ? "-" : "+";
        offset = Math.abs(offset);
        return sign + pad(Math.floor(offset / 60)) + ":" + pad(offset % 60);
    }
    return d.getFullYear() + "-" +
        pad(d.getMonth() + 1) + "-" +
        pad(d.getDate()) + "T" +
        pad(d.getHours()) + ":" +
        pad(d.getMinutes()) + ":" +
        pad(d.getSeconds()) + 
        timezoneOffset(d.getTimezoneOffset());
}

// setReviewCreatetimeNow 함수는 리뷰데이터의 시간을 현재시간으로 설정한다.
function setReviewCreatetimeNow() {
    let date = new Date()
    let time = rfc3339(date)
    $.ajax({
        url: "/api/setreviewcreatetime",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            createtime: time,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("modal-editreview-createtime").value = time
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
    
}

// setReviewMainVersion 함수는 리뷰데이터의 MainVersion을 변경한다.
function setReviewMainVersion() {
    $.ajax({
        url: "/api/setreviewmainversion",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            mainversion: document.getElementById("modal-editreview-mainversion").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.id}-mainversion`).innerText = "v" + data.mainversion
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewSubVersion 함수는 리뷰데이터의 SubVersion을 변경한다.
function setReviewSubVersion() {
    $.ajax({
        url: "/api/setreviewsubversion",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            subversion: document.getElementById("modal-editreview-subversion").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.id}-subversion`).innerText = "v" + data.subversion
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewFps 함수는 리뷰데이터의 Fps를 변경한다.
function setReviewFps() {
    $.ajax({
        url: "/api/setreviewfps",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            fps: document.getElementById("modal-editreview-fps").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            return
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewDescription함수는 리뷰데이터의 Description을 변경한다.
function setReviewDescription() {
    $.ajax({
        url: "/api/setreviewdescription",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            description: document.getElementById("modal-editreview-description").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("description").innerText = data.description
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewCameraInfo 함수는 리뷰데이터의 CameraInfo를 변경한다.
function setReviewCameraInfo() {
    $.ajax({
        url: "/api/setreviewcamerainfo",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            camerainfo: document.getElementById("modal-editreview-camerainfo").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            return
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewOutputDataPath 함수는 리뷰데이터의 OutputDataPath를 변경한다.
function ReviewOutputDataPath() {
    fetch('/api/reviewoutputdatapath', {
        method: 'PATCH',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: new URLSearchParams({
            id: document.getElementById("modal-editreview-id").value,
            outputdatapath: document.getElementById("modal-editreview-outputdatapath").value,
        })
    })
    .then((response) => {
        return response.json()
    })
    .then((data) => {
    })
    .catch((error) => {
        alert(error)
    });
}

// setReviewName 함수는 리뷰데이터의 Name을 변경한다.
function setReviewName() {
    $.ajax({
        url: "/api/setreviewname",
        type: "POST",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            name: document.getElementById("modal-editreview-name").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.id}-name`).innerText = data.name
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}





function addReviewCommentText(text) {
    $.ajax({
        url: "/api/addreviewcomment",
        type: "POST",
        data: {
            id: document.getElementById("current-review-id").value,
            text: text,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            // 데이터가 잘 들어가면 review-comments 에 들어간 데이터를 드로잉한다.
            let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
            let newComment = `<div id="reviewcomment-${data.id}-${data.date}" class="p-1">
            <span class="text-badge">${data.date} / <a href="/user?id=${data.author}" class="text-darkmode">${data.authorname}</a></span>
            <span class="edit" data-toggle="modal" data-target="#modal-editreviewcomment" onclick="setEditReviewCommentModal('${data.id}', '${data.date}')">≡</span>
            <span class="remove" data-toggle="modal" data-target="#modal-rmreviewcomment" onclick="setRmReviewCommentModal('${data.id}','${data.date}')">×</span>
            <br><small class="text-white">${body}</small>`
            if (data.media != "") {
                if (data.media.includes("http")) {
                    newComment += `<div class="row pl-3 pt-3 pb-1">
                        <a href="${data.media}" onclick="copyClipboard('${data.media}')">
                            <img src="/assets/img/link.svg" class="finger">
                        </a>
                        <span class="text-white pl-2 small">${data.mediatitle}</span>
                    </div>`
                } else {
                    newComment += `<div class="row pl-3 pt-3 pb-1">
                        <a href="${data.protocol}://${data.media}" onclick="copyClipboard('${data.media}')">
                            <img src="/assets/img/link.svg" class="finger">
                        </a>
                        <span class="text-white pl-2 small">${data.mediatitle}</span>
                    </div>`
                }
            }
            newComment += `<hr class="my-1 p-0 m-0 divider"></hr></div>`
            document.getElementById("review-comments").innerHTML = newComment + document.getElementById("review-comments").innerHTML;
            document.getElementById("review-comment").value = ""; // 입력한 값을 초기화 한다.
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setTypeAddShot(type, readOnly) {
    document.getElementById("addshot-type").value = type
    document.getElementById('addshot-type').readOnly = readOnly
}

// 리뷰를 위해서 동영상에 그림을 그리기 위해 필요한 글로벌 변수를 셋팅한다.
var drawCanvas, drawCtx;
var mouseStartX=0
var mouseStartY=0
var drawing = false;
var globalClientWidth = 0;
var globalClientHeight = 0;
var globalReviewRenderWidth = 0;
var globalReviewRenderHeight = 0;
var globalReviewRenderWidthOffset = 0;
var globalReviewRenderHeightOffset = 0;
var framelineOffset = 0;
var frameLineMarkHeight = 12; // 프레임 표시라인 높이

function initCanvas() {
    let playerbox = document.getElementById("playerbox"); // player 캔버스를담을 div를 가지고 온다.
    globalClientWidth = playerbox.clientWidth // 클라이언트 사용자의 가로 사이즈를 구한다.
    globalClientHeight = playerbox.clientHeight // 클라이언트 사용자의 세로 사이즈를 구한다.
    // Player 캔버스를 초기화 한다.
    let playerCanvas = document.getElementById("player");
    playerCanvas.setAttribute("width", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    playerCanvas.setAttribute("height", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    playerCanvas.setAttribute("width", globalClientWidth) // 캔버스를 클라이언트 사용자의 가로사이즈로 설정한다.
    playerCanvas.setAttribute("height", globalClientHeight) // 캔버스를 클라이언트 사용자의 세로사이즈로 설정한다.
    // Draw 캔버스를 초기화 한다.
    let drawCanvas = document.getElementById("drawcanvas");
    drawCanvas.setAttribute("width", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    drawCanvas.setAttribute("height", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    drawCanvas.setAttribute("width", globalClientWidth) // 그림을 그리는 캔버스 가로 사이즈를 설정한다.
    drawCanvas.setAttribute("height", globalClientHeight) // 그림을 그리는 캔버스 세로 사이즈를 설정한다.
    // UX 캔버스를 초기화 한다.
    let uxCanvas = document.getElementById("uxcanvas");
    uxCanvas.setAttribute("width", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    uxCanvas.setAttribute("height", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    uxCanvas.setAttribute("width", globalClientWidth) // UX 캔버스 가로 사이즈를 설정한다.
    uxCanvas.setAttribute("height", globalClientHeight) // UX 캔버스 세로 사이즈를 설정한다.
    // Animation UX 캔버스를 초기화 한다.
    let aniuxCanvas = document.getElementById("aniuxcanvas");
    aniuxCanvas.setAttribute("width", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    aniuxCanvas.setAttribute("height", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    aniuxCanvas.setAttribute("width", globalClientWidth) // Animation UX 캔버스 가로 사이즈를 설정한다.
    aniuxCanvas.setAttribute("height", globalClientHeight) // Animation UX 캔버스 세로 사이즈를 설정한다.
    // Screenshot 캔버스를 초기화 한다.
    let screenshotCanvas = document.getElementById("screenshot");
    screenshotCanvas.setAttribute("width", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    screenshotCanvas.setAttribute("height", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    screenshotCanvas.setAttribute("width", globalClientWidth) // 스크린샷 캔버스 가로 사이즈를 설정한다.
    screenshotCanvas.setAttribute("height", globalClientHeight) // 스크린샷 캔버스 세로 사이즈를 설정한다.
}

function selectReviewItem(id) {
    let project, fps, ext, type;
    let playerbox = document.getElementById("playerbox"); // player 캔버스를담을 div를 가지고 온다.
    let clientWidth = playerbox.clientWidth // 클라이언트 사용자의 가로 사이즈를 구한다.
    let clientHeight = playerbox.clientHeight // 클라이언트 사용자의 세로 사이즈를 구한다.
    initCanvas();
    let playerCanvas = document.getElementById("player");
    let playerCtx = playerCanvas.getContext("2d");
    let drawCanvas = document.getElementById("drawcanvas");
    drawCtx = drawCanvas.getContext("2d")
    let uxCanvas = document.getElementById("uxcanvas");
    uxCtx = uxCanvas.getContext("2d")
    let aniuxCanvas = document.getElementById("aniuxcanvas");
    aniuxCtx = aniuxCanvas.getContext("2d")
    let screenshotCanvas = document.getElementById("screenshot");
    screenshotCtx = screenshotCanvas.getContext("2d")

    // 비디오객체의 메타데이터를 로딩하면 실행할 함수를 설정한다.
    let totalFrame = 0
    let sketchesFrame = [];
    // 기존에 드로잉 되어 있는 데이터를 가지고 온다.
    $.ajax({
        url: "/api/review",
        type: "POST",
        data: {
            id: id,
        },
        async: false,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            project = data.project
            fps = data.fps
            ext = data.ext
            type = data.type
            for (let i = 0; i < data.sketches.length; i++) {
                sketchesFrame.push(data.sketches[i].frame)
            }
        },
        error: function(){
            return
        }
    })
    // 입력받은 프로젝트로 웹페이지의 Review Title을 변경한다.
    document.title = "Review: " + project;
    // 브러쉬 설정
    drawCtx.lineWidth = 4; // 브러시 사이즈
    drawCtx.strokeStyle = "#EFEAD6" // 브러시 컬러

    // 마우스 이벤트 처리
    drawCanvas.addEventListener("mousemove", function (e) {move(e)},false);
    drawCanvas.addEventListener("mousedown", function (e) {down(e)}, false);
    drawCanvas.addEventListener("mouseup", function (e) {up(e)}, false);
    drawCanvas.addEventListener("mouseout", function (e) {out(e)}, false);

    // 버튼설정 및 버튼 이벤트
    let playButton = document.getElementById("player-play");
    let pauseButton = document.getElementById("player-pause");
    let playAndPauseButton = document.getElementById("player-playandpause");
    let loopAndLoofOffButton = document.getElementById("player-loopandloopoff");
    let startButton = document.getElementById("player-start");
    let endButton = document.getElementById("player-end");
    let beforeFrameButton = document.getElementById("player-left");
    let afterFrameButton = document.getElementById("player-right");
    let gotoFrameInput = document.getElementById("modal-gotoframe-frame");
    let prevDrawing = document.getElementById("drawing-prev");
    let nextDrawing = document.getElementById("drawing-next");

    // GotoFrame 모달창에서 프레임이 변경되면 해당 프레임으로 이동한다.
    gotoFrameInput.addEventListener("change", function() {
        targetFrame = document.getElementById("modal-gotoframe-frame").value
        video.currentTime = gotoFrame(targetFrame, fps) // video.currentTime이 바뀌기 때문에 video.addEventListener('timeupdate', function () {}) 이벤트가 발생해서 드로잉이 띄워진다.
    });

    // prev Drawing 버튼을 클릭할 때 이벤트
    prevDrawing.addEventListener("click", function() {
        $.ajax({
            url: "/api/reviewdrawingframe",
            type: "POST",
            data: {
                id: id,
                frame: document.getElementById("currentframe").innerHTML,
                mode: "prev",
            },
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value
            },
            dataType: "json",
            success: function(data) {
                video.currentTime = gotoFrame(data.resultframe, fps)
            },
            error: function(){
                return
            }
        })
    });

    // next Drawing 버튼을 클릭할 때 이벤트
    nextDrawing.addEventListener("click", function() {
        $.ajax({
            url: "/api/reviewdrawingframe",
            type: "POST",
            data: {
                id: id,
                frame: document.getElementById("currentframe").innerHTML,
                mode: "next",
            },
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value
            },
            dataType: "json",
            success: function(data) {
                video.currentTime = gotoFrame(data.resultframe, fps)
            },
            error: function(){
                return
            }
        })
    });
    // 플레이 버튼을 클릭할 때 이벤트
    playButton.addEventListener("click", function() {
        playAndPauseButton.className = "player-pause"
        video.play();
    });
    // 일시정지 버튼을 클릭할 때 이벤트
    pauseButton.addEventListener("click", function() {
        playAndPauseButton.className = "player-play"
        video.pause();
    });

    // 재생과 정지가 같이 진행되는 버튼
    playAndPauseButton.addEventListener("click", function() {
        if (!video.paused) {
            playAndPauseButton.className = "player-play"
            video.pause();    
        } else {
            playAndPauseButton.className = "player-pause"
            video.play();
        }
    });
    // Loop버튼 클릭시
    loopAndLoofOffButton.addEventListener("click", function() {
        if (loopAndLoofOffButton.className == "player-loop") {
            loopAndLoofOffButton.className = "player-loopoff"
            video.loop = false
        } else {
            loopAndLoofOffButton.className = "player-loop"
            video.loop = true
        }
    });
    startButton.addEventListener("click", function() {
        if (fps == 25) {
            video.currentTime = video.seekable.start(0)
        } else {
            video.currentTime = video.seekable.start(0) + (1/parseFloat(fps))
        }
    }); // 처음으로 이동하는 버튼을 클릭할 때 이벤트. 맨 앞에서 1/fps만큼 이후로 이동해야 함.
    endButton.addEventListener("click", function() {
        if (fps == 60 || fps == 23.976) {
            video.currentTime = video.seekable.end(0)
        } else {
            // 끝으로 이동하는 버튼을 클릭할 때 이벤트. 맨 뒤에서 1/fps만큼 이전으로 이동해야 함.
            video.currentTime = video.seekable.end(0) - (1/parseFloat(fps))
        }
    });
    // 이전 프레임으로 이동하는 버튼을 클릭할 때 이벤트
    beforeFrameButton.addEventListener("click", function() {
        if (fps == 25) {
            // 25fps를 가지고 있는 미디어는 시작 프레임이 2프레임에서 시작된다.
            if (video.currentTime > 0.0) {
                video.currentTime -= (1/parseFloat(fps));
            } else {
                video.currentTime = video.seekable.start(0)
            }
        } else {
            if (video.currentTime > 0.0) {
                video.currentTime -= (1/parseFloat(fps));
            } else {
                video.currentTime = video.seekable.start(0) + (1/parseFloat(fps))
            }
        }
    }); 
    // 다음 프레임으로 이동하는 버튼을 클릭할 때 이벤트
    afterFrameButton.addEventListener("click", function() {
        if (video.currentTime < video.seekable.end(0)) {
            video.currentTime += (1/parseFloat(fps));
        } else {
            video.currentTime = video.seekable.end(0) - (1/parseFloat(fps))
        }
    });

    // video 객체를 생성한다.
    var video = document.createElement('video');
    video.src = `/reviewdata?id=${id}&ext=${ext}`;
    video.autoplay = true;
    video.loop = true;
    video.defaultMuted = false; // Sound 처리
    video.muted = false; // Sound 처리
    video.volume = 0.5; // Sound 처리
    video.setAttribute("id", "currentvideo");

    // 이미지 객체를 생성한다.
    if (type === "image") {
        reviewImage = new Image();
        reviewImage.src = `/reviewdata?id=${id}&ext=${ext}`;
        reviewImage.onload = function(){
            renderWidth = (this.width * clientHeight) / this.height // 실제로 렌더링되는 너비
            renderHeight = (this.height * clientWidth) / this.width // 실제로 렌더링되는 높이
            if (clientWidth <= renderWidth && renderHeight < clientHeight) {
                // 가로형: 가로비율이 맞고, 높이가 적을 때
                let hOffset = (clientHeight - renderHeight) / 2
                globalReviewRenderWidth = clientWidth
                globalReviewRenderHeight = renderHeight
                globalReviewRenderHeightOffset = hOffset
                playerCtx.drawImage(this, 0, hOffset, clientWidth, renderHeight);
            } else {
                // 세로형: 가로비율이 작고 높이가 맞을 때
                let wOffset = (clientWidth - renderWidth) / 2
                globalReviewRenderWidth = renderWidth
                globalReviewRenderHeight = clientHeight
                globalReviewRenderWidthOffset = wOffset
                playerCtx.drawImage(this, wOffset, 0, renderWidth, clientHeight);
            }
        }
        // 기존 스케치가 있을 수 있다. 드로잉을 지운다.
        removeDrawing()
        
        // 서버에 드로잉이 존재하면 fg 캔버스에 그린다.
        let drawing = new Image()
        let frame = document.getElementById("currentframe").innerHTML
        let url = `/reviewdrawingdata?id=${id}&frame=${frame}&time=${new Date().getTime()}`
        let http = new XMLHttpRequest();
        http.open("HEAD", url, false)
        http.send()
        if (http.status === 200) {
            let fg = document.getElementById("drawcanvas")
            let fgctx = fg.getContext("2d")
            drawing.src = url
            drawing.onload = function() {
                fgctx.drawImage(drawing,
                    0, 0, drawing.width, drawing.height,
                    globalReviewRenderWidthOffset, globalReviewRenderHeightOffset, globalReviewRenderWidth, globalReviewRenderHeight
                );
            };
        } else {
            return
        }
    }
        
    // 플레이창 배경을 검정색으로 채운다.
    playerCtx.fillStyle = "#000000";
    playerCtx.fillRect(0, 0, clientWidth, clientHeight);
    
    // 비디오가 로딩되면 메타데이터로 처리할 수 있는 과정을 처리한다.
    video.onloadedmetadata = function() {
        // Draw 캔버스에 프레임 표기 그림을 그린다.
        totalFrame = Math.round(this.duration * parseFloat(fps)) // round로 해야 23.976fps에서 frame 에러가 발생하지 않는다.
        // totalFrame을 표기한다.
        document.getElementById("totalframe").innerHTML = padNumber(totalFrame);
        // 프레임 표기바의 간격을 구하고 global 변수에 저장한다.
        framelineOffset = clientWidth / totalFrame
        
        // 프레임바를 드로잉 한다. 스케치가 있다면 노란색바를 드로잉한다.
        for (let i = 0; i < totalFrame; i++) {
            uxCtx.beginPath();
            if (sketchesFrame.includes(i+1)) {                
                uxCtx.strokeStyle = '#FFCD31';
            } else {
                uxCtx.strokeStyle = '#333333';
            }
            uxCtx.lineWidth = 2;
            uxCtx.moveTo(i*framelineOffset + (framelineOffset / 2) , clientHeight - frameLineMarkHeight);
            uxCtx.lineTo(i*framelineOffset + (framelineOffset / 2), clientHeight);
            uxCtx.stroke();
            uxCtx.closePath();
        }
        
        // 재생에 필요한 모든 설정이 완료되면 리뷰 데이터를 플레이시킨다.
        playAndPauseButton.className = "player-pause"
        video.play();
    };
    
    video.addEventListener('play', function () {
        let $this = this; //cache
        (function loop() {
            if (!$this.paused && !$this.ended) {
                renderWidth = ($this.videoWidth * clientHeight) / $this.videoHeight // 실제로 렌더링되는 너비
                renderHeight = ($this.videoHeight * clientWidth) / $this.videoWidth // 실제로 렌더링되는 높이
                if (clientWidth <= renderWidth && renderHeight < clientHeight) {
                    // 가로형: 가로비율이 맞고, 높이가 적을 때
                    let hOffset = (clientHeight - renderHeight) / 2
                    globalReviewRenderWidth = clientWidth
                    globalReviewRenderHeight = renderHeight
                    globalReviewRenderHeightOffset = hOffset
                    playerCtx.drawImage($this, 0, hOffset, clientWidth, renderHeight);
                } else {
                    // 세로형: 가로비율이 작고 높이가 맞을 때
                    let wOffset = (clientWidth - renderWidth) / 2
                    globalReviewRenderWidth = renderWidth
                    globalReviewRenderHeight = clientHeight
                    globalReviewRenderWidthOffset = wOffset
                    playerCtx.drawImage($this, wOffset, 0, renderWidth, clientHeight);
                }
                // fps에 맞게 currentFrame을 드로잉한다.
                let currentFrame = Math.floor($this.currentTime * parseFloat(fps))
                if (currentFrame < totalFrame) {
                    document.getElementById("currentframe").innerHTML = padNumber(currentFrame + 1)
                } else {
                    // 재생이 멈추면 표기는 totalFrame이 되어야 하지만 실제 재생시점은 영상의 마지막이 되어야 한다.
                    document.getElementById("currentframe").innerHTML = padNumber(totalFrame)
                }
                // 커서의 위치를 드로잉 한다.
                aniuxCtx.clearRect(0, 0, clientWidth, clientHeight);
                aniuxCtx.strokeStyle = "#FF0000";
                aniuxCtx.lineWidth = 4;
                aniuxCtx.beginPath();
                aniuxCtx.moveTo(currentFrame * framelineOffset + (framelineOffset/2), clientHeight - frameLineMarkHeight);
                aniuxCtx.lineTo(currentFrame * framelineOffset + (framelineOffset/2), clientHeight);
                aniuxCtx.stroke();

                // 다음화면 갱신
                setTimeout(loop, 1000 / parseFloat(fps));
            }
        })();
    }, 0);

    video.addEventListener('timeupdate', function () {
        let $this = this; //cache 화
        if (clientWidth <= globalReviewRenderWidth && globalReviewRenderHeight < clientHeight) {
            // 가로형: 가로비율이 맞고, 높이가 적을 때
            let hOffset = (clientHeight - globalReviewRenderHeight) / 2
            playerCtx.drawImage($this, 0, hOffset, clientWidth, globalReviewRenderHeight);
        } else {
            // 세로형: 가로비율이 작고 높이을 꽉 채울 때
            let wOffset = (clientWidth - globalReviewRenderWidth) / 2
            playerCtx.drawImage($this, wOffset, 0, globalReviewRenderWidth, clientHeight);
        }
        // fps에 맞게 currentFrame을 드로잉한다.
        let currentFrame = Math.floor($this.currentTime * parseFloat(fps))
        if (currentFrame < totalFrame) {
            document.getElementById("currentframe").innerHTML = padNumber(currentFrame + 1)
        } else {
            // 재생이 멈추면 표기는 totalFrame이 되어야 하지만 실제 재생시점은 영상의 마지막이 되어야 한다.
            document.getElementById("currentframe").innerHTML = padNumber(totalFrame)
        }
        // 빨간 커서의 위치를 드로잉 한다.
        aniuxCtx.clearRect(0, 0, clientWidth, clientHeight);
        aniuxCtx.strokeStyle = "#FF0000";
        aniuxCtx.lineWidth = 4;
        aniuxCtx.beginPath();
        aniuxCtx.moveTo(currentFrame * framelineOffset + (framelineOffset/2), clientHeight - frameLineMarkHeight);
        aniuxCtx.lineTo(currentFrame * framelineOffset + (framelineOffset/2), clientHeight);
        aniuxCtx.stroke();
        
        // 드로잉 프레임은 비디오가 정지될 때만 보여야한다.
        if (video.paused) {
            // 프레임을 이동하면 기존 드로잉이 지워져야 한다.
            removeDrawing()
            // 드로잉이 존재하면 fg 캔버스에 그린다.
            let drawing = new Image()
            let frame = document.getElementById("currentframe").innerHTML
            let url = `/reviewdrawingdata?id=${id}&frame=${frame}&time=${new Date().getTime()}`
            let http = new XMLHttpRequest();
            http.open("HEAD", url, false)
            http.send()
            if (http.status === 200) {
                let fg = document.getElementById("drawcanvas")
                let fgctx = fg.getContext("2d")
                drawing.src = url
                drawing.onload = function() {
                    fgctx.drawImage(drawing,
                        0, 0, drawing.width, drawing.height,
                        globalReviewRenderWidthOffset, globalReviewRenderHeightOffset, globalReviewRenderWidth, globalReviewRenderHeight
                    );
                };
            } else {
                return
            }
        }
        
    }, 0);
}

function gotoFrame(frame, fps) {
    return ((parseFloat(frame) - (1.0/parseFloat(fps))) / parseFloat(fps))
}

// draw 함수는 x,y 좌표를 받아 그림을 그린다.
function draw(curX, curY) {
    drawCtx.beginPath();
    drawCtx.moveTo(mouseStartX, mouseStartY);
    drawCtx.lineTo(curX, curY);
    drawCtx.stroke();
}

// down 함수는 마우스를 클릭할 때 현재 위치를 마우스의 시작좌표로 설정하고 그림을 그린다고 설정한다.
function down(e) {
    mouseStartX = e.offsetX;
    mouseStartY = e.offsetY;
    drawing = true;
}

// up 함수는 마우스 버튼을 땔 때 그림그리는 모드를 종료한다.
function up(e) {
    drawing = false;
    saveDrawing()
}

// move 함수는 그림을 그리는 상태이고 마우스가 이동할 때 현재 위치에 그림을 그리고 현재위치를 다시 마우스의 시작위치로 바꾼다.
function move(e) {
    if (!drawing) {
        return; // 마우스가 눌러지지 않았으면 리턴
    }
    var curX = e.offsetX, curY = e.offsetY;
    draw(curX, curY);
    mouseStartX = curX;
    mouseStartY = curY;
}

// out 함수는 화면 밖으로 커서가 나가면 그림그리는 모드를 끈다.
function out(e) {
    drawing = false;
}

// changeYellowDrawingFrame 함수는 프레임을 받아서 노란색 드로잉 마커를 체크한다.
function changeYellowDrawingFrame() {
    let uxCanvas = document.getElementById("uxcanvas");
    uxCtx = uxCanvas.getContext("2d")
    currentFrame = parseInt(document.getElementById("currentframe").innerHTML) - 1
    uxCtx.beginPath();
    uxCtx.strokeStyle = '#FFCD31';
    uxCtx.lineWidth = 2;
    uxCtx.moveTo(currentFrame*framelineOffset + (framelineOffset / 2) , globalClientHeight - frameLineMarkHeight);
    uxCtx.lineTo(currentFrame*framelineOffset + (framelineOffset / 2), globalClientHeight);
    uxCtx.stroke();
    uxCtx.closePath();
}

// changeDrawingGrayFrame 함수는 프레임을 받아서 노란색 드로잉 마커를 체크한다.
function changeDrawingGrayFrame() {
    let uxCanvas = document.getElementById("uxcanvas");
    uxCtx = uxCanvas.getContext("2d")
    currentFrame = parseInt(document.getElementById("currentframe").innerHTML) - 1
    uxCtx.beginPath();
    uxCtx.strokeStyle = '#333333';
    uxCtx.lineWidth = 2;
    uxCtx.moveTo(currentFrame*framelineOffset + (framelineOffset / 2) , globalClientHeight - frameLineMarkHeight);
    uxCtx.lineTo(currentFrame*framelineOffset + (framelineOffset / 2), globalClientHeight);
    uxCtx.stroke();
    uxCtx.closePath();
}


// screenshot 함수는 리뷰중인 스크린을 스크린샷 합니다.
function screenshot(filename) {
    let screenshot = document.getElementById("screenshot");
    let fg = document.getElementById("drawcanvas");
    let bg = document.getElementById("player");
    let screenshotctx = screenshot.getContext('2d')
    screenshotctx.drawImage(bg, 0, 0) // 배경에 bg를 그린다.
    screenshotctx.drawImage(fg, 0, 0) // 배경에 fg를 그린다.
    let dataURL = screenshot.toDataURL("image/png")
    let link = document.createElement('a');
    link.href = dataURL;
    let today = new Date();
    let y = today.getFullYear().toString();
    let m = ("0"+(today.getMonth()+1)).slice(-2);
    let d = ("0"+(today.getDate())).slice(-2);
    let hour = ("0"+(today.getHours())).slice(-2);
    let min = ("0"+(today.getMinutes())).slice(-2);
    let sec = ("0"+(today.getSeconds())).slice(-2);
    let timestamp = y + m + d + 'T' + hour + min + sec; // ISO 8601 format 
    let currentFrame = document.getElementById("currentframe").innerHTML + 'f';
    link.download = filename + '_'+ currentFrame + '_' + timestamp + '.png';
    link.setAttribute("type","hidden") // firefox에서는 꼭 DOM구조를 지켜야 한다.
    document.body.appendChild(link); // firefox에서는 꼭 DOM구조를 지켜야 한다.
    link.click();
    link.remove();
    // 스크린샷이 저장되면 기존에 캔버스에 합성된 이미지를 제거한다.
    let playerbox = document.getElementById("playerbox");
    let clientWidth = playerbox.clientWidth
    let clientHeight = playerbox.clientHeight
    screenshotctx.clearRect(0, 0, clientWidth, clientHeight);
    changeYellowDrawingFrame()
}

// saveDrawing 함수는 리뷰스크린에 드로잉된 이미지를 서버에 저장합니다.
function saveDrawing() {
    let id = document.getElementById("current-review-id").value;
    let token = document.getElementById("token").value;
    // Crop Canvas를 생성한다.
    let cropCanvas = document.createElement("canvas");
    cropCanvas.id = "cropCanvas";
    cropCanvas.width = globalReviewRenderWidth;
    cropCanvas.height = globalReviewRenderHeight;
    let ctx = cropCanvas.getContext('2d');
    let drawingCanvas = document.getElementById("drawcanvas")
    ctx.drawImage(drawingCanvas, globalReviewRenderWidthOffset, globalReviewRenderHeightOffset, globalReviewRenderWidth, globalReviewRenderHeight, 0, 0, globalReviewRenderWidth, globalReviewRenderHeight)

    // canvas의 드로잉을 .png 파일로 파일화 한다.
    let fg = cropCanvas.toDataURL("image/png");
    let blobBin = atob(fg.split(',')[1]); // base64 데이터를 바이너리로 변경한다.
    let array = [];
    for (let i = 0; i < blobBin.length; i++) {
        array.push(blobBin.charCodeAt(i));
    }
    let file = new Blob([new Uint8Array(array)], {type: 'image/png'}); // Blob 생성

    let formData = new FormData();
    formData.append("file", file); // .png 추가
    formData.append("id", id)
    formData.append("frame", document.getElementById("currentframe").innerHTML)
    $.ajax({
        url: "/api/uploadreviewdrawing",
        type: "POST",
        enctype: "multipart/form-data",
        processData: false,
        contentType : false,
        cache: false,
        data: formData,
        headers: {
            "Authorization": "Basic "+ token
        },
        success: function(data) {
            changeYellowDrawingFrame()
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

// removeDrawing 함수는 리뷰 스케치를 제거합니다. 프레임을 갱신할 때 사용합니다.
function removeDrawing() {
    // drawcanvas를 지운다.
    let fg = document.getElementById("drawcanvas");
    let fgctx = fg.getContext("2d");
    fgctx.clearRect(0, 0, globalClientWidth, globalClientHeight);
    // screenshot를 지운다.
    let screenshot = document.getElementById("screenshot");
    let screenshotctx = screenshot.getContext("2d");
    screenshotctx.clearRect(0, 0, globalClientWidth, globalClientHeight);
}

// removeDrawingAndData 함수는 리뷰 스케치를 제거하고 서버의 이미지도 제거합니다.
function removeDrawingAndData() {
    // drawcanvas를 지운다.
    let fg = document.getElementById("drawcanvas");
    let fgctx = fg.getContext("2d");
    fgctx.clearRect(0, 0, globalClientWidth, globalClientHeight);
    // screenshot를 지운다.
    let screenshot = document.getElementById("screenshot");
    let screenshotctx = screenshot.getContext("2d");
    screenshotctx.clearRect(0, 0, globalClientWidth, globalClientHeight);
    // 서버에 파일이 존재하면 삭제한다.
    $.ajax({
        url: "/api/rmreviewdrawing",
        type: "POST",
        data: {
            id: document.getElementById("current-review-id").value,
            frame: parseInt(document.getElementById("currentframe").innerHTML),
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            // 파일이 잘 삭제되면, 그림이 그려진 프레임의 노란바를 회색바로 변경한다.
            changeDrawingGrayFrame()
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

// copyClipboard 는 value 값을 받아서, 클립보드로 복사하는 기능이다.
function copyClipboard(value) {
    let id = document.createElement("input");   // input요소를 만듬
    id.setAttribute("value", value);            // input요소에 값을 추가
    document.body.appendChild(id);              // body에 요소 추가
    id.select();                                // input요소를 선택
    document.execCommand("copy");               // 복사기능 실행
    document.body.removeChild(id);              // body에 요소 삭제
}

// copyClipboardAndMessage 는 value 값을 받아서, 클립보드로 복사하는 기능이다. 이후 메시지를 출력한다.
function copyClipboardAndMessage(value) {
    let id = document.createElement("input");   // input요소를 만듬
    id.setAttribute("value", value);            // input요소에 값을 추가
    document.body.appendChild(id);              // body에 요소 추가
    id.select();                                // input요소를 선택
    document.execCommand("copy");               // 복사기능 실행
    document.body.removeChild(id);              // body에 요소 삭제
    alert(value + "\n값이 클립보드에 복사되었습니다.");
}

// 리뷰페이지 핫키
// Hotkey: http://gcctech.org/csc/javascript/javascript_keycodes.htm
document.onkeydown = function(e) {
    // 인풋창에서는 화살표를 움직였을 때 페이지가 이동되면 안된다.
    if (event.target.tagName === "INPUT") {
        return
    }
    if (event.target.tagName === "TEXTAREA") {
        return
    }
    if (e.which == 37) { // arrow left
        document.getElementById("player-pause").click();
        document.getElementById("player-left").click();
    } else if (e.which == 39) { // arrow right
        document.getElementById("player-pause").click();
        document.getElementById("player-right").click();
    } else if (e.which == 80 || e.which == 83 || e.which == 32) { // p, s, space
        document.getElementById("player-playandpause").click();
    } else if (e.which == 219) { // [
        document.getElementById("player-pause").click();
        document.getElementById("player-start").click();
    } else if (e.which == 221) { // ]
        document.getElementById("player-pause").click();
        document.getElementById("player-end").click();
    } else if (e.which == 84) { // t
        document.getElementById("player-trash").click();
    } else if (e.which == 67) { // c
        document.getElementById("player-screenshot").click();
    } else if (e.which == 76) { // l
        document.getElementById("player-loopandloopoff").click();
    } else if (e.which == 190) { // .
        document.getElementById("drawing-next").click();
    } else if (e.which == 188) { // ,
        document.getElementById("drawing-prev").click();
    }
};

function rvplay(id) {
    // review id의 데이터를 가지고 path값을 구하고 프로토콜을 통해 rv player에 연결한다.
    $.ajax({
        url: "/api/review",
        type: "POST",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            let link = document.createElement('a');
            link.href = document.getElementById("protocol").value + "://" + data.path;
            document.body.appendChild(link); // firefox에서는 꼭 DOM구조를 지켜야 한다.
            link.click();
            link.remove();
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

// playOriginal 함수는 project와 name을 받아서 해당 아이템의 썸네일동영상을 재생한다. 원본 플레이트 영상을 재생하기 위해서 사용한다.
function playOriginal(project, name) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api2/item?project=${project}&name=${name}`,
        type: "get",
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            let link = document.createElement('a');
            link.href = document.getElementById("protocol").value + "://" + data.thummov; // 썸네일 동영상이 보통은 썸네일을 플레이할 때의 Original Plate 이다.
            document.body.appendChild(link); // firefox에서는 꼭 DOM구조를 지켜야 한다.
            link.click();
            link.remove();
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

function rmUser() {
    let id = document.getElementById("modal-rmuser-id").value
    $.ajax({
        url: "/api2/user?id="+id,
        type: "DELETE",
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            location.reload()
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setModalGotoFrame() {
    document.getElementById("modal-gotoframe-frame").value = parseInt(document.getElementById("currentframe").innerHTML);
}

// redirectPage 함수는 page를 받아서 해당 페이지로 리다이렉트 한다.
function redirectPage(page) {
    let href = new URL(window.location.href);
    href.searchParams.set('page', page);
    let url = href.toString();
    window.location.href = url
}

function selectUserID(id) {
    if (document.getElementById(id).style.borderColor === SELECT_COLOR) {
        document.getElementById(id).style.borderColor = NON_SELECT_COLOR
        document.getElementById(id).style.backgroundColor = "rgba(0,0,0,0)";
    } else {
        document.getElementById(id).style.borderColor = SELECT_COLOR
        document.getElementById(id).style.backgroundColor = "rgba(255, 200, 0, 0.1)";
    }
}

function initPasswordUsers() {
    // 선택된 사용자를 출력한다.
    let usercards = document.getElementsByClassName("usercard");
    let users = new Array();
    // 사용자가 선택되었다면 users Array에 넣는다.
    for (let i = 0; i < usercards.length; i++) {
        if (document.getElementById(usercards[i].id).style.borderColor === SELECT_COLOR) {
            users.push(usercards[i].id);
        }
    }
    // 초기화할 사용자가 없다면 종료한다.
    if (users.length === 0) {
        alert(`패스워드를 초기화할 사용자를 선택해주세요.`);
        return;
    }
    // 선택된 각각의 유저를 초기화 한다.
    for (let i = 0; i < users.length; i++) {
        const id = users[i];
        const token = document.getElementById("token").value;
    
        fetch("/api/initpassword", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Basic " + token
            },
            body: JSON.stringify({ id: id })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            // 성공하면 원래 색상으로 돌린다.
            const userElement = document.getElementById(data.id);
            if (userElement) {
                userElement.style.borderColor = NON_SELECT_COLOR;
                userElement.style.backgroundColor = "rgba(0,0,0,0)";
            }
            alert(`${data.id}'s password has been reset.`);
        })
        .catch(error => {
            alert(`Error: ${error.message}`);
        });
    }
    
}

function setReviewAgainForWaitStatusToday() {
    $.ajax({
        url: "/api/setreviewagainforwaitstatustoday",
        type: "POST",
        data: {},
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            alert(`${data.userid}에 의해 ${data.num}개의 Wait 상태 리뷰데이터를 오늘 리뷰항목으로 설정했습니다.`);
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setReviewItemStatus(itemstatus) {
    fetch('/api/setreviewitemstatus', {
        method: 'post',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: new URLSearchParams({
            id: document.getElementById("current-review-id").value,
            itemstatus: itemstatus,
        })
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        // 해당 id의 status 글씨와 색상을 바꾼다.
        let itemStatus = document.getElementById("review-itemstatus-"+data.id)
        itemStatus.innerHTML = data.itemstatus
        itemStatus.setAttribute("class","ml-1 badge badge-"+data.itemstatus)
        // 현재 띄워진 화면의 좌측 Status 상태를 변경한다.
        document.getElementById("current-review-itemstatus").value = data.itemstatus
    })
    .catch((err) => {
        alert(err)
    });
}

function addReviewStatusModeComment() {
    fetch('/api/addreviewstatusmodecomment', {
        method: 'post',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: new URLSearchParams({
            id: document.getElementById("current-review-id").value,
            text: document.getElementById("review-comment").value,
            media: document.getElementById("review-media").value,
            itemstatus: document.getElementById("current-review-itemstatus").value,
            frame: document.getElementById("currentframe").innerHTML,
            framecomment: document.getElementById("review-framecomment").checked,
        })
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        // 데이터가 잘 들어가면 review-comments 에 들어간 데이터를 드로잉한다.
        let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
        let newComment = `<div id="reviewcomment-${data.id}-${data.date}" class="p-1">
        <span class="text-badge">${data.date} / <a href="/user?id=${data.author}" class="text-darkmode">${data.authorname}</a></span>
        <span class="edit" data-toggle="modal" data-target="#modal-editreviewcomment" onclick="setEditReviewCommentModal('${data.id}', '${data.date}')">≡</span>
        <span class="remove" data-toggle="modal" data-target="#modal-rmreviewcomment" onclick="setRmReviewCommentModal('${data.id}','${data.date}')">×</span>
        <br>
        <span class="badge badge-${data.itemstatus} me-1">${data.itemstatus}</span>`
        if (data.framecomment) {
            newComment += `<span class="badge badge-secondary m-1 finger" id="reviewcomment-${data.id}-${data.date}-frame" data-toggle="modal" data-target="#modal-gotoframe" onclick="setModalGotoFrame()">${data.frame}f / ${data.frame+data.productionstartframe-1}f</span>`
        }
        newComment += `<small class="text-white">${body}</small>`
        if (data.media != "") {
            if (data.media.includes("http")) {
                newComment += `<div class="row pl-3 pt-3 pb-1">
                    <a href="${data.media}" onclick="copyClipboard('${data.media}')">
                        <img src="/assets/img/link.svg" class="finger">
                    </a>
                    <span class="text-white pl-2 small">${data.mediatitle}</span>
                </div>`
            } else {
                newComment += `<div class="row pl-3 pt-3 pb-1">
                    <a href="${data.protocol}://${data.media}" onclick="copyClipboard('${data.media}')">
                        <img src="/assets/img/link.svg" class="finger">
                    </a>
                    <span class="text-white pl-2 small">${data.mediatitle}</span>
                </div>`
            }
        }
        newComment += `<hr class="my-1 p-0 m-0 divider"></hr></div>`
        document.getElementById("review-comments").innerHTML = newComment + document.getElementById("review-comments").innerHTML;
        // 입력한 값을 초기화 한다.
        document.getElementById("review-comment").value = ""; 
        document.getElementById("review-media").value = "";
        document.getElementById("review-framecomment").checked = false;
    })
    .catch((error) => {
        alert(error)
    });
}

function addAssetTag(id, assettag) {
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            fetch('/api/addassettag', {
                method: 'POST',
                headers: {
                    "Authorization": "Basic "+ document.getElementById("token").value,
                },
                body: new URLSearchParams({
                    id: id,
                    assettag: assettag,
                })
            })
            .then((response) => {
                if (!response.ok) {
                    throw Error(response.statusText + " - " + response.url);
                }
                return response.json()
            })
            .then((data) => {
                // 기존 Tags에 추가된다.
                let url = `/inputmode?project=${data.project}&searchword=assettag:${data.assettag}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&task=`
                source = `<div id="assettag-${data.id}-${data.assettag}"><a href="${url}" class="badge badge-outline-darkmode ml-1">${data.assettag}</a></div>`;
                document.getElementById("assettags-"+data.id).innerHTML = document.getElementById("assettags-"+data.id).innerHTML + source;
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`assettags-${data.id}`).childElementCount > 0) {
                    document.getElementById("assettags-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addassettag" onclick="setAddAssetTagModal('${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmassettag" onclick="setRmAssetTagModal('${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("assettags-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addassettag" onclick="setAddAssetTagModal('${data.id}')">＋</span>
                    `
                }
            })
            .catch((err) => {
                alert(err)
            });
        }
    } else {
        fetch('/api/addassettag', {
            method: 'POST',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
            body: new URLSearchParams({
                id: id,
                assettag: assettag,
            })
        })
        .then((response) => {
            if (!response.ok) {
                throw Error(response.statusText + " - " + response.url);
            }
            return response.json()
        })
        .then((data) => {
            // 기존 Tags에 추가된다.
            let url = `/inputmode?project=${data.project}&searchword=assettag:${data.assettag}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&task=`
            let source = `<div id="assettag-${data.id}-${data.assettag}"><a href="${url}" class="badge badge-outline-darkmode ml-1">${data.assettag}</a></div>`;
            document.getElementById("assettags-"+data.id).innerHTML = document.getElementById("assettags-"+data.id).innerHTML + source;
            // 요소갯수에 따라 버튼을 설정한다.
            if (document.getElementById(`assettags-${data.id}`).childElementCount > 0) {
                document.getElementById("assettags-button-"+data.id).innerHTML = `
                <span class="add ml-1" data-toggle="modal" data-target="#modal-addassettag" onclick="setAddAssetTagModal('${data.id}')">＋</span>
                <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmassettag" onclick="setRmAssetTagModal('${data.id}')">－</span>
                `
            } else {
                document.getElementById("assettags-button-"+data.id).innerHTML = `
                <span class="add ml-1" data-toggle="modal" data-target="#modal-addassettag" onclick="setAddAssetTagModal('${data.id}')">＋</span>
                `
            }
        })
        .catch((err) => {
            alert(err)
        });
    }
}

function rmAssetTag() {
    let project = document.getElementById('modal-rmassettag-project').value
    let id = document.getElementById('modal-rmassettag-id').value
    let assettag = document.getElementById('modal-rmassettag-tag').value
    let isContain = document.getElementById('modal-rmassettag-iscontain').checked
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            fetch('/api/rmassettag', {
                method: 'POST',
                headers: {
                    "Authorization": "Basic "+ document.getElementById("token").value,
                },
                body: new URLSearchParams({
                    project: project,
                    id: id,
                    assettag: assettag,
                    iscontain: isContain,
                })
            })
            .then((response) => {
                if (!response.ok) {
                    throw Error(response.statusText + " - " + response.url);
                }
                return response.json()
            })
            .then((data) => {
                if (isContain) {
                    document.querySelectorAll(`[id^="assettag-${data.id}-"][id*="${data.assettag}"]`).forEach(el => el.remove());
                } else {
                    document.getElementById(`assettag-${data.id}-${data.assettag}`).remove();
                }
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`assettags-${data.id}`).childElementCount > 0) {
                    document.getElementById("assettags-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addassettag" onclick="setRmAssetTagModal('${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmassettag" onclick="setRmAssetTagModal('${data.id}')">－</span>
                    `;
                } else {
                    document.getElementById("assettags-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addassettag" onclick="setAddAssetTagModal('${data.id}')">＋</span>
                    `;
                }
            })
            .catch((err) => {
                alert(err)
            });
        }
    } else {
        fetch('/api/rmassettag', {
            method: 'POST',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
            body: new URLSearchParams({
                project: project,
                id: id,
                assettag: assettag,
                iscontain: isContain,
            })
        })
        .then((response) => {
            if (!response.ok) {
                throw Error(response.statusText + " - " + response.url);
            }
            return response.json()
        })
        .then((data) => {
            if (isContain) {
                document.querySelectorAll(`[id^="assettag-${data.id}-"][id*="${data.assettag}"]`).forEach(el => el.remove());
            } else {
                document.getElementById(`assettag-${data.id}-${data.assettag}`).remove();
            }
            // 요소갯수에 따라 버튼을 설정한다.
            if (document.getElementById(`assettags-${data.id}`).childElementCount > 0) {
                document.getElementById("assettags-button-"+data.id).innerHTML = `
                <span class="add ml-1" data-toggle="modal" data-target="#modal-addassettag" onclick="setAddAssetTagModal('${data.id}')">＋</span>
                <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmassettag" onclick="setRmAssetTagModal('${data.id}')">－</span>
                `;
            } else {
                document.getElementById("assettags-button-"+data.id).innerHTML = `
                <span class="add ml-1" data-toggle="modal" data-target="#modal-addassettag" onclick="setAddAssetTagModal('${data.id}')">＋</span>
                `;
            }
        })
        .catch((err) => {
            alert(err)
        });
    }
}

// 썸네일 엘리먼트들에 드레그앤 드롭 이벤트 처리
const thumbnails = document.getElementsByClassName("thumbnail");

Array.from(thumbnails).forEach((thumbnail) => {
  thumbnail.addEventListener("dragover", (event) => {
    event.preventDefault();
  });

  thumbnail.addEventListener("drop", (event) => {
    event.preventDefault();
    const files = event.dataTransfer.files;
    if (files.length > 1) {
	    alert("must be drop one image");
	    return;
    }

    const file = files[0];
    if (!file) {
	    alert("No valid file");
	    return;
    }

    const project = thumbnail.getAttribute("data-thumbnail-project");    
    const id = thumbnail.getAttribute("data-thumbnail-id");
    if (project != null && id != null) {
        uploadImage(file, project, id);
    }
  });
});


// 이미지 업로드 함수
function uploadImage(file, project, id) {
  const formData = new FormData();
  formData.append("image", file);
  formData.append("project", project);
  formData.append("id", id);

  fetch(`/api/uploadthumbnail`, {
    method: "POST",
    headers: {
        "Authorization": "Basic "+ document.getElementById("token").value,
    },
    body: formData
  })
  .then(async response => {
  const contentType = response.headers.get("content-type") || "";
  const isJSON = contentType.includes("application/json");

  if (!response.ok) {
    const errorText = isJSON ? await response.json() : await response.text();
    throw new Error(isJSON ? errorText.error : errorText);
  }

  return response.json();
  })
  .then(data => {
    const imgURL = "/thumbnail/"+data.project+"/"+data.id+".jpg?" + new Date().getTime();
    const thumbnail = document.getElementById("thumbnail-" + data.id)
    if (thumbnail) {
      thumbnail.src = imgURL;
      thumbnail.onload = () => {}; //화면갱신
    } else {
      console.warn("Thumbnail element not found for id: ", data.id)
    }
  })
  .catch(error => {
    console.error("Error:", error.message);
  });
}
