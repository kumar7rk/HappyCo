const puppeteer = require('puppeteer');
const CREDS = require('./creds');
const player = require('play-sound')(opts = {});
const {performance} = require('perf_hooks');

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

  var workbook = XLSX.readFile('output.xlsx');
  var first_sheet_name = workbook.SheetNames[0];
  var worksheet = workbook.Sheets[first_sheet_name];
  var address_of_cell = 'A2';
  var desired_cell = worksheet[address_of_cell];

  var t0 = performance.now();
try{
  for (var i = 17; i < 18; i++) {
    console.log("Row: "+i);
    address_of_cell = 'B'+i;
    desired_cell = worksheet[address_of_cell];
    var desired_value = (desired_cell ? desired_cell.v : undefined);

    await page.goto(desired_value); 
    await page.waitFor(2 * 1000);
    await page.evaluate(_ => {
      window.scrollBy(0, window.innerHeight);
    });
  
    await page.waitFor(2 * 1000);

    var data = "";  
    var multiPosition = false
    
    //checking if a job('s holder) exists
    if (await page.evaluate(() => document.getElementsByClassName
      ('pv-entity__position-group-pager ember-view'))==null) {
       console.log("No job exist");
        continue;
    }
    //checking if a job has data
    if (await page.evaluateHandle(() => document.getElementsByClassName
      ('pv-entity__position-group-pager ember-view').textContent)==null) {
        console.log("No data for a job");
        continue;
    }    
    //getting first job - all information
    data = await page.evaluateHandle(() => {
      return Array.from(document.getElementsByClassName('pv-entity__position-group-pager ember-view')).map(elem => elem.textContent.trim()).slice(0,1);
    });
    //converting json to string
    var str = JSON.stringify(await data.jsonValue())
    //if text starts with ["Company Name... it indicates the (first) job has multiple positions
    if (str.trim().startsWith('["Company Name')) {
      multiPosition = true;
    }

    //only one position in the first job
    if (!multiPosition) {
      console.log("single position")

      //getting title, company name, current job duration
      if (await page.evaluate(() => document.querySelector('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3')) != null){
          if (await page.evaluate(() => document.querySelector('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3').textContent) !=null){
              var title = await page.evaluate(() => document.querySelector
              ('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3').textContent)
              writeCell = 'J'+i
              if (!worksheet[writeCell]) {
                 worksheet[writeCell] = {}
              }
              worksheet[writeCell].v = title;
          }
        }
      if (await page.evaluate(() => document.querySelector('div > h4 > span:nth-child(2)')) != null){
          if (await page.evaluate(() => document.querySelector('div > h4 > span:nth-child(2)').textContent) !=null){
            var companyName = await page.evaluate(() => document.querySelector
              ('div > h4 > span:nth-child(2)').textContent)
            writeCell = 'L'+i
            if (!worksheet[writeCell]) {
              worksheet[writeCell] = {}
            }
            worksheet[writeCell].v = companyName;
          }
      }      
      if (await page.evaluate(() => document.querySelector('div.pv-entity__summary-info.pv-entity__summary-info--v2 > h4:nth-child(4) > span.pv-entity__bullet-item-v2')) != null){
        if (await page.evaluate(() => document.querySelector('div.pv-entity__summary-info.pv-entity__summary-info--v2 > h4:nth-child(4) > span.pv-entity__bullet-item-v2').textContent) !=null){
            var currentJobDuration = await page.evaluate(() => document.querySelector
              ('div.pv-entity__summary-info.pv-entity__summary-info--v2 > h4:nth-child(4) > span.pv-entity__bullet-item-v2').textContent)
            writeCell = 'N'+i
            if (!worksheet[writeCell]) {
             worksheet[writeCell] = {}
            }
            //
            worksheet[writeCell].v = currentJobDuration;

          //checking if the current job is less than 3 months old
          var ym = currentJobDuration.split(" ");
          if (ym.length ==2) {
            currentJobDuration = currentJobDuration.replace('mos',"").trim();
            ym = currentJobDuration.split(" ");
            if (ym.length ==1) {
               var num =  parseInt(ym[0]);
               if (num<4) {
                writeCell = 'O'+i
                if (!worksheet[writeCell]) {
                 worksheet[writeCell] = {}
                }
                worksheet[writeCell].v = "Yo"
              }
            }
          }
        }
      }  
    }
    //multiple positions
    else{
      console.log("muliple positions")
      if (await page.evaluate(() => document.querySelector('div > div > div.pv-entity__summary-info-v2.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top.mb2 > h3 > span:nth-child(2)')) != null){
        if (await page.evaluate(() => document.querySelector('div > div > div.pv-entity__summary-info-v2.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top.mb2 > h3 > span:nth-child(2)').textContent) !=null){
          var title = await page.evaluate(() => document.querySelector
          ('div > div > div.pv-entity__summary-info-v2.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top.mb2 > h3 > span:nth-child(2)').textContent)
          title = title.replace('Title','').trim();
          writeCell = 'J'+i
          if (!worksheet[writeCell]) {
           worksheet[writeCell] = {}
          }
          worksheet[writeCell].v = title;
        }
      }

      if (await page.evaluate(() => document.querySelector('div > h3 > span:nth-child(2)')) != null){
        if (await page.evaluate(() => document.querySelector('div > h3 > span:nth-child(2)').textContent) !=null){
          var companyName = await page.evaluate(() => document.querySelector
          ('div > h3 > span:nth-child(2)').textContent)
          writeCell = 'L'+i
          if (!worksheet[writeCell]) {
             worksheet[writeCell] = {}
          }
          worksheet[writeCell].v = companyName;
        }
      }

      if (await page.evaluate(() => document.querySelector('div > h4:nth-child(3) > span.pv-entity__bullet-item-v2')) != null){
        if (await page.evaluate(() => document.querySelector('div > h4:nth-child(3) > span.pv-entity__bullet-item-v2').textContent) !=null){
          var currentJobDuration = await page.evaluate(() => document.querySelector
            ('div > h4:nth-child(3) > span.pv-entity__bullet-item-v2').textContent)
//#ember2170 > div > div > div.pv-entity__summary-info-v2.pv-entity__summary-info--v2.pv-entity__summary-info-margin-top.mb2 > h4:nth-child(3) > span.pv-entity__bullet-item-v2
          writeCell = 'N'+i
          if (!worksheet[writeCell]) {
             worksheet[writeCell] = {}
          }
          worksheet[writeCell].v = currentJobDuration;

          var ym = currentJobDuration.split(" ");
          if (ym.length ==2) {
            currentJobDuration = currentJobDuration.replace('mos',"").trim();
            ym = currentJobDuration.split(" ");
            if (ym.length ==1) {
               var num =  parseInt(ym[0]);
               if (num<4) {
                writeCell = 'O'+i
                if (!worksheet[writeCell]) {
                 worksheet[writeCell] = {}
                }
                worksheet[writeCell].v = "Yo"
              }
            }
          }
        }
      }
    }
    var name = ""
    //getting name location phone birthday
    if (await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1')) != null){
      if (await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent) !=null){
        name = await page.evaluate(() => document.querySelector
            ('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent)
          writeCell = 'I'+i
          if (!worksheet[writeCell]) {
            worksheet[writeCell] = {}
          }
          worksheet[writeCell].v = name.trim();
      }
    }
    
    const clickElement = 'span.pv-top-card-v2-section__entity-name.pv-top-card-v2-section__contact-info.ml2'
    if (name !="") {
      await page.click(clickElement)
    
      await page.waitFor(1 * 1000); 

      if (await page.evaluate(() => document.querySelector('div > section.pv-contact-info__contact-type.ci-phone > ul > li')) != null){
       if (await page.evaluate(() => document.querySelector('div > section.pv-contact-info__contact-type.ci-phone > ul > li').textContent) !=null){
        var phone = await page.evaluate(() => document.querySelector('div > section.pv-contact-info__contact-type.ci-phone > ul > li').textContent);
        phone = phone.trim().replace('(Mobile)','').replace('(Home)','').replace('(Work)','').replace(' ','').trim()
        writeCell = 'G'+i
        if (!worksheet[writeCell]) {
            worksheet[writeCell] = {}
          }
          worksheet[writeCell].v = phone;
      }
    }
    if (await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > h3')) != null){
       if (await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > h3').textContent) !=null){
        var location = await page.evaluate(() => document.querySelector('div.pv-top-card-v2-section__info.mr5 > h3').textContent);
        writeCell = 'H'+i
        if (!worksheet[writeCell]) {
            worksheet[writeCell] = {}
          }
          worksheet[writeCell].v = location;
      }
    }
    if (await page.evaluate(() => document.querySelector('div > section.pv-contact-info__contact-type.ci-birthday > div > span')) != null){
     if (await page.evaluate(() => document.querySelector('div > section.pv-contact-info__contact-type.ci-birthday > div > span').textContent) !=null){
        var birthday = await page.evaluate(() => document.querySelector('div > section.pv-contact-info__contact-type.ci-birthday > div > span').textContent);
        writeCell = 'Q'+i
        if (!worksheet[writeCell]) {
            worksheet[writeCell] = {}
        }
        worksheet[writeCell].v = birthday;
      }
    }
  }
  var today = new Date();
  var dd = today.getDate();
  var mm = today.getMonth()+1; //January is 0!
  var yyyy = today.getFullYear();

  if(dd<10) {
      dd = '0'+dd
  } 
  if(mm<10) {
      mm = '0'+mm
  } 
  var lastVisitedDate = dd + '/' + mm + '/' + yyyy;
  
  writeCell = 'P'+i
  if (!worksheet[writeCell]) {
      worksheet[writeCell] = {}
  }
  worksheet[writeCell].v = lastVisitedDate;
  XLSX.writeFile(workbook ,'output.xlsx')
  }
}
catch(error){
      console.log(error);
      player.play('./files/error.mp3', function(err){
      if (err) throw err
    })
    }
  var t1 = performance.now();
  console.log((t1-t0)/1000+ " seconds");

  await browser.close();
  player.play('./files/completed.mp3', function(err){
      if (err) throw err
    })
})();