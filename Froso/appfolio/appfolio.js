  const puppeteer = require('puppeteer');
  const {performance} = require('perf_hooks');


(async () => {
	
	//const browser = await puppeteer.launch({headless: false}); // default is true

  const browser = await puppeteer.launch();
  const page = await browser.newPage();

  var url =  ["https://vequity.appfolio.com/connect/users/sign_in","https://benchmarkmfr.appfolio.com/connect/users/sign_in","https://realprop.appfolio.com/connect/users/sign_in"
             ,"https://ottosen.appfolio.com/connect/users/sign_in","https://sunpropertymanagement.appfolio.com/connect/users/sign_in"];
/*
  for (var i = 0; i < url.length; i++) {
    await page.goto(url[i]);

    const name = await page.evaluate(() => document.querySelector('.footer__info > p').innerText)
    const phone = await page.evaluate(() => document.querySelector('.footer__info > a').innerText)
    const website = await page.evaluate(() => document.querySelector('.footer__info > p > a ').href)
    
    console.log(name+ "\n" + phone+ "\n" + website)   
  }
*/

if(typeof require !== 'undefined') XLSX = require('xlsx');


  var workbook = XLSX.readFile('Appfolio.xlsx');
  var first_sheet_name = workbook.SheetNames[0];
  var worksheet = workbook.Sheets[first_sheet_name];
  
  var address_of_cell = 'A2';
  var desired_cell = worksheet[address_of_cell];
  var desired_value = (desired_cell ? desired_cell.v : undefined);
//  console.log(desired_value)
  
// worksheet['B2'].v = 'manual';
//  XLSX.writeFile(workbook ,'Appfolio2.xlsx')

  var t0 = performance.now();
  
  for (var i = 2; i < 246; i++) {
    address_of_cell = 'A'+i
    var desired_cell = worksheet[address_of_cell];
    var desired_value = (desired_cell ? desired_cell.v : undefined);
    
    await page.goto(desired_value);

    if (await page.evaluate(() => document.querySelector('.footer__info > p')) != null) {
      
      const name = await page.evaluate(() => document.querySelector('.footer__info > p').innerText)
      const phone = await page.evaluate(() => document.querySelector('.footer__info > a').innerText)
      const website = await page.evaluate(() => document.querySelector('.footer__info > p > a ').href)
      var writeCell = 'B'+i
      worksheet[writeCell].v = name;   

      var writeCell1 = 'C'+i
      worksheet[writeCell1].v = website;   

      //console.log(name+ "\n" + phone+ "\n" + website)   
    }
  }
    XLSX.writeFile(workbook ,'Appfolio4.xlsx')

  var t1 = performance.now();
  console.log((t1-t0)/1000+ " seconds");

  /*var range = {s: {c:0, r:1}, e: {c:0, r:245 }};
  var column = 0;
  for (var R = range.s.r; R < range.e.R; ++R) {
    var cellAddress = {c:column, r:R};
    XLSX.utils.encode_cell(cellAddress);
  }*/

    /*var cell_address = 'B2'
    for (var j = 2; j < 12; j++) {
      cell_address = 'B'+j
      var cell_desired = worksheet[cell_address]; 
      XLSX.write(workbook ,name)
    }*/
  

 /* var range = {s: {c:0, r:1}, e: {c:0, r:245 }};
for(var R = range.s.r; R <= range.e.r; ++R) {
  for(var C = range.s.c; C <= range.e.c; ++C) {
    var cell_address = {c:C, r:R};
    /* if an A1-style address is needed, encode the address 
    var cell_ref = XLSX.utils.encode_cell(cell_address);
    console.log(cell_address)
  }
}*/

/*
  await page.goto('https://vequity.appfolio.com/connect/users/sign_in');
  
  const name = await page.evaluate(() => document.querySelector('.footer__info > p').innerText)
  const phone = await page.evaluate(() => document.querySelector('.footer__info > a').innerText)
  const website = await page.evaluate(() => document.querySelector('.footer__info > p > a ').href)
	
	//console.log(name+ "\n" + phone+ "\n" + website)
*/

  await browser.close();
})();