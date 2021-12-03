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
var generateReportButton = document.getElementById("report-info")

function getOnlyDiv() {
    var div = document.createElement("div")
    var header = document.createElement("header")
    div.style.width = "100%"
    div.style.height = "100%"
    div.style.backgroundColor = "white"
    div.style.position = "absolute"
    div.style.left = "0"
    div.style.top = "0"
    header.style.width = "100%"
    div.appendChild(header)
    app.appendChild(div)
    buttonCloseMenu.click()
    return div
}
function getButtons() {
    var button = document.createElement("button")
    button.style.width = "200px"
    button.style.height = "30px"
    button.style.backgroundColor = "rgb(30, 219, 5)"
    return button
}
function getInputDate() {
    var date = document.createElement("input")
    date.type = "date"
    date.style.textAlign = "center"
    return date
}
function getSpan() {
    var span = document.createElement("span")
    return span
}
function getSelect() {
    return document.createElement("select")
}
function getOptions() {
    return document.createElement("option")
}
generateReportButton.onclick = () => {
    var div = getOnlyDiv()
    div.classList.add("box-generate-info")
    var date1 = getInputDate()
    var date2 = getInputDate()
    var select = getSelect()
    var item1 = getOptions()
    var item2 = getOptions()

    var checkPatientErrors = createCheckboxes("patient-errors", "Generar solo pacientes con errores")

    item1.innerHTML = "Fecha de historia"
    item2.innerHTML = "Fecha de registro"
    item1.value = 1
    item2.value = 0
    select.appendChild(item1)
    select.appendChild(item2)
    select.classList.add("select-of-report")
    var button = getButtons()
    div.children[0].innerHTML = "<i class='fas fa-arrow-left' id='return-arrow'></i>" + "Generación de reporte"
    div.appendChild(getSpan()).innerHTML = "Fecha inicio"
    div.appendChild(date1)
    div.appendChild(getSpan()).innerHTML = "Fecha fin"
    div.appendChild(date2)
    div.appendChild(getSpan()).innerHTML = "Generar por"
    div.appendChild(select)
    div.appendChild(checkPatientErrors)
    div.appendChild(button).innerHTML = "Generar reporte"
    document.getElementById("return-arrow").onclick = () => deleteActualWin(app, div)
    app.style.height = "380px"

    button.onclick = async () => {
        var dateStart = date1.value
        var dateEnd = date2.value
        if (dateStart != "" && dateEnd != "") {
            var checkPatientError = checkPatientErrors.children[0].checked
            if (checkPatientError) {
                checkPatientError = 1
            } else {
                checkPatientError = 0
            }
            var query = "?date-start=" + dateStart + "&date-end=" + dateEnd + "&check-only-p-errors=" + checkPatientError + "&gen-by=" + select.value
            onprogressRequest()
            var state = await new Promise((resolved, rejected) => {
                fetch("/get-information-from-patient" + query, {
                    method: "get"
                })
                    .then(data => data.json())
                    .then(data => resolved(data))
                    .catch(error => rejected(error))
            })
            tableInfo(state)
            removeLastElement()
        } else {
            stateProcessAlert("fa-info-circle", "Faltan campos por llenar, por favor verifique", "orange")
        }
    }

    function createCheckboxes(idCheckbox, contentLabel) {
        var spanCheckbox = document.createElement("span")
        var checkbox = document.createElement("input")
        var label = document.createElement("label")

        spanCheckbox.style.width = "100%"
        checkbox.type = "checkbox"
        checkbox.id = idCheckbox
        label.innerHTML = contentLabel
        label.htmlFor = idCheckbox
        label.style.cursor = "pointer"

        spanCheckbox.appendChild(checkbox)
        spanCheckbox.appendChild(label)
        return spanCheckbox
    }
}
function deleteActualWin(app, div) {
    app.removeChild(div)
    app.style.height = "450px"
}
//const IP_SERVER = "http://192.168.11.105:8005/"
IdPatientBox.onchange = async () => {
    var id = parseInt(IdPatientBox.value)
    if (!isNaN(id)) {
        var stateMessage = " | consultando..."
        IdPatientBox.disabled = true
        IdPatientBox.style.backgroundColor = "rgba(1, 172, 240)"
        IdPatientBox.style.color = "white"
        IdPatientBox.value = IdPatientBox.value + stateMessage
        var getInfoPatient = await new Promise((resolved, rejected) => {
            fetch("/get-data-patient?id-patient=" + id, {
                method: "get"
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
        if (getInfoPatient.Names != "" && getInfoPatient.Lastnames != "") {
            boxNames.value = getInfoPatient.Names
            boxLastnames.value = getInfoPatient.Lastnames
            typeIdPatient.value = getInfoPatient.TypId
        } else {
            boxNames.value = ""
            boxLastnames.value = ""
            stateProcessAlert("fa-address-book", "Sin registros en nuestro sistema interno", "rgb(243, 98, 1)")
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

        switch (typeId) {
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
        if (documentId != "" && documentId.length > 5 && names != "" && lastnames != "" && dateHc != "" && typeId != "") {
            // checking if the patient has errors
            if (checkboxError.checked) {
                patientErrors = true
            } else {
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
                // setting the loading progress
                onprogressRequest()
                // sending the data
                var stateRecord = await new Promise((recorded, rejected) => {
                    fetch("/record-patient",
                        {
                            method: "post",
                            headers: {
                                "content-type": "application/json"
                            },
                            body: sendRecord
                        })
                        .then(resp => resp.json())
                        .then(data => recorded(data))
                        .catch(error => rejected(error))
                })
                stateRecord = stateRecord.ContenMessage
                if (stateRecord.includes("Error 1062") || stateRecord.includes("Duplicate entry")) {
                    stateProcessAlert("fa-user-times", "Usuario existente, se ha denegado el registro", "red")
                } else if (stateRecord.includes("dial tcp: i/o timeout")) {
                    stateProcessAlert("fa-user-times", "Lo sentimos, inténtelo nuevamente (dial/tcp)", "red")
                } else {
                    if (stateRecord == "successfull") {
                        stateProcessAlert("fa-user-check", "Registro éxitoso", "limegreen")
                    }
                }
                IdPatientBox.value = ""
                boxNames.value = ""
                boxLastnames.value = ""
                checkboxError.value = ""
                buttonSender()
                removeLastElement()
            }
            sendRecordToServer()
        } else {
            stateProcessAlert("fa-info-circle", "Faltan campos por llenar, por favor verifique", "orange")
        }
    }
}
function stateProcessAlert(iconClass, message, backgroundColor) {
    var successfull = getDiv()
    var messageP = document.createElement("p")
    var icon = document.createElement("i")
    icon.classList.add("fas", iconClass)
    messageP.innerHTML = message
    successfull.style.backgroundColor = backgroundColor
    successfull.appendChild(icon)
    successfull.appendChild(messageP)
    successfull.classList.add("alert-state", "transition-show")
    document.body.appendChild(successfull)
    setTimeout(() => {
        successfull.classList.remove("transition-show")
        successfull.classList.add("transition-hide")
        setTimeout(() => {
            successfull.classList.remove("transition-hide")
            successfull.style.transform = "translateX(110%)"
            successfull.style.transition = "1s"
            setTimeout(() => {
                document.body.removeChild(successfull)
            }, 1000);
        }, 2000);
    }, 100);
}
function setFailProcessAlert() {
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
    blockMessage.style.height = "60px"
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
    iconHeader("block")
}
function iconHeader(state) {
    var icon = app.getElementsByTagName('header')
    icon[1].children[0].style.display = state
}
function convertNameTo(string) {
    if (string != "") {
        var str1 = string[0]
        var str2 = str1 + string.substring(1, string.length).toLowerCase()
        string = str2
        return string
    } else {
        return ""
    }
}
function tableInfo(contentQuery) {
    if (contentQuery != null) {
        var buttonDownload = document.getElementById("download-report")
        var table = document.getElementById("table-info-view-patient")
        var closeInfo = document.getElementById("close-modal-info")
        for (var i = 0; i < contentQuery.length; i++) {
            var tr = document.createElement("tr")
            var data = JSON.parse(contentQuery[i])
            for (var j = 0; j < 7; j++) {
                var td = document.createElement("td")
                switch (j) {
                    case 0:
                        td.innerHTML = data.IdPatient
                        break;
                    case 1:
                        td.innerHTML = data.TypeId
                        break;
                    case 2:
                        td.innerHTML = data.DateClinicHistory
                        break;
                    case 3:
                        td.innerHTML = data.ActualDateRegistry
                        break;
                    case 4:
                        td.innerHTML = data.PatientNames
                        break;
                    case 5:
                        td.innerHTML = data.PatientLastnames
                        break;
                    case 6:
                        td.innerHTML = data.HasError
                        break;
                }
                tr.appendChild(td)
            }
            table.appendChild(tr)
        }
        showInfoModal()
        closeInfo.onclick = () => {
            closeInfoModal()
        }
        buttonDownload.onclick = async () => {
            var dataExcel = {
                DataExcel: contentQuery
            }
            var downlaodReport = await new Promise((resolved, rejected)=>{
                fetch("/get-report-in-excel", {
                    method:"post",
                    body: JSON.stringify(dataExcel),
                })
                .then(data => data.json())
                .then(data => resolved(data))
                .catch(error => rejected(error))
            })
            location.href = downlaodReport.Link
        }

    }else {
        stateProcessAlert("fa-address-book", "No se han encontrado registros", "rgb(243, 98, 1)")   
    }
}
function closeInfoModal() {
    document.getElementById("information-view-patient").style.display = "none"
    var tbody = document.getElementById("table-info-view-patient")
    while(tbody.children.length > 1) {
        tbody.removeChild(tbody.lastElementChild)
    }
}
function showInfoModal() {
    document.getElementById("information-view-patient").style.display = "block"
}

