  const puppeteer = require('puppeteer');
  const {performance} = require('perf_hooks');


(async () => {
	
	//const browser = await puppeteer.launch({headless: false}); // default is true

  const browser = await puppeteer.launch();
  const page = await browser.newPage();

if(typeof require !== 'undefined') XLSX = require('xlsx');


  var workbook = XLSX.readFile('Appfolio.xlsx');
  var worksheet = workbook.Sheets[workbook.SheetNames[0]];
  
  var t0 = performance.now();
 /*
  kind of like saying which columns do you want to go through
  range.s.c is first column
  range.e.c is last
  instead of looping through all I'm saying take first column which should be indexed 0 and go until range.s.c which is again the first column
  doing this means I don't have to hardcode the rows value
  All I need to know is which column has the urls
  and also what rows does the urls start from - excluding headers
*/
var range = XLSX.utils.decode_range(worksheet['!ref']); // get the range
  
 for(var R = range.s.r+1; R <= range.e.r; ++R){
      console.log("Row: "+R)
      var cellref = XLSX.utils.encode_cell({c:range.s.c, r:R}); //A1..
      if(!worksheet[cellref]) continue; //[object Object]
      var cell = worksheet[cellref] //[object Object]
      console.log("cell value "+cell.v);
      try{
        await page.goto(cell.v);
      }
      catch(error){
        console.log(error);
      }
     
      if (await page.evaluate(() => document.querySelector('.footer__info > p')) != null) {
      
      const name = await page.evaluate(() => document.querySelector('.footer__info > p').innerText)
      const phone = await page.evaluate(() => document.querySelector('.footer__info > a').innerText)
      const website = await page.evaluate(() => document.querySelector('.footer__info > p > a ').href)
      
      var writeCellName = XLSX.utils.encode_cell({c:range.s.c+1, r:R}); //B1..
      var writeCellWebsite = XLSX.utils.encode_cell({c:range.s.c+2, r:R}); //C1..
      var writeCellUnits = XLSX.utils.encode_cell({c:range.s.c+3, r:R}); //D1..
      
      if (!worksheet[writeCellUnits]) {
        worksheet[writeCellUnits] = {}
      }
      worksheet[writeCellName].v = name;   
      worksheet[writeCellWebsite].v = website;   
      worksheet[writeCellUnits].v = "Test";   

    }
    XLSX.writeFile(workbook ,'Appfolio1.xlsx')
  }

  var t1 = performance.now();
  console.log((t1-t0)/1000+ " seconds");

  await browser.close();
})();