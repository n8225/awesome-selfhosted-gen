var entries, langSelect = "";
var tagSelect = new Array;

var requestURL = 'static/list.min.json';
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

const namea = `<article class="media"><div class="media-content"><span class="field is-grouped is-grouped-multiline"><span class="control"><strong>`
const nameb = `</strong></span>`
const propri = `<span class="control"><a class="icon has-text-warning"><i class="fas fa-lg fa-exclamation-triangle"></i></a></span>`
const nonf = `<span class="control"><a class="icon has-text-warning"><i class="fas fa-lg fa-ban"></i></a></span>`
const date = `<span class="control"><span class="tags has-addons"><a class="tag is-light">Updated</a><a class="tag is-info">`
const stara = `</a></span></span><span class="control"><span class="tags has-addons"><a class="tag is-dark icon"><i class="fas fa fa-star"></i></a><a class="tag is-light">`
const starb = `</a></span></span>`
const linka = `<span class="control"><a href="`
const linkb = `" target="_blank" class=""><span class="icon has-text-link"><i class="fas fa-lg fa-`
const src = `code-branch"></i></span></a></span>`
const site = `external-link-alt"></i></span></a></span>`
const demo = `chevron-circle-right"></i></span></a></span>`
const client = `mobile-alt"></i></span></a></span>`
const entriesa = `</span><p>`
const entriesb = `</p></div></article>`

function populateEntries() {
    var x, txt = "";
    for (x in entries) {
        if ((!tagSelect.some(ele => !entries[x].T.includes(ele) || tagSelect === []) && (entries[x].La.includes(langSelect) || langSelect === ""))) {
            txt += namea + entries[x].N + nameb;
            if (entries[x].P !== undefined) {txt += propri;}
            if (entries[x].NF !== undefined) {txt += nonf;}
            txt += parseArr(entries[x].T, "tag");
            if (entries[x].stars !== undefined) {txt += date + entries[x].update + stara + entries[x].stars + starb;}
            txt += getL(entries[x].Li) + parseArr(entries[x].La, "lang");
            txt += linka + entries[x].Sr + linkb + src;
            if (entries[x].Si !== undefined) {txt +=linka + entries[x].Si + linkb + site;}
            if (entries[x].Dem !== undefined) {txt +=linka + entries[x].Dem + linkb + demo;}
            if (entries[x].CL !== undefined) {txt +=linka + entries[x].CL + linkb + client;}
            txt += entriesa + entries[x].D + entriesb
        }
    }
    document.getElementById("demo").innerHTML = txt;
    displayFilters();
}

function getDates(u, s) {

    if (u !== undefined) {
        return "<div class='field is-grouped is-grouped-multiline'><div class='control'><div class='tags has-addons'><span class='tag is-light'>Updated</span><span class='tag is-info'>" + u + "</span></div></div><div class='control'><div class='tags'><span class='icon has-text-dark'><i class='up fas fa-lg fa-star'></i></span><span class='tag is-light is-rounded'>" + s + "</span></div></div></div>"
    } else {
        return ""
    }
}

function populateAllEntries() {
    var x, txt = "";
    for (x in entries) {
        txt += namea + entries[x].N + nameb;
        if (entries[x].P !== undefined) {txt += propri;}
        if (entries[x].NF !== undefined) {txt += nonf;}
        txt += parseArr(entries[x].T, "tag", entries[x].N);
        if (entries[x].stars !== undefined) {txt += date + entries[x].update + stara + entries[x].stars + starb;}
        txt += getL(entries[x].Li, entries[x].N) + parseArr(entries[x].La, "lang", entries[x].N) + linka + entries[x].Sr + linkb + src;
        if (entries[x].Si !== undefined) {txt +=linka + entries[x].Si + linkb + site;}
        if (entries[x].Dem !== undefined) {txt +=linka + entries[x].Dem + linkb + demo;}
        if (entries[x].CL !== undefined) {txt +=linka + entries[x].CL + linkb + client;}
        txt += entriesa + entries[x].D + entriesb
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
function parseArr(e, t, n) {
    if ( e == null) {
        console.log(n + " " + t + " Is null")
        return
    }
    switch (t) {
        case ("tag"):
            oc = `tagPicker`;
            c = `is-link`;
            break;
        case ("lang"):
            oc = `langPicker`;
            c = `is-success`;
            break;
    }
    let res = "";
    e.forEach(function(item) {
        res += `<span class="control"><a class="tag ` + c + `" onclick="` + oc + `(\`` + item + `\`)">` + item + `</a></span>`
    });
    return res;
}

function getL(t, n) {
    let lics = "";
    t.forEach(function(item) {
        lics += `<span class="control"><a class="tag is-primary">` + item + `</a></span>`;
    });
    lics += "";
    return lics
}
