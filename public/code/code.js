var menu = document.getElementById("menu-panel")
var buttonCloseMenu = document.getElementById("close-menu")
var buttonOpenMenu = document.getElementById("open-menu")
var boxNames = document.getElementById("patient-names")
var boxLastnames = document.getElementById("patient-lastnames")
var dateRegistry = document.getElementById("date-reg")
var checkboxError = document.getElementById("error-check")
var IdPatientBox = document.getElementById("id-patient")
var typeIdPatient = document.getElementById("type-id")
var buttonSubmit = document.getElementById("button-submit")
var app = document.getElementById("app-container")


//const IP_SERVER = "http://192.168.11.105:8005/"
IdPatientBox.onchange = async () => {
    var id = parseInt(IdPatientBox.value)   
    if (!isNaN(id)) {
        var stateMessage = " | consultando..."
        IdPatientBox.disabled = true
        IdPatientBox.style.backgroundColor = "rgba(1, 172, 240)"
        IdPatientBox.style.color = "white"
        IdPatientBox.value = IdPatientBox.value + stateMessage
        var getInfoPatient = await new Promise((resolved, rejected)=> {
            fetch("/get-data-patient?id-patient=" + id, {
                method:"get"
            })
            .then(resp => resp.json())
            .then(data => resolved(data))
            .then(error => rejected(error))
        })
        stateMessage = ""
        IdPatientBox.disabled = false
        IdPatientBox.style.backgroundColor = "white"
        IdPatientBox.style.color = "black"
        IdPatientBox.value = id
        if (getInfoPatient.FirstName.trim() != "" && getInfoPatient.FirstLastname.trim() != "") {
            boxNames.value = convertNameTo(getInfoPatient.FirstName.trim()) + " " + convertNameTo(getInfoPatient.SecondName.trim())
            boxLastnames.value = convertNameTo(getInfoPatient.FirstLastname.trim()) + " " + convertNameTo(getInfoPatient.SecondLastname.trim())
        }
        switch(getInfoPatient.TypId.trim()) {
            case "CC":
                typeIdPatient.value = 0
                break;
            case "TI":
                typeIdPatient.value = 1
                break;
            case "CE":
                typeIdPatient.value = 2
                break;
            case "ASI":
                typeIdPatient.value = 3
                break;
            case "CI":
                typeIdPatient.value = 4
                break;
            case "MSI":
                typeIdPatient.value = 5
                break;
            case "NU":
                typeIdPatient.value = 6
                break;
            case "PA":
                typeIdPatient.value = 7
                break;
            case "PE":
                typeIdPatient.value = 8
                break;
            case "RC":
                typeIdPatient.value = 9
                break;
            case "RI":
                typeIdPatient.value = 10
                break;
            default:
                typeIdPatient.value = 0
                break;
        }
    }
}
function boxesForm() {
    

}
buttonSender()
function getDiv() {
    var div = document.createElement("div")
    var i = document.createElement("i")
    i.classList.add("fas", "fa-times-circle")
    div.appendChild(i)
    return div
}
function getAlertDiv() {
    return document.createElement("div")
}
function buttonSender() {
    buttonSubmit.onclick = () => {
        // getting values from boxes in form
        var documentId = IdPatientBox.value 
        var typeId = typeIdPatient.value
        var names = boxNames.value
        var lastnames = boxLastnames.value
        var dateHc = dateRegistry.value
        var patientErrors = checkboxError.value
    
        switch(typeId) {
            case "0":
                typeId = "CC"
                break;
            case "1":
                typeId = "TI"
                break;
            case "2":
                typeId = "CE"
                break;
            case "3":
                typeId = "ASI"
                break;
            case "4":
                typeId = "CI"
                break;
            case "5":
                typeId = "MSI"
                break;
            case "6":
                typeId = "NU"
                break;
            case "8":
                typeId = "PA"
                break;
            case "9":
                typeId = "PE"
                break;
            case "10":
                typeId = "RC"
                break;
            case "11":
                typeId = "RI"
                break;
            default:
                typeId = ""
                break;
        }
        if(documentId != "" && documentId.length > 5 && names != "" && lastnames != "" && dateHc != "" && typeId != "") {
            // checking if the patient has errors
            if(checkboxError.checked) {
                patientErrors = true
            }else {
                patientErrors = false
            }
            // data in object
            var recordDataPatient = {
                idPatient: parseInt(documentId),
                patientNames: names,
                patientLastnames: lastnames,
                dateClinicHistory: dateHc,
                typeId: typeId,
                hasError: patientErrors
            }
            // converting to json
            var sendRecord = JSON.stringify(recordDataPatient)
    
            async function sendRecordToServer() {
                IdPatientBox.value  = ""
                boxNames.value = ""
                boxLastnames.value = ""
                checkboxError.value = ""
                // setting the loading progress
                onprogressRequest()
                // sending the data
                var stateRecord = await new Promise((recorded, rejected) => {
                    fetch("/record-patient", 
                    {
                        method: "post",
                        headers: {
                            "content-type":"application/json"
                        },
                        body: sendRecord
                    })
                    .then(resp => resp.json())
                    .then(data => recorded(data))
                    .catch(error => rejected(error))
                })
                stateRecord = stateRecord.ContenMessage
                if (stateRecord.includes("Error 1062") || stateRecord.includes("Duplicate entry")) {
                    console.log("El usuario ya existe. Se ha denegado el registro")
                }else {
                    if(stateRecord == "successfull") {
                        console.log("Exitoso!")
                    }
                }
                buttonSender()
                removeLastElement()
            }
            sendRecordToServer()
        }else {
            console.log("faltan campos")
        }
    }
}
function stateProcessAlert(iconClass, message, backgroundColor) {
    var successfull = getDiv()
    var messageP =  document.createElement("p")
    var icon = document.createElement("i")
    icon.classList.add("fas", "fa-server")
    messageP.innerHTML = "Error interno de servidor (Código:1000)"
    successfull.appendChild(icon)
    successfull.appendChild(messageP)
    successfull.classList.add("alert-state")
    document.body.appendChild(successfull)
}
function setFailProcessAlert()  {
    var fail = getDiv()
    fail.innerHTML += "<i class='fas fa-exclamation-triangle'></i>El paciente ingresado ya existe. Se ha denegado el registro"
    fail.style.backgroundColor = "red"
    fail.classList.remove("from-top-bottom-alert")
    fail.classList.add("from-bottom-top-alert")
    document.body.appendChild(fail)
}
buttonCloseMenu.onclick = () => {
    menu.classList.remove("show-win-menu")
    menu.classList.add("hide-win-menu")
}
buttonOpenMenu.onclick = () => {
    menu.classList.add("hide-win-menu")
    menu.classList.add("show-win-menu")
}
function onprogressRequest() {
    iconHeader("none")
    var onprogress = getAlertDiv()   
    var blockMessage = getAlertDiv()
    var iconLoad = getAlertDiv()

    onprogress.style.width = "100%"
    onprogress.style.height = "100%"
    onprogress.style.backgroundColor = "rgba(0,0,0, .5)"
    onprogress.style.position = "absolute"
    onprogress.style.zIndex = "100"
    onprogress.style.display = "flex"
    onprogress.style.alignItems = "center"
    onprogress.style.justifyContent = "center"

    blockMessage.style.width = "100px"
    blockMessage.style.height ="60px"
    blockMessage.style.backgroundColor = "#00adef"
    blockMessage.style.marginTop = "-100px"
    blockMessage.style.borderRadius = "10px"
    blockMessage.style.display = "flex"
    blockMessage.style.alignItems = "center"
    blockMessage.style.justifyContent = "center"

    iconLoad.style.width = "40px"
    iconLoad.style.height = "40px"
    iconLoad.style.backgroundImage = "url('../public/Images/loading.gif')"
    iconLoad.style.backgroundSize = "contain"
    blockMessage.appendChild(iconLoad)
    onprogress.appendChild(blockMessage)
    app.appendChild(onprogress)
}
function removeLastElement() {
   app.removeChild(app.lastElementChild)
}
function iconHeader(state) {
    var icon = app.getElementsByTagName('header')
    icon[1].children[0].style.display = state
}
function convertNameTo(string) {
    if(string != "") {
        var str1 = string[0]
        var str2 = str1 + string.substring(1, string.length).toLowerCase()
        string = str2
        return string
    }else {
        return ""
    }
}
stateProcessAlert()