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
console.log(buttonSubmit)
const IP_SERVER = "http://172.16.0.59:8005/"

IdPatientBox.onchange = (e) => {
    
}
function boxesForm() {
    

}
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
            // setting the loading progress
            console.log("enviando...")
            // sending the data
            var stateRecord = await new Promise((recorded, rejected) => {
                fetch(IP_SERVER + "record-patient", 
                {
                    method: "post",
                    headers: {
                        "content-type":"application/json"
                    },
                    body: sendRecord
                })
                .then(data => data.json())
                .then(data => recorded(data))
                .catch(error => rejected(error))
            })
            console.log(stateRecord)
        }
        sendRecordToServer()
    }else {
        console.log("faltan campos")
    }
}

buttonCloseMenu.onclick = () => {
    menu.classList.remove("show-win-menu")
    menu.classList.add("hide-win-menu")
}
buttonOpenMenu.onclick = () => {
    menu.classList.add("hide-win-menu")
    menu.classList.add("show-win-menu")
}