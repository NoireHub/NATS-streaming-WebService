const form = document.getElementById("form-container");
const pageContainer = document.querySelector(".page-container");

const infoContainer = document.createElement("div");
infoContainer.id = "uid-info-container";

const infoSpan = document.createElement("pre");

let isFirstRender = true;

console.log(pageContainer);

function handleUIDValue(form) {
  return form.elements["0"].value;
}

function getUIDInfo(uid) {
  return fetch(`http://localhost:5555/order/${uid}`)
    .then((response) => {
      if (response.status == 200) {
        return response.text();
      } else if (response.status == 204) {
        return "UID doesn't exist";
      }
    })
    .catch((err) => err);
}

function appendInfoWindow(info) {
  infoSpan.innerHTML = info;
  infoContainer.append(infoSpan)
  pageContainer.append(infoContainer);
}

function changeInfoWindowContent(info) {
  infoSpan.innerHTML = info;
}

form.addEventListener("submit", (e) => {
  e.preventDefault();

  let uidInput = handleUIDValue(form);
  getUIDInfo(uidInput).then((data) => {
    data = data.replace(/,/g,"\n").replace(/{/g,"\n{ ")
    if (isFirstRender) {
      appendInfoWindow(data);
      isFirstRender = false;
    } else {
      changeInfoWindowContent(data);
    }
    
  });
});

// localhost:5555/order/
