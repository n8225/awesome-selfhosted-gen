var entries, catSelect = "", tagSelect = "";
var requestURL = 'output.json';
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
        txt2 += "<li><a onClick='catPicker(`" + cObj[y].Cat + "`)' href='#'>" + cObj[y].Cat + "</a></li>"
    }
    document.getElementById("cat").innerHTML = txt2;
}

function populateTags(tObj) {
    var z, txt = "";
    for (z in tObj) {
        txt += "<span class='tag is-link' onClick='tagPicker(`" + tObj[z].Tag + "`)'>" + tObj[z].Tag + "</span>"
    }
    document.getElementById("tagList").innerHTML = txt;
}
function populateEntries(eObj, catSelect, tagSelect) {
    var x, txt = "";
    console.log(catSelect, " ", tagSelect)
    for (x in eObj) {
        if (eObj[x].C == catSelect || eObj[x].T.includes(tagSelect) || (tagSelect == "" && catSelect == "")) {
            txt += "<div class='card' style='margin-bottom:24px'><header class='columns card-header is-marginless has-background-light'>" + "<span class='column'><h5 class='title is-5'>" + addNonFree(eObj[x].F) + eObj[x].N + "</h5></span><span class='column is-narrow tags is-marginless cardTitle'>" + addPdep(eObj[x].P) + getL(eObj[x].Li, "is-primary") + getL(eObj[x].La, "is-success") + "</span></header>"
            txt += "<div class='card-content columns'>" + getTags(eObj[x].T) + "<span class='column'>" + eObj[x].D + "</span></div>"
            txt += "<footer class='card-footer'>" + getLinks(eObj[x].Sr, "Source Code") + getLinks(eObj[x].Si, "Website") + getLinks(eObj[x].Dem, "Demo") + "</footer></div>";
        }
    }
    document.getElementById("demo").innerHTML = txt;
}

function catPicker(c) {
    catSelect = c;

    populateEntries(entries, catSelect, "")
}
function tagPicker(t) {
    tagSelect = t;
    populateEntries(entries, "", tagSelect)
}
function getTags(t) {
    var tags = "<span class='column is-one-quarter tags'>";
    t.forEach(function(item) {
        tags += "<span class='tag is-link'>" + item + "</span>";
    })
    tags += "</span>"
    return tags
}
function getLinks(l, t) {
    if (l !== undefined){
        return "<a href='" + l + "'class='card-footer-item'>" + t + "</a>"
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
    if (f == true) {
        return "<span class='icon is-medium has-text-danger warn'><i class='fas fa-2x fa-ban'></i></span>"
    } else {
        return ""
    }
}

function addPdep(p) {
    if (p == true) {
        return "<span class='icon is-medium has-text-warning warn'><i class='fas fa-2x fa-exclamation-triangle'></i></span>"
    } else {
        return ""
    }
}