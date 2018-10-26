const puppeteer = require('puppeteer');
const CREDS = require('./files/creds');
const player = require('play-sound')(opts = {});
const {performance} = require('perf_hooks');

//puppeteer variables
var browser;
var page;

//SheetJS variables
var workbook;
var first_sheet_name;
var worksheet;

//call main func
run();

//********************************************Main function********************************************
async function run () {
  browser = await puppeteer.launch({
    headless: false
  });
 
  page = await browser.newPage();
  await page.goto("https://www.linkedin.com");

  const USERNAME_SELECTOR = '#login-email';
  const PASSWORD_SELECTOR = '#login-password';
  const BUTTON_SELECTOR = '#login-submit';

  await page.click(USERNAME_SELECTOR);
  await page.keyboard.type(CREDS.username);

  await page.click(PASSWORD_SELECTOR);
  await page.keyboard.type(CREDS.password);
  
  await page.click(BUTTON_SELECTOR);
  await page.waitForNavigation();
    
  var t0 = performance.now();
try{
  if(typeof require !== 'undefined') XLSX = require('xlsx');

  workbook = XLSX.readFile('data.xlsx');
  first_sheet_name = workbook.SheetNames[0];
  worksheet = workbook.Sheets[first_sheet_name];

  for (var i = 2; i < 4; i++) {
    console.log("Row: "+i);

    var address = 'R'+i;
    var cell = workbook[address];
    var profileUnavailable = (cell ? cell.v : undefined);
    //if profile was unavailable last time and wasn't removed- skip the row
    if (profileUnavailable === "TRUE") {
      continue;
    }
    address = 'B'+i;
    cell = worksheet[address];
    var url = (cell ? cell.v : undefined);

    await page.goto(url);
    await page.waitFor(3 * 1000);
    await page.evaluate(_ => {
      window.scrollBy(0, window.innerHeight);
    });
    
    var data = "";  
    var multiPosition = false;
    var jobExists = false;

    //exit if the profile is unavailable/deleted
    if (await page.url().includes("linkedin.com/in/unavailable")) {
      setData('R'+i,"TRUE");
      setTodaysDate(i);
      continue;
    }

    //check if a job exists
    if (await page.evaluate(() => document.getElementsByClassName
      ('pv-entity__position-group-pager ember-view'))!==null) {
      log("Job is not null");
    //if it does, does it have any data
      if (await page.evaluateHandle(() => document.getElementsByClassName
      ('pv-entity__position-group-pager ember-view').textContent)!==null) {
          log("Job has some data too");
          jobExists = true
          //gett all information for the current/most recent job 
          data = await page.evaluateHandle(() => {
          return Array.from(document.getElementsByClassName('pv-entity__position-group-pager ember-view')).map(elem => elem.textContent.trim()).slice(0,1);
        });
      }
    }

    //convert json to string
    var str = JSON.stringify(await data.jsonValue())
    //if text starts with `["Company Name...` --> the (current) job has multiple positions
    log("str:"+str);
    if (str.trim().startsWith('["Company Name')) {
      multiPosition = true;
    }

    var sel = "";
    //if job exists get job specific information (else get basic information)
    if (jobExists && str!=="[]") {
      //only one position in the current job
      if (!multiPosition) {
        console.log("single position")

        //add title, company name, current job duration into excel
        sel = 'div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3';
        await page.waitFor(2 * 1000);
        var title = await getData(sel);
        if (title !== undefined) {
          title = title.replace('Title','').trim();
        }
        setData('J'+i,title);
       
        sel = 'div > h4 > span:nth-child(2)';
        var companyName = await getData(sel);
        log("Company Name:"+companyName)
        setData('L'+i,companyName);
        
        sel = 'div.pv-entity__summary-info.pv-entity__summary-info--v2 > div > h4:nth-child(2) > span.pv-entity__bullet-item-v2';
        var currentJobDuration = await getData(sel);
        log("Current Job Duration:"+currentJobDuration);
        setData('N'+i,currentJobDuration);
      }
      //multiple positions in the current job
      else{
        console.log("multiple positions")
        
        //add title, company name, current job duration into excel
        // sel = 'div > div > div.pv-entity__summary-info-v2.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top.mb2 > h3 > span:nth-child(2)';
        sel = 'div > div > div > h3 > span:nth-child(2)';
        var title = await getData(sel);
        log("I got this title:"+title);
        setData('J'+i,title);
          
        sel = 'div > h3 > span:nth-child(2)';
        var companyName = await getData(sel);
        log("Company Name:"+companyName);
        setData('L'+i,companyName);
        
        sel = 'div > div > div.pv-entity__summary-info-v2.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top.mb2 > div > h4:nth-child(2) > span.pv-entity__bullet-item-v2';
        var currentJobDuration = await getData(sel);
        log("Current Job Duration:"+currentJobDuration);
        setData('N'+i,currentJobDuration);
      }//end of else
    }//jobExists

    //get name
    var name = ""
    sel = 'div.pv-top-card-v2-section__info.mr5 > div:nth-child(1) > h1';
    var name = await getData(sel);
    name = name.trim();
    log("Name:"+name);
    
    //add phone location birthday into excel
    const clickElement = 'span.pv-top-card-v2-section__entity-name.pv-top-card-v2-section__contact-info.ml2'
    if (name !="") {//??note sure why - cause if a profile is unavailable if skip. what else it can be. but better check than not
      await page.click(clickElement);
      await page.waitFor(1 * 1000);

      sel = 'div > section.pv-contact-info__contact-type.ci-phone > ul > li';
      var phone = await getData(sel);
      if (phone !=null) {
        phone = phone.trim().replace('(Mobile)','').replace('(Home)','').replace('(Work)','').replace(' ','').trim();
      }
      log("Phone:"+phone);
      setData('G'+i,phone);

      sel = 'div.pv-top-card-v2-section__info.mr5 > h3';
      var location = await getData(sel);
      location = location.trim();
      log("Location:"+location);
      setData('H'+i,location);

      sel = 'div > section.pv-contact-info__contact-type.ci-birthday > div > span';
      var birthday = await getData(sel);
      log("Birthday:"+birthday);
      setData('Q'+i,birthday);
    }

    //add today's date to track when a profile was last visited
    setTodaysDate(i);

    //add in excel if name title company has changed
    valueChanged('D'+i,name,'I'+i);
    valueChanged('E'+i,title,'K'+i);
    valueChanged('F'+i,companyName,'M'+i);

    //add in excel if the current job is less than 3 months old
    if (currentJobDuration != "" && currentJobDuration != undefined) {
      var ym = currentJobDuration.split(" ");
      if (ym.length ==2) {
        currentJobDuration = currentJobDuration.replace('mos',"").trim();
        ym = currentJobDuration.split(" ");
        if (ym.length ==1) {
           var num =  parseInt(ym[0]);
           if (num<4) {
            setData('O'+i,"TRUE")
          }
        }
      }
    }
    //get today's date for output file name
    var today = new Date();
    var dd = today.getDate();
    var mm = today.getMonth()+1;
    var yyyy = today.getFullYear();

    if(dd<10) {
        dd = '0'+dd
    } 

    if(mm<10) {
        mm = '0'+mm
    } 

    today = mm + '-' + dd + '-' + yyyy;

    if (i%2==0) {
      XLSX.writeFile(workbook ,'output '+today+'.xlsx')
    }
  }//for
}//try
catch(error){
  console.log(error);
  /*player.play('./files/error.mp3', function(err){
  if (err) throw err
})*/
}
var t1 = performance.now();
log((t1-t0)/1000+ " seconds");
await browser.close();
/*player.play('./files/completed.mp3', function(err){
  if (err) throw err
})*/
}//run()

