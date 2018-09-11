const puppeteer = require('puppeteer');
const CREDS = require('./creds');
const player = require('play-sound')(opts = {});

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
    
    //#ember2026 > div > div > div.pv-entity__summary-info-v2.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top.mb2 > h3 > span:nth-child(2)
    //what's this checking?
    if (await page.evaluate(() => document.querySelector
      ('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3'))==null) {
        continue;
    }
    //what's this checking?
    if (await page.evaluate(() => document.querySelector
      ('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3').textContent)==null) {
        continue;
    }
    
    //contains title
    data = await page.evaluate(() => document.querySelector
    ('div.pv-entity__summary-info.pv-entity__summary-info--v2').textContent)
    console.log("Title: "+data);
    multiPosition = data.includes("Title");
    //div.pv-entity__summary-info.pv-entity__summary-info--v2 > h3
    
    //#ember2677 > div > div > div.pv-entity__summary-info-v2.pv-entity__summary-info--v2
    //.pv-entity__summary-info-margin-top.mb2 > h3 > span:nth-child(2)
    
    if (!multiPosition) {
      if ((await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1'))
        || await page.evaluate(() => document.querySelector('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3'))
        || await page.evaluate(() => document.querySelector('div > h4:nth-child(4) > span.pv-entity__bullet-item-v2'))) ===null) {
            continue;
      }

      if ((await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent)
        || await page.evaluate(() => document.querySelector('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3').textContent) 
        || await page.evaluate(() => document.querySelector('div > h4:nth-child(4) > span.pv-entity__bullet-item-v2').textContent)) ===null) {
            continue;
      }
      console.log("one position")
//getting name, company name, title, currentJobDuration
      var name = await page.evaluate(() => document.querySelector
        ('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent)
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

//checking if the current job is less than 3 months old
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
    else{
      if (await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1')) ===null
        || await page.evaluate(() => document.querySelector('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3')) ===null
        || await page.evaluate(() => document.querySelector('div > h4:nth-child(3) > span.pv-entity__bullet-item-v2')) ===null) {
        // console.log("LOC: "+113)
        continue;
      }

      if (await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent) ===null
        || await page.evaluate(() => document.querySelector('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3').textContent) ===null
        || await page.evaluate(() => document.querySelector('div > h4:nth-child(3) > span.pv-entity__bullet-item-v2').textContent) ===null) {
        // console.log("LOC: "+120)
        continue;
      }
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
    XLSX.writeFile(workbook ,'LinkedIn2.xlsx')
  }
}
catch(error){
      console.log(error);
      player.play('error.mp3', function(err){
      if (err) throw err
    })
    }
  await browser.close();
  /*player.play('completed.mp3', function(err){
      if (err) throw err
    })*/
})();