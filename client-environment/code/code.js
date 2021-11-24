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
const IP_SERVER = "http://192.168.1.196/"

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
            idPatient: documentId,
            patientNames: names,
            patientLastnames: lastnames,
            dateClinicHistory: dateHc,
            typeId: typeId,
            hasErrors: patientErrors
        }
        // converting to json
        var sendRecord = JSON.stringify(recordDataPatient)
        
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