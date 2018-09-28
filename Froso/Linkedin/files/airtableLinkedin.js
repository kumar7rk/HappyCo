const puppeteer = require('puppeteer');
const CREDS = require('./creds');

var Airtable = require('airtable');
var base = new Airtable({apiKey: 'keyrYSQEONtptwMth'}).base('appfcatXnrEsiTmFB');
var browser = "";
var webpage = "";

(async () => {
  browser = await puppeteer.launch({
    headless: false
  });
  webpage = await browser.newPage();
  await webpage.goto("https://www.linkedin.com");
  const USERNAME_SELECTOR = '#login-email';
  const PASSWORD_SELECTOR = '#login-password';
  const BUTTON_SELECTOR = '#login-submit';

  await webpage.click(USERNAME_SELECTOR);
  await webpage.keyboard.type(CREDS.username);

  await webpage.click(PASSWORD_SELECTOR);
  await webpage.keyboard.type(CREDS.password);
  
  await webpage.click(BUTTON_SELECTOR);
  await webpage.waitForNavigation();
      
base('scraped_data').select({
    maxRecords: 3,
    view: "Grid view"
}).eachPage(function page(records, fetchNextPage) {
    records.forEach(function(record) {
      (async () => {
        await webpage.goto(record.get('linkedin_url'));

          var name = await webpage.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent);
          var location = await webpage.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > h3').textContent);
          base('scraped_data').update(record.getId(),{
            "location": location 
        },function(err, record) {
            if (err) {console.error(err);return;}
            console.log("Location updated for:"+record.get('full_name'));
        });
          // console.log("The person's who profile you visited is:"+ name);
      })();
          // await browser.close();
    });
    fetchNextPage();
}, function done(err) {
    // await browser.close();
    if (err) { 
      console.error("Error for "+record.get('full_name'))
      console.error(err); 
      return; 
    }
});
})();



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