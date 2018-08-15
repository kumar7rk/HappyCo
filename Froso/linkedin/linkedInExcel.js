const puppeteer = require('puppeteer');
const CREDS = require('./creds');

(async () => {
  const browser = await puppeteer.launch({
    headless: false
  });
  /*const browser = await puppeteer.launch({
        ignoreHTTPSErrors: true
    });*/
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
    
if(typeof require !== 'undefined') XLSX = require('xlsx');

  var workbook = XLSX.readFile('LinkedIn.xlsx');
  var first_sheet_name = workbook.SheetNames[0];
  var worksheet = workbook.Sheets[first_sheet_name];
    var address_of_cell = 'A2';
    var desired_cell = worksheet[address_of_cell];
  
  for (var i = 2; i < 3; i++) {
    address_of_cell = 'A'+i;
    desired_cell = worksheet[address_of_cell];
    var desired_value = (desired_cell ? desired_cell.v : undefined);

  await page.goto(desired_value); 
  
  await page.waitFor(2 * 1000);
  await page.evaluate(_ => {
    window.scrollBy(0, window.innerHeight);
  });

  await page.waitFor(2 * 1000);

    var dur = await page.evaluate(() => document.querySelector
      ('div > h4:nth-child(4) > span.pv-entity__bullet-item-v2').textContent)
        var writeCell = 'C'+i
        worksheet[writeCell].v = dur;
      XLSX.writeFile(workbook ,'LinkedIn2.xlsx')
  }
  await browser.close();
})();