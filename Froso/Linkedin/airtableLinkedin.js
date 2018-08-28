const puppeteer = require('puppeteer');
const CREDS = require('./creds');

var Airtable = require('airtable');
var base = new Airtable({apiKey: 'keyrYSQEONtptwMth'}).base('appfcatXnrEsiTmFB');

base('scraped_data').select({
    maxRecords: 3,
    view: "Grid view"
}).eachPage(function page(records, fetchNextPage) {
    records.forEach(function(record) {
        console.log('Retrieved', record.get('linkedin_url'));
        (async () => {
          const browser = await puppeteer.launch({
                headless: false
            });
          const page = await browser.newPage();

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
          await page.goto(record.get('linkedin_url'));
          var name = await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent);
          base('scraped_data').update(record.getId(),{
            "phone": name 
        },function(err, record) {
            if (err) {console.error(err);return;}
            console.log("Phone updated for "+record.get('full_name'));
        });
          console.log("The person's who profile you visited is:"+ name);
          await browser.close();
        })();
    });

    fetchNextPage();

}, function done(err) {
    if (err) { console.error(err); return; }
});

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