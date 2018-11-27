var entries, catSelect = "", langSelect = "";
var tagSelect = new Array;


var requestURL = 'output.min.json';
var request = new XMLHttpRequest();
request.open('GET', requestURL);
request.responseType = 'json';
request.send();

request.onload = function() {
    entries = request.response.Entries;
    var tags = request.response.Tags;
    var cats = request.response.Cats;
    var langs = request.response.Langs;
    //populateCats(cats);
    populateTags(tags);
    populateLangs(langs);
    populateEntries(entries, "", [], "");
};

function populateLangs(lObj) {
    var z, txt3 = "";
    for (z in lObj) {
        txt3 += "<span id='" + lObj[z].Lang + "' class='tag is-success' onClick='langPicker(`" + lObj[z].Lang + "`)'>" + lObj[z].Lang + "</span>"
    }
    document.getElementById("lang").innerHTML = txt3;
}

/*function populateCats(cObj) {
    var y, txt2 = "";
    for (y in cObj) {
        txt2 += "<li><a id='" + cObj[y].Cat + "' onClick='catPicker(`" + cObj[y].Cat + "`)' href='#'>" + cObj[y].Cat + "</a></li>"
    }
    document.getElementById("cat").innerHTML = txt2;
}*/

function populateTags(tObj) {
    var z, txt = "";
    for (z in tObj) {
        txt += "<span id='" + tObj[z].Tag + "' class='tag is-link' onClick='tagPicker(`" + tObj[z].Tag + "`)'>" + tObj[z].Tag + "</span>"
    }
    document.getElementById("tagList").innerHTML = txt;
}

function populateEntries(eObj, catSelect, tagSelect, langSelect) {
    var x, txt = "";

    for (x in eObj) {
        if ((!tagSelect.some(ele => !eObj[x].T.includes(ele) || tagSelect === []) && (eObj[x].La.includes(langSelect) || langSelect === ""))) {
            txt += "<div class='card' style='margin-bottom:24px'><header class='columns card-header is-marginless has-background-light'>";
            txt += "<span class='column is-narrow'>" + addNonFree(eObj[x].NF) + addPdep(eObj[x].P) + "</span>";
            txt += "<span class='column'><h4 class='title is-4'>" + eObj[x].N + "</h4></span>";
            txt += "<span class='column is-narrow tags'>" + getL(eObj[x].Li, "is-primary") + getL(eObj[x].La, "is-success") + "</span></header>";
            txt += "<span class='card-footer'>" + getTags(eObj[x].T);
            txt += "<span class='column'>" + eObj[x].D + "<span class='level'>" + getLinks(eObj[x].Sr, "Source Code") + getLinks(eObj[x].Si, "Website") + getLinks(eObj[x].Dem, "Demo") + "</span></span></span></div>";
        }
    }
    document.getElementById("demo").innerHTML = txt;
}
function goHome() {
    catSelect = "";
    tagSelect = [];
    langSelect = "";
    rmvActive();
    document.getElementById("home").classList.add("is-active");
    populateEntries(entries, catSelect, tagSelect, langSelect)
}
function rmvActive() {
    let els = document.getElementsByClassName('is-active');
    while (els[0]) {
        els[0].classList.remove('is-active')
    }
}
/*function catPicker(c) {
    rmvActive();
    catSelect = c;
    document.getElementById(c).classList.add("is-active");
    populateEntries(entries, catSelect, "", "")
}*/
function tagPicker(t) {
    tagSelect.push(t)
    document.getElementById(t).classList.add("is-active");
    populateEntries(entries, "", tagSelect, langSelect)
}
function langPicker(l) {
    rmvActive();
    langSelect = l;
    document.getElementById(l).classList.add("is-active");
    populateEntries(entries, "", tagSelect, langSelect)
}
function getTags(t) {
    var tags = "<span class='column is-one-third tags is-marginless'>";
    t.forEach(function(item) {
        tags += "<span class='tag is-link' onclick='tagPicker(`" + item + "`)'>" + item + "</span>";
    });
    tags += "</span>";
    return tags
}
function getLinks(l, t) {
    if (l !== undefined){
        switch (t) {
            case "Source Code":
                return "<a href='" + l + "'class='level-item'><span class='icon has-text-link'><i class='fas fa-lg fa-code-branch'></i></span>" + t + "</a>";
            case "Website":
                return "<a href='" + l + "'class='level-item'><span class='icon has-text-link'><i class='fas fa-lg fa-external-link-alt'></i></span>" + t + "</a>";
            case "Demo":
                return "<a href='" + l + "'class='level-item'><span class='icon has-text-link'><i class='fas fa-lg fa-chevron-circle-right'></i></span>" + t + "</a>";
        }
    } else {
        return ""
    }
}
function getL(t, cl) {
    var tags = "";
    t.forEach(function(item) {
        tags += "<span class='tag is-link " + cl + "' onclick='langPicker(`" + item + "`)'>" + item + "</span>";
    });
    tags += "";
    return tags
}

function addNonFree(nf) {
    if (nf === true) {
        return "<span class='icon is-medium has-text-danger'><i class='fas fa-2x fa-ban'></i></span>"
    } else {
        return ""
    }
}

function addPdep(p) {
    if (p === true) {
        return "<span class='icon is-medium has-text-warning'><i class='fas fa-2x fa-exclamation-triangle'></i></span>"
    } else {
        return ""
    }
}