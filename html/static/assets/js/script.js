function enterSearch(ele){
	if (event.key === 'Enter'){
		search(ele.value)
	}
}

function search(text) {
	searchForm = document.getElementById("searchform")
	queryData = document.getElementById("sQuery")

	document.querySelector("#searchform").setAttribute("action","/search/"+wiki+"/");

	queryData.value = text

	searchForm.submit()
}

// For all those vim lovers
document.addEventListener('keyup', (e) => {
    if (e.keyCode === 191){
      document.getElementById("searchtxt").focus()
    }
});