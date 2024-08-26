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

var mylatesttap;
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
document
  .getElementsByClassName(".toggle")
  .addEventListener("touchstart", toggle_edit_mode, false);
