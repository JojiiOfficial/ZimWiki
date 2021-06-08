document.addEventListener("DOMContentLoaded", function() {
	document.querySelector(".no-js").style.visibility="visible";
});

// For all those vim lovers
document.addEventListener('keyup', (e) => {
    if (e.keyCode === 191) {
      document.getElementById("sQuery").focus();
    }
});

if (wiki == "-") {
	document.getElementById("searchtxt").placeholder = "Search anywhere";
} else {
	document.getElementById("searchtxt").placeholder = "Search in this wiki";
}