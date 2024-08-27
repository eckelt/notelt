// on ctrl-e toggle /e/ in location href
// define a handler
function doc_keyUp(e) {
  // this would test for whichever key is 40 (down arrow) and the ctrl key at the same time
  if (e.ctrlKey && e.key === "e") {
    // call your function to do the thing
    toggle_edit_mode();
  }
}

toggle_edit_mode = function () {
  const pathArr = location.href.split("/");
  if (pathArr.indexOf("e") > -1) {
    location.href = pathArr.filter((item) => item !== "e").join("/");
  } else {
    var note = pathArr.pop();
    location.href = pathArr.join("/") + "/e/" + note;
  }
};

var mylatesttap = 0;
function doubletap() {
  var now = new Date().getTime();
  var timesince = now - mylatesttap;
  if (timesince < 600 && timesince > 0) {
    // double tap
    toggle_edit_mode();
  }
  mylatesttap = new Date().getTime();
}

// register the handler
document.addEventListener("keyup", doc_keyUp, false);

// edit notes
const noteName = location.href.split("/").pop();

function save() {
  const noteValue = document.querySelector(".note").value;
  if (localStorage.getItem("note") !== noteValue) {
    localStorage.setItem("note", noteValue);
    console.log("Saved note ", noteName);
    fetch("/api/v1/note/" + noteName, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: noteValue,
    });
  }

  setTimeout(save, 500);
}

window.onload = async function () {
  const toggleEl = document.getElementsByClassName("toggle")[0];
  toggleEl.addEventListener("touchstart", doubletap);

  let note = document.querySelector(".note");
  if (note !== undefined && note !== null) {
    const response = await fetch("/api/v1/note/" + noteName);
    if (!response.ok) {
      console.log(`Response status: ${response.status}`);
    }
    note.value = await response.text();
    save();
  }
};
