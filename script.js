var entries, catSelect = "", tagSelect = "";

var requestURL = 'output.min.json';
var request = new XMLHttpRequest();
request.open('GET', requestURL);
request.responseType = 'json';
request.send();

request.onload = function() {
    entries = request.response.Entries;
    var tags = request.response.Tags;
    var cats = request.response.Cats;
    populateCats(cats);
    populateTags(tags);
    populateEntries(entries, "", "");
}

function populateCats(cObj) {
    var y, txt2 = "";
    for (y in cObj) {
        txt2 += "<li><a id='" + cObj[y].Cat + "' onClick='catPicker(`" + cObj[y].Cat + "`)' href='#'>" + cObj[y].Cat + "</a></li>"
    }
    document.getElementById("cat").innerHTML = txt2;
}

function populateTags(tObj) {
    var z, txt = "";
    for (z in tObj) {
        txt += "<div class='tags has-addons'><span id='" + tObj[z].Tag + "' class='tag is-link is-rounded' onClick='tagPicker(`" + tObj[z].Tag + "`)'>" + tObj[z].Tag + "</span><span class='tag is-info is-rounded'>" + tObj[z].C + "</span></div>"
    }
    document.getElementById("tagList").innerHTML = txt;
}
function populateEntries(eObj, catSelect, tagSelect) {
    var x, txt = "";
    console.log(catSelect, " ", tagSelect)
    for (x in eObj) {
        if (eObj[x].C == catSelect || eObj[x].T.includes(tagSelect) || (tagSelect == "" && catSelect == "")) {
            txt += "<div class='card' style='margin-bottom:24px'><header class='columns card-header is-marginless has-background-light'>"

            txt += "<span class='column is-narrow'>" + addNonFree(eObj[x].F) + addPdep(eObj[x].P) + "</span>"
            txt += "<span class='column'><h4 class='title is-4'>" + eObj[x].N + "</h4></span>"
            txt += "<span class='column is-narrow tags'>" + getL(eObj[x].Li, "is-primary") + getL(eObj[x].La, "is-success") + "</span></header>"


            txt += "<span class='card-footer'>" + getTags(eObj[x].T)
            txt += "<span class='column'>" + eObj[x].D + "<span class='level'>" + getLinks(eObj[x].Sr, "Source Code") + getLinks(eObj[x].Si, "Website") + getLinks(eObj[x].Dem, "Demo") + "</span></span></span></div>"
            //txt += "<footer class='card-footer'>" + getLinks(eObj[x].Sr, "Source Code") + getLinks(eObj[x].Si, "Website") + getLinks(eObj[x].Dem, "Demo") + "</footer>";
        }
    }
    document.getElementById("demo").innerHTML = txt;
}
function goHome() {
    catSelect = ""
    tagSelect = ""
    rmvActive();
    document.getElementById("home").classList.add("is-active")
    populateEntries(entries, catSelect, tagSelect)
}
function rmvActive() {
    let els = document.getElementsByClassName('is-active');
    console.log(els)
    while (els[0]) {
        els[0].classList.remove('is-active')
    }
}

function catPicker(c) {
    rmvActive();
    catSelect = c;
    document.getElementById(c).classList.add("is-active");
    populateEntries(entries, catSelect, "")
}
function tagPicker(t) {
    rmvActive();
    tagSelect = t;
    document.getElementById(t).classList.add("is-active");
    populateEntries(entries, "", tagSelect)
}
function getTags(t) {
    var tags = "<span class='column is-one-third tags is-marginless'>";
    t.forEach(function(item) {
        tags += "<span class='tag is-link'>" + item + "</span>";
    })
    tags += "</span>"
    return tags
}
function getLinks(l, t) {
    if (l != undefined){
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
        tags += "<span class='tag cardtag " + cl + "'>" + item + "</span>";
    })
    tags += "";
    return tags
}

function addNonFree(f) {
    if (f == false) {
        return "<span class='icon is-medium has-text-danger'><i class='fas fa-2x fa-ban'></i></span>"
    } else {
        return ""
    }
}

function addPdep(p) {
    if (p == true) {
        return "<span class='icon is-medium has-text-warning'><i class='fas fa-2x fa-exclamation-triangle'></i></span>"
    } else {
        return ""
    }
}