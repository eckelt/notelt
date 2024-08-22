// on ctrl-e toggle /e/ in location href
// define a handler
function doc_keyUp(e) {
  // this would test for whichever key is 40 (down arrow) and the ctrl key at the same time
  if (e.ctrlKey && e.key === "e") {
    // call your function to do the thing
    const pathArr = location.href.split("/");
    if (pathArr.indexOf("e") > -1) {
      location.href = pathArr.filter((item) => item !== "e").join("/");
    } else {
      var note = pathArr.pop();
      location.href = pathArr.join("/") + "/e/" + note;
    }
  }
}
// register the handler
document.addEventListener("keyup", doc_keyUp, false);
