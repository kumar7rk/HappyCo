// ==UserScript==
// @name         Quomo
// @namespace    https://app.intercom.io/a/apps/yaqkh6zy
// @version      1.0.2
// @description  See notes below
// @author       You
// @match        https://app.intercom.io/a/apps/yaqkh6zy/*
// @grant        GM_addStyle
// @require      http://code.jquery.com/jquery-3.2.1.js

// ==/UserScript==
//relase notes
//1.0 --> Shows a button for a tier 2 user, hides unwanted sections and some unknown values.
//1.0.1 --> If multiple dialogs are added, removing all at once when the next conversation is non-tier2
//1.0.2 --> Adding dialog only once and also when it's clicked the info dialog is added once. Did some testing working consistently. Not too sure about the later one.
(function() {
//    'use strict';
//--- Style our newly added elements using CSS.
GM_addStyle ( multilineStr ( function () {/*!
    #myContainer {
        position:               relative;
        top:                    20px;
        font-size:              20px;
        background:             #00c389;
        border:                 3px outset blue;
        margin:                 5px;
        opacity:                0.9;
        z-index:                1100;
        padding:                5px 20px;
    }
    #myButton {
        cursor:                 pointer;
    }

*/} ) );
    function multilineStr (dummyFunc) {
    var str = dummyFunc.toString ();
    str = str.replace (/^[^\/]+\/\*!?/, '') // Strip function () { /*!
        .replace (/\s*\*\/\s*\}\s*$/, '') // Strip */ }
        .replace (/\/\/.+$/gm, '') // Double-slash comments wreck CSS. Strip them.
    ;
    return str;
}
    var gmMain = function() {

        var sections_to_be_removed = ["Last viewed","External profiles","Tags","Segments"]
        var all_sections = [];
        // waiting for three seconds and then hides unused sidebar elements such as last viewed external profiles, tags, segments
        // is last viewed exists the code kind of breaks
//        window.onload = function(){});

        var elem = document.getElementById("myContainer");
        if (elem !=null) elem.parentNode.removeChild(elem);


        setTimeout(function() {
            var button_text = ""
            var dialog_text = ""
            var elements = document.getElementsByClassName('profile__sidebar-section ember-view');
            var user_type = "";
            var section_name = "Section Name= ";

            $(".stamp.o__user").filter(function(){
                user_type = $(this).text();
                return null;
            });
            $(".profile-sidebar-section__section-title").filter(function(){
                section_name = $(this).text().trim();
                if (section_name.includes("Notes")) all_sections.push("Notes")
                else if (section_name.includes("Tags")) all_sections.push("Tags")
                else if (section_name.includes("Segments")) all_sections.push("Segments")
                else all_sections.push(section_name)
                return null;
            });
            all_sections.forEach(function(item, index, array) {
                var index_present_at_index = sections_to_be_removed.indexOf(item)
                if(index_present_at_index>-1) elements[index].style.display = 'none';
            });

            // checking for tier 2 support
            // Buildium, ACH, Equity, Colony Starwood, Freddie-Mac
            // checking for due diligence users
            var s = $('[data-attribute-id=email]').parent().parent().find('[data-is-interactive=true]').attr('data-value')
            var bool1 = s.includes("@buildium.com")
            var bool2 = s.includes("@happy.co")

            var bool = bool1 || bool2

            $(".kv__value").filter(function (){
                var text_here = $(this).text().trim()
                if((text_here === "buildium") && !bool){
                    button_text = "Buildium"
                    dialog_text = "I'm Buildium. Only reply to domain @buildium.com :)"
                }
                else if(text_here === "29630"){
                    button_text = "ACH"
                    dialog_text = "I'm ACH. Only reply if I'm JT Bailey. :)"
                }
                else if(text_here === "4722"){
                    button_text = "Equity"
                    dialog_text = "I'm Equity. Tier2 support. Only reply if I'm Loren Lizotte or Ed Leigh, Jennifer Henkel :)"
                }
                else if(text_here === "20477"){
                    button_text = "Colony Starwood"
                    dialog_text = "I'm Colony Starwood. Tier2 support. Only reply If I'm Denise Wesel. CAN ADD/REMOVE USERS: Michael Williams, Leslie Hunt, Malorie Iglesias, Lisa Sasik, Alisha Gardner, Anthony Roy"
                }
                else if(text_here === "36577"){
                    button_text = "Freddie Mac"
                    dialog_text = "I'm Freddie Mac. <We don't have an 'admin' list yet>. Please don't reply to push inspectors from Sabal, Red Capital, ReadyCap, Pinnacle Bank, Hunt, Greystone, Freddie Mac, CPC, CBRE, Capitol One, Basis, Arbor:)"
                }
                else if(text_here === "due_diligence"){
                    button_text = "Due Diligence"
                    dialog_text = "I'm Due Diligence. Please don't advice me to the regular user work like create a report, browse template library etc :)"
                }
            });
            var the_div = document.getElementsByClassName("u__centered-text-block u__mt__20")[0]
            var zNode = document.createElement ('div');
            zNode.innerHTML = '<button id="myButton" type="button" >' +button_text +'</button>';
            zNode.setAttribute ('id', 'myContainer');
            var elem = document.getElementById("myContainer");

/*            if (button_text === ""){
                for (var i = 0; i< 10; i++){
//                    var elem = document.getElementById("myContainer");
                    console.log(i+"");
                    if (elem !=null)
                        elem.parentNode.removeChild(elem);
                }
            }
*/
            if (button_text != ""){
                if(elem === null){
                    the_div.appendChild (zNode);
                }

                document.getElementById ("myButton").addEventListener (
                "click", function(){
                    zNode.innerHTML = dialog_text;
                    if(elem!=null)
                        elem.appendChild (zNode);}
                    , false);
            }
        }, 5000); //Three seconds will elapse and Code will execute.
    };
    // waiting for 12 seconds; hides all the unknown values in details, company details
    // if you click "show x hidden" quick enough it would hide hidden unknown as well
    setTimeout(function() {
        $(".kv__value").filter(function(){
            return $(this).text().trim() === "Unknown";
        }).parent().hide();
    }, 12000); //Twelve seconds will elapse and Code will execute.

    var fireOnHashChangesToo = true;
    var pageURLCheckTimer = setInterval(function () {
        if (this.lastPathStr !== location.pathname || this.lastQueryStr !== location.search || (fireOnHashChangesToo && this.lastHashStr !== location.hash)) {
            this.lastPathStr = location.pathname;
            this.lastQueryStr = location.search;
            this.lastHashStr = location.hash;
            gmMain ();
        }
    } , 111);
})();