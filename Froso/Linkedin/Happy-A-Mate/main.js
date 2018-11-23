const {app, BrowserWindow, ipcMain} = require('electron')

const puppeteer = require('puppeteer');
const {performance} = require('perf_hooks');
const {file} = require('fs');
const {TimeoutError} = require('puppeteer/Errors');

ipcMain.on('email', (event, email) => {
});

ipcMain.on('password', (event, password) => {
});

ipcMain.on('excel', (event, excel) => {
});

ipcMain.on('start row', (event, startRow) => {
});

ipcMain.on('end row', (event, endRow) => {
});

let mainWindow

app.on('ready', createWindow)

app.on('window-all-closed', function () {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

app.on('activate', function () {
  if (mainWindow === null) {
    createWindow()
  }
})

//********************************************Creating main window********************************************
function createWindow () {
  mainWindow = new BrowserWindow({width: 800, height: 600})
  console.log("App launched ");
  mainWindow.loadFile('index.html')
  
  mainWindow.on('closed', function () {  
    mainWindow = null
    console.log("window's closed game over");
  })
}

//********************************************Calling puppeteer code********************************************
function smashIt() {
  run().then(result => {
    log("Running something:"+result)
  }).catch(error => {
    log("Error running something:"+error)
  })
}

//********************************************Puppeteer code********************************************

//puppeteer variables
var browser;
var page;

//SheetJS variables
var workbook;
var first_sheet_name;
var worksheet;

//********************************************Main function********************************************
async function run () {
  log("Run. Run for your life");
  browser = await puppeteer.launch({
    headless: false
  });

  page = await browser.newPage();

  await page.setViewport({
      width: window.screen.availWidth,
      height: window.screen.availHeight
  });

  await page.goto("https://www.linkedin.com");

  const USERNAME_SELECTOR = '#login-email';
  const PASSWORD_SELECTOR = '#login-password';
  const BUTTON_SELECTOR = '#login-submit';

  await page.click(USERNAME_SELECTOR);
  await page.keyboard.type(email);
  await page.click(PASSWORD_SELECTOR);
  await page.keyboard.type(password);
  await page.click(BUTTON_SELECTOR);
  await page.waitForNavigation();

  var t0 = performance.now();
try{
  if(typeof require !== 'undefined') XLSX = require('xlsx');

  workbook = XLSX.readFile(excel);
  first_sheet_name = workbook.SheetNames[0];
  worksheet = workbook.Sheets[first_sheet_name];
  
  var currentProfile = 0;
  var totalNumberOfProfiles = endRow-startRow+1;
  
  //for 
  for (var i = Number(startRow); i <= Number(endRow); i++) {
    var currentProfile1 = document.getElementById("currentProfile");
    currentProfile1.textContent = ++currentProfile;
    
    var timeRemaining = document.getElementById("timeRemaining");
    var time = Math.round(parseInt((totalNumberOfProfiles- currentProfile)/9.5));
    if (time == 0) {
      time = "< 1";
    }
    timeRemaining.textContent = time;



    var totalRecordsDonePercent = document.getElementById("recordsDonePercent");
    totalRecordsDonePercent.textContent = parseInt((currentProfile-1)*100/totalNumberOfProfiles);

    log("Row:"+i)

    await page.setViewport({
      width: window.screen.availWidth,
      height: window.screen.availHeight
    });

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
    await page.evaluate(_ => {
      window.scrollBy(0, window.innerHeight);
    });
    var singleAll = false;
    try{
      await page.waitForSelector('div.pv-entity__summary-info.pv-entity__summary-info--background-section.mb2 > h3',{timeout:2000})
      .then(() => singleAll = true);
    }
    catch(error){
      try{
        await page.waitForSelector('div.pv-entity__summary-info.pv-entity__summary-info--background-section > h3',{timeout:2000})
      }
      catch(error){
        try{
          await page.waitForSelector('div > div > div > h3 > span:nth-child(2)',{timeout:2000})
        }
        catch(error){
        }
      }
    }
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
        if (singleAll) {
          sel = 'div.pv-entity__summary-info.pv-entity__summary-info--background-section.mb2 > h3'
        }
        else{
          sel = 'div.pv-entity__summary-info.pv-entity__summary-info--background-section > h3';
        }
        // await page.waitFor(2 * 1000);
        var title = await getData(sel);
        if (title !== undefined) {
          title = title.replace('Title','').trim();
        }
        setData('J'+i,title);
       
        sel = 'div > h4 > span:nth-child(2)';
        var companyName = await getData(sel);
        log("Company Name:"+companyName)
        setData('L'+i,companyName);
        if (singleAll) {
          sel = 'div.pv-entity__summary-info.pv-entity__summary-info--background-section.mb2 > div > h4:nth-child(2) > span.pv-entity__bullet-item-v2'
        }        
        else {
          sel = 'div.pv-entity__summary-info.pv-entity__summary-info--background-section > div > h4:nth-child(2) > span.pv-entity__bullet-item-v2'
        }
        var currentJobDuration = await getData(sel);
        log("Current Job Duration:"+currentJobDuration);
        setData('N'+i,currentJobDuration);
      }
      //multiple positions in the current job
      else{
        console.log("multiple positions")
        
        //add title, company name, current job duration into excel
        sel = 'div > div > div > h3 > span:nth-child(2)';
        // await page.waitFor(2 * 1000);
        var title = await getData(sel);
        log("I got this title:"+title);
        setData('J'+i,title);
          
        sel = 'div > h3 > span:nth-child(2)';
        var companyName = await getData(sel);
        log("Company Name:"+companyName);
        setData('L'+i,companyName);
        
        sel = 'div > div > div > div > h4:nth-child(2) > span.pv-entity__bullet-item-v2'
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
    if (name !="") {//??not sure why - cause if a profile is unavailable skip. what else can it be. but better check than not
      await page.click(clickElement);
      // await page.waitFor(1 * 1000);
      await page.waitForSelector('div > section.pv-contact-info__contact-type.ci-vanity-url > div > a', {timeout:1000});

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

    log("Today's date:"+today)
    var pathTo = excel.substring(0,excel.lastIndexOf("/")+1);
    XLSX.writeFile(workbook ,pathTo+'/output '+today+'.xlsx')
  }//for
   totalRecordsDonePercent.textContent = 100;
   timeRemaining.textContent = 0;
}//try
catch(error){
  alert("Caught an error:"+error);
}
var t1 = performance.now();
log((t1-t0)/1000+ " seconds");
await browser.close();
if (today !=undefined) {
  alert("Please check output file (output " +today +".xlsx) in source file directory")
}
var elements = ["email", "password", "excel", "start", "end", "signin"]
for (var i = 0; i < elements.length; i++) {
  let element = document.getElementById(elements[i])
  element.disabled = false;
}
}//run()

//********************************************Get data from LinkedIn********************************************
async function getData(selector) {
  var val = "Null";
  const result = await page.evaluate((selector) => {
    if (document.querySelector(selector) != null) {
        if (document.querySelector(selector).textContent != null) {
          val = document.querySelector(selector).textContent;
          return val.trim();
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
    var lastVisitedDate = mm + '/' + dd + '/' + yyyy;
    setData('P'+i,lastVisitedDate);
}

//********************************************Compare values********************************************
async function valueChanged(readCell,data, writeCell) {
  var cell = worksheet[readCell];
  var value = (cell ? cell.v : undefined);
  if (value !== data && data !== "") {
    setData(writeCell, "TRUE")
    setData(readCell, data)
  }
}

//********************************************Log********************************************
async function log(value){
  console.log(value);
}