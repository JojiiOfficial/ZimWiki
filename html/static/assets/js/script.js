document.addEventListener("DOMContentLoaded", function() {
	document.querySelector(".no-js").style.visibility="visible";
});

if (typeof queryText !== 'undefined') {
	document.getElementById("sQuery").value = queryText;
}

if (wiki == "-") {
	document.getElementById("sQuery").placeholder = "Search anywhere";
} else {
	document.getElementById("sQuery").placeholder = "Search in this wiki";
}

// For all those vim lovers
document.addEventListener('keyup', (e) => {
    if (e.keyCode === 191) {
      document.getElementById("sQuery").focus();
    }
});