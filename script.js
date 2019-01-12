var entries, langSelect = "";
var tagSelect = new Array;


var requestURL = 'output.min.json';
var request = new XMLHttpRequest();
request.open('GET', requestURL);
request.responseType = 'json';
request.send();

request.onload = function() {
        if (window.matchMedia("(max-width: 752px)").matches) {
            document.getElementById("panel-l").checked = false;
            document.getElementById("panel-t").checked = false;
        }
    entries = request.response.Entries;
    var tags = request.response.Tags;
    var langs = request.response.Langs;
    populateTags(tags);
    populateLangs(langs);
    populateAllEntries();

};
function clrLang() {
    langSelect = "";
    populateEntries();
}

function remove(array, element) {
    const index = array.indexOf(element);
    array.splice(index, 1);
}

function clrTags(t) {
    remove(tagSelect, t);
    populateEntries();
}

function displayFilters() {
    var del = "<span class='tag is-danger level-item' onclick='goHome()'>Clear All</span>";
    var lang = "<span class='tag is-success'>" + langSelect + "<button onclick='clrLang()' class='delete is-small'></button></span>";
    var tags = "";
    for (i in tagSelect) {
        tags += "<span class='tag is-link'>" + tagSelect[i] + "<button onclick='clrTags(`" + tagSelect[i] + "`)' class='delete is-small'></button></span>";
    }
    switch (true) {
        case (langSelect !== "" && tagSelect[0] !== undefined):
            document.getElementById("filters").innerHTML = del + lang + tags;
            document.getElementById("filters").classList.add("notification", "tags");
            break;
        case langSelect !== "":
            document.getElementById("filters").innerHTML = del + lang;
            document.getElementById("filters").classList.add("notification", "tags");
            break;
        case tagSelect[0] !== undefined:
            document.getElementById("filters").innerHTML = del + tags;
            document.getElementById("filters").classList.add("notification", "tags");
            break;
        default:
            document.getElementById("filters").innerHTML = "";
            document.getElementById("filters").classList.remove("notification", "tags");
    }
};

function populateLangs(lObj) {
    var z, txt3 = "";
    for (z in lObj) {
        txt3 += "<span id='" + lObj[z].Lang + "' class='tag is-success' onClick='langPicker(`" + lObj[z].Lang + "`)'>" + lObj[z].Lang + "</span>"
    }
    document.getElementById("lang").innerHTML = txt3;
}

function populateTags(tObj) {
    var z, txt = "";
    for (z in tObj) {
        txt += "<span id='" + tObj[z].Tag + "' class='tag is-link' onClick='tagPicker(`" + tObj[z].Tag + "`)'>" + tObj[z].Tag + "</span>"
    }
    document.getElementById("tagList").innerHTML = txt;
}

function populateEntries() {
    var x, txt = "";

    for (x in entries) {
        if ((!tagSelect.some(ele => !entries[x].T.includes(ele) || tagSelect === []) && (entries[x].La.includes(langSelect) || langSelect === ""))) {
            txt += "<div class='card' style='margin-bottom:24px'><header class='columns card-header is-marginless has-background-light'>";
            if (entries[x].NF !== undefined ||  entries[x].P !== undefined){txt += "<span class='column is-narrow is-narrow-mobile'>" + addNonFree(entries[x].NF) + addPdep(entries[x].P) + "</span>"};
            txt += "<span class='column'><h4 class='title is-4'>" + entries[x].N + "</h4></span>";
            txt += "<span class='column is-narrow tags'>" + getL(entries[x].Li, "is-primary") + getL(entries[x].La, "is-success") + "</span></header>";
            txt += "<span class='card-footer'><span class='column is-one-third tags'>" + getTags(entries[x].T) + "</span>";
            txt += "<span class='column'><p>" + entries[x].D + "</p>" + getDates(entries[x].create, entries[x].update, entries[x].stars) + "</span></span>";
            txt += "<span class='level'>" + getLinks(entries[x].Sr, "Source Code") + getLinks(entries[x].Si, "Website") + getLinks(entries[x].Dem, "Demo") + getLinks(entries[x].CL, "Clients") + "</span>";
            txt += "</div>";
        }
    }
    document.getElementById("demo").innerHTML = txt;
    displayFilters();
}

