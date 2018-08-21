  const puppeteer = require('puppeteer');
  const {performance} = require('perf_hooks');


(async () => {
	
	//const browser = await puppeteer.launch({headless: false}); // default is true

  const browser = await puppeteer.launch();
  const page = await browser.newPage();

if(typeof require !== 'undefined') XLSX = require('xlsx');


  var workbook = XLSX.readFile('Appfolio.xlsx');
  var first_sheet_name = workbook.SheetNames[0];
  var worksheet = workbook.Sheets[first_sheet_name];
  
  var address_of_cell = 'A2';
  var desired_cell = worksheet[address_of_cell];
  var url = (desired_cell ? desired_cell.v : undefined);

  var t0 = performance.now();

  for (var i = 2; i < 4; i++) {
    address_of_cell = 'A'+i
    desired_cell = worksheet[address_of_cell];
    url = (desired_cell ? desired_cell.v : undefined);
    
    await page.goto(url);

    if (await page.evaluate(() => document.querySelector('.footer__info > p')) != null) {
      
      const name = await page.evaluate(() => document.querySelector('.footer__info > p').innerText)
      const phone = await page.evaluate(() => document.querySelector('.footer__info > a').innerText)
      const website = await page.evaluate(() => document.querySelector('.footer__info > p > a ').href)
      
      var writeCell = 'B'+i
      worksheet[writeCell].v = name;   

      writeCell = 'C'+i
      worksheet[writeCell].v = website;   

      console.log(name+ "\n" + phone+ "\n" + website+ "\n")   
    }
  }
    XLSX.writeFile(workbook ,'Appfolio1.xlsx')

  var t1 = performance.now();
  // console.log((t1-t0)/1000+ " seconds");

  var range = XLSX.utils.decode_range(worksheet['!ref']); // get the range
  console.log(range)
  for (var R = range.s.r; R < range.e.R; ++R) {
      console.log("Stuff");
    for(var C = range.s.c; C <= range.e.c; ++C) {
      console.log("More Stuff");
    }
  }

  await browser.close();
})();