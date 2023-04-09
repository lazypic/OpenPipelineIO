var webpage = "/scenarios" // webpage
var endpoint = "/api/scenario" // restAPI Endpoint
var uxprefix = "scenario-" // UX prefix string

function UxToObject(obj) {
    obj.project = document.getElementById(uxprefix+'project').value
    obj.order = parseFloat(document.getElementById(uxprefix+'order').value) || 0.0
    obj.seq = parseInt(document.getElementById(uxprefix+'seq').value) || 0
    obj.scene = parseInt(document.getElementById(uxprefix+'scene').value) || 0
    obj.cut = parseInt(document.getElementById(uxprefix+'cut').value) || 0
	obj.name = document.getElementById(uxprefix+'name').value
    obj.ispreviz = document.getElementById(uxprefix+'ispreviz').checked
	obj.istechviz = document.getElementById(uxprefix+'istechviz').checked
	obj.isvisuallab = document.getElementById(uxprefix+'isvisuallab').checked
    obj.seedimage = document.getElementById(uxprefix+'seedimage').value
    obj.prompt = document.getElementById(uxprefix+'prompt').value
    obj.negativeprompt = document.getElementById(uxprefix+'negativeprompt').value
    obj.script = document.getElementById(uxprefix+'script').value
    obj.time = document.getElementById(uxprefix+'time').value
	obj.location = document.getElementById(uxprefix+'location').value
    obj.length = document.getElementById(uxprefix+'length').value
	obj.vfxscript = document.getElementById(uxprefix+'vfxscript').value
    obj.vfxsolution = document.getElementById(uxprefix+'vfxsolution').value
	obj.type = document.getElementById(uxprefix+'type').value
    obj.difficult = document.getElementById(uxprefix+'difficult').value
	obj.ea = parseInt(document.getElementById(uxprefix+'ea').value) || 1
    obj.cost = parseInt(document.getElementById(uxprefix+'cost').value) || 0
    obj.pagenum = parseInt(document.getElementById(uxprefix+'pagenum').value) || 0
    obj.linenum = parseInt(document.getElementById(uxprefix+'linenum').value) || 0
    return obj
}

function ObjectToUx(obj) {
    document.getElementById(uxprefix+'id').value = obj.id
    document.getElementById(uxprefix+'project').value = obj.project
    document.getElementById(uxprefix+'order').value = obj.order
    document.getElementById(uxprefix+'seq').value = obj.seq
    document.getElementById(uxprefix+'scene').value = obj.scene
    document.getElementById(uxprefix+'cut').value = obj.cut
	document.getElementById(uxprefix+'name').value = obj.name
    document.getElementById(uxprefix+'ispreviz').checked = obj.ispreviz
	document.getElementById(uxprefix+'istechviz').checked = obj.istechviz
	document.getElementById(uxprefix+'isvisuallab').checked = obj.isvisuallab
    document.getElementById(uxprefix+'seedimage').value = obj.seedimage
    document.getElementById(uxprefix+'prompt').value = obj.prompt
    document.getElementById(uxprefix+'negativeprompt').value = obj.negativeprompt
    document.getElementById(uxprefix+'script').value = obj.script
    document.getElementById(uxprefix+'time').value = obj.time
	document.getElementById(uxprefix+'location').value = obj.location
    document.getElementById(uxprefix+'length').value = obj.length
	document.getElementById(uxprefix+'vfxscript').value = obj.vfxscript
    document.getElementById(uxprefix+'vfxsolution').value = obj.vfxsolution
    document.getElementById(uxprefix+'type').value = obj.type
	document.getElementById(uxprefix+'difficult').value = obj.difficult
    document.getElementById(uxprefix+'ea').value = obj.ea
    document.getElementById(uxprefix+'cost').value = obj.cost
    document.getElementById(uxprefix+'pagenum').value = obj.pagenum
    document.getElementById(uxprefix+'linenum').value = obj.linenum
}

function AddMode() {
    document.getElementById(uxprefix+'postbutton').hidden = false
    document.getElementById(uxprefix+'deletebutton').hidden = true
    document.getElementById(uxprefix+'putbutton').hidden = true
    InitModal()
}

