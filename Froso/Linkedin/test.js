const puppeteer = require('puppeteer');
const CREDS = require('./creds');
const player = require('play-sound')(opts = {});

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
    
if(typeof require !== 'undefined') XLSX = require('xlsx');

  var workbook = XLSX.readFile('LinkedIn.xlsx');
  var first_sheet_name = workbook.SheetNames[0];
  var worksheet = workbook.Sheets[first_sheet_name];
  var address_of_cell = 'A2';
  var desired_cell = worksheet[address_of_cell];

try{  

  for (var i = 3; i < 4; i++) {

    console.log("Row: "+i);
    
    address_of_cell = 'A'+i;
    desired_cell = worksheet[address_of_cell];
    var desired_value = (desired_cell ? desired_cell.v : undefined);

    await page.goto(desired_value); 
    await page.waitFor(2 * 1000);
    await page.evaluate(_ => {
      window.scrollBy(0, window.innerHeight);
    });
  
    await page.waitFor(2 * 1000);

    var data = "";  
    var multiPosition = false;
    data = await page.evaluateHandle(() => {
      return Array.from(document.getElementsByClassName('pv-entity__position-group-pager ember-view')).map(elem => elem.textContent.trim()).slice(0,1);
    });
    var str = JSON.stringify(await data.jsonValue())
    if (str.trim().startsWith('["Company Name')) {
      multiPosition = true;
    } 
    const clickElement = 'span.pv-top-card-v2-section__entity-name.pv-top-card-v2-section__contact-info.ml2'
    await page.waitFor(2 * 1000);
    await page.click(clickElement)

    await page.waitFor(2 * 1000); 
    var phone = await page.evaluate(() => document.querySelector('div > section.pv-contact-info__contact-type.ci-phone > ul > li').textContent);
    console.log(phone.trim().replace('(Mobile)','').replace('(Home)','').replace('(Work)','').replace(' ','').trim())

    console.log("MultiPosition: "+multiPosition)
  }
}
catch(error){
      console.log(error);
      player.play('error.mp3', function(err){
      if (err) throw err
    })
}
  await browser.close();
  player.play('completed.mp3', function(err){
      if (err) throw err
    })
})();