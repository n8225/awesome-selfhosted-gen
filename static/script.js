var list, langSelect = "";
var tagSelect = new Array;

window.onload = function() {
    if (window.matchMedia("(max-width: 752px)").matches) {
        document.getElementById("panel-l").checked = false;
        document.getElementById("panel-t").checked = false;
    }
}
    fetch('static/list.min.json')
        .then(res => res.json())
            .then(res => {
                list = res
                populateTags(list.Tags);
                populateLangs(list.Langs);
                displayFilters();
                populateEntries();
    })
function clrLang() {
    langSelect = "";
    displayFilters();
    populateEntries();
}
function remove(array, element) {
    const index = array.indexOf(element);
   array.splice(index, 1);
}
function clrTags(t) {
   remove(tagSelect, t);
   displayFilters();
    populateEntries();
}
function displayFilters() {
    const del = "<span class='tag is-danger' onclick='goHome()'>Clear All</span>";
    var lang = "<span class='tag is-success'>" + langSelect + "<button onclick='clrLang()' class='delete is-small'></button></span>";
    var tags = "";
    for (i in tagSelect) {
        tags += "<span class='tag is-link'>" + tagSelect[i] + "<button onclick='clrTags(`" + tagSelect[i] + "`)' class='delete is-small'></button></span>";
    }
    switch (true) {
        case (langSelect !== "" && tagSelect[0] !== undefined):
            document.getElementById("filters").innerHTML = del + lang + tags;
            document.getElementById("filters").classList.add("notification");
            break;
        case langSelect !== "":
            document.getElementById("filters").innerHTML = del + lang;
            document.getElementById("filters").classList.add("notification");
            break;
        case tagSelect[0] !== undefined:
            document.getElementById("filters").innerHTML = del + tags;
            document.getElementById("filters").classList.add("notification");
            break;
        default:
            document.getElementById("filters").innerHTML = "";
            document.getElementById("filters").classList.remove("notification");
    }
};
function populateLangs(lObj) {
    var z, txt3 = "";
    for (z in lObj) {
        txt3 += "<span id='" + lObj[z].Lang + "' class='tag is-success' onClick='langPicker(`" + lObj[z].Lang + "`)'>" + lObj[z].Lang + "</span>"
    }
    document.getElementById("lang").innerHTML = txt3;
    displayFilters();
}
function populateTags(tObj) {
    var z, txt = "";
    for (z in tObj) {
        txt += "<span id='" + tObj[z].Tag + "' class='tag is-link' onClick='tagPicker(`" + tObj[z].Tag + "`)'>" + tObj[z].Tag + "</span>"
    }
    document.getElementById("tagList").innerHTML = txt;
    displayFilters();
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
    var x, y, txt = "";
    document.getElementById("demo").innerHTML = ""
    for (y in list.Cats) {
        var ctxt = "<div class='box'><article class='message'><div class='message-header'><p>" + list.Cats[y].Cat + "</p></div></article>"
        var etxt = ""
        for (x in list.Entries) {
            if (list.Entries[x].C === list.Cats[y].Cat) {
                if (!tagSelect.some(ele => !list.Entries[x].T.includes(ele) || tagSelect === []) && (list.Entries[x].La.includes(langSelect) || langSelect === "")) {
                    etxt += namea + list.Entries[x].N + nameb;
                    if (list.Entries[x].P !== undefined) {etxt += propri;}
                    if (list.Entries[x].NF !== undefined) {etxt += nonf;}
                    etxt += parseArr(list.Entries[x].T, "tag");
                    if (list.Entries[x].stars !== undefined) {etxt += date + list.Entries[x].update + stara + list.Entries[x].stars + starb;}
                    etxt += getL(list.Entries[x].Li) + parseArr(list.Entries[x].La, "lang");
                    etxt += linka + list.Entries[x].Sr + linkb + src;
                    if (list.Entries[x].Si !== undefined) {etxt +=linka + list.Entries[x].Si + linkb + site;}
                    if (list.Entries[x].Dem !== undefined) {etxt +=linka + list.Entries[x].Dem + linkb + demo;}
                    if (list.Entries[x].CL !== undefined) {etxt +=linka + list.Entries[x].CL + linkb + client;}
                    etxt += entriesa + list.Entries[x].D + entriesb
                }
            }    
        }
        if (etxt !== "") {
            document.getElementById("demo").innerHTML += ctxt + etxt + "</div>"
        }
    }
}

function getDates(u, s) {

    if (u !== undefined) {
        return "<div class='field is-grouped is-grouped-multiline'><div class='control'><div class='tags has-addons'><span class='tag is-light'>Updated</span><span class='tag is-info'>" + u + "</span></div></div><div class='control'><div class='tags'><span class='icon has-text-dark'><i class='up fas fa-lg fa-star'></i></span><span class='tag is-light is-rounded'>" + s + "</span></div></div></div>"
    } else {
        return ""
    }
}
function goHome() {
    tagSelect = [];
    langSelect = "";
    displayFilters();
    populateEntries()
    //rmvActive();
}
// function rmvActive() {
//     let els = document.getElementsByClassName('is-active');
//      while (els[0]) {
//          els[0].classList.remove('is-active')
//      }
// }
function tagPicker(t) {
    tagSelect.push(t);
    displayFilters();
    populateEntries()
}
function langPicker(l) {
    //rmvActive();
    langSelect = l;
        displayFilters();
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
