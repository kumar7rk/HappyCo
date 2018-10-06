const puppeteer = require('puppeteer');
const CREDS = require('./creds');

var Airtable = require('airtable');
var base = new Airtable({apiKey: 'keyrYSQEONtptwMth'}).base('appfcatXnrEsiTmFB');
var browser = "";
var webpage = "";

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
        maxRecords: 1000,
        pageSize:1,
        sort: [{field: "full_name", direction: "asc"}],
        view: "Master"
    }).eachPage(function page(records, fetchNextPage) {
        records.forEach(function(record) {
          (async () => {
            log("Opening profile");
            await webpage.goto(record.get('linkedin_url'));
            await webpage.waitFor(20000);

            log("Pausing for 30 seconds");
            await new Promise(function(resolve) { 
              setTimeout(resolve, 30000)
            });
            /*log("Scrolling");
            await webpage.evaluate(_ => {
              window.scrollBy(0, window.innerHeight);
            });
            log("Done scrolling. Pausing for 2 seconds")

            await webpage.waitFor(2 * 1000);*/

            log("Getting name and location")

            var name = "";
            name = await webpage.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div:nth-child(1) > h1').textContent);
            var location = "";
            location = await webpage.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > h3').textContent);
            log("Got the name and location. Pausing for 2 seconds")

            await webpage.waitFor(2 * 1000);
            log("calling update func");
            base('Data').update(record.getId(),{
              "location": location,
              // "new_name": name
            },function(err, record) {
                if (err) {
                  console.error("I'm printing an error");
                  console.error(err);
                  return;
                }
                else{ 
                  log("Location updated for:"+record.get('full_name'));
                }
            });
            log("Completed update. Pausing for 2 seconds")
            await webpage.waitFor(2 * 1000);
            await browser.close();
         })();
       });
      fetchNextPage();
      }, function done(err) {
      if (err) { 
        log("I'm Done (w/you)");
        console.error(err); 
        return; 
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


/*OUTPUT
{
    "id": "recz7ddYYhNCqUL1T",
    "fields": {
        "email": "nostensen@olenproperties.com",
        "full_name": "Natalia Ostensen",
        "last_title": "Owner",
        "last_company": "Olen Properties",
        "salesforce_id": "0034A00002kpY1fQAE",
        "linkedin_url": "https://www.linkedin.com/in/natalia-ostensen-2818a61a"
    },
    "createdTime": "2018-08-23T17:33:15.000Z"
}*/