//********************************************Get data from LinkedIn********************************************
async function getData(selector) {
  var val = "Null";
  const result = await page.evaluate((selector) => {
    if (document.querySelector(selector) != null) {
        if (document.querySelector(selector).textContent != null) {
          val = document.querySelector(selector).textContent;
          val = val.trim();
          return val;
       }
    }
  }, selector);
  return result;
}

//********************************************Add data to cell********************************************
async function setData(writeCell, data) {
  if (!worksheet[writeCell]) {
     worksheet[writeCell] = {}
  }
  worksheet[writeCell].v = data;
}
//********************************************Set today's date as a profile's last visited date********************************************
async function setTodaysDate(i) {
  var today = new Date();
  var dd = today.getDate();
  var mm = today.getMonth()+1; //0 indexed
    var yyyy = today.getFullYear();
    if(dd<10) {
        dd = '0'+dd
    } 
    if(mm<10) {
        mm = '0'+mm
    } 
    var lastVisitedDate = dd + '/' + mm + '/' + yyyy;
    setData('P'+i,lastVisitedDate);
}

//********************************************Compare values********************************************
async function valueChanged(readCell,data, writeCell) {
  var cell = worksheet[readCell];
  var value = (cell ? cell.v : undefined);
  if (value !== data && (data !==""|| value !=="")) {
    setData(writeCell,"TRUE")
  }
}
//********************************************Log********************************************
async function log(value){
  // console.log(value);
}