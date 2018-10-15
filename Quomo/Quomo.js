var timeout = null;
var counter = 0;
var gmMain = function() {

var sections_to_be_removed = ["Last viewed","External profiles","Tags","Segments"]
var all_sections = [];

var elem = document.getElementById("myContainer");
if (elem !=null) elem.parentNode.removeChild(elem);

if(timeout) {
    clearTimeout(timeout);
    timeout = null;
}

timeout = setTimeout(function() {
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


    $(".kv__value").each(function (){
	    var s = $('[data-attribute-id=email]').parent().parent().find('[data-is-interactive=true]').attr('data-value')
	    var bool1 = s.includes("@buildium.com")
	    var bool2 = s.includes("@happy.co")

	    var bool = bool1 || bool2
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

    if (button_text != ""){
        if(elem === null){
        console.log("Button")
            the_div.appendChild (zNode);
        }

        var container = document.getElementById ("myContainer");
        container.className +=" myContainer"
        
        document.getElementById ("myButton").addEventListener (
        "click", function(){
            zNode.innerHTML = dialog_text;
            if(elem!=null)
                elem.appendChild (zNode);}
            , false);
    }
    $('[data-attribute-id=email]').parent().parent().find('[data-is-interactive=true]').off('click').on('click',function(e){
        window.open('https://manage.happyco.com/admin/search?utf8=✓&query='+$(this).attr('data-value'));    
    }).css('color','#00c389');
    }, 3000); //Three seconds will elapse and Code will execute.
};
setTimeout(function() {
    $(".kv__value").each(function(){
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