function EditMode() {
    document.getElementById(uxprefix+'postbutton').hidden = true
    document.getElementById(uxprefix+'deletebutton').hidden = true
    document.getElementById(uxprefix+'putbutton').hidden = false
}

function DeleteMode() {
    document.getElementById(uxprefix+'postbutton').hidden = true
    document.getElementById(uxprefix+'deletebutton').hidden = false
    document.getElementById(uxprefix+'putbutton').hidden = true
}

function string2array(str) {
    let newArr = [];
    if (str === "") {
        return newArr
    }
    let arr = str.split(",");
    for (let i = 0; i < arr.length; i += 1) {
        newArr.push(arr[i].trim())
    }
    return newArr;
}

function InitModal() {
    let inputs = document.querySelectorAll("[id^='"+uxprefix+"']")
    for (let i = 0; i < inputs.length; i += 1) {
        if (inputs[i].type === "checkbox") {
            inputs[i].checked = false
        } else {
            inputs[i].value = ""
        }
    }
}

function SetModal(id) {
    EditMode()
    fetch(endpoint+'/'+id, {
        method: 'GET',
        headers: {"Authorization": "Basic "+ document.getElementById("token").value},
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                tata.error('Error', text, {position:'tr',duration: 5000,onClose: null})
                return
            });
        }
        return response.json()
    })
    .then((obj) => {
        ObjectToUx(obj)
    })
    .catch((err) => {
        console.log(err)
    });
}

function Post() {
    let obj = new Object()
    obj = UxToObject(obj)
    if (obj.script === "") {
        tata.error('Error', "Need script.",{position: 'tr',duration: 5000,onClose: null})
        return
    }
    fetch(endpoint, {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify(obj),
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                tata.error('Error',text,{position:'tr',duration: 5000,onClose: null})
                return
            });
        }
        return response.json()
    })
    .then((obj) => {
        tata.success('Add', obj.script + "가 추가되었습니다.", {position: 'tr',duration: 5000,onClick: tataLink,onClose: null})
    })
    .catch((err) => {
        console.log(err)
    });
}

function Put() {
    let obj = new Object()
    obj = UxToObject(obj)
    if (obj.script === "") {
        tata.error('Error',"Need script.",{position: 'tr',duration: 5000,onClose: null})
        return
    }
    fetch(endpoint+'/'+document.getElementById(uxprefix+'id').value, {
        method: 'PUT',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify(obj),
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                tata.error('Error', text, {position: 'tr',duration: 5000,onClose: null})
                return
            });
        }
        return response.json()
    })
    .then((obj) => {
        tata.success('Edit', obj.script + "가 편집되었습니다.", {position: 'tr',duration: 5000,onClick: tataLink,onClose: null})
    })
    .catch((err) => {
        console.log(err)
    });
}

function Delete() {
    fetch(endpoint+'/'+document.getElementById(uxprefix+'id').value, {
        method: 'DELETE',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                tata.error('Error', text, {position: 'tr',duration: 5000,onClose: null})
                return
            });
        }
        return response.json()
    })
    .then((obj) => {
        tata.success('Delete', obj.id + "가 삭제되었습니다.", {position: 'tr',duration: 5000,onClick: tataLink,onClose: null})
    })
    .catch((err) => {
        console.log(err)
    });
}

function tataLink() {
    window.location.replace(webpage)
}

function string2array(str) {
    let newArr = [];
    if (str === "") {
        return newArr
    }
    let arr = str.split(",");
    for (let i = 0; i < arr.length; i += 1) {
        newArr.push(arr[i].trim())
    }
    return newArr;
}

// 이미지 생성 버튼을 클릭하면 이미지를 생성하고 이미지 소스를 업데이트합니다.
document.querySelectorAll('.btn-genimage').forEach(button => {
    button.addEventListener('click', async () => {
      try {
        //const imageUrl = await generateImage();
        //const generatedImage = document.getElementById('generatedImage');
        //generatedImage.src = imageUrl;
        //generatedImage.style.display = 'block';
        console.log("id")
      } catch (error) {
        console.error(error);
      }
    });
});