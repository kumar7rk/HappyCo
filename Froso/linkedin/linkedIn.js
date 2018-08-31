//#Linkedin.js before condition checking
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
try{  
  for (var i = 2; i < 4; i++) {
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
  var multiPosition = "";
  //contains title
  if (await page.evaluate(() => document.querySelector('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3').textContent)!=null) {
    data = await page.evaluate(() => document.querySelector
    ('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3').textContent)

    multiPosition = data.includes("Title");
  }
  if (!multiPosition) {
    if (await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent) !=null) {
      console.log("one position")
      var name = await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent)
          writeCell = 'B'+i
          worksheet[writeCell].v = name.trim();
      var companyName = await page.evaluate(() => document.querySelector
      ('div > h4 > span:nth-child(2)').textContent)
          writeCell = 'C'+i
          worksheet[writeCell].v = companyName;

      var title = await page.evaluate(() => document.querySelector
      ('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3').textContent)
          writeCell = 'D'+i
          worksheet[writeCell].v = title;

      var currentJobDuration = await page.evaluate(() => document.querySelector
        ('div > h4:nth-child(4) > span.pv-entity__bullet-item-v2').textContent)

      writeCell = 'E'+i
      worksheet[writeCell].v = currentJobDuration;

      var ym = currentJobDuration.split(" ");
      if (ym.length ==2) {
        currentJobDuration = currentJobDuration.replace('mos',"").trim();
        ym = currentJobDuration.split(" ");
        if (ym.length ==1) {
             var num =  parseInt(ym[0]);
             if (num<4) {
              writeCell = 'F'+i
              worksheet[writeCell].v = "Yo"
            }
        }
      }
    }
  }
  else{
    if (await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent) !=null) {
      console.log("muliple positions")
      var name = await page.evaluate(() => document.querySelector
        ('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent)
          writeCell = 'B'+i
          worksheet[writeCell].v = name.trim();
      
      var companyName = await page.evaluate(() => document.querySelector
      ('div > h3 > span:nth-child(2)').textContent)
          writeCell = 'C'+i
          worksheet[writeCell].v = companyName;

      var title = await page.evaluate(() => document.querySelector
      ('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3').textContent)
          title = title.replace('Title','').trim();
          writeCell = 'D'+i
          worksheet[writeCell].v = title;

      var currentJobDuration = await page.evaluate(() => document.querySelector
        ('div > h4:nth-child(3) > span.pv-entity__bullet-item-v2').textContent)

      writeCell = 'E'+i
      worksheet[writeCell].v = currentJobDuration;
      var ym = currentJobDuration.split(" ");
      if (ym.length ==2) {
        currentJobDuration = currentJobDuration.replace('mos',"").trim();
        ym = currentJobDuration.split(" ");
        if (ym.length ==1) {
           var num =  parseInt(ym[0]);
           if (num<4) {
            writeCell = 'F'+i
            worksheet[writeCell].v = "Yo"
          }
        }
      }
    }
  }
  
//Main which is title for rob and company name for Dheer
  
//Rob- 
//Dheer#ember5440 > div > div.pv-entity__company-summary-info > h3 > span:nth-child(2)

        //Rob- #ember4297 > div.pv-entity__summary-info.pv-entity__summary-info--v2 > h3
        // Dheer- #ember4535 > div > div > div.pv-entity__summary-info.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top > h3 > span:nth-child(2)


//#ember4535 > div > div > div.pv-entity__summary-info.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top > h4:nth-child(3)
    


    //Rob- #ember2002 > div.pv-entity__summary-info.pv-entity__summary-info--v2 > h4:nth-child(4) > span.pv-entity__bullet-item-v2
    //Dheer#ember4535 > div > div > div.pv-entity__summary-info.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top > h4:nth-child(3) > span.pv-entity__bullet-item-v2
    
      XLSX.writeFile(workbook ,'LinkedIn2.xlsx')
  }
}
catch(error){
      console.log(error);
    }
  await browser.close();
})();