document.addEventListener("DOMContentLoaded", function() {
	document.querySelector(".no-js").style.visibility = "visible";
});

if (typeof queryText !== 'undefined') {
  document.getElementById("sQuery").value = queryText;
}

document.getElementById("sQuery").placeholder = searchPlaceholder;

// Display a spinner if the server takes too long to respond
window.onbeforeunload = function() {
	setTimeout(function() {
		document.getElementById("main").style.transition = "opacity 0.4s";
		document.getElementById("main").style.opacity = "0.35";
		document.getElementById("main").style.pointerEvents = "none";
		document.getElementsByClassName("spinner-container")[0].style.transition = "opacity 0.4s";
		document.getElementsByClassName("spinner-container")[0].style.visibility = "visible";
		document.getElementsByClassName("spinner-container")[0].style.opacity = "1";
	}, 4000);
}

// Remove the spinner when the page is displayed (e.g.: return to the previous page)
window.onpageshow = function() {
	document.getElementById("main").style.transition = "opacity 0.4s";
	document.getElementById("main").style.opacity = "1";
	document.getElementById("main").style.pointerEvents = "all";
	document.getElementsByClassName("spinner-container")[0].style.transition = "opacity 0.4s";
	document.getElementsByClassName("spinner-container")[0].style.opacity = "0";
	document.getElementsByClassName("spinner-container")[0].style.visibility = "hidden";
}

// For all those vim lovers
document.addEventListener('keyup', (e) => {
	if (e.keyCode === 191) {
		document.getElementById("sQuery").focus();
	}
});