function getDates(c, u, s) {

    if (c !== undefined && u !== undefined) {
        return "<div class='field is-grouped is-grouped-multiline'><div class='control'><div class='tags has-addons'><span class='tag is-light'>Created</span><span class='tag is-info'>" + c + "</span></div></div><div class='control'><div class='tags has-addons'><span class='tag is-light'>Updated</span><span class='tag is-info'>" + u + "</span></div></div><div class='control'><div class='tags'><span class='icon has-text-dark'><i class='up fas fa-lg fa-star'></i></span><span class='tag is-light is-rounded'>" + s + "</span></div></div></div>"
    } else {
        return ""
    }
}

function populateAllEntries() {
    var x, txt = "";
    for (x in entries) {
        txt += "<div class='card' style='margin-bottom:24px'><header class='columns card-header is-marginless has-background-light'>";
        if (entries[x].NF !== undefined ||  entries[x].P !== undefined){txt += "<span class='column is-narrow is-narrow-mobile'>" + addNonFree(entries[x].NF) + addPdep(entries[x].P) + "</span>"};
        txt += "<span class='column'><h4 class='title is-4'>" + entries[x].N + "</h4></span>";
        txt += "<span class='column is-narrow tags'>" + getL(entries[x].Li, "is-primary") + getL(entries[x].La, "is-success") + "</span></header>";
        txt += "<span class='card-footer'><span class='column is-one-third tags'>" + getTags(entries[x].T) + "</span>";
        txt += "<span class='column'><p>" + entries[x].D + "</p>" + getDates(entries[x].create, entries[x].update, entries[x].stars) + "</span></span>";
        txt += "<span class='level'>" + getLinks(entries[x].Sr, "Source Code") + getLinks(entries[x].Si, "Website") + getLinks(entries[x].Dem, "Demo") + getLinks(entries[x].CL, "Clients") + "</span>";
        txt += "</div>";
        }
    document.getElementById("demo").innerHTML = txt;
}

function goHome() {
    tagSelect = [];
    langSelect = "";
    rmvActive();
    populateEntries()
}
function rmvActive() {
    let els = document.getElementsByClassName('is-active');
    while (els[0]) {
        els[0].classList.remove('is-active')
    }
}
function tagPicker(t) {
    tagSelect.push(t);
    populateEntries()
}
function langPicker(l) {
    rmvActive();
    langSelect = l;
    populateEntries()
}
function getTags(t) {
    let tags = "";
    t.forEach(function(item) {
        tags += "<span class='tag is-link' onclick='tagPicker(`" + item + "`)'>" + item + "</span>";
    });
    return tags
}
function getLinks(l, t) {
    if (l !== undefined){
        switch (t) {
            case "Source Code":
                return "<a href='" + l + "'target='_blank' class='level-item'><span class='icon has-text-link'><i class='fas fa-lg fa-code-branch'></i></span>" + t + "</a>";
            case "Website":
                return "<a href='" + l + "'target='_blank' class='level-item'><span class='icon has-text-link'><i class='fas fa-lg fa-external-link-alt'></i></span>" + t + "</a>";
            case "Demo":
                return "<a href='" + l + "'target='_blank' class='level-item'><span class='icon has-text-link'><i class='fas fa-lg fa-chevron-circle-right'></i></span>" + t + "</a>";
            case "Clients":
                return "<a href='" + l + "'target='_blank' class='level-item'><span class='icon has-text-link'><i class='fas fa-lg fa-external-link-alt'></i></span>" + t + "</a>";
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