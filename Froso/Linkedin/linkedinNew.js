const puppeteer = require('puppeteer');
const CREDS = require('./creds');
const player = require('play-sound')(opts = {});
const {performance} = require('perf_hooks');

run();

var browser;
var page;
var workbook;
var first_sheet_name;
var worksheet;

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
    
if(typeof require !== 'undefined') XLSX = require('xlsx');

  workbook = XLSX.readFile('data.xlsx');
  first_sheet_name = workbook.SheetNames[0];
  worksheet = workbook.Sheets[first_sheet_name];
  var address_of_cell = 'A2';
  var desired_cell = worksheet[address_of_cell];

  var t0 = performance.now();
try{
  for (var i = 2; i < 11; i++) {
    console.log("Row: "+i);
    address_of_cell = 'B'+i;
    desired_cell = worksheet[address_of_cell];
    var desired_value = (desired_cell ? desired_cell.v : undefined);

    await page.goto(desired_value); 
    await page.evaluate(_ => {
      window.scrollBy(0, window.innerHeight);
    });
    await page.waitFor(2 * 1000);
    var data = "";  
    var multiPosition = false
    var jobExists = false

    //exiting if the profile is unavailable/deleted
    if (await page.url() === "https://www.linkedin.com/in/unavailable/") {
      setData('R'+i,"Yes");
      setTodaysDate(i);
      continue;
    }

    //getting current job - all information
    if (await page.evaluate(() => document.getElementsByClassName
      ('pv-entity__position-group-pager ember-view'))!==null) {
      // console.log("Job is not null");
      if (await page.evaluateHandle(() => document.getElementsByClassName
      ('pv-entity__position-group-pager ember-view').textContent)!==null) {
          // console.log("Job has some data too");
          jobExists = true
          data = await page.evaluateHandle(() => {
          return Array.from(document.getElementsByClassName('pv-entity__position-group-pager ember-view')).map(elem => elem.textContent.trim()).slice(0,1);
        });
      }
    }

    //converting json to string
    var str = JSON.stringify(await data.jsonValue())
    //if text starts with `["Company Name...` --> the (current) job has multiple positions
    // console.log("str:"+str);
    if (str.trim().startsWith('["Company Name')) {
      multiPosition = true;
    }
    var sel = "";
    //if job exists we can get job specific information else get basic information
    if (jobExists) { 
      //only one position in the current job
      if (!multiPosition) {
        console.log("single position")

        //adding title, company name, current job duration
        sel = 'div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3';
        await page.waitFor(2 * 1000);
        var title = await getData(sel);
        title = title.replace('Title','').trim();
        console.log("Title:"+title)
        setData('J'+i,title);
       
        sel = 'div > h4 > span:nth-child(2)';
        var companyName = await getData(sel);
        console.log("Company Name:"+companyName)
        setData('L'+i,companyName);
        
        sel = 'div.pv-entity__summary-info.pv-entity__summary-info--v2 > div > h4:nth-child(2) > span.pv-entity__bullet-item-v2';
        var currentJobDuration = await getData(sel);
        console.log("Current Job Duration:"+currentJobDuration);
        setData('N'+i,currentJobDuration);
      }
    //multiple positions in the current job
      else{
        console.log("muliple positions")

        sel = 'div > div > div.pv-entity__summary-info-v2.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top.mb2 > h3 > span:nth-child(2)';
        var title = await getData(sel);
        console.log("I got this title:"+title);
        setData('J'+i,title);
          
        sel = 'div > h3 > span:nth-child(2)';
        var companyName = await getData(sel);
        console.log("Company Name:"+companyName);
        setData('L'+i,companyName);    

        sel = 'div > div > div.pv-entity__summary-info-v2.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top.mb2 > div > h4:nth-child(2) > span.pv-entity__bullet-item-v2';
        var currentJobDuration = await getData(sel);
        console.log("Current Job Duration:"+currentJobDuration)
        setData('N'+i,currentJobDuration);
      }//end of else
    }//jobExists

    //getting name
    var name = ""
    sel = 'div.pv-top-card-v2-section__info.mr5 > div:nth-child(1) > h1';
    var name = await getData(sel);
    name = name.trim();
    console.log("Name:"+name);
    
    //adding location phone birthday
    const clickElement = 'span.pv-top-card-v2-section__entity-name.pv-top-card-v2-section__contact-info.ml2'
    if (name !="") {
      await page.click(clickElement)
      await page.waitFor(2 * 1000); 

      sel = 'div > section.pv-contact-info__contact-type.ci-phone > ul > li';
      var phone = await getData(sel);
      if (phone !=null) {
        phone = phone.trim().replace('(Mobile)','').replace('(Home)','').replace('(Work)','').replace(' ','').trim();        
      }
      console.log("Phone:"+phone);
      setData('G'+i,phone);

      sel = 'div.pv-top-card-v2-section__info.mr5 > h3';
      var location = await getData(sel);
      location = location.trim();
      console.log("Location:"+location);
      setData('H'+i,location);

      sel = 'div > section.pv-contact-info__contact-type.ci-birthday > div > span';
      var birthday = await getData(sel);
      console.log("Birthday:"+birthday);
      setData('Q'+i,birthday);
    }

    //adding today's date to track when was profile last visited
    setTodaysDate(i);

    //check if name title company job duration has changed
    //if either the name fetched last and now is empty- don't do anything
    
    hasValueChanged('D'+i,name,'I'+i);

    /*address_of_cell = 'D'+i;
    desired_cell = worksheet[address_of_cell];
    desired_value = (desired_cell ? desired_cell.v : undefined);
    if (desired_value !== name && (name !==""|| desired_value !=="")) {
      setData('I'+i,"Yes")
    }*/

    hasValueChanged('E'+i,title,'K'+i);
    /*address_of_cell = 'E'+i;
    desired_cell = worksheet[address_of_cell];
    desired_value = (desired_cell ? desired_cell.v : undefined);

    if (desired_value !== title && (title !==""|| desired_value!=="")) {
     setData('K'+i,"Yes")
    }*/
    hasValueChanged('F'+i,companyName,'M'+i);

/*    address_of_cell = 'F'+i;
    desired_cell = worksheet[address_of_cell];
    desired_value = (desired_cell ? desired_cell.v : undefined);
    if (desired_value !== companyName && (companyName !==""|| desired_value!=="")) {
      setData('M'+i,"Yes")
    }*/

    if (currentJobDuration != "" && currentJobDuration != undefined) {
      // console.log(currentJobDuration);
    //checking if the current job is less than 3 months old
      var ym = currentJobDuration.split(" ");
      if (ym.length ==2) {
        currentJobDuration = currentJobDuration.replace('mos',"").trim();
        ym = currentJobDuration.split(" ");
        if (ym.length ==1) {
           var num =  parseInt(ym[0]);
           if (num<4) {
            setData('O'+i,"Yes")
          }
        }
      }
    }
    // if (i%2==0) {
      XLSX.writeFile(workbook ,'output.xlsx')
    // }
  }
}
catch(error){
      console.log(error);
      player.play('./files/error.mp3', function(err){
      if (err) throw err
    })
}
  var t1 = performance.now();
  // console.log((t1-t0)/1000+ " seconds");

  await browser.close();
  player.play('./files/completed.mp3', function(err){
      if (err) throw err
    })
}

async function getData(selector) {
  // console.log("getData");
  var resultsString = "Null";
  const result = await page.evaluate((selector) => {
    if (document.querySelector(selector) != null) {
        if (document.querySelector(selector).textContent != null) {
          resultsString = document.querySelector(selector).textContent;
          resultsString = resultsString.trim();
          return resultsString;
       }
    }
  }, selector);
  if (result == undefined || result == "undefined") {
      // console.log(result);
  }
  return result;
}

async function setData(writeCell, data) {
  // console.log("setData");

  if (!worksheet[writeCell]) {
     worksheet[writeCell] = {}
  }
  worksheet[writeCell].v = data;
}
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

async function hasValueChanged(readCell,data, writeCell) {
  var cell = worksheet[readCell];
  var value = (cell ? cell.v : undefined);
  if (value !== data && (data !==""|| value !=="")) {
    setData(writeCell,"Yes")
  }
}