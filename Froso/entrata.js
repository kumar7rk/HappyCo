  const puppeteer = require('puppeteer');

(async () => {
  
  //const browser = await puppeteer.launch({headless: false}); // default is true

  const browser = await puppeteer.launch({
        ignoreHTTPSErrors: true
    });
  const page = await browser.newPage();

  var url =  ["https://vequity.appfolio.com/connect/users/sign_in","https://benchmarkmfr.appfolio.com/connect/users/sign_in","https://realprop.appfolio.com/connect/users/sign_in"
             ,"https://ottosen.appfolio.com/connect/users/sign_in","https://sunpropertymanagement.appfolio.com/connect/users/sign_in"];

if(typeof require !== 'undefined') XLSX = require('xlsx');


  var workbook = XLSX.readFile('Entrata.xlsx');
  var first_sheet_name = workbook.SheetNames[0];
  var worksheet = workbook.Sheets[first_sheet_name];
  
  var address_of_cell = 'A113';
  var desired_cell = worksheet[address_of_cell];
  var desired_value = (desired_cell ? desired_cell.v : undefined);

  for (var i = 104; i < 1339; i++) {
    console.log("Row: "+i);
    
    address_of_cell = 'A'+i
    var desired_cell = worksheet[address_of_cell];
    var desired_value = (desired_cell ? desired_cell.v : undefined);
    var name = "";
    var url = "https://"+desired_value
    
    var marketCellAddress = 'C'+i
    var marketCell = worksheet[marketCellAddress];
    var marketValue = (marketCell ? marketCell.v : undefined);
    
    /*if (await page.evaluate(() => document.querySelector('.corporate_logo.corporate-logo-item > a')) != null) {
      name = await page.evaluate(() => document.querySelector('.corporate_logo.corporate-logo-item > a').href)
      if (name = "http://www.trinity-pm.com/") {
        console.log("Trinity: "+i);
        var writeCell = 'B'+i
        worksheet[writeCell].v = "Trinity Property Consultants";
        
        var writeCellMarket = 'C'+i
        worksheet[writeCellMarket].v = "Enterprise";
      }
    }*/
    
    /*if (await page.evaluate(() => document.querySelector('.footer-copyright')) != null) {
      name = await page.evaluate(() => document.querySelector('.footer-copyright').textContent)
      //console.log("Potential Greystar"+ name)
      if (name.includes("Greystar")) {
        console.log("Greystar: "+i);
       
        var writeCell = 'B'+i
        worksheet[writeCell].v = "Greystar";
        
        var writeCellMarket = 'C'+i
        worksheet[writeCellMarket].v = "Enterprise";
      }
    }*/
    
    /*if (await page.evaluate(() => document.querySelector('.corporate_logo_1.corporate-logo-item > a')) != null) {
          name = await page.evaluate(() => document.querySelector('.corporate_logo_1.corporate-logo-item > a').href)
          if (name = "http://www.fairfieldresidential.com/") {
            console.log("Fairfield: "+i);
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Fairfield Residential";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
          }
        }*/
    // console.log(marketValue);
    if (marketValue =="Market type"){
        // console.log(marketValue);
      await page.goto(url); 
      if (await page.evaluate(() => document.querySelector('.corporate-logo-container.corporate-logo > a')) != null) {
        name = await page.evaluate(() => document.querySelector('.corporate-logo-container.corporate-logo > a').href)
        if (name = "http://www.millcreekplaces.com") {
          console.log("Mill creek: "+i);
          var writeCell = 'B'+i
          worksheet[writeCell].v = "Mill Creek";
          
          var writeCellMarket = 'C'+i
          worksheet[writeCellMarket].v = "Enterprise";
        }
      }
      XLSX.writeFile(workbook ,'Entrata2.xlsx')
    }
  }
  await browser.close();
})();