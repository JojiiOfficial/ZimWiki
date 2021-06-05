document.addEventListener("DOMContentLoaded", function() {
	document.querySelector(".no-js").style.visibility="visible";
});

function enterSearch(ele){
	if (event.key === 'Enter') {
		search(ele.value);
	}
}

function search(text) {
	let searchForm = document.getElementById("searchform");
	let queryData = document.getElementById("sQuery");
	document.querySelector("#searchform").setAttribute("action","/search/"+wiki+"/");
	queryData.value = text;
	if (text.replace(/\s/g, "") != "") {
		searchForm.submit();
	}
}

// For all those vim lovers
document.addEventListener('keyup', (e) => {
    if (e.keyCode === 191) {
      document.getElementById("searchtxt").focus();
    }
});