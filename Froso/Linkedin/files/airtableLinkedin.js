const puppeteer = require('puppeteer');
const CREDS = require('./creds');

var Airtable = require('airtable');
var base = new Airtable({apiKey: 'keyrYSQEONtptwMth'}).base('appfcatXnrEsiTmFB');
var browser = "";
var webpage = "";

var allData = new Array();

function data (ID, Location) {
     this.ID = ID;
     this.Location = Location;
 }

run();

async function run () {
  browser = await puppeteer.launch({
    headless: false,
    // slowMo:200
  });
  webpage = await browser.newPage();
  await webpage.goto("https://www.linkedin.com");

  await webpage.setViewport({
    width: 1200,
    height: 800
  });
  const USERNAME_SELECTOR = '#login-email';
  const PASSWORD_SELECTOR = '#login-password';
  const BUTTON_SELECTOR = '#login-submit';

  await webpage.click(USERNAME_SELECTOR);
  await webpage.keyboard.type(CREDS.username);

  await webpage.click(PASSWORD_SELECTOR);
  await webpage.keyboard.type(CREDS.password);
  
  await webpage.click(BUTTON_SELECTOR);
  await webpage.waitForNavigation();
  try{
    base('Data').select({
        maxRecords: 10,
        pageSize:1,
        sort: [{field: "full_name", direction: "asc"}],
        view: "Master"
    }).eachPage(function page(records, fetchNextPage) {
        records.forEach(function(record) {
          (async () => {
            
            await webpage.goto(record.get('linkedin_url'));
            log("fetching: "+record.get('linkedin_url'))
            await new Promise(function(resolve) {
              setTimeout(resolve, 2000);
            });
            // await webpage.waitFor(2 * 1000);
            var location = await webpage.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > h3').textContent);

            await new Promise(function(resolve) {
              setTimeout(resolve, 2000);
            });
            // await webpage.waitFor(2 * 1000);
            var obj = new data(record.getId(),location)
            allData.push(obj)
            log("Pushing data. Very hard:"+record.getId()+record.get('full_name')+location);
            // await browser.close();
            await new Promise(function(resolve) {
              setTimeout(resolve, 2000);
            });
            // await webpage.waitFor(2 * 1000);
         })();
       });
      fetchNextPage();

      }, function done(err) {
          console.log("Done")
      if (err) {
        log("I'm done (w/you)");
        console.error(err); 
        return;
      }

      log("Total record fetched: "+allData.length)
      for (var i = 0; i < allData.length; i++) {
        log(allData[i].Location)
        base('Data').update(allData[i].ID,{
          "location": allData[i].Location,
        },function(err, record) {
            if (err) {
              console.error("I'm printing an error"+ err);
              return;
            }
            else{
              log("Location updated for:"+record.get('full_name'));
            }
          });
      }
    });

  }//try
  catch(error){
    console.log("I'm catch and I'm catching bad boys today");
    console.log(error);
  }
}//run()

//********************************************Log********************************************
async function log(value){
  console.log(value);
}